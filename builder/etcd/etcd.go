/*Package etcd is a library for performing common Etcd tasks.
 */
package etcd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Masterminds/cookoo"
	"github.com/Masterminds/cookoo/log"
	"github.com/Masterminds/cookoo/safely"
	"github.com/coreos/go-etcd/etcd"
)

var (
	retryCycles = 2
	retrySleep  = 200 * time.Millisecond
)

// Getter describes the Get behavior of an Etcd client.
//
// Usually you will want to use go-etcd/etcd.Client to satisfy this.
//
// We use an interface because it is more testable.
type Getter interface {
	Get(string, bool, bool) (*etcd.Response, error)
}

// DirCreator describes etcd's CreateDir behavior.
//
// Usually you will want to use go-etcd/etcd.Client to satisfy this.
type DirCreator interface {
	CreateDir(string, uint64) (*etcd.Response, error)
}

// Watcher watches an etcd entry.
type Watcher interface {
	Watch(string, uint64, bool, chan *etcd.Response, chan bool) (*etcd.Response, error)
}

// Setter sets a value in Etcd.
type Setter interface {
	Set(string, string, uint64) (*etcd.Response, error)
}

// GetterSetter performs get and set operations.
type GetterSetter interface {
	Getter
	Setter
}

// CreateClient creates a new Etcd client and prepares it for work.
//
// Params:
// 	- url (string): A server to connect to.
// 	- retries (int): Number of times to retry a connection to the server
// 	- retrySleep (time.Duration): How long to sleep between retries
//
// Returns:
// 	This puts an *etcd.Client into the context.
func CreateClient(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	url := p.Get("url", "http://localhost:4001").(string)

	// Backed this out because it's unnecessary so far.
	//hosts := p.Get("urls", []string{"http://localhost:4001"}).([]string)
	hosts := []string{url}
	retryCycles = p.Get("retries", retryCycles).(int)
	retrySleep = p.Get("retrySleep", retrySleep).(time.Duration)

	// Support `host:port` format, too.
	for i, host := range hosts {
		if !strings.Contains(host, "://") {
			hosts[i] = "http://" + host
		}
	}

	client := etcd.NewClient(hosts)
	client.CheckRetry = checkRetry

	return client, nil
}

// Get performs an etcd Get operation.
//
// Params:
// 	- client (EtcdGetter): Etcd client
// 	- path (string): The path/key to fetch
//
// Returns:
// - This puts an `etcd.Response` into the context, and returns an error
//   if the client could not connect.
func Get(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	cli, ok := p.Has("client")
	if !ok {
		return nil, errors.New("No Etcd client found.")
	}
	client := cli.(Getter)
	path := p.Get("path", "/").(string)

	res, err := client.Get(path, false, false)
	if err != nil {
		return res, err
	}

	if !res.Node.Dir {
		return res, fmt.Errorf("Expected / to be a dir.")
	}
	return res, nil
}

// IsRunning checks to see if etcd is running.
//
// It will test `count` times before giving up.
//
// Params:
// 	- client (EtcdGetter)
// 	- count (int): Number of times to try before giving up.
//
// Returns:
// 	boolean true if etcd is listening.
func IsRunning(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	client := p.Get("client", nil).(Getter)
	count := p.Get("count", 20).(int)
	for i := 0; i < count; i++ {
		_, err := client.Get("/", false, false)
		if err == nil {
			return true, nil
		}
		log.Infof(c, "Waiting for etcd to come online.")
		time.Sleep(250 * time.Millisecond)
	}
	log.Errf(c, "Etcd is not answering after %d attempts.", count)
	return false, &cookoo.FatalError{"Could not connect to Etcd."}
}

// Set sets a value in etcd.
//
// Params:
// 	- key (string): The key
// 	- value (string): The value
// 	- ttl (uint64): Time to live
// 	- client (EtcdGetter): Client, usually an *etcd.Client.
//
// Returns:
// 	- *etcd.Result
func Set(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	key := p.Get("key", "").(string)
	value := p.Get("value", "").(string)
	ttl := p.Get("ttl", uint64(20)).(uint64)
	client := p.Get("client", nil).(Setter)

	res, err := client.Set(key, value, ttl)
	if err != nil {
		log.Infof(c, "Failed to set %s=%s", key, value)
		return res, err
	}

	return res, nil
}

