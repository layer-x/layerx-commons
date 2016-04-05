package lxfileutils

import (
	"mime/multipart"
	"os"
	"path/filepath"
	"github.com/layer-x/layerx-commons/lxerrors"
	"io"
)

func UntarFileToDirectory(targetDirectory string, sourceTar *multipart.File, header *multipart.FileHeader) (int, error) {
	savedTar, err := os.OpenFile(targetDirectory +filepath.Base(header.Filename), os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return 0, lxerrors.New("creating empty file for copying to", err)
	}
	defer savedTar.Close()
	bytesWritten, err := io.Copy(savedTar, sourceTar)
	if err != nil {
		return 0, lxerrors.New("copying uploaded file to disk", err)
	}
	err = Untar(savedTar.Name(), targetDirectory)
	if err != nil {
		err = UntarNogzip(savedTar.Name(), targetDirectory)
		if err != nil {
			return 0, lxerrors.New("untarring saved tar", err)
		}
	}
	return bytesWritten, nil
}