//go:build !solution

package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Evaluator struct {
	stack []int
	words map[string][]string
}

func NewEvaluator() *Evaluator {
	return &Evaluator{
		stack: []int{},
		words: map[string][]string{},
	}
}

func (e *Evaluator) Process(row string) ([]int, error) {
	row = strings.ToLower(row)

	tokens := strings.Fields(row)
	if len(tokens) == 0 {
		return e.stack, nil
	}

	if tokens[0] == ":" {
		return e.stack, e.parseDefinition(tokens)
	}

	return e.stack, e.evalTokens(tokens, true)
}

func (e *Evaluator) parseDefinition(tokens []string) error {
	if len(tokens) < 4 || tokens[len(tokens)-1] != ";" {
		return errors.New("invalid definition syntax")
	}

	word := tokens[1]
	if _, err := strconv.Atoi(word); err == nil {
		return fmt.Errorf("cannot redefine number: %s", word)
	}

	var finalDef []string
	for _, t := range tokens[2 : len(tokens)-1] {
		if v, ok := e.words[t]; ok {
			finalDef = append(finalDef, v...)
		} else {
			finalDef = append(finalDef, t)
		}
	}

	e.words[word] = finalDef
	return nil
}

func (e *Evaluator) evalTokens(tokens []string, checkDef bool) error {
	for _, token := range tokens {
		if v, ok := e.words[token]; checkDef && ok {
			if err := e.evalTokens(v, false); err != nil {
				return err
			}
			continue
		}

		switch token {
		case "+", "-", "*", "/":
			if err := e.arithmetic(token); err != nil {
				return err
			}
		case "dup", "over", "drop", "swap":
			if err := e.stackOp(token); err != nil {
				return err
			}
		default:
			n, err := strconv.Atoi(token)
			if err != nil {
				return fmt.Errorf("undefined word: %s", token)
			}
			e.stack = append(e.stack, n)
		}
	}
	return nil
}

func (e *Evaluator) arithmetic(op string) error {
	n := len(e.stack)
	if n < 2 {
		return fmt.Errorf("stack underflow for operator %s", op)
	}
	a, b := e.stack[n-2], e.stack[n-1]
	var res int

	switch op {
	case "+":
		res = a + b
	case "-":
		res = a - b
	case "*":
		res = a * b
	case "/":
		if b == 0 {
			return errors.New("division by zero")
		}
		res = a / b
	}

	e.stack = append(e.stack[:n-2], res)
	return nil
}

func (e *Evaluator) stackOp(op string) error {
	n := len(e.stack)
	switch op {
	case "dup":
		if n < 1 {
			return errors.New("stack underflow on dup")
		}
		e.stack = append(e.stack, e.stack[n-1])
	case "drop":
		if n < 1 {
			return errors.New("stack underflow on drop")
		}
		e.stack = e.stack[:n-1]
	case "over":
		if n < 2 {
			return errors.New("stack underflow on over")
		}
		e.stack = append(e.stack, e.stack[n-2])
	case "swap":
		if n < 2 {
			return errors.New("stack underflow on swap")
		}
		e.stack[n-1], e.stack[n-2] = e.stack[n-2], e.stack[n-1]
	}
	return nil
}
