package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
	"os/user"
	"path/filepath"
	"sync"
	"time"
)

var (
	isRead bool
	lock   sync.Mutex
)

// readConfig reads the config file, writing it with default values if it does not exist.
func readConfig(duration chan time.Duration, handleErr func(error)) {
	isRead = false
	lock.Lock()
	defer lock.Unlock()

	setDefaults()

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(getConfigDir())
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err = os.MkdirAll(getConfigDir(), os.ModePerm)

			if err == nil {
				err = viper.SafeWriteConfig()
			}
		}

		if err != nil {
			panic(err)
		}
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		err := sendUpdateDuration(duration)
		if err != nil && handleErr != nil {
			handleErr(err)
		}
	})

	err := sendUpdateDuration(duration)
	if err != nil && handleErr != nil {
		handleErr(err)
	}

	isRead = true
}

// setDefaults sets the default config values.
func setDefaults() {
	viper.SetDefault("picDir", getConfigDir())

	viper.SetDefault("updateInterval", map[string]int{
		"minutes": 30,
		"hours":   0,
		"days":    0,
	})
}

// getConfigDir returns the default config directory.
func getConfigDir() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	return filepath.Join(usr.HomeDir, "/.local/etc/fishcam/")
}

// getPicDir waits until the config file has been completely read
// (if it hasn't been read already) and then return the directory
// where pictures are to be stored.
func getPicDir() string {
	waitUntilConfigIsRead()
	return viper.GetString("picDir")
}

// waitUntilConfigIsRead blocks until the config has been completely read.
func waitUntilConfigIsRead() {
	for !isRead {
		lock.Lock()
		lock.Unlock()
	}
}

// sendUpdateDuration gets the current update interval from
// the config file and send it through the specified channel.
// If and error occurs parsing the duration, nothing is sent.
func sendUpdateDuration(duration chan time.Duration) error {
	m := viper.GetInt("updateInterval.minutes")
	h := viper.GetInt("updateInterval.hours")
	h += viper.GetInt("updateInterval.days") * 24

	d, err := time.ParseDuration(fmt.Sprintf("%dh%dm", h, m))

	if err == nil {
		duration <- d
	}

	return err
}
