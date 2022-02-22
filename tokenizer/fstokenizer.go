// Copyright 2022 Ze-Bin Wang.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package tokenizer

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
)

const (
	minFloat          = -3.14e100
	finalSegProbStart = "fs_pbstart.json"
	finalSegProbTrans = "fs_pbtrans.json"
	finalSegProbEmit  = "fs_pbemit.json"
)

var (
	fsTokenizer = &FinalSeg{
		start: make(map[string]float64),
		trans: make(map[string]map[string]float64),
		emit:  make(map[string]map[string]float64),

		forceSplitWords: &forceSplitWords{
			dict: make(map[string]struct{}),
		},
	}

	prevStatus = map[string][]string{
		"B": {"E", "S"},
		"M": {"M", "B"},
		"S": {"S", "E"},
		"E": {"B", "M"},
	}
	states = []string{"B", "M", "E", "S"}
)

type forceSplitWords struct {
	dict  map[string]struct{}
	mutex sync.RWMutex
}

func (s *forceSplitWords) exist(word string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	_, ok := s.dict[word]
	return ok
}

func (s *forceSplitWords) addForceSplit(word string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.dict[word] = struct{}{}
}

type FinalSeg struct {
	start map[string]float64
	trans map[string]map[string]float64
	emit  map[string]map[string]float64

	forceSplitWords *forceSplitWords
}

func (fs *FinalSeg) Cut(sentence string) []string {
	wordsRet := make([]string, 0, DefaultWordsLen)

	segments := SplitChineseSeg(sentence)
	for _, segment := range segments {
		if IsChineseChars(segment) {
			words := fs.cut(segment)
			for _, v := range words {
				if fs.exist(v) {
					for _, c := range v {
						wordsRet = append(wordsRet, string(c))
					}
				} else {
					wordsRet = append(wordsRet, v)
				}
			}
		} else {
			words := SplitNumberSeg(segment)
			wordsRet = append(wordsRet, words...)
		}
	}
	return wordsRet
}

func (fs *FinalSeg) getMatrixVal(name, key, word string) float64 {
	var m map[string]map[string]float64
	if name == "emit" {
		m = fs.emit
	} else if name == "trans" {
		m = fs.trans
	} else {
		return minFloat
	}
	val, ok := m[key][word]
	if !ok {
		val = minFloat
	}
	return val
}

func (fs *FinalSeg) viterbi(sentence string) []string {
	rs := []rune(sentence)
	n := len(rs)
	if n == 0 {
		return nil
	}

	v := make([]map[string]float64, n)
	for i := 0; i < n; i++ {
		v[i] = make(map[string]float64)
	}
	path := make(map[string][]string, 0)

	word := string(rs[0])
	for _, y := range states {
		v[0][y] = fs.start[y] + fs.getMatrixVal("emit", y, word)
		path[y] = []string{y}
	}

	for i := 1; i < len(rs); i++ {
		word = string(rs[i])

		pathNew := make(map[string][]string, 0)
		for _, y := range states {
			st := ""
			pb := minFloat
			for _, y0 := range prevStatus[y] {
				m := v[i-1][y0] + fs.getMatrixVal("trans", y0, y) + fs.getMatrixVal("emit", y, word)
				if st == "" {
					st = y0
					pb = m
				} else if m > pb {
					st = y0
					pb = m
				}
			}
			v[i][y] = pb
			pathNew[y] = append(path[st], y)
		}
		path = pathNew
	}

	state := "E"
	prob := v[len(rs)-1]["E"]
	if v[len(rs)-1]["S"] > prob {
		prob = v[len(rs)-1]["S"]
		state = "S"
	}
	return path[state]
}

func (fs *FinalSeg) cut(sentence string) []string {
	rs := []rune(sentence)
	wordsRet := make([]string, 0)
	posList := fs.viterbi(sentence)

	begin, next := 0, 0
	for i, word := range rs {
		pos := posList[i]
		if pos == "B" {
			begin = i
		} else if pos == "E" {
			wordsRet = append(wordsRet, string(rs[begin:i+1]))
			next = i + 1
		} else if pos == "S" {
			wordsRet = append(wordsRet, string(word))
			next = i + 1
		}
	}
	if next < len(rs) {
		wordsRet = append(wordsRet, string(rs[next:]))
	}

	return wordsRet
}

func (fs *FinalSeg) exist(word string) bool {
	return fs.forceSplitWords.exist(word)
}

func readJsonFromFile(fn string, fs interface{}) {
	fileProbeStart, err := GetDictFile(fn)
	if err != nil {
		log.Panic(err)
	}
	data, err := ioutil.ReadFile(fileProbeStart)
	if err != nil {
		log.Panic(err)
	}
	err = json.Unmarshal(data, fs)
	if err != nil {
		log.Panic(err)
	}
}

func GetFinalSeg() *FinalSeg {
	return fsTokenizer
}

func InitFSToken() {
	readJsonFromFile(finalSegProbStart, &fsTokenizer.start)
	readJsonFromFile(finalSegProbTrans, &fsTokenizer.trans)
	readJsonFromFile(finalSegProbEmit, &fsTokenizer.emit)
}
