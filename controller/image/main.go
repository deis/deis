package main

import (
	"archive/tar"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/deis/deis/controller/commons"
	"github.com/deis/deis/controller/logger"

	"github.com/docker/docker/pkg/parsers"
	"github.com/fsouza/go-dockerclient"
	"github.com/moraes/config"
)

const (
	workDirRegex string = `\(nop\) WORKDIR (.*)`
)

func main() {
	image := flag.String("image", "", "Application name and version <app>:<version>")

	flag.Parse()

	if flag.NFlag() < 1 || *image == "" {
		returnEmptyResponse()
	}

	logger.Log.Debugf("Checking app %s", *image)

	dockerClient, _ := commons.NewDockerClient()

	logger.Log.Debugf("Image name [%s]", *image)

	imageInfo, err := dockerClient.InspectImage(*image)
	if err != nil {
		logger.Log.Debugf("Error inspecting local image: %v", err)
		logger.Log.Debugf("Pulling image [%s] from registry", *image)

		repo, tag := parsers.ParseRepositoryTag(*image)
		if tag == "" {
			tag = "latest"
		}

		opts := docker.PullImageOptions{Repository: repo, Tag: tag}
		if err = dockerClient.PullImage(opts, docker.AuthConfiguration{}); err != nil {
			logger.Log.Debugf("%v", err)
			logger.Log.Debugf("An error occured pulling the image %s", *image)
			returnEmptyResponse()
		}

		imageInfo, _ = dockerClient.InspectImage(*image)
	}

	procfileContent := getProcfileContent(dockerClient, imageInfo.ID)
	if _, err := config.ParseYaml(procfileContent); err != nil {
		logger.Log.Debugf("%v", err)
		logger.Log.Debug("the procfile does not contains a valid yaml structure")
		returnEmptyResponse()
	}

	commons.ParseProcfile2Structure(procfileContent)
	os.Exit(0)
}

func returnEmptyResponse() {
	// Returna a default value, never an error.
	fmt.Println("{ \"cmd\": 1 }")
	os.Exit(0)
}

func getProcfileContent(client *docker.Client, imageID string) string {
	config := docker.Config{Image: imageID}
	container, err := client.CreateContainer(docker.CreateContainerOptions{
		Config: &config,
	})

	if err != nil {
		logger.Log.Debugf("Error reading Procfile: %v", err)
		returnEmptyResponse()
	}

	defer func() {
		if container != nil {
			client.KillContainer(docker.KillContainerOptions{
				ID: container.ID,
			})
		}
	}()


	var buf bytes.Buffer
	workdir := getProcfileLocation(client, imageID)
	procfile := workdir+"/Procfile"	

	err = client.CopyFromContainer(docker.CopyFromContainerOptions{
		Container:    container.ID,
		Resource:     procfile,
		OutputStream: &buf,
	})

	if err != nil {
		logger.Log.Debugf("Error while copying from %s: %s. Returning None\n", imageID, err)
		returnEmptyResponse()
	}

	content := new(bytes.Buffer)
	r := bytes.NewReader(buf.Bytes())
	tr := tar.NewReader(r)
	tr.Next()
	if err != nil && err != io.EOF {
		logger.Log.Debugf("%v", err)
		returnEmptyResponse()
	}
	if _, err := io.Copy(content, tr); err != nil {
		logger.Log.Debugf("%v", err)
		returnEmptyResponse()
	}

	procfileContent := content.String()
	logger.Log.Debugf("Procfile content \n\n%s\n", procfileContent)
	return procfileContent
}

func getProcfileLocation(client *docker.Client, imageID string) string{
	imageHistory, err := client.ImageHistory(imageID)
	if err != nil {
		return ""
	}

	r := regexp.MustCompile(workDirRegex)
	for _,layer := range imageHistory {
		logger.Log.Debugf("Layer CreatedBy: %s\n", layer.CreatedBy)
		match := r.FindStringSubmatch(layer.CreatedBy)
		if match == nil {
			continue
		}

		logger.Log.Debugf("Workdir: %s\n", match[1])
		return match[1]
	}

	return ""
}