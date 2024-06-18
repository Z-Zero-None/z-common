package internal

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"google.golang.org/protobuf/compiler/protogen"
	"os"
	"path/filepath"
)

type GrpcClientFileSource struct {
	PackageName        string   `json:"package_name"`
	ImportPaths        []string `json:"import_paths"`
	ClientName         string   `json:"client_name"`
	ServiceNames       []string `json:"service_names"`
	GRPCClientFileName string   `json:"grpc_client_file_name"`
}

func GrpcClientsFile(req *GrpcClientFileSource) {
	fset := token.NewFileSet()
	file := &ast.File{}
	// 创建一个包声明
	file.Name = ast.NewIdent(req.PackageName)
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
	//创建固定结构体 type GRPCClient grpc.ClientConn
	clientTypeDecl := &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(req.ClientName),
				Type: &ast.SelectorExpr{
					X:   ast.NewIdent("grpc"),
					Sel: ast.NewIdent("ClientConn"),
				},
			},
		},
	}

	file.Decls = append(file.Decls, importDecl, clientTypeDecl)
	//创建方法和返回值
	for _, name := range req.ServiceNames {
		funcDecl := &ast.FuncDecl{
			Doc:  nil,
			Recv: nil,
			Name: ast.NewIdent(newXXWithoutRPC(name)),
			Type: &ast.FuncType{ //方法 请求参数和返回值
				Params: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								ast.NewIdent("cc"),
							},
							Type: ast.NewIdent(req.ClientName),
						},
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{
							Type: &ast.SelectorExpr{
								X: ast.NewIdent("proto"),

								Sel: ast.NewIdent(xxClient(name)),
							},
						},
					},
				},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ReturnStmt{
						Return: 0,
						Results: []ast.Expr{
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: ast.NewIdent("proto"),

									Sel: ast.NewIdent(newXXClient(name)),
								},
								Args: []ast.Expr{
									ast.NewIdent("cc"),
								},
							},
						},
					},
				},
			},
		}
		file.Decls = append(file.Decls, funcDecl)
	}

	var buf bytes.Buffer
	if err := format.Node(&buf, fset, file); err != nil {
		panic(err)
	}
	if _, err := os.Stat(filepath.Dir(req.GRPCClientFileName)); err == nil {
		// 将源代码写入文件
		if err := os.WriteFile(req.GRPCClientFileName, buf.Bytes(), 0644); err != nil {
			panic(err)
		}

		// 输出到控制台以确认
		fmt.Printf("grpc_clients code written to %s ====", req.GRPCClientFileName)
	}

}
