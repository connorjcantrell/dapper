package dapperfs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Fs struct {
	BoilerplateDir string
	WorkingDir     string
}

/*
New() creates a new FS struct, first by retrieving the environment variable
DAPPER_BOILERPLATE_DIR, which contains a path to local directory containing
boilerplate code. Secondly, it stores a path to the current working directory.
*/
func New() (Fs, error) {
	boilerplateDir := os.Getenv("DAPPER_BOILERPLATE_DIR")
	if boilerplateDir == "" {
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

/*
Mkdir creates a new directory with a provided name. Returns an error.
*/
func (fs *Fs) Mkdir(name string) error {
	err := os.Mkdir(name, 0777)
	if err != nil {
		return err
	}
	return nil
}

/*
WriteConfig marshals any struct to JSON and writes to FS.workingDir/relPath/filename
*/
func (fs *Fs) WriteStructToJSON(s interface{}, relPath, filename string) error {
	b, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		return err
	}
	// p := filepath.Join(fs.WorkingDir, ".dapper", "config.json")
	p := filepath.Join(fs.WorkingDir, relPath, filename)
	if err := os.WriteFile(p, b, 0644); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

/*
CopyFromBoilerplate attempts to copy all files from "~/boilerplate/{dir}" to
"path-to-project/src"
*/
func (fs *Fs) CopyFromBoilerplateDir(relPath string) error {
	boilerplateDir := filepath.Join(fs.BoilerplateDir, relPath)
	dst := filepath.Join(fs.WorkingDir, "src")
	files, err := os.ReadDir(boilerplateDir)

	if err != nil {
		return err
	}

	for _, f := range files {
		srcPath := filepath.Join(boilerplateDir, f.Name())
		dstPath := filepath.Join(dst, f.Name())
		err := CopyFile(srcPath, dstPath)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func CopyFile(src, dst string) error {
	sfi, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return err
		}
	}
	if err = os.Link(src, dst); err == nil {
		return err
	}
	err = copyFileContents(src, dst)
	return err
}

// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
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
