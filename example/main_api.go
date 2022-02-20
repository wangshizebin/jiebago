// Copyright 2022 Ze-Bin Wang.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/wangshizebin/jiebago"
)

func main() {
	sentence := "Shell 位于用户与系统之间，用来帮助用户与操作系统进行沟通。我们通常所说的 Shell 都指的是文字模式的 Shell。"
	fmt.Println("原始语句：", sentence)
	fmt.Println()

	// 默认模式分词
	words := jiebago.Cut(sentence)
	fmt.Println("默认模式分词：", words)

	// 精确模式分词
	words = jiebago.CutAccurate(sentence)
	fmt.Println("精确模式分词：", words)

	// 全模式分词
	words = jiebago.CutFull(sentence)
	fmt.Println("全模式分词：", words)

	// NoHMM模式分词
	words = jiebago.CutNoHMM(sentence)
	fmt.Println("NoHMM模式分词：", words)

	// 搜索引擎模式分词
	words = jiebago.CutForSearch(sentence)
	fmt.Println("搜索引擎模式分词：", words)
	fmt.Println()

	// 提取关键词，即Tag标签
	keywords := jiebago.ExtractKeywords(sentence, 20)
	fmt.Println("提取关键词：", keywords)

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
