#!/bin/bash

# 第五天课程测试脚本

echo "======================================"
echo "    GeeGrep Day 5 - 测试脚本"
echo "    字符组：[abc], [a-z], [^abc]"
echo "======================================"
echo

# 测试 1: 字符枚举 [abc]
echo "测试 1: 字符枚举 [ch]at"
echo "--------------------------------------"
echo "输入: cat, hat, bat, rat"
echo -e "cat\nhat\nbat\nrat" | go run . "[ch]at"
echo
echo "期望输出: cat 和 hat"
echo

# 测试 2: 字符范围 [a-z]
echo "测试 2: 字符范围 [a-z]"
echo "--------------------------------------"
echo "输入: hello, HELLO, 12345"
echo -e "hello\nHELLO\n12345" | go run . "[a-z]"
echo
echo "期望输出: hello"
echo

# 测试 3: 大写字母范围 [A-Z]
echo "测试 3: 大写字母范围 [A-Z]"
echo "--------------------------------------"
echo -e "ABC\nabc\n123" | go run . "[A-Z]"
echo
echo "期望输出: ABC"
echo

# 测试 4: 数字范围 [0-9]
echo "测试 4: 数字范围 [0-9]"
echo "--------------------------------------"
echo -e "abc\na1b\n123" | go run . "[0-9]"
echo
echo "期望输出: a1b 和 123"
echo

# 测试 5: 否定字符组 [^abc]
echo "测试 5: 否定字符组 [^abc]"
echo "--------------------------------------"
echo -e "abc\nxyz\naxy" | go run . "[^abc]"
echo
echo "期望输出: xyz 和 axy (包含非 a/b/c 的字符)"
echo

# 测试 6: 否定数字 [^0-9]
echo "测试 6: 否定数字 [^0-9]"
echo "--------------------------------------"
echo -e "123\nabc\na1b" | go run . "[^0-9]"
echo
echo "期望输出: abc 和 a1b"
echo

# 测试 7: 组合字符组 [a-zA-Z]
echo "测试 7: 组合字符组 [a-zA-Z]"
echo "--------------------------------------"
echo -e "abc\nABC\n123\n@#\$" | go run . "[a-zA-Z]"
echo
echo "期望输出: abc 和 ABC"
echo

# 测试 8: 字母数字 [a-zA-Z0-9]
echo "测试 8: 字母数字 [a-zA-Z0-9]"
echo "--------------------------------------"
echo -e "abc\n123\n@#\$" | go run . "[a-zA-Z0-9]"
echo
echo "期望输出: abc 和 123"
echo

# 测试 9: 十六进制 [0-9a-fA-F]
echo "测试 9: 十六进制 0x[0-9a-fA-F][0-9a-fA-F]"
echo "--------------------------------------"
echo -e "0xff\n0xGG\n0x1a" | go run . "0x[0-9a-fA-F][0-9a-fA-F]"
echo
echo "期望输出: 0xff 和 0x1a"
echo

# 测试 10: 元音字母 [aeiou]
echo "测试 10: 元音字母 [aeiou]"
echo "--------------------------------------"
echo -e "apple\nsky\norange" | go run . "[aeiou]"
echo
echo "期望输出: apple 和 orange"
echo

# 测试 11: 非元音字母 [^aeiou]
echo "测试 11: 非元音字母 [^aeiou]"
echo "--------------------------------------"
echo -e "aaa\nbcd\neee" | go run . "[^aeiou]"
echo
echo "期望输出: bcd"
echo

# 测试 12: 字符组 + -n 选项
echo "测试 12: 字符组 + -n 选项"
echo "--------------------------------------"
echo -e "abc\n123\nxyz" | go run . -n "[0-9]"
echo
echo "期望输出: 2:123"
echo

# 测试 13: 字符组 + -v 选项
echo "测试 13: 字符组 + -v 选项"
echo "--------------------------------------"
echo -e "abc\n123\nxyz" | go run . -v "[0-9]"
echo
echo "期望输出: abc 和 xyz"
echo

# 测试 14: 字符组 + -c 选项
echo "测试 14: 字符组 + -c 选项"
echo "--------------------------------------"
echo -e "a1\nb2\nc3\nde" | go run . -c "[0-9]"
echo
echo "期望输出: 3"
echo

# 测试 15: 复杂模式 - ID 格式
echo "测试 15: ID 格式 ID:[0-9][0-9][0-9]"
echo "--------------------------------------"
echo "ID:123" | go run . "ID:[0-9][0-9][0-9]"
echo
echo "期望输出: ID:123"
echo

# 测试 16: 字符组与字符类混合
echo "测试 16: 字符组与字符类混合"
echo "--------------------------------------"
echo "abc123" | go run . "[a-z]\\d\\d\\d"
echo
echo "期望输出: abc123 (注: 只匹配 c123 部分)"
echo

# 测试 17: 字符组中的特殊字符
echo "测试 17: 字符组中的点号 [.]"
echo "--------------------------------------"
echo -e "3.14\n3a14" | go run . "3[.]14"
echo
echo "期望输出: 3.14 (点号在字符组中是字面字符)"
echo

# 测试 18: 字符组中的连字符
echo "测试 18: 字符组中的连字符 [-a]"
echo "--------------------------------------"
echo -e "a-b\naxb" | go run . "a[-]b"
echo
echo "期望输出: a-b"
echo

# 测试 19: 退出码检查
echo "测试 19: 退出码检查"
echo "--------------------------------------"
echo "测试包含字母的退出码:"
echo "abc" | go run . "[a-z]" > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "✓ 找到匹配: 退出码 0"
else
    echo "✗ 错误的退出码"
fi

echo "测试不包含字母的退出码:"
echo "123" | go run . "[a-z]" > /dev/null 2>&1
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
echo "Day 5 新增功能:"
echo "  ✓ [abc]: 字符枚举"
echo "  ✓ [a-z]: 字符范围"
echo "  ✓ [A-Z]: 大写字母范围"
echo "  ✓ [0-9]: 数字范围"
echo "  ✓ [^abc]: 否定字符组"
echo "  ✓ [a-zA-Z0-9]: 组合字符组"
echo "  ✓ 与命令行选项组合使用"