// FindSSHUser finds an SSH user by public key.
//
// Some parts of the system require that we know not only the SSH key, but also
// the name of the user. That information is stored in etcd.
//
// Params:
// 	- client (EtcdGetter)
// 	- fingerprint (string): The fingerprint of the SSH key.
//
// Returns:
// - username (string)
func FindSSHUser(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	client := p.Get("client", nil).(Getter)
	fingerprint := p.Get("fingerprint", nil).(string)

	res, err := client.Get("/deis/builder/users", false, true)
	if err != nil {
		log.Warnf(c, "Error querying etcd: %s", err)
		return "", err
	} else if res.Node == nil || !res.Node.Dir {
		log.Warnf(c, "No users found in etcd.")
		return "", errors.New("Users not found")
	}
	for _, user := range res.Node.Nodes {
		log.Infof(c, "Checking user %s", user.Key)
		for _, keyprint := range user.Nodes {
			if strings.HasSuffix(keyprint.Key, fingerprint) {
				parts := strings.Split(user.Key, "/")
				username := parts[len(parts)-1]
				log.Infof(c, "Found user %s for fingerprint %s", username, fingerprint)
				return username, nil
			}
		}
	}

	return "", fmt.Errorf("User not found for fingerprint %s", fingerprint)
}

// StoreHostKeys stores SSH hostkeys locally.
//
// First it tries to fetch them from etcd. If the keys are not present there,
// it generates new ones and then puts them into etcd.
//
// Params:
// 	- client(EtcdGetterSetter)
// 	- ciphers([]string): A list of ciphers to generate. Defaults are dsa,
// 		ecdsa, ed25519 and rsa.
// 	- basepath (string): Base path in etcd (ETCD_PATH).
// Returns:
//
func StoreHostKeys(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	defaultCiphers := []string{"rsa", "dsa", "ecdsa", "ed25519"}
	client := p.Get("client", nil).(GetterSetter)
	ciphers := p.Get("ciphers", defaultCiphers).([]string)
	basepath := p.Get("basepath", "/deis/builder").(string)

	res, err := client.Get("sshHostKey", false, false)
	if err != nil || res.Node == nil {
		log.Infof(c, "Could not get SSH host key from etcd. Generating new ones.")
		if err := genSSHKeys(c); err != nil {
			log.Err(c, "Failed to generate SSH keys. Aborting.")
			return nil, err
		}
		if err := keysToEtcd(c, client, ciphers, basepath); err != nil {
			return nil, err
		}
	} else if err := keysToLocal(c, client, ciphers, basepath); err != nil {
		log.Infof(c, "Fetching SSH host keys from etcd.")
		return nil, err
	}

	return nil, nil
}

// keysToLocal copies SSH host keys from etcd to the local file system.
//
// This only fails if the main key, sshHostKey cannot be stored or retrieved.
func keysToLocal(c cookoo.Context, client Getter, ciphers []string, etcdPath string) error {
	lpath := "/etc/ssh/ssh_host_%s_key"
	privkey := "%s/sshHost%sKey"
	for _, cipher := range ciphers {
		path := fmt.Sprintf(lpath, cipher)
		key := fmt.Sprintf(privkey, etcdPath, cipher)
		res, err := client.Get(key, false, false)
		if err != nil || res.Node == nil {
			continue
		}

		content := res.Node.Value
		if err := ioutil.WriteFile(path, []byte(content), 0600); err != nil {
			log.Errf(c, "Error writing ssh host key file: %s", err)
		}
	}

	// Now get generic key.
	res, err := client.Get("sshHostKey", false, false)
	if err != nil || res.Node == nil {
		return fmt.Errorf("Failed to get sshHostKey from etcd. %v", err)
	}

	content := res.Node.Value
	if err := ioutil.WriteFile("/etc/ssh/ssh_host_key", []byte(content), 0600); err != nil {
		log.Errf(c, "Error writing ssh host key file: %s", err)
		return err
	}
	return nil
}

