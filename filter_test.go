package contour

import (
	"strings"
	"testing"
)

//func TestGetFilterNames(t *testing.T) {
//	appCfg = newTestCfg()
//	filterNames := appCfg.GetBoolFilterNames()
//
//}

func TestGetBoolFilter(t *testing.T) {
	tests := []struct {
		name     string
		short    string
		value    bool
		expected bool
		err      string
	}{
		{"", "", false, false, ""},
		{"corebool", "c", true, true, ""},
		{"bool", "b", false, false, ""},
		{"cfgbool", "", true, true, ""},
		{"flagbool", "f", false, false, ""},
		{"cb", "", false, false, ""},
	}
	appCfg = newTestCfg()
	appCfg.SetName("rancher-test")
	// check that an unfiltered boolflag's var is nil
	if appCfg.filterVars["flagbool-tst"] != nil {
		t.Errorf("flags that weren't present in the args should be nil, flagbool-test wasn't")
	}
	args, err := appCfg.FilterArgs([]string{"-flagbool=false", "-fake", "command"})
	if err != nil {
		if !strings.HasPrefix(err.Error(), "parse of command-line arguments failed: flag provided but not defined: -fake") {
			t.Errorf("expected \" parse of command-line arguments failed: flag provided but not defined: -fake\", got %q", err.Error())
		}
	}
	appCfg = newTestCfg()
	appCfg.SetName("rancher-test")
	args, err = appCfg.FilterArgs([]string{"-flagbool=false", "fake", "command"})
	if err != nil {
		t.Errorf("expected no error, got %q", err.Error())
		return
	}
	if len(args) != 2 {
		t.Errorf("Expected 2 arg to be returned, got %d", len(args))
	}
	for i, test := range tests {
		for _, v := range appCfg.boolFilterNames {
			if v == test.name {
				if appCfg.filterVars[v] == nil {
					t.Errorf("%d: %s was nil", i, v)
				} else {
					switch appCfg.filterVars[v].(type) {
					case bool:
						if appCfg.filterVars[v].(bool) != test.expected {
							t.Errorf("%d expected %v, got %v", i, test.expected, appCfg.filterVars[v].(bool))
						}
					case *bool:
						if *appCfg.filterVars[v].(*bool) != test.expected {
							t.Errorf("%d-%s expected %v, got %v", i, test.name, test.expected, *appCfg.filterVars[v].(*bool))
						}
					}
				}
				break
			}
		}
	}
	// args not filtered shouldn't be affected
	b, err := appCfg.GetBoolE("flagbool-tst")
	if err != nil {
		t.Errorf("Expected no error, got %s", err.Error())
	} else {
		if !b {
			t.Errorf("Expected \"flagbool-tst\" to be true, got %v", b)
		}
	}
}
