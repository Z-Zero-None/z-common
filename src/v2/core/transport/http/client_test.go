package http

import (
	"fmt"
	"net/http"
	"testing"
)

type Params struct {
	Param1 string `json:"param1"`
	Param2 int    `json:"param2"`
}

func Test_Client(t *testing.T) {
	// 要请求的URL
	url := "http://test.act.pago.tv"
	method := "h5/ivt/invitation/rewards"

	newClient := NewClient(
		Host(url),
	)
	res := &Response{}
	responseFunc := RespFunc[*Response](res)
	if err := newClient.Request(http.MethodGet, method, DefaultParamsFunc, responseFunc); err != nil {
		fmt.Println("err:", err.Error())
	}

	fmt.Printf("res：%v\ncode:%d,msg:%s \n", res, res.Code, res.Msg)
}
