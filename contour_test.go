package contour

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	rand.Seed(int64(time.Now().Nanosecond()))
	std = newTestSettings()
	std.name = "contour.test"
}

type basic struct {
	name        string
	settingType SettingType
	value       string
	expected    string
	expectedErr string
}

var tomlExample = []byte(`
appVar1 = true
appVar2 = false
appVar3 = "42"
appVar4 = "zip"
appVar5 = [
	"less",
	"sass",
	"scss"
]

[logging]
Logging = true
LogCfg = "test/test.toml"
LogFileLevel = "debug"
LogStdoutLevel = "error"
`)

var tomlResults = map[string]interface{}{
	"appVar1": true,
	"appVar2": false,
	"appVar3": "42",
	"appVar4": "zip",
	"appVar5": []string{"less", "sass", "scss"},
	"logging": map[string]interface{}{
		"Logging":        true,
		"LogCfg":         "test/test.toml",
		"LogFileLevel":   "debug",
		"LogStdoutLevel": "error",
	},
}

var jsonExample = []byte(`
{
	"appVar1": true,
	"appVar2": false,
	"appVar3": "42",
	"appVar4": "zip",
	"appVar5": [
		"less",
		"sass",
		"scss"
	],
	"logging": {
		"logging": true,
		"logcfg": "test/test.toml",
		"logfilelevel": "debug",
		"logstdoutlevel": "error"
	}
}
`)

var jsonResults = map[string]interface{}{
	"appVar1": true,
	"appVar2": false,
	"appVar3": "42",
	"appVar4": "zip",
	"appVar5": []string{"less", "sass", "scss"},
	"logging": map[string]interface{}{
		"logging":        true,
		"logcfg":         "test/test.toml",
		"logfilelevel":   "debug",
		"logstdoutlevel": "error",
	},
}

var jsonTest = []byte(`
{
	"cfgbool": true,
	"flagbool": true,
	"cfgint": 42,
	"flagint": 1999,
	"cfgstring": "foo",
	"flagstring": "bar",
	"cfgslice": [
		"mysql",
		"pgsql"
	],
	"flagslice": [
		"less",
		"sass",
		"scss"
	],
	"cfgmap": {
		"faz": 42,
		"fiz": true,
		"fuz": "buz"
	},
	"flagmap": {
		"log": true,
		"logcfg": "test/test.toml",
		"logfilelevel": "debug",
		"logstdoutlevel": "error"
	}
}
`)

var jsonTestResults = map[string]interface{}{
	"cfgbool":    true,
	"flagbool":   true,
	"cfgint":     42,
	"flagint":    1999,
	"cfgstring":  "foo",
	"flagstring": "bar",
	"cfgslice":   []string{"mysql", "pgsql"},
	"flagslice":  []string{"less", "sass", "scss"},
	"cfgmap": map[string]interface{}{
		"faz": 41,
		"fiz": true,
		"fuz": "buz",
	},
	"flagmap": map[string]interface{}{
		"log":            true,
		"logcfg":         "test/test.toml",
		"logfilelevel":   "debug",
		"logstdoutlevel": "error",
	},
}

var yamlExample = []byte(`appVar1: true
appVar2: false
appVar3: 42
appVar4: zip
appVar5:
  - less
  - iass
  - scss

logging:
  - Logging: true
  - LogCfg: test/test.yaml
  - LogFileLevel: debug
  - LogStdoutLevel: error
`)

var yamlResults = map[interface{}]interface{}{
	"appVar1": true,
	"appVar2": false,
	"appVar3": 42,
	"appVar4": "zip",
	"appVar5": []string{"less", "sass", "scss"},
	"logging": map[interface{}]interface{}{
		"logging":        true,
		"logcfg":         "test/test.toml",
		"logfilelevel":   "debug",
		"logstdoutlevel": "error",
	},
}

var xmlExample = []byte(`<cfg>
	<appVar1>true</appVar1>
	<appVar2>false</appVar2>
	<appVar3>42</appVar3>
	<appVar4>zip</appVar4>
	<appVar5>less</appVar5>
	<appVar5>sass</appVar5>
	<appVar5>scss</appVar5>
	<logging>
		<logging>true</logging>
		<logcfg>test/test.toml</logcfg>
		<logfilelevel>debug</logfilelevel>
		<logstdoutlevel>error</logstdoutlevel>
	</logging>
</cfg>
`)

