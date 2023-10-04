package instagram

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/cecobask/instagram-insights/pkg/file"
)

func FetchInstagramInformation(archiveURL string) error {
	archiveURL, err := parseArchiveURL(archiveURL)
	if err != nil {
		return err
	}
	if err = file.DownloadFile(archiveURL, pathDataArchive); err != nil {
		return err
	}
	return file.UnzipFile(pathDataArchive, pathData)
}

func CleanupInstagramInformation() error {
	filePaths := []string{
		pathDataArchive,
		pathData,
	}
	return file.CleanupFilePaths(filePaths)
}

func parseArchiveURL(archiveURL string) (string, error) {
	parsedURL, err := url.Parse(archiveURL)
	if err != nil {
		return "", err
	}
	switch parsedURL.Host {
	case googleDriveHost:
		pathSegments := strings.Split(parsedURL.Path, "/")
		if len(pathSegments) < 4 {
			return "", fmt.Errorf("received invalid google drive url %s - it must be similar to this https://drive.google.com/file/d/8FOVUK1cYjgMnocmf7gMqXYdhBKHWLdnP", archiveURL)
		}
		return fmt.Sprintf(googleDriveParsedUrlFormat, pathSegments[3]), err
	default:
		return archiveURL, nil
	}
}
