package compiler

import (
	"c--/ast"
	"fmt"
	"strings"
)

type Compiler struct {
	indent  int
	builder *strings.Builder
	program *ast.Program
}

func NewCompiler(program *ast.Program) *Compiler {
	return &Compiler{
		program: program,
		builder: &strings.Builder{},
	}
}

func (c *Compiler) Compile() string {
	c.emit("#include <stdio.h>\n\n")
	for _, stmt := range c.program.Statements {
		c.compile(stmt)
	}

	return c.builder.String()
}

func (c *Compiler) compile(node interface{}) {
	for i := 0; i < c.indent; i++ {
		c.emit(" ")
	}

	switch node := node.(type) {
	case ast.FnStmt:
		c.emit(node.Type.Kind)
		c.emit(" ")
		c.emit(node.Name)
		c.emit("() ")
		c.emit("{\n")

		c.pushIndent()
		for _, stmt := range node.Body {
			c.compile(stmt)
		}
		c.popIndent()

		c.emit("}")

	case ast.VarStmt:
		c.emit(node.Type)
		c.emit(" ")
		c.emit(node.Name)
		c.emit(" = ")
		c.compileExpr(node.Expr)
		c.emit(";")

	case ast.ReturnStmt:
		c.emit("return ")
		if node.Value != nil {
			c.compileExpr(*node.Value)
		}
		c.emit(";")

	case ast.ExprStmt:
		c.compileExpr(node.Expr)
		c.emit(";")

	default:
		panic(fmt.Sprintf("Unimplemented %T", node))
	}

	c.emit("\n")
}

func (c *Compiler) compileExpr(node interface{}) {
	switch node := node.(type) {

	case ast.BinaryExpr:
		c.compileExpr(node.Left)

	case ast.CallAtom:
		c.emit(node.Name)
		c.emit("(")
		for i, expr := range node.Args {
			c.compileExpr(expr)
			if i != len(node.Args)-1 {
				c.emit(",")
			}
		}
		c.emit(")")

	case ast.StringAtom:
		c.emit(node.Value)

	case ast.IntAtom:
		c.emit(fmt.Sprint(node.Value))

	case ast.IdentAtom:
		c.emit(node.Value)

	default:
		panic(fmt.Sprintf("Unimplemented %T", node))
	}
}

func (c *Compiler) emit(s string) {
	c.builder.WriteString(s)
}

func (c *Compiler) pushIndent() {
	c.indent += 2
}

func (c *Compiler) popIndent() {
	c.indent -= 2
}
