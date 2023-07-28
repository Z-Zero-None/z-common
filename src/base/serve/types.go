package serve

type MethodType uint8
type ServerType uint8

type HandlerInfo struct {
	Method             MethodType
	Path               string
	Handler            interface{}
	MiddlewareHandlers []interface{}
	RpcRegister        interface{}
}

const (
	_ = iota
	TYPE_GET
	TYPE_POST
	TYPE_PUT
	TYPE_DELETE
)

const (
	_ = iota
	TYPE_HTTP
	TYPE_RPC
)

type RequestInfo struct {
	Host           string
	ClientIP       string
	Method         string
	Path           string
	RequestContent string
}

type ResponseMap map[string]interface{}
