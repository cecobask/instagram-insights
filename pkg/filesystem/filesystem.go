package filesystem

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

type Fs interface {
	CopyToFile(destination io.Writer, source io.Reader) (int64, error)
	CreateDirectory(path string, perm os.FileMode) error
	CreateFile(name string) (io.WriteCloser, error)
	FindFiles(pattern string) ([]string, error)
	OpenFile(name string, flag int, perm os.FileMode) (*os.File, error)
	OpenZip(name string) (*zip.ReadCloser, error)
	ReadFile(name string) ([]byte, error)
	ReadZipFile(file *zip.File) (io.ReadCloser, error)
	RemoveDirectory(path string) error
	Unzip(source, destination string) error
	UnzipFile(zipFile *zip.File, destination string) error
}

type fileSystem struct{}

func NewFs() Fs {
	return new(fileSystem)
}

func (fs *fileSystem) CopyToFile(destination io.Writer, source io.Reader) (int64, error) {
	return io.Copy(destination, source)
}

func (fs *fileSystem) CreateDirectory(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (fs *fileSystem) CreateFile(name string) (io.WriteCloser, error) {
	return os.Create(name)
}

func (fs *fileSystem) FindFiles(pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}

func (fs *fileSystem) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

func (fs *fileSystem) OpenZip(name string) (*zip.ReadCloser, error) {
	return zip.OpenReader(name)
}

func (fs *fileSystem) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (fs *fileSystem) ReadZipFile(file *zip.File) (io.ReadCloser, error) {
	return file.Open()
}

func (fs *fileSystem) RemoveDirectory(path string) error {
	return os.RemoveAll(path)
}

func (fs *fileSystem) Unzip(source, destination string) error {
	archive, err := fs.OpenZip(source)
	if err != nil {
		return err
	}
	defer archive.Close()
	if err = fs.CreateDirectory(destination, 0755); err != nil {
		return err
	}
	for _, file := range archive.File {
		if err = fs.UnzipFile(file, destination); err != nil {
			return err
		}
	}
	return nil
}

func (fs *fileSystem) UnzipFile(zipFile *zip.File, destination string) error {
	reader, err := fs.ReadZipFile(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	path := filepath.Join(destination, zipFile.Name)
	if zipFile.FileInfo().IsDir() {
		return fs.CreateDirectory(path, zipFile.Mode())
	}
	if err = fs.CreateDirectory(filepath.Dir(path), zipFile.Mode()); err != nil {
		return err
	}
	file, err := fs.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zipFile.Mode())
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err = fs.CopyToFile(file, reader); err != nil {
		return err
	}
	return nil
}
