package token

import (
	"fmt"

	"github.com/DDP-Projekt/Kompilierer/src/ddptypes"
)

type TokenType int

// a single ddp token
type Token struct {
	Type      TokenType               // type of the token
	Literal   string                  // the literal from which it was scanned
	Indent    uint                    // how many levels it is indented
	Range     Range                   // the range the token spans
	AliasInfo *ddptypes.ParameterType // only present in ALIAS_PARAMETERs, holds type information, nil otherwise
}

func (t *Token) String() string {
	if t == nil {
		return "<nil>"
	}
	return t.Literal
}

func (t *Token) StringVerbose() string {
	if t == nil {
		return "<nil>"
	}
	return fmt.Sprintf("[L: %d C: %d I: %d Lit: \"%s\"] Type: %s", t.Range.Start.Line, t.Range.Start.Column, t.Indent, t.Literal, t.Type)
}

// t.Range.Start.Line
func (t *Token) Line() uint {
	return t.Range.Start.Line
}

// t.Range.Start.Column
func (t *Token) Column() uint {
	return t.Range.Start.Column
}
