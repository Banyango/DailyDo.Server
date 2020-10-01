package os

import (
	"os"
	"path"
)

type OSInterface interface {
	GetEnv(variable string) string
	GetFile(dir string, filename string) (*os.File, error)
}

type OSService struct {
}

func NewOSService() *OSService {
	m := new(OSService)
	return m
}

func (e *OSService) GetEnv(variable string) string {
	return os.Getenv(variable)
}

func (e *OSService) GetFile(dir string, filename string) (*os.File, error) {
	join := path.Join(dir, filename)
	return os.Open(join)
}
