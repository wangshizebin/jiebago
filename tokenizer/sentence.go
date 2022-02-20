// Copyright 2022 Ze-Bin Wang.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package tokenizer

import "math"

type NodeDAG struct {
	X float64
	Y int
}

type Sentence struct {
	sentenceRune []rune
}

func NewSentence(s string) *Sentence {
	return &Sentence{
		sentenceRune: []rune(s),
	}
}

func (s *Sentence) Len() int {
	if s.sentenceRune == nil {
		return 0
	}
	return len(s.sentenceRune)
}

func (s *Sentence) GetWord(start, end int) string {
	if (start < 0 || start >= s.Len()) || (end <= 0 || end > s.Len()) {
		return ""
	}
	return string(s.sentenceRune[start:end])
}

func (s *Sentence) GetChar(i int) string {
	if i < 0 || i >= s.Len() {
		return ""
	}
	return string(s.sentenceRune[i])
}

func (s *Sentence) GetDAG() [][]int {
	dag := make([][]int, 0)

	dictionary := GetDictionary()
	n := s.Len()
	for k := 0; k < n; k++ {
		l := make([]int, 0)

		i := k
		for i < n {
			word := string(s.sentenceRune[k : i+1])
			if freq, ok := dictionary.GetWord(word); ok {
				if freq > 0 {
					l = append(l, i)
				}
				i++
				continue
			}
			break
		}

		if len(l) == 0 {
			l = append(l, k)
		}
		dag = append(dag, l)
	}
	return dag
}

func (s *Sentence) CalcDAG() []NodeDAG {
	n := s.Len()
	route := make([]NodeDAG, n+1)
	route[n] = NodeDAG{0, 0}

	dictionary := GetDictionary()
	logTotal := math.Log(dictionary.GetTotalFreq())

	dag := s.GetDAG()
	for k := n - 1; k >= 0; k-- {
		score := float64(0)
		idx := -1
		for _, x := range dag[k] {
			word := s.GetWord(k, x+1)
			freq := 0
			if v, ok := dictionary.GetWord(word); ok {
				freq = v
			}
			if freq == 0 {
				freq = 1
			}
			val := math.Log(float64(freq)) - logTotal + route[x+1].X
			if idx == -1 {
				score = val
				idx = x
			} else if val >= score {
				score = val
				idx = x
			}
		}
		route[k] = NodeDAG{score, idx}
	}
	return route
}
