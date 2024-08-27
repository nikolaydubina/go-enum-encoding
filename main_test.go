package main_test

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path"
	"testing"
)

func FuzzBadFile(f *testing.F) {
	testbin := path.Join(f.TempDir(), "go-enum-encoding-test")
	exec.Command("go", "build", "-cover", "-o", testbin, "main.go").Run()

	f.Fuzz(func(t *testing.T, orig string) {
		t.Run("when bad go file, then error", func(t *testing.T) {
			fname := path.Join(t.TempDir(), "fuzz-test-file.go")
			os.WriteFile(fname, []byte(orig), 0644)

			cmd := exec.Command(testbin, "--type", "Color")
			cmd.Env = append(cmd.Environ(), "GOFILE="+fname, "GOPACKAGE=color")
			if err := cmd.Run(); err == nil {
				t.Fatal("must be error")
			}
		})
	})
}

func TestMain(t *testing.T) {
	coverdir := t.TempDir()
	testbin := path.Join(t.TempDir(), "go-enum-encoding-test")
	exec.Command("go", "build", "-cover", "-o", testbin, "main.go").Run()
	defer exec.Command("go", "tool", "covdata", "textfmt", "-i="+coverdir, "-o", os.Getenv("GOCOVERPROFILE")).Run()

	t.Run("when ok, then file matches expected", func(t *testing.T) {
		t.Run("when short mode, then file matches expected", func(t *testing.T) {
			cmd := exec.Command(testbin, "--type", "Color")
			cmd.Env = append(cmd.Environ(), "GOFILE=internal/testdata/color.go", "GOLINE=5", "GOPACKAGE=color", "GOCOVERDIR="+coverdir)
			cmd.Run()

			assertEqFile(t, "internal/testdata/color_enum_encoding.go", "internal/testdata/exp/color_enum_encoding.go")
			assertEqFile(t, "internal/testdata/color_enum_encoding_test.go", "internal/testdata/exp/color_enum_encoding_test.go")
		})

		t.Run("when auto mode, then long can be detected and file matches expected", func(t *testing.T) {
			cmd := exec.Command(testbin, "--type", "Currency")
			cmd.Env = append(cmd.Environ(), "GOFILE=internal/testdata/currency.go", "GOLINE=5", "GOPACKAGE=color", "GOCOVERDIR="+coverdir)
			cmd.Run()

			assertEqFile(t, "internal/testdata/currency_enum_encoding.go", "internal/testdata/exp/currency_enum_encoding.go")
			assertEqFile(t, "internal/testdata/currency_enum_encoding_test.go", "internal/testdata/exp/currency_enum_encoding_test.go")
		})

		t.Run("string", func(t *testing.T) {
			t.Run("short", func(t *testing.T) {
				var errb bytes.Buffer

				cmd := exec.Command(testbin, "--type", "ColorString", "--string")
				cmd.Env = append(cmd.Environ(), "GOFILE=internal/testdata/color.go", "GOLINE=20", "GOPACKAGE=color", "GOCOVERDIR="+coverdir)
				cmd.Stderr = &errb
				cmd.Run()

				assertEqFile(t, "internal/testdata/colorstring_enum_encoding.go", "internal/testdata/exp/colorstring_enum_encoding.go")
				assertEqFile(t, "internal/testdata/colorstring_enum_encoding_test.go", "internal/testdata/exp/colorstring_enum_encoding_test.go")
				if errb, err := io.ReadAll(&errb); len(errb) > 0 || err != nil {
					t.Errorf("stderr: %s", string(errb))
				}
			})

			t.Run("long", func(t *testing.T) {
				cmd := exec.Command(testbin, "--type", "CurrencyString", "--mode", "long", "--string")
				cmd.Env = append(cmd.Environ(), "GOFILE=internal/testdata/currency_string.go", "GOLINE=5", "GOPACKAGE=color", "GOCOVERDIR="+coverdir)
				cmd.Run()

				assertEqFile(t, "internal/testdata/currencystring_enum_encoding.go", "internal/testdata/exp/currencystring_enum_encoding.go")
				assertEqFile(t, "internal/testdata/currencystring_enum_encoding_test.go", "internal/testdata/exp/currencystring_enum_encoding_test.go")
			})

			t.Run("custom method", func(t *testing.T) {
				cmd := exec.Command(
					testbin,
					"--type", "CurrencyStringCustom",
					"--mode", "long",
					"--string",
					"--encode-method", "MarshalTextName",
					"--decode-method", "UnmarshalTextName",
					"--string-method", "StringName",
				)
				cmd.Env = append(cmd.Environ(), "GOFILE=internal/testdata/currency_string_custom.go", "GOLINE=5", "GOPACKAGE=color", "GOCOVERDIR="+coverdir)
				cmd.Run()

				assertEqFile(t, "internal/testdata/currencystringcustom_enum_encoding.go", "internal/testdata/exp/currencystringcustom_enum_encoding.go")
				assertEqFile(t, "internal/testdata/currencystringcustom_enum_encoding_test.go", "internal/testdata/exp/currencystringcustom_enum_encoding_test.go")
			})
		})

		t.Run("when multiple enums in same file, then file matches expected for each", func(t *testing.T) {
			cmd := exec.Command(testbin, "--type", "Color2")
			cmd.Env = append(cmd.Environ(), "GOFILE=internal/testdata/multiple.go", "GOLINE=5", "GOPACKAGE=color", "GOCOVERDIR="+coverdir)
			cmd.Run()

			assertEqFile(t, "internal/testdata/color2_enum_encoding.go", "internal/testdata/exp/color2_enum_encoding.go")
			assertEqFile(t, "internal/testdata/color2_enum_encoding_test.go", "internal/testdata/exp/color2_enum_encoding_test.go")

			cmd = exec.Command(testbin, "--type", "Currency2")
			cmd.Env = append(cmd.Environ(), "GOFILE=internal/testdata/multiple.go", "GOLINE=14", "GOPACKAGE=color", "GOCOVERDIR="+coverdir)
			cmd.Run()

			assertEqFile(t, "internal/testdata/currency2_enum_encoding.go", "internal/testdata/exp/currency2_enum_encoding.go")
			assertEqFile(t, "internal/testdata/currency2_enum_encoding_test.go", "internal/testdata/exp/currency2_enum_encoding_test.go")
		})
	})

	t.Run("when bad go file, then error", func(t *testing.T) {
		cmd := exec.Command(testbin, "--type", "Color")
		cmd.Env = append(cmd.Environ(), "GOFILE=README.md", "GOLINE=5", "GOPACKAGE=color", "GOCOVERDIR="+coverdir)
		if err := cmd.Run(); err == nil {
			t.Fatal("must be error")
		}
	})

	t.Run("when enum values not immediately after go:generate line, then error", func(t *testing.T) {
		cmd := exec.Command(testbin, "--type", "Color")
		cmd.Env = append(cmd.Environ(), "GOFILE=internal/testdata/color.go", "GOLINE=1", "GOPACKAGE=color", "GOCOVERDIR="+coverdir)
		if err := cmd.Run(); err == nil {
			t.Fatal("must be error")
		}
	})

	t.Run("when invalid package name, then error", func(t *testing.T) {
		cmd := exec.Command(testbin, "--type", "Color")
		cmd.Env = append(cmd.Environ(), "GOFILE=internal/testdata/color.go", "GOLINE=5", "GOPACKAGE=\"", "GOCOVERDIR="+coverdir)
		if err := cmd.Run(); err == nil {
			t.Fatal("must be error")
		}
	})

	t.Run("when not found file, then error", func(t *testing.T) {
		cmd := exec.Command(testbin, "--type", "Color")
		cmd.Env = append(cmd.Environ(), "GOFILE=asdf.asdf", "GOPACKAGE=color", "GOLINE=5", "GOCOVERDIR="+coverdir)
		if err := cmd.Run(); err == nil {
			t.Fatal("must be error")
		}
	})

	t.Run("when wrong params, then error", func(t *testing.T) {
		cmd := exec.Command(testbin)
		cmd.Env = append(cmd.Environ(), "GOFILE=internal/testdata/color.go", "GOLINE=5", "GOPACKAGE=color", "GOCOVERDIR="+coverdir)
		if err := cmd.Run(); err == nil {
			t.Fatal("must be error")
		}
	})
}

func assertEqFile(t *testing.T, a, b string) {
	fa, _ := os.ReadFile(a)
	fb, _ := os.ReadFile(b)
	if string(fa) != string(fb) {
		t.Error("files are different (" + a + " <> " + b + ")")
	}
}
