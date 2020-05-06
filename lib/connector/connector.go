package connector

import (
	"io"
)

type Connector interface {
	Connect() (io.ReadWriteCloser, error)
}
