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
	"reflect"
)

// shell holder
type ParameterHandler struct {
}

func (h *ParameterHandler) GetInterfacesParams(p map[string][]string, s []byte) []byte {
	// some work shoud go here
	return s
}

func (h *ParameterHandler) GetBgpv4NeighborsParams(p map[string][]string, s []byte) []byte {
	// some work should go here
	return s
}

func (h *ParameterHandler) GetIpv4RoutesParams(p map[string][]string, s []byte) []byte {
	// some work should go here
	/*
		tmp := "HEY THERE!"
		s = append(s, tmp...)
	*/
	return s
}

func CallParameterHandlerFunc(c interface{}, funcName string, params ...interface{}) (out []reflect.Value, err error) {
	function := reflect.ValueOf(c)
	m := function.MethodByName(funcName)
	if !m.IsValid() {
		return make([]reflect.Value, 0), fmt.Errorf("Method not found \"%s\"\n", funcName)
	}

	in := make([]reflect.Value, len(params))
	for i, param := range params {
		in[i] = reflect.ValueOf(param)
	}

	out = m.Call(in)
	return
}
