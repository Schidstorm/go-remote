package lib

import "io"

// StdInOutPipe a readwritercloser for sstdin and stdout and exits on close
type StdInOutPipe struct {
	Reader io.Reader
	Writer io.Writer
	Closer io.Closer
}

// Close does nothing
func (stdInOutPipe StdInOutPipe) Close() error {
	return stdInOutPipe.Closer.Close()
}

func (stdInOutPipe StdInOutPipe) Read(p []byte) (n int, err error) {
	return stdInOutPipe.Reader.Read(p)
}

func (stdInOutPipe StdInOutPipe) Write(p []byte) (n int, err error) {
	return stdInOutPipe.Writer.Write(p)
}
