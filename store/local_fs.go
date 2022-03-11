package store

import "os"

// LocalFs is Saver to local disk
type LocalFs struct{}

// Save performs writing file to disk
func (l LocalFs) Save(filename string, bs []byte) error {
	err := os.WriteFile(filename, bs, 0o0600)
	return err
}
