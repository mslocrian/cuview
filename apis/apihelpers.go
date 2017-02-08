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
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net"
	"os/exec"
	"regexp"
	"strings"
	"sync"

	//"github.com/jeffail/gabs"
	"github.com/tdewolff/minify"
	minify_json "github.com/tdewolff/minify/json"
)

func Minify(s []byte) []byte {
	var (
		err    error
		output bytes.Buffer
	)
	m := minify.New()
	writer := bufio.NewWriter(&output)
	m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), minify_json.Minify)
	if err = m.Minify("application/json", writer, bytes.NewReader(s)); err != nil {
		return s
	} else {
		writer.Flush()
		return output.Bytes()
	}
	return s
}

func minifyOutput(p map[string][]string, s []byte) (map[string][]string, []byte) {
	if _, ok := p["minify"]; ok {
		if strings.ToLower(p["minify"][0]) == "false" {
			delete(p, "minify")
			return p, s
		} else {
			delete(p, "minify")
			return p, Minify(s)
		}
	}
	delete(p, "minify")
	return p, Minify(s)
}

func runNetdCommand(netCmdSock string, netCmd string, c string) ([]byte, error) {
	var (
		err    error
		cmdOut bytes.Buffer
	)
	cmd := strings.Split(c, " ")
	cmd = append([]string{netCmd}, cmd...)
	cmd = append(cmd, "json")

	command, err := json.Marshal(cmd)
	if err != nil {
		return cmdOut.Bytes(), err
	}

	conn, err := net.Dial("unix", netCmdSock)
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

func runVtyshCommand(cCmd string, cmd string) ([]byte, error) {
	var (
		err       error
		cmdOut    bytes.Buffer
		runCmd    *exec.Cmd
		waitGroup sync.WaitGroup
	)
	cmd = "-c " + cmd + " json"
	runCmd = exec.Command(cCmd, cmd)
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
