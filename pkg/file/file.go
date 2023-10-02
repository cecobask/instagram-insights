package file

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func DownloadFile(url, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failure downloading instagram information from google drive... http status: %s", response.Status)
	}
	if _, err = io.Copy(file, response.Body); err != nil {
		return err
	}
	return nil
}

func UnzipFile(source, destination string) error {
	r, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer r.Close()
	if err = os.MkdirAll(destination, 0755); err != nil {
		return err
	}
	for _, f := range r.File {
		if err = extractAndWriteFile(f, destination); err != nil {
			return err
		}
	}
	return nil
}

func CleanupFilePaths(filePaths []string) error {
	for _, path := range filePaths {
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}
	return nil
}

func extractAndWriteFile(f *zip.File, destination string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer func() {
		if err = rc.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	path := filepath.Join(destination, f.Name)
	if !strings.HasPrefix(path, filepath.Clean(destination)+string(os.PathSeparator)) {
		return fmt.Errorf("illegal file path: %s", path)
	}
	if f.FileInfo().IsDir() {
		if err = os.MkdirAll(path, f.Mode()); err != nil {
			return err
		}
	} else {
		if err = os.MkdirAll(filepath.Dir(path), f.Mode()); err != nil {
			return err
		}
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer func() {
			if err = file.Close(); err != nil {
				log.Fatal(err)
			}
		}()
		if _, err = io.Copy(file, rc); err != nil {
			return err
		}
	}
	return nil
}
