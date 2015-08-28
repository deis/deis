package cmd

import (
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/deis/deis/client/controller/api"
)

func TestGetKey(t *testing.T) {
	t.Parallel()

	file, err := ioutil.TempFile("", "deis-key")

	if err != nil {
		t.Fatal(err)
	}

	toWrite := []byte("ssh-rsa abc test@example.com")

	expected := api.KeyCreateRequest{
		ID:     "test@example.com",
		Public: string(toWrite),
		Name:   file.Name(),
	}

	if _, err = file.Write(toWrite); err != nil {
		t.Fatal(err)
	}

	key, err := getKey(file.Name())

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, key) {
		t.Errorf("Expected %v, Got %v", expected, key)
	}
}

func TestGetKeyNoComment(t *testing.T) {
	t.Parallel()

	file, err := ioutil.TempFile("", "deis-key")

	if err != nil {
		t.Fatal(err)
	}

	toWrite := []byte("ssh-rsa abc")

	expected := api.KeyCreateRequest{
		ID:     path.Base(file.Name()),
		Public: string(toWrite),
		Name:   file.Name(),
	}

	if _, err = file.Write(toWrite); err != nil {
		t.Fatal(err)
	}

	key, err := getKey(file.Name())

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, key) {
		t.Errorf("Expected %v, Got %v", expected, key)
	}
}

func TestListKeys(t *testing.T) {
	name, err := ioutil.TempDir("", "deis-key")

	if err != nil {
		t.Fatal(err)
	}

	os.Setenv("HOME", name)

	folder := path.Join(name, ".ssh")

	if err = os.Mkdir(folder, 0755); err != nil {
		t.Fatal(err)
	}

	toWrite := []byte("ssh-rsa abc test@example.com")
	fileNames := []string{"test1.pub", "test2.pub"}

	expected := []api.KeyCreateRequest{
		api.KeyCreateRequest{
			ID:     "test@example.com",
			Public: string(toWrite),
			Name:   path.Join(folder, fileNames[0]),
		},
		api.KeyCreateRequest{
			ID:     "test@example.com",
			Public: string(toWrite),
			Name:   path.Join(folder, fileNames[1]),
		},
	}

	for _, file := range fileNames {
		if err = ioutil.WriteFile(path.Join(folder, file), toWrite, 0775); err != nil {
			t.Fatal(err)
		}
	}

	keys, err := listKeys()

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, keys) {
		t.Errorf("Expected %v, Got %v", expected, keys)
	}
}
