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

	"github.com/gorilla/mux"
	//"github.com/prometheus/common/log"
)

type ApiRoute struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// type ApiRoutes []ApiRoute

type ApiMgr struct {
	//logger		*logging.Writer
	//clientMgr	*clients.ClientMgr
	//objectMgr	*objects.ObjectMgr
	//actionMgr	*actions.ActionMgr
	apiVer        string
	apiBase       string
	apiBaseState  string
	pRestRtr      *mux.Router
	restRoutes    []ApiRoute
	apiCallSeqNum uint32
}

var gApiMgr *ApiMgr

func InitializeApiMgr() *ApiMgr {
	//var err error
	mgr := new(ApiMgr)
	mgr.apiVer = "v1"
	mgr.apiBase = "/api/" + mgr.apiVer + "/"
	mgr.apiBaseState = mgr.apiBase + "state/"
	gApiMgr = mgr
	return mgr
}

func (mgr *ApiMgr) InitializeRestRoutes() bool {
	var rt ApiRoute

	// Need to get a better way of defining these routes.

	/* /public/v1/state/bgpv4neighbors */
	rt = ApiRoute{"bgpv4neighbors",
		"GET",
		mgr.apiBaseState + "bgpv4neighbors",
		GetCumulusBGPv4Neighbors,
	}
	mgr.restRoutes = append(mgr.restRoutes, rt)

	/* /public/v1/state/ipv4routes/{id} */
	rt = ApiRoute{"getroutesbyid",
		"GET",
		mgr.apiBaseState + "ipv4routes" + "/" + "{objId}",
		GetCumulusIPv4RoutesById,
	}
	mgr.restRoutes = append(mgr.restRoutes, rt)

	/* /public/v1/state/ipv4routes */
	rt = ApiRoute{"getroutes",
		"GET",
		mgr.apiBaseState + "ipv4routes",
		GetCumulusIPv4Routes,
	}
	mgr.restRoutes = append(mgr.restRoutes, rt)

	/* /public/v1/state/interfaces/{id} */
	rt = ApiRoute{"getinterfacesbyid",
		"GET",
		mgr.apiBaseState + "interfaces" + "/" + "{objId}",
		GetCumulusInterfacesById,
	}
	mgr.restRoutes = append(mgr.restRoutes, rt)

	/* /public/v1/state/interfaces */
	rt = ApiRoute{"getinterfaces",
		"GET",
		mgr.apiBaseState + "interfaces",
		GetCumulusInterfaces,
	}
	mgr.restRoutes = append(mgr.restRoutes, rt)

	return true
}

func (mgr *ApiMgr) InstantiateRestRtr() *mux.Router {
	mgr.pRestRtr = mux.NewRouter().StrictSlash(true)
	for _, route := range mgr.restRoutes {
		//var handler http.Handler
		mgr.pRestRtr.Methods(route.Method).Path(route.Pattern).Handler(route.HandlerFunc)
	}
	return mgr.pRestRtr
}

func (mgr *ApiMgr) GetRestRtr() *mux.Router {
	return mgr.pRestRtr
}
