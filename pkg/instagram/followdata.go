package instagram

import (
	"encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"path/filepath"
)

type User struct {
	ProfileURL string `json:"href"`
	Username   string `json:"value"`
	Timestamp  int64  `json:"timestamp"`
}

type FollowData struct {
	Following map[string]User
	Followers map[string]User
}

func NewFollowData() *FollowData {
	return &FollowData{
		Following: make(map[string]User),
		Followers: make(map[string]User),
	}
}

type userData struct {
	UserData []User `json:"string_list_data"`
}

func (fd *FollowData) ExtractFollowing() error {
	data, err := os.ReadFile(pathFollowing)
	if err != nil {
		return err
	}
	jsonStruct := make(map[string][]userData)
	if err = json.Unmarshal(data, &jsonStruct); err != nil {
		return err
	}
	for _, following := range jsonStruct["relationships_following"] {
		userData := following.UserData[0]
		fd.Following[userData.Username] = userData
	}
	return nil
}

func (fd *FollowData) ExtractFollowers() error {
	files, err := filepath.Glob(pathFollowers)
	if err != nil {
		return err
	}
	for i := range files {
		data, err := os.ReadFile(files[i])
		if err != nil {
			return err
		}
		var jsonArray []userData
		err = json.Unmarshal(data, &jsonArray)
		if err != nil {
			return err
		}
		for _, follower := range jsonArray {
			userData := follower.UserData[0]
			fd.Followers[userData.Username] = userData
		}
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
	var tableRows []table.Row
	for username, userData := range fd.Following {
		if _, found := fd.Followers[username]; !found {
			tableRows = append(tableRows, table.Row{
				userData.Username,
				userData.ProfileURL,
			})
		}
	}
	if len(tableRows) == 0 {
		fmt.Println("No unfollowers found!")
		return
	}
	unfollowersTable := table.NewWriter()
	unfollowersTable.SetAutoIndex(true)
	unfollowersTable.SetOutputMirror(os.Stdout)
	unfollowersTable.SetStyle(table.StyleBold)
	unfollowersTable.AppendHeader(table.Row{
		tableHeaderUsername,
		tableHeaderProfileUrl,
	})
	unfollowersTable.AppendRows(tableRows)
	unfollowersTable.Render()
}
