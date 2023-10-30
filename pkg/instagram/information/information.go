package information

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/cecobask/instagram-insights/pkg/filesystem"
	"github.com/cecobask/instagram-insights/pkg/instagram"
)

type Interface interface {
	Cleanup() error
	Load(source string) error
}

type handler struct {
	fileSystem filesystem.Fs
}

func NewHandler() Interface {
	return &handler{
		fileSystem: filesystem.NewFs(),
	}
}

func (h *handler) Cleanup() error {
	paths := []string{
		instagram.PathDataArchive,
		instagram.PathData,
	}
	for _, path := range paths {
		if err := h.fileSystem.RemoveDirectory(path); err != nil {
			return err
		}
	}
	return nil
}

func (h *handler) Load(source string) error {
	archiveURL, err := validateArchiveSource(source)
	if err != nil {
		return err
	}
	if archiveURL.Scheme == "file" {
		return h.fileSystem.Unzip(archiveURL.Path, instagram.PathData)
	}
	archiveURL, err = transformHttpUrl(archiveURL)
	if err != nil {
		return err
	}
	file, err := h.fileSystem.CreateFile(instagram.PathDataArchive)
	if err != nil {
		return err
	}
	defer file.Close()
	response, err := http.Get(archiveURL.String())
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failure downloading instagram data... http status: %s", response.Status)
	}
	if _, err = h.fileSystem.CopyToFile(file, response.Body); err != nil {
		return err
	}
	return h.fileSystem.Unzip(instagram.PathDataArchive, instagram.PathData)
}

func validateArchiveSource(source string) (*url.URL, error) {
	parsedURL, err := url.Parse(source)
	if err != nil {
		return nil, err
	}
	switch parsedURL.Scheme {
	case "http", "https", "file":
		return parsedURL, nil
	default:
		return nil, fmt.Errorf("unsupported source scheme: %s", parsedURL.Scheme)
	}

}

func transformHttpUrl(httpUrl *url.URL) (*url.URL, error) {
	switch httpUrl.Host {
	case instagram.GoogleDriveHost:
		pathSegments := strings.Split(httpUrl.Path, "/")
		if len(pathSegments) < 4 {
			return nil, fmt.Errorf("received invalid google drive source %s - it must be similar to this https://drive.google.com/file/d/8FOVUK1cYjgMnocmf7gMqXYdhBKHWLdnP", httpUrl.String())
		}
		return url.Parse(fmt.Sprintf(instagram.GoogleDriveParsedUrlFormat, pathSegments[3]))
	default:
		return httpUrl, nil
	}
}
