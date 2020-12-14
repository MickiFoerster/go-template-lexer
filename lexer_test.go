package main

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestIsAlphaNumeric(t *testing.T) {
	type test struct {
		input           rune
		expected_result bool
	}

	for _, test := range []test{
		{input: 'a', expected_result: true},
		{input: 'b', expected_result: true},
		{input: 'c', expected_result: true},
		{input: 'd', expected_result: true},
		{input: 'e', expected_result: true},
		{input: 'f', expected_result: true},
		{input: 'g', expected_result: true},
		{input: 'h', expected_result: true},
		{input: 'i', expected_result: true},
		{input: 'j', expected_result: true},
		{input: 'k', expected_result: true},
		{input: 'l', expected_result: true},
		{input: 'm', expected_result: true},
		{input: 'n', expected_result: true},
		{input: 'o', expected_result: true},
		{input: 'p', expected_result: true},
		{input: 'q', expected_result: true},
		{input: 'r', expected_result: true},
		{input: 's', expected_result: true},
		{input: 't', expected_result: true},
		{input: 'u', expected_result: true},
		{input: 'v', expected_result: true},
		{input: 'w', expected_result: true},
		{input: 'x', expected_result: true},
		{input: 'y', expected_result: true},
		{input: 'z', expected_result: true},
		{input: '0', expected_result: true},
		{input: '1', expected_result: true},
		{input: '2', expected_result: true},
		{input: '3', expected_result: true},
		{input: '4', expected_result: true},
		{input: '5', expected_result: true},
		{input: '6', expected_result: true},
		{input: '7', expected_result: true},
		{input: '8', expected_result: true},
		{input: '9', expected_result: true},
	} {
		t.Logf("Testing %v\n", string(test.input))
		res := isAlphaNumeric(test.input)
		if res != test.expected_result {
			t.Errorf("test error: unexpected result %v for input %v\n", res, string(test.input))
		}
		if 'a' <= test.input && test.input <= 'z' {
			upper := strings.ToUpper(string(test.input))
			inp, _ := utf8.DecodeRuneInString(upper)
			t.Logf("Testing %v\n", string(inp))
			res = isAlphaNumeric(inp)
			if res != test.expected_result {
				t.Errorf("test error: unexpected result %v for input %v\n", res, string(test.input))
			}
		}
	}
}
