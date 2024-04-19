#!/bin/bash
#监控使用cpu前十的进程
#通过ps和top命令进行查看

#统计内存
memory(){
    #创建临时文件存放信息
    temp_file=`mktemp memeory.XXX`
    top -b -n1 >$temp_file
    echo "mem:"
    #查看文件，按进程统计内存使用大小
    tail -n +8 $temp_file | awk '{array[$NF]+=$6}END{for (i in array) print array[i],i}'| sort -k 1 -n -r | head -10
    #针对ps吻技安处理分析
    # tail -n +2 $temp_file | awk '{array[$11]+=$6}END{for (i in array) print i,array[i]}'| sort -k 2 -n -r | head -10

    rm -f $temp_file
}

# 统计CPU
cpu(){
     #创建临时文件存放信息
    temp_file=`mktemp cpu.XXX`
    ps aux >$temp_file
    echo "cpu:"
    #查看文件，按进程统计内存使用大小 获取最后一列和第三列的为数组key与value 遍历排序打印
    tail -n +2 $temp_file | awk '{array[$NF]+=$3}END{for (i in array) print array[i],i}'| sort -k 1 -n -r | head -10

    rm -f $temp_file
}

memory
cpu