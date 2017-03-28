package contour

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
)

func TestSettings(t *testing.T) {
	s := New("test")
	if s == nil {
		t.Errorf("New test Settings was nil")
	} else {
		if s.name != "test" {
			t.Errorf("Expected test got %s", s.name)
		}
		if s.UseEnvVars() != false {
			t.Errorf("Expected true got %v", s.UseEnvVars())
		}
	}
	if settings == nil {
		t.Errorf("global Settings was nil")
	} else {
		if settings.name != "contour.test" {
			t.Errorf("Expected contour.test got %s", settings.name)
		}
		if settings.UseEnvVars() != false {
			t.Errorf("Expected true got %v", settings.UseEnvVars())
		}
	}
}

func TestLoadEnv(t *testing.T) {
	tests := []struct {
		name        string
		envValue    string
		typ         string
		origValue   interface{}
		expValue    interface{}
		expectedErr string
	}{
		{"tcfgbool", "true", "cfgbool", false, true, ""},
		{"tcfgint", "99", "cfgint", 42, 99, ""},
		{"tcfgstring", "bar", "cfgstring", "foo", "bar", ""},
		{"tflagbool", "true", "fkagbool", false, true, ""},
		{"tflagint", "88", "flagint", 42, 88, ""},
		{"tflagstring", "biz", "flagstring", "fiz", "biz", ""},
	}
	testCfg := New("contourtest")
	for _, test := range tests {
		switch test.typ {
		case "cfgbool":
			testCfg.RegisterBoolConfFileVar(test.name, test.origValue.(bool))
		case "cfgint":
			testCfg.RegisterIntConfFileVar(test.name, test.origValue.(int))
		case "cfgstring":
			testCfg.RegisterStringConfFileVar(test.name, test.origValue.(string))
		case "flagbool":
			testCfg.RegisterBoolFlag(test.name, "", test.origValue.(bool), "", "")
		case "flagint":
			testCfg.RegisterIntFlag(test.name, "", test.origValue.(int), "", "")
		case "flagstring":
			testCfg.RegisterStringFlag(test.name, "", test.origValue.(string), "", "")
		}
		os.Setenv(strings.ToUpper(fmt.Sprintf("%s_%s", testCfg.Name(), test.name)), test.envValue)
	}
	testCfg.SetFromEnvVars()
	for _, test := range tests {
		tmp := testCfg.Get(test.name)
		switch test.typ {
		case "cfgbool", "flagbool":
			if test.expValue != tmp.(bool) {
				t.Errorf("expected %v, got %v", test.expValue, tmp)
			}
		case "cfgint", "flagint":
			if test.expValue != tmp.(int) {
				t.Errorf("expected %v, got %v", test.expValue, tmp)
			}
		case "cfgstring", "flagstring":
			if test.expValue != tmp.(string) {
				t.Errorf("expected %v, got %v", test.expValue, tmp)
			}
		}
	}
}

func TestCfgBools(t *testing.T) {
	bTests := []struct {
		val      bool
		expected bool
	}{
		{true, true},
		{false, false},
		{true, true},
	}
	tstSettings := New("test")
	for _, test := range bTests {
		tstSettings.SetErrOnMissingConfFile(test.val)
		b := tstSettings.ErrOnMissingConfFile()
		if b != test.expected {
			t.Errorf("ErrOnMissingCfg:  expected %v, got %v", test.expected, b)
		}
		tstSettings.SetSearchPath(test.val)
		b = tstSettings.SearchPath()
		if b != test.expected {
			t.Errorf("SearchPath:  expected %v, got %v", test.expected, b)
		}
		tstSettings.SetUseConfFile(test.val)
		b = tstSettings.UseConfFile()
		if b != test.expected {
			t.Errorf("SetUseCfgFile:  expected %v, got %v", test.expected, b)
		}
		tstSettings.SetUseEnvVars(test.val)
		b = tstSettings.UseEnvVars()
		if b != test.expected {
			t.Errorf("SetUseEnv:  expected %v, got %v", test.expected, b)
		}
	}
}

//func TestSetFromCfg(t *testing.T) {
//
//}

