package main_test

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"
)

func FuzzBadFile(f *testing.F) {
	testdir := f.TempDir()
	testbin := path.Join(testdir, "go-enum-encoding-test")
	exec.Command("go", "build", "-cover", "-o", testbin, "main.go").Run()

	f.Fuzz(func(t *testing.T, orig string) {
		t.Run("when bad go file, then error", func(t *testing.T) {
			fname := path.Join(testdir, "fuzz-test-file.go")
			os.WriteFile(fname, []byte(orig), 0644)

			cmd := exec.Command(testbin, "--type", "Color")
			cmd.Env = append(cmd.Environ(), "GOFILE="+fname, "GOPACKAGE=image")
			if err := cmd.Run(); err == nil {
				t.Fatal("must be error")
			}
		})
	})
}

func TestMain(t *testing.T) {
	testdir := t.TempDir()
	testbin := path.Join(testdir, "go-enum-encoding-test")
	exec.Command("go", "build", "-cover", "-o", testbin, "main.go").Run()
	defer exec.Command("go", "tool", "covdata", "textfmt", "-i="+testdir, "-o", os.Getenv("GOCOVERPROFILE")).Run()

	t.Run("ok", func(t *testing.T) {
		exec.Command("cp", filepath.Join("internal", "testdata", "image.go"), filepath.Join(testdir, "image.go")).Run()

		t.Run("struct", func(t *testing.T) {
			cmd := exec.Command(testbin, "--type", "Color")
			cmd.Env = append(cmd.Environ(), "GOFILE="+filepath.Join(testdir, "image.go"), "GOLINE=4", "GOPACKAGE=image", "GOTESTDIR="+testdir)
			if err := cmd.Run(); err != nil {
				t.Error(err)
			}

			assertEqFile(t, filepath.Join(testdir, "color_enum_encoding.go"), filepath.Join("internal", "testdata", "color_enum_encoding.go"))
			assertEqFile(t, filepath.Join(testdir, "color_enum_encoding_test.go"), filepath.Join("internal", "testdata", "color_enum_encoding_test.go"))
		})

		t.Run("iota, string", func(t *testing.T) {
			cmd := exec.Command(testbin, "--type", "ImageSize", "--string")
			cmd.Env = append(cmd.Environ(), "GOFILE="+filepath.Join(testdir, "image.go"), "GOLINE=18", "GOPACKAGE=image", "GOTESTDIR="+testdir)
			if err := cmd.Run(); err != nil {
				t.Error(err)
			}

			assertEqFile(t, filepath.Join(testdir, "image_size_enum_encoding.go"), filepath.Join("internal", "testdata", "image_size_enum_encoding.go"))
			assertEqFile(t, filepath.Join(testdir, "image_size_enum_encoding_test.go"), filepath.Join("internal", "testdata", "image_size_enum_encoding_test.go"))
		})
	})

	t.Run("when bad go file, then error", func(t *testing.T) {
		exec.Command("cp", filepath.Join("internal", "README.md"), filepath.Join(testdir, "README.md")).Run()

		cmd := exec.Command(testbin, "--type", "Color")
		cmd.Env = append(cmd.Environ(), "GOFILE=README.md", "GOLINE=5", "GOPACKAGE=image", "GOTESTDIR="+testdir)
		if err := cmd.Run(); err == nil {
			t.Fatal("must be error")
		}
	})

	t.Run("when enum values not immediately after go:generate line, then error", func(t *testing.T) {
		cmd := exec.Command(testbin, "--type", "Color")
		cmd.Env = append(cmd.Environ(), "GOFILE="+filepath.Join(testdir, "image.go"), "GOLINE=1", "GOPACKAGE=image", "GOTESTDIR="+testdir)
		if err := cmd.Run(); err == nil {
			t.Fatal("must be error")
		}
	})

	t.Run("when invalid package name, then error", func(t *testing.T) {
		cmd := exec.Command(testbin, "--type", "Color")
		cmd.Env = append(cmd.Environ(), "GOFILE="+filepath.Join(testdir, "image.go"), "GOLINE=5", "GOPACKAGE=\"", "GOTESTDIR="+testdir)
		if err := cmd.Run(); err == nil {
			t.Fatal("must be error")
		}
	})

	t.Run("when not found file, then error", func(t *testing.T) {
		cmd := exec.Command(testbin, "--type", "Color")
		cmd.Env = append(cmd.Environ(), "GOFILE=asdf.asdf", "GOPACKAGE=image", "GOLINE=5", "GOTESTDIR="+testdir)
		if err := cmd.Run(); err == nil {
			t.Fatal("must be error")
		}
	})

	t.Run("when wrong params, then error", func(t *testing.T) {
		cmd := exec.Command(testbin)
		cmd.Env = append(cmd.Environ(), "GOFILE="+filepath.Join(testdir, "image.go"), "GOLINE=5", "GOPACKAGE=color", "GOTESTDIR="+testdir)
		if err := cmd.Run(); err == nil {
			t.Fatal("must be error")
		}
	})
}

func assertEqFile(t *testing.T, a, b string) {
	t.Helper()
	fa, err := os.ReadFile(a)
	if err != nil {
		t.Error(err)
	}
	fb, err := os.ReadFile(b)
	if err != nil {
		t.Error(err)
	}
	if len(fa) == 0 || len(fb) == 0 {
		t.Error("empty file: " + a + " or " + b)
	}
	if string(fa) != string(fb) {
		t.Error("files different: " + a + " != " + b + ": " + string(fa) + " <> " + string(fb))
	}
}
