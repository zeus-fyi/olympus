package hera_v1_codegen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"strings"
)

type FunctionInfo struct {
	Name       string
	Parameters string
	ReturnType string
	Body       string
}

func extractGoFileInfo(src interface{}) (*GoCodeFile, error) {
	fset := token.NewFileSet() // positions are relative to fset
	fileInfo := &GoCodeFile{
		Functions: make(map[string]FunctionInfo),
		Variables: make(map[string]GoVar),
		Constants: make(map[string]GoVar),
		Structs:   make(map[string]string),
	}
	var reader *strings.Reader
	switch s := src.(type) {
	case string:
		reader = strings.NewReader(s)
		fileInfo.Contents = s
	case []byte:
		reader = strings.NewReader(string(s))
		fileInfo.Contents = string(s)
	default:
		return nil, fmt.Errorf("src must be a string or []byte")
	}
	// Parse the file given in arguments
	// Parse the source code
	node, err := parser.ParseFile(fset, "", reader, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// Package name
	fileInfo.PackageName = node.Name.Name

	// Inspect the AST and find all imports and functions
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.ImportSpec:
			if x.Path != nil {
				fileInfo.Imports = append(fileInfo.Imports, x.Path.Value)
			}
		case *ast.FuncDecl:
			var funcInfo FunctionInfo
			if x.Name != nil {
				funcInfo.Name = x.Name.Name
			}
			// Parameters and Return Type
			funcInfo.Parameters = fieldListToString(fset, x.Type.Params)
			funcInfo.ReturnType = fieldListToString(fset, x.Type.Results)
			// Function body
			var body bytes.Buffer
			if x.Body != nil {
				err = printer.Fprint(&body, fset, x.Body)
				if err != nil {
					panic(err)
				}
				funcInfo.Body = body.String()
			}
			fileInfo.Functions[funcInfo.Name] = funcInfo
		case *ast.GenDecl:
			switch x.Tok {
			case token.VAR, token.CONST:
				for _, spec := range x.Specs {
					if valueSpec, ok := spec.(*ast.ValueSpec); ok {
						varType := ""
						if valueSpec.Type != nil {
							varType = exprToString(fset, valueSpec.Type)
						}

						for i, name := range valueSpec.Names {
							var goVar GoVar

							// Handle variable value and infer type if necessary
							if valueSpec.Values != nil && i < len(valueSpec.Values) {
								switch v := valueSpec.Values[i].(type) {
								case *ast.CompositeLit:
									goVar.Value = compositeLitToString(fset, v)
									if varType == "" {
										varType = exprToString(fset, v.Type)
									}
								case *ast.BasicLit:
									goVar.Value = v.Value
									if varType == "" {
										switch v.Kind {
										case token.INT:
											varType = "int"
										case token.FLOAT:
											varType = "float"
										case token.IMAG:
											varType = "complex"
										case token.CHAR:
											varType = "char"
										case token.STRING:
											varType = "string"
										default:
											varType = "unknown"
										}
									}
								default:
									goVar.Value = exprToString(fset, valueSpec.Values[i])
								}
							}
							goVar.Type = varType

							// Get the raw text of the variable declaration
							var buf bytes.Buffer
							err = printer.Fprint(&buf, fset, valueSpec)
							if err != nil {
								panic(err)
							}
							switch x.Tok {
							case token.CONST:
								goVar.Content = "const " + buf.String()
								fileInfo.Constants[name.Name] = goVar
							case token.VAR:
								goVar.Content = "var " + buf.String()
								fileInfo.Variables[name.Name] = goVar
							}
						}
					}
				}
			case token.TYPE:
				for _, spec := range x.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						if _, ok := typeSpec.Type.(*ast.StructType); ok {
							fileInfo.Structs[typeSpec.Name.Name] = exprToString(fset, typeSpec.Type)
						}
					}
				}
			}
		}
		return true
	})

	return fileInfo, nil
}

func fieldListToString(fset *token.FileSet, fl *ast.FieldList) string {
	if fl == nil {
		return ""
	}

	var buf bytes.Buffer
	for _, field := range fl.List {
		if len(field.Names) > 0 {
			for _, name := range field.Names {
				buf.WriteString(name.Name)
				buf.WriteString(" ")
			}
		}

		err := printer.Fprint(&buf, fset, field.Type)
		if err != nil {
			panic(err)
		}
		buf.WriteString(", ")
	}

	// Remove trailing comma and space
	result := buf.String()
	if len(result) >= 2 {
		result = result[:len(result)-2]
	}
	return result
}

func compositeLitToString(fset *token.FileSet, cl *ast.CompositeLit) string {
	if cl == nil {
		return ""
	}
	var buf bytes.Buffer
	err := printer.Fprint(&buf, fset, cl)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func exprToString(fset *token.FileSet, expr ast.Expr) string {
	if expr == nil {
		return ""
	}
	var buf bytes.Buffer
	err := printer.Fprint(&buf, fset, expr)
	if err != nil {
		panic(err)
	}
	return buf.String()
}
