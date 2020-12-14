package storage

type Storage interface {
	Save(descriptor string, content []byte) error
	Read(descriptor string) ([]byte, error)
	Remove(descriptor string)
}
