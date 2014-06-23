package mockserviceutils

import (
	"github.com/deis/deis/tests/dockercliutils"
	"github.com/deis/deis/tests/etcdutils"
	"github.com/deis/deis/tests/utils"
	"testing"
	"time"
	"fmt"
)

func RunMockDatabase(t *testing.T, uid string,port string) {
	cli, stdout, stdoutPipe := dockercliutils.GetNewClient()
	done := make(chan bool, 1)
	var imageId string
	var imageTag string
	IPAddress :=  utils.GetHostIpAddress()
	fmt.Println("starting Mock Database")
	done <-true
	go func() {
		<-done
		dockercliutils.PullImage(t, cli, "paintedfox/postgresql")
		imageId = dockercliutils.GetImageId(t, "paintedfox/postgresql")
		imageTag = "deis/database:" + uid
		cli.CmdTag(imageId, imageTag)
		dockercliutils.RunContainer(t, cli, "--name", "deis-database-"+uid, "-p", "5432:5432", "-e", "PUBLISH=5432", "-e", "HOST="+IPAddress, "-e", "USER=deis", "-e", "DB=deis", "-e", "PASS=deis", "deis/database:"+uid)
	}()
	time.Sleep(1000 * time.Millisecond)
	dockercliutils.PrintToStdout(t, stdout, stdoutPipe, "Starting")
	// docker run --name="deis-database"  -p -e PUBLISH=5432 -e HOST=172.17.8.100 -e USER="super" -e DB="deis" -e PASS="jaffa"  deis/database
	setkeys := []string{"/deis/database/user",
		"/deis/database/password",
		"/deis/database/name"}
	setdir := []string{}
	dbhandler := etcdutils.InitetcdValues(setdir, setkeys, port)
	etcdutils.Publishvalues(t, dbhandler)
	IPAddress = dockercliutils.GetInspectData(t, "{{ .NetworkSettings.IPAddress }}", "deis-database-"+uid)
	etcdutils.SetEtcdValues(t, []string{"/deis/database/host", "/deis/database/port", "/deis/database/engine"}, []string{IPAddress, "5432", "postgresql_psycopg2"}, dbhandler.C)
}
