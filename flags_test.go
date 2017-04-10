package contour

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestGetBoolFilter(t *testing.T) {
	tests := []struct {
		name     string
		short    string
		value    bool
		expected bool
		err      string
	}{
		{"", "", false, false, ""},
		{"corebool", "", true, true, ""},
		{"bool", "", false, false, ""},
		{"cfgbool", "", true, true, ""},
		{"flagbool", "b", false, false, ""},
		{"cb", "", false, false, ""},
	}
	_ = tests
	appCfg := newTestSettings()
	appCfg.useFlags = true
	// check that an unfiltered boolflag's var is nil
	if appCfg.flagVars["flagbool-tst"] != nil {
		t.Errorf("flags that weren't present in the args should be nil, flagbool-test wasn't")
	}
	args, err := appCfg.ParseFlags([]string{"cmdname", "-flagbool=false", "-fake", "command"})
	if err != nil {
		if !strings.HasPrefix(err.Error(), "parse of command-line arguments failed: flag provided but not defined: -fake") {
			t.Errorf("expected \" parse of command-line arguments failed: flag provided but not defined: -fake\", got %q", err.Error())
		}
	}
	appCfg = newTestSettings()
	appCfg.useFlags = true
	args, err = appCfg.ParseFlags([]string{"-flagbool=false", "fake", "command"})
	if err != nil {
		t.Errorf("expected no error, got %q", err.Error())
		return
	}
	if len(args) != 2 {
		t.Errorf("Expected 2 arg to be returned, got %d", len(args))
	}
	for i, test := range tests {
		for _, v := range appCfg.settings {
			if v.Name == test.name && v.IsFlag {
				if appCfg.flagVars[v.Name] == nil {
					t.Errorf("%d: flagVar %s was nil", i, v.Name)
				} else {
					switch appCfg.flagVars[v.Name].(type) {
					case bool:
						if appCfg.flagVars[v.Name].(bool) != test.expected {
							t.Errorf("%d:%s expected %v, got %v", i, v.Name, test.expected, appCfg.flagVars[v.Name].(bool))
						}
					case *bool:
						if *appCfg.flagVars[v.Name].(*bool) != test.expected {
							t.Errorf("%d-%s:%s: expected %v, got %v", i, test.name, v.Name, test.expected, *appCfg.flagVars[v.Name].(*bool))
						}
					default:
						t.Errorf("%d-%s:%s: is type %d wanted bool or *bool", i, test.name, v.Name, reflect.TypeOf(v.Value))
					}
				}
				break
			}
		}
	}

	// args not filtered shouldn't be affected
	b, err := appCfg.BoolE("flagbool-tst")
	if err != nil {
		t.Errorf("Expected no error, got %s", err.Error())
	} else {
		if !b {
			t.Errorf("Expected \"flagbool-tst\" to be true, got %v", b)
		}
	}
}

func TestParseFlags(t *testing.T) {
	tst := newTestSettings()
	tst.useFlags = true
	args := []string{"-b=false", "-i=1999", "-flagbool-tst=false", "-flagint-tst=11", "-flagstring-tst=updated", "cmd"}
	expected := map[string]string{
		"b":              "false",
		"i":              "1999",
		"flagbool-tst":   "false",
		"flagint-tst":    "11",
		"flagstring-tst": "updated",
	}
	vals, err := tst.ParseFlags(args)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if len(vals) != 1 {
		t.Errorf("expected 1 arg to be returned, got %d: %v", len(vals), vals)
		t.Errorf("visited: %v", tst.Visited())
		return
	}
	if vals[0] != "cmd" {
		t.Errorf("vals: got %s; want cmd", vals[0])
		return
	}
	for k, v := range expected {
		s, ok := tst.settings[k]
		if !ok {
			s = tst.settings[tst.shortFlags[k]]
		}
		if fmt.Sprintf("%s", s.Value) != v {
			t.Errorf("%s: got %v; want %v", k, s.Value, v)
		}
	}
}

func TestParsedFlags(t *testing.T) {
	tst := newTestSettings()
	tst.useFlags = true
	args := []string{"-b=false", "-i=1999", "-flagbool-tst=false", "-flagint-tst=11", "-flagstring-tst=updated", "cmd"}
	expected := []string{"flagbool", "flagint", "flagbool-tst", "flagint-tst", "flagstring-tst"}

	_, err := tst.ParseFlags(args)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	var checked []int
	if len(tst.parsedFlags) != len(expected) {
		t.Errorf("expected %d flags to be parsed, %d were", len(expected), len(tst.parsedFlags))
		return
	}
	// see if the parsed flags are tracked
	for _, v := range tst.parsedFlags {
		var found bool
		for i, x := range expected {
			if v == x {
				found = true
				checked = append(checked, i)
				break
			}
		}
		if !found {
			t.Errorf("%s was parsed but not expected", v)
			return
		}
	}
	sort.Ints(checked)
	for i, v := range checked {
		if i != v {
			t.Errorf("expected %s to be one of the parsed flags, it wasn't", expected[v])
			return
		}
	}
}
