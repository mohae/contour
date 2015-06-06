package contour

import (
	"bytes"
	"io/ioutil"
	"log"
	"strconv"
	"testing"

	"github.com/mohae/customjson"
	. "github.com/smartystreets/goconvey/convey"
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

func newTestCfg() Cfg {
	return Cfg{settings: map[string]setting{
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
	}}
}

// helper function
func checkTestReturn(test basic, format string, err error) {
	if err != nil {
		if test.expectedErr != "" {
			Convey("Should result in an error", func() {
				So(err.Error(), ShouldEqual, test.expectedErr)
			})
		} else {
			Convey("Should not error", func() {
				So(err, ShouldBeNil)
			})
		}
	} else {
		if test.expectedErr != "" {
			Convey("Should not result in an error", func() {
				So(err.Error(), ShouldEqual, test.expectedErr)
			})
		} else {
			Convey("Should result in the file's extenstion", func() {
				So(format, ShouldEqual, test.expected)
			})
		}
	}
}

// Testing
func TestFormatFromFilename(t *testing.T) {
	tests := []basic{
		{"an empty cfgfilename", "", "", "no config filename"},
		{"a cfgfilename without an extension", "cfg", "", "unable to determine cfg's config format: no extension"},
		{"a cfgfilename with an invalid extension", "cfg.bmp", "", "unsupported cfg format: bmp"},
		{"a cfgfilename with a json extension", "cfg.json", "json", ""},
		{"a path and multi dot cfgfilename with a json extension", "path/to/custom.cfg.json", "json", ""},
		{"a cfgfilename with a toml extension", "cfg.toml", "toml", ""},
		{"a cfgfilename with a toml extension", "cfg.yaml", "yaml", ""},
		{"a cfgfilename with a toml extension", "cfg.yml", "yaml", ""},
		{"a cfgfilename with a toml extension", "cfg.xml", "xml", "unsupported cfg format: xml"},
		{"a cfgfilename with a toml extension", "cfg.ini", "", "unsupported cfg format: ini"},
	}
	for _, test := range tests {
		Convey("Given "+test.name+" Test", t, func() {
			Convey("Getting the cfg format", func() {
				format, err := formatFromFilename(test.value)
				checkTestReturn(test, format.String(), err)
			})
		})
	}

}

func TestIsSupportedFormat(t *testing.T) {
	tests := []basic{
		{"empty format test", "", "false", ""},
		{"invalid format test", "bmp", "false", ""},
		{"json format test", "json", "true", ""},
		{"tom format testl", "toml", "true", ""},
		{"yaml format test", "yaml", "true", ""},
		{"yml format test", "yml", "true", ""},
		{"xml format test", "xml", "false", ""},
	}
	for _, test := range tests {
		Convey("Given some supported format tests", t, func() {

			Convey(test.name, func() {
				formatString := ParseFormat(test.value)
				is := formatString.isSupported()
				Convey("Should result in "+test.expected, func() {
					So(strconv.FormatBool(is), ShouldEqual, test.expected)
				})
			})

		})

	}

}

func TestMarshalFormatReader(t *testing.T) {
	tests := []struct {
		name        string
		format      Format
		value       []byte
		expected    interface{}
		expectedErr string
	}{
		{"json cfg", JSON, jsonExample, jsonResults, ""},
		{"toml cfg", TOML, tomlExample, tomlResults, ""},
		{"yaml cfg", YAML, yamlExample, []byte(""), "unsupported cfg format: yaml"},
		{"xml cfg", XML, xmlExample, []byte(""), "unsupported cfg format: xml"},
		{"unsupported cfg", Unsupported, []byte(""), []byte(""), "unsupported cfg format: unsupported"},
	}
	for _, test := range tests {
		r := bytes.NewReader([]byte(test.value))
		ires, err := unmarshalFormatReader(test.format, r)
		if err != nil {
			if test.expectedErr == "" {
				t.Errorf("%s: expected nil for error; got %q", test.name, err)
				continue
			}
			if err.Error() != test.expectedErr {
				t.Errorf("%s: expected %q; got %q", test.name, test.expectedErr, err)
				continue
			}
		} else {
			val, ok := ires.(map[string]interface{})["appVar1"]
			if !ok {
				t.Errorf("appVar1 not found")
			} else {
				if val != test.expected.(map[string]interface{})["appVar1"] {
					t.Errorf("appVar1: expected %q, got %q", test.expected.(map[string]interface{})["appVar1"], val)
				}
			}
			val, ok = ires.(map[string]interface{})["appVar2"]
			if !ok {
				t.Errorf("appVar2 not found")
			} else {
				if val != test.expected.(map[string]interface{})["appVar2"] {
					t.Errorf("appVar2: expected %q, got %q", test.expected.(map[string]interface{})["appVar2"], val)
				}
			}
			val, ok = ires.(map[string]interface{})["appVar3"]
			if !ok {
				t.Errorf("appVar3 not found")
			} else {
				if val != test.expected.(map[string]interface{})["appVar3"] {
					t.Errorf("appVar3: expected %q, got %q", test.expected.(map[string]interface{})["appVar3"], val)
				}
			}
			val, ok = ires.(map[string]interface{})["appVar4"]
			if !ok {
				t.Errorf("appVar4 not found")
			} else {
				if val != test.expected.(map[string]interface{})["appVar4"] {
					t.Errorf("appVar4: expected %q, got %q", test.expected.(map[string]interface{})["appVar4"], val)
				}
			}
			val, ok = ires.(map[string]interface{})["appVar5"]
			if !ok {
				t.Errorf("appVar5 not found")
			}
		}
	}
}

func TestCanUpdate(t *testing.T) {
	tests := []struct {
		value       string
		expected    string
		expectedErr string
	}{
		{"corestring", "false", "cannot update \"corestring\": core settings cannot be updated"},
		{"flagstring", "true", ""},
		{"cfgstring", "true", ""},
		{"string", "true", ""},
		{"arr", "false", "cannot update \"arr\": not found"},
		{"", "false", "cannot update \"\": not found"},
	}
	tstCfg := newTestCfg()
	for i, test := range tests {
		res, err := tstCfg.canUpdate(test.value)
		if err != nil {
			if err.Error() != test.expectedErr {
				t.Errorf("%d\texpected %q, got %q", i, test.expectedErr, err.Error())
			}
			continue
		}
		if strconv.FormatBool(res) != test.expected {
			t.Errorf("%d\texpected %q, got %q", i, test.expected, strconv.FormatBool(res))
		}
	}
}

func TestNotFoundErr(t *testing.T) {
	tests := []basic{
		basic{"notFoundErr test1", "setting", "not found: setting", ""},
		basic{"notFoundErr test2", "grail", "not found: grail", ""},
	}
	for _, test := range tests {
		Convey(test.name+"  given a string", t, func() {
			Convey("calling notFoundErr with it", func() {
				err := notFoundErr(test.value)
				Convey("should result in an error", func() {
					So(err, ShouldNotBeNil)
					Convey("with the error message", func() {
						So(err.Error(), ShouldEqual, test.expected)
					})
				})
			})
		})
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
