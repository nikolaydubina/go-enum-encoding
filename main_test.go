package main

import (
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	if err := process("Color", "internal/testdata/color.go", "main"); err != nil {
		t.Fatal(err)
	}

	assertEqFile(t, "internal/testdata/color_enum_encoding.go", "internal/testdata/exp/color_enum_encoding.go")
	assertEqFile(t, "internal/testdata/color_enum_encoding_test.go", "internal/testdata/exp/color_enum_encoding_test.go")

	t.Run("when wrong params, then error", func(t *testing.T) {
		if err := process("", "internal/testdata/color.go", "main"); err == nil {
			t.Fatal("must be error")
		}
	})

	t.Run("when no undefined, then error", func(t *testing.T) {
		if err := process("NoUndefined", "internal/testdata/no_undefined.go", "main"); err == nil {
			t.Fatal("must be error")
		}
	})

	t.Run("when undefined has wrong tag, then error", func(t *testing.T) {
		if err := process("BadUndefined1", "internal/testdata/wrong_tag_undefined.go", "main"); err == nil {
			t.Fatal("must be error")
		}
	})
}

func assertEqFile(t *testing.T, a, b string) {
	fa, _ := os.ReadFile(a)
	fb, _ := os.ReadFile(b)
	if string(fa) != string(fb) {
		t.Error("files are different")
	}
}
