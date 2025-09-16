// SPDX-FileCopyrightText: 2022 Since 2011 Authors of OpenSlides, see https://github.com/OpenSlides/OpenSlides/blob/master/AUTHORS
//
// SPDX-License-Identifier: MIT

package search

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/OpenSlides/openslides-search-service/pkg/config"
	"github.com/jackc/pgx/v5"
)

const (
	selectCollectionSizesSQL = `
SELECT
	n_live_tup,relname
FROM
	pg_stat_user_tables
ORDER BY n_live_tup DESC
`

	selectAllTableNames = `
SELECT
	tablename
FROM
	pg_tables
WHERE
	schemaname = 'public'
`

	selectTableContentTemplate = `
SELECT
	*
FROM
	$1
`

	selectLatestUpdate = `
SELECT
	timestamp
FROM
	os_notify_log_t
WHERE
	fqid = '$1'
ORDER BY xact_id DESC
LIMIT 1
	`

	/* TODO selectDiffSQL = `
	SELECT
	  fqid,
	  CASE WHEN updated > $1 THEN data::text ELSE NULL END,
	  updated
	FROM models
	WHERE NOT deleted`*/
)

type entry struct {
	updated time.Time
	gen     uint16
}

// Database manages the updates needed to drive the text index.
type Database struct {
	cfg         *config.Config
	last        time.Time
	gen         uint16
	collections map[string]map[int32]*entry
}

// NewDatabase creates a new database,
func NewDatabase(cfg *config.Config) *Database {
	return &Database{
		cfg: cfg,
	}
}

func (db *Database) run(fn func(context.Context, *pgx.Conn) error) error {
	ctx := context.Background()
	config, err := pgx.ParseConfig(db.cfg.Database.ConnectionConfig())
	if err != nil {
		return err
	}

	// Simple protocol is used for PGBouncer compatibility
	config.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	con, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		return err
	}
	defer con.Close(ctx)
	return fn(ctx, con)
}

func (db *Database) numEntries() int {
	if db.collections == nil {
		return 0
	}
	var sum int
	for _, col := range db.collections {
		sum += len(col)
	}
	return sum
}

func splitFqid(fqid string) (string, int, error) {
	col, idS, ok := strings.Cut(fqid, "/")
	if !ok {
		return "", 0, fmt.Errorf("invalid fqid: %q", fqid)
	}
	id, err := strconv.Atoi(idS)
	if err != nil {
		return "", 0, fmt.Errorf("invalid fqid: %q: %v", fqid, err)
	}
	return col, id, nil
}

type updateEventType int

const (
	addedEvent updateEventType = iota
	changedEvent
	removeEvent
)

type eventHandler func(evtType updateEventType, collection string, id int32, data map[string]any) error

func nullEventHandler(updateEventType, string, int32, map[string]any) error { return nil }

func (db *Database) update(handler eventHandler) error {
	start := time.Now()

	// Do not update if it is young enough.
	if !db.last.IsZero() && !start.After(db.last.Add(db.cfg.Index.Age)) {
		return nil
	}

	if handler == nil {
		handler = nullEventHandler
	}

	defer func() {
		log.Debugf("updating database took %v\n", time.Since(start))
	}()

	return db.run(func(ctx context.Context, conn *pgx.Conn) error {

		queryMap, err := db.generateTableQueryMap(ctx, conn)
		if err != nil {
			return err
		}

		before := db.numEntries()
		var unchanged, added, entries int

		ngen := db.gen + 1 // may overflow but thats okay.

		for tablename, query := range queryMap {

			rows, err := conn.Query(ctx, query)
			if err != nil {
				return err
			}
			defer rows.Close()

			// Get column names of table
			descriptions := rows.FieldDescriptions()
			columns := make([]string, len(descriptions))

			for i, description := range descriptions {
				columns[i] = description.Name
			}

			for rows.Next() {
				values, err := rows.Values()
				if err != nil {
					return err
				}

				// Assign data
				data := make(map[string]any, len(values))
				var id int32
				id = -1

				log.Info(len(values))

				for i, v := range values {
					log.Info(columns[i])
					if columns[i] == "id" {
						id = v.(int32)
						continue
					}
					data[columns[i]] = v
				}

				if id == -1 {
					// Discard this table
					log.Info(tablename + " discarded, for there is no id column found")
					break
				}

				entries++

				// handle changed and new
				collection := db.collections[tablename]
				if collection == nil {
					collection = make(map[int32]*entry)
					db.collections[tablename] = collection
				}
				e := collection[id]
				if e == nil {
					if err := handler(addedEvent, tablename, id, data); err != nil {
						return err
					}
					collection[id] = &entry{
						// TODO: updated: updated,
						gen: ngen,
					}
					added++
				} else {
					// TODO: e.updated = updated
					e.gen = ngen
					// TODO: data is always non-nil, because we are no longer using the diff-SQL which would've returned nil for data if it shouldn't be updated
					if data != nil {
						if err := handler(changedEvent, tablename, id, data); err != nil {
							return err
						}
					} else {
						unchanged++
					}
				}
			}
			if err := rows.Err(); err != nil {
				return err
			}
		}
		// TODO: Do some clever arithmetics based on
		// before, entries, added and unchanged to
		// early stop this.
		var removed int
		if unchanged != before {
			for k, col := range db.collections {
				for id, e := range col {
					if e.gen != ngen {
						removed++
						delete(col, id)
						if err := handler(removeEvent, k, id, nil); err != nil {
							return err
						}
					}
				}
			}
		}

		log.Debugf("entries: %d / before: %d\n",
			entries, before)
		log.Debugf("unchanged: %d / added: %d / removed: %d\n",
			unchanged, added, removed)

		db.last = start
		db.gen = ngen
		return nil
	})
}