// keysToEtcd copies local keys into etcd.
//
// It only fails if it cannot copy ssh_host_key to sshHostKey. All other
// abnormal conditions are logged, but not considered to be failures.
func keysToEtcd(c cookoo.Context, client Setter, ciphers []string, etcdPath string) error {
	lpath := "/etc/ssh/ssh_host_%s_key"
	privkey := "%s/sshHost%sKey"
	for _, cipher := range ciphers {
		path := fmt.Sprintf(lpath, cipher)
		key := fmt.Sprintf(privkey, etcdPath, cipher)
		content, err := ioutil.ReadFile(path)
		if err != nil {
			log.Infof(c, "No key named %s", path)
		} else if _, err := client.Set(key, string(content), 0); err != nil {
			log.Errf(c, "Could not store ssh key in etcd: %s", err)
		}
	}
	// Now we set the generic key:
	if content, err := ioutil.ReadFile("/etc/ssh/ssh_host_key"); err != nil {
		log.Errf(c, "Could not read the ssh_host_key file.")
		return err
	} else if _, err := client.Set("sshHostKey", string(content), 0); err != nil {
		log.Errf(c, "Failed to set sshHostKey in etcd.")
		return err
	}
	return nil
}

// genSshKeys generates the default set of SSH host keys.
func genSSHKeys(c cookoo.Context) error {
	// Generate a new key
	out, err := exec.Command("ssh-keygen", "-A").CombinedOutput()
	if err != nil {
		log.Infof(c, "ssh-keygen: %s", out)
		log.Errf(c, "Failed to generate SSH keys: %s", err)
		return err
	}
	return nil
}

// UpdateHostPort intermittently notifies etcd of the builder's address.
//
// If `port` is specified, this will notify etcd at 10 second intervals that
// the builder is listening at $HOST:$PORT, setting the TTL to 20 seconds.
//
// This will notify etcd as long as the local sshd is running.
//
// Params:
// 	- base (string): The base path to write the data: $base/host and $base/port.
// 	- host (string): The hostname
// 	- port (string): The port
// 	- client (Setter): The client to use to write the data to etcd.
// 	- sshPid (int): The PID for SSHD. If SSHD dies, this stops notifying.
func UpdateHostPort(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	base := p.Get("base", "").(string)
	host := p.Get("host", "").(string)
	port := p.Get("port", "").(string)
	client := p.Get("client", nil).(Setter)
	sshd := p.Get("sshdPid", 0).(int)

	// If no port is specified, we don't do anything.
	if len(port) == 0 {
		log.Infof(c, "No external port provided. Not publishing details.")
		return false, nil
	}

	var ttl uint64 = 20

	if err := setHostPort(client, base, host, port, ttl); err != nil {
		log.Errf(c, "Etcd error setting host/port: %s", err)
		return false, err
	}

	// Update etcd every ten seconds with this builder's host/port.
	safely.GoDo(c, func() {
		ticker := time.NewTicker(10 * time.Second)
		for range ticker.C {
			//log.Infof(c, "Setting SSHD host/port")
			if _, err := os.FindProcess(sshd); err != nil {
				log.Errf(c, "Lost SSHd process: %s", err)
				break
			} else {
				if err := setHostPort(client, base, host, port, ttl); err != nil {
					log.Errf(c, "Etcd error setting host/port: %s", err)
					break
				}
			}
		}
		ticker.Stop()
	})

	return true, nil
}

func setHostPort(client Setter, base, host, port string, ttl uint64) error {
	if _, err := client.Set(base+"/host", host, ttl); err != nil {
		return err
	}
	if _, err := client.Set(base+"/port", port, ttl); err != nil {
		return err
	}
	return nil
}

