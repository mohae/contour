package contour

import (
	"encoding/json"
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
		if s.UseConfFile() {
			t.Errorf("Expected false got %v", s.UseConfFile())
		}
		if s.UseEnvVars() {
			t.Errorf("Expected false got %v", s.UseEnvVars())
		}
		if s.UseFlags() {
			t.Errorf("Expected false got %v", s.UseFlags())
		}
	}
}

func TestLoadEnv(t *testing.T) {
	tests := []struct {
		sTyp        SettingType
		dTyp        dataType
		envValue    string
		name        string
		origValue   interface{}
		origString  string
		expValue    interface{}
		expectedErr string
	}{
		{ConfFileVar, _bool, "true", "cfgbool", false, "false", false, ""},
		{ConfFileVar, _int, "99", "cfgint", 42, "42", 42, ""},
		{ConfFileVar, _string, "bar", "cfgstring", "foo", "foo", "foo", ""},
		{EnvVar, _bool, "true", "envbool", false, "false", true, ""},
		{EnvVar, _int, "99", "envint", 42, "42", 99, ""},
		{EnvVar, _string, "bar", "envstring", "foo", "foo", "bar", ""},
		{Flag, _bool, "true", "flagbool", false, "false", true, ""},
		{Flag, _int, "88", "flagint", 42, "42", 88, ""},
		{Flag, _string, "biz", "flagstring", "fiz", "fiz", "biz", ""},
	}
	testCfg := New("contourtest")
	for _, test := range tests {
		switch test.sTyp {
		case ConfFileVar:
			testCfg.registerConfFileVar(test.dTyp, test.name, test.origValue, test.origString)
		case EnvVar:
			testCfg.registerEnvVar(test.dTyp, test.name, test.origValue, test.origString)
		case Flag:
			testCfg.registerFlag(test.dTyp, test.name, "", test.origValue, test.origString, "")
		}
		if test.sTyp == ConfFileVar { // ConfFileVars cannot be environment variables.
			continue
		}
		os.Setenv(strings.ToUpper(fmt.Sprintf("%s_%s", testCfg.Name(), test.name)), test.envValue)
	}
	testCfg.SetFromEnvVars()
	for _, test := range tests {
		tmp := testCfg.Get(test.name)
		switch test.dTyp {
		case _bool:
			if test.expValue != tmp.(bool) {
				t.Errorf("expected %v, got %v", test.expValue, tmp)
			}
		case _int:
			if test.expValue != tmp.(int) {
				t.Errorf("expected %v, got %v", test.expValue, tmp)
			}
		case _string:
			if test.expValue != tmp.(string) {
				t.Errorf("expected %v, got %v", test.expValue, tmp)
			}
		}
	}
}

func TestCheckPaths(t *testing.T) {
	fname := "abc.xyz"
	paths := []string{"aaaaa", "bbbbb", "ccccc"}
	cfg := New("test")
	b, err := cfg.checkPaths(fname, paths)
	if b != nil {
		t.Errorf("got %v; want nil", b)
	} else {
		if err != os.ErrNotExist {
			t.Errorf("got %s; want %s", err, os.ErrNotExist)
		}
	}
}

func TestConfFilePaths(t *testing.T) {
	tests := [][]string{
		{"testpath", "/test/path"},
		nil,
		{"path"},
	}
	cfg := New("test")
	for _, v := range tests {
		cfg.SetConfFilePaths(v)
		p := cfg.ConfFilePaths()
		if !reflect.DeepEqual(v, p) {
			t.Errorf("got %v; want %v", p, v)
		}
	}
}

