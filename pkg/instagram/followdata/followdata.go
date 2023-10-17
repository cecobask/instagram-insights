package followdata

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/cecobask/instagram-insights/pkg/filesystem"
	"github.com/cecobask/instagram-insights/pkg/instagram"
	"github.com/jedib0t/go-pretty/v6/table"
	"gopkg.in/yaml.v3"
)

type Interface interface {
	Followers(opts instagram.Options) (*string, error)
	Following(opts instagram.Options) (*string, error)
	Unfollowers(opts instagram.Options) (*string, error)
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

func (h *handler) Followers(opts instagram.Options) (*string, error) {
	files, err := h.fileSystem.FindFiles(instagram.PathFollowers)
	if err != nil {
		return nil, err
	}
	for i := range files {
		data, err := h.fileSystem.ReadFile(files[i])
		if err != nil {
			return nil, err
		}
		if err = h.followData.hydrateFollowers(data); err != nil {
			return nil, err
		}
	}
	return h.followData.Followers.output(opts.Output)
}

func (h *handler) Following(opts instagram.Options) (*string, error) {
	data, err := h.fileSystem.ReadFile(instagram.PathFollowing)
	if err != nil {
		return nil, err
	}
	if err = h.followData.hydrateFollowing(data); err != nil {
		return nil, err
	}
	return h.followData.Following.output(opts.Output)
}

func (h *handler) Unfollowers(opts instagram.Options) (*string, error) {
	childOptions := instagram.NewOptions(instagram.OutputNone)
	if _, err := h.Followers(childOptions); err != nil {
		return nil, err
	}
	if _, err := h.Following(childOptions); err != nil {
		return nil, err
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
	UserData []userOriginal `json:"string_list_data"`
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
		fd.Followers.users[ud.Value] = user{
			ProfileUrl: ud.Href,
			Username:   ud.Value,
			Timestamp: &timestamp{
				Time: time.Unix(int64(ud.Timestamp), 0),
			},
		}
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
		fd.Following.users[ud.Value] = user{
			ProfileUrl: ud.Href,
			Username:   ud.Value,
			Timestamp: &timestamp{
				Time: time.Unix(int64(ud.Timestamp), 0),
			},
		}
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

func (t *timestamp) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(t.String())), nil
}

func (t *timestamp) UnmarshalJSON(b []byte) error {
	var unixTimestamp int64
	if err := json.Unmarshal(b, &unixTimestamp); err != nil {
		return err
	}
	t.Time = time.Unix(unixTimestamp, 0)
	return nil
}

func (t *timestamp) MarshalYAML() (interface{}, error) {
	return t.String(), nil
}

func (t *timestamp) String() string {
	return t.Format(time.DateOnly)
}

type userOriginal struct {
	Href      string `json:"href"`
	Value     string `json:"value"`
	Timestamp int    `json:"timestamp"`
}

type user struct {
	ProfileUrl string     `json:"profileUrl" yaml:"profileUrl"`
	Username   string     `json:"username" yaml:"username"`
	Timestamp  *timestamp `json:"timestamp" yaml:"timestamp"`
}

type userList struct {
	users         map[string]user
	showTimestamp bool
}

func (u *userList) output(format string) (*string, error) {
	switch format {
	case instagram.OutputJson:
		return u.outputJson()
	case instagram.OutputNone:
		return u.outputNone()
	case instagram.OutputTable:
		return u.outputTable()
	case instagram.OutputYaml:
		return u.outputYaml()
	default:
		return nil, fmt.Errorf("invalid output format: %s", format)
	}
}

func (u *userList) outputNone() (*string, error) {
	output := ""
	return &output, nil
}

func (u *userList) outputJson() (*string, error) {
	var users []user
	for i := range u.users {
		users = append(users, u.users[i])
	}
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return nil, err
	}
	output := string(data)
	return &output, nil
}

func (u *userList) outputTable() (*string, error) {
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
	header := table.Row{
		instagram.TableHeaderUsername,
		instagram.TableHeaderProfileUrl,
	}
	if u.showTimestamp {
		header = append(header, instagram.TableHeaderTimestamp)
	}
	usersTable := table.NewWriter()
	usersTable.SetAutoIndex(true)
	usersTable.SetStyle(table.StyleBold)
	usersTable.AppendHeader(header)
	usersTable.AppendRows(rows)
	output := usersTable.Render()
	return &output, nil
}

func (u *userList) outputYaml() (*string, error) {
	var users []user
	for i := range u.users {
		users = append(users, u.users[i])
	}
	data, err := yaml.Marshal(users)
	if err != nil {
		return nil, err
	}
	output := string(data)
	return &output, nil
}