func TestCfgProcessed(t *testing.T) {
	tests := []struct {
		useConfFile     bool
		confFileVarsSet bool
		useEnvVars      bool
		envVarsSet      bool
		useFlags        bool
		flagsParsed     bool
		expected        bool
	}{
		// 0
		{false, false, false, false, false, false, true},
		{false, false, false, false, false, true, true},
		{false, false, false, false, true, false, false},
		{false, false, false, false, true, true, true},
		{false, false, false, true, false, false, true},
		// 5
		{false, false, false, true, false, true, true},
		{false, false, false, true, true, false, false},
		{false, false, false, true, true, true, true},
		{false, false, true, false, false, false, false},
		{false, false, true, false, false, true, false},
		// 10
		{false, false, true, false, true, false, false},
		{false, false, true, false, true, true, false},
		{false, false, true, true, false, false, true},
		{false, false, true, true, true, false, false},
		{false, false, true, true, false, true, true},
		// 15
		{false, false, true, true, true, true, true},
		{false, true, false, false, false, false, true},
		{false, true, false, false, false, true, true},
		{false, true, false, false, true, false, false},
		{false, true, false, false, true, true, true},
		{false, true, false, true, false, false, true},
		// 20
		{false, true, false, true, false, true, true},
		{false, true, false, true, true, false, false},
		{false, true, false, true, true, true, true},
		{false, true, true, false, false, false, false},
		{false, true, true, false, false, true, false},
		// 25
		{false, true, true, false, true, false, false},
		{false, true, true, false, true, true, false},
		{false, true, true, true, false, false, true},
		{false, true, true, true, true, false, false},
		{false, true, true, true, false, true, true},
		// 30
		{false, true, true, true, true, true, true},
		{true, false, false, false, false, false, false},
		{true, false, false, false, false, true, false},
		{true, false, false, false, true, false, false},
		{true, false, false, false, true, true, false},
		// 35
		{true, false, false, true, false, false, false},
		{true, false, false, true, false, true, false},
		{true, false, false, true, true, false, false},
		{true, false, false, true, true, true, false},
		{true, false, true, false, false, false, false},
		// 40
		{true, false, true, false, false, true, false},
		{true, false, true, false, true, false, false},
		{true, false, true, false, true, true, false},
		{true, false, true, true, false, false, false},
		{true, false, true, true, true, false, false},
		// 45
		{true, false, true, true, false, true, false},
		{true, false, true, true, true, true, false},
		{true, true, false, false, false, false, true},
		{true, true, false, false, false, true, true},
		{true, true, false, false, true, false, false},
		// 50
		{true, true, false, false, true, true, true},
		{true, true, false, true, false, false, true},
		{true, true, false, true, false, true, true},
		{true, true, false, true, true, false, false},
		{true, true, false, true, true, true, true},
		// 55
		{true, true, true, false, false, false, false},
		{true, true, true, false, false, true, false},
		{true, true, true, false, true, false, false},
		{true, true, true, false, true, true, false},
		{true, true, true, true, false, false, true},
		// 60
		{true, true, true, true, true, false, false},
		{true, true, true, true, false, true, true},
		{true, true, true, true, true, true, true},
	}
	appCfg := New("test")
	for i, test := range tests {
		appCfg.SetUseConfFile(test.useConfFile)
		appCfg.confFileVarsSet = test.confFileVarsSet
		appCfg.envVarsSet = test.envVarsSet
		appCfg.SetUseEnvVars(test.useEnvVars)
		appCfg.useFlags = test.useFlags
		appCfg.flagsParsed = test.flagsParsed
		b := appCfg.IsSet()
		if b != test.expected {
			t.Errorf("%d expected %v, got %v", i, test.expected, b)
		}
	}
}

func TestSetUsage(t *testing.T) {
	f := func() { fmt.Println("hello world") }
	tstSettings := New("app")
	tstSettings.SetUsage(f)
	if tstSettings.flagSet.Usage == nil {
		t.Error("expected Cfg.flagSet.Usage to not be nil, it was nil")
	}
}

