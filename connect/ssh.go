package connect

import (
	"github.com/helloyi/go-sshclient"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type SshClient struct {
	Hostname string
	Port     int
	Username string
	Password string
	KeyFile  string
}

func (c SshClient) SshDiskCheck() string {
	var useKeyFile = false
	var client *sshclient.Client
	var err error
	panicMessage := ""
	switch c.Password {
	case "":
		useKeyFile = true
	default:
		useKeyFile = false
	}
	addr := c.Hostname + ":" + strconv.Itoa(c.Port)
	if useKeyFile {
		client, err = sshclient.DialWithKey(addr, c.Username, c.KeyFile)
		if err != nil {
			client, err = sshclient.DialWithPasswd(addr, c.Username, c.Password)
		}
	} else {
		client, err = sshclient.DialWithPasswd(addr, c.Username, c.Password)
	}
	if err != nil {
		log.Fatalln("SSH connection failed: " + err.Error())
	}
	defer client.Close()
	output, _ := client.Script("echo ----------; df -h").SmartOutput()
	lines := strings.Split(strings.Split(string(output), "----------")[1], "\n")
	for _, line := range lines {
		if strings.Contains(line, "/dev/") {
			fields := strings.Fields(line)
			if loopback, _ := regexp.Match("loop", []byte(fields[0])); loopback {
				continue
			}
			if tmpfs, _ := regexp.Match("tmpfs", []byte(fields[0])); tmpfs {
				continue
			}
			if len(fields) == 6 {
				percentage, _ := strconv.Atoi(fields[4][:len(fields[4])-1])
				if percentage > 90 {
					panicMessage += "Disk usage is over 90% on " + fields[0]
				}
			}
		}
	}
	return panicMessage
}
