package main_test

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	var covdirs []string
	testdir := t.TempDir()
	testbin := path.Join(testdir, "go-enum-encoding-test")

	if b, err := exec.Command("go", "build", "-cover", "-o", testbin, "main.go").CombinedOutput(); err != nil {
		t.Fatal(err, string(b))
	}

	defer func() {
		if b, err := exec.Command("go", "tool", "covdata", "textfmt", "-i="+strings.Join(covdirs, ","), "-o", os.Getenv("GOCOVERPROFILE")).CombinedOutput(); err != nil {
			t.Error(err, string(b))
		}
	}()

	t.Run("ok", func(t *testing.T) {
		if b, err := exec.Command("cp", filepath.Join("internal", "testdata", "image.go"), filepath.Join(testdir, "image.go")).CombinedOutput(); err != nil {
			t.Fatal(err, string(b))
		}

		t.Run("struct", func(t *testing.T) {
			covdir := "cov_struct"
			exec.Command("mkdir", covdir).Run()
			covdirs = append(covdirs, covdir)

			cmd := exec.Command(testbin, "--type", "Color")
			cmd.Env = append(cmd.Environ(), "GOFILE="+filepath.Join(testdir, "image.go"), "GOLINE=4", "GOPACKAGE=image", "GOTESTDIR="+testdir, "GOCOVERDIR="+covdir)
			if err := cmd.Run(); err != nil {
				t.Error(err)
			}

			assertEqFile(t, filepath.Join(testdir, "color_enum_encoding.go"), filepath.Join("internal", "testdata", "color_enum_encoding.go"))
			assertEqFile(t, filepath.Join(testdir, "color_enum_encoding_test.go"), filepath.Join("internal", "testdata", "color_enum_encoding_test.go"))
		})

		t.Run("iota, string", func(t *testing.T) {
			covdir := "cov_iota"
			exec.Command("mkdir", covdir).Run()
			covdirs = append(covdirs, covdir)

			cmd := exec.Command(testbin, "--type", "ImageSize", "--string")
			cmd.Env = append(cmd.Environ(), "GOFILE="+filepath.Join(testdir, "image.go"), "GOLINE=18", "GOPACKAGE=image", "GOTESTDIR="+testdir, "GOCOVERDIR="+covdir)
			if err := cmd.Run(); err != nil {
				t.Error(err)
			}

			assertEqFile(t, filepath.Join(testdir, "image_size_enum_encoding.go"), filepath.Join("internal", "testdata", "image_size_enum_encoding.go"))
			assertEqFile(t, filepath.Join(testdir, "image_size_enum_encoding_test.go"), filepath.Join("internal", "testdata", "image_size_enum_encoding_test.go"))
		})

		t.Run("run tests within generated code", func(t *testing.T) {
			from, _ := os.Getwd()
			os.Chdir(testdir)
			t.Cleanup(func() { os.Chdir(from) })

			os.WriteFile(filepath.Join(testdir, "go.mod"), []byte("module test\ngo 1.24"), 0644)

			if b, err := exec.Command("go", "test", ".").CombinedOutput(); err != nil {
				t.Error(err, string(b))
			}
		})
	})

	t.Run("when bad go file, then error", func(t *testing.T) {
		covdir := "cov_err_bad_file"
		exec.Command("mkdir", covdir).Run()
		covdirs = append(covdirs, covdir)

		exec.Command("cp", filepath.Join("internal", "README.md"), filepath.Join(testdir, "README.md")).Run()

		cmd := exec.Command(testbin, "--type", "Color")
		cmd.Env = append(cmd.Environ(), "GOFILE=README.md", "GOLINE=5", "GOPACKAGE=image", "GOTESTDIR="+testdir, "GOCOVERDIR="+covdir)
		if err := cmd.Run(); err == nil {
			t.Fatal("must be error")
		}
	})

	t.Run("when enum values not immediately after go:generate line, then error", func(t *testing.T) {
		covdir := "cov_err_enum"
		exec.Command("mkdir", covdir).Run()
		covdirs = append(covdirs, covdir)

		cmd := exec.Command(testbin, "--type", "Color")
		cmd.Env = append(cmd.Environ(), "GOFILE="+filepath.Join(testdir, "image.go"), "GOLINE=1", "GOPACKAGE=image", "GOTESTDIR="+testdir, "GOCOVERDIR="+covdir)
		if err := cmd.Run(); err == nil {
			t.Fatal("must be error")
		}
	})

	t.Run("when invalid package name, then error", func(t *testing.T) {
		covdir := "cov_err_invalid_pkg"
		exec.Command("mkdir", covdir).Run()
		covdirs = append(covdirs, covdir)

		cmd := exec.Command(testbin, "--type", "Color")
		cmd.Env = append(cmd.Environ(), "GOFILE="+filepath.Join(testdir, "image.go"), "GOLINE=5", "GOPACKAGE=\"", "GOTESTDIR="+testdir, "GOCOVERDIR="+covdir)
		if err := cmd.Run(); err == nil {
			t.Fatal("must be error")
		}
	})

	t.Run("when not found file, then error", func(t *testing.T) {
		covdir := "cov_err_file_not_found"
		exec.Command("mkdir", covdir).Run()
		covdirs = append(covdirs, covdir)

		cmd := exec.Command(testbin, "--type", "Color")
		cmd.Env = append(cmd.Environ(), "GOFILE=asdf.asdf", "GOPACKAGE=image", "GOLINE=5", "GOTESTDIR="+testdir, "GOCOVERDIR="+covdir)
		if err := cmd.Run(); err == nil {
			t.Fatal("must be error")
		}
	})

	t.Run("when wrong params, then error", func(t *testing.T) {
		covdir := "cov_err_wrong_params"
		exec.Command("mkdir", covdir).Run()
		covdirs = append(covdirs, covdir)

		cmd := exec.Command(testbin)
		cmd.Env = append(cmd.Environ(), "GOFILE="+filepath.Join(testdir, "image.go"), "GOLINE=5", "GOPACKAGE=color", "GOTESTDIR="+testdir, "GOCOVERDIR="+covdir)
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
