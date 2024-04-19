#!/bin/bash
#监控一个服务端口,案例监控命令

#监控方法
#1.通过查询服务启动状态——systemctl status xxx/service xxx status
#2.查看端口是否存在——lsof -i :port
#3.查看进程是否存在——ps aux process
#无法排查服务假死、压力过大无法响应等问题
#4.通过测试端口是否响应——telnet

#main
temp_file=$(mktemp port_status.XXX)

#判断命令是否存在
[ ! -x /usr/bin/telnet ] && echo "telnet command is not found!" && exit 1

#test code $1为ip地址，$2为端口号
(
    telnet $1 $2 <<EOF
quit
EOF
) &>$temp_file

#分析文件内容，判断结果
if egrep "\^]" $temp_file &>/dev/null; then
    echo "$1 $2 is open"
else
    echo "$1 $2 is close"
fi

rm -f $temp_file
