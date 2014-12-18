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
		if c.Code() != "" {
			t.Errorf("Expected \"\" got %s", c.code)
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
		if Code() != "" {
			t.Errorf("Expected \"\" got %s", c.code)
		}
		if c.UseEnv() != false {
			t.Error("Expected false got %s", strconv.FormatBool(c.UseEnv()))
		}
	}

}

func TestSetCode(t *testing.T) {
	initCfgs()

	tests := []struct {
		name        string
		code        string
		expected    string
		expectedErr string
	}{
		{"set empty", "", "", ""},
		{"set code", "val", "val", ""},
		{"set already set code", "newval", "val", "this configuration's code is already set and cannot be overridden"},
	}

	tstCfg := NewCfg("testCfg")

	// Cfg
	for _, test := range tests {
		err := tstCfg.SetCode(test.code)
		if err != nil {
			if test.expectedErr == err.Error() {
				continue
			}
			t.Errorf("Test %s: Expected %q got %q", test.name, test.expectedErr, err)
		}
		if tstCfg.Code() != test.expected {
			t.Errorf("Test %s: Expected %q got %q", test.name, test.expected, testCfg.Code())
		}
	}

	// test the funcs, app config
	for _, test := range tests {
		err := SetCode(test.code)
		if err != nil {
			if test.expectedErr == err.Error() {
				continue
			}
			t.Errorf("Test %s: Expected %q got %q", test.name, test.expectedErr, err)
		}
		if appCfg.Code() != test.expected {
			t.Errorf("Test %s: Expected %q got %q", test.name, test.expected, testCfg.Code())
		}
	}
}

/*
func TestCfgProcessed(t *testing.T) {
	tests := []struct{
		name string
		value bool
	}
}
*/
