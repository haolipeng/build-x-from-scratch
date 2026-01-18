package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Matcher 定义匹配器接口
// 这个接口为后续扩展正则表达式匹配做准备
type Matcher interface {
	Match(line string) bool
	Pattern() string
}

// LiteralMatcher 实现字面字符串匹配
// 简单检查行中是否包含指定的模式字符串
type LiteralMatcher struct {
	pattern string
}

// NewLiteralMatcher 创建一个新的字面字符串匹配器
func NewLiteralMatcher(pattern string) *LiteralMatcher {
	return &LiteralMatcher{pattern: pattern}
}

// Match 检查一行是否匹配模式
func (m *LiteralMatcher) Match(line string) bool {
	return strings.Contains(line, m.pattern)
}

// Pattern 返回匹配模式
func (m *LiteralMatcher) Pattern() string {
	return m.pattern
}

// SearchResult 存储搜索结果
type SearchResult struct {
	HasMatches bool // 是否找到匹配
	MatchCount int  // 匹配行数
}

// Search 在输入流中搜索匹配的行
// reader: 输入流
// matcher: 匹配器
// 返回: 搜索结果和可能的错误
func Search(reader io.Reader, matcher Matcher) (*SearchResult, error) {
	scanner := bufio.NewScanner(reader)
	result := &SearchResult{
		HasMatches: false,
		MatchCount: 0,
	}

	// 逐行读取并检查是否匹配
	for scanner.Scan() {
		line := scanner.Text()

		if matcher.Match(line) {
			fmt.Println(line)
			result.HasMatches = true
			result.MatchCount++
		}
	}

	// 检查扫描过程中是否发生错误
	if err := scanner.Err(); err != nil {
		return result, fmt.Errorf("error reading input: %v", err)
	}

	return result, nil
}

// SearchFiles 在多个文件中搜索
// files: 文件列表
// matcher: 匹配器
// 返回: 搜索结果和可能的错误
func SearchFiles(files []string, matcher Matcher) (*SearchResult, error) {
	totalResult := &SearchResult{
		HasMatches: false,
		MatchCount: 0,
	}

	for _, filename := range files {
		file, err := os.Open(filename)
		if err != nil {
			return totalResult, fmt.Errorf("cannot open file %s: %v", filename, err)
		}
		defer file.Close()

		result, err := Search(file, matcher)
		if err != nil {
			return totalResult, err
		}

		// 合并结果
		if result.HasMatches {
			totalResult.HasMatches = true
			totalResult.MatchCount += result.MatchCount
		}
	}

	return totalResult, nil
}
