package tokens

import (
	"fmt"

	// "encoding/json"
	"regexp"
	// "strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/nikolalohinski/gonja/v2/config"
	log "github.com/sirupsen/logrus"

	"github.com/nikolalohinski/gonja/v2/logging"
)

// EOF is an arbitraty value for End Of File
const rEOF = -1

var escapedStrings = map[string]string{
	`\"`: `"`,
	`\'`: `'`,
}

// var pattern = regexp.MustCompile(`(?m)(?P<key>\w+):\s+(?P<value>\w+)$`)

// lexFn represents the state of the scanner
// as a function that returns the next state.
type lexFn func() lexFn

// Lexer holds the state of the scanner.
type Lexer struct {
	Input string // the string being scanned.
	Start int    // start position of this item.
	Pos   int    // current position in the input.
	Width int    // width of last rune read from input.
	Line  int    // Current line in the input
	Col   int    // Current position in the line
	// Position Position // Current lexing position in the input
	Config               *config.Config // The lexer configuration
	Tokens               chan *Token    // channel of scanned tokens.
	delimiters           []rune
	RawControlStructures rawControlStructure
	rawEnd               *regexp.Regexp
	expressionEnd        Type
	lineStatement        bool
	lineOffsets          []int // precomputed line start offsets for O(log N) position lookups
	collected            []*Token // when non-nil, tokens are collected here instead of sent to channel
	tokenSlab            []Token  // pre-allocated token storage to reduce heap allocations
	tokenSlabIdx         int      // next free index in tokenSlab
}

// TODO: set from env
type rawControlStructure map[string]*regexp.Regexp

func normalizeInput(input string, cfg *config.Config) string {
	input = strings.ReplaceAll(input, "\r\n", "\n")
	input = strings.ReplaceAll(input, "\r", "\n")
	if !cfg.KeepTrailingNewline && strings.HasSuffix(input, "\n") {
		input = input[:len(input)-1]
	}
	return input
}

func normalizeNewlines(input, newline string) string {
	if newline == "" || newline == "\n" {
		return input
	}
	return strings.ReplaceAll(input, "\n", newline)
}

func trimLeadingWhitespaceFromLastLine(input string) string {
	lineStart := strings.LastIndex(input, "\n") + 1
	for i := lineStart; i < len(input); i++ {
		if input[i] != ' ' && input[i] != '\t' {
			return input
		}
	}
	return input[:lineStart]
}

// NewLexer creates a new scanner for the input string.
func NewLexer(input string, config *config.Config) *Lexer {
	normalizedInput := normalizeInput(input, config)
	return &Lexer{
		Input:       normalizedInput,
		Tokens:      make(chan *Token),
		Config:      config,
		lineOffsets: PrecomputeLineOffsets(normalizedInput),
		RawControlStructures: rawControlStructure{
			"raw":     regexp.MustCompile(fmt.Sprintf(`%s[-+]?\s*endraw`, regexp.QuoteMeta(config.BlockStartString))),
			"comment": regexp.MustCompile(fmt.Sprintf(`%s[-+]?\s*endcomment`, regexp.QuoteMeta(config.BlockStartString))),
		},
	}
}

func Lex(input string, config *config.Config) *Stream {
	l := NewLexer(input, config)
	go l.Run()
	return NewStream(l.Tokens)
}

// LexAll lexes the input synchronously, collecting all tokens into a slice.
// This avoids goroutine/channel overhead.
func LexAll(input string, cfg *config.Config) *Stream {
	l := NewLexer(input, cfg)
	// Estimate ~1 token per 3 bytes as initial capacity (measured ratio is ~3.8)
	estTokens := len(l.Input)/3 + 16
	l.collected = make([]*Token, 0, estTokens)
	l.tokenSlab = make([]Token, estTokens)
	l.runSync()
	return NewStream(l.collected)
}

// allocToken returns a pointer to a Token from the pre-allocated slab,
// growing it if necessary. Falls back to a new allocation in channel mode.
func (l *Lexer) allocToken() *Token {
	if l.tokenSlab != nil {
		if l.tokenSlabIdx >= len(l.tokenSlab) {
			// Grow slab
			newSlab := make([]Token, len(l.tokenSlab)*2)
			l.tokenSlab = newSlab
			l.tokenSlabIdx = 0
		}
		tok := &l.tokenSlab[l.tokenSlabIdx]
		l.tokenSlabIdx++
		return tok
	}
	return &Token{}
}

