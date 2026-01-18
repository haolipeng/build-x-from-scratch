package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	// 解析命令行参数
	args, err := ParseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}

	// 验证参数
	if err := args.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}

	// 编译正则表达式
	pattern, err := Compile(args.Pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error compiling pattern: %v\n", err)
		os.Exit(2)
	}

	// 创建正则表达式匹配器
	matcher := NewRegexMatcher(pattern)

	// 执行搜索
	var hasMatch bool
	if len(args.Files) == 0 {
		// 从标准输入读取
		hasMatch, err = searchReader(os.Stdin, matcher)
	} else {
		// 从文件读取
		hasMatch, err = searchFiles(args.Files, matcher)
	}

	// 处理错误
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}

	// 根据是否找到匹配设置退出码
	if hasMatch {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

// searchReader 在输入流中搜索匹配的行
func searchReader(reader io.Reader, matcher *RegexMatcher) (bool, error) {
	scanner := bufio.NewScanner(reader)
	hasMatch := false

	for scanner.Scan() {
		line := scanner.Text()

		if matcher.Match(line) {
			fmt.Println(line)
			hasMatch = true
		}
	}

	if err := scanner.Err(); err != nil {
		return hasMatch, fmt.Errorf("error reading input: %v", err)
	}

	return hasMatch, nil
}

// searchFiles 在多个文件中搜索
func searchFiles(files []string, matcher *RegexMatcher) (bool, error) {
	hasMatch := false

	for _, filename := range files {
		file, err := os.Open(filename)
		if err != nil {
			return hasMatch, fmt.Errorf("cannot open file %s: %v", filename, err)
		}
		defer file.Close()

		matched, err := searchReader(file, matcher)
		if err != nil {
			return hasMatch, err
		}

		if matched {
			hasMatch = true
		}
	}

	return hasMatch, nil
}
