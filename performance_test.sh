#!/bin/bash

echo "=== Log Parser 性能测试 ==="
echo

# 测试基本功能
echo "1. 测试基本功能..."
echo "2024-01-01 12:00:00 [INFO] 测试日志" | ./logparser
echo

# 测试不同格式
echo "2. 测试不同日志格式..."
echo -e "2024-01-01 12:00:00 [INFO] 标准格式\n2024-01-01T12:00:00Z [ERROR] ISO格式\n[DEBUG] 简单格式\n未知格式日志" | ./logparser
echo

# 测试大文件性能
echo "3. 测试大文件性能..."
echo "生成测试日志文件..."
for i in {1..1000}; do
    echo "2024-01-01 12:00:0$((i%60)) [INFO] 这是第 $i 条测试日志消息，包含一些随机内容来测试性能"
done > test_large.log

echo "处理 1000 行日志..."
time cat test_large.log | ./logparser > /dev/null
echo

# 测试带性能分析的版本
echo "4. 测试带性能分析的版本..."
echo "处理 100 行日志（带性能分析）..."
time head -100 test_large.log | ./logparser-profiler -profile -stats > /dev/null
echo

# 清理测试文件
echo "5. 清理测试文件..."
rm -f test_large.log
echo "测试完成！"
