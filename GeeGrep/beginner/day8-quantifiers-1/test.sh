#!/bin/bash

# 第八天课程测试脚本

echo "======================================"
echo "    GeeGrep Day 8 - 测试脚本"
echo "    量词：?, +, *"
echo "======================================"
echo

# 测试 1: ? 零或一次
echo "测试 1: ? 零或一次 - colou?r"
echo "--------------------------------------"
echo "输入: color, colour, colouur"
echo -e "color\ncolour\ncolouur" | go run . "colou?r"
echo
echo "期望输出: color 和 colour"
echo

# 测试 2: ? 零或一次 - ab?c
echo "测试 2: ? 零或一次 - ab?c"
echo "--------------------------------------"
echo -e "ac\nabc\nabbc" | go run . "ab?c"
echo
echo "期望输出: ac 和 abc"
echo

# 测试 3: + 一或多次
echo "测试 3: + 一或多次 - a+"
echo "--------------------------------------"
echo -e "a\naa\naaa\nb" | go run . "a+"
echo
echo "期望输出: a, aa, aaa"
echo

# 测试 4: + 一或多次 - a+bc
echo "测试 4: + 一或多次 - a+bc"
echo "--------------------------------------"
echo -e "abc\naabc\naaabc\nbc" | go run . "a+bc"
echo
echo "期望输出: abc, aabc, aaabc"
echo

# 测试 5: * 零或多次
echo "测试 5: * 零或多次 - ab*c"
echo "--------------------------------------"
echo -e "ac\nabc\nabbc\nabbbc" | go run . "ab*c"
echo
echo "期望输出: 所有行"
echo

# 测试 6: * 零或多次 - a*bc
echo "测试 6: * 零或多次 - a*bc"
echo "--------------------------------------"
echo -e "bc\nabc\naabc\naaabc" | go run . "a*bc"
echo
echo "期望输出: 所有行"
echo

# 测试 7: 字符类 + 量词 - \d+
echo "测试 7: 字符类 + 量词 - a\\d+"
echo "--------------------------------------"
echo -e "a\na1\na12\na123" | go run . "a\\d+"
echo
echo "期望输出: a1, a12, a123"
echo

# 测试 8: \w+ 匹配单词
echo "测试 8: \\w+ 匹配单词"
echo "--------------------------------------"
echo -e "abc\nabc123\n123\n@#$" | go run . "\\w+"
echo
echo "期望输出: abc, abc123, 123"
echo

# 测试 9: 通配符 + 量词 - .+
echo "测试 9: 通配符 + 量词 - a.+c"
echo "--------------------------------------"
echo -e "ac\naXc\naXXc\naXXXc" | go run . "a.+c"
echo
echo "期望输出: aXc, aXXc, aXXXc (不包括 ac)"
echo

# 测试 10: 通配符 + 量词 - .*
echo "测试 10: 通配符 + 量词 - a.*c"
echo "--------------------------------------"
echo -e "ac\naXc\naXXc" | go run . "a.*c"
echo
echo "期望输出: 所有行"
echo

# 测试 11: 锚点 + 量词 - ^hello.*
echo "测试 11: 锚点 + 量词 - ^hello.*"
echo "--------------------------------------"
echo -e "hello\nhello world\nworld hello" | go run . "^hello.*"
echo
echo "期望输出: hello 和 hello world"
echo

# 测试 12: 量词 + 行尾锚点 - \d+$
echo "测试 12: 量词 + 行尾锚点 - \\d+\$"
echo "--------------------------------------"
echo -e "123\nabc123\n123abc" | go run . "\\d+$"
echo
echo "期望输出: 123 和 abc123"
echo

# 测试 13: 字符组 + 量词 - [a-z]+
echo "测试 13: 字符组 + 量词 - [a-z]+"
echo "--------------------------------------"
echo -e "abc\nABC\n123\nabc123" | go run . "[a-z]+"
echo
echo "期望输出: abc 和 abc123"
echo

# 测试 14: https?:// 匹配协议
echo "测试 14: https?:// 匹配协议"
echo "--------------------------------------"
echo -e "http://example.com\nhttps://example.com\nftp://example.com" | go run . "https?://"
echo
echo "期望输出: http://example.com 和 https://example.com"
echo

# 测试 15: 与 -n 选项组合
echo "测试 15: 量词 + -n 选项"
echo "--------------------------------------"
echo -e "a\naa\naaa\nb" | go run . -n "a+"
echo
echo "期望输出: 1:a, 2:aa, 3:aaa"
echo

# 测试 16: 与 -c 选项组合
echo "测试 16: 量词 + -c 选项"
echo "--------------------------------------"
echo -e "a\naa\naaa\nb" | go run . -c "a+"
echo
echo "期望输出: 3"
echo

# 测试 17: 退出码检查
echo "测试 17: 退出码检查"
echo "--------------------------------------"
echo "测试匹配的退出码:"
echo "aaa" | go run . "a+" > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "✓ 找到匹配: 退出码 0"
else
    echo "✗ 错误的退出码"
fi

echo "测试不匹配的退出码:"
echo "bbb" | go run . "a+" > /dev/null 2>&1
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
echo "Day 8 新增功能:"
echo "  ✓ ?: 零或一次"
echo "  ✓ +: 一或多次"
echo "  ✓ *: 零或多次"
echo "  ✓ 贪婪匹配（默认）"
echo "  ✓ 量词与字符类组合"
echo "  ✓ 量词与锚点组合"
echo "  ✓ 回溯机制"
