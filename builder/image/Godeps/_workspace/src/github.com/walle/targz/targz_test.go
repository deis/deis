package targz

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// Check if path exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	switch {
	case err == nil:
		return true, nil
	case os.IsNotExist(err):
		return false, nil
	default:
		return false, err
	}
}

func Test_CompressAndExtract(t *testing.T) {
	tmpDir, dirToCompress := createTestData()
	defer os.RemoveAll(tmpDir)

	structureBefore := directoryStructureString(dirToCompress)

	err := Compress(dirToCompress, filepath.Join(tmpDir, "my_archive.tar.gz"))
	if err != nil {
		t.Errorf("Compress error: %s", err)
	}

	err = Extract(filepath.Join(tmpDir, "my_archive.tar.gz"), filepath.Join(tmpDir, "extracted"))
	if err != nil {
		t.Errorf("Extract error: %s", err)
	}

	structureAfter := directoryStructureString(filepath.Join(tmpDir, "extracted", "my_folder"))

	if structureAfter != structureBefore {
		t.Errorf("Directory structure before compress and after extract does not match. Before {%s}, After {%s}", structureBefore, structureAfter)
	}
}

func Test_CompressAndExtractWithTrailingSlash(t *testing.T) {
	tmpDir, dirToCompress := createTestData()
	defer os.RemoveAll(tmpDir)

	structureBefore := directoryStructureString(dirToCompress)

	err := Compress(dirToCompress+"/", filepath.Join(tmpDir, "my_archive.tar.gz"))
	if err != nil {
		t.Errorf("Compress error: %s", err)
	}

	err = Extract(filepath.Join(tmpDir, "my_archive.tar.gz"), filepath.Join(tmpDir, "extracted/"))
	if err != nil {
		t.Errorf("Extract error: %s", err)
	}

	structureAfter := directoryStructureString(filepath.Join(tmpDir, "extracted", "my_folder"))

	if structureAfter != structureBefore {
		t.Errorf("Directory structure before compress and after extract does not match. Before {%s}, After {%s}", structureBefore, structureAfter)
	}
}

func Test_GivesErrorIfOutputIsFile(t *testing.T) {
	tmpDir, dirToCompress := createTestData()
	defer os.RemoveAll(tmpDir)

	err := Compress(dirToCompress, filepath.Join(tmpDir, "my_archive.tar.gz"))
	if err != nil {
		t.Errorf("Compress error: %s", err)
	}

	err = Extract(filepath.Join(tmpDir, "my_archive.tar.gz"), filepath.Join(tmpDir, "my_archive.tar.gz"))
	if err == nil {
		t.Errorf("Should say that my_archive.tar.gz isn't a directory")
	}
}

func Test_GivesErrorIfInputDirDoesNotExist(t *testing.T) {
	tmpDir, dirToCompress := createTestData()
	defer os.RemoveAll(tmpDir)

	os.RemoveAll(dirToCompress)

	err := Compress(dirToCompress, filepath.Join(tmpDir, "my_archive.tar.gz"))
	if err == nil {
		t.Errorf("Should say that %s doesn't exist", dirToCompress)
	}
}

func Test_GivesErrorIfInputDirIsEmpty(t *testing.T) {
	tmpDir, dirToCompress := createTestData()
	defer os.RemoveAll(tmpDir)

	os.RemoveAll(filepath.Join(dirToCompress, "my_sub_folder"))

	err := Compress(dirToCompress, filepath.Join(tmpDir, "my_archive.tar.gz"))
	if err == nil {
		t.Errorf("Should say that %s is empty", dirToCompress)
	}
}

func Test_CompressAndExtractWithMultipleFiles(t *testing.T) {
	tmpDir, dirToCompress := createTestData()
	defer os.RemoveAll(tmpDir)

	createFiles(dirToCompress, "file1.txt", "file2.txt", "file3.txt")

	structureBefore := directoryStructureString(dirToCompress)

	err := Compress(dirToCompress, filepath.Join(tmpDir, "my_archive.tar.gz"))
	if err != nil {
		t.Errorf("Compress error: %s", err)
	}

	err = Extract(filepath.Join(tmpDir, "my_archive.tar.gz"), filepath.Join(tmpDir, "extracted"))
	if err != nil {
		t.Errorf("Extract error: %s", err)
	}

	structureAfter := directoryStructureString(filepath.Join(tmpDir, "extracted", "my_folder"))

	if structureAfter != structureBefore {
		t.Errorf("Directory structure before compress and after extract does not match. Before {%s}, After {%s}", structureBefore, structureAfter)
	}
}

