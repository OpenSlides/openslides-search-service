// SPDX-FileCopyrightText: 2022 Since 2011 Authors of OpenSlides, see https://github.com/OpenSlides/OpenSlides/blob/master/AUTHORS
//
// SPDX-License-Identifier: MIT

package search

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"sort"
	"testing"
)

// For data in mock_data.sql and request string "q=test", response is
// {"meeting/1":{"Score":0.013346666139263209,"MatchedWords":{"welcome_text":["text"]}},"meeting/2":{"Score":0.013346666139263209,"MatchedWords":{"welcome_text":["text"]}},"topic/2":{"Score":2.4873344398209953,"MatchedWords":{"_title_original":["test"],"text":["test","west"],"title":["test"]}}}

// For data in mock_data.sql and request string "q=test&c=topic,meeting", response  is
// {"meeting/1":{"Score":0.045900890677894324,"MatchedWords":{"_bleve_type":["meeting"],"welcome_text":["text"]}},"meeting/2":{"Score":0.045900890677894324,"MatchedWords":{"_bleve_type":["meeting"],"welcome_text":["text"]}},"topic/2":{"Score":1.1264835358858345,"MatchedWords":{"_bleve_type":["topic"],"_title_original":["test"],"text":["test","west"],"title":["test"]}}}

// For data in mock_data.sql and request string "q=test&c=topic", response is
// {"topic/2":{"Score":1.327796051982089,"MatchedWords":{"_bleve_type":["topic"],"_title_original":["test"],"text":["test","west"],"title":["test"]}}}

// For data in mock_data.sql and request string "q=test&c=motion", response is
// {}

// For data in mock_data.sql and request string "q=teams", response is
// {"topic/2":{"Score":0.8773653826510427,"MatchedWords":{"text":["team"]}}}

const localSearchAddress = "http://localhost:9050/system/search?"

type OutputData struct {
	Query      string
	OutputJSON string
}

func TestUnrestrictedOutput(t *testing.T) {
	outputs := []OutputData{
		{
			"q=test",
			`{"meeting/1":{"Score":0.013346666139263209,"MatchedWords":{"welcome_text":["text"]}},"meeting/2":{"Score":0.013346666139263209,"MatchedWords":{"welcome_text":["text"]}},"topic/2":{"Score":2.4873344398209953,"MatchedWords":{"_title_original":["test"],"text":["test","west"],"title":["test"]}}}`,
		},
		{
			"q=test&c=topic,meeting",
			`{"meeting/1":{"Score":0.045900890677894324,"MatchedWords":{"_bleve_type":["meeting"],"welcome_text":["text"]}},"meeting/2":{"Score":0.045900890677894324,"MatchedWords":{"_bleve_type":["meeting"],"welcome_text":["text"]}},"topic/2":{"Score":1.1264835358858345,"MatchedWords":{"_bleve_type":["topic"],"_title_original":["test"],"text":["test","west"],"title":["test"]}}}`,
		},
		{
			"q=test&c=topic",
			`{"topic/2":{"Score":1.327796051982089,"MatchedWords":{"_bleve_type":["topic"],"_title_original":["test"],"text":["test","west"],"title":["test"]}}}`,
		},
		{
			"q=test&c=motion",
			`{}`,
		},
		{
			"q=teams",
			`{"topic/2":{"Score":0.8773653826510427,"MatchedWords":{"text":["team"]}}}`,
		},
	}

	t.Run("Check output of unrestricted search queries", func(t *testing.T) {
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

			byteWantedOutput := []byte(output.OutputJSON + "\n") // Line feed character necessary for deep equal

			if !byteEqualityByCharCount(byteBody, byteWantedOutput) {
				t.Errorf("\nOutput of unrestricted query \"%s\" is\n%s\n  should be\n%s", output.Query, byteBody, byteWantedOutput)
			}
		}
	})
}

func TestRestrictedOutput(t *testing.T) {
	outputs := []OutputData{
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
			}
			defer response.Body.Close()

			byteBody, err := io.ReadAll(response.Body)

			if err != nil {
				t.Errorf("Reading response body: %s", err)
			}

			byteWantedOutput := []byte(output.OutputJSON)

			if !byteEqualityByCharCount(byteBody, byteWantedOutput) {
				t.Errorf("\nOutput of reestricted query \"%s\" is\n%s\n  should be\n%s", output.Query, byteBody, byteWantedOutput)
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