func TestIsFuncs(t *testing.T) {
	tests := []struct {
		name          string
		IsCore        bool
		IsConfFileVar bool
		IsEnvVar      bool
		IsFlag        bool
		err           string
	}{
		{"", false, false, false, false, " setting not found"},
		{"string", false, false, false, false, ""},
		{"corebool", true, false, false, false, ""},
		{"cfgint", false, true, false, false, ""},
		{"flagint64", false, true, true, true, ""},
	}
	tstSettings := newTestSettings()
	for i, test := range tests {
		// Core
		b, err := tstSettings.IsCoreE(test.name)
		if err != nil {
			if err.Error() != fmt.Sprintf(": core%s", test.err) {
				t.Errorf("%d: expected %q got %q", i, fmt.Sprintf(": core%s", test.err), err.Error())
			}
		} else {
			if b != test.IsCore {
				t.Errorf("%d: expected %v, got %v", i, test.IsCore, b)
			}
			b = tstSettings.IsCore(test.name)
			if b != test.IsCore {
				t.Errorf("%d: expected %v, got %v", i, test.IsCore, b)
			}
		}
		// ConfFileVars
		b, err = tstSettings.IsConfFileVarE(test.name)
		if err != nil {
			if err.Error() != fmt.Sprintf(": configuration file var%s", test.err) {
				t.Errorf("%d: expected %q got %q", i, fmt.Sprintf(": file%s", test.err), err.Error())
			}
		} else {
			if b != test.IsConfFileVar {
				t.Errorf("%d: expected %v, got %v", i, test.IsConfFileVar, b)
			}
			b = tstSettings.IsConfFileVar(test.name)
			if b != test.IsConfFileVar {
				t.Errorf("%d: expected %v, got %v", i, test.IsConfFileVar, b)
			}
		}
		// Env
		b, err = tstSettings.IsEnvVarE(test.name)
		if err != nil {
			if err.Error() != fmt.Sprintf(": env var%s", test.err) {
				t.Errorf("%d: expected %q got %s", i, fmt.Sprintf(": env var%s", test.err), err.Error())
			}
		} else {
			if b != test.IsEnvVar {
				t.Errorf("%d: expected %v, got %v", i, test.IsEnvVar, b)
			}
			b = tstSettings.IsEnvVar(test.name)
			if b != test.IsEnvVar {
				t.Errorf("%d: expected %v, got %v", i, test.IsEnvVar, b)
			}
		}
		// Flag
		b, err = tstSettings.IsFlagE(test.name)
		if err != nil {
			if err.Error() != fmt.Sprintf(": flag%s", test.err) {
				t.Errorf("%d: expected %q got %q", i, fmt.Sprintf(": flag%s", test.err), err.Error())
			}
		} else {
			if b != test.IsFlag {
				t.Errorf("%d: expected %v, got %v", i, test.IsFlag, b)
			}
			b = tstSettings.IsFlag(test.name)
			if b != test.IsFlag {
				t.Errorf("%d: expected %v, got %v", i, test.IsFlag, b)
			}
		}
	}
}

