#!/bin/bash

# 第六天课程测试脚本

echo "======================================"
echo "    GeeGrep Day 6 - 测试脚本"
echo "    组合字符类：[\\d\\s], [\\w-], 等"
echo "======================================"
echo

# 测试 1: [\d\s] 匹配数字或空白
echo "测试 1: [\\d\\s] 匹配数字或空白"
echo "--------------------------------------"
echo "输入: abc, 123, (空格行)"
echo -e "abc\n123\n " | go run . "[\\d\\s]"
echo
echo "期望输出: 123 和空格行"
echo

# 测试 2: [\w-] 匹配单词字符或连字符
echo "测试 2: [\\w-] 匹配单词字符或连字符"
echo "--------------------------------------"
echo -e "hello-world\nhello@world\nhello_world" | go run . "hello[\\w-]world"
echo
echo "期望输出: hello-world 和 hello_world"
echo

# 测试 3: [a-z\d] 匹配小写字母或数字
echo "测试 3: [a-z\\d] 匹配小写字母或数字"
echo "--------------------------------------"
echo -e "abc\nABC\n123\na1b" | go run . "[a-z\\d]"
echo
echo "期望输出: abc, 123, a1b"
echo

# 测试 4: [^\d] 匹配非数字
echo "测试 4: [^\\d] 匹配非数字"
echo "--------------------------------------"
echo -e "123\nabc\na1b" | go run . "[^\\d]"
echo
echo "期望输出: abc 和 a1b"
echo

# 测试 5: [\d.] 匹配数字或点号
echo "测试 5: [\\d.] 匹配数字或点号"
echo "--------------------------------------"
echo -e "3.14\n3a14\n3.14.15" | go run . "[\\d.][\\d.][\\d.][\\d.]"
echo
echo "期望输出: 3.14 和 3.14.15"
echo

# 测试 6: [^\w\s] 匹配非单词字符且非空白
echo "测试 6: [^\\w\\s] 匹配特殊字符"
echo "--------------------------------------"
echo -e "hello\nhello@world\n123" | go run . "[^\\w\\s]"
echo
echo "期望输出: hello@world (包含 @)"
echo

# 测试 7: [\w\d] 组合（\w 已包含 \d）
echo "测试 7: [\\w\\d] 组合"
echo "--------------------------------------"
echo -e "abc\n123\n@#$" | go run . "[\\w\\d]"
echo
echo "期望输出: abc 和 123"
echo

# 测试 8: 复杂组合 [a-zA-Z\d_-]
echo "测试 8: 复杂组合 [a-zA-Z\\d_-]"
echo "--------------------------------------"
echo -e "my-var_123\nmy var\nmy@var" | go run . "[a-zA-Z\\d_-][a-zA-Z\\d_-][a-zA-Z\\d_-][a-zA-Z\\d_-][a-zA-Z\\d_-][a-zA-Z\\d_-][a-zA-Z\\d_-][a-zA-Z\\d_-][a-zA-Z\\d_-][a-zA-Z\\d_-]"
echo
echo "期望输出: my-var_123"
echo

# 测试 9: 与 -n 选项组合
echo "测试 9: 组合字符类 + -n 选项"
echo "--------------------------------------"
echo -e "abc\n123\nxyz" | go run . -n "[\\d]"
echo
echo "期望输出: 2:123"
echo

# 测试 10: 与 -v 选项组合
echo "测试 10: 组合字符类 + -v 选项"
echo "--------------------------------------"
echo -e "abc\n123\nxyz" | go run . -v "[\\d]"
echo
echo "期望输出: abc 和 xyz"
echo

# 测试 11: 与 -c 选项组合
echo "测试 11: 组合字符类 + -c 选项"
echo "--------------------------------------"
echo -e "a1\nb2\nc3\nde" | go run . -c "[\\d]"
echo
echo "期望输出: 3"
echo

# 测试 12: 否定组合字符类 [^\d\s]
echo "测试 12: 否定组合 [^\\d\\s]"
echo "--------------------------------------"
echo -e "123\n   \nabc" | go run . "[^\\d\\s]"
echo
echo "期望输出: abc"
echo

# 测试 13: 十六进制匹配
echo "测试 13: 十六进制 [\\da-fA-F]"
echo "--------------------------------------"
echo -e "0xff\n0xGG\n0x1a" | go run . "0x[\\da-fA-F][\\da-fA-F]"
echo
echo "期望输出: 0xff 和 0x1a"
echo

# 测试 14: 退出码检查
echo "测试 14: 退出码检查"
echo "--------------------------------------"
echo "测试匹配的退出码:"
echo "123" | go run . "[\\d\\s]" > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "✓ 找到匹配: 退出码 0"
else
    echo "✗ 错误的退出码"
fi

echo "测试不匹配的退出码:"
echo "abc" | go run . "[\\d\\s]" > /dev/null 2>&1
if [ $? -eq 1 ]; then
    echo "✓ 未找到匹配: 退出码 1"
else
    echo "✗ 错误的退出码"
fi
echo

echo "======================================"
echo "    测试完成！"
echo "======================================"
echo
echo "Day 6 新增功能:"
echo "  ✓ [\\d\\s]: 字符组内使用字符类"
echo "  ✓ [\\w-]: 字符类与字面字符组合"
echo "  ✓ [a-z\\d]: 范围与字符类组合"
echo "  ✓ [^\\d]: 否定字符类"
echo "  ✓ 复杂组合模式"
echo "  ✓ 与命令行选项组合使用"
