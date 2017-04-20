package contour

import (
	"fmt"
	"reflect"
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
	tests := []struct {
		k       string
		v       interface{}
		visited bool
		t       dataType
	}{
		{"corebool", true, false, _bool},
		{"cfgbool", true, false, _bool},
		{"envbool", true, false, _bool},
		{"flagbool", false, true, _bool},
		{"flagbool-tst", false, true, _bool},
		{"coreint", 42, false, _int},
		{"cfgint", 42, false, _int},
		{"envint", 42, false, _int},
		{"flagint", 1999, true, _int},
		{"flagint-tst", 11, true, _int},
		{"coreint64", int64(42), false, _int64},
		{"cfgint64", int64(42), false, _int64},
		{"envint64", int64(42), false, _int64},
		{"flagint64", int64(42), false, _int64},
		{"flagint64-tst", int64(42), false, _int64},
		{"corestring", "a core string", false, _string},
		{"cfgstring", "a cfg string", false, _string},
		{"envstring", "an env string", false, _string},
		{"flagstring", "a flag string", false, _string},
		{"flagstring-tst", "updated", true, _string},
	}
	a, err := tst.ParseFlags(args)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if len(a) != 1 {
		t.Errorf("after parsing: got %d args; want 1", len(a))
	}

	if len(tst.parsedFlags) != 5 {
		t.Errorf("expected 5 flags to be parsed, %d were", len(tst.parsedFlags))
		return
	}
	visited := tst.Visited()
	// see if the parsed flags are tracked
TESTS:
	for _, test := range tests {
		for _, v := range visited {
			if v == test.k {
				if !test.visited {
					t.Errorf("%s was unexpectedly visited", v)
					continue TESTS
				}
				break
			}
		}
		switch test.t {
		case _bool:
			val, err := tst.BoolE(test.k)
			if err != nil {
				t.Errorf("%s: unexpected error getting bool: %s", test.k, err)
				continue
			}
			if val != test.v {
				t.Errorf("%s: got %v; want %v", test.k, val, test.v)
			}
		case _int:
			val, err := tst.IntE(test.k)
			if err != nil {
				t.Errorf("%s: unexpected error getting int: %s", test.k, err)
				continue
			}
			if val != test.v {
				t.Errorf("%s: got %v; want %v", test.k, val, test.v)
			}
		case _int64:
			val, err := tst.Int64E(test.k)
			if err != nil {
				t.Errorf("%s: unexpected error getting int64: %s", test.k, err)
				continue
			}
			if val != test.v {
				t.Errorf("%s: got %v; want %v", test.k, val, test.v)
			}
		case _string:
			val, err := tst.StringE(test.k)
			if err != nil {
				t.Errorf("%s: unexpected error getting string: %s", test.k, err)
				continue
			}
			if val != test.v {
				t.Errorf("%s: got %v; want %v", test.k, val, test.v)
			}
		default:
			t.Errorf("%s: unexpected type: %s", test.k, test.t)
		}

	}
	for _, v := range visited {
		var b bool
		for _, test := range tests {
			if test.k == v {
				b = true
				break
			}
		}
		if !b {
			t.Errorf("expected flag %q to be parsed; it wasn't", v)
		}
	}
}
