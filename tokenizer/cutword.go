// Copyright 2022 Ze-Bin Wang.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package tokenizer

func CutFullW(s string, words *[]string) {
	bufEnglish := ""
	pos := -1

	sentence := NewSentence(s)
	dag := sentence.GetDAG()
	for k, listPos := range dag {
		if len(bufEnglish) > 0 && !IsEnglishChars(sentence.GetChar(k)) {
			*words = append(*words, bufEnglish)
			bufEnglish = ""
		}

		if len(listPos) == 1 && k > pos {
			word := sentence.GetWord(k, listPos[0]+1)
			if IsEnglishChars(word) {
				bufEnglish += word
			}
			if len(bufEnglish) == 0 {
				*words = append(*words, word)
			}
			pos = listPos[0]
		} else {
			for _, j := range listPos {
				if j > k {
					*words = append(*words, sentence.GetWord(k, j+1))
					pos = j
				}
			}
		}
	}

	if len(bufEnglish) > 0 {
		*words = append(*words, bufEnglish)
	}
}

func CutAccurateW(s string, words *[]string) {
	sentence := NewSentence(s)
	route := sentence.CalcDAG()
	dictionary := GetDictionary()
	buf := ""
	for i := 0; i < sentence.Len(); {
		y := route[i].Y + 1
		leftWord := sentence.GetWord(i, y)
		if y-i == 1 {
			buf += leftWord
			i = y
			continue
		}

		if len(buf) > 0 {
			if len([]rune(buf)) == 1 {
				*words = append(*words, buf)
			} else {
				if !dictionary.Exist(buf) {
					wordsRecognized := GetFinalSeg().Cut(buf)
					for _, w := range wordsRecognized {
						*words = append(*words, w)
					}
				} else {
					for _, v := range buf {
						*words = append(*words, string(v))
					}
				}
			}
			buf = ""
		}
		*words = append(*words, leftWord)
		i = y
	}

	if len(buf) > 0 {
		if len([]rune(buf)) == 1 {
			*words = append(*words, buf)
		} else {
			if !dictionary.Exist(buf) {
				wordsRecognized := GetFinalSeg().Cut(buf)
				for _, w := range wordsRecognized {
					*words = append(*words, w)
				}
			} else {
				for _, v := range buf {
					*words = append(*words, string(v))
				}
			}
		}
	}
}

func CutNoHMMW(s string, words *[]string) {
	sentence := NewSentence(s)
	route := sentence.CalcDAG()

	bufEnglish := ""
	for i := 0; i < sentence.Len(); {
		y := route[i].Y + 1
		leftWord := sentence.GetWord(i, y)
		if IsEnglishChars(leftWord) && len(leftWord) == 1 {
			bufEnglish += leftWord
			i = y
			continue
		}

		if len(bufEnglish) > 0 {
			*words = append(*words, bufEnglish)
			bufEnglish = ""
		}
		*words = append(*words, leftWord)
		i = y
	}

	if len(bufEnglish) > 0 {
		*words = append(*words, bufEnglish)
	}
}

func CutSymbolW(s string, words *[]string) {
	n := len(s)
	if n == 0 {
		return
	}

	buf := ""
	word := ""
	prePos := 0
	for {
		loc := reDelimiter.FindStringIndex(s[prePos:])
		if loc == nil {
			word = s[prePos:]
			prePos = n
		} else {
			loc[0] += prePos
			loc[1] += prePos
			if loc[0] > prePos {
				buf = s[prePos:loc[0]]
			}

			word = s[loc[0]:loc[1]]
			prePos = loc[1]
		}

		if buf == "\r" && word == "\n" {
			*words = append(*words, "\r\n")
			buf = ""
		} else {
			if buf != "" {
				*words = append(*words, buf)
			}
			if word != "" {
				buf = word
			}
		}

		if prePos == n {
			if buf != "" {
				*words = append(*words, buf)
			}
			return
		}
	}
}

func Init(dictPath string) {
	SetDictPath(dictPath)
	InitDictionary()
	InitTFIDF()
	InitFSToken()
}
