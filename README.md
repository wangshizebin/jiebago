[![logo](http://static.codebaoku.com/images/blogo.png)](http://www.codebaoku.com)

JiebaGo 是 jieba 中文分词的 Go 语言版本。

## 功能特点

+ 支持多种分词方式，包括: 最大概率模式, HMM新词发现模式, 搜索引擎模式, 全模式
+ 支持抽取关键词，包括: 无权重关键词, 权重关键词
+ 支持多种使用方式，包括: Go语言包, Windows Dll, Web API, Docker
+ 支持在线并行添加字典词库和停止词
+ 全部代码使用 go 语言实现，全面兼容 jieba python 词库

## 引用方法

不使用包管理工具：
```bash
go get github.com/wangshizebin/jiebagou
```

使用 go mod 管理：
代码中直接引用 github.com/wangshizebin/jiebagou 即可。

## 功能示例

```golang
package main

import (
	"fmt"
	"strings"

	"github.com/wangshizebin/jiebago"
)

func main() {
	sentence := "Shell 位于用户与系统之间，用来帮助用户与操作系统进行沟通。通常都是文字模式的 Shell。"
	fmt.Println("原始语句：", sentence)
	fmt.Println()

	// 默认模式分词
	words := jiebago.Cut(sentence)
	fmt.Println("默认模式分词：", strings.Join(words,"/"))

	// 精确模式分词
	words = jiebago.CutAccurate(sentence)
	fmt.Println("精确模式分词：", strings.Join(words,"/"))

	// 全模式分词
	words = jiebago.CutFull(sentence)
	fmt.Println("全模式分词：", strings.Join(words,"/"))

	// NoHMM模式分词
	words = jiebago.CutNoHMM(sentence)
	fmt.Println("NoHMM模式分词：", strings.Join(words,"/"))

	// 搜索引擎模式分词
	words = jiebago.CutForSearch(sentence)
	fmt.Println("搜索引擎模式分词：", strings.Join(words,"/"))
	fmt.Println()

	// 提取关键词，即Tag标签
	keywords := jiebago.ExtractKeywords(sentence, 20)
	fmt.Println("提取关键词：", strings.Join(keywords,"/"))

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
```

```
原始语句： Shell位于用户与系统之间，用来帮助用户与操作系统进行沟通。

默认模式分词： Shell/位于/用户/与/系统/之间/，/用来/帮助/用户/与/操作系统/进行/沟通/。
精确模式分词： Shell/位于/用户/与/系统/之间/，/用来/帮助/用户/与/操作系统/进行/沟通/。
全模式分词： Shell/位于/用户/与/系统/之间/，/用来/帮助/用户/与/操作/操作系统/系统/进行/沟通/。
NoHMM模式分词： Shell/位于/用户/与/系统/之间/，/用来/帮助/用户/与/操作系统/进行/沟通/。
搜索引擎模式分词： Shell/位于/用户/与/系统/之间/，/用来/帮助/用户/与/操作/系统/操作系/操作系统/进行/沟通/。

提取关键词： 用户/Shell/操作系统/沟通/帮助/位于/系统/之间/进行
提取带权重的关键词： [{用户 1.364467214484} {Shell 1.19547675029} {操作系统 0.9265948663750001} {沟通 0.694890548758} {帮助 0.5809050240370001} {位于 0.496609078159} {系统 0.49601794343199995} {之间 0.446152979906} {进行 0.372712479502}]

向字典加入单词：编程宝库
向字典加入停止词：the
```

更详细的例子参照 example/main.go, jiebago_test.go, api/iebago_test.go 中的代码。

## 单元测试
go 包

```bash
go test
```

Web API

```bash
cd api
go test 
```

## 注意问题


## Contact

+ Email: `wangzebin@vip.163.com`
+ QQ: `931058240`