func TestConfFileEnvVars(t *testing.T) {
	tests := [][]string{
		{"CONTOURPATH"},
		nil,
		{"CONTOURPATH", "ZPATH"},
	}
	cfg := New("test")
	for _, v := range tests {
		cfg.SetConfFilePathEnvVars(v)
		p := cfg.ConfFilePathEnvVars()
		if !reflect.DeepEqual(v, p) {
			t.Errorf("got %v; want %v", p, v)
		}
	}
}

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
		appCfg.useConfFile = test.useConfFile
		appCfg.confFileVarsSet = test.confFileVarsSet
		appCfg.envVarsSet = test.envVarsSet
		appCfg.useEnvVars = test.useEnvVars
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

// TestSet
func TestSet(t *testing.T) {
	tmpDir, err := ioutil.TempDir(os.TempDir(), "feedlot")
	if err != nil {
		t.Errorf("cannot do tests, %s", err.Error())
		return
	}
	// clean up on exit
	defer os.RemoveAll(tmpDir)
	fname := filepath.Join(tmpDir, "testcfg.json")
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
		{"", "", Unsupported, false, false, false, "", nil, "configuration filename: set failed: no name provided"},
		{"", "", Unsupported, false, true, false, "", nil, "configuration filename: set failed: no name provided"},
		{"", fname, JSON, false, false, false, "", nil, ""},
		{"", fname, JSON, true, false, false, "", nil, ""},
		{"", fname, JSON, false, true, false, "", nil, ""},
		// 5
		{"", fname, JSON, true, true, false, "", nil, ""},
		{"cfgstring", fname, JSON, false, false, false, "envstring", nil, ""},
		{"cfgstring", fname, JSON, true, false, false, "envstring", nil, ""},
		{"cfgstring", fname, JSON, false, true, false, "envstring", nil, ""},
		{"cfgstring", fname, JSON, true, true, false, "envstring", nil, ""},
		// 10
		{"cfgbool", fname, JSON, false, false, false, "true", nil, ""},
		{"cfgbool", fname, JSON, true, false, false, "true", nil, ""},
		{"cfgbool", fname, JSON, false, true, false, "true", nil, ""},
		{"cfgbool", fname, JSON, true, true, false, "true", nil, ""},
		{"cfgint", fname, JSON, false, false, false, "55", nil, ""},
		// 15
		{"cfgint", fname, JSON, true, false, false, "55", nil, ""},
		{"cfgint", fname, JSON, false, true, false, "55", nil, ""},
		{"cfgint", fname, JSON, true, true, false, "55", nil, ""},
		{"cfgint64", fname, JSON, false, false, false, "5564", nil, ""},
		{"cfgint64", fname, JSON, true, false, false, "5564", nil, ""},
		// 20
		{"cfgint64", fname, JSON, false, true, false, "5564", nil, ""},
		{"cfgint64", fname, JSON, true, true, false, "5564", nil, ""},
		{"flagstring", fname, JSON, false, false, false, "envstring", nil, ""},
		{"flagstring", fname, JSON, true, false, false, "envstring", nil, ""},
		{"flagstring", fname, JSON, false, true, false, "envstring", nil, ""},
		// 25
		{"flagstring", fname, JSON, true, true, false, "envstring", nil, ""},
		{"flagbool", fname, JSON, false, false, false, "true", nil, ""},
		{"flagbool", fname, JSON, true, false, false, "true", nil, ""},
		{"flagbool", fname, JSON, false, true, false, "true", nil, ""},
		{"flagbool", fname, JSON, true, true, false, "true", nil, ""},
		// 30
		{"flagint", fname, JSON, false, false, false, "22", nil, ""},
		{"flagint", fname, JSON, true, false, false, "22", nil, ""},
		{"flagint", fname, JSON, false, true, false, "22", nil, ""},
		{"flagint", fname, JSON, true, true, false, "22", nil, ""},
		{"flagint64", fname, JSON, false, false, false, "5564", nil, ""},
		// 35
		{"flagint64", fname, JSON, true, false, false, "5564", nil, ""},
		{"flagint64", fname, JSON, false, true, false, "5564", nil, ""},
		{"flagint64", fname, JSON, true, true, false, "5564", nil, ""},
		{"cfgnotthere", fname, JSON, false, false, false, "", nil, ""},
		{"cfgnotthere", fname, JSON, true, false, false, "", nil, ""},
		// 40
		{"cfgnotthere", fname, JSON, false, true, false, "", nil, ""},
		{"cfgnotthere", fname, JSON, true, true, false, "", nil, ""},
		{"flagnotthere", fname, JSON, false, false, false, "", nil, ""},
		{"flagnotthere", fname, JSON, true, false, false, "", nil, ""},
		{"flagnotthere", fname, JSON, false, true, false, "", nil, ""},
		// 45
		{"flagnotthere", fname, JSON, true, true, false, "", nil, ""},
		{"bool", fname, JSON, false, false, false, "true", nil, ""},
		{"bool", fname, JSON, true, false, false, "true", nil, ""},
		{"bool", fname, JSON, false, true, false, "true", nil, ""},
		{"bool", fname, JSON, true, true, false, "true", nil, ""},
		// 50
		{"int", fname, JSON, false, false, false, "33", nil, ""},
		{"int", fname, JSON, true, false, false, "33", nil, ""},
		{"int", fname, JSON, false, true, false, "33", nil, ""},
		{"int", fname, JSON, true, true, false, "33", nil, ""},
		{"int64", fname, JSON, false, false, false, "5564", nil, ""},
		// 55
		{"int64", fname, JSON, true, false, false, "5564", nil, ""},
		{"int64", fname, JSON, false, true, false, "5564", nil, ""},
		{"int64", fname, JSON, true, true, false, "5564", nil, ""},
		{"string", fname, JSON, false, false, false, "envstring", nil, ""},
		{"string", fname, JSON, true, false, false, "envstring", nil, ""},
		// 60
		{"string", fname, JSON, false, true, false, "envstring", nil, ""},
		{"string", fname, JSON, true, true, false, "envstring", nil, ""},
		{"corebool", fname, JSON, false, false, false, "true", nil, ""},
		{"corebool", fname, JSON, true, false, false, "true", nil, ""},
		{"corebool", fname, JSON, false, true, false, "true", nil, ""},
		// 65
		{"corebool", fname, JSON, true, true, false, "true", nil, ""},
		{"coreint", fname, JSON, false, false, false, "44", nil, ""},
		{"coreint", fname, JSON, true, false, false, "44", nil, ""},
		{"coreint", fname, JSON, false, true, false, "44", nil, ""},
		{"coreint", fname, JSON, true, true, false, "44", nil, ""},
		// 70
		{"coreint64", fname, JSON, false, false, false, "5564", nil, ""},
		{"coreint64", fname, JSON, true, false, false, "5564", nil, ""},
		{"coreint64", fname, JSON, false, true, false, "5564", nil, ""},
		{"coreint64", fname, JSON, true, true, false, "5564", nil, ""},
		{"string", fname, JSON, false, false, false, "envstring", nil, ""},
		// 75
		{"string", fname, JSON, true, false, false, "envstring", nil, ""},
		{"string", fname, JSON, false, true, false, "envstring", nil, ""},
		{"string", fname, JSON, true, true, false, "envstring", nil, ""},
	}
	// create temp file names
	// write the tmp json cfg file
	err = ioutil.WriteFile(fname, jsonTest, 0777)
	if err != nil {
		t.Errorf("cannot do tests, %s", err.Error())
		return
	}
	tstCfg := newTestSettings()
	tstCfg.name = "feedlot"
	for i, test := range tests {
		tstCfg.confFilename = ""
		tstCfg.confFileVarsSet = false
		tstCfg.envVarsSet = false
		err := tstCfg.SetConfFilename(test.fullPath)
		if err != nil {
			if err.Error() != test.err {
				t.Errorf("%d: setconffilename: got %q want %q", i, err, test.err)
			}
			continue
		}
		if test.err != "" {
			t.Errorf("%d: setconffilename: got no error; want %q", i, test.err)
			continue
		}
		tstCfg.useConfFile = test.useCfg
		tstCfg.useEnvVars = test.useEnvVars
		err = tstCfg.Set()
		if err != nil {
			if test.err != err.Error() {
				t.Errorf("%d: expected %q, got %q", i, test.err, err.Error())
			}
			continue
		}
		if test.err != "" {
			t.Errorf("%d expected %q, got nil", i, test.err)
			continue
		}
	}
}

