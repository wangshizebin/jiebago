// Copyright 2022 Ze-Bin Wang.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package tokenizer

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	DefaultIDFSize = 300000
)

var tfIDF = &TFIDF{
	idfLoader: &IDFLoader{
		idfFreq: make(map[string]float64),
	},
	stopWords: &StopWords{
		dictMap: make(map[string]struct{}),
	},
}

type Keyword struct {
	Word   string  `json:"word"`
	Weight float64 `json:"weight"`
}

type Keywords []Keyword

func (k Keywords) Len() int {
	return len(k)
}

func (k Keywords) Less(i, j int) bool {
	if k[i].Weight > k[j].Weight {
		return true
	}
	return false
}

func (k Keywords) Swap(i, j int) {
	k[i], k[j] = k[j], k[i]
}

type StopWords struct {
	dictMap map[string]struct{}
	mu      sync.RWMutex
}

func (d *StopWords) load(fileStopWords string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	timeStart := time.Now()

	f, err := os.Open(fileStopWords)
	if err != nil {
		log.Println(err)
		return errors.New("unable to load the stop words library:" + filepath.Base(fileStopWords))
	}
	defer func() {
		_ = f.Close()
	}()

	itemCount := 0
	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			break
		}

		elem := strings.Fields(line)
		if len(elem) == 0 {
			if err == io.EOF {
				break
			}
			continue
		}

		for _, v := range elem {
			if v == "" {
				continue
			}

			itemCount++
			d.dictMap[strings.ToLower(v)] = struct{}{}
		}

		if err == io.EOF {
			break
		}
	}

	log.Printf("%v stop words are loaded, and take %v\n",
		itemCount, time.Now().Sub(timeStart))
	return nil
}

func (d *StopWords) exist(s string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	_, ok := d.dictMap[strings.ToLower(s)]
	return ok
}

func (d *StopWords) add(s string) (exist bool, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	s = strings.ToLower(strings.TrimSpace(s))
	if s == "" {
		return
	}
	if _, ok := d.dictMap[s]; ok {
		exist = true
		return
	}

	stopWordsStdFile, err := GetDictFile(StopWordsStdFile)
	if err != nil {
		return
	}

	stopWordsUserFile := filepath.Dir(stopWordsStdFile)
	stopWordsUserFile += string(filepath.Separator) + StopWordsUserFile
	f, err := os.OpenFile(stopWordsUserFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	defer func() {
		_ = f.Close()
	}()

	stat, err := f.Stat()
	if err != nil {
		return
	}

	line := ""
	n := stat.Size()
	if n > 0 {
		buf := make([]byte, 1, 1)
		_, err = f.ReadAt(buf, n-1)
		if err != nil {
			return
		}
		if buf[0] != '\n' {
			line += "\n"
		}
	}
	line += s + "\n"
	_, err = f.Write([]byte(line))
	if err != nil {
		log.Println(err)
		return
	}

	d.dictMap[s] = struct{}{}
	return
}

type IDFLoader struct {
	idfFreq   map[string]float64
	idfMedian float64
}

func (d *IDFLoader) load(idfFile string) error {
	timeStart := time.Now()

	f, err := os.Open(idfFile)
	if err != nil {
		log.Println(err)
		return errors.New("unable to load the IDF library")
	}
	defer func() {
		_ = f.Close()
	}()

	itemCount := 0
	freqSlice := make([]float64, 0, DefaultIDFSize)
	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			break
		}

		elem := strings.Fields(line)
		if len(elem) != 2 {
			if err == io.EOF {
				break
			}
			continue
		}

		itemCount++
		freq, err := strconv.ParseFloat(elem[1], 64)
		if err != nil {
			log.Println(err)
			freq = float64(0)
		}

		d.idfFreq[strings.ToLower(elem[0])] = freq
		freqSlice = append(freqSlice, freq)

		if err == io.EOF {
			break
		}
	}

	sort.Float64s(freqSlice)
	d.idfMedian = freqSlice[itemCount/2]

	log.Printf("%v idfs are loaded, and take %v\n",
		itemCount, time.Now().Sub(timeStart))
	return nil
}

type TFIDF struct {
	idfLoader *IDFLoader
	stopWords *StopWords
}

func (t *TFIDF) ExtractKeywords(s string, count int, withWeight bool) interface{} {
	words := make([]string, 0, DefaultWordsLen)
	segments := SplitTextSeg(s)
	for _, segment := range segments {
		if IsTextChars(segment) {
			CutAccurateW(segment, &words)
		} else {
			CutSymbolW(segment, &words)
		}
	}

	freqMap, freqMedian := t.idfLoader.idfFreq, t.idfLoader.idfMedian

	freqTotal := 0
	freqWords := make(map[string]int)
	for _, word := range words {
		if len([]rune(word)) < 2 || t.ExistStopWord(word) {
			continue
		}
		if val, ok := freqWords[word]; ok {
			freqWords[word] = val + 1
			freqTotal++
			continue
		}
		freqTotal++
		freqWords[word] = 1
	}

	i := 0
	wordsRet := make(Keywords, len(freqWords))
	for word, s := range freqWords {
		val := freqMedian
		if v, ok := freqMap[word]; ok {
			val = v
		}
		wordsRet[i] = Keyword{
			Word:   word,
			Weight: float64(s) * (val / float64(freqTotal)),
		}
		i++
	}

	sort.Sort(wordsRet)
	if count == 0 {
		count = 20
	}
	if count < len(wordsRet) {
		wordsRet = wordsRet[:count]
	}
	if withWeight {
		return wordsRet
	}
	stringSet := make([]string, len(wordsRet))
	for i, v := range wordsRet {
		stringSet[i] = v.Word
	}
	return stringSet
}

func (t *TFIDF) ExistStopWord(word string) bool {
	return t.stopWords.exist(word)
}

func (t *TFIDF) AddStopWord(word string) (exist bool, err error) {
	return t.stopWords.add(word)
}

func GetTFIDF() *TFIDF {
	return tfIDF
}

func init() {
	// load the tf-idf library
	idfFile, err := GetDictFile(IDFStdFile)
	if err != nil {
		log.Panic(err)
	}

	err = tfIDF.idfLoader.load(idfFile)
	if err != nil {
		log.Panic(err)
	}

	// load the standard stop words library
	stopWordsStdFile, err := GetDictFile(StopWordsStdFile)
	if err == nil {
		tfIDF.stopWords.load(stopWordsStdFile)
	}

	// load the user-defined stop words library
	stopWordsUserFile, err := GetDictFile(StopWordsUserFile)
	if err == nil {
		tfIDF.stopWords.load(stopWordsUserFile)
	}
}
