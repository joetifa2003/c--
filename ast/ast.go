package ast

type Program struct {
	Statements []Statement `@@*`
}

type Statement interface {
	stmt()
}

type VarStmt struct {
	Name string     `"let" @Ident`
	Type string     `":" @Ident`
	Expr BinaryExpr `"=" @@`
}

func (s VarStmt) stmt() {}

type FnStmt struct {
	Type Type          `"fn" @@`
	Name string        `@Ident`
	Args []FnParameter `"(" ( @@ ("," @@)* )? ")"`
	Body []Statement   `"{" @@* "}"`
}

type FnParameter struct {
	Type Type   `@@`
	Name string `@Ident`
}

func (t FnStmt) stmt() {}

type ExprStmt struct {
	Expr BinaryExpr `@@`
}

func (t ExprStmt) stmt() {}

type ReturnStmt struct {
	Value *BinaryExpr `"return" @@?`
}

func (t ReturnStmt) stmt() {}

type Type struct {
	Kind string `@Ident`
}

type BinaryExpr struct {
	Left  Atom        `@@`
	Op    *string     `[ @"+"`
	Right *BinaryExpr `@@ ]?`
}

type Atom interface {
	atom()
}

type CallAtom struct {
	Name string       `@Ident"("`
	Args []BinaryExpr `@@*")"`
}

func (t CallAtom) atom() {}

type IntAtom struct {
	Value int `@Number`
}

func (t IntAtom) atom() {}

type StringAtom struct {
	Value string `@String`
}

func (t StringAtom) atom() {}

type IdentAtom struct {
	Value string `@Ident`
}

func (t IdentAtom) atom() {}
