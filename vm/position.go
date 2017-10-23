// Copyright 2016 Google Inc. All Rights Reserved.
// This file is available under the Apache license.

package vm

import "fmt"

// A position is the location in the source program that a token appears.
type position struct {
	filename string
	line     int // Line in the source for this token.
	startcol int // Starting and ending columns in the source for this token.
	endcol   int
}

func (p position) String() string {
	r := fmt.Sprintf("%s:%d:%d", p.filename, p.line+1, p.startcol+1)
	if p.endcol > p.startcol {
		r += fmt.Sprintf("-%d", p.endcol+1)
	}
	return r
}

func MergePosition(a, b *position) *position {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a.filename != b.filename {
		return a
	}
	// TODO(jaq): handle multi-line positions
	if a.line != b.line {
		return a
	}
	r := *a
	if b.startcol < r.startcol {
		r.startcol = b.startcol
	}
	if b.endcol > r.endcol {
		r.endcol = b.endcol
	}
	return &r
}

// mergepositionlist is a helper that merges the positions of all the nodes in a list
func mergepositionlist(l []astNode) *position {
	if len(l) == 0 {
		return nil
	}
	if len(l) == 1 {
		if l[0] != nil {
			return l[0].Pos()
		} else {
			return nil
		}
	}
	return MergePosition(l[0].Pos(), mergepositionlist(l[1:]))
}
