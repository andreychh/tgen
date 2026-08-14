// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package types

import (
	"errors"
	"fmt"
	"strings"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/prose"
	"github.com/andreychh/tgen/model/typeexpr"
)

// Lexer is a tokenizer over the prose of a type cell, recognizing built-in type
// words through a primitive vocabulary.
type Lexer struct {
	phrase     prose.Phrase
	primitives typeexpr.Primitives
}

// NewLexer constructs a Lexer over the prose of a type cell, using the default
// primitive vocabulary.
func NewLexer(phrase prose.Phrase) Lexer {
	return Lexer{phrase: phrase, primitives: typeexpr.NewPrimitives()}
}

// Tokens returns the tokens the phrase lexes to, in the order the prose writes
// them. It fails on a link that is not an anchor, an inline that is neither
// text nor a link, an "Array" that "of" does not follow, and a word that names
// no built-in type.
func (l Lexer) Tokens() ([]Token, error) {
	var out []Token
	for _, node := range l.phrase.Inlines() {
		tokens, err := l.inline(node)
		if err != nil {
			return nil, err
		}
		out = append(out, tokens...)
	}
	return out, nil
}

// inline returns the tokens node lexes to. It fails on a link that is not an
// anchor and on an inline that is neither text nor a link.
func (l Lexer) inline(node prose.Inline) ([]Token, error) {
	switch node := node.(type) {
	case prose.Link:
		target, ok := node.Anchor()
		if !ok {
			return nil, fmt.Errorf("type link %q is not an anchor", node.Href())
		}
		return []Token{NewRef(model.Reference(target))}, nil
	case prose.Text:
		return l.words(node.Content())
	default:
		return nil, fmt.Errorf("unexpected inline %T in type expression", node)
	}
}

// words returns the structural and primitive tokens the content of a text run
// lexes to. It fails on an "Array" that "of" does not follow and on a word that
// names no built-in type.
func (l Lexer) words(content string) ([]Token, error) {
	var out []Token
	cursor := NewCursor(split(content))
	for {
		word, ok := cursor.Take()
		if !ok {
			break
		}
		switch word {
		case "Array":
			next, found := cursor.Take()
			if !found || next != "of" {
				return nil, errors.New(`expected "of" after "Array"`)
			}
			out = append(out, NewArrayOf())
		case "or":
			out = append(out, NewOr())
		case "and", ",":
			out = append(out, NewSeries())
		default:
			kind, known := l.primitives.Kind(word)
			if !known {
				return nil, fmt.Errorf("unknown type word %q", word)
			}
			out = append(out, NewPrimitive(kind))
		}
	}
	return out, nil
}

// split breaks a text run into words, treating each comma as its own word.
func split(content string) []string {
	return strings.Fields(strings.ReplaceAll(content, ",", " , "))
}
