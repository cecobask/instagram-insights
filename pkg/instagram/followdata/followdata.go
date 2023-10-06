package followdata

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cecobask/instagram-insights/pkg/filesystem"
	"github.com/cecobask/instagram-insights/pkg/instagram"
	"github.com/jedib0t/go-pretty/v6/table"
)

type Interface interface {
	Followers(opts instagram.Options) error
	Following(opts instagram.Options) error
	Unfollowers(opts instagram.Options) error
}

type handler struct {
	fileSystem filesystem.Fs
	followData *followData
}

func NewHandler() Interface {
	return &handler{
		fileSystem: filesystem.NewFs(),
		followData: newFollowData(),
	}
}

func (h *handler) Followers(opts instagram.Options) error {
	files, err := h.fileSystem.FindFiles(instagram.PathFollowers)
	if err != nil {
		return err
	}
	for i := range files {
		data, err := h.fileSystem.ReadFile(files[i])
		if err != nil {
			return err
		}
		if err = h.followData.hydrateFollowers(data); err != nil {
			return err
		}
	}
	return h.followData.Followers.output(opts.Output)
}

func (h *handler) Following(opts instagram.Options) error {
	data, err := h.fileSystem.ReadFile(instagram.PathFollowing)
	if err != nil {
		return err
	}
	if err = h.followData.hydrateFollowing(data); err != nil {
		return err
	}
	return h.followData.Following.output(opts.Output)
}

func (h *handler) Unfollowers(opts instagram.Options) error {
	childOptions := instagram.NewOptions(instagram.OutputNone)
	if err := h.Followers(childOptions); err != nil {
		return err
	}
	if err := h.Following(childOptions); err != nil {
		return err
	}
	h.followData.hydrateUnfollowers()
	return h.followData.Unfollowers.output(opts.Output)
}

type followData struct {
	Following   users
	Followers   users
	Unfollowers users
}

type userData struct {
	UserData []user `json:"string_list_data"`
}

func newFollowData() *followData {
	return &followData{
		Following:   make(users),
		Followers:   make(users),
		Unfollowers: make(users),
	}
}

func (fd *followData) hydrateFollowers(data []byte) error {
	var jsonData []userData
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return err
	}
	for _, follower := range jsonData {
		userData := follower.UserData[0]
		fd.Followers[userData.Username] = userData
	}
	return nil
}

func (fd *followData) hydrateFollowing(data []byte) error {
	jsonData := make(map[string][]userData)
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return err
	}
	for _, following := range jsonData["relationships_following"] {
		userData := following.UserData[0]
		fd.Following[userData.Username] = userData
	}
	return nil
}

func (fd *followData) hydrateUnfollowers() {
	for username, user := range fd.Following {
		if _, found := fd.Followers[username]; !found {
			fd.Unfollowers[username] = user
		}
	}
}

type user struct {
	ProfileUrl string `json:"href"`
	Username   string `json:"value"`
	Timestamp  int    `json:"timestamp"`
}

type users map[string]user

func (u users) output(format string) error {
	switch format {
	case instagram.OutputTable:
		return u.outputTable()
	case instagram.OutputNone:
		return nil
	default:
		return fmt.Errorf("invalid output format: %s", format)
	}
}

func (u users) outputTable() error {
	var rows []table.Row
	for _, user := range u {
		rows = append(rows, table.Row{
			user.Username,
			user.ProfileUrl,
		})
	}
	usersTable := table.NewWriter()
	usersTable.SetAutoIndex(true)
	usersTable.SetOutputMirror(os.Stdout)
	usersTable.SetStyle(table.StyleBold)
	usersTable.AppendHeader(table.Row{
		instagram.TableHeaderUsername,
		instagram.TableHeaderProfileUrl,
	})
	usersTable.AppendRows(rows)
	usersTable.Render()
	return nil
}
