package contour

import (
	"strings"
	"testing"
)

func TestGetFilterNames(t *testing.T) {
	appCfg = newTestCfg()
	filterNames := appCfg.GetBoolFilterNames()
	if len(filterNames) != 2 {
		t.Errorf("Expected 1 bool filter, got %d", len(filterNames))
	}
	if filterNames[0] != "flagbool" && filterNames[0] != "flagbool-tst" {
		t.Errorf("Expected the bool flag to be \"flagbool\" or \"flagbool-tst\", got %q", filterNames[0])
	}
	filterNames = appCfg.GetIntFilterNames()
	if len(filterNames) != 2 {
		t.Errorf("Expected 1 int filter, got %d", len(filterNames))
	}
	if filterNames[0] != "flagint" && filterNames[0] != "flagint-tst" {
		t.Errorf("Expected the int flag to be \"flagint\" or \"flagint-tst\", got %q", filterNames[0])
	}
	filterNames = appCfg.GetInt64FilterNames()
	if len(filterNames) != 2 {
		t.Errorf("Expected 1 int64 filter, got %d", len(filterNames))
	}
	if filterNames[0] != "flagint64" && filterNames[0] != "flagint64-tst" {
		t.Errorf("Expected the int64 flag to be \"flagint64\" or \"flagint64-tst\", got %q", filterNames[0])
	}
	filterNames = appCfg.GetStringFilterNames()
	if len(filterNames) != 2 {
		t.Errorf("Expected 1 string filter, got %d", len(filterNames))
	}
	if filterNames[0] != "flagstring" && filterNames[0] != "flagstring-tst" {
		t.Errorf("Expected the string flag to be \"flagstring\" or \"flagstring-tst\", got %q", filterNames[0])
	}
}

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

func TestGetIntFilter(t *testing.T) {
	tests := []struct {
		name     string
		short    string
		value    int
		expected int
		err      string
	}{
		{"", "", 42, 42, ""},
		{"coreint", "c", 42, 42, ""},
		{"int", "b", 42, 42, ""},
		{"cfgint", "", 42, 42, ""},
		{"flagint", "f", 427, 427, ""},
		{"ci", "", 42, 42, ""},
	}
	appCfg = newTestCfg()
	appCfg.SetName("rancher-test")
	args, err := appCfg.FilterArgs([]string{"-flagint=427", "fake", "command"})
	if err != nil {
		t.Errorf("expected no error, got %q", err.Error())
		return
	}
	if len(args) != 2 {
		t.Errorf("Expected 2 arg to be returned, got %d", len(args))
	}
	for i, test := range tests {
		for _, v := range appCfg.intFilterNames {
			if v == test.name {
				if appCfg.filterVars[v] == nil {
					t.Errorf("%d: %s was nil", i, v)
				} else {
					switch appCfg.filterVars[v].(type) {
					case int:
						if appCfg.filterVars[v].(int) != test.expected {
							t.Errorf("%d expected %v, got %v", i, test.expected, appCfg.filterVars[v].(int64))
						}
					case *int:
						if *appCfg.filterVars[v].(*int) != test.expected {
							t.Errorf("%d-%s expected %v, got %v", i, test.name, test.expected, *appCfg.filterVars[v].(*int))
						}
					}
				}
				break
			}
		}
	}
	// args not filtered shouldn't be affected
	i, err := appCfg.GetIntE("flagint-tst")
	if err != nil {
		t.Errorf("Expected no error, got %s", err.Error())
	} else {
		if i != 42 {
			t.Errorf("Expected \"flagint-tst\" to be 42, got %v", i)
		}
	}
}

func TestGetInt64Filter(t *testing.T) {
	tests := []struct {
		name     string
		short    string
		value    int64
		expected int64
		err      string
	}{
		{"", "", int64(42), int64(42), ""},
		{"coreint64", "c", int64(42), int64(42), ""},
		{"int64", "6", int64(42), int64(42), ""},
		{"cfgint64", "", int64(42), int64(42), ""},
		{"flagint64", "f", int64(427), int64(427), ""},
		{"ci64", "", int64(42), int64(42), ""},
	}
	appCfg = newTestCfg()
	appCfg.SetName("rancher-test")
	args, err := appCfg.FilterArgs([]string{"-flagint64=427", "fake", "command"})
	if err != nil {
		t.Errorf("expected no error, got %q", err.Error())
		return
	}
	if len(args) != 2 {
		t.Errorf("Expected 2 arg to be returned, got %d", len(args))
	}
	for i, test := range tests {
		for _, v := range appCfg.int64FilterNames {
			if v == test.name {
				if appCfg.filterVars[v] == nil {
					t.Errorf("%d: %s was nil", i, v)
				} else {
					switch appCfg.filterVars[v].(type) {
					case int64:
						if appCfg.filterVars[v].(int64) != test.expected {
							t.Errorf("%d expected %v, got %v", i, test.expected, appCfg.filterVars[v].(int64))
						}
					case *int64:
						if *appCfg.filterVars[v].(*int64) != test.expected {
							t.Errorf("%d-%s expected %v, got %v", i, test.name, test.expected, *appCfg.filterVars[v].(*int64))
						}
					}
				}
				break
			}
		}
	}
	// args not filtered shouldn't be affected
	i64, err := appCfg.GetInt64E("flagint64-tst")
	if err != nil {
		t.Errorf("Expected no error, got %s", err.Error())
	} else {
		if i64 != int64(42) {
			t.Errorf("Expected \"flagint64-tst\" to be 42, got %v", i64)
		}
	}

}

func TestGetStringFilter(t *testing.T) {
	tests := []struct {
		name     string
		short    string
		value    string
		expected string
		err      string
	}{
		{"", "", "a flag string", "a flag string", ""},
		{"corestring", "c", "a flag string", "a flag string", ""},
		{"string", "s", "a flag string", "a flag string", ""},
		{"cfgstring", "", "a flag string", "a flag string", ""},
		{"flagstring", "f", "xanadu", "xanadu", ""},
		{"cs", "", "a flag string", "a flag string", ""},
	}
	appCfg = newTestCfg()
	appCfg.SetName("rancher-test")
	args, err := appCfg.FilterArgs([]string{"-flagstring=xanadu", "fake", "command"})
	if err != nil {
		t.Errorf("expected no error, got %q", err.Error())
		return
	}
	if len(args) != 2 {
		t.Errorf("Expected 2 arg to be returned, got %d", len(args))
	}
	for i, test := range tests {
		for _, v := range appCfg.int64FilterNames {
			if v == test.name {
				if appCfg.filterVars[v] == nil {
					t.Errorf("%d: %s was nil", i, v)
				} else {
					switch appCfg.filterVars[v].(type) {
					case string:
						if appCfg.filterVars[v].(string) != test.expected {
							t.Errorf("%d expected %v, got %v", i, test.expected, appCfg.filterVars[v].(string))
						}
					case *string:
						if *appCfg.filterVars[v].(*string) != test.expected {
							t.Errorf("%d-%s expected %v, got %v", i, test.name, test.expected, *appCfg.filterVars[v].(*string))
						}
					}
				}
				break
			}
		}
	}
	// args not filtered shouldn't be affected
	s, err := appCfg.GetStringE("flagstring-tst")
	if err != nil {
		t.Errorf("Expected no error, got %s", err.Error())
	} else {
		if s != "a flag string" {
			t.Errorf("Expected \"flagstring-tst\" to be a flag string, got %v", s)
		}
	}

}
