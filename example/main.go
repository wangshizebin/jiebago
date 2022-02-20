// Copyright 2022 Ze-Bin Wang.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strings"

	"github.com/wangshizebin/jiebago"
)

func main() {
	sentence := "Shell位于用户与系统之间，用来帮助用户与操作系统进行沟通。"
	fmt.Println("原始语句：", sentence)
	fmt.Println()

	// 默认模式分词
	words := jiebago.Cut(sentence)
	fmt.Println("默认模式分词：", strings.Join(words, "/"))

	// 精确模式分词
	words = jiebago.CutAccurate(sentence)
	fmt.Println("精确模式分词：", strings.Join(words, "/"))

	// 全模式分词
	words = jiebago.CutFull(sentence)
	fmt.Println("全模式分词：", strings.Join(words, "/"))

	// NoHMM模式分词
	words = jiebago.CutNoHMM(sentence)
	fmt.Println("NoHMM模式分词：", strings.Join(words, "/"))

	// 搜索引擎模式分词
	words = jiebago.CutForSearch(sentence)
	fmt.Println("搜索引擎模式分词：", strings.Join(words, "/"))
	fmt.Println()

	// 提取关键词，即Tag标签
	keywords := jiebago.ExtractKeywords(sentence, 20)
	fmt.Println("提取关键词：", strings.Join(keywords, "/"))

	// 提取带权重的关键词，即Tag标签
	keywordsWeight := jiebago.ExtractKeywordsWeight(sentence, 20)
	fmt.Println("提取带权重的关键词：", keywordsWeight)
	fmt.Println()

	// 向字典加入单词
	exist, err := jiebago.AddDictWord("编程宝库", 3, "n")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("向字典加入单词：编程宝库")
		if exist {
			fmt.Println("单词已经存在")
		}
	}

	// 向字典加入停止词
	exist, err = jiebago.AddStopWord("the")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("向字典加入停止词：the")
		if exist {
			fmt.Println("单词已经存在")
		}
	}
}
