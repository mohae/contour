package contour

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
		if s.UseEnv() != true {
			t.Errorf("Expected true got %v", s.UseEnv())
		}
	}
	if settings == nil {
		t.Errorf("global Settings was nil")
	} else {
		if settings.name != "app" {
			t.Errorf("Expected app got %s", settings.name)
		}
		if settings.UseEnv() != true {
			t.Errorf("Expected true got %v", settings.UseEnv())
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
			testCfg.RegisterBoolCfg(test.name, test.origValue.(bool))
		case "cfgint":
			testCfg.RegisterIntCfg(test.name, test.origValue.(int))
		case "cfgstring":
			testCfg.RegisterStringCfg(test.name, test.origValue.(string))
		case "flagbool":
			testCfg.RegisterBoolFlag(test.name, "", test.origValue.(bool), "", "")
		case "flagint":
			testCfg.RegisterIntFlag(test.name, "", test.origValue.(int), "", "")
		case "flagstring":
			testCfg.RegisterStringFlag(test.name, "", test.origValue.(string), "", "")
		}
		os.Setenv(strings.ToUpper(fmt.Sprintf("%s_%s", testCfg.Name(), test.name)), test.envValue)
	}
	testCfg.SetFromEnv()
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
		tstSettings.SetErrOnMissingFile(test.val)
		b := tstSettings.ErrOnMissingFile()
		if b != test.expected {
			t.Errorf("ErrOnMissingCfg:  expected %v, got %v", test.expected, b)
		}
		tstSettings.SetSearchPath(test.val)
		b = tstSettings.SearchPath()
		if b != test.expected {
			t.Errorf("SearchPath:  expected %v, got %v", test.expected, b)
		}
		tstSettings.SetUseCfg(test.val)
		b = tstSettings.UseCfg()
		if b != test.expected {
			t.Errorf("SetUseCfgFile:  expected %v, got %v", test.expected, b)
		}
		tstSettings.SetUseEnv(test.val)
		b = tstSettings.UseEnv()
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
		useCfg       bool
		cfgSet       bool
		useEnv       bool
		envSet       bool
		useFlags     bool
		argsFiltered bool
		expected     bool
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
		appCfg.SetUseCfg(test.useCfg)
		appCfg.cfgSet = test.cfgSet
		appCfg.envSet = test.envSet
		appCfg.SetUseEnv(test.useEnv)
		appCfg.useFlags = test.useFlags
		appCfg.argsFiltered = test.argsFiltered
		b := appCfg.CfgProcessed()
		if b != test.expected {
			t.Errorf("%d expected %v, got %v", i, test.expected, b)
		}
	}
}

func TestCanUpdate(t *testing.T) {
	tests := []struct {
		name         string
		argsFiltered bool
		expected     bool
		err          string
	}{
		// 0
		{"corebool", false, false, "corebool: core settings cannot be updated"},
		{"x-corebool", false, false, "x-corebool: setting not found"},
		{"coreint", false, false, "coreint: core settings cannot be updated"},
		{"x-coreint", false, false, "x-coreint: setting not found"},
		{"coreint64", false, false, "coreint64: core settings cannot be updated"},
		// 5
		{"x-coreint64", false, false, "x-coreint64: setting not found"},
		{"corestring", false, false, "corestring: core settings cannot be updated"},
		{"x-corestring", false, false, "x-corestring: setting not found"},
		{"cfgbool", false, true, ""},
		{"x-cfgbool", false, false, "x-cfgbool: setting not found"},
		// 10
		{"cfgint", false, true, ""},
		{"x-cfgint", false, false, "x-cfgint: setting not found"},
		{"cfgint64", false, true, ""},
		{"x-cfgint64", false, false, "x-cfgint64: setting not found"},
		{"cfgstring", false, true, ""},
		// 15
		{"x-cfgstring", false, false, "x-cfgstring: setting not found"},
		{"flagbool", false, true, ""},
		{"x-flagbool", false, false, "x-flagbool: setting not found"},
		{"flagint", false, true, ""},
		{"x-flagint", false, false, "x-flagint: setting not found"},
		// 20
		{"flagint64", false, true, ""},
		{"x-flagint64", false, false, "x-flagint64: setting not found"},
		{"flagstring", false, true, ""},
		{"x-flagstring", false, false, "x-flagstring: setting not found"},
		{"flagbool", true, false, "flagbool: flag settings cannot be updated after arg filtering"},
		// 25
		{"x-flagbool", true, false, "x-flagbool: setting not found"},
		{"flagint", true, false, "flagint: flag settings cannot be updated after arg filtering"},
		{"x-flagint", true, false, "x-flagint: setting not found"},
		{"flagint64", true, false, "flagint64: flag settings cannot be updated after arg filtering"},
		{"x-flagint64", true, false, "x-flagint64: setting not found"},
		// 30
		{"flagstring", true, false, "flagstring: flag settings cannot be updated after arg filtering"},
		{"x-flagstring", true, false, "x-flagstring: setting not found"},
		{"bool", false, true, ""},
		{"x-bool", false, false, "x-bool: setting not found"},
		{"int", false, true, ""},
		// 35
		{"x-int", false, false, "x-int: setting not found"},
		{"int64", false, true, ""},
		{"x-int64", false, false, "x-int64: setting not found"},
		{"string", false, true, ""},
		{"x-string", false, false, "x-string: setting not found"},
	}
	appCfg := newTestSettings()
	for i, test := range tests {
		appCfg.argsFiltered = test.argsFiltered
		b, err := appCfg.canUpdate(test.name)
		if err != nil {
			if err.Error() != test.err {
				t.Errorf("%d expected %q got %q", i, test.err, err.Error())
			}
			if b {
				t.Errorf("%d: expected returned value to be false on an error, it was not false", i)
			}
			continue
		}
		if b != test.expected {
			t.Errorf("%d: expected %v got %v", i, test.expected, b)
		}
	}
}

