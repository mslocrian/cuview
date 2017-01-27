package definitions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var (
	BaseDirectory *string
)

type Envelope struct {
	Type string
	Msg  interface{}
}

type SwaggerDef struct {
	SwaggerVer      string `json:"swagger"`
	Paths           map[string]*Endpoint
	InterfacePaths  map[string]interface{} `json:"paths"`
	CumulusCommands CumulusCommands        `json:"x-cumulus-commands"`
	Info            Info                   `json:"info"`
}

type CumulusCommands struct {
	NetdSocket  string `json:"netdSocket"`
	NetdCommand  string `json:"netdCommand"`
	Vtysh string `json:"vtysh"`
}

type CumulusOption struct {
	Netd    bool
	Command string
}

type Info struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type Endpoint struct {
	RequestMethods map[string]*RequestMethod
}

type RequestMethod struct {
	Summary        string
	Description    string
	CumulusOptions *CumulusOption
	//Parameters     []*Parameter
	Parameters     map[string]*Parameter
}

type Parameter struct {
	Name        string
	In          string
	Description string
	Type        string
	Required    bool
}

func (s *SwaggerDef) InitPaths() error {
	var err error
	s.Paths = make(map[string]*Endpoint)
	return err
}

func (s *SwaggerDef) AddPath(p string) error {
	var err error
	if s.Paths == nil {
		s.InitPaths()
		s.Paths[p] = &Endpoint{}
	} else {
		s.Paths[p] = &Endpoint{}
	}
	return err
}

func (s *Endpoint) InitRequestMethods(p string) error {
	var err error
	s.RequestMethods = make(map[string]*RequestMethod)
	return err
}

func (s *Endpoint) AddRequestMethod(m string) error {
	var err error
	if s.RequestMethods == nil {
		s.InitRequestMethods(m)
		s.RequestMethods[m] = &RequestMethod{}
	} else {
		s.RequestMethods[m] = &RequestMethod{}
	}
	return err
}

func (s *RequestMethod) InitParameters() error {
	var err error
	s.Parameters = make(map[string]*Parameter)
	return err
}

func (s *RequestMethod) AddParameters(opts interface{}) error {
	var err error

	if s.Parameters == nil {
		s.InitParameters()
	}

	for _, entry := range opts.([]interface{}) {
		p := NewParameter()
		for k, v := range entry.(map[string]interface{}) {
			switch k {
			case "name":
				p.Name = v.(string)
			case "in":
				p.In = v.(string)
			case "description":
				p.Description = v.(string)
			case "type":
				p.Type = v.(string)
			case "required":
				p.Required = v.(bool)
			default:
				continue
			}
		}
		//s.Parameters = append(s.Parameters, p)
		s.Parameters[p.Name] = p
	}
	return err
}

func (s *RequestMethod) AddCumulusOpts(opts interface{}) error {
	var err error
	copts := &CumulusOption{}
	for k, v := range opts.(map[string]interface{}) {
		switch k {
		case "command":
			copts.Command = v.(string)
		case "netd":
			copts.Netd = v.(bool)
		default:
			continue
		}
	}
	s.CumulusOptions = copts
	return err
}

func (s *RequestMethod) AddOptions(key string, opts interface{}) error {
	var err error
	switch key {
	case "summary":
		s.Summary = opts.(string)
	case "description":
		s.Description = opts.(string)
	case "parameters":
		s.AddParameters(opts)
	case "x-cumulus-options":
		s.AddCumulusOpts(opts)
		return err
	default:
		return err
	}
	return err
}

func (s *SwaggerDef) GetRoutes() []string {
	var routes []string
	for k, _ := range s.Paths {
		routes = append(routes, strings.TrimSpace(k))
	}
	return routes
}

func (s *SwaggerDef) GetRequestMethods(route string) []string {
	var methods []string
	for k, _ := range s.Paths[route].RequestMethods {
		methods = append(methods, strings.TrimSpace(k))
	}
	return methods
}

func (s *SwaggerDef) GetURLParameters(route string, method string) map[string]*Parameter {
	return s.Paths[route].RequestMethods[method].Parameters
}

func (s *SwaggerDef) GetCumulusOptions(route string, method string) *CumulusOption {
	return s.Paths[route].RequestMethods[method].CumulusOptions
}

func NewParameter() *Parameter {
	return &Parameter{Required: false}
}

func LoadAPIDefs() SwaggerDef {
	var swagDef SwaggerDef
	jsonData, err := ioutil.ReadFile(*BaseDirectory + "/definitions/swagger.json")
	if err != nil {
		fmt.Printf("Could not open definition file: %s\n", err.Error())
		os.Exit(1)
	}

	err = json.Unmarshal(jsonData, &swagDef)
	if err != nil {
		fmt.Printf("Could not Unmarshal JSON input: %s\n", err.Error())
		os.Exit(1)
	}

	for pathKey, _ := range swagDef.InterfacePaths {
		swagDef.AddPath(pathKey)
		for reqMethodKey, _ := range swagDef.InterfacePaths[pathKey].(map[string]interface{}) {
			swagDef.Paths[pathKey].AddRequestMethod(reqMethodKey)
			for reqOptionsKey, reqOptionsVal := range swagDef.InterfacePaths[pathKey].(map[string]interface{})[reqMethodKey].(map[string]interface{}) {
				swagDef.Paths[pathKey].RequestMethods[reqMethodKey].AddOptions(reqOptionsKey, reqOptionsVal)
			}
		}
	}

	return swagDef
}
