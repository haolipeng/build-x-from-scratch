#!/bin/bash

# 第一天课程测试脚本

echo "======================================"
echo "    GeeGrep Day 1 - 测试脚本"
echo "======================================"
echo

# 测试 1: 从标准输入搜索
echo "测试 1: 从标准输入搜索 'hello'"
echo "--------------------------------------"
echo -e "hello world\nfoo bar\nhello again" | go run . "hello"
echo
echo "期望输出: hello world 和 hello again"
echo

# 测试 2: 搜索包含字符的行
echo "测试 2: 从标准输入搜索 'ap'"
echo "--------------------------------------"
echo -e "apple pie\nbanana\napricot\ngrape" | go run . "ap"
echo
echo "期望输出: apple pie, apricot, grape"
echo

# 测试 3: 从文件搜索
echo "测试 3: 从文件搜索"
echo "--------------------------------------"
echo -e "hello world\ntest line\nhello again" > /tmp/test-day1.txt
echo "创建测试文件: /tmp/test-day1.txt"
echo "文件内容:"
cat /tmp/test-day1.txt
echo
echo "搜索结果:"
go run . "hello" /tmp/test-day1.txt
echo
echo "期望输出: hello world 和 hello again"
echo

# 测试 4: 多文件搜索
echo "测试 4: 多文件搜索"
echo "--------------------------------------"
echo "match in file 1" > /tmp/file1.txt
echo "match in file 2" > /tmp/file2.txt
echo "no pattern here" > /tmp/file3.txt
echo "创建 3 个测试文件"
echo
echo "搜索结果:"
go run . "match" /tmp/file1.txt /tmp/file2.txt /tmp/file3.txt
echo
echo "期望输出: match in file 1 和 match in file 2"
echo

# 测试 5: 无匹配
echo "测试 5: 无匹配（检查退出码）"
echo "--------------------------------------"
echo "foo" | go run . "bar" > /dev/null 2>&1
exitcode=$?
echo "退出码: $exitcode"
if [ $exitcode -eq 1 ]; then
    echo "✓ 正确：未找到匹配，退出码为 1"
else
    echo "✗ 错误：期望退出码 1，实际为 $exitcode"
fi
echo

# 测试 6: 有匹配的退出码
echo "测试 6: 有匹配（检查退出码）"
echo "--------------------------------------"
echo "hello" | go run . "hello" > /dev/null 2>&1
exitcode=$?
echo "退出码: $exitcode"
if [ $exitcode -eq 0 ]; then
    echo "✓ 正确：找到匹配，退出码为 0"
else
    echo "✗ 错误：期望退出码 0，实际为 $exitcode"
fi
echo

# 清理临时文件
rm -f /tmp/test-day1.txt /tmp/file1.txt /tmp/file2.txt /tmp/file3.txt

echo "======================================"
echo "    测试完成！"
echo "======================================"
