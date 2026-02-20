// SPDX-FileCopyrightText: 2022 Since 2011 Authors of OpenSlides, see https://github.com/OpenSlides/OpenSlides/blob/master/AUTHORS
//
// SPDX-License-Identifier: MIT

package search

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"testing"

	"github.com/OpenSlides/openslides-go/auth"
	"github.com/OpenSlides/openslides-go/datastore/pgtest"
	"github.com/OpenSlides/openslides-search-service/pkg/config"
	"github.com/OpenSlides/openslides-search-service/pkg/meta"
	log "github.com/sirupsen/logrus"
)

const localSearchAddress = "http://localhost:9050/system/search?"

type OutputDataHTMLQuery struct {
	Query      string
	OutputJSON string
}

type OutputDataIndexQuery struct {
	WordQuery     string
	Collections   []string
	OutputAnswers map[string]Answer
}

type testTextIndexController struct {
	*TextIndex
	*pgtest.PostgresTest
	context.Context
}

type mockController struct {
	cfg       *config.Config
	auth      *auth.Auth
	qs        *QueryServer
	reqFields map[string]map[string]*meta.CollectionRelation
	collRel   map[string]map[string]struct{}
}

func TestUnrestrictedOutput(t *testing.T) {
	outputs := []OutputDataIndexQuery{
		{
			"test",
			[]string{},
			map[string]Answer{
				"topic/2": {2.4873344398209953, map[string][]string{
					"_title_original": {"test"},
					"text":            {"test", "west"},
					"title":           {"test"},
				},
				},
				"meeting/2": {0.013346666139263209, map[string][]string{
					"welcome_text": {"text"},
				},
				},
				"meeting/1": {0.013346666139263209, map[string][]string{
					"welcome_text": {"text"},
				},
				},
			},
		},
		{
			"test",
			[]string{"topic", "meeting"},
			map[string]Answer{
				"topic/2": {2.5441687942241002, map[string][]string{
					"_bleve_type":     {"topic"},
					"_title_original": {"test"},
					"text":            {"test", "west"},
					"title":           {"test"},
				},
				},
				"meeting/2": {0.47219033407422906, map[string][]string{
					"_bleve_type":  {"meeting"},
					"welcome_text": {"text"},
				},
				},
				"meeting/1": {0.47219033407422906, map[string][]string{
					"_bleve_type":  {"meeting"},
					"welcome_text": {"text"},
				},
				},
			},
		},
		{
			"test",
			[]string{"topic"},
			map[string]Answer{
				"topic/2": {3.2582204751744155, map[string][]string{
					"_bleve_type":     {"topic"},
					"_title_original": {"test"},
					"text":            {"test", "west"},
					"title":           {"test"},
				},
				},
			},
		},
		{
			"test",
			[]string{"motion"},
			map[string]Answer{},
		},
		{
			"teams",
			[]string{},
			map[string]Answer{
				"topic/2": {0.8773653826510427, map[string][]string{
					"text": {"team"},
				},
				},
			},
		},
	}

	// Setup text index & database
	ctrl, err := initIndex(t)

	if err != nil {
		t.Errorf("Couldn't init index %s", err)
	}
	defer ctrl.PostgresTest.Close()

	t.Run("Check output of unrestricted search queries", func(t *testing.T) {
		for _, output := range outputs {
			answers, err := ctrl.TextIndex.Search(output.WordQuery, output.Collections, 0)

			if err != nil {
				t.Errorf("Error searching in text index: %s", err)
			}

			if !compareAnswers(answers, output.OutputAnswers) {
				t.Errorf("\nOutput of unrestricted text index search should be \n%v\nis\n%v", output.OutputAnswers, answers)
			}
		}
	})

	t.Run("Trying to get info that doesn't exist", func(t *testing.T) {
		answers, err := ctrl.TextIndex.Search("qwertyuiop", []string{""}, 0)

		if err != nil {
			t.Errorf("Error searching in text index: %s", err)
		}

		if !compareAnswers(answers, map[string]Answer{}) {
			t.Errorf("\nOutput of text index search for a non-existing object should be \n%v\nis\n%v", map[string]Answer{}, answers)
		}
	})
}

