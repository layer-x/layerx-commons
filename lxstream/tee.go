package lxstream
import (
	"os"
	"io"
"bytes"
	"bufio"
"github.com/layer-x/layerx-commons/lxerrors"
)

func Tee(file *os.File, buf *bytes.Buffer) error {
	r, w, err := os.Pipe()
	if err != nil {
		return lxerrors.New("creating pipe", err)
	}
	stdout := file
	file = w
	multi := io.MultiWriter(stdout, bufio.NewWriter(buf))
	reader := bufio.NewReader(r)
	go func() {
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				return
			}
			_, err = multi.Write(append(line, byte('\n')))
			if err != nil {
				return
			}
		}
	}()
	return nil
}
