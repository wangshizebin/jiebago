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

var (
	jieBaGo *jiebago.JieBaGo
)

//export Init
func Init(path string) {
	jieBaGo = jiebago.NewJieBaGo(path)
}

//export Cut
func Cut(sentence string) string {
	if jieBaGo == nil {
		return ""
	}
	words := jieBaGo.Cut(sentence)
	return wordsToJson(&words)
}

//export CutFull
func CutFull(sentence string) string {
	if jieBaGo == nil {
		return ""
	}
	words := jieBaGo.Cut(sentence)
	return wordsToJson(&words)
}

//export CutAccurate
func CutAccurate(sentence string) string {
	if jieBaGo == nil {
		return ""
	}
	words := jieBaGo.Cut(sentence)
	return wordsToJson(&words)
}

//export CutNoHMM
func CutNoHMM(sentence string) string {
	if jieBaGo == nil {
		return ""
	}
	words := jieBaGo.Cut(sentence)
	return wordsToJson(&words)
}

//export CutForSearch
func CutForSearch(sentence string) string {
	if jieBaGo == nil {
		return ""
	}
	words := jieBaGo.Cut(sentence)
	return wordsToJson(&words)
}

//export ExtractKeywords
func ExtractKeywords(s string, count int) string {
	if jieBaGo == nil {
		return ""
	}
	words := jieBaGo.ExtractKeywords(s, count)
	return wordsToJson(&words)
}

//export ExtractKeywordsWeight
func ExtractKeywordsWeight(s string, count int) string {
	if jieBaGo == nil {
		return ""
	}
	tags := jieBaGo.ExtractKeywordsWeight(s, count)
	return wordsWeightToJson(&tags)
}

//export AddDictWord
func AddDictWord(word string, freq int, prop string) bool {
	if jieBaGo == nil {
		return false
	}
	_, err := jieBaGo.AddDictWord(word, freq, prop)
	if err != nil {
		return false
	}
	return true
}

//export AddStopWord
func AddStopWord(word string) bool {
	if jieBaGo == nil {
		return false
	}
	_, err := jieBaGo.AddStopWord(word)
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
