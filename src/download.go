package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	url = "https://www.fishcam.com/live_images/liveimage_1280.jpg?camera=0"
)

func init() {
	// This is usually not a good idea, but is required as fishcam.com has certificate issues.
	// Since all we're doing is downloading a photo from the site and nothing sensitive, we'll take the risk.
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

// downloadPic downloads a photo from fishcam.com, and save it to the specified picDir.
func downloadPic(picDir string) (string, error) {
	filename := fmt.Sprintf("fish-%d.jpg", time.Now().UnixNano()/10000)
	filename = filepath.Join(picDir, filename)

	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()

		if out, err := os.Create(filename); err == nil {
			defer out.Close()
			_, err = io.Copy(out, resp.Body)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	} else {
		return "", err
	}

	return filename, nil
}
