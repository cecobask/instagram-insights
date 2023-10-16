package followdata

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

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
	Following   *userList
	Followers   *userList
	Unfollowers *userList
}

type userData struct {
	UserData []user `json:"string_list_data"`
}

func newFollowData() *followData {
	return &followData{
		Following: &userList{
			users:         make(map[string]user),
			showTimestamp: true,
		},
		Followers: &userList{
			users:         make(map[string]user),
			showTimestamp: true,
		},
		Unfollowers: &userList{
			users:         make(map[string]user),
			showTimestamp: false,
		},
	}
}

func (fd *followData) hydrateFollowers(data []byte) error {
	var jsonData []userData
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return err
	}
	for i := range jsonData {
		ud := jsonData[i].UserData[0]
		fd.Followers.users[ud.Username] = ud
	}
	return nil
}

func (fd *followData) hydrateFollowing(data []byte) error {
	jsonData := make(map[string][]userData)
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return err
	}
	for i := range jsonData["relationships_following"] {
		ud := jsonData["relationships_following"][i].UserData[0]
		fd.Following.users[ud.Username] = ud
	}
	return nil
}

func (fd *followData) hydrateUnfollowers() {
	for username, ud := range fd.Following.users {
		if _, found := fd.Followers.users[username]; !found {
			fd.Unfollowers.users[username] = ud
		}
	}
}

type timestamp struct {
	time.Time
}

func (t *timestamp) UnmarshalJSON(b []byte) error {
	var unixTimestamp int64
	if err := json.Unmarshal(b, &unixTimestamp); err != nil {
		return err
	}
	t.Time = time.Unix(unixTimestamp, 0)
	return nil
}

func (t *timestamp) String() string {
	return t.Format(time.DateOnly)
}

type user struct {
	ProfileUrl string     `json:"href"`
	Username   string     `json:"value"`
	Timestamp  *timestamp `json:"timestamp"`
}

type userList struct {
	users         map[string]user
	showTimestamp bool
}

func (u *userList) output(format string) error {
	switch format {
	case instagram.OutputTable:
		return u.outputTable()
	case instagram.OutputNone:
		return nil
	default:
		return fmt.Errorf("invalid output format: %s", format)
	}
}

func (u *userList) outputTable() error {
	var rows []table.Row
	for _, user := range u.users {
		row := table.Row{
			user.Username,
			user.ProfileUrl,
		}
		if u.showTimestamp {
			row = append(row, user.Timestamp)
		}
		rows = append(rows, row)
	}
	headers := table.Row{
		instagram.TableHeaderUsername,
		instagram.TableHeaderProfileUrl,
	}
	if u.showTimestamp {
		headers = append(headers, instagram.TableHeaderTimestamp)
	}
	usersTable := table.NewWriter()
	usersTable.SetAutoIndex(true)
	usersTable.SetOutputMirror(os.Stdout)
	usersTable.SetStyle(table.StyleBold)
	usersTable.AppendHeader(headers)
	usersTable.AppendRows(rows)
	usersTable.Render()
	return nil
}
