package store

import (
	"os"
	"path/filepath"
)

const downloadDir = "backups"

// LocalFs is Saver to local disk
type LocalFs struct{}

// Save performs writing file to disk
func (l LocalFs) Save(filename string, bs []byte) error {
	err := createDownloadDir()
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(".", downloadDir, filename), bs, 0o0600)
	return err
}

func createDownloadDir() error {
	return os.MkdirAll(filepath.Join(".", downloadDir), os.ModePerm)
}