var emptySettings map[string]Settings
var testSettings = map[string]Settings{
	Exe:     Settings{settings: map[string]setting{}},
	"test1": Settings{settings: map[string]setting{}},
}

func newTestSettings() *Settings {
	s := New("contour.test")
	s.settings = map[string]setting{
		"corebool": setting{
			Type:   _bool,
			Name:   "corebool",
			Value:  true,
			IsCore: true,
		},
		"coreint": setting{
			Type:   _int,
			Name:   "coreint",
			Value:  42,
			IsCore: true,
		},
		"coreint64": setting{
			Type:   _int64,
			Name:   "coreint64",
			Value:  int64(42),
			IsCore: true,
		},
		"corestring": setting{
			Type:   _string,
			Name:   "corestring",
			Value:  "a core string",
			IsCore: true,
		},
		"coreslice": setting{
			// Type:   "string-slice",
			Type:   _interface,
			Name:   "coreslice",
			Value:  []string{},
			IsCore: true,
		},
		"coremap": setting{
			//Type:   "map",
			Type:   _interface,
			Name:   "coremap",
			Value:  map[string]interface{}{},
			IsCore: true,
		},
		"cfgbool": setting{
			Type:          _bool,
			Name:          "cfgbool",
			Value:         true,
			Short:         "",
			IsConfFileVar: true,
		},
		"cfgint": setting{
			Type:          _int,
			Name:          "cfgint",
			Value:         42,
			IsConfFileVar: true,
		},
		"cfgint64": setting{
			Type:          _int64,
			Name:          "cfgint64",
			Value:         int64(42),
			IsConfFileVar: true,
		},
		"cfgstring": setting{
			Type:          _string,
			Name:          "cfgstring",
			Value:         "a cfg string",
			Short:         "",
			IsConfFileVar: true,
		},
		"cfgslice": setting{
			//Type:  "string-slice",
			Type:          _interface,
			Name:          "cfgslice",
			Value:         []string{},
			Short:         "",
			IsConfFileVar: true,
		},
		"cfgmap": setting{
			//Type:  "map",
			Type:          _interface,
			Name:          "cfgmap",
			Value:         map[string]interface{}{},
			Short:         "",
			IsConfFileVar: true,
		},
		"envbool": setting{
			Type:          _bool,
			Name:          "envbool",
			Value:         true,
			Short:         "",
			IsConfFileVar: true,
			IsEnvVar:      true,
		},
		"envint": setting{
			Type:          _int,
			Name:          "envint",
			Value:         42,
			IsConfFileVar: true,
			IsEnvVar:      true,
		},
		"envint64": setting{
			Type:          _int64,
			Name:          "envint64",
			Value:         int64(42),
			IsConfFileVar: true,
			IsEnvVar:      true,
		},
		"envstring": setting{
			Type:          _string,
			Name:          "envstring",
			Value:         "an env string",
			Short:         "",
			IsConfFileVar: true,
			IsEnvVar:      true,
		},
		"envslice": setting{
			//Type:  "string-slice",
			Type:          _interface,
			Name:          "envslice",
			Value:         []string{},
			Short:         "",
			IsConfFileVar: true,
			IsEnvVar:      true,
		},
		"envmap": setting{
			//Type:  "map",
			Type:          _interface,
			Name:          "envmap",
			Value:         map[string]interface{}{},
			Short:         "",
			IsConfFileVar: true,
			IsEnvVar:      true,
		},
		"flagbool": setting{
			Type:          _bool,
			Name:          "flagbool",
			Value:         true,
			Short:         "b",
			IsConfFileVar: true,
			IsEnvVar:      true,
			IsFlag:        true,
		},
		"flagbool-tst": setting{
			Type:          _bool,
			Name:          "flagbool-tst",
			Value:         true,
			Short:         "",
			IsConfFileVar: true,
			IsEnvVar:      true,
			IsFlag:        true,
		},
		"flagint": setting{
			Type:          _int,
			Name:          "flagint",
			Value:         42,
			Short:         "i",
			IsConfFileVar: true,
			IsEnvVar:      true,
			IsFlag:        true,
		},
		"flagint-tst": setting{
			Type:          _int,
			Name:          "flagint-tst",
			Value:         42,
			Short:         "",
			IsConfFileVar: true,
			IsEnvVar:      true,
			IsFlag:        true,
		},
		"flagint64": setting{
			Type:          _int64,
			Name:          "flagint64",
			Value:         int64(42),
			Short:         "",
			IsConfFileVar: true,
			IsEnvVar:      true,
			IsFlag:        true,
		},
		"flagint64-tst": setting{
			Type:          _int64,
			Name:          "flagint64-tst",
			Value:         int64(42),
			Short:         "",
			IsConfFileVar: true,
			IsEnvVar:      true,
			IsFlag:        true,
		},
		"flagstring": setting{
			Type:          _string,
			Name:          "flagstring",
			Value:         "a flag string",
			Short:         "s",
			IsConfFileVar: true,
			IsEnvVar:      true,
			IsFlag:        true,
		},
		"flagstring-tst": setting{
			Type:          _string,
			Name:          "flagstring-tst",
			Value:         "a flag string",
			Short:         "",
			IsConfFileVar: true,
			IsEnvVar:      true,
			IsFlag:        true,
		},
		"flagslice": setting{
			//Type:   "string-slice",
			Type:          _interface,
			Name:          "flagslice",
			Value:         []string{},
			Short:         "",
			IsConfFileVar: true,
			IsEnvVar:      true,
			IsFlag:        true,
		},
		"flagmap": setting{
			//Type:   "map",
			Type:          _interface,
			Name:          "flagmap",
			Value:         map[string]interface{}{},
			Short:         "",
			IsConfFileVar: true,
			IsEnvVar:      true,
			IsFlag:        true,
		},
		"bool": setting{
			Type:  _bool,
			Name:  "bool",
			Value: true,
			Short: "b",
		},
		"int": setting{
			Type:  _int,
			Name:  "int",
			Value: 42,
			Short: "i",
		},
		"int64": setting{
			Type:  _int64,
			Name:  "int64",
			Value: int64(42),
			Short: "",
		},
		"string": setting{
			Type:  _string,
			Name:  "string",
			Value: "a string",
			Short: "s",
		},
		"slice": setting{
			//Type:  "string-slice",
			Type:  _interface,
			Name:  "slice",
			Value: []string{},
			Short: "s",
		},
		"map": setting{
			//Type:  "map",
			Type:  _interface,
			Name:  "map",
			Value: map[string]interface{}{},
			Short: "s",
		},
	}

	s.shortFlags = map[string]string{
		"b": "flagbool",
		"i": "flagint",
		"s": "flagstring",
	}
	return s
}

