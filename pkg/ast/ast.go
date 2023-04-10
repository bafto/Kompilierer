package ast

import (
	"fmt"

	"github.com/DDP-Projekt/Kompilierer/pkg/ddperror"
	"github.com/DDP-Projekt/Kompilierer/pkg/ddptypes"
	"github.com/DDP-Projekt/Kompilierer/pkg/token"
)

// represents an Abstract Syntax Tree for a token.DDP program
type Ast struct {
	Statements []Statement   // the top level statements
	Comments   []token.Token // all the comments in the source code
	Symbols    *SymbolTable
	Faulty     bool // set if the ast has any errors (doesn't matter what from which phase they came)
}

// invoke the Visitor for each top level statement in the Ast
func WalkAst(ast *Ast, visitor FullVisitor) {
	for _, stmt := range ast.Statements {
		stmt.Accept(visitor)
	}
}

// wrapper for an alias
type FuncAlias struct {
	Tokens   []token.Token                     // tokens of the alias
	Original token.Token                       // the original string
	Func     *FuncDecl                         // the function it refers to (if it is used outside a FuncDecl)
	Args     map[string]ddptypes.ParameterType // types of the arguments (used for funcCall parsing)
}

// basic Node interfaces
type (
	Node interface {
		fmt.Stringer
		Token() token.Token
		GetRange() token.Range
		Accept(FullVisitor)
	}

	Expression interface {
		Node
		expressionNode() // dummy function for the interface
	}

	Statement interface {
		Node
		statementNode() // dummy function for the interface
	}

	Declaration interface {
		Node
		declarationNode() // dummy function for the interface
		Name() string     // returns the name of the declaration or "" for BadDecls
		Public() bool     // returns wether the declaration is public. always false for BadDecls
	}

	// *Ident or *Indexing
	// Nodes that fulfill this interface can be
	// on the left side of an assignement (meaning, variables or references)
	Assigneable interface {
		Expression
		assigneable() // dummy function for the interface
	}
)

// Declarations
type (
	// an invalid Declaration
	BadDecl struct {
		Tok token.Token
		Err ddperror.Error
	}

	VarDecl struct {
		Range    token.Range
		Comment  *token.Token  // optional comment (also contained in ast.Comments)
		Type     ddptypes.Type // type of the variable
		NameTok  token.Token   // identifier name
		IsPublic bool          // wether the function is marked with öffentliche
		// wether the variable was declared in the global scope
		// used pretty much only in the compiler
		IsGlobal bool
		InitVal  Expression // initial value
	}

	FuncDecl struct {
		Range         token.Range
		Comment       *token.Token             // optional comment (also contained in ast.Comments)
		Tok           token.Token              // Die
		NameTok       token.Token              // token of the name
		IsPublic      bool                     // wether the function is marked with öffentliche
		ParamNames    []token.Token            // x, y und z
		ParamTypes    []ddptypes.ParameterType // type, and wether the argument is a reference
		ParamComments []*token.Token           // comments for the parameters, the slice or its elements may be nil
		Type          ddptypes.Type            // Zahl Kommazahl nichts ...
		Body          *BlockStmt               // nil for extern functions
		ExternFile    token.Token              // string literal with filepath (only pesent if Body is nil)
		Aliases       []FuncAlias
	}
)

func (decl *BadDecl) String() string  { return "BadDecl" }
func (decl *VarDecl) String() string  { return "VarDecl" }
func (decl *FuncDecl) String() string { return "FuncDecl" }

func (decl *BadDecl) Token() token.Token  { return decl.Tok }
func (decl *VarDecl) Token() token.Token  { return decl.NameTok }
func (decl *FuncDecl) Token() token.Token { return decl.Tok }

func (decl *BadDecl) GetRange() token.Range  { return decl.Err.Range }
func (decl *VarDecl) GetRange() token.Range  { return decl.Range }
func (decl *FuncDecl) GetRange() token.Range { return decl.Range }