// TestSetCfg, and by proxy UpdateFromEnv
func TestSetCfg(t *testing.T) {
	tmpDir, err := ioutil.TempDir(os.TempDir(), "rancher")
	if err != nil {
		t.Errorf("cannot do tests, %s", err.Error())
		return
	}
	// clean up on exit
	defer os.RemoveAll(tmpDir)
	fname := "testcfg.json"
	tests := []struct {
		name          string
		fullPath      string
		format        Format
		useCfg        bool
		useEnvVars    bool
		updateEnvVars bool
		envValue      string
		expected      interface{}
		err           string
	}{
		// 0
		{"", "", Unsupported, false, false, false, "", nil, ""},
		{"", "", Unsupported, false, true, false, "", nil, ""},
		{"", filepath.Join(tmpDir, fname), JSON, false, false, false, "", nil, ""},
		{"", filepath.Join(tmpDir, fname), JSON, true, false, false, "", nil, ""},
		{"", filepath.Join(tmpDir, fname), JSON, false, true, false, "", nil, ""},
		// 5
		{"", filepath.Join(tmpDir, fname), JSON, true, true, false, "", nil, ""},
		{"cfgstring", filepath.Join(tmpDir, fname), JSON, false, false, false, "envstring", nil, ""},
		{"cfgstring", filepath.Join(tmpDir, fname), JSON, true, false, false, "envstring", nil, ""},
		{"cfgstring", filepath.Join(tmpDir, fname), JSON, false, true, false, "envstring", nil, ""},
		{"cfgstring", filepath.Join(tmpDir, fname), JSON, true, true, false, "envstring", nil, ""},
		// 10
		{"cfgbool", filepath.Join(tmpDir, fname), JSON, false, false, false, "true", nil, ""},
		{"cfgbool", filepath.Join(tmpDir, fname), JSON, true, false, false, "true", nil, ""},
		{"cfgbool", filepath.Join(tmpDir, fname), JSON, false, true, false, "true", nil, ""},
		{"cfgbool", filepath.Join(tmpDir, fname), JSON, true, true, false, "true", nil, ""},
		{"cfgint", filepath.Join(tmpDir, fname), JSON, false, false, false, "55", nil, ""},
		// 15
		{"cfgint", filepath.Join(tmpDir, fname), JSON, true, false, false, "55", nil, ""},
		{"cfgint", filepath.Join(tmpDir, fname), JSON, false, true, false, "55", nil, ""},
		{"cfgint", filepath.Join(tmpDir, fname), JSON, true, true, false, "55", nil, ""},
		{"cfgint64", filepath.Join(tmpDir, fname), JSON, false, false, false, "5564", nil, ""},
		{"cfgint64", filepath.Join(tmpDir, fname), JSON, true, false, false, "5564", nil, ""},
		// 20
		{"cfgint64", filepath.Join(tmpDir, fname), JSON, false, true, false, "5564", nil, ""},
		{"cfgint64", filepath.Join(tmpDir, fname), JSON, true, true, false, "5564", nil, ""},
		{"flagstring", filepath.Join(tmpDir, fname), JSON, false, false, false, "envstring", nil, ""},
		{"flagstring", filepath.Join(tmpDir, fname), JSON, true, false, false, "envstring", nil, ""},
		{"flagstring", filepath.Join(tmpDir, fname), JSON, false, true, false, "envstring", nil, ""},
		// 25
		{"flagstring", filepath.Join(tmpDir, fname), JSON, true, true, false, "envstring", nil, ""},
		{"flagbool", filepath.Join(tmpDir, fname), JSON, false, false, false, "true", nil, ""},
		{"flagbool", filepath.Join(tmpDir, fname), JSON, true, false, false, "true", nil, ""},
		{"flagbool", filepath.Join(tmpDir, fname), JSON, false, true, false, "true", nil, ""},
		{"flagbool", filepath.Join(tmpDir, fname), JSON, true, true, false, "true", nil, ""},
		// 30
		{"flagint", filepath.Join(tmpDir, fname), JSON, false, false, false, "22", nil, ""},
		{"flagint", filepath.Join(tmpDir, fname), JSON, true, false, false, "22", nil, ""},
		{"flagint", filepath.Join(tmpDir, fname), JSON, false, true, false, "22", nil, ""},
		{"flagint", filepath.Join(tmpDir, fname), JSON, true, true, false, "22", nil, ""},
		{"flagint64", filepath.Join(tmpDir, fname), JSON, false, false, false, "5564", nil, ""},
		// 35
		{"flagint64", filepath.Join(tmpDir, fname), JSON, true, false, false, "5564", nil, ""},
		{"flagint64", filepath.Join(tmpDir, fname), JSON, false, true, false, "5564", nil, ""},
		{"flagint64", filepath.Join(tmpDir, fname), JSON, true, true, false, "5564", nil, ""},
		{"cfgnotthere", filepath.Join(tmpDir, fname), JSON, false, false, false, "", nil, ""},
		{"cfgnotthere", filepath.Join(tmpDir, fname), JSON, true, false, false, "", nil, ""},
		// 40
		{"cfgnotthere", filepath.Join(tmpDir, fname), JSON, false, true, false, "", nil, ""},
		{"cfgnotthere", filepath.Join(tmpDir, fname), JSON, true, true, false, "", nil, ""},
		{"flagnotthere", filepath.Join(tmpDir, fname), JSON, false, false, false, "", nil, ""},
		{"flagnotthere", filepath.Join(tmpDir, fname), JSON, true, false, false, "", nil, ""},
		{"flagnotthere", filepath.Join(tmpDir, fname), JSON, false, true, false, "", nil, ""},
		// 45
		{"flagnotthere", filepath.Join(tmpDir, fname), JSON, true, true, false, "", nil, ""},
		{"bool", filepath.Join(tmpDir, fname), JSON, false, false, false, "true", nil, ""},
		{"bool", filepath.Join(tmpDir, fname), JSON, true, false, false, "true", nil, ""},
		{"bool", filepath.Join(tmpDir, fname), JSON, false, true, false, "true", nil, ""},
		{"bool", filepath.Join(tmpDir, fname), JSON, true, true, false, "true", nil, ""},
		// 50
		{"int", filepath.Join(tmpDir, fname), JSON, false, false, false, "33", nil, ""},
		{"int", filepath.Join(tmpDir, fname), JSON, true, false, false, "33", nil, ""},
		{"int", filepath.Join(tmpDir, fname), JSON, false, true, false, "33", nil, ""},
		{"int", filepath.Join(tmpDir, fname), JSON, true, true, false, "33", nil, ""},
		{"int64", filepath.Join(tmpDir, fname), JSON, false, false, false, "5564", nil, ""},
		// 55
		{"int64", filepath.Join(tmpDir, fname), JSON, true, false, false, "5564", nil, ""},
		{"int64", filepath.Join(tmpDir, fname), JSON, false, true, false, "5564", nil, ""},
		{"int64", filepath.Join(tmpDir, fname), JSON, true, true, false, "5564", nil, ""},
		{"string", filepath.Join(tmpDir, fname), JSON, false, false, false, "envstring", nil, ""},
		{"string", filepath.Join(tmpDir, fname), JSON, true, false, false, "envstring", nil, ""},
		// 60
		{"string", filepath.Join(tmpDir, fname), JSON, false, true, false, "envstring", nil, ""},
		{"string", filepath.Join(tmpDir, fname), JSON, true, true, false, "envstring", nil, ""},
		{"corebool", filepath.Join(tmpDir, fname), JSON, false, false, false, "true", nil, ""},
		{"corebool", filepath.Join(tmpDir, fname), JSON, true, false, false, "true", nil, ""},
		{"corebool", filepath.Join(tmpDir, fname), JSON, false, true, false, "true", nil, ""},
		// 65
		{"corebool", filepath.Join(tmpDir, fname), JSON, true, true, false, "true", nil, ""},
		{"coreint", filepath.Join(tmpDir, fname), JSON, false, false, false, "44", nil, ""},
		{"coreint", filepath.Join(tmpDir, fname), JSON, true, false, false, "44", nil, ""},
		{"coreint", filepath.Join(tmpDir, fname), JSON, false, true, false, "44", nil, ""},
		{"coreint", filepath.Join(tmpDir, fname), JSON, true, true, false, "44", nil, ""},
		// 70
		{"coreint64", filepath.Join(tmpDir, fname), JSON, false, false, false, "5564", nil, ""},
		{"coreint64", filepath.Join(tmpDir, fname), JSON, true, false, false, "5564", nil, ""},
		{"coreint64", filepath.Join(tmpDir, fname), JSON, false, true, false, "5564", nil, ""},
		{"coreint64", filepath.Join(tmpDir, fname), JSON, true, true, false, "5564", nil, ""},
		{"string", filepath.Join(tmpDir, fname), JSON, false, false, false, "envstring", nil, ""},
		// 75
		{"string", filepath.Join(tmpDir, fname), JSON, true, false, false, "envstring", nil, ""},
		{"string", filepath.Join(tmpDir, fname), JSON, false, true, false, "envstring", nil, ""},
		{"string", filepath.Join(tmpDir, fname), JSON, true, true, false, "envstring", nil, ""},
	}
	// create temp file names
	// write the tmp json cfg file
	err = ioutil.WriteFile(filepath.Join(tmpDir, fname), jsonTest, 0777)
	if err != nil {
		t.Errorf("cannot do tests, %s", err.Error())
		return
	}
	tstCfg := newTestSettings()
	tstCfg.name = "rancher"
	tstCfg.SetConfFilename(tests[5].fullPath)
	for i, test := range tests {
		tstCfg.SetUseConfFile(test.useCfg)
		tstCfg.SetUseEnvVars(test.useEnvVars)
		os.Setenv(EnvVarName(test.name), test.envValue)
		err := tstCfg.SetFromConfFile()
		if err != nil {
			if test.err != err.Error() {
				t.Errorf("%d: expected %s, got %s", i, test.err, err.Error())
			}
			continue
		}
		if test.err != "" {
			t.Errorf("%d expected %s, got nil", i, test.err)
			continue
		}
	}
}

