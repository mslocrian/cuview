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
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"time"

	"github.com/mslocrian/cuview/apis"
	"github.com/mslocrian/cuview/definitions"

	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
)

var (
	listenAddress = flag.String("listen.address", "127.0.0.1", "IP Address to bind webserver to.")
	listenPort    = flag.Int("listen.port", 9000, "Port to bind webserver to.")
	baseDirectory = flag.String("base.dir", "/usr/local/cuview", "Path to installation base")
)

func LogWebRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ts string
		t := time.Now()
		ts = fmt.Sprintf("%s", t.Format("Mon Jan _2 15:04:05 2006"))
		log.Infof("%s [%s] %s %s (%s)", r.RemoteAddr, ts, r.Method, r.URL, r.Proto)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	flag.Parse()
	log.Infoln("Starting cuview", version.Info())
	log.Infoln("Build Context", version.BuildContext())
	apis.BaseDirectory = baseDirectory
	definitions.BaseDirectory = baseDirectory

	mgr := apis.InitializeApiMgr()
	mgr.InitializeRestRoutes(definitions.LoadAPIDefs())
	mgr.InstantiateRestRtr()
	restRtr := mgr.GetRestRtr()

	bindTuple := *listenAddress + ":" + strconv.Itoa(*listenPort)
	log.Infof("Listening on %s", bindTuple)
	log.Fatal(http.ListenAndServe(bindTuple, LogWebRequest(restRtr)))
}
