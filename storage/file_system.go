package storage

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type fileSystemPasswordStorage struct {
	basePath string
}

func (fileSystemStorage *fileSystemPasswordStorage) Save(descriptor string, content []byte) error {
	if len(descriptor) == 0 {
		return fmt.Errorf("Descriptor length cannot be empty")
	}
	if len(content) == 0 {
		return fmt.Errorf("File content cannot be empty")
	}
	return ioutil.WriteFile(fileSystemStorage.buildFullDescriptor(descriptor), content, 0644)
}

func (fileSystemStorage *fileSystemPasswordStorage) Read(descriptor string) ([]byte, error) {
	if len(descriptor) == 0 {
		return nil, fmt.Errorf("Descriptor length cannot be empty")
	}
	contentBytes, err := ioutil.ReadFile(fileSystemStorage.buildFullDescriptor(descriptor))
	if err != nil {
		return nil, err
	}
	return contentBytes, nil
}

func (fileSystemStorage *fileSystemPasswordStorage) Remove(descriptor string) {
	os.Remove(fileSystemStorage.buildFullDescriptor(descriptor))
}

func (fileSystemStorage *fileSystemPasswordStorage) buildFullDescriptor(descriptor string) string {
	return fileSystemStorage.basePath + string(filepath.Separator) + descriptor
}

func NewFileSystemPasswordStorage(basePath string) Storage {
	if len(basePath) == 0 {
		panic("Base path should be provided!")
	}
	return &fileSystemPasswordStorage{basePath}
}
