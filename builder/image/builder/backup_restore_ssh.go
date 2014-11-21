package builder

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/coreos/go-etcd/etcd"
	"github.com/crowdmob/goamz/s3"
	"github.com/walle/targz"

	"github.com/deis/deis/boot/commons"
	"github.com/deis/deis/boot/logger"
)

const (
	sshTarFile   string = "ssh.tar.gz"
	bucketName   string = "builder"
	sshDirectory string = "/etc/ssh"
	outputFile   string = "/tmp/ssh.tar.gz"
)

func CheckSSHKeysInStore(client *etcd.Client) {
	logger.Log.Debug("obtaining information of ceph data store")
	accessKey := commons.GetEtcd(client, "/deis/store/gateway/accessKey")
	secretKey := commons.GetEtcd(client, "/deis/store/gateway/secretKey")

	// to access the store we use one of the routers
	storeHosts := commons.GetListEtcd(client, "/deis/router/hosts")
	if storeHosts == nil {
		logger.Log.Debug("there is no hosts (routers) in /deis/router/hosts. check the deis-router/s status")
		return
	}

	// use a fake name
	storeHost := "deis-store.local.deisapp.com"
	logger.Log.Debug("connecting to the ceph store %s (fake name)", storeHost)

	// TODO: change this when deis use some internal dns (skydns).
	// HACK: copy current /etc/hosts, add a host (deis-store) and after the
	// interaction restore the original file. This only is use at boot time
	// to be able to interact wirh the deis-store gateway.
	hostsFile := "/etc/hosts"
	backFile := hostsFile + ".back"

	if _, err := commons.CopyFile(backFile, hostsFile); err != nil {
		logger.Log.Debugf("%v", err)
		return
	}

	file, err := os.OpenFile(hostsFile, os.O_APPEND|os.O_WRONLY, 0600)
	defer file.Close()

	ip := strings.Split(storeHosts[0].Value, ":")
	if _, err = file.WriteString(ip[0] + " " + storeHost); err != nil {
		logger.Log.Debugf("%v", err)
		return
	}
	// end hack

	s3Connection := commons.ConnectS3Store(accessKey, secretKey, "http://"+storeHost)
	builderBucket := s3Connection.Bucket(bucketName)
	builderBucket.PutBucket(s3.BucketOwnerFull)

	logger.Log.Info("checking if there is a backup of the ssh files in the store")
	result, err := builderBucket.Exists(sshTarFile)

	if err != nil {
		logger.Log.Info("an error occurred trying to access the data store. Using local configuration...")
		logger.Log.Debugf("%v", err)
		return
	}

	if !result {
		logger.Log.Info("there is no ssh builder keys in the ceph data store. Uploading local files...")
		err := targz.Compress(sshDirectory, outputFile)
		if err != nil {
			logger.Log.Info("is not possible to upload the local keys")
			logger.Log.Debugf("%v", err)
			return
		}

		data, _ := ioutil.ReadFile(outputFile)

		err = builderBucket.Put(
			sshTarFile,
			data,
			"application/x-compressed",
			s3.BucketOwnerFull, s3.Options{},
		)

		if err != nil {
			logger.Log.Info("an error occurred uploading the ssh keys. Skipping...")
			logger.Log.Debugf("%v", err)
			return
		}
	} else {
		logger.Log.Info("restoring ssh keys from the data store...")
		data, err := builderBucket.Get(sshTarFile)
		if err != nil {
			logger.Log.Info("an error occurred downloading ssh keys from the store. Skipping...")
			logger.Log.Debugf("%v", err)
			return
		}

		if err := ioutil.WriteFile(outputFile, data, 0644); err != nil {
			logger.Log.Info("an error occurred decompresing ssh keys from the store. Skipping...")
			logger.Log.Debugf("%v", err)
			return
		}

		if err := targz.Extract(outputFile, "/etc"); err != nil {
			logger.Log.Info("an error occurred decompressing ssh keys from the store. Skipping...")
			logger.Log.Debugf("%v", err)
			return
		}
	}

	// restore the original /etc/hosts
	if _, err := commons.CopyFile(hostsFile, backFile); err != nil {
		logger.Log.Debugf("an error occurred restoring the original /etc/hosts file: %v", err)
		return
	}
	// remove temporal hosts file
	os.Remove(backFile)
}
