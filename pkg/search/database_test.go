// SPDX-FileCopyrightText: 2022 Since 2011 Authors of OpenSlides, see https://github.com/OpenSlides/OpenSlides/blob/master/AUTHORS
//
// SPDX-License-Identifier: MIT

package search

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"sort"
	"testing"

	"github.com/OpenSlides/openslides-go/auth"
	"github.com/OpenSlides/openslides-search-service/pkg/config"
	"github.com/OpenSlides/openslides-search-service/pkg/meta"
	"golang.org/x/sys/unix"
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

type mockController struct {
	cfg       *config.Config
	auth      *auth.Auth
	qs        *QueryServer
	reqFields map[string]map[string]*meta.CollectionRelation
	collRel   map[string]map[string]struct{}
}

func signalContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, unix.SIGTERM)
		<-sig
		cancel()
		<-sig
		os.Exit(2)
	}()
	return ctx, cancel
}

func initIndex() (*TextIndex, error) {
	err := os.Setenv("RESTRICTER_URL", "...")

	if err != nil {
		fmt.Println("Error setting environment variable:", err)
	}

	cfg, _ := config.GetConfig()

	_, cancel := signalContext()
	defer cancel()

	models, err := meta.Fetch[meta.Collections]("../../meta/models.yml")
	if err != nil {
		return nil, fmt.Errorf("loading models failed: %w", err)
	}

	// For text indexing we can only use string fields.
	searchModels := models.Clone()

	// If there are search filters configured cut search models further down.
	if cfg.Models.Search != "" {
		searchFilter, err := meta.Fetch[meta.Filters]("../../meta/search.yml")
		if err != nil {
			return nil, fmt.Errorf("loading search filters failed. %w", err)
		}
		searchModels.Retain(searchFilter.Retain(false))
	} else {
		searchModels.Retain(meta.RetainStrings())
	}

	db := NewDatabase(cfg)
	ti, err := NewTextIndex(cfg, db, searchModels)
	if err != nil {
		return nil, fmt.Errorf("creating text index failed: %w", err)
	}

	return ti, nil
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

	ti, err := initIndex()

	if err != nil {
		t.Errorf("Couldn't init index %s", err)
	}

	if err != nil {
		t.Errorf("Error in search index %s", err)
	}

	t.Run("Check output of unrestricted search queries", func(t *testing.T) {
		for _, output := range outputs {
			answers, err := ti.Search(output.WordQuery, output.Collections, 0)

			if err != nil {
				t.Errorf("Error searching in text index: %s", err)
			}

			if !compareAnswers(answers, output.OutputAnswers) {
				t.Errorf("\nOutput of unrestricted text index search should be \n%v\nis\n%v", output.OutputAnswers, answers)
			}
		}
	})
}

func TestRestrictedOutput(t *testing.T) {
	outputs := []OutputDataHTMLQuery{
		{
			"q=test",
			`{"meeting/1":{"content":{"id":1,"name":"meeting"},"matched_by":{"welcome_text":["text"]},"score":0.013346666139263209},"meeting/2":{"content":{"id":2,"name":"name"},"matched_by":{"welcome_text":["text"]},"score":0.013346666139263209}}`,
		},
		{
			"q=test&c=topic,meeting",
			`{"meeting/1":{"content":{"id":1,"name":"meeting"},"matched_by":{"_bleve_type":["meeting"],"welcome_text":["text"]},"score":0.045900890677894324},"meeting/2":{"content":{"id":2,"name":"name"},"matched_by":{"_bleve_type":["meeting"],"welcome_text":["text"]},"score":0.045900890677894324}}`,
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

	t.Run("Check output of reestricted search queries", func(t *testing.T) {
		for _, output := range outputs {
			address := fmt.Sprintf("%s%s", localSearchAddress, output.Query)
			response, err := http.Get(address)
			if err != nil {
				t.Errorf("Couldn't establish connection with Search Service: %s", err)
				response.Body.Close()
				continue
			}
			defer response.Body.Close()

			byteBody, err := io.ReadAll(response.Body)

			if err != nil {
				t.Errorf("Reading response body: %s", err)
			}

			byteWantedOutput := []byte(output.OutputJSON)

			if !byteEqualityByCharCount(byteBody, byteWantedOutput) {
				t.Errorf("\nOutput of restricted query \"%s\" is\n%s\n  should be\n%s", output.Query, byteBody, byteWantedOutput)
			}
		}
	})
}

func debugPrintByteArrayAsInt(t *testing.T, a []byte) {
	var s string
	for _, x := range a {
		s += fmt.Sprint(int(x))
	}
	t.Log(s)
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
		byteA = append(byteA, []byte(fmt.Sprint(answer.Score))...)
	}
	return byteA
}
