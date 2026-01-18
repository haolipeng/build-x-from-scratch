#!/bin/bash

# 第二天课程测试脚本

echo "======================================"
echo "    GeeGrep Day 2 - 测试脚本"
echo "   正则表达式基础：字面字符 + 通配符"
echo "======================================"
echo

# 测试 1: 字面字符匹配
echo "测试 1: 字面字符匹配"
echo "--------------------------------------"
echo "模式: 'hello'"
echo -e "hello world\ngoodbye\nhello again" | go run . "hello"
echo
echo "期望输出: 包含 'hello' 的两行"
echo

# 测试 2: 通配符 . 匹配任意单个字符
echo "测试 2: 通配符匹配 - 'c.t'"
echo "--------------------------------------"
echo "输入: cat, cot, cut, coat"
echo -e "cat\ncot\ncut\ncoat" | go run . "c.t"
echo
echo "期望输出: cat, cot, cut (coat 有4个字符，但包含 'coa' 所以也匹配)"
echo

# 测试 3: 多个通配符
echo "测试 3: 多个通配符 - 'h...o'"
echo "--------------------------------------"
echo -e "hello\nhallo\nhalo\nhellboy" | go run . "h...o"
echo
echo "期望输出: hello, hallo (5个字符的匹配)"
echo

# 测试 4: 纯通配符模式
echo "测试 4: 纯通配符 - '...'"
echo "--------------------------------------"
echo -e "ab\nabc\nabcd\nxy" | go run . "..."
echo
echo "期望输出: abc, abcd (至少3个字符的行)"
echo

# 测试 5: 混合模式
echo "测试 5: 混合模式 - 'a.c'"
echo "--------------------------------------"
echo -e "abc\nabc123\naxc\na  c\nadc" | go run . "a.c"
echo
echo "期望输出: 所有包含 'a[任意字符]c' 的行"
echo

# 测试 6: 从文件读取
echo "测试 6: 从文件读取"
echo "--------------------------------------"
echo -e "pattern123test\npattern456test\npatternXYZtest" > /tmp/test-day2.txt
echo "文件内容:"
cat /tmp/test-day2.txt
echo
echo "搜索模式: 'pattern...test'"
go run . "pattern...test" /tmp/test-day2.txt
echo
echo "期望输出: 所有行（都匹配 pattern[3个字符]test）"
rm /tmp/test-day2.txt
echo

# 测试 7: 特殊字符
echo "测试 7: 特殊字符"
echo "--------------------------------------"
echo -e "hello@world\nhello#world\nhello world" | go run . "hello.world"
echo
echo "期望输出: 所有行（. 匹配 @, #, 空格）"
echo

# 测试 8: 行首匹配
echo "测试 8: 包含匹配（grep 默认行为）"
echo "--------------------------------------"
echo -e "prefix_cat_suffix\ndog\ncat" | go run . "cat"
echo
echo "期望输出: prefix_cat_suffix 和 cat（包含匹配）"
echo

# 测试 9: 无匹配情况
echo "测试 9: 无匹配（检查退出码）"
echo "--------------------------------------"
echo "foo" | go run . "bar" > /dev/null 2>&1
exitcode=$?
if [ $exitcode -eq 1 ]; then
    echo "✓ 正确：未找到匹配，退出码为 1"
else
    echo "✗ 错误：期望退出码 1，实际为 $exitcode"
fi
echo

# 测试 10: 有匹配的退出码
echo "测试 10: 有匹配（检查退出码）"
echo "--------------------------------------"
echo "hello" | go run . "h...o" > /dev/null 2>&1
exitcode=$?
if [ $exitcode -eq 0 ]; then
    echo "✓ 正确：找到匹配，退出码为 0"
else
    echo "✗ 错误：期望退出码 0，实际为 $exitcode"
fi
echo

# 测试 11: 空模式检查
echo "测试 11: 空模式错误处理"
echo "--------------------------------------"
echo "test" | go run . "" 2>&1 | head -1
echo "期望输出: 错误信息"
echo

# 测试 12: 对比 Day 1 和 Day 2
echo "测试 12: Day 1 vs Day 2 比较"
echo "--------------------------------------"
echo "输入: 'hello', 模式: 'h.llo'"
echo "hello" | go run . "h.llo"
echo "Day 1（字面匹配）无法匹配 'h.llo'"
echo "Day 2（正则匹配）可以匹配！"
echo

echo "======================================"
echo "    测试完成！"
echo "======================================"
echo
echo "总结:"
echo "✓ 支持字面字符精确匹配"
echo "✓ 支持通配符 . 匹配任意字符"
echo "✓ 支持混合使用字面字符和通配符"
echo "✓ 正确的退出码（0=匹配，1=无匹配，2=错误）"
