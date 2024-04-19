#!/bin/bash
#用脚本判断远程主机是否存活
#监控方法ping 底层原理ICMP
#禁用非白名单IP请求

#满足条件
#网络延迟 导致ping失败出现假告警——多次尝试失败，才发出告警
#ping的频率设置

#main
for ((i = 0; i < 4; i++)); do
    #test code
    echo ping $1
    echo ping_count"$i"
    if ping -c1 $1 &>/dev/null; then
        export ping_count"$i"=1
    else
        export ping_count"$i"=0
    fi
    #时间间隔
    sleep 1
done

#3次ping失败告警
if [ $ping_count1 -eq $ping_count2 ] && [ $ping_count2 -eq $ping_count3 ] && [ $ping_count1 -eq 0 ]; then
    echo "$1 is down!"
elif [ $ping_count1 -eq $ping_count2 ] && [ $ping_count2 -eq $ping_count3 ] && [ $ping_count1 -eq 1 ]; then
    echo "$1 is up!"
else
    echo "$1 network anomaly!"
fi
unset ping_count[1,2,3]
