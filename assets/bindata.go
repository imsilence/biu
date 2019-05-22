// Code generated by go-bindata.
// sources:
// templates/task.html
// DO NOT EDIT!

package assets

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _templatesTaskHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\xd1\xb1\x4a\xc4\x40\x10\x06\xe0\xde\xa7\x18\xa7\x14\xcc\x5e\x82\xe8\x29\xbb\x69\xd4\xfa\x0e\xbc\xc6\x72\x93\xcc\x91\x60\x36\x09\xc9\x44\xbc\x47\x10\x04\xd1\xda\xc6\xca\xce\x4e\x11\x0f\xdf\xe6\x88\x3e\x86\x24\x8b\x5e\x4e\x14\xe4\xa6\xd9\x19\x98\xfd\x58\xf6\x97\x9b\x47\xa3\xc3\xc9\xe9\xf8\x18\x62\x36\xa9\xbf\x21\xed\x01\x00\x20\x63\xd2\x91\x6d\xbb\xd1\x10\x6b\x08\x63\x5d\x56\xc4\x0a\x6b\x9e\x6e\x0f\x11\x44\x6f\x81\x13\x4e\xc9\x0f\x92\x5a\x0a\xdb\x5a\x46\x2c\x1d\x19\xe4\xd1\xac\x77\x63\x9a\x97\x06\x0c\x71\x9c\x47\x0a\xc7\xa3\x93\x09\x82\x0e\x39\xc9\x33\x85\x82\x75\x75\x26\x70\xb9\xdc\x56\x73\x7d\xbb\x98\x3f\x1f\x80\x4c\xb2\xa2\x66\xe0\x59\x41\x0a\x99\x2e\x18\x21\xd3\x86\x14\x16\x08\xe7\x3a\xad\x49\xe1\x16\x42\x91\xea\x90\xe2\x3c\x8d\xa8\xec\x66\xe1\xcb\xa0\xec\xbf\xb8\xad\xc5\xcb\xbc\xb9\x7b\xfd\x17\xf9\x43\x74\xbd\x3d\x67\xe0\x0c\x1c\xf7\x0f\xf9\xfd\xed\xa6\x79\x7c\x5a\x4b\xde\xf7\x1c\x77\x77\xd8\xda\xc2\xdb\xc1\x5f\xf5\x15\xb5\xaa\x03\x93\xf0\x37\xd7\x5c\x3e\x7c\xdc\x5f\xad\x86\x23\xda\xbf\xfe\x4a\xc4\xc6\x20\x45\x17\xf6\x67\x00\x00\x00\xff\xff\xc8\x66\x97\x0b\x03\x02\x00\x00")

func templatesTaskHtmlBytes() ([]byte, error) {
	return bindataRead(
		_templatesTaskHtml,
		"templates/task.html",
	)
}

func templatesTaskHtml() (*asset, error) {
	bytes, err := templatesTaskHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/task.html", size: 515, mode: os.FileMode(438), modTime: time.Unix(1558543264, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"templates/task.html": templatesTaskHtml,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"templates": &bintree{nil, map[string]*bintree{
		"task.html": &bintree{templatesTaskHtml, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

