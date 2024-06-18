package internal

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"google.golang.org/protobuf/compiler/protogen"
	"os"
	"sort"
	"strings"
)

type LocalClientsFileSource struct {
	ModuleName     string                     `json:"module_name"`
	FileDir        string                     `json:"file_dir"`
	Name           string                     `json:"name"`
	BaseName       string                     `json:"base_name"`
	ImportPaths    []string                   `json:"import_paths"`
	ServiceMethods map[string][]*MethodSource `json:"service_methods"`
}
type MethodSource struct {
	Name   string       `json:"name"`
	Input  *ParamsNames `json:"input"`
	Output *ParamsNames `json:"output"`
}
type ParamsNames struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
}

func (obj *LocalClientsFileSource) CheckGRPCServicesDirAndFile() (string, bool) {
	if len(obj.FileDir) == 0 || len(obj.ModuleName) == 0 {
		return "", false
	}

	dir := fmt.Sprintf(FMTGRPCServersFilePath, obj.ModuleName)
	fi, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return "", false
		}
		return "", false // 返回其他错误
	}
	if fi.IsDir() {
		path := fmt.Sprintf("%s/%s.go", dir, obj.BaseName)
		if _, err := os.Stat(path); err == nil {
			return fmt.Sprintf(FMTGRPCServersImportPath, obj.ModuleName), true
		}
	}
	return "", false
}

func (obj *LocalClientsFileSource) LocalClientsGoFileName() (string, bool) {
	//"./apps/%s/handlers/localclients/%s.go"
	fi, err := os.Stat(obj.FileDir)
	if err != nil {
		if os.IsNotExist(err) {
			return "", false
		}
		return "", false // 返回其他错误
	}
	if fi.IsDir() {
		return fmt.Sprintf("./%s/%s.go", obj.FileDir, obj.BaseName), true
	}
	return "", false
}

func (obj *LocalClientsFileSource) ImportProtoName() string {
	return fmt.Sprintf(FMTProtoImportPath, obj.ModuleName)
}

func LocalClientsFile(req *LocalClientsFileSource, wireReq *WireFileSource) {
	fset := token.NewFileSet()
	file := &ast.File{}
	// 固定包名
	file.Name = &ast.Ident{Name: "localclients"}
	grpcservicesName := "srv"
	contextName := "ctx"
	inName := "in"
	objName := "obj"
	grpcservicesImport, ok := req.CheckGRPCServicesDirAndFile()
	if !ok {
		return
	}
	protoORptypesFn := func(param *ParamsNames) *ast.Ident {
		if strings.Contains(param.FullName, "proto") {
			return ast.NewIdent("proto")
		}
		if strings.Contains(param.FullName, "ptypes") {
			return ast.NewIdent("ptypes")
		}
		return nil
	}
	req.ImportPaths = append(req.ImportPaths, grpcservicesImport, req.ImportProtoName())
	sort.Strings(req.ImportPaths)
	var importsSpecs []ast.Spec
	for _, path := range req.ImportPaths {
		importsSpecs = append(importsSpecs, &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: protogen.GoImportPath(path).String(),
			},
		})
	}
	//导入包
	importDecl := &ast.GenDecl{
		Tok:   token.IMPORT,
		Specs: importsSpecs,
	}
	file.Decls = append(file.Decls, importDecl)
	for structName, methods := range req.ServiceMethods {
		//service结构体
		clientStructDecl := &ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name: &ast.Ident{Name: xxWithoutRPC(structName)},
					Type: &ast.StructType{
						Fields: &ast.FieldList{
							List: []*ast.Field{
								{
									Type: &ast.SelectorExpr{
										X:   ast.NewIdent("proto"),
										Sel: ast.NewIdent(xxServer(structName)),
									},
								},
							},
						},
					},
				},
			},
		}
		//newService方法
		newStructDecl := &ast.FuncDecl{
			Name: ast.NewIdent(newXXWithoutRPC(structName)),
			Type: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								{
									Name: grpcservicesName,
								},
							},
							Type: &ast.StarExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("grpcservers"),
									Sel: ast.NewIdent(xxWithoutRPC(structName)),
								},
							},
						},
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.SelectorExpr{
								X:   ast.NewIdent("proto"),
								Sel: ast.NewIdent(xxClient(structName)),
							},
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Results: []ast.Expr{
							&ast.UnaryExpr{
								Op: token.AND,
								X: &ast.CompositeLit{
									Type: ast.NewIdent(xxWithoutRPC(structName)),
									Elts: []ast.Expr{
										&ast.KeyValueExpr{
											Key:   ast.NewIdent(xxServer(structName)),
											Value: ast.NewIdent(grpcservicesName),
										},
									},
								},
							},
						},
					},
				},
			},
		}
		wireReq.ServiceNames = append(wireReq.ServiceNames, structName)
		file.Decls = append(file.Decls, clientStructDecl, newStructDecl)
		//生成每个类的方法
		for _, method := range methods {
			methodFuncDecl := &ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								ast.NewIdent(objName),
							},
							Type: &ast.StarExpr{
								X: ast.NewIdent(xxWithoutRPC(structName)),
							},
						},
					},
				},
				Name: ast.NewIdent(method.Name),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent(contextName),
								},
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("context"),
									Sel: ast.NewIdent("Context"),
								},
							},
							{
								Names: []*ast.Ident{
									ast.NewIdent(inName),
								},
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X:   protoORptypesFn(method.Input),
										Sel: ast.NewIdent(method.Input.Name),
									},
								},
							},
							{
								Names: []*ast.Ident{
									ast.NewIdent("opts"),
								},
								Type: &ast.Ellipsis{
									Elt: &ast.SelectorExpr{
										X:   ast.NewIdent("grpc"),
										Sel: ast.NewIdent("CallOption"),
									},
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X:   protoORptypesFn(method.Output),
										Sel: ast.NewIdent(method.Output.Name),
									},
								},
							},
							{
								Type: ast.NewIdent("error"),
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent(objName),
											Sel: ast.NewIdent(xxServer(structName)),
										},
										Sel: ast.NewIdent(method.Name),
									},
									Args: []ast.Expr{
										ast.NewIdent(contextName),
										ast.NewIdent(inName),
									},
								},
							},
						},
					},
				},
			}
			file.Decls = append(file.Decls, methodFuncDecl)
		}
	}

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, file); err != nil {
		panic(err)
	}
	if fp, ok := req.LocalClientsGoFileName(); ok {
		// 将源代码写入文件
		if err := os.WriteFile(fp, buf.Bytes(), 0644); err != nil {
			panic(err)
		}
		fmt.Printf("locClient code written to %s ====", fp)
	}
}
