package store

import (
	"io"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func isDirExist(dirPath string) bool {
	_, err := os.Stat(dirPath)
	return !os.IsNotExist(err)
}

// tearDown removes temporary files and folders.
func tearDown() {
	_ = os.RemoveAll(downloadDir)
}

func TestLocalFs_Save(t *testing.T) {
	lfs := LocalFs{}

	fName, content := "file1.txt", []byte("content 1")
	err := lfs.Save(fName, content)
	assert.NoError(t, err)
	assert.True(t, isDirExist(downloadDir))
	assert.True(t, isDirExist(path.Join(downloadDir, fName)))
	// #nosec G304 - This is a test file with controlled input
	f, _ := os.OpenFile(path.Join(downloadDir, fName), os.O_RDONLY, 0o0600)
	bs, _ := io.ReadAll(f)
	assert.Equal(t, string(content), string(bs))

	fName, content = "file2.txt", []byte("content 2")
	err = lfs.Save(fName, content)
	assert.NoError(t, err)
	assert.True(t, isDirExist(downloadDir))
	assert.True(t, isDirExist(path.Join(downloadDir, fName)))
	// #nosec G304 - This is a test file with controlled input
	f, _ = os.OpenFile(path.Join(downloadDir, fName), os.O_RDONLY, 0o0600)
	bs, _ = io.ReadAll(f)
	assert.Equal(t, string(content), string(bs))

	tearDown()
}

func Test_createDownloadDir(t *testing.T) {
	assert.False(t, isDirExist(downloadDir))
	err := createDownloadDir()
	assert.NoError(t, err)
	assert.True(t, isDirExist(downloadDir))

	// already exist
	err = createDownloadDir()
	assert.NoError(t, err)

	tearDown()
}
