package main

import (
	"strconv"

	"github.com/deis/deis/database/bindata"

	"github.com/deis/deis/pkg/boot"
	"github.com/deis/deis/pkg/etcd"
	Log "github.com/deis/deis/pkg/log"
	"github.com/deis/deis/pkg/os"
	"github.com/deis/deis/pkg/types"
)

const (
	servicePort = 5432
)

var (
	etcdPath     = os.Getopt("ETCD_PATH", "/deis/database")
	externalPort = os.Getopt("EXTERNAL_PORT", strconv.Itoa(servicePort))
	bucketName   = os.Getopt("BUCKET_NAME", "db_wal")
	log          = Log.New()
)

func init() {
	boot.RegisterComponent(new(DatabaseBoot), "deis-component")
}

func main() {
	boot.Start(etcdPath, externalPort, false)
}

type DatabaseBoot struct{}

func (dbb *DatabaseBoot) MkdirsEtcd() []string {
	return []string{}
}

func (dbb *DatabaseBoot) EtcdDefaults() map[string]string {
	adminUser := os.Getopt("PG_ADMIN_USER", "postgres")
	adminPass := os.Getopt("PG_ADMIN_PASS", "changeme123")
	user := os.Getopt("PG_USER_NAME", "deis")
	password := os.Getopt("PG_USER_PASS", "changeme123")
	name := os.Getopt("PG_USER_DB", "deis")

	keys := make(map[string]string)
	keys[etcdPath+"/engine"] = "postgresql_psycopg2"
	keys[etcdPath+"/adminUser"] = adminUser
	keys[etcdPath+"/adminPass"] = adminPass
	keys[etcdPath+"/user"] = user
	keys[etcdPath+"/password"] = password
	keys[etcdPath+"/name"] = name
	keys[etcdPath+"/bucketName"] = bucketName
	return keys
}

func (dbb *DatabaseBoot) PreBootScripts(currentBoot *types.CurrentBoot) []*types.Script {
	params := make(map[string]string)
	params["BUCKET_NAME"] = bucketName
	params["ETCD_PATH"] = etcdPath
	params["ETCD"] = currentBoot.Host.String() + ":" + currentBoot.EtcdPort
	return []*types.Script{
		&types.Script{Name: "bash/postgres-init.bash", Content: bindata.Asset},
		&types.Script{Name: "bash/postgres.bash", Params: params, Content: bindata.Asset},
	}
}

func (dbb *DatabaseBoot) PreBoot(currentBoot *types.CurrentBoot) {
	log.Info("database: starting...")
	os.RunScript("bash/create-dummy-scripts.bash", nil, bindata.Asset)
}

func (dbb *DatabaseBoot) BootDaemons(currentBoot *types.CurrentBoot) []*types.ServiceDaemon {
	pgConfig := os.Getopt("PG_CONFIG", "/etc/postgresql/9.3/main/postgresql.conf")
	listenAddress := os.Getopt("PG_LISTEN", "*")
	postgresCommand := "sudo -i -u postgres /usr/lib/postgresql/9.3/bin/postgres" +
		" -c config-file=" + pgConfig +
		" -c listen-addresses=" + listenAddress
	cmd, args := os.BuildCommandFromString(postgresCommand)
	return []*types.ServiceDaemon{&types.ServiceDaemon{Command: cmd, Args: args}}
}

func (dbb *DatabaseBoot) WaitForPorts() []int {
	return []int{servicePort}
}

func (dbb *DatabaseBoot) PostBootScripts(currentBoot *types.CurrentBoot) []*types.Script {
	initialBackup := etcd.Get(currentBoot.EtcdClient, etcdPath+"/initialBackup")
	if initialBackup == "1" {
		// do not delete backups from this run
		backupsToRetain := os.Getopt("BACKUPS_TO_RETAIN", "100000")
		params := make(map[string]string)
		params["BACKUPS_TO_RETAIN"] = backupsToRetain
		return []*types.Script{
			&types.Script{Name: "bash/backup.bash", Params: params, Content: bindata.Asset},
		}
	} else {
		return []*types.Script{}
	}
}

func (dbb *DatabaseBoot) PostBoot(currentBoot *types.CurrentBoot) {
	log.Info("database: postgres is running...")
}

func (dbb *DatabaseBoot) ScheduleTasks(currentBoot *types.CurrentBoot) []*types.Cron {
	backupsToRetain := os.Getopt("BACKUPS_TO_RETAIN", "5")
	backupFrequency := os.Getopt("BACKUP_FREQUENCY", "3h")
	params := make(map[string]string)
	params["BACKUPS_TO_RETAIN"] = backupsToRetain

	return []*types.Cron{
		&types.Cron{
			Frequency: "@every " + backupFrequency,
			Code: func() {
				log.Debug("creating database backup with wal-e...")
				os.RunScript("bash/backup.bash", params, bindata.Asset)
			},
		},
	}
}

func (dbb *DatabaseBoot) UseConfd() bool {
	return true
}
