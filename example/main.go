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
	jieBaGo := jiebago.NewJieBaGo()
	// 可以指定字典库的位置
	// jieBaGo := jiebago.NewJieBaGo("/data/mydict")

	sentence := "Shell 位于用户与系统之间，用来帮助用户与操作系统进行沟通。通常都是文字模式的 Shell。"
	fmt.Println("原始语句：", sentence)
	fmt.Println()

	// 默认模式分词
	words := jieBaGo.Cut(sentence)
	fmt.Println("默认模式分词：", strings.Join(words, "/"))

	// 精确模式分词
	words = jieBaGo.CutAccurate(sentence)
	fmt.Println("精确模式分词：", strings.Join(words, "/"))

	// 全模式分词
	words = jieBaGo.CutFull(sentence)
	fmt.Println("全模式分词：", strings.Join(words, "/"))

	// NoHMM模式分词
	words = jieBaGo.CutNoHMM(sentence)
	fmt.Println("NoHMM模式分词：", strings.Join(words, "/"))

	// 搜索引擎模式分词
	words = jieBaGo.CutForSearch(sentence)
	fmt.Println("搜索引擎模式分词：", strings.Join(words, "/"))
	fmt.Println()

	// 提取关键词，即Tag标签
	keywords := jieBaGo.ExtractKeywords(sentence, 20)
	fmt.Println("提取关键词：", strings.Join(keywords, "/"))

	// 提取带权重的关键词，即Tag标签
	keywordsWeight := jieBaGo.ExtractKeywordsWeight(sentence, 20)
	fmt.Println("提取带权重的关键词：", keywordsWeight)
	fmt.Println()

	// 向字典加入单词
	exist, err := jieBaGo.AddDictWord("编程宝库", 3, "n")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("向字典加入单词：编程宝库")
		if exist {
			fmt.Println("单词已经存在")
		}
	}

	// 向字典加入停止词
	exist, err = jieBaGo.AddStopWord("the")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("向字典加入停止词：the")
		if exist {
			fmt.Println("单词已经存在")
		}
	}
}