/* This requires a running autoupdate service
func TestRestrictedOutput(t *testing.T) {
	outputs := []OutputDataHTMLQuery{
		{
			"q=test",
			`{"meeting/1":{"content":{"id":1,"name":"meeting"},"matched_by":{"welcome_text":["text"]}},"meeting/2":{"content":{"id":2,"name":"name"},"matched_by":{"welcome_text":["text"]}}}`,
		},
		{
			"q=test&c=topic,meeting",
			`{"meeting/1":{"content":{"id":1,"name":"meeting"},"matched_by":{"_bleve_type":["meeting"],"welcome_text":["text"]}},"meeting/2":{"content":{"id":2,"name":"name"},"matched_by":{"_bleve_type":["meeting"],"welcome_text":["text"]}}}`,
		},
		{
			"q=test&c=topic",
			`{}`,
		},
		{
			"q=test&c=motion",
			`{}`,
		},
		{
			"q=teams",
			`{}`,
		},
	}

	// Setup text index & database
	ti, err := initIndex(t)

	if err != nil {
		t.Errorf("Couldn't init index %s", err)
	}

	t.Run("Check output of restricted search queries", func(t *testing.T) {
		for _, output := range outputs {
			address := fmt.Sprintf("%s%s", localSearchAddress, output.Query)
			response, err := http.Get(address)
			if err != nil {
				t.Errorf("Couldn't establish connection with Search Service: %s", err)
			}
			defer response.Body.Close()

			byteBody, err := io.ReadAll(response.Body)

			if err != nil {
				t.Errorf("Reading response body: %s", err)
			}

			// Remove score from response
			outputWithoutScore := regexp.MustCompile(`,"score":\d+(\.\d+)?`).ReplaceAllString(string(byteBody), "")

			byteWantedOutput := []byte(output.OutputJSON)

			if !byteEqualityByCharCount([]byte(outputWithoutScore), byteWantedOutput) {
				t.Errorf("\nOutput of restricted query \"%s\" is\n%s\n  should be\n%s", output.Query, outputWithoutScore, byteWantedOutput)
			}
		}
	})
}*/

func TestDatabaseUpdate(t *testing.T) {
	outputBeforeUdpate := OutputDataIndexQuery{

		"test",
		[]string{},
		map[string]Answer{
			"topic/2": {2.4873344398209953, map[string][]string{
				"_title_original": {"test"},
				"text":            {"test", "west"},
				"title":           {"test"},
			},
			},
			"meeting/2": {0.013346666139263209, map[string][]string{
				"welcome_text": {"text"},
			},
			},
			"meeting/1": {0.013346666139263209, map[string][]string{
				"welcome_text": {"text"},
			},
			},
		},
	}

	outputAfterUdpate := OutputDataIndexQuery{
		"test",
		[]string{},
		map[string]Answer{
			"topic/2": {1.8763260236206487, map[string][]string{
				"_title_original": {"test"},
				"text":            {"test", "west"},
				"title":           {"test"},
			},
			},
			"meeting/2": {0.7814626926547352, map[string][]string{
				"welcome_text": {"text", "test"},
			},
			},
			"meeting/1": {0.013398034798872952, map[string][]string{
				"welcome_text": {"text"},
			},
			},
		},
	}

	outputAfterAdd := OutputDataIndexQuery{
		"test",
		[]string{},
		map[string]Answer{
			"topic/2": {2.0287553566700622, map[string][]string{
				"_title_original": {"test"},
				"text":            {"test", "west"},
				"title":           {"test"},
			},
			},
			"topic/3": {0.04828040627900243, map[string][]string{
				"_title_original": {"west"},
				"text":            {"west"},
				"title":           {"west"},
			},
			},
			"meeting/2": {0.8690472848365689, map[string][]string{
				"welcome_text": {"text", "test"},
			},
			},
			"meeting/1": {0.014899656597235321, map[string][]string{
				"welcome_text": {"text"},
			},
			},
		},
	}

	// Setup text index & database
	ctrl, err := initIndex(t)

	if err != nil {
		t.Errorf("Couldn't init index %s", err)
	}
	defer ctrl.PostgresTest.Close()

	// Update Textindex
	ctrl.TextIndex.db.cfg.Index.Age = 0 // Force update
	err = ctrl.TextIndex.update()

	if err != nil {
		t.Errorf("Error updating text index: %s", err)
	}

	t.Run("Check output before updating database", func(t *testing.T) {
		answers, err := ctrl.TextIndex.Search(outputBeforeUdpate.WordQuery, outputBeforeUdpate.Collections, 0)

		if err != nil {
			t.Errorf("Error searching in text index: %s", err)
		}

		if !compareAnswers(answers, outputBeforeUdpate.OutputAnswers) {
			t.Errorf("\nOutput of unrestricted text index search should be \n%v\nis\n%v", outputBeforeUdpate.OutputAnswers, answers)
		}

	})

	// Update database - change searchable object
	err = pgConnCommand(t, ctrl.PostgresTest, ctrl.Context, "UPDATE meeting_t SET welcome_text = 'text test' WHERE id = 2", true)

	if err != nil {
		t.Errorf("Error updating postgres database: %s", err)
	}

	// Update Textindex
	err = ctrl.TextIndex.update()

	if err != nil {
		t.Errorf("Error updating text index: %s", err)
	}

	t.Run("Check output after updating database", func(t *testing.T) {
		answers, err := ctrl.TextIndex.Search(outputAfterUdpate.WordQuery, outputAfterUdpate.Collections, 0)

		if err != nil {
			t.Errorf("Error searching in text index: %s", err)
		}

		if !compareAnswers(answers, outputAfterUdpate.OutputAnswers) {
			t.Errorf("\nOutput of unrestricted text index search should be \n%v\nis\n%v", outputAfterUdpate.OutputAnswers, answers)
		}
	})

	// Update database - add new searchable object
	err = pgConnCommand(t, ctrl.PostgresTest, ctrl.Context, `
INSERT INTO topic_t (id, title, text, sequential_number, meeting_id)
VALUES (3, 'West Land', 'A western', 2, 2);
SELECT nextval('topic_t_id_seq');
INSERT INTO list_of_speakers_t (
    id, content_object_id, sequential_number, meeting_id
)
VALUES (3, 'topic/3', 2, 2);
INSERT INTO agenda_item_t (content_object_id, meeting_id)
VALUES ('topic/3', 2);`, true)

	if err != nil {
		t.Errorf("Error updating postgres database: %s", err)
	}

	// Update Textindex
	err = ctrl.TextIndex.update()

	if err != nil {
		t.Errorf("Error updating text index: %s", err)
	}

	t.Run("Check output after updating database", func(t *testing.T) {
		answers, err := ctrl.TextIndex.Search(outputAfterAdd.WordQuery, outputAfterAdd.Collections, 0)

		if err != nil {
			t.Errorf("Error searching in text index: %s", err)
		}

		if !compareAnswers(answers, outputAfterAdd.OutputAnswers) {
			t.Errorf("\nOutput of text index search after newly added object to database should be \n%v\nis\n%v", outputAfterAdd.OutputAnswers, answers)
		}
	})

	// Update database - delete searchable object
	err = pgConnCommand(t, ctrl.PostgresTest, ctrl.Context, `
DELETE FROM agenda_item_t WHERE content_object_id = 'topic/3';
DELETE FROM list_of_speakers_t WHERE id = 3;
DELETE FROM topic_t WHERE id = 3;
		`, true)

	if err != nil {
		t.Errorf("Error updating postgres database: %s", err)
	}

	// Update Textindex
	err = ctrl.TextIndex.update()

	if err != nil {
		t.Errorf("Error updating text index: %s", err)
	}

	t.Run("Check output after added object has been deleted again from database", func(t *testing.T) {
		answers, err := ctrl.TextIndex.Search(outputAfterUdpate.WordQuery, outputAfterUdpate.Collections, 0)

		if err != nil {
			t.Errorf("Error searching in text index: %s", err)
		}

		if !compareAnswers(answers, outputAfterUdpate.OutputAnswers) {
			t.Errorf("\nOutput of text index search after deleting newly added object should be \n%v\nis\n%v", outputAfterUdpate.OutputAnswers, answers)
		}
	})
}

