package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

func bash_example_bash() ([]byte, error) {
	return bindata_read([]byte{
		0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x00, 0xff, 0x4c, 0xce,
		0xb1, 0x4e, 0x45, 0x21, 0x0c, 0xc6, 0xf1, 0xf9, 0xf4, 0x29, 0x1a, 0xe2,
		0xa0, 0x03, 0x97, 0xfd, 0x4e, 0xea, 0xe4, 0x6b, 0x70, 0x8f, 0x3d, 0x50,
		0x03, 0x94, 0xb4, 0xe0, 0x19, 0xd4, 0x77, 0x17, 0xe3, 0xa0, 0x63, 0x9b,
		0x7f, 0x7e, 0xf9, 0x8c, 0x06, 0x7a, 0x12, 0xec, 0xdc, 0xe9, 0x88, 0x5c,
		0x00, 0x32, 0x95, 0x22, 0xfe, 0x16, 0x2d, 0xdf, 0x3f, 0xe0, 0x07, 0x6c,
		0xb4, 0x67, 0x41, 0xf7, 0xf2, 0xf3, 0xc5, 0x53, 0xb4, 0xbc, 0xe2, 0xa1,
		0x52, 0xf1, 0x79, 0x05, 0x0e, 0xbe, 0x00, 0x6a, 0xe4, 0xf6, 0xbf, 0x7c,
		0xd2, 0x34, 0x2b, 0xb5, 0x61, 0x57, 0x87, 0xee, 0xee, 0xd1, 0xc1, 0xf6,
		0x27, 0xe2, 0x27, 0x2a, 0xbd, 0x93, 0x1a, 0xc1, 0xb6, 0x4f, 0x2d, 0xe8,
		0x0d, 0xf3, 0x18, 0xdd, 0xae, 0x21, 0xc4, 0xce, 0x97, 0xc4, 0x23, 0xcf,
		0xdb, 0x65, 0x97, 0x1a, 0x94, 0xba, 0x58, 0xe8, 0x2a, 0x49, 0x79, 0xd6,
		0x90, 0x7e, 0x01, 0xd2, 0x45, 0xbc, 0x99, 0x34, 0xdf, 0x85, 0xdb, 0x58,
		0x67, 0x90, 0xb3, 0x91, 0x86, 0x22, 0x89, 0xdb, 0x9a, 0xf3, 0x1d, 0x00,
		0x00, 0xff, 0xff, 0xca, 0x1c, 0xba, 0x45, 0xd0, 0x00, 0x00, 0x00,
	},
		"bash/example.bash",
	)
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
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
var _bindata = map[string]func() ([]byte, error){
	"bash/example.bash": bash_example_bash,
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
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() ([]byte, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"bash": &_bintree_t{nil, map[string]*_bintree_t{
		"example.bash": &_bintree_t{bash_example_bash, map[string]*_bintree_t{
		}},
	}},
}}