func (decl *BadDecl) Accept(visitor FullVisitor)  { visitor.VisitBadDecl(decl) }
func (decl *VarDecl) Accept(visitor FullVisitor)  { visitor.VisitVarDecl(decl) }
func (decl *FuncDecl) Accept(visitor FullVisitor) { visitor.VisitFuncDecl(decl) }

func (decl *BadDecl) declarationNode()  {}
func (decl *VarDecl) declarationNode()  {}
func (decl *FuncDecl) declarationNode() {}

func (decl *BadDecl) Name() string  { return "" }
func (decl *VarDecl) Name() string  { return decl.NameTok.Literal }
func (decl *FuncDecl) Name() string { return decl.NameTok.Literal }

func (decl *BadDecl) Public() bool  { return false }
func (decl *VarDecl) Public() bool  { return decl.IsPublic }
func (decl *FuncDecl) Public() bool { return decl.IsPublic }

// Expressions
type (
	BadExpr struct {
		Tok token.Token
		Err ddperror.Error
	}

	Ident struct {
		Literal token.Token
		// the variable declaration this identifier refers to
		// is set by the resolver, or nil if the name was not found
		Declaration *VarDecl
	}

	// also exists as Binary expression for Literals
	// this one can count as Reference, and may be used
	// inplace of Ident (may be assigned to etc.)
	Indexing struct {
		Lhs   Assigneable // variable Name
		Index Expression
	}

	IntLit struct {
		Literal token.Token
		Value   int64
	}

	FloatLit struct {
		Literal token.Token
		Value   float64 // the parsed float
	}

	BoolLit struct {
		Literal token.Token
		Value   bool
	}

	CharLit struct {
		Literal token.Token
		Value   rune
	}

	StringLit struct {
		Literal token.Token
		Value   string // the evaluated string
	}

	ListLit struct {
		Tok   token.Token
		Range token.Range
		// type of the empty list if Values is nil
		// the typechecker fills this field if Values is non-nil
		Type   ddptypes.Type
		Values []Expression // the values in the Literal
		// if Values, Count and Value are nil, the list is empty
		Count Expression // for big list initializations
		Value Expression // the default value for big list initializations
	}

	UnaryExpr struct {
		Range    token.Range
		Tok      token.Token
		Operator UnaryOperator
		Rhs      Expression
	}

	BinaryExpr struct {
		Range    token.Range
		Tok      token.Token
		Lhs      Expression
		Operator BinaryOperator
		Rhs      Expression
	}

	// currently only used for von bis
	TernaryExpr struct {
		Range    token.Range
		Tok      token.Token
		Lhs      Expression
		Mid      Expression
		Rhs      Expression
		Operator TernaryOperator
	}

	// als Expressions cannot be unary
	// because the type operator might be multiple
	// tokens long
	CastExpr struct {
		Range token.Range
		Type  ddptypes.Type
		Lhs   Expression
	}

	Grouping struct {
		Range  token.Range
		LParen token.Token // (
		Expr   Expression
	}

	FuncCall struct {
		Range token.Range
		Tok   token.Token // first token of the call
		Name  string      // name of the function
		// the function declaration this call refers to
		// is set by the parser, or nil if the name was not found
		Func *FuncDecl
		Args map[string]Expression
	}
)

func (expr *BadExpr) String() string     { return "BadExpr" }
func (expr *Ident) String() string       { return "Ident" }
func (expr *Indexing) String() string    { return "Indexing" }
func (expr *IntLit) String() string      { return "IntLit" }
func (expr *FloatLit) String() string    { return "FloatLit" }
func (expr *BoolLit) String() string     { return "BoolLit" }
func (expr *CharLit) String() string     { return "CharLit" }
func (expr *StringLit) String() string   { return "StringLit" }
func (expr *ListLit) String() string     { return "ListLit" }
func (expr *UnaryExpr) String() string   { return "UnaryExpr" }
func (expr *BinaryExpr) String() string  { return "BinaryExpr" }
func (expr *TernaryExpr) String() string { return "BinaryExpr" }
func (expr *CastExpr) String() string    { return "CastExpr" }
func (expr *Grouping) String() string    { return "Grouping" }
func (expr *FuncCall) String() string    { return "FuncCall" }

