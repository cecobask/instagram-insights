package information

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/cecobask/instagram-insights/pkg/filesystem"
	"github.com/cecobask/instagram-insights/pkg/instagram"
)

func Test_informationHandler_Cleanup(t *testing.T) {
	type fields struct {
		fileSystem *filesystem.MockFs
	}
	tests := []struct {
		name         string
		expectations func(f *fields)
		assertions   func(t *testing.T, f *fields)
		wantErr      bool
	}{
		{
			name: "succeeds to cleanup paths",
			expectations: func(f *fields) {
				f.fileSystem.On("RemoveDirectory", instagram.PathDataArchive).Return(nil).Once()
				f.fileSystem.On("RemoveDirectory", instagram.PathData).Return(nil).Once()
			},
			assertions: func(t *testing.T, f *fields) {
				f.fileSystem.AssertNumberOfCalls(t, "RemoveDirectory", 2)
			},
			wantErr: false,
		},
		{
			name: "fails to remove directory",
			expectations: func(f *fields) {
				f.fileSystem.On("RemoveDirectory", instagram.PathDataArchive).Return(fmt.Errorf("fails to remove directory"))
			},
			assertions: func(t *testing.T, f *fields) {
				f.fileSystem.AssertNumberOfCalls(t, "RemoveDirectory", 1)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fields{
				fileSystem: &filesystem.MockFs{},
			}
			h := &handler{
				fileSystem: f.fileSystem,
			}
			if tt.expectations != nil {
				tt.expectations(f)
			}
			if err := h.Cleanup(); (err != nil) != tt.wantErr {
				t.Errorf("Cleanup() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.assertions != nil {
				tt.assertions(t, f)
			}
		})
	}
}

func Test_informationHandler_Download(t *testing.T) {
	type fields struct {
		fileSystem *filesystem.MockFs
	}
	type args struct {
		url string
	}
	tests := []struct {
		name           string
		args           args
		httpStatusCode int
		expectations   func(f *fields)
		assertions     func(t *testing.T, f *fields)
		wantErr        bool
	}{
		{
			name: "succeeds to download file",
			args: args{
				url: "",
			},
			httpStatusCode: http.StatusOK,
			expectations: func(f *fields) {
				dummyFile := &os.File{}
				f.fileSystem.On("CreateFile", instagram.PathDataArchive).Return(dummyFile, nil)
				f.fileSystem.On("CopyToFile", dummyFile, http.NoBody).Return(int64(0), nil)
				f.fileSystem.On("Unzip", instagram.PathDataArchive, instagram.PathData).Return(nil)
			},
			assertions: func(t *testing.T, f *fields) {
				f.fileSystem.AssertNumberOfCalls(t, "CreateFile", 1)
				f.fileSystem.AssertNumberOfCalls(t, "CopyToFile", 1)
				f.fileSystem.AssertNumberOfCalls(t, "Unzip", 1)
			},
			wantErr: false,
		},
		{
			name: "fails to parse archive url",
			args: args{
				url: string(rune(0x7f)),
			},
			wantErr: true,
		},
		{
			name: "fails to create file",
			args: args{
				url: "",
			},
			expectations: func(f *fields) {
				f.fileSystem.On("CreateFile", instagram.PathDataArchive).Return(nil, fmt.Errorf("fails to create file"))
			},
			assertions: func(t *testing.T, f *fields) {
				f.fileSystem.AssertNumberOfCalls(t, "CreateFile", 1)
			},
			wantErr: true,
		},
		{
			name: "fails to issue http get request",
			args: args{
				url: "example.com",
			},
			expectations: func(f *fields) {
				f.fileSystem.On("CreateFile", instagram.PathDataArchive).Return(&os.File{}, nil)
			},
			assertions: func(t *testing.T, f *fields) {
				f.fileSystem.AssertNumberOfCalls(t, "CreateFile", 1)
			},
			wantErr: true,
		},
		{
			name: "fails to get healthy http status code",
			args: args{
				url: "",
			},
			httpStatusCode: http.StatusNotFound,
			expectations: func(f *fields) {
				f.fileSystem.On("CreateFile", instagram.PathDataArchive).Return(&os.File{}, nil)

			},
			assertions: func(t *testing.T, f *fields) {
				f.fileSystem.AssertNumberOfCalls(t, "CreateFile", 1)
			},
			wantErr: true,
		},
		{
			name: "fails to copy http response body to file",
			args: args{
				url: "",
			},
			httpStatusCode: http.StatusOK,
			expectations: func(f *fields) {
				dummyFile := &os.File{}
				f.fileSystem.On("CreateFile", instagram.PathDataArchive).Return(dummyFile, nil)
				f.fileSystem.On("CopyToFile", dummyFile, http.NoBody).Return(int64(0), fmt.Errorf("fails to copy http response body to file"))
			},
			assertions: func(t *testing.T, f *fields) {
				f.fileSystem.AssertNumberOfCalls(t, "CreateFile", 1)
				f.fileSystem.AssertNumberOfCalls(t, "CopyToFile", 1)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.url == "" {
				server := createHttpServerWithStatus(tt.httpStatusCode)
				defer server.Close()
				tt.args.url = server.URL
			}
			f := &fields{
				fileSystem: &filesystem.MockFs{},
			}
			h := &handler{
				fileSystem: f.fileSystem,
			}
			if tt.expectations != nil {
				tt.expectations(f)
			}
			if err := h.Download(tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("Download() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.assertions != nil {
				tt.assertions(t, f)
			}
		})
	}
}

func Test_parseArchiveURL(t *testing.T) {
	type args struct {
		archiveURL string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "valid generic url",
			args: args{
				archiveURL: "https://example.com",
			},
			want:    "https://example.com",
			wantErr: false,
		},
		{
			name: "google drive url",
			args: args{
				archiveURL: "https://drive.google.com/file/d/gMFOVXYdhBK8gMnYjqcmocf7HWLUK1dnP",
			},
			want:    "https://drive.google.com/u/0/uc?id=gMFOVXYdhBK8gMnYjqcmocf7HWLUK1dnP&export=download&confirm=t",
			wantErr: false,
		},
		{
			name: "google drive url with invalid path",
			args: args{
				archiveURL: "https://drive.google.com/file/d",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "invalid url",
			args: args{
				archiveURL: string(rune(0x7f)),
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseArchiveURL(tt.args.archiveURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseArchiveURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseArchiveURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func createHttpServerWithStatus(statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
	}))
}
