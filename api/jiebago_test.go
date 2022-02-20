// Copyright 2022 Ze-Bin Wang.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/wangshizebin/jiebago/tokenizer"
)

func TestCutWordsGet(t *testing.T) {
	sentence := "Shell位于用户与系统之间，用来帮助用户与操作系统进行沟通。"
	t.Log(sentence)

	url := "http://localhost:8118/cut_words?s=" + sentence
	modes := []string{"", "accurate", "full", "nohmm", "search"}
	for _, mode := range modes {
		t.Log("=== mode: " + mode)
		result, err := Get(url + "&mode=" + mode)
		if err != nil {
			t.Error(err)
			return
		}
		var w struct {
			Words []string `json:"words"`
		}
		err = json.Unmarshal([]byte(result), &w)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("结果：", strings.Join(w.Words, "/"))

		words := []string{"Shell", "操作系统", "沟通"}
		for _, word := range words {
			ok := false
			for _, v := range w.Words {
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
}

func TestCutWordsPost(t *testing.T) {
	sentence := "Shell位于用户与系统之间，用来帮助用户与操作系统进行沟通。"
	t.Log(sentence)

	url := "http://localhost:8118/cut_words"
	modes := []string{"", "accurate", "full", "nohmm", "search"}
	for _, mode := range modes {
		t.Log("=== mode: " + mode)
		data := fmt.Sprintf(`{"s":"%s", "mode":"%s"}`, sentence, mode)
		result, err := Post(url, data, "application/json")
		if err != nil {
			t.Error(err)
		}
		var w struct {
			Words []string `json:"words"`
		}
		err = json.Unmarshal([]byte(result), &w)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("结果：", strings.Join(w.Words, "/"))

		words := []string{"Shell", "操作系统", "沟通"}
		for _, word := range words {
			ok := false
			for _, v := range w.Words {
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
}

func TestExtractKeywordsGet(t *testing.T) {
	sentence := "Shell位于用户与系统之间，用来帮助用户与操作系统进行沟通。"
	t.Log(sentence)

	url := "http://localhost:8118/extract_keywords?s=" + sentence + "&count=3"
	result, err := Get(url)
	if err != nil {
		t.Error(err)
		return
	}

	var w struct {
		Tags []string `json:"tags"`
	}

	err = json.Unmarshal([]byte(result), &w)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("结果：", strings.Join(w.Tags, "/"))

	words := []string{"Shell", "操作系统", "用户"}
	for _, word := range words {
		ok := false
		for _, v := range w.Tags {
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

func TestExtractKeywordsPost(t *testing.T) {
	sentence := "Shell位于用户与系统之间，用来帮助用户与操作系统进行沟通。"
	t.Log(sentence)

	url := "http://localhost:8118/extract_keywords"
	data := fmt.Sprintf(`{"s":"%s", "count":%d}`, sentence, 3)
	result, err := Post(url, data, "application/json")
	if err != nil {
		t.Error(err)
		return
	}

	var w struct {
		Tags []string `json:"tags"`
	}

	err = json.Unmarshal([]byte(result), &w)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("结果：", strings.Join(w.Tags, "/"))

	words := []string{"Shell", "操作系统", "用户"}
	for _, word := range words {
		ok := false
		for _, v := range w.Tags {
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

func TestExtractKeywordsWeightGet(t *testing.T) {
	sentence := "Shell位于用户与系统之间，用来帮助用户与操作系统进行沟通。"
	t.Log(sentence)

	url := "http://localhost:8118/extract_keywords?s=" + sentence + "&mode=weight&count=3"
	result, err := Get(url)
	if err != nil {
		t.Error(err)
		return
	}

	var w struct {
		Tags []tokenizer.Keyword `json:"tags"`
	}

	err = json.Unmarshal([]byte(result), &w)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("结果：", w)

	words := []string{"Shell", "操作系统", "用户"}
	for _, word := range words {
		ok := false
		for _, v := range w.Tags {
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

func TestExtractKeywordsWeightPost(t *testing.T) {
	sentence := "Shell位于用户与系统之间，用来帮助用户与操作系统进行沟通。"
	t.Log(sentence)

	url := "http://localhost:8118/extract_keywords"
	data := fmt.Sprintf(`{"s":"%s", "mode":"%s","count":%d}`, sentence, "weight", 3)
	result, err := Post(url, data, "application/json")
	if err != nil {
		t.Error(err)
		return
	}

	var w struct {
		Tags []tokenizer.Keyword `json:"tags"`
	}

	err = json.Unmarshal([]byte(result), &w)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("结果：", w)

	words := []string{"Shell", "操作系统", "用户"}
	for _, word := range words {
		ok := false
		for _, v := range w.Tags {
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

func TestAddDictWordsGet(t *testing.T) {
	word := "编程宝库"
	t.Log("=== 添加字典单词: " + word)
	url := fmt.Sprintf(`http://localhost:8118/add_dict_word?s=%s&weight=%d&prop=%s`, word, 3, "n")
	result, err := Get(url)
	if err != nil {
		t.Error(err)
		return
	}
	var response struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	err = json.Unmarshal([]byte(result), &response)
	if err != nil {
		t.Error(err)
		return
	}
	if response.ErrMsg != "" {
		t.Log(response.ErrMsg)
	}
}

func TestAddDictWordsPost(t *testing.T) {
	url := "http://localhost:8118/add_dict_word"

	word := "编程宝库"
	t.Log("=== 添加字典单词: " + word)
	data := fmt.Sprintf(`{"s":"%s", "weight":%d,"prop":"%s"}`, word, 3, "n")
	result, err := Post(url, data, "application/json")
	if err != nil {
		t.Error(err)
		return
	}
	var response struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	err = json.Unmarshal([]byte(result), &response)
	if err != nil {
		t.Error(err)
		return
	}
	if response.ErrMsg != "" {
		t.Log(response.ErrMsg)
	}
}

func TestAddStopWordsGet(t *testing.T) {
	word := "the"
	t.Log("=== 添加停止词: " + word)
	url := fmt.Sprintf(`http://localhost:8118/add_stop_word?s=%s`, word)
	result, err := Get(url)
	if err != nil {
		t.Error(err)
		return
	}
	var response struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	err = json.Unmarshal([]byte(result), &response)
	if err != nil {
		t.Error(err)
		return
	}
	if response.ErrMsg != "" {
		t.Log(response.ErrMsg)
	}
}

func TestAddStopWordsPost(t *testing.T) {
	url := "http://localhost:8118/add_stop_word"

	word := "the"
	t.Log("=== 添加停止词: " + word)
	data := fmt.Sprintf(`{"s":"%s"}`, word)
	result, err := Post(url, data, "application/json")
	if err != nil {
		t.Error(err)
		return
	}
	var response struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	err = json.Unmarshal([]byte(result), &response)
	if err != nil {
		t.Error(err)
		return
	}
	if response.ErrMsg != "" {
		t.Log(response.ErrMsg)
	}
}

// 发送GET请求
// url：		请求地址
// response：	请求返回的内容
func Get(url string) (string, error) {
	// 超时时间：5秒
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("status code:" + strconv.Itoa(resp.StatusCode))
	}
	var buffer [256]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}
	}

	return result.String(), nil
}

// 发送POST请求
// url：         请求地址
// data：        POST请求提交的数据
// contentType： 请求体格式，如：application/json
// content：     请求放回的内容
func Post(url string, data string, contentType string) (string, error) {
	// 超时时间：5秒
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(url, contentType, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
