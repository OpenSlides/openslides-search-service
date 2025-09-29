// SPDX-FileCopyrightText: 2022 Since 2011 Authors of OpenSlides, see https://github.com/OpenSlides/OpenSlides/blob/master/AUTHORS
//
// SPDX-License-Identifier: MIT

package search

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
