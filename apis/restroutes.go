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
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	defs "github.com/mslocrian/cuview/definitions"
)

var (
	BaseDirectory  *string
)

type ApiRoute struct {
	Name        string
	Method      string
	Pattern     string
	Parameters  []*defs.Parameter
	Options     *defs.CumulusOption
	Commands    defs.CumulusCommands
	HandlerFunc http.HandlerFunc
}

// type ApiRoutes []ApiRoute

type ApiMgr struct {
	apiVer        string
	apiBase       string
	apiBaseState  string
	pRestRtr      *mux.Router
	restRoutes    []ApiRoute
	apiCallSeqNum uint32
}

var gApiMgr *ApiMgr

func InitializeApiMgr() *ApiMgr {
	mgr := new(ApiMgr)
	mgr.apiVer = "v1"
	mgr.apiBase = "/api/" + mgr.apiVer + "/"
	mgr.apiBaseState = mgr.apiBase + "state/"
	gApiMgr = mgr
	return mgr
}

func (mgr *ApiMgr) InitializeRestRoutes(defs defs.SwaggerDef) bool {
	var rt ApiRoute

	for _, route := range defs.GetRoutes() {
		for _, method := range defs.GetRequestMethods(route) {
			rt = ApiRoute{strings.Replace(route, "/", method+"_", -1),
				strings.ToUpper(method),
				mgr.apiBaseState + strings.Replace(route, "/", "", -1),
				defs.GetURLParameters(route, method),
				defs.GetCumulusOptions(route, method),
				defs.CumulusCommands,
				CumulusHandler,
			}
			mgr.restRoutes = append(mgr.restRoutes, rt)
		}
	}
	mgr.restRoutes = append(mgr.restRoutes, rt)

	return true
}

func (mgr *ApiMgr) InstantiateRestRtr() *mux.Router {
	mgr.pRestRtr = mux.NewRouter().StrictSlash(true)
	mgr.pRestRtr.PathPrefix("/v2/api-spec/").Handler(http.StripPrefix("/v2/api-spec/", http.FileServer(http.Dir(*BaseDirectory + "/definitions"))))
	mgr.pRestRtr.PathPrefix("/api-docs/").Handler(http.StripPrefix("/api-docs/", http.FileServer(http.Dir(*BaseDirectory + "/api-docs"))))

	for _, route := range mgr.restRoutes {
		//ch := GetCumulusHTTPHandler(route.HandlerFunc, &route)
		ch := GetCumulusHTTPHandler(route.HandlerFunc, route)
		mgr.pRestRtr.Methods(route.Method).Path(route.Pattern).Handler(ch)
	}
	return mgr.pRestRtr
}

func (mgr *ApiMgr) GetRestRtr() *mux.Router {
	return mgr.pRestRtr
}
