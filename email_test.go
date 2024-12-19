package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_formatEmailName(t *testing.T) {

	testCases := map[string]string{
		"Bingus Online - Signup":   "bingusonline",
		"Bingus Online - Register": "bingusonline",
		"Bingus Online":            "bingusonline",
		"Bingus! Register":         "bingus",
	}

	for c, expected := range testCases {
		result := formatEmailName(c)
		fmt.Println(c, "->", result)
		assert.Equal(t, expected, result)
	}
}

func Test_emailHash(t *testing.T) {
	testCases := map[string]string{
		"Bingus Online - Signup":   "2b",
		"Bingus Online - Register": "2e",
		"Bingus Online":            "aa",
		"Bingus! Register":         "3f",
	}

	for c, expected := range testCases {
		result := emailHash(c)
		fmt.Println(c, "->", result)
		assert.Equal(t, expected, result)
	}
}