func (expr *BadExpr) Token() token.Token     { return expr.Tok }
func (expr *Ident) Token() token.Token       { return expr.Literal }
func (expr *Indexing) Token() token.Token    { return expr.Lhs.Token() }
func (expr *IntLit) Token() token.Token      { return expr.Literal }
func (expr *FloatLit) Token() token.Token    { return expr.Literal }
func (expr *BoolLit) Token() token.Token     { return expr.Literal }
func (expr *CharLit) Token() token.Token     { return expr.Literal }
func (expr *StringLit) Token() token.Token   { return expr.Literal }
func (expr *ListLit) Token() token.Token     { return expr.Tok }
func (expr *UnaryExpr) Token() token.Token   { return expr.Tok }
func (expr *BinaryExpr) Token() token.Token  { return expr.Tok }
func (expr *TernaryExpr) Token() token.Token { return expr.Tok }
func (expr *CastExpr) Token() token.Token    { return expr.Lhs.Token() }
func (expr *Grouping) Token() token.Token    { return expr.LParen }
func (expr *FuncCall) Token() token.Token    { return expr.Tok }

func (expr *BadExpr) GetRange() token.Range { return expr.Err.Range }
func (expr *Ident) GetRange() token.Range   { return token.NewRange(expr.Literal, expr.Literal) }
func (expr *Indexing) GetRange() token.Range {
	return token.Range{Start: expr.Lhs.GetRange().Start, End: expr.Index.GetRange().End}
}
func (expr *IntLit) GetRange() token.Range      { return expr.Literal.Range }
func (expr *FloatLit) GetRange() token.Range    { return expr.Literal.Range }
func (expr *BoolLit) GetRange() token.Range     { return expr.Literal.Range }
func (expr *CharLit) GetRange() token.Range     { return expr.Literal.Range }
func (expr *StringLit) GetRange() token.Range   { return expr.Literal.Range }
func (expr *ListLit) GetRange() token.Range     { return expr.Range }
func (expr *UnaryExpr) GetRange() token.Range   { return expr.Range }
func (expr *BinaryExpr) GetRange() token.Range  { return expr.Range }
func (expr *TernaryExpr) GetRange() token.Range { return expr.Range }
func (expr *CastExpr) GetRange() token.Range    { return expr.Range }
func (expr *Grouping) GetRange() token.Range    { return expr.Range }
func (expr *FuncCall) GetRange() token.Range    { return expr.Range }

func (expr *BadExpr) Accept(v FullVisitor)     { v.VisitBadExpr(expr) }
func (expr *Ident) Accept(v FullVisitor)       { v.VisitIdent(expr) }
func (expr *Indexing) Accept(v FullVisitor)    { v.VisitIndexing(expr) }
func (expr *IntLit) Accept(v FullVisitor)      { v.VisitIntLit(expr) }
func (expr *FloatLit) Accept(v FullVisitor)    { v.VisitFloatLit(expr) }
func (expr *BoolLit) Accept(v FullVisitor)     { v.VisitBoolLit(expr) }
func (expr *CharLit) Accept(v FullVisitor)     { v.VisitCharLit(expr) }
func (expr *StringLit) Accept(v FullVisitor)   { v.VisitStringLit(expr) }
func (expr *ListLit) Accept(v FullVisitor)     { v.VisitListLit(expr) }
func (expr *UnaryExpr) Accept(v FullVisitor)   { v.VisitUnaryExpr(expr) }
func (expr *BinaryExpr) Accept(v FullVisitor)  { v.VisitBinaryExpr(expr) }
func (expr *TernaryExpr) Accept(v FullVisitor) { v.VisitTernaryExpr(expr) }
func (expr *CastExpr) Accept(v FullVisitor)    { v.VisitCastExpr(expr) }
func (expr *Grouping) Accept(v FullVisitor)    { v.VisitGrouping(expr) }
func (expr *FuncCall) Accept(v FullVisitor)    { v.VisitFuncCall(expr) }

