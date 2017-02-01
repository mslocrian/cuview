/*
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package apis

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"

	defs "github.com/mslocrian/cuview/definitions"
)

var (
	vtyshCommand = "/usr/bin/vtysh"
	netCommand   = "/usr/bin/net"
	netdSocket   = "/var/run/nclu/uds"

	// Trailing spaces should remain
	vtyshRouteCommandArgs       = "-c show ip route "
	netIfaceCommandArgs         = []string{netCommand, "show", "interface"}
	netBGPv4NeighborCommandArgs = []string{netCommand, "show", "bgp", "ipv4", "unicast", "summary"}
	netLldpArgs                 = []string{netCommand, "show", "lldp"}
)

type CumulusHTTPHandler struct {
	handler   http.Handler
	routeData ApiRoute
}

func (ch *CumulusHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var s string
	var output []byte
	//var validParams map[string]interface{}
	var validParams map[string][]string

	rec := httptest.NewRecorder()
	ch.handler.ServeHTTP(rec, r)
	data_before := []byte("text shoved in before\n")
	data_after := []byte("text shoved in after\n")

	// Set headers
	for k, v := range rec.Header() {
		w.Header()[k] = v
	}
	w.Header().Set("X-Cumulus-API", "True")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	clen, _ := strconv.Atoi(r.Header.Get("Content-Length"))

	queryParams := r.URL.Query()
	validParams = make(map[string][]string)
	for k, v := range queryParams {
		if _, ok := ch.routeData.Parameters[k]; ok {
			validParams[k] = v
		}
	}

	switch ch.routeData.Method {
	case "GET":
		var err error
		output, err = runCommand(ch.routeData.Parameters, ch.routeData.Options, ch.routeData.Commands)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "could not run command!\n")
			return
		}
		w.WriteHeader(200)
	default:
		w.WriteHeader(405)
	}

	clen += len(data_before) + len(data_after) + len(s) + len(output)

	/*
		w.Write(data_before)
		w.Write(rec.Body.Bytes())
		w.Write(data_after)
		w.Write([]byte(s))
	*/

	// holding off for now.. some bug
	validParams, output = minifyOutput(validParams, output)
	ph := &ParameterHandler{}
	out, err := CallParameterHandlerFunc(ph, ch.routeData.Options.ParamHandler, validParams, output)
	if err != nil {
		fmt.Printf("Caught error in handler func: %s\n", err)
		w.Write(output)
		return
	}
	output = out[0].Bytes()
	w.Write(output)
	return
}

func runCommand(params map[string]*defs.Parameter, co *defs.CumulusOption, cc defs.CumulusCommands) ([]byte, error) {
	var (
		cmdOut []byte
		err    error
	)
	if co.Netd == true {
		cmdOut, err = runNetdCommand(cc.NetdSocket, cc.NetdCommand, co.Command)
	} else {
		cmdOut, err = runVtyshCommand(cc.Vtysh, co.Command)
	}
	return cmdOut, err
}

func GetCumulusIPv4Routes(w http.ResponseWriter, r *http.Request) {
	var (
		cmdOut []byte
		err    error
	)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	cmdOut, err = runCumulusVtyshCommand(vtyshCommand, vtyshRouteCommandArgs+"json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "There was an error running %s: %s\n", vtyshCommand, err)
	} else {
		if doMinifyOutput(r) == false {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(cmdOut))
		} else {
			w.WriteHeader(http.StatusOK)
			minifyOutputOrig(w, cmdOut)
		}
	}
	return
}

func GetCumulusInterfaces(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		cmdOut []byte
	)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	cmdOut, err = runCumulusNetdCommand(netIfaceCommandArgs)

	if err != nil {
		fmt.Fprintf(w, "There was an error running netd command: %s\n", err)
		return
	}

	if doMinifyOutput(r) {
		w.WriteHeader(http.StatusOK)
		minifyOutputOrig(w, cmdOut)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(cmdOut)
	}
	return
}

func GetCumulusBGPv4Neighbors(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		cmdOut []byte
	)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	cmdOut, err = runCumulusNetdCommand(netBGPv4NeighborCommandArgs)

	if err != nil {
		fmt.Fprintf(w, "There was an error running netd command: %s\n", err)
		return
	}

	if doMinifyOutput(r) {
		w.WriteHeader(http.StatusOK)
		minifyOutputOrig(w, cmdOut)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(cmdOut)
	}
	return
}

func GetCumulusIPv4RoutesById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	output := `{"GetCumulusIPv4RoutesById": "OK"}`
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(output))
	return
}

func GetCumulusInterfacesById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	output := `{"GetCumulusIPv4RoutesById": "OK"}`
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(output))
	return
}

func CumulusHandler(w http.ResponseWriter, r *http.Request) {
	return
}

func GetCumulusHTTPHandler(handler http.Handler, route ApiRoute) *CumulusHTTPHandler {
	return &CumulusHTTPHandler{handler: handler, routeData: route}
}