// sendToken either appends to the collected slice or sends over the channel.
func (l *Lexer) sendToken(tok *Token) {
	if l.collected != nil {
		l.collected = append(l.collected, tok)
	} else {
		l.Tokens <- tok
	}
}

// errorf returns an error token and terminates the scan
// by passing back a nil pointer that will be the next
// state, terminating Lexer.Run.
func (l *Lexer) errorf(format string, args ...any) lexFn {
	tok := l.allocToken()
	tok.Type = Error
	tok.Val = fmt.Sprintf(format, args...)
	tok.Pos = l.Pos
	l.sendToken(tok)
	return nil
}

// Position return the current position in the input
func (l *Lexer) Position() *Position {
	return &Position{
		Offset: l.Pos,
		Line:   l.Line,
		Column: l.Col,
	}
}

func (l *Lexer) Current() string {
	return l.Input[l.Start:l.Pos]
}

// Run lexes the input by executing state functions until
// the state is nil.
func (l *Lexer) Run() {
	for state := l.lexData; state != nil; {
		state = state()
	}
	close(l.Tokens) // No more tokens will be delivered.
}

// runSync lexes without using the channel (for synchronous collection mode).
func (l *Lexer) runSync() {
	for state := l.lexData; state != nil; {
		state = state()
	}
}

// next returns the next rune in the input.
func (l *Lexer) next() (rune rune) {
	if l.Pos >= len(l.Input) {
		l.Width = 0
		return rEOF
	}
	rune, l.Width = utf8.DecodeRuneInString(l.Input[l.Pos:])
	l.Pos += l.Width
	if rune == '\n' {
		l.Line++
		l.Col = 1
	}
	return rune
}

// emit passes a Token back to the client.
func (l *Lexer) emit(t Type) {
	l.processAndEmit(t, nil)
}

func (l *Lexer) processAndEmit(t Type, fn func(string) string) {
	line, col := ReadablePositionFromOffsets(l.Start, l.lineOffsets)
	val := l.Input[l.Start:l.Pos]
	if fn != nil {
		val = fn(val)
	}
	tok := l.allocToken()
	tok.Type = t
	tok.Val = val
	tok.Pos = l.Start
	tok.Line = line
	tok.Col = col
	l.sendToken(tok)
	l.Start = l.Pos
}

// backup steps back one rune.
// Can be called only once per call of next.
func (l *Lexer) backup() {
	l.Pos -= l.Width
}

// peek returns but does not consume
// the next rune in the input.
func (l *Lexer) peek() rune {
	rune := l.next()
	l.backup()
	return rune
}

