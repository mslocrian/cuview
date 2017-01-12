package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os/exec"

	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
)

var (
	routeCommand = "/usr/bin/vtysh"
	routeCommandArgs = []string{"-c", "show ip route json"}
)

func handler(w http.ResponseWriter, r *http.Request) {
	var (
		cmdOut []byte
		err error
	)
	if cmdOut, err = exec.Command(routeCommand, routeCommandArgs...).Output(); err != nil {
		fmt.Fprintf(w, "There was an error running %s: %s", routeCommand, err)
	} else {
		w.Write([]byte(cmdOut))
	}
}

func main() {
	flag.Parse()
	log.Infoln("Starting cuview", version.Info())
	log.Infoln("Build Context", version.BuildContext())
        http.HandleFunc("/test", handler)

        listenAddress := "127.0.0.1:9000"
	log.Infof("Listening on %s", listenAddress)
        log.Fatal(http.ListenAndServe(listenAddress, nil))
}
