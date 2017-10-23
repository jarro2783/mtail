// Copyright 2017 Google Inc. All Rights Reserved.
// This file is available under the Apache license.

package vm

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/mtail/metrics"
)

// Sexp is for converting program syntax trees into typed s-expression for printing
type Sexp struct {
	output string // Accumulator for the result

	emitTypes bool

	col  int // column to indent current line to
	line string
}

func (s *Sexp) indent() {
	s.col += 2
}

func (s *Sexp) outdent() {
	s.col -= 2
}

func (s *Sexp) prefix() (r string) {
	for i := 0; i < s.col; i++ {
		r += " "
	}
	return
}

func (s *Sexp) emit(str string) {
	s.line += str
}

func (s *Sexp) newline() {
	if s.line != "" {
		s.output += s.prefix() + s.line
	}
	s.output += "\n"
	s.line = ""
}

func (s *Sexp) VisitBefore(n astNode) Visitor {
	s.emit(fmt.Sprintf("( ;;%T ", n))
	if s.emitTypes {
		s.emit(fmt.Sprintf("<%s> ", n.Type()))
	}
	s.emit(fmt.Sprintf("@ %s", n.Pos()))
	s.newline()
	s.indent()
	switch v := n.(type) {

	case *regexNode:
		s.emit("\"" + v.pattern + "\"")

	case *binaryExprNode:
		switch v.op {
		case LT:
			s.emit("<")
		case GT:
			s.emit(">")
		case LE:
			s.emit("<=")
		case GE:
			s.emit(">=")
		case EQ:
			s.emit("==")
		case NE:
			s.emit("!=")
		case SHL:
			s.emit("<<")
		case SHR:
			s.emit(">>")
		case AND:
			s.emit("&")
		case OR:
			s.emit("|")
		case XOR:
			s.emit("^")
		case NOT:
			s.emit("~")
		case PLUS:
			s.emit("+")
		case MINUS:
			s.emit("-")
		case MUL:
			s.emit("*")
		case DIV:
			s.emit("/")
		case POW:
			s.emit("**")
		case ASSIGN:
			s.emit("=")
		case MOD:
			s.emit("%")
		}
		s.newline()
		s.indent()

	case *idNode:
		s.emit("\"" + v.name + "\"")

	case *caprefNode:
		s.emit("\"" + v.name + "\"")

	case *builtinNode:
		s.emit("\"" + v.name + "\"")
		s.newline()

	case *declNode:
		switch v.kind {
		case metrics.Counter:
			s.emit("counter ")
		case metrics.Gauge:
			s.emit("gauge ")
		case metrics.Timer:
			s.emit("timer ")
		}
		s.emit(v.name)
		if len(v.keys) > 0 {
			s.emit(" (")
			s.emit(strings.Join(v.keys, " "))
			s.emit(")")
		}

	case *unaryExprNode:
		switch v.op {
		case INC:
			s.emit("++")
		case NOT:
			s.emit("~")
		}

	case *stringConstNode:
		s.emit("\"" + v.text + "\"")

	case *intConstNode:
		s.emit(strconv.FormatInt(v.i, 10))

	case *floatConstNode:
		s.emit(strconv.FormatFloat(v.f, 'g', -1, 64))

	case *nextNode:
		s.emit("next")
	case *otherwiseNode:
		s.emit("otherwise")
	case *delNode:
		s.emit("del")

	case *indexedExprNode, *stmtlistNode, *exprlistNode, *condNode, *defNode, *decoNode: // normal walk

	default:
		panic(fmt.Sprintf("unparser found undefined type %T", n))
	}
	return s
}

func (s *Sexp) VisitAfter(node astNode) {
	switch node.(type) {
	case *binaryExprNode:
		s.outdent()
	}
	s.outdent()
	s.emit(")")
	s.newline()
}

// Dump begins the dumping of the syntax tree, returning the s-expression as a single string
func (s *Sexp) Dump(n astNode) string {
	Walk(s, n)
	return s.output
}
