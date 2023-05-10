package enum_mapping

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"
)

const showFlag = "notShow"

func TranEnum2MapFromFile(filePath string) (map[string]string, error) {
	fset := token.NewFileSet()
	src, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	ret := map[string]string{}
	for _, v := range f.Scope.Objects {
		if v.Kind == ast.Con { // parse comment
			d := v.Decl.(*ast.ValueSpec)
			if len(d.Values) > 0 {
				if e, ok := d.Values[0].(*ast.BasicLit); ok {
					if strings.Contains(d.Comment.Text(), showFlag) { // delete not show enum
						continue
					}
					ret[e.Value] = strings.TrimSuffix(d.Comment.Text(), "\n") // key -> enum value,value -> mark in comment
				}
			}
		}
	}
	return ret, nil
}
