package followdata

import (
	"fmt"
	"testing"
	"time"

	"github.com/cecobask/instagram-insights/pkg/filesystem"
	"github.com/cecobask/instagram-insights/pkg/instagram"
	"github.com/stretchr/testify/assert"
)

func Test_handler_Followers(t *testing.T) {
	type fields struct {
		fileSystem *filesystem.MockFs
		followData *followData
	}
	tests := []struct {
		name         string
		expectations func(f *fields)
		assertions   func(t *testing.T, f *fields)
		wantErr      bool
	}{
		{
			name: "succeeds to output followers",
			expectations: func(f *fields) {
				f.fileSystem.On("FindFiles", instagram.PathFollowers).Return([]string{"file1"}, nil)
				f.fileSystem.On("ReadFile", "file1").Return([]byte(`[{"string_list_data":[{"href":"https://www.instagram.com/username","value":"username","timestamp":0}]}]`), nil)
			},
			assertions: func(t *testing.T, f *fields) {
				f.fileSystem.AssertNumberOfCalls(t, "FindFiles", 1)
				f.fileSystem.AssertNumberOfCalls(t, "ReadFile", 1)
			},
			wantErr: false,
		},
		{
			name: "fails to find files",
			expectations: func(f *fields) {
				f.fileSystem.On("FindFiles", instagram.PathFollowers).Return(nil, fmt.Errorf("fails to find files"))
			},
			assertions: func(t *testing.T, f *fields) {
				f.fileSystem.AssertNumberOfCalls(t, "FindFiles", 1)
			},
			wantErr: true,
		},
		{
			name: "fails to read file",
			expectations: func(f *fields) {
				f.fileSystem.On("FindFiles", instagram.PathFollowers).Return([]string{"file1"}, nil)
				f.fileSystem.On("ReadFile", "file1").Return(nil, fmt.Errorf("fails to read file"))
			},
			assertions: func(t *testing.T, f *fields) {
				f.fileSystem.AssertNumberOfCalls(t, "FindFiles", 1)
				f.fileSystem.AssertNumberOfCalls(t, "ReadFile", 1)
			},
			wantErr: true,
		},
		{
			name: "fails to hydrate followers",
			expectations: func(f *fields) {
				f.fileSystem.On("FindFiles", instagram.PathFollowers).Return([]string{"file1"}, nil)
				f.fileSystem.On("ReadFile", "file1").Return([]byte(""), nil)
			},
			assertions: func(t *testing.T, f *fields) {
				f.fileSystem.AssertNumberOfCalls(t, "FindFiles", 1)
				f.fileSystem.AssertNumberOfCalls(t, "ReadFile", 1)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fields{
				fileSystem: &filesystem.MockFs{},
				followData: &followData{},
			}
			h := &handler{
				fileSystem: f.fileSystem,
				followData: f.followData,
			}
			if tt.expectations != nil {
				tt.expectations(f)
			}
			if _, err := h.Followers(instagram.NewEmptyOptions()); (err != nil) != tt.wantErr {
				t.Errorf("Followers() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.assertions != nil {
				tt.assertions(t, f)
			}
		})
	}
}

func Test_handler_Following(t *testing.T) {
	type fields struct {
		fileSystem *filesystem.MockFs
		followData *followData
	}
	tests := []struct {
		name         string
		expectations func(f *fields)
		assertions   func(t *testing.T, f *fields)
		wantErr      bool
	}{
		{
			name: "succeeds to output following",
			expectations: func(f *fields) {
				f.fileSystem.On("ReadFile", instagram.PathFollowing).Return([]byte(`{"relationships_following":[{"string_list_data":[{"href":"https://www.instagram.com/username","value":"username","timestamp":0}]}]}`), nil)
			},
			assertions: func(t *testing.T, f *fields) {
				f.fileSystem.AssertNumberOfCalls(t, "ReadFile", 1)
			},
			wantErr: false,
		},
		{
			name: "fails to read file",
			expectations: func(f *fields) {
				f.fileSystem.On("ReadFile", instagram.PathFollowing).Return(nil, fmt.Errorf("fails to read file"))
			},
			assertions: func(t *testing.T, f *fields) {
				f.fileSystem.AssertNumberOfCalls(t, "ReadFile", 1)
			},
			wantErr: true,
		},
		{
			name: "fails to hydrate following",
			expectations: func(f *fields) {
				f.fileSystem.On("ReadFile", instagram.PathFollowing).Return([]byte(""), nil)
			},
			assertions: func(t *testing.T, f *fields) {
				f.fileSystem.AssertNumberOfCalls(t, "ReadFile", 1)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fields{
				fileSystem: &filesystem.MockFs{},
				followData: &followData{},
			}
			h := &handler{
				fileSystem: f.fileSystem,
				followData: f.followData,
			}
			if tt.expectations != nil {
				tt.expectations(f)
			}
			if _, err := h.Following(instagram.NewEmptyOptions()); (err != nil) != tt.wantErr {
				t.Errorf("Following() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.assertions != nil {
				tt.assertions(t, f)
			}
		})
	}
}

func Test_handler_Unfollowers(t *testing.T) {
	type fields struct {
		fileSystem *filesystem.MockFs
		followData *followData
	}
	tests := []struct {
		name         string
		expectations func(f *fields)
		assertions   func(t *testing.T, f *fields)
		wantErr      bool
	}{
		{
			name: "succeeds to output unfollowers",
			expectations: func(f *fields) {
				f.fileSystem.On("FindFiles", instagram.PathFollowers).Return([]string{"file1"}, nil)
				f.fileSystem.On("ReadFile", "file1").Return([]byte(`[{"string_list_data":[{"href":"https://www.instagram.com/username1","value":"username1","timestamp":0}]}]`), nil).Once()
				f.fileSystem.On("ReadFile", instagram.PathFollowing).Return([]byte(`{"relationships_following":[{"string_list_data":[{"href":"https://www.instagram.com/username2","value":"username2","timestamp":0}]}]}`), nil).Once()
			},
			assertions: func(t *testing.T, f *fields) {
				f.fileSystem.AssertNumberOfCalls(t, "FindFiles", 1)
				f.fileSystem.AssertNumberOfCalls(t, "ReadFile", 2)
			},
			wantErr: false,
		},
		{
			name: "fails to get followers",
			expectations: func(f *fields) {
				f.fileSystem.On("FindFiles", instagram.PathFollowers).Return(nil, fmt.Errorf("fails to find files"))
			},
			assertions: func(t *testing.T, f *fields) {
				f.fileSystem.AssertNumberOfCalls(t, "FindFiles", 1)
			},
			wantErr: true,
		},
		{
			name: "fails to get following",
			expectations: func(f *fields) {
				f.fileSystem.On("FindFiles", instagram.PathFollowers).Return([]string{"file1"}, nil)
				f.fileSystem.On("ReadFile", "file1").Return([]byte(`[{"string_list_data":[{"href":"https://www.instagram.com/username","value":"username","timestamp":0}]}]`), nil).Once()
				f.fileSystem.On("ReadFile", instagram.PathFollowing).Return(nil, fmt.Errorf("fails to read file")).Once()
			},
			assertions: func(t *testing.T, f *fields) {
				f.fileSystem.AssertNumberOfCalls(t, "FindFiles", 1)
				f.fileSystem.AssertNumberOfCalls(t, "ReadFile", 2)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fields{
				fileSystem: &filesystem.MockFs{},
				followData: &followData{},
			}
			h := &handler{
				fileSystem: f.fileSystem,
				followData: f.followData,
			}
			if tt.expectations != nil {
				tt.expectations(f)
			}
			if _, err := h.Unfollowers(instagram.NewEmptyOptions()); (err != nil) != tt.wantErr {
				t.Errorf("Unfollowers() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.assertions != nil {
				tt.assertions(t, f)
			}
		})
	}
}

func Test_userList_output(t *testing.T) {
	type args struct {
		format string
	}
	dummyUser := user{
		ProfileUrl: "https://www.instagram.com/username",
		Username:   "username",
		Timestamp:  &timestamp{},
	}
	u := userList{
		users:         []user{dummyUser},
		showTimestamp: true,
	}
	uNoTimestamp := userList{
		users:         []user{dummyUser},
		showTimestamp: false,
	}
	tests := []struct {
		name    string
		u       userList
		args    args
		wantErr bool
	}{
		{
			name: "succeeds to output json",
			u:    u,
			args: args{
				format: instagram.OutputJson,
			},
			wantErr: false,
		},
		{
			name: "succeeds to output json without timestamp",
			u:    uNoTimestamp,
			args: args{
				format: instagram.OutputJson,
			},
			wantErr: false,
		},
		{
			name: "succeeds to output table",
			u:    u,
			args: args{
				format: instagram.OutputTable,
			},
			wantErr: false,
		},
		{
			name: "succeeds to output yaml",
			u:    u,
			args: args{
				format: instagram.OutputYaml,
			},
			wantErr: false,
		},
		{
			name: "succeeds to output yaml without timestamp",
			u:    uNoTimestamp,
			args: args{
				format: instagram.OutputYaml,
			},
			wantErr: false,
		},
		{
			name: "succeeds to output none",
			u:    u,
			args: args{
				format: instagram.OutputNone,
			},
			wantErr: false,
		},
		{
			name: "fails to output invalid format",
			u:    u,
			args: args{
				format: "invalid",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.u.output(tt.args.format); (err != nil) != tt.wantErr {
				t.Errorf("output() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_timestamp_UnmarshalJSON(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "succeeds to unmarshal json",
			args: args{
				b: []byte("1697474963"),
			},
			wantErr: false,
		},
		{
			name: "fails to unmarshal json",
			args: args{
				b: []byte("invalid"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := &timestamp{}
			if err := ts.UnmarshalJSON(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userList_Sort(t *testing.T) {
	type fields struct {
		users []user
	}
	type args struct {
		field string
		order string
	}
	timeNow := time.Now()
	testUsers := []user{
		{
			ProfileUrl: "https://www.instagram.com/username1",
			Username:   "username1",
			Timestamp: &timestamp{
				Time: timeNow,
			},
		},
		{
			ProfileUrl: "https://www.instagram.com/username2",
			Username:   "username2",
			Timestamp: &timestamp{
				Time: timeNow.Add(time.Hour),
			},
		},
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		assertions func(t *testing.T, ul *userList)
	}{
		{
			name: "succeeds to sort by timestamp ascending",
			fields: fields{
				users: testUsers,
			},
			args: args{
				field: instagram.FieldTimestamp,
				order: instagram.OrderAsc,
			},
			assertions: func(t *testing.T, ul *userList) {
				assert.Equal(t, "username1", ul.users[0].Username)
				assert.Equal(t, "username2", ul.users[1].Username)
			},
		},
		{
			name: "succeeds to sort by timestamp descending",
			fields: fields{
				users: testUsers,
			},
			args: args{
				field: instagram.FieldTimestamp,
				order: instagram.OrderDesc,
			},
			assertions: func(t *testing.T, ul *userList) {
				assert.Equal(t, "username2", ul.users[0].Username)
				assert.Equal(t, "username1", ul.users[1].Username)
			},
		},
		{
			name: "succeeds to sort by username ascending",
			fields: fields{
				users: testUsers,
			},
			args: args{
				field: instagram.FieldUsername,
				order: instagram.OrderAsc,
			},
			assertions: func(t *testing.T, ul *userList) {
				assert.Equal(t, "username1", ul.users[0].Username)
				assert.Equal(t, "username2", ul.users[1].Username)
			},
		},
		{
			name: "succeeds to sort by username descending",
			fields: fields{
				users: testUsers,
			},
			args: args{
				field: instagram.FieldUsername,
				order: instagram.OrderDesc,
			},
			assertions: func(t *testing.T, ul *userList) {
				assert.Equal(t, "username2", ul.users[0].Username)
				assert.Equal(t, "username1", ul.users[1].Username)
			},
		},
		{
			name: "succeeds to sort by timestamp when field is unknown",
			fields: fields{
				users: testUsers,
			},
			args: args{
				field: "unknown",
				order: instagram.OrderAsc,
			},
			assertions: func(t *testing.T, ul *userList) {
				assert.Equal(t, "username1", ul.users[0].Username)
				assert.Equal(t, "username2", ul.users[1].Username)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ul := &userList{
				users: tt.fields.users,
			}
			ul.Sort(tt.args.field, tt.args.order)
			tt.assertions(t, ul)
		})
	}
}
