#!/bin/bash

# 第三天课程测试脚本

echo "======================================"
echo "    GeeGrep Day 3 - 测试脚本"
echo "    命令行选项：-n, -i, -v, -c"
echo "======================================"
echo

# 测试 1: -n 显示行号
echo "测试 1: -n 显示行号"
echo "--------------------------------------"
echo "输入: apple, banana, apricot"
echo "模式: 'ap'"
echo -e "apple\nbanana\napricot" | go run . -n "ap"
echo
echo "期望输出: 1:apple 和 3:apricot"
echo

# 测试 2: -i 忽略大小写
echo "测试 2: -i 忽略大小写"
echo "--------------------------------------"
echo "输入: Hello, WORLD, hello"
echo "模式: 'hello' (忽略大小写)"
echo -e "Hello\nWORLD\nhello" | go run . -i "hello"
echo
echo "期望输出: Hello 和 hello"
echo

# 测试 3: -v 反向匹配
echo "测试 3: -v 反向匹配"
echo "--------------------------------------"
echo "输入: error, success, error"
echo "模式: 'error' (反向匹配)"
echo -e "error\nsuccess\nerror" | go run . -v "error"
echo
echo "期望输出: success"
echo

# 测试 4: -c 统计数量
echo "测试 4: -c 统计数量"
echo "--------------------------------------"
echo "输入: cat, dog, cat, bird"
echo "模式: 'cat'"
echo -e "cat\ndog\ncat\nbird" | go run . -c "cat"
echo
echo "期望输出: 2"
echo

# 测试 5: -i + -n 组合
echo "测试 5: -i + -n 组合"
echo "--------------------------------------"
echo "输入: Hello, world, HELLO"
echo "模式: 'hello' (忽略大小写 + 显示行号)"
echo -e "Hello\nworld\nHELLO" | go run . -i -n "hello"
echo
echo "期望输出: 1:Hello 和 3:HELLO"
echo

# 测试 6: -c + -v 组合
echo "测试 6: -c + -v 组合"
echo "--------------------------------------"
echo "输入: a, b, c, a"
echo "模式: 'a' (反向匹配 + 统计)"
echo -e "a\nb\nc\na" | go run . -c -v "a"
echo
echo "期望输出: 2 (不包含 a 的行数)"
echo

# 测试 7: -i + 正则表达式
echo "测试 7: -i + 正则表达式"
echo "--------------------------------------"
echo "输入: Cat, cot, CUT"
echo "模式: 'c.t' (忽略大小写)"
echo -e "Cat\ncot\nCUT" | go run . -i "c.t"
echo
echo "期望输出: 所有三行"
echo

# 测试 8: -v + -n 组合
echo "测试 8: -v + -n 组合"
echo "--------------------------------------"
echo "输入: error, info, warning, error"
echo "模式: 'error' (反向 + 行号)"
echo -e "error\ninfo\nwarning\nerror" | go run . -v -n "error"
echo
echo "期望输出: 2:info 和 3:warning"
echo

# 测试 9: 从文件读取 + -c
echo "测试 9: 从文件读取 + -c"
echo "--------------------------------------"
echo "创建测试文件..."
echo -e "hello\nworld\nhello" > /tmp/test1.txt
echo -e "foo\nbar" > /tmp/test2.txt
echo "搜索 'hello' 并统计"
go run . -c "hello" /tmp/test1.txt /tmp/test2.txt
echo
echo "期望输出:"
echo "  /tmp/test1.txt:2"
echo "  /tmp/test2.txt:0"
rm /tmp/test1.txt /tmp/test2.txt
echo

# 测试 10: 从文件读取 + -n
echo "测试 10: 从文件读取 + -n"
echo "--------------------------------------"
echo -e "apple\nbanana\napricot" > /tmp/test-day3.txt
echo "文件内容:"
cat /tmp/test-day3.txt
echo
echo "搜索 'ap' 并显示行号"
go run . -n "ap" /tmp/test-day3.txt
rm /tmp/test-day3.txt
echo

# 测试 11: 通配符 + -i
echo "测试 11: 通配符 + -i"
echo "--------------------------------------"
echo "输入: Hello, Hallo, HXLLO"
echo "模式: 'h.llo' (忽略大小写)"
echo -e "Hello\nHallo\nHXLLO" | go run . -i "h.llo"
echo
echo "期望输出: 所有三行"
echo

# 测试 12: 退出码检查
echo "测试 12: 退出码检查"
echo "--------------------------------------"
echo "测试无匹配的退出码:"
echo "foo" | go run . "bar" > /dev/null 2>&1
if [ $? -eq 1 ]; then
    echo "✓ 无匹配: 退出码 1"
else
    echo "✗ 错误的退出码"
fi

echo "测试有匹配的退出码:"
echo "hello" | go run . "hello" > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "✓ 有匹配: 退出码 0"
else
    echo "✗ 错误的退出码"
fi

echo "测试 -v 反向匹配的退出码:"
echo "foo" | go run . -v "bar" > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "✓ 反向匹配成功: 退出码 0"
else
    echo "✗ 错误的退出码"
fi
echo

echo "======================================"
echo "    测试完成！"
echo "======================================"
echo
echo "Day 3 新增功能:"
echo "  ✓ -n: 显示行号"
echo "  ✓ -i: 忽略大小写"
echo "  ✓ -v: 反向匹配"
echo "  ✓ -c: 统计数量"
echo "  ✓ 选项可以组合使用"
echo "  ✓ 多文件支持"
