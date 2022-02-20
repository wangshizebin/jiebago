// Copyright 2022 Ze-Bin Wang.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package jiebago

import (
	"strings"

	"github.com/wangshizebin/jiebago/tokenizer"
)

func Cut(sentence string) []string {
	return CutAccurate(sentence)
}

func CutFull(s string) []string {
	wordsRet := make([]string, 0, tokenizer.DefaultWordsLen)

	segments := tokenizer.SplitTextSeg(s)
	for _, segment := range segments {
		if strings.Trim(segment, " ") == "" {
			continue
		}
		if tokenizer.IsTextChars(segment) {
			tokenizer.CutFullW(segment, &wordsRet)
		} else {
			tokenizer.CutSymbolW(segment, &wordsRet)
		}
	}
	return wordsRet
}

func CutAccurate(s string) []string {
	wordsRet := make([]string, 0, tokenizer.DefaultWordsLen)

	segments := tokenizer.SplitTextSeg(s)
	for _, segment := range segments {
		if strings.Trim(segment, " ") == "" {
			continue
		}
		if tokenizer.IsTextChars(segment) {
			tokenizer.CutAccurateW(segment, &wordsRet)
		} else {
			tokenizer.CutSymbolW(segment, &wordsRet)
		}
	}

	return wordsRet
}

func CutNoHMM(s string) []string {
	wordsRet := make([]string, 0, tokenizer.DefaultWordsLen)

	segments := tokenizer.SplitTextSeg(s)
	for _, segment := range segments {
		if strings.Trim(segment, " ") == "" {
			continue
		}
		if tokenizer.IsTextChars(segment) {
			tokenizer.CutNoHMMW(segment, &wordsRet)
		} else {
			tokenizer.CutSymbolW(segment, &wordsRet)
		}
	}

	return wordsRet
}

func CutForSearch(s string) []string {
	wordsRet := make([]string, 0, tokenizer.DefaultWordsLen)

	segments := tokenizer.SplitTextSeg(s)
	for _, segment := range segments {
		if strings.Trim(segment, " ") == "" {
			continue
		}
		if tokenizer.IsTextChars(segment) {
			cutForSearchW(segment, &wordsRet)
		} else {
			tokenizer.CutSymbolW(segment, &wordsRet)
		}
	}

	return wordsRet
}

func cutForSearchW(s string, words *[]string) {
	dictionary := tokenizer.GetDictionary()

	for _, word := range CutAccurate(s) {
		wordRune := []rune(word)
		if len(wordRune) > 2 {
			for i := 0; i < len(wordRune)-1; i++ {
				s := string(wordRune[i : i+2])
				if dictionary.Exist(s) {
					*words = append(*words, s)
				}
			}
		}
		if len(wordRune) > 3 {
			for i := 0; i < len(wordRune)-2; i++ {
				s := string(wordRune[i : i+3])
				if dictionary.Exist(s) {
					*words = append(*words, s)
				}
			}
		}
		*words = append(*words, word)
	}
}

func ExtractKeywords(s string, count int) []string {
	keywords := tokenizer.GetTFIDF().ExtractKeywords(s, count, false)
	return keywords.([]string)
}

func ExtractKeywordsWeight(s string, count int) []tokenizer.Keyword {
	keywords := tokenizer.GetTFIDF().ExtractKeywords(s, count, true)
	return []tokenizer.Keyword(keywords.(tokenizer.Keywords))
}

func AddDictWord(word string, freq int, prop string) (exist bool, err error) {
	return tokenizer.GetDictionary().AddWord(word, freq, prop)
}

func AddStopWord(word string) (exist bool, err error) {
	return tokenizer.GetTFIDF().AddStopWord(word)
}
