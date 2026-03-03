package main

import (
	"crypto/sha256"
	"fmt"
)

func hex(b []byte) string {
	s := ""
	for _, v := range b {
		s += fmt.Sprintf("%02x", v)
	}
	return s
}

func main() {
	tests := []struct {
		input    string
		expected string
	}{
		{"", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{"abc", "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"},
		{"Hello, SHA-256!", "3d61375cc279a6fe4c887bc4eee58ec78bc32cf2b026252918b9b0d74eaffeb5"},
	}

	pass := 0
	for _, t := range tests {
		h := sha256.Sum256([]byte(t.input))
		got := hex(h[:])
		if got == t.expected {
			pass++
			fmt.Printf("PASS: sha256(%q)\n", t.input)
		} else {
			fmt.Printf("FAIL: sha256(%q) = %s (expected %s)\n", t.input, got, t.expected)
		}
	}

	// Incremental hashing
	h := sha256.New()
	h.Write([]byte("Hello, "))
	h.Write([]byte("SHA-256!"))
	incremental := hex(h.Sum(nil))
	expected := "3d61375cc279a6fe4c887bc4eee58ec78bc32cf2b026252918b9b0d74eaffeb5"
	if incremental == expected {
		pass++
		fmt.Println("PASS: incremental hash")
	} else {
		fmt.Printf("FAIL: incremental = %s\n", incremental)
	}

	fmt.Printf("sha256 tests: %d/4 passed\n", pass)
	if pass == 4 {
		fmt.Println("result: OK")
	} else {
		fmt.Println("result: FAIL")
	}
}
