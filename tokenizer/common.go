// Copyright 2022 Ze-Bin Wang.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package tokenizer

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	DictStdFile       = "dict_std_utf8.txt"        // standard dictionary file
	DictUserFile      = "dict_user_utf8.txt"       // user-defined dictionary file
	IDFStdFile        = "idf_std_utf8.txt"         // standard IDF file
	StopWordsStdFile  = "stop_words_std_utf8.txt"  // standard stop words file
	StopWordsUserFile = "stop_words_user_utf8.txt" // user-defined stop words file

	RegExpEnglish   = "([a-zA-Z0-9])+"                     // English regular expression
	RegExpChinese   = "([\u4e00-\u9fa5])+"                 // Chinese regular expression
	RegExpText      = "([\u4e00-\u9fa5a-zA-Z0-9+#&._%-])+" // text regular expression
	RegExpNumber    = "[a-zA-Z0-9]+(\\.\\d+)?%?"           // numeric regular expression
	RegExpDelimiter = "[\\r\\n\\s\\t]"                     // delimiter regular expression

	DefaultWordsLen = 32 // default slice size of the word segmentation result
)

var (
	reEnglish, _   = regexp.Compile(RegExpEnglish)   // precompiled English regular expression
	reChinese, _   = regexp.Compile(RegExpChinese)   // precompiled Chinese regular expression
	reText, _      = regexp.Compile(RegExpText)      // precompiled text regular expression
	reNumber, _    = regexp.Compile(RegExpNumber)    // precompiled numeric regular expression
	reDelimiter, _ = regexp.Compile(RegExpDelimiter) // precompiled delimiter regular expression

	dictPath string // dictionary directory, default is current work directory
)

func IsEnglishChars(s string) bool {
	return reEnglish.Match([]byte(s))
}

func IsChineseChars(s string) bool {
	return reChinese.Match([]byte(s))
}

func IsTextChars(s string) bool {
	return reText.Match([]byte(s))
}

// Split sentence according to normal text
func SplitTextSeg(s string) []string {
	return splitRegExp(s, reText)
}

// Split sentence according to Chinese
func SplitChineseSeg(s string) []string {
	return splitRegExp(s, reChinese)
}

// Split sentence according to number
func SplitNumberSeg(s string) []string {
	return splitRegExp(s, reNumber)
}

// Get the dictionary file directory
func GetDictFile(file string) (string, error) {
	errFileNotFound := errors.New("unable to load the dictionary file")

	dictPath := ""
	if GetDictPath() != "" {
		dictPath = filepath.Join(GetDictPath(), file)
		if !fileExist(dictPath) {
			return "", errFileNotFound
		}
		return dictPath, nil
	}

	dictFile := fmt.Sprintf("%cdictionary%c%s", os.PathSeparator, os.PathSeparator, file)

	// check exe file directory
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Println(err)
		return "", errFileNotFound
	}

	dictPath = path + dictFile
	if !fileExist(dictPath) {
		path, err = os.Getwd()
		if err != nil {
			log.Println(err)
			return "", errFileNotFound
		}
	}

	// check work directory
	dictPath = path + dictFile
	if !fileExist(dictPath) {
		path = getParentPath(path)
		if path == "" {
			return "", errFileNotFound
		}
	}

	// check parent of work directory
	dictPath = path + dictFile
	if !fileExist(dictPath) {
		return "", errFileNotFound
	}

	return dictPath, nil
}

// Split sentence according to the specified regular expression
func splitRegExp(s string, re *regexp.Regexp) []string {
	sentences := make([]string, 0)

	n := len(s)
	prePos := 0
	for {
		loc := re.FindStringIndex(s[prePos:])
		if loc == nil {
			sentences = append(sentences, s[prePos:])
			return sentences
		}
		loc[0] += prePos
		loc[1] += prePos
		if loc[0] > prePos {
			sentences = append(sentences, s[prePos:loc[0]])
		}
		sentences = append(sentences, s[loc[0]:loc[1]])
		prePos = loc[1]
		if prePos == n {
			break
		}
	}
	return sentences
}

func fileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func getParentPath(path string) string {
	return substrRune(path, 0, strings.LastIndex(path, string(os.PathSeparator)))
}

func substrRune(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func GetDictPath() string {
	return dictPath
}

func SetDictPath(path string) {
	dictPath = path
}
