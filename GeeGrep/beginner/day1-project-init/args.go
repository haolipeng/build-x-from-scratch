package main

import (
	"fmt"
)

// Args 存储解析后的命令行参数
type Args struct {
	Pattern string   // 搜索模式
	Files   []string // 文件列表（为空则从标准输入读取）
}

// ParseArgs 解析命令行参数
// 参数格式: grep PATTERN [FILE...]
func ParseArgs(args []string) (*Args, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("usage: grep PATTERN [FILE...]")
	}

	result := &Args{
		Pattern: args[0],
		Files:   []string{},
	}

	// 如果提供了文件参数
	if len(args) > 1 {
		result.Files = args[1:]
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
	if len(a.Files) == 0 {
		return fmt.Sprintf("Pattern: %q, Source: stdin", a.Pattern)
	}
	return fmt.Sprintf("Pattern: %q, Files: %v", a.Pattern, a.Files)
}
