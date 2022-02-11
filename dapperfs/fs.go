package dapperfs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Fs struct {
	BoilerplateDir string
	WorkingDir     string
}

func NewFs() (Fs, error) {
	boilerplateDir := os.Getenv("DAPPER_BOILERPLATE_DIR")
	exists, err := exists(boilerplateDir)
	if err != nil {
		return Fs{}, err
	}
	if !exists {
		return Fs{}, errors.New("source directory does not exist")
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return Fs{}, nil
	}

	fs := Fs{
		BoilerplateDir: boilerplateDir,
		WorkingDir:     workingDir,
	}
	return fs, nil
}

func (fs *Fs) Mkdir(s string) error {
	err := os.Mkdir(s, 0777)
	if err != nil {
		return err
	}
	return nil
}

func (fs *Fs) WriteConfig(s interface{}) error {
	configJson, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		return err
	}
	p := filepath.Join(fs.WorkingDir, ".dapper", "config.json")
	if err := ioutil.WriteFile(p, configJson, 0644); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

/*
CopyFromBoilerplate attempts to copy all files from "~/boilerplate/{dir}" to
"path-to-project/src"
*/
func (fs *Fs) CopyFromBoilerplate(dir string) error {
	boilerplate := filepath.Join(fs.BoilerplateDir, dir)
	if _, err := exists(boilerplate); err != nil {
		return err
	}
	dst := filepath.Join(fs.WorkingDir, "src")
	if _, err := exists(dst); err != nil {
		return err
	}

	files, err := ioutil.ReadDir(boilerplate)
	if err != nil {
		return err
	}

	for _, f := range files {
		srcPath := filepath.Join(boilerplate, f.Name())
		dstPath := filepath.Join(dst, f.Name())
		err := CopyFile(srcPath, dstPath)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return true, nil
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	err = out.Sync()
	return err
}