// Testing
func TestFormatFromFilename(t *testing.T) {
	tests := []basic{
		{"an empty cfgfilename", 0, "", "", "no configuration filename"},
		{"a cfgfilename without an extension", 0, "cfg", "", "unable to determine cfg's format: no extension"},
		{"a cfgfilename with an invalid extension", 0, "cfg.bmp", "", "bmp: unsupported configuration format"},
		{"a cfgfilename with a json extension", 0, "cfg.json", "json", ""},
		{"a path and multi dot cfgfilename with a json extension", 0, "path/to/custom.cfg.json", "json", ""},
		{"a cfgfilename with a toml extension", 0, "cfg.toml", "toml", ""},
		{"a cfgfilename with a toml extension", 0, "cfg.yaml", "yaml", ""},
		{"a cfgfilename with a toml extension", 0, "cfg.yml", "yaml", ""},
		{"a cfgfilename with a toml extension", 0, "cfg.xml", "", "xml: unsupported configuration format"},
		{"a cfgfilename with a toml extension", 0, "cfg.ini", "", "ini: unsupported configuration format"},
	}
	for _, test := range tests {
		format, err := formatFromFilename(test.value)
		if err != nil {
			if err.Error() != test.expectedErr {
				t.Errorf("%s: expected %s, got %s", test.name, test.expectedErr, err.Error())
			}
			continue
		}
		if test.expectedErr != "" {
			t.Errorf("%s: expected %s, got no error", test.name, test.expectedErr)
		}
		if format.String() != test.expected {
			t.Errorf("%s: expected %s, got %s", test.name, test.expected, format.String())
		}
	}
}

