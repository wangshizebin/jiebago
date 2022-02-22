// Copyright 2022 Ze-Bin Wang.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package jiebago

import (
	"strings"
	"testing"
)

var (
	sentence   = "Shell位于用户与系统之间，用来帮助用户与操作系统进行沟通。"
	resultTest = []string{"Shell", "操作系统", "沟通"}
	jieBaGo    = NewJieBaGo()
)

func TestCut(t *testing.T) {
	testCutWords(jieBaGo.Cut, t)
}

func TestCutFull(t *testing.T) {
	testCutWords(jieBaGo.CutFull, t)
}

func TestCutAccurate(t *testing.T) {
	testCutWords(jieBaGo.CutAccurate, t)
}

func TestCutNoHMM(t *testing.T) {
	testCutWords(jieBaGo.CutNoHMM, t)
}

func TestCutForSearch(t *testing.T) {
	testCutWords(jieBaGo.CutForSearch, t)
}

func TestExtractKeywords(t *testing.T) {
	t.Log("原始语句： " + sentence)

	words := jieBaGo.ExtractKeywords(sentence, 20)
	t.Log("提取关键字：", strings.Join(words, "/"))
	for _, word := range resultTest {
		ok := false
		for _, v := range words {
			if word == v {
				ok = true
			}
		}
		if !ok {
			t.Error(word + " not pass")
		} else {
			t.Log(word + " OK")
		}
	}
}

func TestExtractKeywordsWeight(t *testing.T) {
	t.Log("原始语句： " + sentence)

	words := jieBaGo.ExtractKeywordsWeight(sentence, 20)
	t.Log("提取关键字：", words)
	for _, word := range resultTest {
		ok := false
		for _, v := range words {
			if word == v.Word {
				ok = true
			}
		}
		if !ok {
			t.Error(word + " not pass")
		} else {
			t.Log(word + " OK")
		}
	}
}

func TestAddWordToDict(t *testing.T) {
	words := []string{"编程宝库", "王泽宾", "codebaoku"}
	t.Log("加入词典：", words)
	for _, word := range words {
		exist, err := jieBaGo.AddDictWord(word, 3, "n")
		if err != nil {
			t.Error(err)
		} else {
			if exist {
				t.Log(word + " 已经存在")
			} else {
				t.Log(word + " 添加入库")
			}
		}
	}
}

func TestAddStopWord(t *testing.T) {
	words := []string{"the", "of", "is"}
	t.Log("加入停止词：", words)
	for _, word := range words {
		exist, err := jieBaGo.AddStopWord(word)
		if err != nil {
			t.Error(err)
		} else {
			if exist {
				t.Log(word + " 已经存在")
			} else {
				t.Log(word + " 添加入库")
			}
		}
	}
}

func testCutWords(f func(string) []string, t *testing.T) {
	t.Log("原始语句： " + sentence)

	wordsResult := f(sentence)
	t.Log("分词结果：", strings.Join(wordsResult, "/"))
	for _, word := range resultTest {
		ok := false
		for _, v := range wordsResult {
			if word == v {
				ok = true
			}
		}
		if !ok {
			t.Error(word + " not pass")
		} else {
			t.Log(word + " OK")
		}
	}
}
