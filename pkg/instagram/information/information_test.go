package information

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
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

func Test_informationHandler_Load(t *testing.T) {
	type fields struct {
		fileSystem *filesystem.MockFs
	}
	type args struct {
		source string
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
			name: "succeeds to load information from file source",
			args: args{
				source: "file:///home/username/Desktop/instagram_data.zip",
			},
			expectations: func(f *fields) {
				f.fileSystem.On("Unzip", "/home/username/Desktop/instagram_data.zip", instagram.PathData).Return(nil)
			},
			assertions: func(t *testing.T, f *fields) {
				f.fileSystem.AssertNumberOfCalls(t, "Unzip", 1)
			},
			wantErr: false,
		},
		{
			name: "succeeds to load information from http source",
			args: args{
				source: "", // dynamically set at runtime
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
			name: "fails to transform http url",
			args: args{
				source: "https://drive.google.com/file/d",
			},
			wantErr: true,
		},
		{
			name: "fails to validate archive source",
			args: args{
				source: "ftp://example.com",
			},
			wantErr: true,
		},
		{
			name: "fails to create file",
			args: args{
				source: "",
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
				source: "",
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
				source: "",
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
				source: "",
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
			if tt.args.source == "" {
				server := createHttpServerWithStatus(tt.httpStatusCode)
				defer server.Close()
				tt.args.source = server.URL
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
			if err := h.Load(tt.args.source); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.assertions != nil {
				tt.assertions(t, f)
			}
		})
	}
}

func Test_transformHttpUrl(t *testing.T) {
	type args struct {
		httpUrl *url.URL
	}
	tests := []struct {
		name    string
		args    args
		want    *url.URL
		wantErr bool
	}{
		{
			name: "parses generic url",
			args: args{
				httpUrl: &url.URL{
					Scheme: "https",
					Host:   "example.com",
				},
			},
			want: &url.URL{
				Scheme: "https",
				Host:   "example.com",
			},
			wantErr: false,
		},
		{
			name: "transforms google drive url",
			args: args{
				httpUrl: &url.URL{
					Scheme: "https",
					Host:   instagram.GoogleDriveHost,
					Path:   "/file/d/gMFOVXYdhBK8gMnYjqcmocf7HWLUK1dnP",
				},
			},
			want: &url.URL{
				Scheme:   "https",
				Host:     instagram.GoogleDriveHost,
				Path:     "/u/0/uc",
				RawQuery: "id=gMFOVXYdhBK8gMnYjqcmocf7HWLUK1dnP&export=download&confirm=t",
			},
			wantErr: false,
		},
		{
			name: "fails to transform google drive url",
			args: args{
				httpUrl: &url.URL{
					Scheme: "https",
					Host:   instagram.GoogleDriveHost,
					Path:   "/file/d",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := transformHttpUrl(tt.args.httpUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("transformHttpUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && got.String() != tt.want.String() {
				t.Errorf("transformHttpUrl() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateArchiveSource(t *testing.T) {
	type args struct {
		source string
	}
	tests := []struct {
		name    string
		args    args
		want    *url.URL
		wantErr bool
	}{
		{
			name: "validates archive source",
			args: args{
				source: "https://example.com",
			},
			want: &url.URL{
				Scheme: "https",
				Host:   "example.com",
			},
			wantErr: false,
		},
		{
			name: "fails to parse url",
			args: args{
				source: string(rune(0x7f)),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "fails to validate scheme",
			args: args{
				source: "ftp://example.com",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateArchiveSource(tt.args.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateArchiveSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validateArchiveSource() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func createHttpServerWithStatus(statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
	}))
}
