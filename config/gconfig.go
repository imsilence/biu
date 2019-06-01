package config

import (
	"os"
	"path/filepath"

	"github.com/kardianos/osext"
)

var POC_NAME_PATTERN string = "*.json"

var HOME string
var CWD string
var POC_DIR string

func init() {
	HOME, _ = osext.ExecutableFolder()
	CWD, _ := os.Getwd()
	for _, dir := range []string{HOME, CWD} {
		path := filepath.Join(dir, "poc")
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			POC_DIR = path
			break
		}
	}
}
