package confd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
	"time"

	logger "github.com/deis/deis/mesos/pkg/log"
	oswrapper "github.com/deis/deis/mesos/pkg/os"
)

const (
	confdInterval       = 5                // seconds
	errorTickInterval   = 60 * time.Second // 1 minute
	maxErrorsInInterval = 5                // up to 5 errors per time interval
)

var (
	log                = logger.New()
	templateErrorRegex = "(\\d{4})-(\\d{2})-(\\d{2})T(\\d{2}):(\\d{2}):\\d{2}Z.*ERROR template:"
)

// WaitForInitialConf wait until the compilation of the templates is correct
func WaitForInitialConf(etcd []string, timeout time.Duration) {
	log.Info("waiting for confd to write initial templates...")
	for {
		cmdAsString := fmt.Sprintf("confd -onetime -node %v -confdir /app", strings.Join(etcd, ","))
		log.Debugf("running %s", cmdAsString)
		cmd, args := oswrapper.BuildCommandFromString(cmdAsString)
		err := oswrapper.RunCommand(cmd, args)
		if err == nil {
			break
		}

		time.Sleep(timeout)
	}
}

// Launch launch confd as a daemon process.
func Launch(signalChan chan os.Signal, etcd []string) {
	confdLogLevel := "error"
	if log.Level.String() == "debug" {
		confdLogLevel = "debug"
	}
	cmdAsString := fmt.Sprintf("confd -node %v -confdir /app --interval %v --log-level %v", confdInterval, strings.Join(etcd, ","), confdLogLevel)
	cmd, args := oswrapper.BuildCommandFromString(cmdAsString)
	go runConfdDaemon(signalChan, cmd, args)
}

func runConfdDaemon(signalChan chan os.Signal, command string, args []string) {
	cmd := exec.Command(command, args...)

	stdout, err := cmd.StdoutPipe()
	checkError(signalChan, err)
	// stderr, err := cmd.StderrPipe()
	// checkError(signalChan, err)

	go io.Copy(os.Stdout, stdout)
	// go io.Copy(os.Stderr, stderr)

	go checkNumberOfErrors(stdout, maxErrorsInInterval, errorTickInterval, signalChan)

	err = cmd.Start()
	if err != nil {
		log.Errorf("an error ocurred executing confd: [%s params %v], %v", command, args, err)
		signalChan <- syscall.SIGKILL
	}

	err = cmd.Wait()
	log.Errorf("confd command finished with error: %v", err)
	signalChan <- syscall.SIGKILL
}

func checkError(signalChan chan os.Signal, err error) {
	if err != nil {
		log.Errorf("%v", err)
		signalChan <- syscall.SIGKILL
	}
}

func checkNumberOfErrors(std io.ReadCloser, count uint64, tick time.Duration, signalChan chan os.Signal) {
	testRegex := regexp.MustCompile(templateErrorRegex)
	var tickErrors uint64
	lines := make(chan string)
	go func() {
		scanner := bufio.NewScanner(std)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			lines <- scanner.Text()
		}
	}()

	timer := time.NewTicker(tick)

	for {
		select {
		case <-timer.C:
			if tickErrors > count {
				log.Debugf("number of errors %v", tickErrors)
				log.Error("too many confd errors in the last minute. restarting component")
				signalChan <- syscall.SIGKILL
				return
			}

			tickErrors = 0
		case line := <-lines:
			match := testRegex.FindStringSubmatch(line)
			if match != nil {
				tickErrors++
			}
		}
	}
}
