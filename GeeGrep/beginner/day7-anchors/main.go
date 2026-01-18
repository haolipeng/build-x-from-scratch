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

	// 创建匹配器并设置选项
	matcher := NewRegexMatcher(pattern)
	matcher.SetIgnoreCase(args.Options.IgnoreCase)

	// 创建输出格式化器
	formatter := NewOutputFormatter(args.Options)

	// 执行搜索
	var hasMatch bool
	if len(args.Files) == 0 {
		// 从标准输入读取
		hasMatch, err = searchReader(os.Stdin, matcher, formatter, "")
	} else {
		// 从文件读取
		hasMatch, err = searchFiles(args.Files, matcher, formatter)
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
func searchReader(reader io.Reader, matcher *RegexMatcher, formatter *OutputFormatter, filename string) (bool, error) {
	scanner := bufio.NewScanner(reader)
	lineNumber := 0
	matchCount := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		// 执行匹配
		matched := matcher.Match(line)

		// 创建匹配结果
		result := &MatchResult{
			LineNumber: lineNumber,
			Line:       line,
			Matched:    matched,
		}

		// 如果应该打印这一行（考虑 -v 选项）
		if formatter.ShouldPrint(matched) {
			matchCount++
			formatter.PrintLine(result)
		}
	}

	if err := scanner.Err(); err != nil {
		return false, fmt.Errorf("error reading input: %v", err)
	}

	// 如果是 -c 选项，打印统计结果
	if filename != "" {
		formatter.PrintFileCount(filename, matchCount)
	} else {
		formatter.PrintCount(matchCount)
	}

	return matchCount > 0, nil
}

// searchFiles 在多个文件中搜索
func searchFiles(files []string, matcher *RegexMatcher, formatter *OutputFormatter) (bool, error) {
	hasMatch := false

	for _, filename := range files {
		file, err := os.Open(filename)
		if err != nil {
			return hasMatch, fmt.Errorf("cannot open file %s: %v", filename, err)
		}
		defer file.Close()

		matched, err := searchReader(file, matcher, formatter, filename)
		if err != nil {
			return hasMatch, err
		}

		if matched {
			hasMatch = true
		}
	}

	return hasMatch, nil
}