func TestSettingExistsError(t *testing.T) {
	tests := []basic{
		basic{name: "test1", value: "dinosaur", expected: "", expectedErr: "dinosaur: setting exists"},
		basic{name: "test2", settingType: Core, value: "swallow", expected: "", expectedErr: "swallow: core setting exists"},
	}

	for _, test := range tests {
		err := error(SettingExistsError{typ: test.settingType, k: test.value})
		if err.Error() != test.expectedErr {
			t.Errorf("%s: expected %q got %q", test.name, test.expectedErr, err)
		}
	}
}

func TestShortFlagErr(t *testing.T) {
	err := ShortFlagExistsError{k: "bar", short: "b", shortName: "biz"}
	exp := "bar: short flag \"b\" already exists for \"biz\""
	if err.Error() != exp {
		t.Errorf("got %s; want %s", err.Error(), exp)
	}
}

func TestSettingNotFoundError(t *testing.T) {
	tests := []basic{
		basic{name: "notFoundErr test1", value: "dinosaur", expected: "", expectedErr: "dinosaur: setting not found"},
		basic{name: "notFoundErr test2", settingType: Core, value: "swallow", expected: "", expectedErr: "swallow: core setting not found"},
	}

	for _, test := range tests {
		err := error(SettingNotFoundError{settingType: test.settingType, k: test.value})
		if err.Error() != test.expectedErr {
			t.Errorf("%s: expected %q got %q", test.name, test.expectedErr, err)
		}
	}
}