func Test_ThatOutputDirIsRemovedIfCompressFails(t *testing.T) {
	tmpDir, dirToCompress := createTestData()
	defer os.RemoveAll(tmpDir)

	os.RemoveAll(filepath.Join(dirToCompress, "my_sub_folder"))

	baseDir := "dir_to_be_removed"
	for _, dir := range []string{"", "sub1", "sub1/sub2"} {
		dstDir := filepath.Join(tmpDir, baseDir, dir)
		err := Compress(dirToCompress, filepath.Join(dstDir, "my_archive.tar.gz"))
		if err == nil {
			t.Errorf("Should say that %s is empty", dirToCompress)
		}

		d := filepath.Join(tmpDir, baseDir)
		exist, err := exists(d)
		if err != nil {
			panic(err)
		}
		if exist {
			t.Errorf("%s should be removed", d)
		}
		os.RemoveAll(dstDir)
	}
}

func Test_ThatOutputDirIsRemovedIfExtractFails(t *testing.T) {
	tmpDir, _ := createTestData()
	defer os.RemoveAll(tmpDir)

	err := Extract(filepath.Join(tmpDir, "my_archive.tar.gz"), filepath.Join(tmpDir, "extracted"))
	if err == nil {
		t.Errorf("Should say that %s doesn't exist", filepath.Join(tmpDir, "my_archive.tar.gz"))
	}

	exist, err := exists(filepath.Join(tmpDir, "extracted"))
	if err != nil {
		panic(err)
	}
	if exist {
		t.Errorf("%s should be removed", filepath.Join(tmpDir, "extracted"))
	}
}

func Test_CompabilityWithTar(t *testing.T) {
	_, err := exec.LookPath("tar")
	if err == nil {
		tmpDir, dirToCompress := createTestData()
		defer os.RemoveAll(tmpDir)

		structureBefore := directoryStructureString(dirToCompress)

		err := Compress(dirToCompress, filepath.Join(tmpDir, "my_archive.tar.gz"))
		if err != nil {
			t.Errorf("Compress error: %s", err)
		}

		os.MkdirAll(filepath.Join(tmpDir, "extracted"), 0755)
		cmd := exec.Command("tar", "xfz", filepath.Join(tmpDir, "my_archive.tar.gz"), "-C", filepath.Join(tmpDir, "extracted"))
		err = cmd.Run()
		if err != nil {
			fmt.Println("Run error")
			panic(err)
		}

		structureAfter := directoryStructureString(filepath.Join(tmpDir, "extracted", "my_folder"))

		if structureAfter != structureBefore {
			t.Errorf("Directory structure before compress and after extract does not match. Before {%s}, After {%s}", structureBefore, structureAfter)
		}
	} else {
		t.Skip("Skipping test because tar command was not found.")
	}
}

func createTestData() (string, string) {
	tmpDir, err := ioutil.TempDir("", "targz-test")
	if err != nil {
		fmt.Println("TempDir error")
		panic(err)
	}

	directory := filepath.Join(tmpDir, "my_folder")
	subDirectory := filepath.Join(directory, "my_sub_folder")
	err = os.MkdirAll(subDirectory, 0755)
	if err != nil {
		fmt.Println("MkdirAll error")
		panic(err)
	}

	_, err = os.Create(filepath.Join(subDirectory, "my_file.txt"))
	if err != nil {
		fmt.Println("Create file error")
		panic(err)
	}

	return tmpDir, directory
}

func createFiles(dir string, names ...string) {
	for _, name := range names {
		_, err := os.Create(filepath.Join(dir, name))
		if err != nil {
			fmt.Println("Create file error")
			panic(err)
		}
	}
}

func directoryStructureString(directory string) string {
	structure := ""

	file, err := os.Open(directory)
	if err != nil {
		fmt.Println("Open file error")
		panic(err)
	}
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Stat file error")
		panic(err)
	}

	if fileInfo.IsDir() {
		structure += "-" + filepath.Base(file.Name())

		files, err := ioutil.ReadDir(file.Name())
		if err != nil {
			fmt.Println("ReadDir error")
			panic(err)
		}
		for _, f := range files {
			structure += directoryStructureString(filepath.Join(directory, f.Name()))
		}
	} else {
		structure += "*" + filepath.Base(file.Name())
	}

	return structure
}
