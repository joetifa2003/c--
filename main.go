package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"

	"c--/ast"
	"c--/compiler"
)

func main() {
	lex := lexer.MustSimple([]lexer.SimpleRule{
		{Name: "Keyword", Pattern: `return`},
		{Name: "Comment", Pattern: `(?i)rem[^\n]*`},
		{Name: "String", Pattern: `"(\\"|[^"])*"`},
		{Name: "Number", Pattern: `[-+]?(\d*\.)?\d+`},
		{Name: "Ident", Pattern: `[a-zA-Z_]\w*`},
		{Name: "Punct", Pattern: `[-[!@#$%^&*()+_={}\|:;"'<,>.?/]|]`},
		{Name: "whitespace", Pattern: `[ \t\n]+`},
		{Name: "EOL", Pattern: `[\n\r]+`},
	})

	parser, err := participle.Build[ast.Program](
		participle.Lexer(lex),
		participle.Elide("whitespace"),
		participle.Union[ast.Statement](ast.FnStmt{}, ast.VarStmt{}, ast.ReturnStmt{}, ast.ExprStmt{}),
		participle.Union[ast.Atom](ast.CallAtom{}, ast.IntAtom{}, ast.StringAtom{}, ast.IdentAtom{}),
	)
	if err != nil {
		panic(err)
	}

	p, err := parser.ParseString("main.cmm", `
		fn int add(int a, int b) {
			return a + b
		}

  	fn int main() {
			printf("hello world\n")	
  		return 0
  	}
  `)
	if err != nil {
		panic(err)
	}

	fmt.Println(parser.String())
	fmt.Println("\n\n")

	c := compiler.NewCompiler(p)
	output := c.Compile()

	f, err := os.Create("out/out.c")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString(output)

	fmt.Println(output)
}
