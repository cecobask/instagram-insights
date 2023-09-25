package main

import (
	"archive/zip"
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	envVarKeyFileID = "FILE_ID"
	fileUrlFormat   = "https://drive.google.com/file/d/1B3cPMMROyUYpn_ZGAku8aVye2MGCkPfs"
			 //https://drive.google.com/drive/folders/1B3cPMMROyUYpn_ZGAku8aVye2MGCkPfs?usp=sharing
	pathData        = "instagram_data"
	pathDataArchive = pathData + ".zip"
	pathFollowers   = pathData + "/followers_and_following/followers_*.json"
	pathFollowing   = pathData + "/followers_and_following/following.json"
)

func main() {
	fileID, ok := os.LookupEnv(envVarKeyFileID)
	if !ok || fileID == "" {
		log.Fatalf("environment variable %s not set", envVarKeyFileID)
	}
	defer cleanup()
	if err := downloadFile(fmt.Sprintf(fileUrlFormat, fileID), pathDataArchive); err != nil {
		log.Fatal(err)
	}
	if err := unzip(pathDataArchive, pathData); err != nil {
		log.Fatal(err)
	}
	followData := NewFollowData()
	if err := followData.ExtractAllData(); err != nil {
		log.Fatal(err)
	}
	followData.FindUnfollowers()
}

func cleanup() error {
	if err := os.RemoveAll(pathDataArchive); err != nil {
		return err
	}
	if err := os.RemoveAll(pathData); err != nil {
		return err
	}
	return nil
}

func downloadFile(url, fileName string) error {
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
		return fmt.Errorf("failure downloading instagram data from google drive... http status: %s", response.Status)
	}
	if _, err = io.Copy(file, response.Body); err != nil {
		return err
	}
	return nil
}

func unzip(source, destination string) error {
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
