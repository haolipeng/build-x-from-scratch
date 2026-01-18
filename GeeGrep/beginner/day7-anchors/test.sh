#!/bin/bash

# 第七天课程测试脚本

echo "======================================"
echo "    GeeGrep Day 7 - 测试脚本"
echo "    锚点：^, $, \\b, \\B"
echo "======================================"
echo

# 测试 1: 行首锚点 ^
echo "测试 1: 行首锚点 ^hello"
echo "--------------------------------------"
echo "输入: hello world, world hello, hello"
echo -e "hello world\nworld hello\nhello" | go run . "^hello"
echo
echo "期望输出: hello world 和 hello"
echo

# 测试 2: 行尾锚点 $
echo "测试 2: 行尾锚点 world$"
echo "--------------------------------------"
echo -e "hello world\nworld hello\nworld" | go run . "world$"
echo
echo "期望输出: hello world 和 world"
echo

# 测试 3: 同时使用 ^ 和 $
echo "测试 3: ^hello$ 完整行匹配"
echo "--------------------------------------"
echo -e "hello\nhello world\nworld" | go run . "^hello$"
echo
echo "期望输出: hello"
echo

# 测试 4: ^..$ 匹配两个字符的行
echo "测试 4: ^..$ 匹配两个字符的行"
echo "--------------------------------------"
echo -e "a\nab\nabc" | go run . "^..$"
echo
echo "期望输出: ab"
echo

# 测试 5: 单词边界 \b
echo "测试 5: 单词边界 \\bthe\\b"
echo "--------------------------------------"
echo -e "the cat\nthere\nother" | go run . "\\bthe\\b"
echo
echo "期望输出: the cat"
echo

# 测试 6: 单词边界匹配独立的 cat
echo "测试 6: 单词边界 \\bcat\\b"
echo "--------------------------------------"
echo -e "cat\ncatalog\ncat food" | go run . "\\bcat\\b"
echo
echo "期望输出: cat 和 cat food"
echo

# 测试 7: 非单词边界 \B
echo "测试 7: 非单词边界 \\Bthe"
echo "--------------------------------------"
echo -e "there\nthe\nother" | go run . "\\Bthe"
echo
echo "期望输出: other (匹配 'other' 中的 'the')"
echo

# 测试 8: 非单词边界 cat\B
echo "测试 8: 非单词边界 cat\\B"
echo "--------------------------------------"
echo -e "catalog\ncat\nconcat" | go run . "cat\\B"
echo
echo "期望输出: catalog"
echo

# 测试 9: 空行匹配 ^$
echo "测试 9: 空行匹配 ^$"
echo "--------------------------------------"
printf "hello\n\nworld\n" | go run . "^$"
echo "(空行已匹配)"
echo
echo "期望输出: 空行"
echo

# 测试 10: 行首 + 字符类
echo "测试 10: 行首 + 字符类 ^[a-z]"
echo "--------------------------------------"
echo -e "hello\n123\nWorld" | go run . "^[a-z]"
echo
echo "期望输出: hello"
echo

# 测试 11: 行尾 + 数字
echo "测试 11: 行尾 + 数字 \\d$"
echo "--------------------------------------"
echo -e "abc\nabc123\n123abc" | go run . "\\d$"
echo
echo "期望输出: abc123"
echo

# 测试 12: 与 -n 选项组合
echo "测试 12: ^hello + -n 选项"
echo "--------------------------------------"
echo -e "hello\nworld\nhello world" | go run . -n "^hello"
echo
echo "期望输出: 1:hello 和 3:hello world"
echo

# 测试 13: 与 -v 选项组合
echo "测试 13: world$ + -v 选项"
echo "--------------------------------------"
echo -e "hello\nworld\nhello world" | go run . -v "world$"
echo
echo "期望输出: hello"
echo

# 测试 14: 与 -c 选项组合
echo "测试 14: ^hello + -c 选项"
echo "--------------------------------------"
echo -e "hello\nhello world\nworld hello" | go run . -c "^hello"
echo
echo "期望输出: 2"
echo

# 测试 15: 复杂模式 - 以数字开头以字母结尾
echo "测试 15: ^\\d.*[a-z]$ 以数字开头以字母结尾"
echo "--------------------------------------"
echo -e "1abc\nabc1\n123\n1a" | go run . "^\\d.*[a-z]$"
echo
echo "期望输出: 1abc 和 1a (注: .* 待后续实现)"
echo

# 测试 16: 退出码检查
echo "测试 16: 退出码检查"
echo "--------------------------------------"
echo "测试匹配的退出码:"
echo "hello" | go run . "^hello" > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "✓ 找到匹配: 退出码 0"
else
    echo "✗ 错误的退出码"
fi

echo "测试不匹配的退出码:"
echo "world" | go run . "^hello" > /dev/null 2>&1
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
echo "Day 7 新增功能:"
echo "  ✓ ^: 行首锚点"
echo "  ✓ $: 行尾锚点"
echo "  ✓ \\b: 单词边界"
echo "  ✓ \\B: 非单词边界"
echo "  ✓ 锚点优化（^ 开头只在位置 0 匹配）"
echo "  ✓ 与命令行选项组合使用"