func initIndex(t *testing.T) (*testTextIndexController, error) {
	err := os.Setenv("RESTRICTER_URL", "...")

	if err != nil {
		t.Errorf("Error setting environment variable: %s", err)
		return nil, err
	}

	cfg, _ := config.GetConfig()

	ctx := t.Context()

	models, err := meta.Fetch[meta.Collections]("../../meta/models.yml")
	if err != nil {
		t.Errorf("loading models failed: %s", err)
		return nil, err
	}

	// For text indexing we can only use string fields.
	searchModels := models.Clone()

	// If there are search filters configured cut search models further down.
	if cfg.Models.Search != "" {
		searchFilter, err := meta.Fetch[meta.Filters]("../../meta/search.yml")
		if err != nil {
			t.Errorf("loading search filters failed. %s", err)
			return nil, err
		}
		searchModels.Retain(searchFilter.Retain(false))
	} else {
		searchModels.Retain(meta.RetainStrings())
	}

	// Create test postgres container
	closePG := true
	pg, err := pgtest.NewPostgresTest(ctx)
	if err != nil {
		t.Errorf("Error starting postgres: %s", err)
		return nil, err
	}
	defer func() {
		if closePG {
			pg.Close()
		}
	}()

	// Alter cfg to refer to test postgres container
	cfg.Database.User = pg.Env["DATABASE_USER"]
	cfg.Database.Database = pg.Env["DATABASE_NAME"]
	cfg.Database.Host = pg.Env["DATABASE_HOST"]
	cfg.Database.Port, err = strconv.Atoi(pg.Env["DATABASE_PORT"])

	if err != nil {
		t.Errorf("converting test postgres post to int: %s", err)
		return nil, err
	}

	// Add mock data
	sqlFromFile(t, ctx, pg, "../../meta/dev/sql/test_data.sql")
	sqlFromFile(t, ctx, pg, "../../dev/mock_data.sql")

	// Create database and text index
	db := NewDatabase(cfg)
	ti, err := NewTextIndex(cfg, db, searchModels)
	if err != nil {
		t.Errorf("creating text index failed: %s", err)
		return nil, err
	}

	closePG = false
	return &testTextIndexController{ti, pg, ctx}, nil
}