func TestIsSupportedFormat(t *testing.T) {
	tests := []basic{
		{"empty format test", 0, "", "false", ""},
		{"invalid format test", 0, "bmp", "false", ""},
		{"json format test", 0, "json", "true", ""},
		{"tom format testl", 0, "toml", "true", ""},
		{"yaml format test", 0, "yaml", "true", ""},
		{"yml format test", 0, "yml", "true", ""},
		{"xml format test", 0, "xml", "false", ""},
	}
	for i, test := range tests {
		// we don't care about error on this, only the supported part
		f, _ := ParseFormat(test.value)
		is := f.isSupported()
		if strconv.FormatBool(is) != test.expected {
			t.Errorf("%d: expected %v, got %v", i, test.expected, is)
		}
	}
}

func TestUnmarshalCfgBytes(t *testing.T) {
	tests := []struct {
		name        string
		format      Format
		value       []byte
		expected    interface{}
		expectedErr string
	}{
		{"json cfg", JSON, jsonExample, jsonResults, ""},
		{"toml cfg", TOML, tomlExample, tomlResults, ""},
		{"yaml cfg", YAML, yamlExample, yamlResults, ""},
		{"unsupported cfg", Unsupported, []byte(""), []byte(""), "unsupported: unsupported configuration format"},
	}

	for _, test := range tests {
		bites := []byte(test.value)
		ires, err := unmarshalConfBytes(test.format, bites)
		if err != nil {
			if test.expectedErr == "" {
				t.Errorf("%s: expected nil for error; got %q", test.name, err)
				continue
			}
			if err.Error() != test.expectedErr {
				t.Errorf("%s: expected %q; got %q", test.name, test.expectedErr, err)
			}
			continue
		}
		if test.format == YAML {
			val, ok := ires.(map[interface{}]interface{})["appVar1"]
			if !ok {
				t.Errorf("appVar1 not found")
			} else {
				if val != test.expected.(map[interface{}]interface{})["appVar1"] {
					t.Errorf("appVar1: expected %v, got %v", test.expected.(map[interface{}]interface{})["appVar1"], val)
				}
			}
			val, ok = ires.(map[interface{}]interface{})["appVar2"]
			if !ok {
				t.Errorf("appVar2 not found")
			} else {
				if val != test.expected.(map[interface{}]interface{})["appVar2"] {
					t.Errorf("appVar2: expected %v, got %v", test.expected.(map[interface{}]interface{})["appVar2"], val)
				}
			}
			val, ok = ires.(map[interface{}]interface{})["appVar3"]
			if !ok {
				t.Errorf("appVar3 not found")
			} else {
				if val != test.expected.(map[interface{}]interface{})["appVar3"] {
					t.Errorf("appVar3: expected %v, got %v", test.expected.(map[interface{}]interface{})["appVar3"], val)
				}
			}
			val, ok = ires.(map[interface{}]interface{})["appVar4"]
			if !ok {
				t.Errorf("appVar4 not found")
			} else {
				if val != test.expected.(map[interface{}]interface{})["appVar4"] {
					t.Errorf("appVar4: expected %v, got %v", test.expected.(map[interface{}]interface{})["appVar4"], val)
				}
			}
			val, ok = ires.(map[interface{}]interface{})["appVar5"]
			if !ok {
				t.Errorf("appVar5 not found")
			}
			continue
		}
		val, ok := ires.(map[string]interface{})["appVar1"]
		if !ok {
			t.Errorf("appVar1 not found")
		} else {
			if val != test.expected.(map[string]interface{})["appVar1"] {
				t.Errorf("appVar1: expected %v, got %v", test.expected.(map[string]interface{})["appVar1"], val)
			}
		}
		val, ok = ires.(map[string]interface{})["appVar2"]
		if !ok {
			t.Errorf("appVar2 not found")
		} else {
			if val != test.expected.(map[string]interface{})["appVar2"] {
				t.Errorf("appVar2: expected %v, got %v", test.expected.(map[string]interface{})["appVar2"], val)
			}
		}
		val, ok = ires.(map[string]interface{})["appVar3"]
		if !ok {
			t.Errorf("appVar3 not found")
		} else {
			if val != test.expected.(map[string]interface{})["appVar3"] {
				t.Errorf("appVar3: expected %v, got %v", test.expected.(map[string]interface{})["appVar3"], val)
			}
		}
		val, ok = ires.(map[string]interface{})["appVar4"]
		if !ok {
			t.Errorf("appVar4 not found")
		} else {
			if val != test.expected.(map[string]interface{})["appVar4"] {
				t.Errorf("appVar4: expected %v, got %v", test.expected.(map[string]interface{})["appVar4"], val)
			}
		}
		val, ok = ires.(map[string]interface{})["appVar5"]
		if !ok {
			t.Errorf("appVar5 not found")
		}
	}
}

