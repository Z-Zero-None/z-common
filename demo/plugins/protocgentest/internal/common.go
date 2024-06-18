package internal

import (
	"fmt"
	"path/filepath"
	"strings"
)

var (
	DefaultClientName       = "GRPCClient"
	DefaultImportPath       = "luya/transports/grpc"
	WireImportPath          = "github.com/google/wire"
	LocalClientsImportPaths = []string{
		"context",
		"google.golang.org/grpc",
	}
	LocalClientsPTypesImportPath = "luya/api/ptypes"
	KeyImportPath                = "api/ptypes/type.proto"
)

const (
	FMTProtoImportPath       = "luya/api/%s/proto"
	FMTGRPCServersImportPath = "luya/apps/%s/handlers/grpcservers"
	FMTGRPCServersFilePath   = "./apps/%s/handlers/grpcservers"
)

func xxClient(xx string) string {
	return fmt.Sprintf("%sClient", xx)
}

func xxServer(xx string) string {
	return fmt.Sprintf("%sServer", xx)
}

func xxWithoutRPC(xx string) string {
	xx = strings.Replace(xx, "RPC", "", 1)
	return xx
}

func newXXClient(xx string) string {
	return fmt.Sprintf("New%sClient", xx)
}

func newXXWithoutRPC(xx string) string {
	xx = strings.Replace(xx, "RPC", "", 1)
	return fmt.Sprintf("New%s", xx)
}

func BaseNameWithoutExtension(path string) string {
	base := filepath.Base(path)          // 获取文件名，例如 "xx.json"
	ext := filepath.Ext(base)            // 获取扩展名，例如 ".json"
	return strings.TrimSuffix(base, ext) // 去除扩展名，得到 "xx"
}
