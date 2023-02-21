package main

import (
	"dockerClient/udsClient"
	"flag"
	"strings"
)

func main() {
	socketUrl := "http://localhost"
	sockAddr := "/var/run/docker.sock"

	container := flag.String("con", "", "container's name")
	cmdStr := flag.String("cmd", "", "container's name")
	flag.Parse()
	cmdArr := strings.Split(*cmdStr, " ")

	httpc := udsClient.NewClient(sockAddr, socketUrl)
	execId := httpc.CreateExec(*container, cmdArr)
	httpc.StartExec(execId)
}
