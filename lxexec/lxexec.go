package lxexec
import (
	"bytes"
	"github.com/layer-x/layerx-commons/lxerrors"
	"os/exec"
)

func RunCommand(args ...string) (string, error) {
	if len(args) < 1 {
		return "", lxerrors.New("must provide path to command", nil)
	}
	cmd := exec.Command(args...)

	stdErr, err := cmd.StderrPipe()
	if err != nil {
		return "", lxerrors.New("piping stderr of command", err)
	}

	errBuf := new(bytes.Buffer)
	_, err = errBuf.ReadFrom(stdErr)
	if err != nil {
		return "", lxerrors.New("reading stdout from pipe", err)
	}

	out, err := cmd.Output()
	if err != nil {
		return "", lxerrors.New("command exited with error: "+errBuf.String(), err)
	}
	return out, nil
}
