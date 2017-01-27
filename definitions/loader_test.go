package definitions

import "testing"

func TestLoadApiDefs(t *testing.T) {
	defs := LoadAPIDefs()
	expected_version := "2.0"
	if expected_version != defs.SwaggerVer {
		t.Errorf("expected_version == %q, want %q",
			defs.SwaggerVer, expected_version)
	}
}
