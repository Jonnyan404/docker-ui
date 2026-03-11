//go:build !embed

package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	. "github.com/gohutool/boot4go-docker-ui/log"
	util4go "github.com/gohutool/boot4go-util"
)

func setLicenseJS(li string) {
	var txt string
	if util4go.IsEmpty(li) {
		txt = ""
	} else {
		txt = `
				myConfig.li="%v";
			`
		txt = fmt.Sprintf(txt, li)
	}

	err := ioutil.WriteFile("./html/static/public/js/cubeui.li.js", []byte(txt), 0666)
	if err != nil {
		panic("License check error:" + err.Error())
	}
}

func setEndpointJS(endpoint string) {
	if util4go.IsEmpty(endpoint) {
		Logger.Info("Local Docker Endpoint is attached")
		return
	}

	v2 := strings.Split(endpoint, ":")
	host, port := "unix", "2375"
	if len(v2) >= 1 {
		host = v2[0]
	}
	if len(v2) >= 2 {
		port = v2[1]
	}

	txt := `
				local_node.node_host = "%v";
				local_node.node_port = "%v";
			`
	txt = fmt.Sprintf(txt, host, port)

	Logger.Info("%v Docker Endpoint is attached", endpoint)

	err := ioutil.WriteFile("./html/api/node.config.js", []byte(txt), 0666)
	if err != nil {
		panic("Endpoint check error:" + err.Error())
	}
}
