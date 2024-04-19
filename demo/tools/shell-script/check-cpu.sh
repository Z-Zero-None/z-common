#!/bin/bash
#监控内存使用率——使用free命令获取

#main
#获取内存总量
#获取Mem数据，进行” “保留，切除多余的空格，获取第一个剪切数据
memory_totle=$(free -m | grep -i "mem" | tr -s " " | cut -d " " -f2)
swap_totle=$(free -m | grep -i "swap" | tr -s " " | cut -d " " -f2)
#获取内存的使用量
memory_use=$(free -m | grep -i "mem" | tr -s " " | cut -d " " -f3)
swap_use=$(free -m | grep -i "swap" | tr -s " " | cut -d " " -f3)
#获取buffer占用量
buffer_totle=$(free -m | grep -i "mem" | tr -s " " | cut -d " " -f6)

#计算输出（注意浮点数运算）
#内存使用率= memory_use*100/memory_totle
echo "memory_percentage:$(echo "scale=2;$memory_use*100/$memory_totle" | bc)% and buffer:$buffer_totle MB"
#区别判断是否为zero
if [ $swap_totle -eq 0 ]; then
    echo "swap is zero!!!"
else
    echo "swap_percentage:$(echo "scale=2;$swap_use*100/$swap_totle" | bc)%"
fi
#计算方式二
#cat /proc/meminfo 查看内存文件
memory_used() {
    #获取文件前两行数据，取第一行的第二个参数赋值给t
    mem_use=$(head -2 /proc/meminfo | awk 'NR==1{t=$2}NR==2{f=$2;print(t-f)*100/t}')
    mem_cache=$(head -5 /proc/meminfo | awk 'NR==1{t=$2}NR==5{f=$2;print(t-f)*100/t}')
    mem_buffer=$(head -4 /proc/meminfo | awk 'NR==1{t=$2}NR==4{f=$2;print(t-f)*100/t}')

    echo "mem_use:$mem_use%"
    echo "mem_cache:$mem_cache%"
    echo "mem_buffer:$mem_buffer%"
}

memory_used