func TestCanOverride(t *testing.T) {
	tests := []struct {
		name         string
		argsFiltered bool
		expected     bool
	}{
		{"", false, false},
		{"", true, false},
		{"xyz", false, false},
		{"xyz", true, false},
		{"bool", false, true},
		{"bool", true, true},
		{"coreint", false, false},
		{"coreint", true, false},
		{"cfgint64", false, true},
		{"cfgint64", true, true},
		{"flagstring", false, true},
		{"flagstring", true, false},
	}
	appCfg := newTestSettings()
	for i, test := range tests {
		appCfg.argsFiltered = test.argsFiltered
		b := appCfg.canOverride(test.name)
		if b != test.expected {
			t.Errorf("%d: expected %v, got %v", i, test.expected, b)
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
		name   string
		IsCore bool
		IsCfg  bool
		IsEnv  bool
		IsFlag bool
		err    string
	}{
		{"", false, false, false, false, " setting not found"},
		{"string", false, false, false, false, ""},
		{"corebool", true, false, false, false, ""},
		{"cfgint", false, true, true, false, ""},
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
		// Cfg
		b, err = tstSettings.IsCfgE(test.name)
		if err != nil {
			if err.Error() != fmt.Sprintf(": file%s", test.err) {
				t.Errorf("%d: expected %q got %q", i, fmt.Sprintf(": file%s", test.err), err.Error())
			}
		} else {
			if b != test.IsCfg {
				t.Errorf("%d: expected %v, got %v", i, test.IsCfg, b)
			}
			b = tstSettings.IsCfg(test.name)
			if b != test.IsCfg {
				t.Errorf("%d: expected %v, got %v", i, test.IsCfg, b)
			}
		}
		// Env
		b, err = tstSettings.IsEnvE(test.name)
		if err != nil {
			if err.Error() != fmt.Sprintf(": env%s", test.err) {
				t.Errorf("%d: expected %q got %s", i, fmt.Sprintf(": env%s", test.err), err.Error())
			}
		} else {
			if b != test.IsEnv {
				t.Errorf("%d: expected %v, got %v", i, test.IsEnv, b)
			}
			b = tstSettings.IsEnv(test.name)
			if b != test.IsEnv {
				t.Errorf("%d: expected %v, got %v", i, test.IsEnv, b)
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
		name      string
		fullPath  string
		format    Format
		useCfg    bool
		useEnv    bool
		updateEnv bool
		envValue  string
		expected  interface{}
		err       string
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
	tstCfg.RegisterCfgFile("cfg_file", tests[5].fullPath)
	for i, test := range tests {
		tstCfg.UpdateCfgFile(test.fullPath)
		tstCfg.SetUseCfg(test.useCfg)
		tstCfg.SetUseEnv(test.useEnv)
		os.Setenv(GetEnvName(test.name), test.envValue)
		err := tstCfg.SetFromFile()
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
		ires, err := unmarshalCfgBytes(test.format, bites)
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
