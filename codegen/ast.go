package codegen

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func ParseFile(filePath string) (*ast.File, error) {
	fset := token.NewFileSet()
	fileName := filepath.Base(filePath)
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return parser.ParseFile(fset, fileName, fileData, parser.ParseComments)
}

type FileMeta struct {
	Package string
	Structs []Struct
}

type Struct struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name string
	Type string
}

func ParseFileMeta(node *ast.File, structNames ...string) FileMeta {
	nameMap := make(map[string]bool)
	for _, name := range structNames {
		nameMap[strings.TrimSpace(name)] = true
	}

	fileMeta := FileMeta{
		Structs: make([]Struct, 0, len(structNames)),
	}
	ast.Inspect(node, func(n ast.Node) bool {
		switch ret := n.(type) {
		case *ast.TypeSpec:
			if nameMap[ret.Name.String()] {
				stVisitor := &structVisitor{
					st: Struct{
						Name: ret.Name.String(),
					},
				}
				ast.Walk(stVisitor, n)
				fileMeta.Structs = append(fileMeta.Structs, stVisitor.st)
			}
		case *ast.File:
			fileMeta.Package = ret.Name.String()
		}
		return true
	})
	return fileMeta
}

type structVisitor struct {
	st  Struct
	err error
}

func (v *structVisitor) Visit(node ast.Node) ast.Visitor {
	if v.err != nil {
		return v
	}

	if st, ok := node.(*ast.StructType); ok {
		v.st.Fields = make([]Field, 0, len(st.Fields.List))
		for _, field := range st.Fields.List {
			typeName := getTypeString(field.Type)
			v.st.Fields = append(v.st.Fields, Field{
				Type: typeName,
				Name: field.Names[0].Name,
			})
		}
	}
	return v
}

func getTypeString(tye ast.Expr) string {
	switch fType := tye.(type) {
	case *ast.Ident:
		if fType.Obj != nil {
			if objT, ok := fType.Obj.Decl.(*ast.TypeSpec); ok {
				return getTypeString(objT.Type)
			}
		}
		return fType.Name
	case *ast.SelectorExpr:
		return fType.Sel.Name
	case *ast.ArrayType:
		return getTypeString(fType.Elt)
	default:
		return ""
	}
}
