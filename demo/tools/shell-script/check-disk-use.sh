#!/bin/bash

# 设置磁盘使用率的阈值（以百分比表示）
THRESHOLD=90

# 获取当前日期和时间
DATE=$(date +"%Y-%m-%d %H:%M:%S")

# 获取磁盘使用率
DISK_USAGE=$(df -h | awk '$NF=="/"{print $5}' | sed 's/%//')

# 检查磁盘使用率是否超过阈值
if [ $DISK_USAGE -gt $THRESHOLD ]; then
    # 发送警报
    echo "磁盘使用率超过阈值：$DISK_USAGE%" | mail -s "磁盘使用率警报" 1006746222@qq.com
    echo "[$DATE] 磁盘使用率超过阈值：$DISK_USAGE%"
else
    echo "[$DATE] 磁盘使用率正常：$DISK_USAGE%"