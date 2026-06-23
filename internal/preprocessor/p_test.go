package preprocessor

import (
	"fmt"
	"testing"
)

func TestDefine(t *testing.T) {
	a := New()
	var b = []byte(`
A

#define A = 10

A
B

#define B = hello

B`)

	a.Define(&b)

	fmt.Println(string(b))
}
