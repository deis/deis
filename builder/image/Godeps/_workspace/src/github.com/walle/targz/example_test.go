package targz_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/walle/targz"
)

func Example() {
	// Create a temporary file structiure we can use
	tmpDir, dirToCompress := createExampleData()

	// Comress a folder to my_archive.tar.gz
	err := targz.Compress(dirToCompress, filepath.Join(tmpDir, "my_archive.tar.gz"))
	if err != nil {
		panic(err)
	}

	// Extract my_archive.tar.gz to a new folder called extracted
	err = targz.Extract(filepath.Join(tmpDir, "my_archive.tar.gz"), filepath.Join(tmpDir, "extracted"))
	if err != nil {
		panic(err)
	}
}

func createExampleData() (string, string) {
	tmpDir, err := ioutil.TempDir("", "targz-example")
	if err != nil {
		panic(err)
	}

	directory := filepath.Join(tmpDir, "my_folder")
	subDirectory := filepath.Join(directory, "my_sub_folder")
	err = os.MkdirAll(subDirectory, 0755)
	if err != nil {
		panic(err)
	}

	_, err = os.Create(filepath.Join(subDirectory, "my_file.txt"))
	if err != nil {
		panic(err)
	}

	return tmpDir, directory
}
