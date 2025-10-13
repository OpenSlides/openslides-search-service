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

	selectLatestUpdates = `
SELECT DISTINCT ON (fqid)
	fqid,
	operation
FROM
	os_notify_log_t
WHERE
	timestamp >= $1
ORDER BY fqid, timestamp DESC
	`

	selectElementFromTableTemplate = `
SELECT
	*
FROM
	$1
WHERE
	id = $2
`
)

type updateOperation struct {
	id        int
	operation string
}

// Database manages the updates needed to drive the text index.
type Database struct {
	cfg  *config.Config
	last time.Time
	gen  uint16
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

type eventHandler func(evtType updateEventType, collection string, id int, data map[string]any) error

func nullEventHandler(updateEventType, string, int, map[string]any) error { return nil }

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

		updateLogs, err := conn.Query(ctx, selectLatestUpdates, db.last)
		if err != nil {
			return err
		}
		defer updateLogs.Close()

		var added, removed, entries int

		ngen := db.gen + 1 // may overflow but thats okay.

		changeMap := make(map[string][]updateOperation)

		// Convert each fqid to a tablename - id mapping
		for updateLogs.Next() {
			var fqid string
			var operation string
			err = updateLogs.Scan(&fqid, &operation)
			if err != nil {
				return err
			}

			tableName, id, err := splitFqid(fqid)

			if err != nil {
				return err
			}

			if changeMap[tableName] == nil {
				changeMap[tableName] = []updateOperation{}
			}

			changeMap[tableName] = append(changeMap[tableName], updateOperation{
				id,
				operation,
			})
		}

		updateLogs.Close()

		// For each tablename - id mapping, query the corresponding row for data and insert it update the search index
		for tablename, updateOperations := range changeMap {
			for _, updateOperation := range updateOperations {
				// Create SQL Query
				constructedSQLStatement := strings.Replace(selectElementFromTableTemplate, "$1", tablename+"_t", -1)
				constructedSQLStatement = strings.Replace(constructedSQLStatement, "$2", fmt.Sprint(updateOperation.id), -1)

				// log.Infof("A change has been registered: %s with id %d and operation %s", tablename+"_t", updateOperation.id, updateOperation.operation)

				// Skip data collection if it was a delete event
				if updateOperation.operation == "delete" {
					removed++
					if err := handler(removeEvent, tablename, updateOperation.id, nil); err != nil {
						return err
					}
					continue
				}

				// Query
				rows, err := conn.Query(ctx, constructedSQLStatement)
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

					if id == -1 {
						// Discard this table
						log.Info(tablename + " discarded, for there is no id column found")
						continue
					}

					entries++

					// Act based on operation
					switch updateOperation.operation {
					case "insert":
						if err := handler(addedEvent, tablename, updateOperation.id, data); err != nil {
							return err
						}
						added++
					case "update":

						if err := handler(changedEvent, tablename, updateOperation.id, data); err != nil {
							return err
						}
					}
				}
				if err := rows.Err(); err != nil {
					return err
				}
			}
		}

		log.Debugf("entries: %d",
			entries)
		log.Debugf("added: %d / removed: %d\n",
			added, removed)

		db.last = start
		db.gen = ngen
		return nil
	})
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
		queryMap, err := db.generateTableQueryMap(ctx, conn)
		if err != nil {
			return err
		}

		for tablename, query := range queryMap {
			var numEntries, size int

			// Alter tablename to conform meta models
			tablename := strings.TrimSuffix(tablename, "_t")

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
				log.Debugf("Fill %s: %d datapoints for id %d ", tablename, len(data), id)

				if id == -1 {
					// Discard this table
					log.Info(tablename + " discarded, for there is no id column found")
					continue
				}

				// Handle Data
				if err := handler(addedEvent, tablename, int(id), data); err != nil {
					return err
				}

				size += len(data)

				numEntries++

				if err := rows.Err(); err != nil {
					return err
				}

				log.Debugf("num entries: %d / size: %d (%.2fMiB)\n",
					numEntries,
					size, float64(size)/(1024*1024))
			}
		}

		db.last = start
		return nil
	})
}
