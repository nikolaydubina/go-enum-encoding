package main_test

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
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
	coverdir, outdir := t.TempDir(), t.TempDir()
	testbin := path.Join(t.TempDir(), "go-enum-encoding-test")
	exec.Command("go", "build", "-cover", "-o", testbin, "main.go").Run()
	defer exec.Command("go", "tool", "covdata", "textfmt", "-i="+coverdir, "-o", os.Getenv("GOCOVERPROFILE")).Run()

	assertFile := func(t *testing.T, a string) {
		t.Helper()
		assertEqFile(t, filepath.Join(outdir, a), filepath.Join("internal", "exp", a))
	}

	t.Run("when ok, then file matches expected", func(t *testing.T) {
		t.Run("when short mode, then file matches expected", func(t *testing.T) {
			exec.Command("cp", filepath.Join("internal", "color.go"), filepath.Join(outdir, "color.go")).Run()

			cmd := exec.Command(testbin, "--type", "Color")
			cmd.Env = append(cmd.Environ(), "GOFILE="+filepath.Join(outdir, "color.go"), "GOLINE=4", "GOPACKAGE=color", "GOCOVERDIR="+coverdir)
			cmd.Run()

			assertFile(t, "color_enum_encoding.go")
			assertFile(t, "color_enum_encoding_test.go")
			assertFile(t, "color_enum_encoding_json_test.go")
		})

		t.Run("when auto mode, then long can be detected and file matches expected", func(t *testing.T) {
			exec.Command("cp", filepath.Join("internal", "currency.go"), filepath.Join(outdir, "currency.go")).Run()

			cmd := exec.Command(testbin, "--type", "Currency")
			cmd.Env = append(cmd.Environ(), "GOFILE="+filepath.Join(outdir, "currency.go"), "GOLINE=4", "GOPACKAGE=color", "GOCOVERDIR="+coverdir)
			cmd.Run()

			assertFile(t, "currency_enum_encoding.go")
			assertFile(t, "currency_enum_encoding_test.go")
			assertFile(t, "currency_enum_encoding_json_test.go")
		})

		t.Run("string", func(t *testing.T) {
			t.Run("short", func(t *testing.T) {
				exec.Command("cp", filepath.Join("internal", "color.go"), filepath.Join(outdir, "color.go")).Run()

				cmd := exec.Command(testbin, "--type", "ColorString", "--string")
				cmd.Env = append(cmd.Environ(), "GOFILE="+filepath.Join(outdir, "color.go"), "GOLINE=18", "GOPACKAGE=color", "GOCOVERDIR="+coverdir)
				cmd.Run()

				assertFile(t, "colorstring_enum_encoding.go")
				assertFile(t, "colorstring_enum_encoding_test.go")
				assertFile(t, "colorstring_enum_encoding_string_test.go")
				assertFile(t, "colorstring_enum_encoding_json_test.go")
			})

			t.Run("long", func(t *testing.T) {
				exec.Command("cp", filepath.Join("internal", "currency_string.go"), filepath.Join(outdir, "currency_string.go")).Run()

				cmd := exec.Command(testbin, "--type", "CurrencyString", "--mode", "long", "--string")
				cmd.Env = append(cmd.Environ(), "GOFILE="+filepath.Join(outdir, "currency_string.go"), "GOLINE=4", "GOPACKAGE=color", "GOCOVERDIR="+coverdir)
				cmd.Run()

				assertFile(t, "currencystring_enum_encoding.go")
				assertFile(t, "currencystring_enum_encoding_test.go")
				assertFile(t, "currencystring_enum_encoding_string_test.go")
				assertFile(t, "currencystring_enum_encoding_json_test.go")
			})

			t.Run("custom method", func(t *testing.T) {
				exec.Command("cp", filepath.Join("internal", "currency_string_custom.go"), filepath.Join(outdir, "currency_string_custom.go")).Run()

				cmd := exec.Command(
					testbin,
					"--type", "CurrencyStringCustom",
					"--mode", "long",
					"--string",
					"--encode-method", "MarshalTextName",
					"--decode-method", "UnmarshalTextName",
					"--string-method", "StringName",
				)
				cmd.Env = append(cmd.Environ(), "GOFILE="+filepath.Join(outdir, "currency_string_custom.go"), "GOLINE=4", "GOPACKAGE=color", "GOCOVERDIR="+coverdir)
				cmd.Run()

				assertFile(t, "currencystringcustom_enum_encoding.go")
				assertFile(t, "currencystringcustom_enum_encoding_test.go")
			})
		})

		t.Run("when multiple enums in same file, then file matches expected for each", func(t *testing.T) {
			exec.Command("cp", filepath.Join("internal", "multiple.go"), filepath.Join(outdir, "multiple.go")).Run()

			cmd := exec.Command(testbin, "--type", "Color2")
			cmd.Env = append(cmd.Environ(), "GOFILE="+filepath.Join(outdir, "multiple.go"), "GOLINE=4", "GOPACKAGE=color", "GOCOVERDIR="+coverdir)
			cmd.Run()

			assertFile(t, "color2_enum_encoding.go")
			assertFile(t, "color2_enum_encoding_test.go")
			assertFile(t, "color2_enum_encoding_string_test.go")
			assertFile(t, "color2_enum_encoding_json_test.go")

			cmd = exec.Command(testbin, "--type", "Currency2")
			cmd.Env = append(cmd.Environ(), "GOFILE="+filepath.Join(outdir, "multiple.go"), "GOLINE=12", "GOPACKAGE=color", "GOCOVERDIR="+coverdir)
			cmd.Run()

			assertFile(t, "currency2_enum_encoding.go")
			assertFile(t, "currency2_enum_encoding_test.go")
			assertFile(t, "currency2_enum_encoding_string_test.go")
			assertFile(t, "currency2_enum_encoding_json_test.go")
		})
	})

	t.Run("when bad go file, then error", func(t *testing.T) {
		exec.Command("cp", filepath.Join("internal", "README.md"), filepath.Join(outdir, "README.md")).Run()

		cmd := exec.Command(testbin, "--type", "Color")
		cmd.Env = append(cmd.Environ(), "GOFILE=README.md", "GOLINE=5", "GOPACKAGE=color", "GOCOVERDIR="+coverdir)
		if err := cmd.Run(); err == nil {
			t.Fatal("must be error")
		}
	})

	t.Run("when enum values not immediately after go:generate line, then error", func(t *testing.T) {
		cmd := exec.Command(testbin, "--type", "Color")
		cmd.Env = append(cmd.Environ(), "GOFILE="+filepath.Join(outdir, "color.go"), "GOLINE=1", "GOPACKAGE=color", "GOCOVERDIR="+coverdir)
		if err := cmd.Run(); err == nil {
			t.Fatal("must be error")
		}
	})

	t.Run("when invalid package name, then error", func(t *testing.T) {
		cmd := exec.Command(testbin, "--type", "Color")
		cmd.Env = append(cmd.Environ(), "GOFILE="+filepath.Join(outdir, "color.go"), "GOLINE=5", "GOPACKAGE=\"", "GOCOVERDIR="+coverdir)
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
		cmd.Env = append(cmd.Environ(), "GOFILE="+filepath.Join(outdir, "color.go"), "GOLINE=5", "GOPACKAGE=color", "GOCOVERDIR="+coverdir)
		if err := cmd.Run(); err == nil {
			t.Fatal("must be error")
		}
	})
}

func assertEqFile(t *testing.T, a, b string) {
	fa, _ := os.ReadFile(a)
	fb, _ := os.ReadFile(b)
	if string(fa) != string(fb) {
		t.Error("files are different (" + a + " <> " + b + "), " + string(fa) + " <> " + string(fb))
	}
}
