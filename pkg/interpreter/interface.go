package interpreter

import (
	"github.com/DDP-Projekt/Kompilierer/pkg/ast"
	"github.com/DDP-Projekt/Kompilierer/pkg/ddperror"
	"github.com/DDP-Projekt/Kompilierer/pkg/parser"
)

// interpret the given ddp file
func InterpretFile(path string, errorHandler ddperror.Handler) error {
	Ast, err := parser.ParseFile(path, errorHandler)
	if err != nil {
		return err
	}

	return New(Ast, errorHandler).Interpret()
}

// interpret the given ddp source code
func InterpretSource(name string, src []byte, errorHandler ddperror.Handler) error {
	Ast, err := parser.ParseSource(name, src, errorHandler)
	if err != nil {
		return err
	}

	return New(Ast, errorHandler).Interpret()
}

// interpret the given AST
func InterpretAst(Ast *ast.Ast, errorHandler ddperror.Handler) error {
	return New(Ast, errorHandler).Interpret()
}