// MakeDir makes a directory in Etcd.
//
// Params:
// 	- client (EtcdDirCreator): Etcd client
//  - path (string): The name of the directory to create.
// 	- ttl (uint64): Time to live.
// Returns:
// 	*etcd.Response
func MakeDir(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	name := p.Get("path", "").(string)
	ttl := p.Get("ttl", uint64(0)).(uint64)
	cli, ok := p.Has("client")
	if !ok {
		return nil, errors.New("No Etcd client found.")
	}
	client := cli.(DirCreator)

	if len(name) == 0 {
		return false, errors.New("Expected directory name to be more than zero characters.")
	}

	res, err := client.CreateDir(name, ttl)
	if err != nil {
		return res, &cookoo.RecoverableError{err.Error()}
	}

	return res, nil
}

// Watch watches a given path, and executes a git check-repos for each event.
//
// It starts the watcher and then returns. The watcher runs on its own
// goroutine. To stop the watching, send the returned channel a bool.
//
// Params:
// - client (Watcher): An Etcd client.
// - path (string): The path to watch
//
// Returns:
// 	- chan bool: Send this a message to stop the watcher.
func Watch(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	// etcdctl -C $ETCD watch --recursive /deis/services
	path := p.Get("path", "/deis/services").(string)
	cli, ok := p.Has("client")
	if !ok {
		return nil, errors.New("No etcd client found.")
	}
	client := cli.(Watcher)

	// Stupid hack because etcd watch seems to be broken, constantly complaining
	// that the JSON it received is malformed.
	safely.GoDo(c, func() {
		for {
			response, err := client.Watch(path, 0, true, nil, nil)
			if err != nil {
				log.Errf(c, "Etcd Watch failed: %s", err)
				time.Sleep(50 * time.Millisecond)
				continue
			}

			if response.Node == nil {
				log.Infof(c, "Unexpected Etcd message: %v", response)
			}
			git := exec.Command("/home/git/check-repos")
			if out, err := git.CombinedOutput(); err != nil {
				log.Errf(c, "Failed git check-repos: %s", err)
				log.Infof(c, "Output: %s", out)
			}
		}

	})

	return nil, nil

	/* Watch seems to be broken. So we do this stupid watch loop instead.
	receiver := make(chan *etcd.Response)
	stop := make(chan bool)
	// Buffer the channels so that we don't hang waiting for go-etcd to
	// read off the channel.
	stopetcd := make(chan bool, 1)
	stopwatch := make(chan bool, 1)


	// Watch for errors.
	safely.GoDo(c, func() {
		// When a receiver is passed in, no *Response is ever returned. Instead,
		// Watch acts like an error channel, and receiver gets all of the messages.
		_, err := client.Watch(path, 0, true, receiver, stopetcd)
		if err != nil {
			log.Infof(c, "Watcher stopped with error '%s'", err)
			stopwatch <- true
			//close(stopwatch)
		}
	})
	// Watch for events
	safely.GoDo(c, func() {
		for {
			select {
			case msg := <-receiver:
				if msg.Node != nil {
					log.Infof(c, "Received notification %s for %s", msg.Action, msg.Node.Key)
				} else {
					log.Infof(c, "Received unexpected etcd message: %v", msg)
				}
				git := exec.Command("/home/git/check-repos")
				if out, err := git.CombinedOutput(); err != nil {
					log.Errf(c, "Failed git check-repos: %s", err)
					log.Infof(c, "Output: %s", out)
				}
			case <-stopwatch:
				c.Logf("debug", "Received signal to stop watching events.")
				return
			}
		}
	})
	// Fan out stop requests.
	safely.GoDo(c, func() {
		<-stop
		stopwatch <- true
		stopetcd <- true
		close(stopwatch)
		close(stopetcd)
	})

	return stop, nil
	*/
}

// checkRetry overrides etcd.DefaultCheckRetry.
//
// It adds configurable number of retries and configurable timesouts.
func checkRetry(c *etcd.Cluster, numReqs int, last http.Response, err error) error {
	if numReqs > retryCycles*len(c.Machines) {
		return fmt.Errorf("Tried and failed %d cluster connections: %s", retryCycles, err)
	}

	switch last.StatusCode {
	case 0:
		return nil
	case 500:
		time.Sleep(retrySleep)
		return nil
	case 200:
		return nil
	default:
		return fmt.Errorf("Unhandled HTTP Error: %s %d", last.Status, last.StatusCode)
	}
}
