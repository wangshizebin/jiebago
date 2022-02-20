// Copyright 2022 Ze-Bin Wang.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"C"
	"encoding/json"

	"github.com/wangshizebin/jiebago"
	"github.com/wangshizebin/jiebago/tokenizer"
)

//export Cut
func Cut(sentence string) string {
	words := jiebago.Cut(sentence)
	return wordsToJson(&words)
}

//export CutFull
func CutFull(sentence string) string {
	words := jiebago.Cut(sentence)
	return wordsToJson(&words)
}

//export CutAccurate
func CutAccurate(sentence string) string {
	words := jiebago.Cut(sentence)
	return wordsToJson(&words)
}

//export CutNoHMM
func CutNoHMM(sentence string) string {
	words := jiebago.Cut(sentence)
	return wordsToJson(&words)
}

//export CutForSearch
func CutForSearch(sentence string) string {
	words := jiebago.Cut(sentence)
	return wordsToJson(&words)
}

//export ExtractKeywords
func ExtractKeywords(s string, count int) string {
	keywords := tokenizer.GetTFIDF().ExtractKeywords(s, count, false)
	words := keywords.([]string)
	return wordsToJson(&words)
}

//export ExtractKeywordsWeight
func ExtractKeywordsWeight(s string, count int) string {
	keywords := tokenizer.GetTFIDF().ExtractKeywords(s, count, true)
	tags := []tokenizer.Keyword(keywords.(tokenizer.Keywords))
	return wordsWeightToJson(&tags)
}

//export AddDictWord
func AddDictWord(word string, freq int, prop string) bool {
	_, err := tokenizer.GetDictionary().AddWord(word, freq, prop)
	if err != nil {
		return false
	}
	return true
}

//export AddStopWord
func AddStopWord(word string) bool {
	_, err := tokenizer.GetTFIDF().AddStopWord(word)
	if err != nil {
		return false
	}
	return true
}

func wordsToJson(words *[]string) string {
	w := struct {
		Words *[]string `json:"words"`
	}{
		Words: words,
	}
	v, _ := json.Marshal(w)
	return string(v)
}

func wordsWeightToJson(tags *[]tokenizer.Keyword) string {
	w := struct {
		Tags *[]tokenizer.Keyword `json:"tags"`
	}{
		Tags: tags,
	}
	v, _ := json.Marshal(w)
	return string(v)
}

func main() {
	// Need a main function to make CGO compile package as C shared library
}