func sqlFromFile(t *testing.T, ctx context.Context, pg *pgtest.PostgresTest, path string) error {

	// Read sql content
	file, err := os.ReadFile(path)
	if err != nil {
		t.Errorf("reading sql file for path %s: %s", path, err)
		return err
	}

	conn, err := pg.Conn(ctx)

	if err != nil {
		t.Errorf("getting pgx connection for path %s: %s", path, err)
		return err
	}

	_, err = conn.Begin(ctx)
	if err != nil {
		t.Errorf("starting pgx connection for path %s: %s", path, err)
		return err
	}
	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, string(file))

	if err != nil {
		t.Errorf("adding mock data for path %s: %s", path, err)
		return err
	}

	return nil
}

func pgConnCommand(t *testing.T, pg *pgtest.PostgresTest, ctx context.Context, cmd string, execMode bool) error {
	conn, err := pg.Conn(ctx)

	if err != nil {
		t.Errorf("getting pgx connection for command %s: %s", cmd, err)
		return err
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		t.Errorf("starting pgx connection for command %s: %s", cmd, err)
		return err
	}
	defer conn.Close(ctx)

	if execMode {
		_, err := tx.Exec(ctx, cmd)

		if err != nil {
			t.Errorf("adding mock data for command %s: %s", cmd, err)
			return err
		}

		err = tx.Commit(ctx)

		if err != nil {
			t.Errorf("adding mock data for command %s: %s", cmd, err)
			return err
		}
	} else {
		// Query command for debugging purposes
		rows, err := tx.Query(ctx, cmd)

		if err != nil {
			t.Errorf("adding mock data for command %s: %s", cmd, err)
			return err
		}
		defer rows.Close()

		// Get column names of table
		descriptions := rows.FieldDescriptions()
		columns := make([]string, len(descriptions))

		for i, description := range descriptions {
			columns[i] = description.Name
		}

		var ticker int
		for rows.Next() {
			values, err := rows.Values()
			if err != nil {
				return err
			}

			// Assign data
			data := make(map[string]any, len(values))
			log.Infof("New Row %d", ticker)
			ticker += 1

			for i, v := range values {
				data[columns[i]] = v
				log.Infof("Column #%d - %s has value %v", i, columns[i], data[columns[i]])
			}
		}
	}
	return nil
}

func sortByteArray(a []byte) []byte {
	newSlice := make([]byte, len(a))
	copy(newSlice, a)
	sort.Slice(newSlice, func(i, j int) bool {
		return newSlice[i] < newSlice[j]
	})

	bytes.ReplaceAll(newSlice, []byte{'\n'}, []byte{})
	return newSlice
}

func byteEqualityByCharCount(a, b []byte) bool {
	return reflect.DeepEqual(sortByteArray(a), sortByteArray(b))
}

// Can not use reflect.DeepEqual, since maps and arrays may have different orderings of elements
// Instead, convert both Answer-maps to byte arrays, sort those arrays and check for equality on sorted byte arrays
func compareAnswers(a, b map[string]Answer) bool {
	byteA := convertAnswerMapToByteArray(a)
	byteB := convertAnswerMapToByteArray(b)

	return byteEqualityByCharCount(byteA, byteB)
}

func convertAnswerMapToByteArray(a map[string]Answer) []byte {
	var byteA []byte

	for _, answer := range a {
		for _, words := range answer.MatchedWords {
			for _, word := range words {
				byteA = append(byteA, []byte(word)...)
			}
		}
		// Score seems to arbitrarily change between databank resets, so it's taken out of the comparison function for now
		// byteA = append(byteA, []byte(fmt.Sprint(answer.Score))...)
	}
	return byteA
}

func debugPrintByteArrayAsInt(t *testing.T, a []byte) {
	var s string
	for _, x := range a {
		s += fmt.Sprint(int(x))
	}
	t.Log(s)
}

func debugFindLocationOfDifference(a, b []byte) {
	log.Info("Comparing")
	log.Info(string(a))
	log.Info(string(b))
	for i, v := range a {
		if len(b) <= i {
			log.Infof("B is too short, missing %d at the end", v)
		} else if b[i] != a[i] {
			log.Infof("Pos %d character is %d for a and %d for b", i, v, b[i])
		}
	}
}
