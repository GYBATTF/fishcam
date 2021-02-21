package main

import (
	"github.com/reujab/wallpaper"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// setBackground collects a list of pictures currently in the picDir,
// downloads the newest picture, sets it as background, and if it succeeds so far,
// sets cleans up the photo directory so only five photos exist.
func setBackground(picDir string, onDownload func(filename string)) (err error) {
	pictures, err := getPictures(picDir)
	if err != nil {
		return err
	}

	filename, err := downloadPic(picDir)
	if onDownload != nil {
		onDownload(filename)
	}

	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			err = cleanupPictures(pictures)
		}
	}()

	err = wallpaper.SetFromFile(filename)
	if err != nil {
		return err
	}

	return
}

// cleanupPictures takes an array of picture file name, that should be in sorted order from oldest->newest.
// It then deletes the oldest on until only four remain.
func cleanupPictures(pictures []string) error {
	for len(pictures) >= 5 {
		if err := os.Remove(pictures[0]); err != nil {
			return err
		}

		pictures = pictures[1:]
	}
	return nil
}

// getPictures gets the pictures in the specified directory.
func getPictures(picDir string) ([]string, error) {
	files, err := ioutil.ReadDir(picDir)
	if err != nil {
		return nil, err
	}

	var pictures []string

	for _, f := range files {
		if strings.Contains(f.Name(), "jpg") {
			_, filename := filepath.Split(f.Name())
			filename = filepath.Join(picDir, filename)
			pictures = append(pictures, filename)
		}
	}

	sort.Strings(pictures)
	return pictures, nil
}
