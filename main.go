package main

import (
	"archive/zip"
	"encoding/json"
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
	instagramDataFolder  = "instagram_data"
	instagramDataArchive = instagramDataFolder + ".zip"
	fileUrlFormat        = "https://drive.google.com/u/0/uc?id=%s&export=download&confirm=t"
	envVarKeyFileID      = "FILE_ID"
)

type StringListData struct {
	Href      string
	Value     string
	Timestamp int64
}

type Followers struct {
	StringListData []StringListData `json:"string_list_data"`
}

type Following struct {
	RelationshipsFollowing []struct {
		StringListData []StringListData `json:"string_list_data"`
	} `json:"relationships_following"`
}

type FollowData struct {
	ProfileURL string
	Username   string
	Timestamp  int64
}

func main() {
	fileID, ok := os.LookupEnv(envVarKeyFileID)
	if !ok || fileID == "" {
		log.Fatalf("environment variable %s not set", envVarKeyFileID)
	}
	defer func() {
		err := cleanup()
		if err != nil {
			log.Fatal(err)
		}
	}()
	fileUrl := fmt.Sprintf(fileUrlFormat, fileID)
	err := downloadFile(fileUrl, instagramDataArchive)
	if err != nil {
		log.Fatal(err)
	}
	err = unzip(instagramDataArchive, instagramDataFolder)
	if err != nil {
		log.Fatal(err)
	}
	following, err := extractFollowData("instagram_data/followers_and_following/following.json")
	if err != nil {
		log.Fatal(err)
	}
	followers, err := extractFollowData("instagram_data/followers_and_following/followers_1.json")
	if err != nil {
		log.Fatal(err)
	}
	findUnfollowers(following, followers)
}

func cleanup() error {
	err := os.RemoveAll(instagramDataArchive)
	if err != nil {
		return err
	}
	err = os.RemoveAll(instagramDataFolder)
	if err != nil {
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
		return fmt.Errorf("bad status: %s", response.Status)
	}
	_, err = io.Copy(file, response.Body)
	if err != nil {
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
	err = os.MkdirAll(destination, 0755)
	if err != nil {
		return err
	}
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()
		path := filepath.Join(destination, f.Name)
		if !strings.HasPrefix(path, filepath.Clean(destination)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			err = os.MkdirAll(path, f.Mode())
			if err != nil {
				return err
			}
		} else {
			err = os.MkdirAll(filepath.Dir(path), f.Mode())
			if err != nil {
				return err
			}
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}
	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}
	return nil
}

func extractFollowData(filePath string) ([]FollowData, error) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	var followData []FollowData
	if strings.HasSuffix(filePath, "following.json") {
		var following Following
		err = json.Unmarshal(byteValue, &following)
		if err != nil {
			return nil, err
		}
		for i := range following.RelationshipsFollowing {
			entry := following.RelationshipsFollowing[i].StringListData[0]
			followData = append(followData, FollowData{
				ProfileURL: entry.Href,
				Username:   entry.Value,
				Timestamp:  entry.Timestamp,
			})
		}
	} else {
		var followers []Followers
		err = json.Unmarshal(byteValue, &followers)
		if err != nil {
			return nil, err
		}
		for i := range followers {
			entry := followers[i].StringListData[0]
			followData = append(followData, FollowData{
				ProfileURL: entry.Href,
				Username:   entry.Value,
				Timestamp:  entry.Timestamp,
			})
		}
	}
	return followData, nil
}

func findUnfollowers(following, followers []FollowData) {
	fmt.Println("=================================================================================================")
	for a := range following {
		var found bool
		for b := range followers {
			if following[a].Username == followers[b].Username {
				found = true
			}
		}
		if !found {
			fmt.Printf("User %s does not follow you back... Profile URL: %s\n", following[a].Username, following[a].ProfileURL)
		}
	}
	fmt.Println("=================================================================================================")
}
