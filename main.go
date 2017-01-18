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

package main

import (
	"flag"
	"net/http"
	_ "net/http/pprof"
	"strconv"

	"github.com/mslocrian/cuview/apis"

	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
)

var (
	listenAddress	= flag.String("listen.address", "127.0.0.1", "IP Address to bind webserver to.")
	listenPort	= flag.Int("listen.port", 9000, "Port to bind webserver to.")
)

func main() {
	flag.Parse()
	log.Infoln("Starting cuview", version.Info())
	log.Infoln("Build Context", version.BuildContext())

        mgr := apis.InitializeApiMgr()
	mgr.InitializeRestRoutes()
	mgr.InstantiateRestRtr()
	restRtr := mgr.GetRestRtr()

	bindTuple := *listenAddress + ":" + strconv.Itoa(*listenPort)
	log.Infof("Listening on %s", bindTuple)
	log.Fatal(http.ListenAndServe(bindTuple, restRtr))
}
