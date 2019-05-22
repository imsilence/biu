package config

import (
	"os"
	"path/filepath"

	"github.com/kardianos/osext"
)

var PLUGIN_NAME_PATTERN string = "*.json"

var HOME string
var CWD string
var PLUGIN_DIR string
var PLUGIN_PATH_PATTERN string

func init() {
	HOME, _ = osext.ExecutableFolder()
	CWD, _ := os.Getwd()
	for _, dir := range []string{HOME, CWD} {
		path := filepath.Join(dir, "plugins")
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			PLUGIN_DIR = path
			break
		}
	}
}
