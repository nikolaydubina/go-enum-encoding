package main_test

import (
	"os"
	"os/exec"
	"path"
	"testing"
)

func TestMain(t *testing.T) {
	coverdir := t.TempDir()
	testbin := path.Join(t.TempDir(), "go-enum-encoding-test")
	exec.Command("go", "build", "-cover", "-o", testbin, "main.go").Run()
	defer exec.Command("go", "tool", "covdata", "textfmt", "-i="+coverdir, "-o", os.Getenv("GOCOVERPROFILE")).Run()

	cmd := exec.Command(testbin, "--type", "Color")
	cmd.Env = append(cmd.Environ(), "GOFILE=internal/testdata/color.go", "GOPACKAGE=main", "GOCOVERDIR="+coverdir)
	cmd.Run()

	assertEqFile(t, "internal/testdata/color_enum_encoding.go", "internal/testdata/exp/color_enum_encoding.go")
	assertEqFile(t, "internal/testdata/color_enum_encoding_test.go", "internal/testdata/exp/color_enum_encoding_test.go")

	t.Run("when bad go file, then error", func(t *testing.T) {
		cmd := exec.Command(testbin, "--type", "Color")
		cmd.Env = append(cmd.Environ(), "GOFILE=README.md", "GOPACKAGE=main", "GOCOVERDIR="+coverdir)
		if err := cmd.Run(); err == nil {
			t.Fatal("must be error")
		}
	})

	t.Run("when invalid package name, then error", func(t *testing.T) {
		cmd := exec.Command(testbin, "--type", "Color")
		cmd.Env = append(cmd.Environ(), "GOFILE=internal/testdata/color.go", "GOPACKAGE=m_aIn", "GOCOVERDIR="+coverdir)
		if err := cmd.Run(); err == nil {
			t.Fatal("must be error")
		}
	})

	t.Run("when not found file, then error", func(t *testing.T) {
		cmd := exec.Command(testbin, "--type", "Color")
		cmd.Env = append(cmd.Environ(), "GOFILE=asdf.asdf", "GOPACKAGE=main", "GOCOVERDIR="+coverdir)
		if err := cmd.Run(); err == nil {
			t.Fatal("must be error")
		}
	})

	t.Run("when wrong params, then error", func(t *testing.T) {
		cmd := exec.Command(testbin)
		cmd.Env = append(cmd.Environ(), "GOFILE=internal/testdata/color.go", "GOPACKAGE=main", "GOCOVERDIR="+coverdir)
		if err := cmd.Run(); err == nil {
			t.Fatal("must be error")
		}
	})

	t.Run("when no undefined, then error", func(t *testing.T) {
		cmd := exec.Command(testbin, "--type", "NoUndefined")
		cmd.Env = append(cmd.Environ(), "GOFILE=internal/testdata/no_undefined.go", "GOPACKAGE=main", "GOCOVERDIR="+coverdir)
		if err := cmd.Run(); err == nil {
			t.Fatal("must be error")
		}
	})

	t.Run("when undefined has wrong tag, then error", func(t *testing.T) {
		cmd := exec.Command(testbin, "--type", "BadUndefined1")
		cmd.Env = append(cmd.Environ(), "GOFILE=internal/testdata/wrong_tag_undefined.go", "GOPACKAGE=main", "GOCOVERDIR="+coverdir)
		if err := cmd.Run(); err == nil {
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
