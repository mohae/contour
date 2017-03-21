package contour

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"testing"
	"time"
)

func init() {
	log.SetOutput(ioutil.Discard)
	rand.Seed(int64(time.Now().Nanosecond()))
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

// temporary one; once slice and map support has been added use the original one.
var jsonTest = []byte(`
{
	"cfgbool": true,
	"flagbool": true,
	"cfgint": 42,
	"flagint": 1999,
	"cfgstring": "foo",
	"flagstring": "bar"
}
`)

/*
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
*/
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
	app:     Settings{settings: map[string]setting{}},
	"test1": Settings{settings: map[string]setting{}},
}

func newTestSettings() *Settings {
	return &Settings{
		flagSet:           flag.NewFlagSet(fmt.Sprintf("rancher-%d", rand.Int63()), flag.ContinueOnError),
		useCfg:            true,
		useEnv:            true,
		filterVars:        map[string]interface{}{},
		boolFilterNames:   []string{},
		intFilterNames:    []string{},
		int64FilterNames:  []string{},
		stringFilterNames: []string{},
		settings: map[string]setting{
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
			/*
				"coreslice": setting{
					Type:   "string-slice",
					Name:   "coreslice",
					Value:  []string{},
					IsCore: true,
				},
				"coremap": setting{
					Type:   "map",
					Name:   "coremap",
					Value:  map[string]interface{}{},
					IsCore: true,
				},
			*/
			"cfgbool": setting{
				Type:  _bool,
				Name:  "cfgbool",
				Value: true,
				Short: "",
				IsCfg: true,
				IsEnv: true,
			},
			"cfgint": setting{
				Type:  _int,
				Name:  "cfgint",
				Value: 42,
				IsCfg: true,
				IsEnv: true,
			},
			"cfgint64": setting{
				Type:  _int64,
				Name:  "cfgint64",
				Value: int64(42),
				IsCfg: true,
				IsEnv: true,
			},
			"cfgstring": setting{
				Type:  _string,
				Name:  "cfgstring",
				Value: "a cfg string",
				Short: "",
				IsCfg: true,
				IsEnv: true,
			},
			/*
				"cfgslice": setting{
					Type:  "string-slice",
					Name:  "cfgslice",
					Value: []string{},
					Short: "",
					IsCfg: true,
					IsEnv: true,
				},
				"cfgmap": setting{
					Type:  "map",
					Name:  "cfgmap",
					Value: map[string]interface{}{},
					Short: "",
					IsCfg: true,
					IsEnv: true,
				},
			*/
			"flagbool": setting{
				Type:   _bool,
				Name:   "flagbool",
				Value:  true,
				Short:  "",
				IsCfg:  true,
				IsEnv:  true,
				IsFlag: true,
			},
			"flagbool-tst": setting{
				Type:   _bool,
				Name:   "flagbool-tst",
				Value:  true,
				Short:  "",
				IsCfg:  true,
				IsEnv:  true,
				IsFlag: true,
			},
			"flagint": setting{
				Type:   _int,
				Name:   "flagint",
				Value:  42,
				Short:  "",
				IsCfg:  true,
				IsEnv:  true,
				IsFlag: true,
			},
			"flagint-tst": setting{
				Type:   _int,
				Name:   "flagint-tst",
				Value:  42,
				Short:  "",
				IsCfg:  true,
				IsEnv:  true,
				IsFlag: true,
			},
			"flagint64": setting{
				Type:   _int64,
				Name:   "flagint64",
				Value:  int64(42),
				Short:  "",
				IsCfg:  true,
				IsEnv:  true,
				IsFlag: true,
			},
			"flagint64-tst": setting{
				Type:   _int64,
				Name:   "flagint64-tst",
				Value:  int64(42),
				Short:  "",
				IsCfg:  true,
				IsEnv:  true,
				IsFlag: true,
			},
			"flagstring": setting{
				Type:   _string,
				Name:   "flagstring",
				Value:  "a flag string",
				Short:  "",
				IsCfg:  true,
				IsEnv:  true,
				IsFlag: true,
			},
			"flagstring-tst": setting{
				Type:   _string,
				Name:   "flagstring-tst",
				Value:  "a flag string",
				Short:  "",
				IsCfg:  true,
				IsEnv:  true,
				IsFlag: true,
			},
			/*
				"flagslice": setting{
					Type:   "string-slice",
					Name:   "flagslice",
					Value:  []string{},
					Short:  "",
					IsCfg:  true,
					IsEnv:  true,
					IsFlag: true,
				},
				"flagmap": setting{
					Type:   "map",
					Name:   "flagmap",
					Value:  map[string]interface{}{},
					Short:  "",
					IsCfg:  true,
					IsEnv:  true,
					IsFlag: true,
				},
			*/
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
			/*
				"slice": setting{
					Type:  "string-slice",
					Name:  "slice",
					Value: []string{},
					Short: "s",
				},
				"map": setting{
					Type:  "map",
					Name:  "map",
					Value: map[string]interface{}{},
					Short: "s",
				},
			*/
		}}
}

func TestNotFoundErr(t *testing.T) {
	tests := []basic{
		basic{"notFoundErr test1", 0, "setting", "", "setting: not found"},
		basic{"notFoundErr test2", 0, "grail", "", "grail: not found"},
	}
	for _, test := range tests {
		err := error(NotFoundErr{test.value})
		if err.Error() != test.expectedErr {
			t.Errorf("%s: expected %s, got %s", test.name, test.expectedErr, err.Error())
		}
	}

}

func TestSettingNotFoundErr(t *testing.T) {
	tests := []basic{
		basic{name: "notFoundErr test1", value: "dinosaur", expected: "", expectedErr: "dinosaur: setting not found"},
		basic{name: "notFoundErr test2", settingType: Core, value: "swallow", expected: "", expectedErr: "swallow: core setting not found"},
	}

	for _, test := range tests {
		err := error(SettingNotFoundErr{settingType: test.settingType, name: test.value})
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
		err      string
	}{
		{"", 0, ": not a supported data type"},
		{"strung", 0, "strung: not a supported data type"},
		{"Int32", 0, "Int32: not a supported data type"},
		{"string", _string, ""},
		{"Int", _int, ""},
		{"int", _int, ""},
		{"int64", _int64, ""},
		{"bool", _bool, ""},
		{"BOOL", _bool, ""},
	}
	for i, test := range tests {
		v, err := parseDataType(test.v)
		if err != nil {
			if err.Error() != test.err {
				t.Errorf("%d: got %s, want %s", i, err, test.err)
			}
			continue
		}
		if test.err != "" {
			t.Errorf("%d: wanted %s, got no error", i, test.err)
			continue
		}
		if v != test.expected {
			t.Errorf("%d: got %s want %s", i, v, test.v)
		}
	}
}