func TestSetFromConfFile(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "contourTest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	pathVars := []string{"ABCPATH", "DEFPATH"}
	pathVarStr := "$ABCPATH; $DEFPATH; $PATH"
	badJSON := []byte(`{"var1": true, var2: false}`)
	goodJSON := []byte(`
{
	"var1": true,
	"var2": 42,
	"var3": "pan-galactic gargle blaster",
	"var4": [
		11,
		42
	],
	"var5": {
		"log": true
	}
}
`)
	tests := []struct {
		fname            string
		writeConfFile    bool
		data             []byte
		isSet            bool
		settings         map[string]setting
		expectedErr      error
		expectedSettings map[string]setting
	}{
		{"", false, nil, false, nil, error(&os.PathError{Op: "open file", Path: fmt.Sprintf("abcdef.json: %s", pathVarStr), Err: os.ErrNotExist}), nil},
		{"contour_test.yml", false, nil, false, nil, error(&os.PathError{Op: "open file", Path: fmt.Sprintf("%s: %s", filepath.Join(tmpDir, "contour_test.yml"), pathVarStr), Err: os.ErrNotExist}), nil},
		{"contour_test.json", false, nil, false, nil, error(&os.PathError{Op: "open file", Path: fmt.Sprintf("%s: %s", filepath.Join(tmpDir, "contour_test.json"), pathVarStr), Err: os.ErrNotExist}), nil},
		{
			"contour_test.json", true, badJSON, false, nil,
			fmt.Errorf("%s: invalid character 'v' looking for beginning of object key string", filepath.Join(tmpDir, "contour_test.json")),
			nil,
		},
		{"countour_test.json", true, goodJSON, true,
			map[string]setting{
				"var1": setting{Type: _bool, Name: "var1", IsConfFileVar: true},
				"var2": setting{Type: _int, Name: "var2", IsConfFileVar: true},
				"var3": setting{Type: _string, Name: "var3", IsConfFileVar: true},
				"var4": setting{Type: _interface, Name: "var4", IsConfFileVar: true},
				"var5": setting{Type: _interface, Name: "var5", IsConfFileVar: true},
				"var6": setting{Type: _int, Name: "var6", Value: interface{}(11), IsConfFileVar: false},
			},
			nil,
			map[string]setting{
				"var1": setting{Type: _bool, Name: "var1", Value: interface{}(true), IsConfFileVar: true},
				"var2": setting{Type: _int, Name: "var2", Value: interface{}(42), IsConfFileVar: true},
				"var3": setting{Type: _string, Name: "var3", Value: interface{}("pan-galactic gargle blaster"), IsConfFileVar: true},
				"var4": setting{Type: _interface, Name: "var4", Value: interface{}([]int{11, 42}), IsConfFileVar: true},
				"var5": setting{Type: _interface, Name: "var5", Value: interface{}(map[string]bool{"log": true}), IsConfFileVar: true},
				"var6": setting{Type: _int, Name: "var6", Value: interface{}(11), IsConfFileVar: false},
			},
		},
	}

	cfg := New("abcdef")
	cfg.confFilePathEnvVars = pathVars
	cfg.useConfFile = true
	cfg.confFileVarsSet = true
	cfg.settings = nil
	var fname string
	for i, test := range tests {
		cfg.confFileVarsSet = false
		if test.fname == "" {
			fname = ""
		} else {
			fname = filepath.Join(tmpDir, test.fname)
		}
		cfg.confFilename = fname
		cfg.settings = test.settings
		if test.writeConfFile {
			err = ioutil.WriteFile(fname, test.data, 0777)
			if err != nil {
				t.Fatalf("%d: write test conf file: %s", i, err)
			}
		}

		err = cfg.SetFromConfFile()
		if err != nil {
			if test.expectedErr == nil {
				t.Errorf("%d: expected no error; got %q", i, err)
			} else {
				if err.Error() != test.expectedErr.Error() {
					t.Errorf("%d: got %q; want %q", i, err.Error(), test.expectedErr.Error())
				}
			}
			continue
		}
		if test.expectedErr != nil {
			t.Errorf("%d: got no error; want %q", i, test.expectedErr)
			continue
		}
		got, _ := json.MarshalIndent(cfg.settings, "", "\t")
		want, _ := json.MarshalIndent(test.expectedSettings, "", "\t")
		if string(got) != string(want) {
			t.Errorf("%d: got\n%s\n; want\n%s", i, got, want)
			continue
		}
		// check set flags
		if cfg.confFileVarsSet != test.isSet {
			t.Errorf("%d: got %v; want %v", i, cfg.confFileVarsSet, test.isSet)
		}
	}
}

