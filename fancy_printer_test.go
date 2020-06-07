package fancy_printer

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
)

func TestBasicPrinting(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	for i := 1; i < 10; i++ {
		fmt.Println(i)
		tokens := generateTokens(i, 3, 20)
		PrintArrayAsGrid(tokens, false, false)
		fmt.Println()
	}

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	fmt.Println(string(out))
}

func generateTokens(numTokens, minTokenLength, maxTokenLength int) []string {
	// tokens := []string{"Ello there, govnah", "Bartholomew", "A", "B", "C"}
	tokens := make([]string, numTokens)
	for i := 0; i < numTokens; i++ {
		tokens[i] = generateToken(minTokenLength, maxTokenLength)
	}
	return tokens
}

func generateToken(minTokenLength, maxTokenLength int) string {
	tokenLength := rand.Intn(maxTokenLength-minTokenLength) + minTokenLength
	// tokenLength = 120
	token := ""
	for i := 0; i < tokenLength; i++ {
		token += "A"
	}
	return token
}
