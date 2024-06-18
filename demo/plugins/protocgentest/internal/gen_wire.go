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

type WireFileSource struct {
	PackageName  string   `json:"package_name"`
	ImportPaths  []string `json:"import_paths"`
	ServiceNames []string `json:"service_names"`
	WireFileName string   `json:"wire_file_name"`
}

func WireFile(req *WireFileSource) {
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

	valDecl := &ast.GenDecl{
		Tok: token.VAR,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names: []*ast.Ident{
					{
						Name: "ProviderSet",
						Obj: &ast.Object{
							Kind: ast.Var,
							Name: "ProviderSet",
						},
					},
				},
				Values: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("wire"),
							Sel: ast.NewIdent("NewSet"),
						},
						Args: func() []ast.Expr {
							var res []ast.Expr
							for _, name := range req.ServiceNames {
								res = append(res, ast.NewIdent(newXXWithoutRPC(name)))
							}
							return res
						}(),
					},
				},
			},
		},
	}
	file.Decls = append(file.Decls, importDecl, valDecl)
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, file); err != nil {
		panic(err)
	}
	if _, err := os.Stat(filepath.Dir(req.WireFileName)); err == nil {
		// 将源代码写入文件
		if err := os.WriteFile(req.WireFileName, buf.Bytes(), 0644); err != nil {
			panic(err)
		}
		fmt.Printf("wire code written to %s ====", req.WireFileName)
	}
}
