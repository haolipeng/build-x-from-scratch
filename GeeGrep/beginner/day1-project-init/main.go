package main

import (
	"fmt"
	"os"
)

func main() {
	// 解析命令行参数
	// os.Args[0] 是程序名，os.Args[1:] 是实际参数
	args, err := ParseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2) // 2 表示参数错误
	}

	// 验证参数
	if err := args.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}

	// 创建匹配器
	matcher := NewLiteralMatcher(args.Pattern)

	// 执行搜索
	var result *SearchResult
	if len(args.Files) == 0 {
		// 从标准输入读取
		result, err = Search(os.Stdin, matcher)
	} else {
		// 从文件读取
		result, err = SearchFiles(args.Files, matcher)
	}

	// 处理错误
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}

	// 根据是否找到匹配设置退出码
	// 0: 找到匹配
	// 1: 未找到匹配
	if result.HasMatches {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
