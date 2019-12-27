package connector

type Connector interface {
	Open() error
	Close()
	ReadUntilCR() ([]byte, error)
	Read(terminator byte) ([]byte, error)
	Write(bytes []byte) error
}
