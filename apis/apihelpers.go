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
        "bytes"
        "encoding/json"
        "io"
        "net"
        "net/http"
        "os/exec"
        "regexp"
        "strings"
	"sync"

        "github.com/tdewolff/minify"
        minify_json "github.com/tdewolff/minify/json"

)

func doMinifyOutput(r *http.Request) bool {
        minifyParam := r.URL.Query().Get("minify")
        if minifyParam == "" {
                return true
        } else if strings.ToLower(minifyParam) == "true" {
                return true
        } else {
                return false
        }
}

func minifyOutput(w http.ResponseWriter, s []byte) {
        var (
                err error
        )
        m := minify.New()
        m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), minify_json.Minify)
        if err = m.Minify("application/json", w, bytes.NewReader(s)); err != nil {
                w.Write([]byte(s))
        }
        return
}

func runCumulusNetdCommand(cmd []string) ([]byte, error) {
        var (
                err     error
                cmdOut  bytes.Buffer
        )
        cmd = append(cmd, "json")
        command, err := json.Marshal(cmd)
        if err != nil {
                return cmdOut.Bytes(), err
        }

        conn, err := net.Dial("unix", netdSocket)
        if err != nil {
                return cmdOut.Bytes(), err
        }
        defer conn.Close()

        _, err = conn.Write(command)
        if err != nil {
                return cmdOut.Bytes(), err
        }

        io.Copy(&cmdOut, conn)
        return cmdOut.Bytes(), err
}

func runCumulusVtyshCommand(cmd string, args string) ([]byte, error) {
	var (
		err		error
		cmdOut		bytes.Buffer
		runCmd		*exec.Cmd
		waitGroup	sync.WaitGroup
	)
	runCmd = exec.Command(cmd, args)
	stdout, _ := runCmd.StdoutPipe()

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		io.Copy(&cmdOut, stdout)
	}()

	err = runCmd.Run()

	waitGroup.Wait()

	return cmdOut.Bytes(), err
}
