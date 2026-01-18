package main

import (
	"flag"
	"fmt"
)

// Options 存储命令行选项
type Options struct {
	LineNumber  bool // -n: 显示行号
	IgnoreCase  bool // -i: 忽略大小写
	InvertMatch bool // -v: 反向匹配
	Count       bool // -c: 统计数量
}

// Args 存储解析后的命令行参数
type Args struct {
	Pattern string   // 搜索模式（正则表达式）
	Files   []string // 文件列表（为空则从标准输入读取）
	Options *Options // 命令行选项
}

// ParseArgs 解析命令行参数
// 支持的选项：
//   -n: 显示行号
//   -i: 忽略大小写
//   -v: 反向匹配
//   -c: 统计数量
func ParseArgs(args []string) (*Args, error) {
	// 创建 FlagSet
	fs := flag.NewFlagSet("grep", flag.ContinueOnError)

	// 定义选项
	opts := &Options{}
	fs.BoolVar(&opts.LineNumber, "n", false, "显示行号")
	fs.BoolVar(&opts.IgnoreCase, "i", false, "忽略大小写")
	fs.BoolVar(&opts.InvertMatch, "v", false, "反向匹配")
	fs.BoolVar(&opts.Count, "c", false, "统计数量")

	// 解析选项
	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	// 获取非选项参数
	remaining := fs.Args()

	if len(remaining) < 1 {
		return nil, fmt.Errorf("usage: grep [OPTIONS] PATTERN [FILE...]")
	}

	result := &Args{
		Pattern: remaining[0],
		Files:   []string{},
		Options: opts,
	}

	// 如果提供了文件参数
	if len(remaining) > 1 {
		result.Files = remaining[1:]
	}

	return result, nil
}

// Validate 验证参数有效性
func (a *Args) Validate() error {
	if a.Pattern == "" {
		return fmt.Errorf("pattern cannot be empty")
	}
	return nil
}

// String 返回 Args 的字符串表示（用于调试）
func (a *Args) String() string {
	source := "stdin"
	if len(a.Files) > 0 {
		source = fmt.Sprintf("%v", a.Files)
	}

	return fmt.Sprintf("Pattern: %q, Source: %s, Options: %+v",
		a.Pattern, source, a.Options)
}
