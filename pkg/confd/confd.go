package confd

import (
	//"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	// "regexp"
	"syscall"
	"time"

	"github.com/pmylund/go-cache"

	logger "github.com/deis/deis/pkg/log"
	oswrapper "github.com/deis/deis/pkg/os"
)

var (
	log                = logger.New()
	templateErrorRegex = "(\\d{4})-(\\d{2})-(\\d{2})T(\\d{2}):(\\d{2}):\\d{2}Z.*ERROR template:"
	confdCacheError    = cache.New(60*time.Minute, 60*time.Second)
)

// WaitForInitialConf wait until the compilation of the templates is correct
func WaitForInitialConf(signalChan chan os.Signal, etcd string, timeout time.Duration) {
	log.Info("waiting for confd to write initial templates...")
	for {
		cmdAsString := fmt.Sprintf("confd -onetime -node %s -confdir /app", etcd)
		cmd, args := oswrapper.BuildCommandFromString(cmdAsString)
		err := oswrapper.RunCommand(signalChan, cmd, args, false)
		if err == nil {
			break
		}

		time.Sleep(timeout)
	}
}

// Launch launch confd as a daemon process.
func Launch(signalChan chan os.Signal, etcd string) {
	cmdAsString := fmt.Sprintf("confd -node %s -confdir /app --interval 5 --log-level error", etcd)
	cmd, args := oswrapper.BuildCommandFromString(cmdAsString)
	go runConfdDaemon(signalChan, cmd, args)
}

func runConfdDaemon(signalChan chan os.Signal, command string, args []string) {
	// testRegex := regexp.MustCompile(templateErrorRegex)

	cmd := exec.Command(command, args...)

	stdout, err := cmd.StdoutPipe()
	checkError(signalChan, err)
	stderr, err := cmd.StderrPipe()
	checkError(signalChan, err)

	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)

	// we check if there's more than 5 errors per minute
	// in that case we need to exist and restart the component
	// this is to avoid the endless wait for some keys in etcd
	// go func() {
	// 	scanner := bufio.NewScanner(stderr)
	// 	for scanner.Scan() {
	// 		match := testRegex.FindStringSubmatch(scanner.Text())
	// 		if match != nil {
	// 			ts := match[1] + match[2] + match[3] + match[4] + match[5]
	// 			if _, found := confdCacheError.Get(ts); found {
	// 				confdCacheError.IncrementInt(ts, 1)
	// 			} else {
	// 				confdCacheError.Set(ts, 1, cache.DefaultExpiration)
	// 			}
	// 			errorCount, _ := confdCacheError.Get(ts)
	// 			log.Errorf("confd template error (%v in the last minute)", errorCount)
	// 			if errorCount.(int) > 4 {
	// 				log.Error("too many confd errors in the last minute. restarting deis component")
	// 				signalChan <- syscall.SIGKILL
	// 			}
	// 		}
	// 	}
	// }()

	err = cmd.Start()
	if err != nil {
		log.Errorf("an error ocurred executing command: [%s params %v], %v", command, args, err)
		signalChan <- syscall.SIGKILL
	}

	err = cmd.Wait()
	log.Errorf("command finished with error: %v", err)
	signalChan <- syscall.SIGKILL
}

func checkError(signalChan chan os.Signal, err error) {
	if err != nil {
		log.Errorf("%v", err)
		signalChan <- syscall.SIGKILL
	}
}
