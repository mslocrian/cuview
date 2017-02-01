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
	"fmt"
	"io"
	"net"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"github.com/jeffail/gabs"
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

func minifyOutputOrig(w http.ResponseWriter, s []byte) {
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

func handleParams(p map[string][]string, s []byte) []byte {
	jsonParsed, err := gabs.ParseJSON(s)
	if err != nil {
		return s
	}

	exists := jsonParsed.Exists("swp1")
	fmt.Printf("exists=%s\n", exists)
	exists = jsonParsed.Exists("heya")
	fmt.Printf("exists=%s\n", exists)

	value1, ok1 := jsonParsed.Path("hi").Data().(string)
	fmt.Printf("value1=%s\nok1=%s\n", value1, ok1)
	value2, ok2 := jsonParsed.Path("swp1").Data().(string)
	fmt.Printf("value2=%s\nok2=%s\n", value2, ok2)

	return s
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

func runCumulusNetdCommand(cmd []string) ([]byte, error) {
	var (
		err    error
		cmdOut bytes.Buffer
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
		err       error
		cmdOut    bytes.Buffer
		runCmd    *exec.Cmd
		waitGroup sync.WaitGroup
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