func preAllocCollections(ctx context.Context, conn *pgx.Conn) (map[string]map[int32]*entry, error) {
	cols := make(map[string]map[int32]*entry)
	rows, err := conn.Query(ctx, selectCollectionSizesSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var size int
		var col string
		if err := rows.Scan(&size, &col); err != nil {
			return nil, err
		}
		cols[col] = make(map[int32]*entry, size)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return cols, nil
}

func (db *Database) generateTableQueryMap(ctx context.Context, conn *pgx.Conn) (map[string]string, error) {
	// Get all tablenames
	tablenames, err := conn.Query(ctx, selectAllTableNames)
	if err != nil {
		return nil, err
	}
	defer tablenames.Close()

	// Construct Query Map
	queryMap := make(map[string]string)

	for tablenames.Next() {
		// Scan tablename
		var tablename string
		tablenames.Scan(&tablename)

		// Create SQL Query
		constructedSQLStatement := strings.Replace(selectTableContentTemplate, "$1", tablename, -1)

		// Add to queryMap
		queryMap[tablename] = constructedSQLStatement
	}

	return queryMap, nil
}

func (db *Database) fill(handler eventHandler) error {
	start := time.Now()
	defer func() {
		log.Infof("initial database fill took %v\n", time.Since(start))
	}()

	if handler == nil {
		handler = nullEventHandler
	}

	return db.run(func(ctx context.Context, conn *pgx.Conn) error {
		cols, err := preAllocCollections(ctx, conn)
		if err != nil {
			return err
		}

		queryMap, err := db.generateTableQueryMap(ctx, conn)
		if err != nil {
			return err
		}

		for tablename, query := range queryMap {
			var numEntries, size int

			rows, err := conn.Query(ctx, query)
			if err != nil {
				return err
			}
			defer rows.Close()

			// Get column names of table
			descriptions := rows.FieldDescriptions()
			columns := make([]string, len(descriptions))

			for i, description := range descriptions {
				columns[i] = description.Name
			}

			for rows.Next() {
				values, err := rows.Values()
				if err != nil {
					return err
				}
				// Assign data
				data := make(map[string]any, len(values))
				var id int32
				id = -1

				for i, v := range values {
					if columns[i] == "id" {
						id = v.(int32)
						continue
					}
					data[columns[i]] = v
				}
				log.Infof("Fill %s: %d datapoints for id %d ", tablename, len(data), id)

				if id == -1 {
					// Discard this table
					log.Info(tablename + " discarded, for there is no id column found")
					break
				}

				collection := cols[tablename]
				if collection == nil {
					log.Warnf("alloc collection %q. This should has happend before.\n", tablename)
					collection = make(map[int32]*entry)
					cols[tablename] = collection
				}
				if err := handler(addedEvent, tablename, id, data); err != nil {
					return err
				}

				size += len(data)

				/*collection[id] = &entry{
					updated: updated,
				}*/

				numEntries++

				if err := rows.Err(); err != nil {
					return err
				}

				log.Debugf("num entries: %d / size: %d (%.2fMiB)\n",
					numEntries,
					size, float64(size)/(1024*1024))
			}
		}

		log.Debugf("num collections: %d\n", len(cols))
		db.collections = cols
		db.last = start
		return nil
	})
}
