package main

import (
	"github.com/robfig/cron"

	"github.com/deis/deis/database/bindata"

	"github.com/deis/deis/pkg/boot"
	"github.com/deis/deis/pkg/etcd"
	"github.com/deis/deis/pkg/os"

	Log "github.com/deis/deis/pkg/log"
)

var log = Log.New()

func main() {
	externalPort := os.Getopt("EXTERNAL_PORT", "5432")
	etcdPath := os.Getopt("ETCD_PATH", "/deis/database")
	adminUser := os.Getopt("PG_ADMIN_USER", "postgres")
	adminPass := os.Getopt("PG_ADMIN_PASS", "changeme123")
	user := os.Getopt("PG_USER_NAME", "deis")
	password := os.Getopt("PG_USER_PASS", "changeme123")
	name := os.Getopt("PG_USER_DB", "deis")
	bucketName := os.Getopt("BUCKET_NAME", "db_wal")
	backupsToRetain := os.Getopt("BACKUPS_TO_RETAIN", "5")
	backupFrequency := os.Getopt("BACKUP_FREQUENCY", "3h")
	pgConfig := os.Getopt("PG_CONFIG", "/etc/postgresql/9.3/main/postgresql.conf")
	listenAddress := os.Getopt("PG_LISTEN", "*")

	bootProcess := boot.New(etcdPath, externalPort)

	log.Debug("creating required defaults in etcd...")
	etcd.Mkdir(bootProcess.Etcd, etcdPath)
	etcd.SetDefault(bootProcess.Etcd, etcdPath+"/engine", "postgresql_psycopg2")
	etcd.SetDefault(bootProcess.Etcd, etcdPath+"/adminUser", adminUser)
	etcd.SetDefault(bootProcess.Etcd, etcdPath+"/adminPass", adminPass)
	etcd.SetDefault(bootProcess.Etcd, etcdPath+"/user", user)
	etcd.SetDefault(bootProcess.Etcd, etcdPath+"/password", password)
	etcd.SetDefault(bootProcess.Etcd, etcdPath+"/name", name)
	etcd.SetDefault(bootProcess.Etcd, etcdPath+"/bucketName", bucketName)

	bootProcess.StartConfd()

	bootProcess.RunScript("bash/postgres-init.bash", nil, bindata.Asset)

	postgresCommand := "sudo -i -u postgres /usr/lib/postgresql/9.3/bin/postgres" +
		" -c config-file=" + pgConfig +
		" -c listen-addresses=" + listenAddress
	bootProcess.RunProcessAsDaemon(os.BuildCommandFromString(postgresCommand))
	bootProcess.WaitForLocalConnection("5432")
	bootProcess.Publish()

	params := make(map[string]string)
	params["BUCKET_NAME"] = bucketName
	bootProcess.RunScript("bash/postgres.bash", params, bindata.Asset)

	// schedule periodic backups using wal-e
	scheduleBackup := cron.New()
	scheduleBackup.AddFunc("@every "+backupFrequency,
		func() {
			log.Debug("creating database backup with wal-e...")
			params := make(map[string]string)
			params["BACKUPS_TO_RETAIN"] = backupsToRetain
			bootProcess.RunScript("bash/backup.bash", params, bindata.Asset)
		})
	scheduleBackup.Start()

	log.Info("database: postgres is running...")
	bootProcess.Wait()
}