func (expr *BadExpr) expressionNode()     {}
func (expr *Ident) expressionNode()       {}
func (expr *Indexing) expressionNode()    {}
func (expr *IntLit) expressionNode()      {}
func (expr *FloatLit) expressionNode()    {}
func (expr *BoolLit) expressionNode()     {}
func (expr *CharLit) expressionNode()     {}
func (expr *StringLit) expressionNode()   {}
func (expr *ListLit) expressionNode()     {}
func (expr *UnaryExpr) expressionNode()   {}
func (expr *BinaryExpr) expressionNode()  {}
func (expr *TernaryExpr) expressionNode() {}
func (expr *CastExpr) expressionNode()    {}
func (expr *Grouping) expressionNode()    {}
func (expr *FuncCall) expressionNode()    {}

func (expr *Ident) assigneable()    {}
func (expr *Indexing) assigneable() {}

// Statements
type (
	BadStmt struct {
		Tok token.Token
		Err ddperror.Error
	}

	DeclStmt struct {
		Decl Declaration
	}

	ExprStmt struct {
		Expr Expression
	}

	// import statement for meta-information in the ast
	// will be already resolved by the parser
	ImportStmt struct {
		Range token.Range
		// the string literal which specified the filename
		FileName token.Token
		// the module that was imported because of this
		// nil if it does not exist or a similar error occured while importing
		Module *Module
		// slice of identifiers which specify
		// the individual symbols imported
		// if nil, all symbols are imported
		ImportedSymbols []token.Token
	}

	AssignStmt struct {
		Range token.Range
		Tok   token.Token
		Var   Assigneable // the variable to assign to
		Rhs   Expression  // literal assign value
	}

	BlockStmt struct {
		Range      token.Range
		Colon      token.Token
		Statements []Statement
		Symbols    *SymbolTable
	}

	IfStmt struct {
		Range     token.Range
		If        token.Token // wenn/aber
		Condition Expression
		Then      Statement
		Else      Statement
	}

	WhileStmt struct {
		Range     token.Range
		While     token.Token // solange, mache, mal
		Condition Expression
		Body      Statement
	}

	ForStmt struct {
		Range       token.Range
		For         token.Token // Für
		Initializer *VarDecl    // Zahl (name) von (Initializer.InitVal)
		To          Expression  // bis (To)
		StepSize    Expression  // Schrittgröße
		Body        *BlockStmt
	}

	ForRangeStmt struct {
		Range       token.Range
		For         token.Token // Für
		Initializer *VarDecl    // InitVal is the same pointer as In
		In          Expression  // the string/list to range over
		Body        *BlockStmt
	}

	ReturnStmt struct {
		Range  token.Range
		Return token.Token // Gib
		Func   string
		Value  Expression // nil for void return
	}
)

func (stmt *BadStmt) String() string      { return "BadStmt" }
func (stmt *DeclStmt) String() string     { return "DeclStmt" }
func (stmt *ExprStmt) String() string     { return "ExprStmt" }
func (stmt *ImportStmt) String() string   { return "ImportStmt" }
func (stmt *AssignStmt) String() string   { return "AssignStmt" }
func (stmt *BlockStmt) String() string    { return "BlockStmt" }
func (stmt *IfStmt) String() string       { return "IfStmt" }
func (stmt *WhileStmt) String() string    { return "WhileStmt" }
func (stmt *ForStmt) String() string      { return "ForStmt" }
func (stmt *ForRangeStmt) String() string { return "ForRangeStmt" }
func (stmt *ReturnStmt) String() string   { return "ReturnStmt" }

