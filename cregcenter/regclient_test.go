package cregcenter_test

import (
	"fmt"
	"go_common/cregcenter"
	"net/http"
	"os"
	"testing"
)

var (
	port           = "8888"
	callTestServer *cregcenter.CallServerClient
)

func init() {
	err := cregcenter.InitRegInfo(&cregcenter.RegInfo{
		LoginName:     "hdikjagi3",
		LoginPassword: "jroi32hd987e",
		RegCenterAddr: "10.0.0.26:2748",
		Workspace:     "dev",
		ServerName:    "testRegCenter",
		Port:          port,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	callTestServer, err = cregcenter.NewCallClient("testServer")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

func TestReg(t *testing.T) {
	fmt.Println(111)
	bytes, err := callTestServer.Call("POST", "/test", nil, []byte("{}"), 10)
	fmt.Println(string(bytes), err)
	http.ListenAndServe(":"+port, nil)
}
