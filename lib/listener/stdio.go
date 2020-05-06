package listener

import (
	"io"
	"os"

	"github.com/schidstorm/go-remote/lib"
)

type Io struct {
}

func (io *Io) Listen() (io.ReadWriteCloser, error) {
	return lib.StdInOutPipe{
		Reader: os.Stdin,
		Writer: os.Stdout,
		Closer: os.Stdin,
	}, nil
}
