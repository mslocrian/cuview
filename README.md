cuview (Pronounced: Queue View) is a web frontend to expose Cumulus Linux command output in JSON format via a web interface

# Building cuview
## Via go
`go get github.com/mslocrian/cuview`
## Via source
`git clone https://github.com/mslocrian/cuview`
`cd cuview && make`

# Using cuview
cuview routes can be defined, and redefined, via the swagger.json file under the definitions/ directory.  This swagger.json is required in order for cuview to run.  The swagger.json is expected to be found under a definitions/ directory.  

## swagger.json isms
The swagger.json is fully swagger 2.0 compliant, and there are some custom fields which we have added.  Under the global configuration there are these changs:
```
"x-cumulus-commands": {
    "netdSocket": "/var/run/nclu/uds",
    "netdCommand": "/usr/bin/net",
    "vtysh": "/usr/bin/vtysh"
},
```
These entries define the paths to various cumulus commands.  This is something that should not be changed.

Under each path endpoint definition, there are the following definitions:
```
"x-cumulus-options": {
    "netd": true,
    "command": "show interface",
    "parameter_handler": "GetInterfaceParams"
}
```
These definitions define whether or not Netd should be called, or vtysh should be called.  It also defines the command which should be executed.  The parameter_handler is a function which will need to be written under `apis/parameterhandlers.go`.  If a handler is not written, then no action will be taken.  The parameter handler is if you want to do anything additional with the output via parameters from within the query string.  A prime example would be filtering output.

# Running cuview
To run cuview, run the following.
`/path/to/cuview -base.dir /some/path` where `/some/path` would contain the `definitions/swagger.json`.
