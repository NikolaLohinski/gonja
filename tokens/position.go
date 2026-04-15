package tokens

import (
	"fmt"
	"strings"
)

type Pos interface {
	Pos() int
}

// Position describes an arbitrary source position
// including the file, line, and column location.
// A Position is valid if the line number is > 0.
//
type Position struct {
	Filename string // filename, if any
	Offset   int    // offset, starting at 0
	Line     int    // line number, starting at 1
	Column   int    // column number, starting at 1 (byte count)
}

// IsValid reports whether the position is valid.
func (pos *Position) IsValid() bool { return pos.Line > 0 }

// Pos return the current offset starting at 0
func (pos *Position) Pos() int { return pos.Offset }

// String returns a string in one of several forms:
//
//	file:line:column    valid position with file name
//	file:line           valid position with file name but no column (column == 0)
//	line:column         valid position without file name
//	line                valid position without file name and no column (column == 0)
//	file                invalid position with file name
//	-                   invalid position without file name
//
func (pos Position) String() string {
	s := pos.Filename
	if pos.IsValid() {
		if s != "" {
			s += ":"
		}
		s += fmt.Sprintf("%d", pos.Line)
		if pos.Column != 0 {
			s += fmt.Sprintf(":%d", pos.Column)
		}
	}
	if s == "" {
		s = "-"
	}
	return s
}

func ReadablePosition(pos int, input string) (int, int) {
	before := input[:pos]
	lines := strings.Split(before, "\n")
	length := len(lines)
	return length, len(lines[length-1]) + 1
}

// PrecomputeLineOffsets builds a table of byte offsets where each line starts.
// offsets[0] = 0 (first line starts at byte 0), offsets[1] = position after
// the first '\n', etc. This is an O(N) pass done once per input.
func PrecomputeLineOffsets(input string) []int {
	offsets := make([]int, 1, len(input)/60+1)
	offsets[0] = 0
	for i := 0; i < len(input); i++ {
		if input[i] == '\n' {
			offsets = append(offsets, i+1)
		}
	}
	return offsets
}

// ReadablePositionFromOffsets returns (line, col) for a byte offset using
// a precomputed line offset table. O(log N) per call, zero allocations.
func ReadablePositionFromOffsets(pos int, lineOffsets []int) (int, int) {
	lo, hi := 0, len(lineOffsets)-1
	for lo < hi {
		mid := (lo + hi + 1) / 2
		if lineOffsets[mid] <= pos {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	return lo + 1, pos - lineOffsets[lo] + 1
}

// ReadablePositionHinted is like ReadablePositionFromOffsets but accepts a hint
// (the 0-based line index from the last lookup). Since lexer positions increase
// monotonically, the hint lets us skip the binary search entirely when the
// position is on the same or next line, which is the common case.
// Returns (line, col, newHint).
func ReadablePositionHinted(pos int, lineOffsets []int, hint int) (int, int, int) {
	n := len(lineOffsets)
	// Clamp hint
	if hint < 0 || hint >= n {
		hint = 0
	}
	// Fast path: check if pos is on the hinted line
	if lineOffsets[hint] <= pos {
		// Check if pos is before the next line (or hint is the last line)
		if hint+1 >= n || lineOffsets[hint+1] > pos {
			return hint + 1, pos - lineOffsets[hint] + 1, hint
		}
		// Check next line (very common: token crosses to next line)
		if hint+2 >= n || lineOffsets[hint+2] > pos {
			return hint + 2, pos - lineOffsets[hint+1] + 1, hint + 1
		}
		// Scan forward a few lines before falling back to binary search
		for i := hint + 2; i < n && i < hint+8; i++ {
			if i+1 >= n || lineOffsets[i+1] > pos {
				return i + 1, pos - lineOffsets[i] + 1, i
			}
		}
	}
	// Fall back to binary search
	lo, hi := 0, n-1
	for lo < hi {
		mid := (lo + hi + 1) / 2
		if lineOffsets[mid] <= pos {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	return lo + 1, pos - lineOffsets[lo] + 1, lo
}
