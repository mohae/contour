package contour

import (
	"strconv"
	"testing"
)

func TestConfig(t *testing.T) {
	initCfgs()

	c := NewCfg("test")
	if c == nil {
		t.Errorf("New test cfg was nil")
	} else {
		if c.name != "test" {
			t.Errorf("Expected test got %s", c.name)
		}
		if c.UseEnv() != false {
			t.Error("Expected false got %s", strconv.FormatBool(c.UseEnv()))
		}
	}

	a := AppCfg()
	if a == nil {
		t.Errorf("New test cfg was nil")
	} else {
		if a.name != "app" {
			t.Errorf("Expected app got %s", a.name)
		}
		if c.UseEnv() != false {
			t.Error("Expected false got %s", strconv.FormatBool(c.UseEnv()))
		}
	}

}