func TestExists(t *testing.T) {
	tests := []struct {
		k      string
		exists bool
	}{
		{"", false},
		{"x", false},
		{"y", true},
		{"z", true},
	}
	appCfg := newTestSettings()
	appCfg.RegisterInt("y", 11)
	appCfg.RegisterInt("z", 42)
	for i, test := range tests {
		exists := appCfg.Exists(test.k)
		if exists != test.exists {
			t.Errorf("%d: got %v, want %v", i, exists, test.exists)
		}
	}
}

func TestVisited(t *testing.T) {
	tst := newTestSettings()
	args := []string{"-b=false", "-i=1999", "-flagbool-tst=false", "-flagint-tst=11", "-flagstring-tst=updated", "cmd"}
	expected := []string{"flagbool", "flagint", "flagbool-tst", "flagint-tst", "flagstring-tst"}
	tests := []struct {
		name    string
		visited bool
	}{
		{"", false},
		{"flagstring", false},
		{"flagbool", true},
		{"flagint", true},
		{"flagbool-tst", true},
		{"flagint-tst", true},
		{"flagstring-tst", true},
	}

	_, err := tst.ParseFlags(args)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	for _, test := range tests {
		visited := tst.WasVisited(test.name)
		if visited != test.visited {
			t.Errorf("%s: got %t; want %t", test.name, visited, test.visited)
		}
	}
	sort.Strings(expected)
	if !reflect.DeepEqual(tst.Visited(), expected) {
		t.Errorf("visited: got %v; want %v", tst.Visited(), expected)
	}
}