func TestSetFromConfFileMissingFile(t *testing.T) {
	tests := []struct {
		errOnMissingConfFile bool
		err                  error
	}{
		{false, nil},
		{true, &os.PathError{Op: "open file", Path: "does.not.exist.json: testpath; another/path; $ABCPATH; $DEFPATH; $PATH", Err: os.ErrNotExist}},
	}
	cfg := New("abcdef")
	cfg.confFilePaths = []string{"testpath", "another/path"}
	cfg.confFilePathEnvVars = []string{"ABCPATH", "DEFPATH"}
	cfg.confFilename = "does.not.exist.json"
	cfg.useConfFile = true
	for i, test := range tests {
		cfg.SetErrOnMissingConfFile(test.errOnMissingConfFile)
		err := cfg.SetFromConfFile()
		if err == nil {
			if test.err != nil {
				t.Errorf("%d: got nil; want %q", i, test.err)
			}
			continue
		}
		if !reflect.DeepEqual(err, test.err) {
			t.Errorf("%d: got %q; want %q", i, err, test.err)
		}
	}
}

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
	appCfg.AddInt("y", 11)
	appCfg.AddInt("z", 42)
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
	tst.useFlags = true
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

func TestGetEnvVarPaths(t *testing.T) {
	tests := []struct {
		value    string
		expected []string
	}{
		{"", []string{}},
		{"path", []string{"path"}},
		{strings.Join([]string{"path", filepath.Join("another", "path")}, string(os.PathListSeparator)), []string{"path", filepath.Join("another", "path")}},
		{strings.Join([]string{"path", filepath.Join("$TCONTOURPATHT", "another", "path")}, string(os.PathListSeparator)), []string{"path", filepath.Join("abc", "another", "path")}},
		{strings.Join([]string{"path", filepath.Join("another", "path"), filepath.Join("yellow", "brick", "road")}, string(os.PathListSeparator)), []string{"path", filepath.Join("another", "path"), filepath.Join("yellow", "brick", "road")}},
	}
	FeedlotPath := "TESTCONTOURPATH"
	os.Setenv("TCONTOURPATHT", "abc")
	// test empty key case first
	v := GetEnvVarPaths("")
	if v != nil {
		t.Errorf("got %#v; want nil", v)
	}
	for _, test := range tests {
		err := os.Setenv(FeedlotPath, test.value)
		if err != nil {
			t.Fatal(err)
		}
		p := GetEnvVarPaths(FeedlotPath)
		if !reflect.DeepEqual(p, test.expected) {
			t.Errorf("got %#v; want %#v", p, test.expected)
		}
	}
	os.Unsetenv("TCONTOURPATHT")
}
