package contour

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/mohae/customjson"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

type basic struct {
	name        string
	value       string
	expected    string
	expectedErr string
}

var toString = customjson.NewMarshalString()

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

var yamlExample = []byte(`
appVar1: true
appVar2: false
appVar3: 42
appVar4: zip
appVar5:
  - less
  -	iass
  -	scss


logging:
  - Logging: true
  - LogCfg: test/test.yaml
  - LogFileLevel: debug
  - LogStdoutLevel: error
`)

var xmlExample = []byte(`
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

var yamlResults = map[string]interface{}{
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

var xmlResults = map[string]interface{}{
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

var emptyCfgs map[string]Cfg
var testCfgs = map[string]Cfg{
	app:     Cfg{settings: map[string]setting{}},
	"test1": Cfg{settings: map[string]setting{}},
}

func newTestCfg() *Cfg {
	return &Cfg{settings: map[string]setting{
		"corebool": setting{
			Type:   "bool",
			Value:  true,
			IsCore: true,
		},
		"coreint": setting{
			Type:   "int",
			Value:  42,
			IsCore: true,
		},
		"coreint64": setting{
			Type:   "int64",
			Value:  int64(42),
			IsCore: true,
		},
		"corestring": setting{
			Type:   "string",
			Value:  "a core string",
			IsCore: true,
		},
		"coreslice": setting{
			Type:   "string-slice",
			Value:  []string{},
			IsCore: true,
		},
		"coremap": setting{
			Type:   "map",
			Value:  map[string]interface{}{},
			IsCore: true,
		},
		"cfgbool": setting{
			Type:  "bool",
			Value: true,
			Short: "",
			IsCfg: true,
			IsEnv: true,
		},
		"cfgint": setting{
			Type:  "int",
			Value: 42,
			IsCfg: true,
			IsEnv: true,
		},
		"cfgint64": setting{
			Type:  "int64",
			Value: int64(42),
			IsCfg: true,
			IsEnv: true,
		},
		"cfgstring": setting{
			Type:  "string",
			Value: "a cfg string",
			Short: "",
			IsCfg: true,
			IsEnv: true,
		},
		"cfgslice": setting{
			Type:  "string-slice",
			Value: []string{},
			Short: "",
			IsCfg: true,
			IsEnv: true,
		},
		"cfgmap": setting{
			Type:  "map",
			Value: map[string]interface{}{},
			Short: "",
			IsCfg: true,
			IsEnv: true,
		},
		"flagbool": setting{
			Type:   "bool",
			Value:  true,
			Short:  "",
			IsCfg:  true,
			IsEnv:  true,
			IsFlag: true,
		},
		"flagint": setting{
			Type:   "int",
			Value:  42,
			Short:  "",
			IsCfg:  true,
			IsEnv:  true,
			IsFlag: true,
		},
		"flagint64": setting{
			Type:   "int64",
			Value:  int64(42),
			Short:  "",
			IsCfg:  true,
			IsEnv:  true,
			IsFlag: true,
		},
		"flagstring": setting{
			Type:   "string",
			Value:  "a flag string",
			Short:  "",
			IsCfg:  true,
			IsEnv:  true,
			IsFlag: true,
		},
		"flagslice": setting{
			Type:   "string-slice",
			Value:  []string{},
			Short:  "",
			IsCfg:  true,
			IsEnv:  true,
			IsFlag: true,
		},
		"flagmap": setting{
			Type:   "map",
			Value:  map[string]interface{}{},
			Short:  "",
			IsCfg:  true,
			IsEnv:  true,
			IsFlag: true,
		},
		"bool": setting{
			Type:  "bool",
			Value: true,
			Short: "b",
		},
		"int": setting{
			Type:  "int",
			Value: 42,
			Short: "i",
		},
		"int64": setting{
			Type:  "int64",
			Value: int64(42),
			Short: "",
		},
		"string": setting{
			Type:  "string",
			Value: "a string",
			Short: "s",
		},
		"slice": setting{
			Type:  "string-slice",
			Value: []string{},
			Short: "s",
		},
		"map": setting{
			Type:  "map",
			Value: map[string]interface{}{},
			Short: "s",
		},
	}}
}

func TestNotFoundErr(t *testing.T) {
	tests := []basic{
		basic{"notFoundErr test1", "setting", "not found: setting", ""},
		basic{"notFoundErr test2", "grail", "not found: grail", ""},
	}
	for _, test := range tests {
		err := notFoundErr(test.value)
		if err != nil {
			if err.Error() != test.expected {
				t.Errorf("%s: expected %s, got %s", test.name, test.expected, err.Error())
			}
			continue
		}
		if test.expected != "" {
			t.Errorf("%s: expected %s, got no error", test.name, test.expected)
		}
	}

}

func TestSettingNotFoundErr(t *testing.T) {
	tests := []basic{
		basic{"notFoundErr test1", "dinosaur", "setting not found: dinosaur", ""},
		basic{"notFoundErr test2", "swallow type", "setting not found: swallow type", ""},
	}
	for _, test := range tests {
		err := settingNotFoundErr(test.value)
		if err.Error() != test.expected {
			t.Errorf("%s: expected %q got %q", test.name, test.expected, err)
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
		{"xml", XML, "xml"},
	}
	for _, test := range tests {
		s := test.format.String()
		if s != test.expected {
			t.Errorf("format %s: expected %s got %s", test.name, test.expected, s)
		}
	}
}
