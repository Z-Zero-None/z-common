package main

import (
	"flag"
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"luya/plugins/protocgentest/internal"
	"slices"
	"sort"
	"strings"
	"sync"
)

// window: go build -o $GOPATH/bin/protoc-gen-test.exe main.go
// linux: go build -o $GOPATH/bin/protoc-gen-test main.go
// gen_sh_test: protoc --test_out=. --test_opt=api_file_dir=api/auth,api_module_name=auth,api_import_prefix=luya,loc_client_type=RPC,loc_file_dir=apps/auth/handlers/localclients   api/auth/src/*.proto
// gen_sh: protoc --test_out=. --test_opt=api_file_dir=api/${dir_name},api_module_name=${dir_name},api_import_prefix=luya,loc_client_type=RPC,loc_file_dir=apps/${dir_name}/handlers/localclients api/${dir_name}/src/*.proto
func main() {
	var (
		flags           flag.FlagSet
		apiFileDir      = flags.String("api_file_dir", "", "file_dir/xx.go")
		apiModuleName   = flags.String("api_module_name", "name", "package name")
		apiImportPrefix = flags.String("api_import_prefix", "", "import")
		locClientType   = flags.String("loc_client_type", "", "RPC:xxxRPC,Svc:xxxSvc")
		locFileDir      = flags.String("loc_file_dir", "", "file_dir/xx.go")
		isOpen          = flags.Bool("is_open", true, "generate grpc_client switch")
	)
	//尝试使用插件生成
	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(p *protogen.Plugin) error {
		req := &internal.GrpcClientFileSource{
			PackageName:        *apiModuleName,
			ImportPaths:        []string{internal.DefaultImportPath},
			ClientName:         internal.DefaultClientName,
			ServiceNames:       []string{},
			GRPCClientFileName: "./grpc_clients.go",
		}
		wireReq := &internal.WireFileSource{
			PackageName:  *apiModuleName,
			ImportPaths:  []string{internal.WireImportPath},
			ServiceNames: []string{},
			WireFileName: "./wire.go",
		}
		if len(*apiFileDir) > 0 {
			req.GRPCClientFileName = fmt.Sprintf("./%s/grpc_clients.go", strings.Trim(*apiFileDir, "/"))
			wireReq.WireFileName = fmt.Sprintf("./%s/wire.go", strings.Trim(*apiFileDir, "/"))
		}
		once := sync.Once{}
		var serviceNames []string

		var localClients []*internal.LocalClientsFileSource
		for _, f := range p.FilesByPath {
			if !f.Generate {
				continue
			}
			client := &internal.LocalClientsFileSource{
				ModuleName:     *apiModuleName,
				FileDir:        strings.Trim(*locFileDir, "/"),
				Name:           f.GeneratedFilenamePrefix,
				BaseName:       internal.BaseNameWithoutExtension(f.GeneratedFilenamePrefix),
				ImportPaths:    internal.LocalClientsImportPaths,
				ServiceMethods: make(map[string][]*internal.MethodSource),
			}
			if slices.Contains(f.Proto.Dependency, internal.KeyImportPath) {
				client.ImportPaths = append(client.ImportPaths, internal.LocalClientsPTypesImportPath)
			}
			for _, s := range f.Services {
				serviceNames = append(serviceNames, s.GoName)
				if len(*locClientType) > 0 && !strings.HasSuffix(s.GoName, *locClientType) {
					continue
				}
				for _, method := range s.Methods {
					desc := method.Desc
					client.ServiceMethods[s.GoName] = append(client.ServiceMethods[s.GoName], &internal.MethodSource{
						Name: string(desc.Name()),
						Input: &internal.ParamsNames{
							Name:     string(desc.Input().Name()),
							FullName: string(desc.Input().FullName()),
						},
						Output: &internal.ParamsNames{
							Name:     string(desc.Output().Name()),
							FullName: string(desc.Output().FullName()),
						},
					})

				}

			}
			if len(client.ServiceMethods) > 0 {
				localClients = append(localClients, client)
			}
			once.Do(func() {
				path := string(f.GoImportPath)
				if len(*apiImportPrefix) > 0 {
					path = fmt.Sprintf("%s/%s", *apiImportPrefix, path)
				}
				req.ImportPaths = append(req.ImportPaths, path)
			})
		}
		if *isOpen && len(serviceNames) > 0 {
			req.ServiceNames = serviceNames
			wireReq.ServiceNames = serviceNames
			sort.Strings(req.ImportPaths)
			//生成对应api下grpc_client以及wire
			internal.GrpcClientsFile(req)
			internal.WireFile(wireReq)
			//生成对应apps下localclients以及wire
			wireLocReq := &internal.WireFileSource{
				PackageName:  "localclients",
				ImportPaths:  []string{internal.WireImportPath},
				ServiceNames: []string{},
				WireFileName: "./wire.go",
			}
			for _, local := range localClients {
				internal.LocalClientsFile(local, wireLocReq)
			}
			if len(*locFileDir) > 0 {
				wireLocReq.WireFileName = fmt.Sprintf("./%s/wire.go", strings.Trim(*locFileDir, "/"))
			}
			if len(wireLocReq.ServiceNames) > 0 {
				internal.WireFile(wireLocReq)
			}
			//debug查看数据

			//var buf bytes.Buffer
			//buf.WriteString(time.Now().Format(time.DateTime))
			//buf.WriteString("\n")
			//marshal, _ := json.Marshal(req)
			//buf.Write(marshal)
			//buf.WriteString("\n")
			//marshal1, _ := json.Marshal(wireReq)
			//buf.Write(marshal1)
			//buf.WriteString("\n")
			//marshal2, _ := json.Marshal(localClients)
			//buf.Write(marshal2)
			//buf.WriteString("\n")
			//marshal3, _ := json.Marshal(wireLocReq)
			//buf.Write(marshal3)
			//buf.WriteString("\n")
			//if err := os.WriteFile("./test.json", buf.Bytes(), 0644); err != nil {
			//	panic(err)
			//}
		}
		return nil
	})

}
