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
	Download(url string) error
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

func (h *handler) Download(url string) error {
	archiveURL, err := parseArchiveURL(url)
	if err != nil {
		return err
	}
	file, err := h.fileSystem.CreateFile(instagram.PathDataArchive)
	if err != nil {
		return err
	}
	defer file.Close()
	response, err := http.Get(archiveURL)
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

func parseArchiveURL(archiveURL string) (string, error) {
	parsedURL, err := url.Parse(archiveURL)
	if err != nil {
		return "", err
	}
	switch parsedURL.Host {
	case instagram.GoogleDriveHost:
		pathSegments := strings.Split(parsedURL.Path, "/")
		if len(pathSegments) < 4 {
			return "", fmt.Errorf("received invalid google drive url %s - it must be similar to this https://drive.google.com/file/d/8FOVUK1cYjgMnocmf7gMqXYdhBKHWLdnP", archiveURL)
		}
		return fmt.Sprintf(instagram.GoogleDriveParsedUrlFormat, pathSegments[3]), err
	default:
		return archiveURL, nil
	}
}