func (stmt *BadStmt) Token() token.Token      { return stmt.Tok }
func (stmt *DeclStmt) Token() token.Token     { return stmt.Decl.Token() }
func (stmt *ExprStmt) Token() token.Token     { return stmt.Expr.Token() }
func (stmt *ImportStmt) Token() token.Token   { return stmt.FileName }
func (stmt *AssignStmt) Token() token.Token   { return stmt.Tok }
func (stmt *BlockStmt) Token() token.Token    { return stmt.Colon }
func (stmt *IfStmt) Token() token.Token       { return stmt.If }
func (stmt *WhileStmt) Token() token.Token    { return stmt.While }
func (stmt *ForStmt) Token() token.Token      { return stmt.For }
func (stmt *ForRangeStmt) Token() token.Token { return stmt.For }
func (stmt *ReturnStmt) Token() token.Token   { return stmt.Return }

func (stmt *BadStmt) GetRange() token.Range      { return stmt.Err.Range }
func (stmt *DeclStmt) GetRange() token.Range     { return stmt.Decl.GetRange() }
func (stmt *ExprStmt) GetRange() token.Range     { return stmt.Expr.GetRange() }
func (stmt *ImportStmt) GetRange() token.Range   { return stmt.Range }
func (stmt *AssignStmt) GetRange() token.Range   { return stmt.Range }
func (stmt *BlockStmt) GetRange() token.Range    { return stmt.Range }
func (stmt *IfStmt) GetRange() token.Range       { return stmt.Range }
func (stmt *WhileStmt) GetRange() token.Range    { return stmt.Range }
func (stmt *ForStmt) GetRange() token.Range      { return stmt.Range }
func (stmt *ForRangeStmt) GetRange() token.Range { return stmt.Range }
func (stmt *ReturnStmt) GetRange() token.Range   { return stmt.Range }

func (stmt *BadStmt) Accept(v FullVisitor)      { v.VisitBadStmt(stmt) }
func (stmt *DeclStmt) Accept(v FullVisitor)     { v.VisitDeclStmt(stmt) }
func (stmt *ExprStmt) Accept(v FullVisitor)     { v.VisitExprStmt(stmt) }
func (stmt *ImportStmt) Accept(v FullVisitor)   { v.VisitImportStmt(stmt) }
func (stmt *AssignStmt) Accept(v FullVisitor)   { v.VisitAssignStmt(stmt) }
func (stmt *BlockStmt) Accept(v FullVisitor)    { v.VisitBlockStmt(stmt) }
func (stmt *IfStmt) Accept(v FullVisitor)       { v.VisitIfStmt(stmt) }
func (stmt *WhileStmt) Accept(v FullVisitor)    { v.VisitWhileStmt(stmt) }
func (stmt *ForStmt) Accept(v FullVisitor)      { v.VisitForStmt(stmt) }
func (stmt *ForRangeStmt) Accept(v FullVisitor) { v.VisitForRangeStmt(stmt) }
func (stmt *ReturnStmt) Accept(v FullVisitor)   { v.VisitReturnStmt(stmt) }

func (stmt *BadStmt) statementNode()      {}
func (stmt *DeclStmt) statementNode()     {}
func (stmt *ExprStmt) statementNode()     {}
func (stmt *ImportStmt) statementNode()   {}
func (stmt *AssignStmt) statementNode()   {}
func (stmt *BlockStmt) statementNode()    {}
func (stmt *IfStmt) statementNode()       {}
func (stmt *WhileStmt) statementNode()    {}
func (stmt *ForStmt) statementNode()      {}
func (stmt *ForRangeStmt) statementNode() {}
func (stmt *ReturnStmt) statementNode()   {}
