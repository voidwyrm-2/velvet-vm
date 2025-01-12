package sectioner

import (
	"github.com/voidwyrm-2/velvet-vm/velvc/lexer/tokens"
)

func SectionIntoLines(unsectionsTokens []tokens.Token) [][]tokens.Token {
	toks := [][]tokens.Token{}
	acc := []tokens.Token{}
	currentLine := 1

	for _, t := range unsectionsTokens {
		if t.GetLn() == currentLine {
			acc = append(acc, t)
		} else {
			toks = append(toks, acc)
			acc = []tokens.Token{t}
			currentLine = t.GetLn()
		}
	}

	if len(acc) > 0 {
		toks = append(toks, acc)
	}

	return toks
}
