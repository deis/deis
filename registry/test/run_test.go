package verbose

import (
		"github.com/dotcloud/docker/api/client"
        "bufio"
        "fmt"
        "io"
        "time"
        "testing"
        "strings"
				"os"

)



const (
        unitTestImageName        = "docker-test-image"
        unitTestImageID          = "83599e29c455eb719f77d799bc7c51521b9551972f5a850d7ad265bc1b5292f6" // 1.0
        unitTestImageIDShort     = "83599e29c455"
        unitTestNetworkBridge    = "testdockbr0"
        unitTestStoreBase        = "/var/lib/docker/unit-tests"
        testDaemonAddr           = "192.168.59.103:2375"
        testDaemonProto          = "tcp"
        testDaemonHttpsProto     = "tcp"
        testDaemonHttpsAddr      = "localhost:4271"
        testDaemonRogueHttpsAddr = "localhost:4272"
)

func closeWrap(args ...io.Closer) error {
	e := false
	ret := fmt.Errorf("Error closing elements")
	for _, c := range args {
		if err := c.Close(); err != nil {
			e = true
			ret = fmt.Errorf("%s\n%s", ret, err)
		}
	}
	if e {
		return ret
	}
	return nil
}


func getDetails(/*format string,t *testing.T*/) string{
	stdin, _   := io.Pipe()
	stdout, stdoutPipe := io.Pipe()
	fmt.Println("1")
	cli := client.NewDockerCli(nil, stdoutPipe, nil, testDaemonProto, testDaemonAddr, nil)
	var IPAdress string
	go func(){
		err1:= cli.CmdInspect("--format","'{{ .NetworkSettings.IPAddress }}'","deis-etcd")
		if err1 != nil {
			t.Fatalf("%s",err1)
		}
		if err := closeWrap(stdout, stdoutPipe,stdin); err != nil {
			t.Fatalf("%s",err)
		}
	}()

	for{
		if cmdBytes,err:= bufio.NewReader(stdout).ReadString('\n'); err==nil{
			IPAdress=cmdBytes
			fmt.Println(cmdBytes)
		}else{
			break
		}
		return IPAdress
		fmt.Println("look1")
	}


}



func TestBuild(t *testing.T) {
        stdin, _   := io.Pipe()
        stdout, stdoutPipe := io.Pipe()
        fmt.Println("1")
        cli := client.NewDockerCli(nil, stdoutPipe, nil, testDaemonProto, testDaemonAddr, nil)
        fmt.Println("2")
				c := make(chan int)


				go func (){
					err1:= cli.CmdRun("--name","deis-etcd","coreos/etcd")
					if err1 != nil {
						t.Fatalf("%s",err1)
					}
					c<-1
				}( )


				go func (){
					<-c
        	err1:= cli.CmdBuild("-t","../")
        		if err1 != nil {
							t.Fatalf("%s",err1)
						}
						if err := closeWrap(stdout, stdoutPipe,stdin); err != nil {
							t.Fatalf("%s",err)
						}
				}( )

        fmt.Println("3")
        time.Sleep(3000 * time.Millisecond)
        var imageId string
        fmt.Println("4")

				for{
					if cmdBytes,err:= bufio.NewReader(stdout).ReadString('\n'); err==nil{
				    imageId=cmdBytes
		  			fmt.Println(cmdBytes)
		  		}else{
		  			break
		  		}
		  		fmt.Println("look2")
				}


				words := strings.Fields(imageId)
				fmt.Println("Stirng ID is "+words[2])

				CmdIp="docker inspect --format '{{ .NetworkSettings.IPAddress }}' deis-etcd"
				cmd = exec.Command("docker", "-c", cmdIp)

				cmdString1="docker inspect deis-registry-data >/dev/null 2>&1 || docker run --name deis-registry-data -v /data deis/base /bin/true"
				cmdString2="docker run --name deis-registry -p 5000:5000 -e PUBLISH=5000 -e HOST=${COREOS_PRIVATE_IPV4} --volumes-from deis-registry-data deis/registry"
				cmd := exec.Command("sh", "-c", cmdString1)
				out,err1 := runCommandWithStdoutStderr(cmd)
				if err1 !=nil {
					t.fatalf(err1)
				}

				cmd = exec.Command("sh", "-c", cmdString2)
				out,err1 = runCommandWithStdoutStderr(cmd)
				if err1 !=nil {
					t.fatalf(err1)
				}


}
