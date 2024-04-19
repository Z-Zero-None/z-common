#!/bin/bash
#监控io队列？磁盘信息
#iostat
#   -c cpu信息
#   -d 设备信息
#   -x 详细信息

#判断命令是否存在 安装命令apt install sysstat
[ ! -x /usr/bin/iostat ] && echo "iostat command is not found!" && exit 1

#main aqu_sz对列长度所在列
io() {
    device_num=$(iostat -x | egrep "^sd[a-z]" | wc -l)
    iostat -x 1 3 | egrep "^sd[a-z]" | tail -n +$((device_num + 1)) | awk '{io_long[$1]+=$19}END{for(i in io_long)print io_long[i],i}'
}

#阈值判断
# while true; do
    io
#     sleep 5
# done
