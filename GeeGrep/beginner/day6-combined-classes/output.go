package main

import (
	"fmt"
)

// OutputFormatter 处理不同格式的输出
type OutputFormatter struct {
	opts *Options
}

// NewOutputFormatter 创建一个新的输出格式化器
func NewOutputFormatter(opts *Options) *OutputFormatter {
	return &OutputFormatter{opts: opts}
}

// MatchResult 存储单行的匹配结果
type MatchResult struct {
	LineNumber int    // 行号（从 1 开始）
	Line       string // 行内容
	Matched    bool   // 是否匹配
}

// ShouldPrint 判断是否应该打印这一行
// 考虑 -v (反向匹配) 选项
func (f *OutputFormatter) ShouldPrint(matched bool) bool {
	if f.opts.InvertMatch {
		return !matched
	}
	return matched
}

// PrintLine 打印单行结果
func (f *OutputFormatter) PrintLine(result *MatchResult) {
	if !f.ShouldPrint(result.Matched) {
		return
	}

	// -c 选项下不打印具体行，只统计
	if f.opts.Count {
		return
	}

	// -n 选项：显示行号
	if f.opts.LineNumber {
		fmt.Printf("%d:%s\n", result.LineNumber, result.Line)
	} else {
		fmt.Println(result.Line)
	}
}

// PrintCount 打印统计结果
func (f *OutputFormatter) PrintCount(count int) {
	if f.opts.Count {
		fmt.Println(count)
	}
}

// PrintFileCount 打印文件的统计结果（多文件模式）
func (f *OutputFormatter) PrintFileCount(filename string, count int) {
	if f.opts.Count {
		fmt.Printf("%s:%d\n", filename, count)
	}
}
