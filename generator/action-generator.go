package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

func read(reader io.ReaderAt, fset *token.FileSet, node ast.Node) string {
	startOffset := fset.Position(node.Pos()).Offset
	endOffset := fset.Position(node.End()).Offset
	var buffer []byte = make([]byte, endOffset-startOffset)
	reader.ReadAt(buffer, int64(startOffset))
	return string(buffer)
}

func readFieldListTypes(reader io.ReaderAt, fset *token.FileSet, fields *ast.FieldList) string {
	typeStrings := make([]string, 0)
	for _, field := range fields.List {
		typeStrings = append(typeStrings, read(reader, fset, field.Type))
	}
	return "(" + strings.Join(typeStrings, ", ") + ")"
}

func main() {
	if len(os.Args) < 1 {
		panic("You need to specify the action type name as the first argument.")
	}
	//actionStructName := os.Args[1]

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	targetFile := path.Join(dir, os.Getenv("GOFILE"))
	targetFileReader, err := os.OpenFile(targetFile, os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	log.Println("Processing " + targetFile)

	fset := token.NewFileSet() // positions are relative to fset

	// Parse src but stop after processing the imports.
	node, err := parser.ParseFile(fset, targetFile, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	functionHeader := `package %s

import (
	"errors"
	"reflect"

	"github.com/schidstorm/go-remote/lib"
)
`

	ast.Inspect(node, func(n ast.Node) bool {

		// Find Return Statements
		// (TInput, TOutput, TOutputError, structName, methodName)
		ret, ok := n.(*ast.FuncDecl)
		if ok && ret.Recv.NumFields() > 0 {
			tInput := read(targetFileReader, fset, ret.Type.Params.List[0].Type)
			tOutput := read(targetFileReader, fset, ret.Type.Params.List[1].Type)
			methodName := ret.Name.String()
			starReceiver, isStar := ret.Recv.List[0].Type.(*ast.StarExpr)
			receiver := ret.Recv.List[0].Type
			if isStar {
				receiver = starReceiver.X
			}
			structName := read(targetFileReader, fset, receiver)

			functionBody := fmt.Sprintf(`


type receiverType = %sRemote
type inputType = %s
type resultType = %s

type %sError struct {
	Data resultType
	Error error
}

type errorResultType = %sError

type %sRemote struct {}


func (s *receiverType) %s(client lib.Callable, opts inputType) []errorResultType {
	interfaceResults := client.Call("%s.%s", opts, reflect.TypeOf((resultType)(nil)).Elem())

	results := []errorResultType{}
	for _, interfaceResult := range interfaceResults {
		switch v := interfaceResult.Data.(type) {
		case resultType:
			results = append(results, errorResultType{v, interfaceResult.Error})
		default:
			err := interfaceResult.Error
			if err == nil {
				err = errors.New("Unknown result type")
			}
			results = append(results, errorResultType{nil, err})
		}
	}

	return results
}
			`, node.Name.Name, structName, tInput, tOutput, structName, structName, structName, methodName, structName, methodName)

			ioutil.WriteFile(path.Join(path.Dir(targetFile), "remote."+path.Base(targetFile)), []byte(functionBody), os.ModeExclusive)
		}
		return true
	})
}