func TestDataTypeErr(t *testing.T) {
	tests := []struct {
		name     string
		typ      dataType
		expected string
	}{
		{"corebool", _int, "corebool is bool, not int"},
		{"corestring", _int64, "corestring is string, not int64"},
		{"coreint", _bool, "coreint is int, not bool"},
		{"coreint64", _string, "coreint64 is int64, not string"},
	}

	var err error
	testSettings := newTestSettings()
	for _, test := range tests {
		switch test.typ {
		case _bool:
			_, err = testSettings.BoolE(test.name)
		case _int:
			_, err = testSettings.IntE(test.name)
		case _int64:
			_, err = testSettings.Int64E(test.name)
		case _string:
			_, err = testSettings.StringE(test.name)
		}
		if err == nil {
			t.Errorf("%s: expected error, got none", test.name)
			continue
		}
		if err.Error() != test.expected {
			t.Errorf("%s: got %s want %s", test.name, err, test.expected)
		}
	}
}

func TestFormatString(t *testing.T) {
	tests := []struct {
		name     string
		format   Format
		expected string
	}{
		{"Unsupported", Unsupported, "unsupported"},
		{"json", JSON, "json"},
		{"toml", TOML, "toml"},
		{"yaml", YAML, "yaml"},
	}
	for _, test := range tests {
		s := test.format.String()
		if s != test.expected {
			t.Errorf("format %s: expected %s got %s", test.name, test.expected, s)
		}
	}
}

func TestParseDataType(t *testing.T) {
	tests := []struct {
		v        string
		expected dataType
	}{
		// any type that isn't matched is considered an interface{}
		{"", _interface},
		{"strung", _interface},
		{"string", _string},
		{"Int", _int},
		{"int", _int},
		{"int64", _int64},
		{"bool", _bool},
		{"BOOL", _bool},
	}
	for i, test := range tests {
		v := parseDataType(test.v)
		if v != test.expected {
			t.Errorf("%d: got %s want %s", i, v, test.v)
		}
	}
}

func TestParseFilename(t *testing.T) {
	tests := []struct {
		name   string
		format Format
		err    error
	}{
		{"", Unsupported, UnsupportedFormatError{""}},
		{"file", Unsupported, UnsupportedFormatError{""}},
		{"file.xml", Unsupported, UnsupportedFormatError{"xml"}},
		{"file.jpeg", Unsupported, UnsupportedFormatError{"jpeg"}},
		{"file.json", JSON, nil},
		{"file.jsn", JSON, nil},
		{"file.cjson", JSON, nil},
		{"file.cjsn", JSON, nil},
		{"file.toml", TOML, nil},
		{"file.toml", TOML, nil},
		{"file.yaml", YAML, nil},
		{"file.yml", YAML, nil},
		{"pat/to/file.yaml", YAML, nil},
	}

	for _, test := range tests {
		format, err := ParseFilenameFormat(test.name)
		if err != nil {
			if err != test.err {
				t.Errorf("%s: got %s; want %s", test.name, err, test.err)
			}
			continue
		}
		if test.err != nil {
			t.Errorf("%s: got no error; want %s", test.name, test.err)
			continue
		}
		if format != test.format {
			t.Errorf("%s: got %s; want %s", test.name, format, test.format)
		}
	}
}

func TestPathsFromEnvVar(t *testing.T) {
	ps := string(os.PathListSeparator)
	tests := []struct {
		k        string
		v        string
		expected []string
	}{
		{"", "", nil},
		{"FEEDTEST_1", "", nil},
		{"FEEDTEST_2", "x", []string{"x"}},
		{"FEEDTEST_3", fmt.Sprintf("a%[1]sb", ps), []string{"a", "b"}},
		{"FEEDTEST_4", fmt.Sprintf("a%[1]sb%[1]sc", ps), []string{"a", "b", "c"}},
		{"FEEDTEST_5", fmt.Sprintf("a%[1]sb%[1]s%[2]s", ps, filepath.Join("y", "$FEEDTEST_2")), []string{"a", "b", filepath.Join("y", "x")}},
	}
	for _, test := range tests {
		os.Setenv(test.k, test.v)
		v := PathsFromEnvVar(test.k)
		if len(v) != len(test.expected) {
			t.Errorf("%s: got %d elements; want %d", test.k, len(v), len(test.expected))
			continue
		}
		sort.Strings(v)
		sort.Strings(test.expected)
		if !reflect.DeepEqual(v, test.expected) {
			t.Errorf("%s: got %v; want %v", test.k, v, test.expected)
		}
	}

	for _, test := range tests {
		os.Unsetenv(test.k)
	}
}
