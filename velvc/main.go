package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/voidwyrm-2/velvet-vm/velvc/lexer"
	"github.com/voidwyrm-2/velvet-vm/velvc/sectioner"
)

func readFile(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	return string(b), err
}

func main() {
	showTokens := flag.Bool("tokens", false, "Print the generated tokens")
	showSectioned := flag.Bool("sectioned", false, "Print the sectioned tokens")

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("expected 'velvc <file>'")
		os.Exit(1)
	}

	content, err := readFile(args[0])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	lexer := lexer.New(strings.TrimSpace(content))

	tokens, err := lexer.Lex()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	} else if *showTokens {
		for _, t := range tokens {
			fmt.Println(t.Str())
		}
	}

	sectionedTokens := sectioner.SectionIntoLines(tokens)
	if *showSectioned {
		formatted := []string{}
		for _, tl := range sectionedTokens {
			line := []string{}
			for _, t := range tl {
				line = append(line, t.Str())
			}
			formatted = append(formatted, "["+strings.Join(line, " ")+"]")
		}
		fmt.Println(strings.Join(formatted, "\n"))
	}
}
