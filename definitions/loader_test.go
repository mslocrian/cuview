package definitions

import "testing"
import "flag"
import "fmt"
import "log"
import "os/user"

func init() {
	flag.Parse()

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	str := fmt.Sprintf("%s/git/cuview", usr.HomeDir)
	BaseDirectory = &str

}

func TestLoadApiDefs(t *testing.T) {
	defs := LoadAPIDefs()
	expected_version := "2.0"
	if expected_version != defs.SwaggerVer {
		t.Errorf("expected_version == %q, want %q",
			defs.SwaggerVer, expected_version)
	}
}
