/*
实现ping操作
*/
package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

var (
	timeout int64
	size    int
	count   int
)

type ICMP struct {
	Type        uint8
	Code        uint8
	CheckSum    uint16
	ID          uint8
	SequenceNum uint16
}

func getCommandArgs() {
	flag.Int64Var(&timeout, "w", 1, "请求超时时长，单位秒")
	flag.IntVar(&size, "l", 32, "请求发送缓冲区大小，单位字节")
	flag.IntVar(&count, "n", 4, "发送请求数")
	flag.Parse()
}

func Ping() {
	//获取参数
	getCommandArgs()
	//获取传参ip
	desIp := os.Args[len(os.Args)-1]

	// desIp:= "www.baidu.com"

	//请求ip的icmp方法
	conn, err := net.DialTimeout("ip4:icmp", desIp, time.Duration(timeout)*time.Second)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()
	fmt.Printf("正在 Ping %s [%s] 具有 %d 字节的数据:\n", desIp, conn.RemoteAddr(), size)
	for i := 0; i < count; i++ {
		now := time.Now()
		data := make([]byte, size)
		var buf bytes.Buffer
		icmp := &ICMP{
			Type: 8, Code: 0, CheckSum: 0, ID: 1, SequenceNum: 1,
		}
		binary.Write(&buf, binary.BigEndian, icmp)
		buf.Write(data)
		data = buf.Bytes()

		checkSum := checkSum(data)
		data[2] = byte(checkSum >> 8)
		data[3] = byte(checkSum)
		conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
		_, err := conn.Write(data)
		if err != nil {
			log.Println(err)
			continue
		}
		b := make([]byte, 65535)
		n, err := conn.Read(b)
		if err != nil {
			log.Println(err)
			continue
		}
		ts := time.Since(now).Milliseconds()
		fmt.Printf("来自 %d.%d.%d.%d 的回复: 字节=%d 时间=%dms TTL=%d\n", b[12], b[13], b[14], b[15], n-27, ts, b[8])
	}
}

// ICMP校验
func checkSum(data []byte) uint16 {
	length := len(data)
	index := 0
	var sum uint32
	//将相邻两个字节拼接到一起组成16bit数，累加求和
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])<<8
		length -= 2
		index += 2
	}
	//若长度为奇数 累加
	if length != 0 {
		sum += uint32(data[index])
	}
	// 将值的高16位和低16位不断求和，知道高16位为0
	high := sum >> 16
	for high != 0 {
		sum = high + uint32(uint16(sum))
		high = sum >> 16
	}
	// 0取余
	return uint16(^sum)
}

// go run goping.go -w 1 -l 32 -n 8 www.baidu.com
func main() {
	Ping()
}
