package main

import "testing"

func TestisAlphaNumeric(t *testing.T) bool {
	type test struct {
		input           rune
		expected_result bool
	}

	for _, t := range []test{
		input:           'a',
		expected_result: true,
	} {
		res := isAlphaNumeric(t.input)
		if res != t.expected_result {
			t.Errorf("test error: unexpected result %v for input %v\n", res, t.input)
		}
	}
}
