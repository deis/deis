package mockserviceutils

import (
	"github.com/deis/deis/tests/dockercliutils"
	"github.com/deis/deis/tests/etcdutils"
	//"github.com/deis/deis/tests/utils"
	"fmt"
	"testing"
	"time"
)

func RunMockDatabase(t *testing.T, uid string) {
	cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
	done := make(chan bool, 1)
	done1 := make(chan bool, 1)
	done2 := make(chan bool, 1)
	var imageId string
	var imageTag string
	go func() {
		fmt.Println("inside pull postgresql")
		dockercliutils.PullImage(t, cli, "paintedfox/postgresql")
		done <- true
	}()
	go func() {
		<-done
		fmt.Println("inside getting imageId")
		imageId = dockercliutils.GetImageId(t, "paintedfox/postgresql")
		imageTag = "deis/database:" + uid
		cli.CmdTag(imageId, imageTag)
		done1 <- true
	}()
	go func() {
		<-done1
		done2 <- true
		fmt.Println("inside run etcd")
		dockercliutils.RunContainer(t, cli, "--name", "deis-database-"+uid, "-p", "5432:5432", "-e", "PUBLISH=5432", "-e", "HOST=172.17.8.100", "-e", "USER=deis", "-e", "DB=deis", "-e", "PASS=deis", "deis/database:"+uid)

	}()
	time.Sleep(1000 * time.Millisecond)
	dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "Starting")
	// docker run --name="deis-database"  -p -e PUBLISH=5432 -e HOST=172.17.8.100 -e USER="super" -e DB="deis" -e PASS="jaffa"  deis/database
	setkeys := []string{"/deis/database/user",
		"/deis/database/password",
		"/deis/database/name"}
	setdir := []string{}
	dbhandler := etcdutils.InitetcdValues(setdir, setkeys)
	etcdutils.PublishControllervalues(t, dbhandler)
	IPAddress := dockercliutils.GetInspectData(t, "{{ .NetworkSettings.IPAddress }}", "deis-database-"+uid)
	etcdutils.SetEtcdValues(t, []string{"/deis/database/host", "/deis/database/port", "/deis/database/engine"}, []string{IPAddress, "5432", "postgresql_psycopg2"}, dbhandler.C)
}