// accept consumes the next rune
// if it's from the valid set.
func (l *Lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

func (l *Lexer) pushDelimiter(r rune) {
	l.delimiters = append(l.delimiters, r)
}

func (l *Lexer) hasPrefix(prefix string) bool {
	return strings.HasPrefix(l.Input[l.Pos:], prefix)
}

func (l *Lexer) hasPrefixAt(pos int, prefix string) bool {
	if pos < 0 || pos > len(l.Input) {
		return false
	}
	return strings.HasPrefix(l.Input[pos:], prefix)
}

func (l *Lexer) atLineStart(pos int) bool {
	return pos == 0 || l.Input[pos-1] == '\n'
}

func (l *Lexer) emitData() {
	l.emitDataValue(l.Input[l.Start:l.Pos])
}

func (l *Lexer) emitLeftStrippedData() {
	l.emitDataValue(trimLeadingWhitespaceFromLastLine(l.Input[l.Start:l.Pos]))
}

func (l *Lexer) emitDataValue(val string) {
	if l.Pos <= l.Start {
		return
	}
	line, col := ReadablePositionFromOffsets(l.Start, l.lineOffsets)
	val = normalizeNewlines(val, l.Config.NewlineSequence)
	if val == "" {
		l.Start = l.Pos
		return
	}
	tok := l.allocToken()
	tok.Type = Data
	tok.Val = val
	tok.Pos = l.Start
	tok.Line = line
	tok.Col = col
	l.sendToken(tok)
	l.Start = l.Pos
}

func (l *Lexer) consumeWhitespaceControl(allowed string) byte {
	if l.Pos >= len(l.Input) {
		return 0
	}
	control := l.Input[l.Pos]
	if !strings.ContainsRune(allowed, rune(control)) {
		return 0
	}
	l.Pos++
	return control
}

func (l *Lexer) getOpeningWhitespaceControl(prefix string) byte {
	pos := l.Pos + len(prefix)
	if pos >= len(l.Input) {
		return 0
	}
	control := l.Input[pos]
	if control != '-' && control != '+' {
		return 0
	}
	return control
}

func (l *Lexer) consumeFollowingNewline() {
	if l.hasPrefix("\n") {
		l.Pos++
		l.Start = l.Pos
	}
}

type rootTokenKind int

const (
	rootTokenNone rootTokenKind = iota
	rootTokenComment
	rootTokenVariable
	rootTokenBlock
)

func (l *Lexer) currentRootToken() rootTokenKind {
	var (
		longest int
		kind    rootTokenKind
	)
	for _, candidate := range []struct {
		prefix string
		kind   rootTokenKind
	}{
		{l.Config.CommentStartString, rootTokenComment},
		{l.Config.VariableStartString, rootTokenVariable},
		{l.Config.BlockStartString, rootTokenBlock},
	} {
		if candidate.prefix == "" || !l.hasPrefix(candidate.prefix) {
			continue
		}
		if len(candidate.prefix) > longest {
			longest = len(candidate.prefix)
			kind = candidate.kind
		}
	}
	return kind
}

type linePrefixKind int

const (
	linePrefixNone linePrefixKind = iota
	linePrefixStatement
	linePrefixComment
)

func (l *Lexer) currentLinePrefix() (linePrefixKind, bool) {
	if !l.atLineStart(l.Pos) {
		return linePrefixNone, false
	}
	prefixPos := l.Pos
	for prefixPos < len(l.Input) && isLinePrefixSpace(rune(l.Input[prefixPos])) {
		prefixPos++
	}
	var (
		longest int
		kind    linePrefixKind
	)
	for _, candidate := range []struct {
		prefix string
		kind   linePrefixKind
	}{
		{l.Config.LineStatementPrefix, linePrefixStatement},
		{l.Config.LineCommentPrefix, linePrefixComment},
	} {
		if candidate.prefix == "" || !l.hasPrefixAt(prefixPos, candidate.prefix) {
			continue
		}
		if len(candidate.prefix) > longest {
			longest = len(candidate.prefix)
			kind = candidate.kind
		}
	}
	return kind, kind != linePrefixNone
}

func (l *Lexer) hasInlineLineComment() bool {
	if l.Config.LineCommentPrefix == "" || l.Pos == 0 || l.atLineStart(l.Pos) {
		return false
	}
	prev, _ := utf8.DecodeLastRuneInString(l.Input[:l.Pos])
	if prev == '\n' || isLinePrefixSpace(prev) {
		return false
	}
	prefixPos := l.Pos
	for prefixPos < len(l.Input) && isLinePrefixSpace(rune(l.Input[prefixPos])) {
		prefixPos++
	}
	return l.hasPrefixAt(prefixPos, l.Config.LineCommentPrefix)
}

func (l *Lexer) lineStatementEndLength() (int, bool) {
	remaining := l.remaining()
	newlineIdx := strings.IndexByte(remaining, '\n')
	tail := remaining
	endLength := len(remaining)
	if newlineIdx >= 0 {
		tail = remaining[:newlineIdx]
		endLength = newlineIdx + 1
	}
	if !isLineStatementClosingTail(tail) {
		return 0, false
	}
	return endLength, true
}

func isLineStatementClosingTail(tail string) bool {
	index := 0
	for index < len(tail) && isLinePrefixSpace(rune(tail[index])) {
		index++
	}
	if index == len(tail) {
		return true
	}
	colons := index
	for colons < len(tail) && tail[colons] == ':' {
		colons++
	}
	if colons == index {
		return false
	}
	for colons < len(tail) {
		if !isLinePrefixSpace(rune(tail[colons])) {
			return false
		}
		colons++
	}
	return true
}

func (l *Lexer) popDelimiter(r rune) bool {
	if len(l.delimiters) == 0 {
		l.errorf(`Unexpected delimiter "%c"`, r)
		return false
	}
	last := len(l.delimiters) - 1
	expected := l.delimiters[last]
	if r != expected {
		l.errorf(`Unbalanced delimiters, expected "%c", got "%c"`, expected, r)
		return false
	}
	// l.delimiters[last] = nil // Erase element (write zero value)
	l.delimiters = l.delimiters[:last]
	return true
}

// return whether or not we are expecting r as the next delimiter
func (l *Lexer) expectDelimiter(r rune) bool {
	if len(l.delimiters) == 0 {
		return false
	}
	expected := l.delimiters[len(l.delimiters)-1]
	return r == expected
}

func (l *Lexer) lexData() lexFn {
	for {
		switch l.currentRootToken() {
		case rootTokenComment:
			if l.Config.LeftStripBlocks && l.getOpeningWhitespaceControl(l.Config.CommentStartString) != '+' {
				l.emitLeftStrippedData()
			} else {
				l.emitData()
			}
			return l.lexComment
		case rootTokenVariable:
			l.emitData()
			return l.lexVariable
		case rootTokenBlock:
			if l.Config.LeftStripBlocks && l.getOpeningWhitespaceControl(l.Config.BlockStartString) != '+' {
				l.emitLeftStrippedData()
			} else {
				l.emitData()
			}
			return l.lexBlock
		}

		if prefix, ok := l.currentLinePrefix(); ok {
			l.emitData()
			if prefix == linePrefixComment {
				return l.lexLineComment
			}
			return l.lexLineStatement
		}

		if l.hasInlineLineComment() {
			l.emitData()
			return l.lexLineComment
		}

		if l.next() == rEOF {
			break
		}
	}
	// Correctly reached EOF.
	l.emitData()
	l.emit(EOF) // Useful to make EOF a token.
	return nil  // Stop the run loop.
}

func (l *Lexer) remaining() string {
	return l.Input[l.Pos:]
}

func (l *Lexer) lexRaw() lexFn {
	loc := l.rawEnd.FindStringIndex(l.remaining())
	if loc == nil {
		return l.errorf(`Unable to find raw closing controlStructure`)
	}
	l.Pos += loc[0]
	l.emitData()
	l.rawEnd = nil
	return l.lexBlock
	// regexp.MustCompile(`(?m)(?P<key>\w+):\s+(?P<value>\w+)$`)
	// idx := pattern
}

func (l *Lexer) lexComment() lexFn {
	l.Pos += len(l.Config.CommentStartString)
	l.consumeWhitespaceControl("-+")
	l.emit(CommentBegin)
	i := strings.Index(l.Input[l.Pos:], l.Config.CommentEndString)
	if i < 0 {
		return l.errorf("unclosed comment")
	}
	contentEnd := l.Pos + i
	if contentEnd > l.Start && (l.Input[contentEnd-1] == '-' || l.Input[contentEnd-1] == '+') {
		contentEnd--
	}
	l.Pos = contentEnd
	l.emitData()
	l.Pos = contentEnd
	control := l.consumeWhitespaceControl("-+")
	l.Pos += len(l.Config.CommentEndString)
	l.emit(CommentEnd)
	if l.Config.TrimBlocks && control != '+' {
		l.consumeFollowingNewline()
	}
	return l.lexData
}

func (l *Lexer) lexVariable() lexFn {
	if logging.Enabled() {
		log.WithFields(log.Fields{
			"pos":       l.Pos,
			"input":     l.Input,
			"remaining": l.remaining(),
		}).Trace("Lexer.lexVariable")
	}
	l.Pos += len(l.Config.VariableStartString)
	l.accept("-")
	l.expressionEnd = VariableEnd
	l.emit(VariableBegin)
	return l.lexExpression
}

func (l *Lexer) lexVariableEnd() lexFn {
	l.accept("-")
	l.Pos += len(l.Config.VariableEndString)
	l.emit(VariableEnd)
	return l.lexData
}

func (l *Lexer) lexBlock() lexFn {
	l.Pos += len(l.Config.BlockStartString)
	l.consumeWhitespaceControl("-+")
	l.expressionEnd = BlockEnd
	l.emit(BlockBegin)
	for isSpace(l.peek()) {
		l.next()
	}
	if len(l.Current()) > 0 {
		l.emit(Whitespace)
	}
	controlStructure := l.nextIdentifier()
	l.emit(Name)
	re, exists := l.RawControlStructures[controlStructure]
	if exists {
		l.rawEnd = re
	}
	return l.lexExpression
}

func (l *Lexer) lexLineStatement() lexFn {
	for isLinePrefixSpace(l.peek()) {
		l.next()
	}
	l.Start = l.Pos
	l.Pos += len(l.Config.LineStatementPrefix)
	l.expressionEnd = BlockEnd
	l.lineStatement = true
	l.emit(BlockBegin)
	for isSpace(l.peek()) {
		l.next()
	}
	if len(l.Current()) > 0 {
		l.emit(Whitespace)
	}
	controlStructure := l.nextIdentifier()
	l.emit(Name)
	if re, exists := l.RawControlStructures[controlStructure]; exists {
		l.rawEnd = re
	}
	return l.lexExpression
}

func (l *Lexer) lexLineComment() lexFn {
	for isLinePrefixSpace(l.peek()) {
		l.next()
	}
	l.Pos += len(l.Config.LineCommentPrefix)
	for {
		r := l.peek()
		if r == '\n' || r == rEOF {
			break
		}
		l.next()
	}
	l.Start = l.Pos
	return l.lexData
}

func (l *Lexer) lexBlockEnd() lexFn {
	control := l.consumeWhitespaceControl("-+")
	l.Pos += len(l.Config.BlockEndString)
	l.emit(BlockEnd)
	if l.Config.TrimBlocks && control != '+' {
		l.consumeFollowingNewline()
	}
	if l.rawEnd != nil {
		return l.lexRaw
	} else {
		return l.lexData
	}
}

func (l *Lexer) lexExpression() lexFn {
	if logging.Enabled() {
		log.WithFields(log.Fields{
			"pos":       l.Pos,
			"input":     l.Input,
			"remaining": l.remaining(),
		}).Trace("lexExpression")
	}
	for {
		if len(l.delimiters) == 0 && l.lineStatement {
			if endLength, ok := l.lineStatementEndLength(); ok {
				l.Pos += endLength
				l.lineStatement = false
				l.emit(BlockEnd)
				return l.lexData
			}
		}

		if !l.expectDelimiter(l.peek()) {
			if l.expressionEnd == VariableEnd && l.hasPrefix(l.Config.VariableEndString) {
				return l.lexVariableEnd
			}

			if l.expressionEnd == BlockEnd && l.hasPrefix(l.Config.BlockEndString) {
				return l.lexBlockEnd
			}
		}

		r := l.next()
		if logging.Enabled() {
			log.WithFields(log.Fields{"rune": r}).Trace("lexExpression")
		}
		switch {
		case isEOF(r):
			return l.lexEOF
		case isSpace(r):
			return l.lexSpace
		case isNumeric(r):
			l.backup()
			return l.lexNumber
		case r == '"' || r == '\'':
			l.backup()
			return l.lexString
		case r == ',':
			l.emit(Comma)
		case r == '|':
			l.emit(Pipe)
		case r == '+':
			if l.hasPrefix(l.Config.BlockEndString) {
				l.backup()
				return l.lexBlockEnd
			} else {
				l.emit(Addition)
			}
		case r == '-':
			if l.hasPrefix(l.Config.BlockEndString) {
				l.backup()
				return l.lexBlockEnd
			} else if l.hasPrefix(l.Config.VariableEndString) {
				l.backup()
				return l.lexVariableEnd
			} else {
				l.emit(Subtraction)
			}
		case r == '~':
			l.emit(Tilde)
		case r == ':':
			l.emit(Colon)
		case r == '.':
			l.emit(Dot)
		case r == '%':
			l.emit(Modulo)
		case r == '/':
			if l.accept("/") {
				l.emit(FloorDivision)
			} else {
				l.emit(Division)
			}
		case r == '<':
			if l.accept("=") {
				l.emit(LowerThanOrEqual)
			} else {
				l.emit(LowerThan)
			}
		case r == '>':
			if l.accept("=") {
				l.emit(GreaterThanOrEqual)
			} else {
				l.emit(GreaterThan)
			}
		case r == '*':
			if l.accept("*") {
				l.emit(Power)
			} else {
				l.emit(Multiply)
			}
		case r == '!':
			if l.accept("=") {
				l.emit(Ne)
			} else {
				// l.emit(Not)
				l.errorf(`Unexpected "!"`)
			}
		case r == '=':
			if l.accept("=") {
				l.emit(Equals)
			} else {
				l.emit(Assign)
			}
		case r == '(':
			l.emit(LeftParenthesis)
			l.pushDelimiter(')')
		case r == '{':
			l.emit(LeftBrace)
			l.pushDelimiter('}')
		case r == '[':
			l.emit(LeftBracket)
			l.pushDelimiter(']')
		case r == ')':
			if !l.popDelimiter(')') {
				return nil
			}
			l.emit(RightParenthesis)
		case r == '}':
			if !l.popDelimiter('}') {
				return nil
			}
			l.emit(RightBrace)
		case r == ']':
			if !l.popDelimiter(']') {
				return nil
			}
			l.emit(RightBracket)
		// in
		case r == 'i' && l.accept("n"):
			if !isSpace(l.peek()) {
				return l.lexIdentifier
			}
			l.emit(In)
		// is
		case r == 'i' && l.accept("s"):
			if !isSpace(l.peek()) {
				return l.lexIdentifier
			}
			l.emit(Is)
		// and
		case r == 'a' && l.accept("n"):
			if !(l.accept("d") && isSpace(l.peek())) {
				return l.lexIdentifier
			}
			l.emit(And)
		// or
		case r == 'o' && l.accept("r"):
			if !isSpace(l.peek()) {
				return l.lexIdentifier
			}
			l.emit(Or)
		// not
		case r == 'n' && l.accept("o"):
			if !(l.accept("t") && (isSpace(l.peek()) || l.peek() == '(')) {
				return l.lexIdentifier
			}
			l.emit(Not)
		case isAlphaNumeric(r):
			return l.lexIdentifier
		}
	}
}

func (l *Lexer) lexSpace() lexFn {
	for isSpace(l.peek()) {
		l.next()
	}
	l.emit(Whitespace)
	return l.lexExpression
}

func (l *Lexer) lexEOF() lexFn {
	if len(l.delimiters) > 0 {
		last := len(l.delimiters) - 1
		expected := l.delimiters[last]
		l.errorf(`Unbalanced delimiters, expected "%c", got EOF`, expected)
	}
	l.emit(EOF)
	return nil
}

func (l *Lexer) nextIdentifier() string {
	for {
		switch r := l.next(); {
		case isAlphaNumeric(r):
			// absorb.
		default:
			l.backup()
			// l.emit(Name)
			return l.Current()
		}
	}
}

func (l *Lexer) lexIdentifier() lexFn {
	l.nextIdentifier()
	l.emit(Name)
	return l.lexExpression
}

func (l *Lexer) lexNumber() lexFn {
	tokType, err := l.scanNumber()
	if err != nil {
		return l.errorf("%s", err)
	}
	if tokType == Integer && isAlphaNumeric(l.peek()) {
		return l.lexIdentifier
	}
	l.emit(tokType)
	return l.lexExpression
}

func (l *Lexer) scanNumber() (Type, error) {
	if !isNumeric(l.peek()) {
		return Integer, fmt.Errorf("invalid numeric token")
	}

	l.next()
	if l.Current() == "0" {
		switch unicode.ToLower(l.peek()) {
		case 'b':
			l.next()
			hasDigits, err := l.scanDigits(isBinaryDigit, false)
			if err != nil {
				return Integer, err
			}
			if !hasDigits {
				return Integer, fmt.Errorf("invalid numeric token")
			}
			return Integer, nil
		case 'o':
			l.next()
			hasDigits, err := l.scanDigits(isOctalDigit, false)
			if err != nil {
				return Integer, err
			}
			if !hasDigits {
				return Integer, fmt.Errorf("invalid numeric token")
			}
			return Integer, nil
		case 'x':
			l.next()
			hasDigits, err := l.scanDigits(isHexDigit, false)
			if err != nil {
				return Integer, err
			}
			if !hasDigits {
				return Integer, fmt.Errorf("invalid numeric token")
			}
			return Integer, nil
		}
	}

	if _, err := l.scanDigits(isNumeric, true); err != nil {
		return Integer, err
	}

	tokType := Integer
	if l.peek() == '.' {
		l.next()
		switch next := l.peek(); {
		case isNumeric(next):
			tokType = Float
			if _, err := l.scanDigits(isNumeric, false); err != nil {
				return tokType, err
			}
		case isSpace(next) || isNumericTerminator(next) || isEOF(next):
			tokType = Float
		default:
			l.backup()
		}
	}

	hasExponent, err := l.scanExponent()
	if err != nil {
		return Float, err
	}
	if hasExponent {
		tokType = Float
	}

	if tokType == Float && l.peek() == '.' {
		l.next()
		next := l.peek()
		l.backup()
		if isNumeric(next) || isSpace(next) || isNumericTerminator(next) || isEOF(next) {
			return tokType, fmt.Errorf("two dots in numeric token")
		}
	}

	return tokType, nil
}

func (l *Lexer) scanDigits(valid func(rune) bool, sawDigit bool) (bool, error) {
	underscorePending := false
	for {
		switch r := l.peek(); {
		case valid(r):
			l.next()
			sawDigit = true
			underscorePending = false
		case r == '_':
			if !sawDigit || underscorePending {
				return sawDigit, fmt.Errorf("invalid numeric token")
			}
			l.next()
			underscorePending = true
		default:
			if underscorePending {
				return sawDigit, fmt.Errorf("invalid numeric token")
			}
			return sawDigit, nil
		}
	}
}

func (l *Lexer) scanExponent() (bool, error) {
	if next := l.peek(); next != 'e' && next != 'E' {
		return false, nil
	}

	savedPos := l.Pos
	savedWidth := l.Width

	l.next()
	if sign := l.peek(); sign == '+' || sign == '-' {
		l.next()
	}

	hasDigits, err := l.scanDigits(isNumeric, false)
	if err != nil || !hasDigits {
		l.Pos = savedPos
		l.Width = savedWidth
		return false, nil
	}
	return true, nil
}

func unescape(str string) string {
	str = str[1 : len(str)-1]
	for escaped, unescaped := range escapedStrings {
		str = strings.ReplaceAll(str, escaped, unescaped)
	}
	return str
}

func (l *Lexer) lexString() lexFn {
	quote := l.next() // should be either ' or "
	var prev rune
	near := make([]rune, 0)
	for r := l.next(); r != quote || prev == '\\'; r, prev = l.next(), r {
		// only keep near context of the current line in case of error
		if (len(near) == 0 || near[len(near)-1] != '\n') && r != rEOF {
			near = append(near, r)
		}
		if r == rEOF {
			return l.errorf(`%s`, string(near))
		}
	}
	l.processAndEmit(String, func(value string) string {
		return normalizeNewlines(unescape(value), l.Config.NewlineSequence)
	})
	return l.lexExpression
}

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

func isLinePrefixSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

func isNumeric(r rune) bool {
	return unicode.IsDigit(r)
}

func isBinaryDigit(r rune) bool {
	return r == '0' || r == '1'
}

func isOctalDigit(r rune) bool {
	return r >= '0' && r <= '7'
}

func isHexDigit(r rune) bool {
	return isNumeric(r) || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
}

func isNumericTerminator(r rune) bool {
	switch r {
	case ',', '|', ':', ')', ']', '}', '+', '-', '*', '/', '%', '<', '>', '=', '!':
		return true
	default:
		return false
	}
}

func isEOF(r rune) bool {
	return r == rEOF
}
