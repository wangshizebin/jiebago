// Copyright 2022 Ze-Bin Wang.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package jiebago

import (
	"strings"

	"github.com/wangshizebin/jiebago/tokenizer"
)

type JieBaGo struct {
}

func NewJieBaGo(path ...string) *JieBaGo {
	configPath := ""
	if len(path) > 0 {
		configPath = path[0]
	}
	tokenizer.Init(configPath)
	jieBaGo := &JieBaGo{}
	return jieBaGo
}

func (g *JieBaGo) Cut(sentence string) []string {
	return g.CutAccurate(sentence)
}

func (g *JieBaGo) CutFull(s string) []string {
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

func (g *JieBaGo) CutAccurate(s string) []string {
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

func (g *JieBaGo) CutNoHMM(s string) []string {
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

func (g *JieBaGo) CutForSearch(s string) []string {
	wordsRet := make([]string, 0, tokenizer.DefaultWordsLen)

	segments := tokenizer.SplitTextSeg(s)
	for _, segment := range segments {
		if strings.Trim(segment, " ") == "" {
			continue
		}
		if tokenizer.IsTextChars(segment) {
			g.cutForSearchW(segment, &wordsRet)
		} else {
			tokenizer.CutSymbolW(segment, &wordsRet)
		}
	}

	return wordsRet
}

func (g *JieBaGo) cutForSearchW(s string, words *[]string) {
	dictionary := tokenizer.GetDictionary()

	for _, word := range g.CutAccurate(s) {
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

func (g *JieBaGo) ExtractKeywords(s string, count int) []string {
	keywords := tokenizer.GetTFIDF().ExtractKeywords(s, count, false)
	return keywords.([]string)
}

func (g *JieBaGo) ExtractKeywordsWeight(s string, count int) []tokenizer.Keyword {
	keywords := tokenizer.GetTFIDF().ExtractKeywords(s, count, true)
	return []tokenizer.Keyword(keywords.(tokenizer.Keywords))
}

func (g *JieBaGo) AddDictWord(word string, freq int, prop string) (exist bool, err error) {
	return tokenizer.GetDictionary().AddWord(word, freq, prop)
}

func (g *JieBaGo) AddStopWord(word string) (exist bool, err error) {
	return tokenizer.GetTFIDF().AddStopWord(word)
}
