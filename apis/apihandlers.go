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
			minifyOutput(w, cmdOut)
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
		minifyOutput(w, cmdOut)
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
		minifyOutput(w, cmdOut)
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
