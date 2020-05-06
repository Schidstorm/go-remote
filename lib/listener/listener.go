package listener

import "io"

type Listener interface {
	Listen() (io.ReadWriteCloser, error)
}
