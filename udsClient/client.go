package udsClient

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

type SocketClient struct {
	Client    http.Client
	SocketUrl string
}

func (Client *SocketClient) CreateExec(containerName string, cmd []string) string {
	jsonStr, _ := json.Marshal(map[string]interface{}{
		"AttachStdin":  false,
		"AttachStdout": true,
		"AttachStderr": true,
		"Tty":          true,
		"Cmd":          cmd,
	})
	url := Client.SocketUrl + "/containers/" + containerName + "/exec"
	resp := Client.post(url, string(jsonStr))
	defer resp.Body.Close()

	// 读取id
	body, _ := ioutil.ReadAll(resp.Body)
	res := make(map[string]interface{})
	json.Unmarshal(body, &res)

	return res["Id"].(string)
}

func (Client *SocketClient) StartExec(execId string) {
	jsonStr, _ := json.Marshal(map[string]interface{}{
		"Detach":      false,
		"Tty":         true,
		"ConsoleSize": []int{80, 64},
	})
	url := Client.SocketUrl + "/exec/" + execId + "/start"

	resp := Client.post(url, string(jsonStr))
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func (Client *SocketClient) post(url string, payloadStr string) *http.Response {
	payload := strings.NewReader(payloadStr)
	resp, err := Client.Client.Post(url, "application/json", payload)
	if err != nil {
		panic(err)
	}
	return resp
}

// unix domain socket 客户端
func NewClient(sockAddr, sockUrl string) SocketClient {
	return SocketClient{
		Client: http.Client{
			Transport: &http.Transport{
				DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
					return net.Dial("unix", sockAddr)
				},
			},
		},
		SocketUrl: sockUrl,
	}
}
