package main

import "testing"

func TestisAlphaNumeric(t *testing.T) {
	type test struct {
		input           rune
		expected_result bool
	}

	for _, test := range []test{
		{input: 'a', expected_result: true},
		{input: 'b', expected_result: false},
	} {
		res := isAlphaNumeric(test.input)
		if res != test.expected_result {
			t.Errorf("test error: unexpected result %v for input %v\n", res, test.input)
		}
	}
}
