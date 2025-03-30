package scanner

import (
	"github.com/demouth/mini-typescript-go/internal/ast"
	"github.com/demouth/mini-typescript-go/internal/jsnum"
	"github.com/demouth/mini-typescript-go/internal/stringutil"
)

type ScannerState struct {
	pos          int // Current position in text (and ending position of current token)
	fullStartPos int // Starting position of current token including preceding whitespace
	token        ast.Kind
	tokenValue   string
}

type Scanner struct {
	text string

	ScannerState
}

func NewScanner() *Scanner {
	return &Scanner{}
}

func (s *Scanner) TokenFullStart() int {
	return s.fullStartPos
}

func (s *Scanner) TokenValue() string {
	return s.tokenValue
}

func (s *Scanner) SetText(text string) {
	s.text = text
	s.ScannerState = ScannerState{}
}

func (s *Scanner) char() rune {
	if s.pos < len(s.text) {
		return rune(s.text[s.pos])
	}
	return -1
}

func (s *Scanner) Scan() ast.Kind {
	s.fullStartPos = s.pos
	for {
		ch := s.char()
		switch ch {
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			s.token = s.scanNumber()
		default:
			if ch < 0 {
				s.token = ast.KindEndOfFile
				break
			}
		}
		return s.token
	}
}

func (s *Scanner) scanNumber() ast.Kind {
	start := s.pos
	if s.char() == '0' {
		s.pos++
		return ast.KindNumericLiteral
	} else {
		s.scanNumberFragment()
	}
	end := s.pos
	s.tokenValue = s.text[start:end]

	var result ast.Kind
	result = s.scanBigIntSuffix()
	return result
}

func (s *Scanner) scanNumberFragment() string {
	start := s.pos
	for {
		ch := s.char()
		if stringutil.IsDigit(ch) {
			s.pos++
			continue
		}
		break
	}
	return s.text[start:s.pos]
}

func (s *Scanner) scanBigIntSuffix() ast.Kind {
	s.tokenValue = jsnum.FromString(s.tokenValue).String()
	return ast.KindNumericLiteral
}
