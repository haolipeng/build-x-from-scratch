#!/bin/bash

# 第四天课程测试脚本

echo "======================================"
echo "    GeeGrep Day 4 - 测试脚本"
echo "    字符类：\\d, \\w, \\s 和转义"
echo "======================================"
echo

# 测试 1: \d 匹配数字
echo "测试 1: \\d 匹配数字"
echo "--------------------------------------"
echo "输入: abc, 123, a1b"
echo -e "abc\n123\na1b" | go run . "\\d"
echo
echo "期望输出: 123 和 a1b (包含数字的行)"
echo

# 测试 2: \d 匹配电话号码模式
echo "测试 2: \\d 匹配电话号码"
echo "--------------------------------------"
echo "Call: 555-1234" | go run . "\\d\\d\\d-\\d\\d\\d\\d"
echo
echo "期望输出: Call: 555-1234"
echo

# 测试 3: \w 匹配单词字符
echo "测试 3: \\w 匹配单词字符"
echo "--------------------------------------"
echo "输入: hello, hello_world, hello-world"
echo -e "hello\nhello_world\nhello-world" | go run . "hello\\wworld"
echo
echo "期望输出: hello_world (下划线是单词字符)"
echo

# 测试 4: \w 匹配变量名
echo "测试 4: \\w 匹配变量名"
echo "--------------------------------------"
echo "var_123" | go run . "\\w\\w\\w_\\d\\d\\d"
echo
echo "期望输出: var_123"
echo

# 测试 5: \s 匹配空白
echo "测试 5: \\s 匹配空白"
echo "--------------------------------------"
echo -e "hello world\nhelloworld" | go run . "hello\\sworld"
echo
echo "期望输出: hello world"
echo

# 测试 6: \s 匹配制表符
echo "测试 6: \\s 匹配制表符"
echo "--------------------------------------"
printf "tab\there\nnotab" | go run . "tab\\shere"
echo
echo "期望输出: tab	here"
echo

# 测试 7: 字面转义 \.
echo "测试 7: 字面转义 \\."
echo "--------------------------------------"
echo "输入: 3.14, 3a14"
echo -e "3.14\n3a14" | go run . "3\\.14"
echo
echo "期望输出: 3.14 (只匹配字面点号)"
echo

# 测试 8: 字面转义 \\
echo "测试 8: 字面转义 \\\\"
echo "--------------------------------------"
echo 'path\to\file' | go run . "path\\\\to"
echo
echo "期望输出: path\\to\\file"
echo

# 测试 9: \D 匹配非数字
echo "测试 9: \\D 匹配非数字"
echo "--------------------------------------"
echo -e "123\nabc\n1a2" | go run . "\\D"
echo
echo "期望输出: abc 和 1a2 (包含非数字字符)"
echo

# 测试 10: \W 匹配非单词字符
echo "测试 10: \\W 匹配非单词字符"
echo "--------------------------------------"
echo -e "hello\nhello@world\nabc_123" | go run . "\\W"
echo
echo "期望输出: hello@world (包含 @)"
echo

# 测试 11: \S 匹配非空白
echo "测试 11: \\S 匹配非空白"
echo "--------------------------------------"
echo -e "   \nhello\n\t" | go run . "\\S"
echo
echo "期望输出: hello"
echo

# 测试 12: 混合使用字符类
echo "测试 12: 混合使用字符类"
echo "--------------------------------------"
echo "Date: 12/31/2024" | go run . "\\d\\d/\\d\\d/\\d\\d\\d\\d"
echo
echo "期望输出: Date: 12/31/2024"
echo

# 测试 13: 字符类 + 通配符
echo "测试 13: 字符类 + 通配符"
echo "--------------------------------------"
echo "user@example.com" | go run . "\\w\\w\\w\\w.\\w\\w\\w\\w\\w\\w\\w\\.\\w\\w\\w"
echo
echo "期望输出: user@example.com"
echo

# 测试 14: 字符类 + 选项 -n
echo "测试 14: 字符类 + -n 选项"
echo "--------------------------------------"
echo -e "no digits\nline 123\nmore text" | go run . -n "\\d"
echo
echo "期望输出: 2:line 123"
echo

# 测试 15: 字符类 + 选项 -i
echo "测试 15: 字符类 + -i 选项"
echo "--------------------------------------"
echo -e "ABC\n123\nxyz" | go run . -i "\\w\\w\\w"
echo
echo "期望输出: 所有三行"
echo

# 测试 16: 字符类 + 选项 -v
echo "测试 16: 字符类 + -v 选项"
echo "--------------------------------------"
echo -e "abc\n123\nxyz" | go run . -v "\\d"
echo
echo "期望输出: abc 和 xyz (不包含数字)"
echo

# 测试 17: 字符类 + 选项 -c
echo "测试 17: 字符类 + -c 选项"
echo "--------------------------------------"
echo -e "a1\nb2\nc3\nde" | go run . -c "\\d"
echo
echo "期望输出: 3"
echo

# 测试 18: 复杂模式
echo "测试 18: 复杂模式 - IPv4 地址"
echo "--------------------------------------"
echo "IP: 192.168.1.1" | go run . "\\d\\d\\d\\.\\d\\d\\d\\.\\d\\.\\d"
echo
echo "期望输出: IP: 192.168.1.1"
echo

# 测试 19: 文件扩展名匹配
echo "测试 19: 文件扩展名匹配"
echo "--------------------------------------"
echo -e "file.txt\nfile.doc\nfiletxt" | go run . "file\\.\\w\\w\\w"
echo
echo "期望输出: file.txt 和 file.doc"
echo

# 测试 20: 退出码检查
echo "测试 20: 退出码检查"
echo "--------------------------------------"
echo "测试包含数字的退出码:"
echo "123" | go run . "\\d" > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "✓ 找到匹配: 退出码 0"
else
    echo "✗ 错误的退出码"
fi

echo "测试不包含数字的退出码:"
echo "abc" | go run . "\\d" > /dev/null 2>&1
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
echo "Day 4 新增功能:"
echo "  ✓ \\d: 匹配数字 [0-9]"
echo "  ✓ \\w: 匹配单词字符 [a-zA-Z0-9_]"
echo "  ✓ \\s: 匹配空白字符"
echo "  ✓ \\D, \\W, \\S: 对应的否定形式"
echo "  ✓ 字面转义: \\., \\\\, \\*, 等"
echo "  ✓ 字符类与选项组合使用"
