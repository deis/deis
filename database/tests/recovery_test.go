package tests

import (
	"database/sql"
	"fmt"
	"github.com/deis/deis/tests/dockercli"
	"github.com/deis/deis/tests/mock"
	"github.com/deis/deis/tests/utils"
	"github.com/lib/pq"
	"testing"
	"time"
)

func OpenDeisDatabase(t *testing.T, host string, port string) *sql.DB {
	db, err := sql.Open("postgres", "postgres://deis:changeme123@"+host+":"+port+"/deis?sslmode=disable&connect_timeout=4")
	if err != nil {
		t.Fatal(err)
	}
	WaitForDatabase(t, db)
	return db
}

func WaitForDatabase(t *testing.T, db *sql.DB) {
	fmt.Println("--- Waiting for pg to be ready")
	for {
		err := db.Ping()
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == "cannot_connect_now" {
				fmt.Println(err.Message)
				time.Sleep(1000 * time.Millisecond)
				continue
			}
			t.Fatal(err)
		}
		fmt.Println("Ready")
		break
	}
}

func TryTableSelect(t *testing.T, db *sql.DB, tableName string, expectFailure bool) {
	_, err := db.Query("select * from " + tableName)

	if expectFailure {
		if err == nil {
			t.Fatal("The table should not exist")
		}
	} else {
		if err != nil {
			t.Fatal(err)
		}
	}
}

func execSql(t *testing.T, db *sql.DB, q string) {
	_, err := db.Query(q)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDatabaseRecovery(t *testing.T) {
	var err error
	tag, etcdPort := utils.BuildTag(), utils.RandomPort()
	cli, stdout, _ := dockercli.NewClient()
	imageName := utils.ImagePrefix() + "database" + ":" + tag

	// start etcd container
	etcdName := "deis-etcd-" + tag
	dockercli.RunTestEtcd(t, etcdName, etcdPort)
	defer cli.CmdRm("-f", etcdName)

	// run mock ceph containers
	cephName := "deis-ceph-" + tag
	mock.RunMockCeph(t, cephName, cli, etcdPort)
	defer cli.CmdRm("-f", cephName)

	// create volumes
	databaseVolumeA := "deis-database-data-a-" + tag
	databaseVolumeB := "deis-database-data-b-" + tag
	defer cli.CmdRm("-f", databaseVolumeA)
	defer cli.CmdRm("-f", databaseVolumeB)
	go func() {
		fmt.Printf("--- Creating Volume A\n")
		_ = cli.CmdRm("-f", "-v", databaseVolumeA)
		dockercli.CreateVolume(t, cli, databaseVolumeA, "/var/lib/postgresql")

		fmt.Printf("--- Creating Volume B\n")

		_ = cli.CmdRm("-f", databaseVolumeB)
		dockercli.CreateVolume(t, cli, databaseVolumeB, "/var/lib/postgresql")
	}()
	dockercli.WaitForLine(t, stdout, databaseVolumeB, true)

	// setup database container start/stop routines
	host, port := utils.HostAddress(), utils.RandomPort()
	fmt.Printf("--- Run deis/database:%s at %s:%s\n", tag, host, port)
	name := "deis-database-" + tag
	defer cli.CmdRm("-f", name)
	startDatabase := func(volumeName string) {
		_ = cli.CmdRm("-f", name)
		err = dockercli.RunContainer(cli,
			"--name", name,
			"--volumes-from", volumeName,
			"--rm",
			"-p", port+":5432",
			"-e", "EXTERNAL_PORT="+port,
			"-e", "HOST="+host,
			"-e", "ETCD_PORT="+etcdPort,
			"-e", "ETCD_TTL=2",
			"-e", "BACKUP_FREQUENCY=1s",
			"-e", "BACKUPS_TO_RETAIN=100",
			imageName)
	}

	stopDatabase := func() {
		fmt.Println("--- Stopping data-database... ")
		if err = stdout.Close(); err != nil {
			t.Fatal("Failed to closeStdout")
		}
		_ = cli.CmdStop(name)
		fmt.Println("Done")
	}

	//ACTION

	//STEP 1: start db with volume A and wait for init to complete
	fmt.Println("--- Starting database with Volume A... ")
	go startDatabase(databaseVolumeA)
	dockercli.WaitForLine(t, stdout, "database: postgres is running...", true)
	fmt.Println("Done")

	db := OpenDeisDatabase(t, host, port)
	TryTableSelect(t, db, "api_foo", true)

	stopDatabase()

	//STEP 2a: start db with volume B, wait for init and create the table
	cli, stdout, _ = dockercli.NewClient()
	fmt.Printf("--- Starting database with Volume B... ")
	go startDatabase(databaseVolumeB)
	dockercli.WaitForLine(t, stdout, "database: postgres is running...", true)
	fmt.Println("Done")

	db = OpenDeisDatabase(t, host, port)
	TryTableSelect(t, db, "api_foo", true)

	fmt.Println("--- Creating the table")
	execSql(t, db, "create table api_foo(t text)")

	//STEP 2b: make sure we observed full backup cycle after forced checkpoint
	fmt.Println("--- Waiting for the change to be backed up... ")
	dockercli.WaitForLine(t, stdout, "database: performing a backup...", true)
	dockercli.WaitForLine(t, stdout, "database: backup has been completed.", true)
	fmt.Println("Done")

	stopDatabase()

	//STEP 3: start db with volume A again and assert table existence
	cli, stdout, _ = dockercli.NewClient()
	fmt.Printf("--- Starting database with Volume A again... ")
	go startDatabase(databaseVolumeA)
	dockercli.WaitForLine(t, stdout, "database: postgres is running...", true)
	fmt.Println("Done")

	db = OpenDeisDatabase(t, host, port)
	TryTableSelect(t, db, "api_foo", false)

}
