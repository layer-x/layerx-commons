package lxfileutils
import (
	"os"
	"fmt"
"io"
	"github.com/layer-x/layerx-commons/lxerrors"
"os/exec"
	"io/ioutil"
	"path/filepath"
)


func WriteFile(path string, data []byte) error {
	err := ioutil.WriteFile(path, data, 0777)
	if err != nil {
		err := os.MkdirAll(filepath.Dir(path), 0777)
		if err != nil {
			return err
		}
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.Write(data)
		if err != nil {
			return err
		}
	}
	return nil
}

func ReadFile(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return data, err
	}
	return data, nil
}


func Untar(src, dest string) error {
	tarPath, err := exec.LookPath("tar")

	if err != nil {
		return lxerrors.New("tar not found in path", nil)
	}

	err = os.MkdirAll(dest, 0755)
	if err != nil {
		return err
	}

	command := exec.Command(tarPath, "pzxf", src, "-C", dest)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}

func UntarNogzip(src, dest string) error {
	tarPath, err := exec.LookPath("tar")

	if err != nil {
		return lxerrors.New("tar not found in path", nil)
	}

	err = os.MkdirAll(dest, 0755)
	if err != nil {
		return err
	}

	command := exec.Command(tarPath, "pxf", src, "-C", dest)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}


//from http://stackoverflow.com/questions/21060945/simple-way-to-copy-a-file-in-golang

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

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
