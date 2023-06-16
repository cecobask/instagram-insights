package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type UserData struct {
	ProfileURL string `json:"href"`
	Username   string `json:"value"`
	Timestamp  int64  `json:"timestamp"`
}

type FollowData struct {
	Following []UserData
	Followers []UserData
}

type UserDataJson struct {
	UserData []UserData `json:"string_list_data"`
}

func (fd *FollowData) ExtractFollowing() error {
	data, err := os.ReadFile(pathFollowing)
	if err != nil {
		return err
	}
	jsonStruct := make(map[string][]UserDataJson)
	if err = json.Unmarshal(data, &jsonStruct); err != nil {
		return err
	}
	for _, following := range jsonStruct["relationships_following"] {
		fd.Following = append(fd.Following, following.UserData[0])
	}
	return nil
}

func (fd *FollowData) ExtractFollowers() error {
	data, err := os.ReadFile(pathFollowers)
	if err != nil {
		return err
	}
	var jsonArray []UserDataJson
	err = json.Unmarshal(data, &jsonArray)
	if err != nil {
		return err
	}
	for _, follower := range jsonArray {
		fd.Followers = append(fd.Followers, follower.UserData[0])
	}
	return nil
}

func (fd *FollowData) ExtractAllData() error {
	err := fd.ExtractFollowing()
	if err != nil {
		return err
	}
	err = fd.ExtractFollowers()
	if err != nil {
		return err
	}
	return nil
}

func (fd *FollowData) FindUnfollowers() {
	fmt.Println("=================================================================================================")
	for _, following := range fd.Following {
		var found bool
		for _, follower := range fd.Followers {
			if following.Username == follower.Username {
				found = true
			}
		}
		if !found {
			fmt.Printf("User %s does not follow you back... Profile URL: %s\n", following.Username, following.ProfileURL)
		}
	}
	fmt.Println("=================================================================================================")
}
