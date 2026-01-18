#!/bin/bash

# GeeGrep 演进对比：Day 1 → Day 2 → Day 3

echo "=============================================="
echo "  GeeGrep 功能演进对比"
echo "  Day 1 → Day 2 → Day 3"
echo "=============================================="
echo

echo "测试数据: cat, Cat, cot"
echo

# Day 1: 基础字符串匹配
echo "----------------------------------------"
echo "Day 1: 基础字符串匹配"
echo "----------------------------------------"
echo "模式: 'cat'"
echo -e "cat\nCat\ncot" | (cd day1-project-init && go run . "cat" 2>/dev/null)
echo "✓ 只匹配精确的 'cat'"
echo

# Day 2: 正则表达式
echo "----------------------------------------"
echo "Day 2: 正则表达式基础"
echo "----------------------------------------"
echo "模式: 'c.t' (通配符)"
echo -e "cat\nCat\ncot" | (cd day2-regex-basics && go run . "c.t" 2>/dev/null)
echo "✓ 通配符 . 匹配任意字符"
echo "✗ 但区分大小写，不匹配 'Cat'"
echo

# Day 3: 命令行选项
echo "----------------------------------------"
echo "Day 3: 命令行选项"
echo "----------------------------------------"
echo "模式: 'c.t' + -i (忽略大小写)"
echo -e "cat\nCat\ncot" | (cd day3-command-options && go run . -i "c.t" 2>/dev/null)
echo "✓ 通配符 + 忽略大小写"
echo "✓ 匹配所有三行"
echo

echo "=============================================="
echo "功能对比表"
echo "=============================================="
echo

cat << 'EOF'
| 功能                | Day 1 | Day 2 | Day 3 |
|---------------------|-------|-------|-------|
| 字面字符匹配        | ✓     | ✓     | ✓     |
| 通配符 .           | ✗     | ✓     | ✓     |
| 正则表达式          | ✗     | ✓     | ✓     |
| -n 显示行号        | ✗     | ✗     | ✓     |
| -i 忽略大小写      | ✗     | ✗     | ✓     |
| -v 反向匹配        | ✗     | ✗     | ✓     |
| -c 统计数量        | ✗     | ✗     | ✓     |
| 文件搜索            | ✓     | ✓     | ✓     |
| 标准输入            | ✓     | ✓     | ✓     |
EOF

echo
echo "=============================================="
echo "实际演示"
echo "=============================================="
echo

echo "场景: 在日志中查找错误，忽略大小写，显示行号"
echo

# 创建测试日志
echo -e "INFO: Application started\nERROR: Connection failed\ninfo: Processing data\nerror: Timeout" > /tmp/test.log

echo "测试日志内容:"
cat -n /tmp/test.log
echo

echo "Day 1 搜索 'error' (精确匹配):"
(cd day1-project-init && go run . "error" /tmp/test.log 2>/dev/null)
echo "→ 只找到小写的 'error'"
echo

echo "Day 2 搜索 'error' (正则表达式):"
(cd day2-regex-basics && go run . "error" /tmp/test.log 2>/dev/null)
echo "→ 还是只找到小写的 'error'"
echo

echo "Day 3 搜索 'error' -i -n (忽略大小写 + 行号):"
(cd day3-command-options && go run . -i -n "error" /tmp/test.log 2>/dev/null)
echo "→ 找到所有错误，并显示行号！"
echo

rm /tmp/test.log

echo "=============================================="
echo "总结"
echo "=============================================="
echo
echo "Day 1: 简单实用"
echo "  • 适合精确字符串搜索"
echo "  • 实现简单，性能好"
echo "  • 功能有限"
echo
echo "Day 2: 强大灵活"
echo "  • 支持正则表达式模式匹配"
echo "  • 可扩展性强"
echo "  • 为后续功能打下基础"
echo
echo "Day 3: 实用完整"
echo "  • 添加常用命令行选项"
echo "  • 更接近真实的 grep 工具"
echo "  • 可以处理实际的日常任务"
echo
