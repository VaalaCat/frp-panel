package pty

type PTYInterface interface {
	Write(p []byte) (n int, err error)
	Read(p []byte) (n int, err error)
	Getsize() (uint16, uint16, error)
	Setsize(cols, rows uint32) error
	Close() error
}
