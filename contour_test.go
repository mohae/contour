package contour

import (
	_ "bytes"
	_ "os"
	"strconv"
	"testing"

	"github.com/mohae/customjson"
	. "github.com/smartystreets/goconvey/convey"
)

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
appVar3 = 42
appVar4 = "zip"
appVar5 = [
	"less",
	"sass",
	"scss"
]

[logging]
Logging = true
LogConfig = "test/test.toml"
LogFileLevel = "debug"
LogStdoutLevel = "error"
`)

var jsonExample = []byte(`
{
	"appVar1": true,
	"appVar2": false,
	"appVar3": 42,
	"appVar4": "zip",
	"appVar5": [
		"less",
		"sass",
		"scss"
	],
	"logging": {
		"logging": true,
		"logconfig": "test/test.toml",
		"logfilelevel": "debug",
		"logstdoutlevel": "error"
	}
}
`)

var tomlResults = map[string]interface{}{
	"appVar1": true,
	"appVar2": false,
	"appVar3": 42,
	"appVar4": "zip",
	"appVar5": []string{"less", "sass", "scss"},
	"logging": map[string]interface{}{
		"Logging":        true,
		"LogConfig":      "test/test.toml",
		"LogFileLevel":   "debug",
		"LogStdoutLevel": "error",
	},
}

var jsonResults = map[string]interface{}{
	"appVar1": true,
	"appVar2": false,
	"appVar3": 42,
	"appVar4": "zip",
	"appVar5": []string{"less", "sass", "scss"},
	"logging": map[string]interface{}{
		"logging":        true,
		"logconfig":      "test/test.toml",
		"logfilelevel":   "debug",
		"logstdoutlevel": "error",
	},
}

var testConfig = &Cfg{settings: map[string]*setting{
	"corebool": &setting{
		Type:   "bool",
		Value:  true,
		IsCore: true,
	},
	"coreint": &setting{
		Type:   "int",
		Value:  42,
		IsCore: true,
	},
	"corestring": &setting{
		Type:   "string",
		Value:  "a core string",
		IsCore: true,
	},
	"configbool": &setting{
		Type:  "bool",
		Value: true,
		Short: "",
		IsCfg: true,
	},
	"configint": &setting{
		Type:  "int",
		Value: 42,
		IsCfg: true,
	},
	"configstring": &setting{
		Type:  "string",
		Value: "a config string",
		Short: "",
		IsCfg: true,
	},
	"flagbool": &setting{
		Type:   "bool",
		Value:  true,
		Short:  "b",
		IsFlag: true,
		IsCfg:  true,
	},
	"flagint": &setting{
		Type:   "int",
		Value:  42,
		Short:  "i",
		IsFlag: true,
		IsCfg:  true,
	},
	"flagstring": &setting{
		Type:   "string",
		Value:  "a flag string",
		Short:  "s",
		IsFlag: true,
		IsCfg:  true,
	},
	"bool": &setting{
		Type:  "bool",
		Value: true,
		Short: "b",
	},
	"int": &setting{
		Type:  "int",
		Value: 42,
		Short: "i",
	},
	"string": &setting{
		Type:  "string",
		Value: "a string",
		Short: "s",
	},
}}

var emptyConfigs map[string]*Cfg
var testConfigs = map[string]*Cfg{
	app:     &Cfg{settings: map[string]*setting{}},
	"test1": &Cfg{settings: map[string]*setting{}},
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
func TestConfigFormat(t *testing.T) {
	tests := []basic{
		{"an empty configfilename", "", "", "a config filename was expected, none received"},
		{"a configfilename without an extension", "config", "", "unable to determine config format, the configuration file, config, doesn't have an extension"},
		{"a configfilename with an invalid extension", "config.bmp", "", "bmp is not a supported configuration file format"},
		{"a configfilename with a json extension", "config.json", "json", ""},
		{"a configfilename with a toml extension", "config.toml", "toml", ""},
		{"a configfilename with a toml extension", "config.yaml", "yaml", ""},
		{"a configfilename with a toml extension", "config.yml", "yaml", ""},
		{"a configfilename with a toml extension", "config.xml", "xml", ""},
		{"a configfilename with a toml extension", "config.ini", "", "ini is not a supported configuration file format"},
	}

	for _, test := range tests {
		Convey("Given "+test.name+" Test", t, func() {

			Convey("Getting the config format", func() {
				format, err := cfgFormat(test.value)
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
		{"xml format test", "xml", "true", ""},
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

/*
func TestLoadEnvs(t *testing.T) {

}


// Only testing failure for now
func TestLoadConfigFile(t *testing.T) {

	Convey("Given an empty config filename", t, func() {
		AppConfig.settings[EnvConfigFilename] = &setting{Value: ""}
		AppConfig.settings[EnvConfigFormat] = &setting{Value: ""}

		Convey("loading the config file", func() {
			err := loadConfigFile()

			Convey("Should not result in an error", func() {
				So(err, ShouldBeNil)
			})

		})

	})

	Convey("Given an invalid config filename", t, func() {
		AppConfig.settings[EnvConfigFilename] = &setting{Value: "holygrail"}
		AppConfig.settings[EnvConfigFormat] = &setting{Value: ""}
		Convey("loading the config file", func() {
			err := loadConfigFile()

			Convey("Should result in an error", func() {
				So(err.Error(), ShouldEqual, "open holygrail: no such file or directory")
			})

		})

	})

}

func TestMarshalFormatReader(t *testing.T) {

	Convey("Given an JSON config", t, func() {

		Convey("Given a []byte", func() {

			Convey("marshalling it should result in", func() {
				r := bytes.NewReader(jsonExample)
				err := marshalFormatReader("json", r)

				Convey("Should not error", func() {
					So(err, ShouldBeNil)
				})

				Convey("Should equal our expectations", func() {
					So(toString.Get(configFile), ShouldEqual, toString.Get(jsonResults))
				})

			})

		})

	})

	Convey("Given an TOML config", t, func() {

		Convey("Given a []byte", func() {

			Convey("marshalling it should result in", func() {
				r := bytes.NewReader(tomlExample)
				err := marshalFormatReader("toml", r)

				Convey("Should not error", func() {
					So(err, ShouldBeNil)
				})

				Convey("Should equal our expectations", func() {
					So(toString.Get(configFile), ShouldEqual, toString.Get(tomlResults))
				})

			})

		})

	})

}

func TestCanUpdate(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		value2   string
		expected string
	}{
		{"update a core setting", "configfilename", "another.file", "false"},
		{"update an immutable setting with value", "logconfigfilename", "logconfig.xml", "false"},
		{"update an immutable setting without a value", "unsetimmutable", "json", "true"},
		{"update a setting", "logging", "true", "true"},
		{"update a setting that does not exist", "arr", "", "false"},
	}

	AppConfig = testConfig

	Convey("Given some CanUpdate tests", t, func() {
		for _, test := range tests {

			Convey("Given a setting test: "+test.name, func() {
				res := utils.BoolToString(canUpdate(test.name))

				Convey("Should result in "+test.expected, func() {
					So(res, ShouldEqual, test.expected)
				})

			})

		}
	})

}

/*
// Since SetIdemString wraps SetIdempotentString, it is called instead-2for1!
func TestSetIdempotentString(t *testing.T) {
	tests := []struct{
		name string
		key string
		value string
		expected *setting
	}{
		{name: "test empty idempotent", key: "rock", value: "roll", expected:
			&setting{
				Code: "",
				Type: "string",
				Value: "roll",
				IsFlag: false,
				IsIdempotent: true,
				SourceIsEnv: false,
				IsCore:	false,
			},
		},
	}


	Convey("Given a range of tests", t, func() {

		for _, test := range tests {

			Convey("setting them should not error", func() {
				err := SetIdempotentString(test.key, test.value)
				So(err, ShouldBeNil)

				Convey("and getting it should result it", func() {
					res := os.Getenv(test.key)
					So(res, ShouldEqual, test.value)
				})

				Convey("and the AppConfig settings for it", func() {
					So(AppConfig.settings[test.key], ShouldResemble, test.expected)
				})

			})

		}

	})

}

func TestSetBoolFlag(t *testing.T) {
	tests := []struct{
		name string
		key string
		value string
		b bool
		expected setting
	}{
		{name: "setboolflag", key: "bool-t-true", value: "t", b: true, expected: setting{}},
		{name: "setboolflag", key: "bool-t-false", value: "t", b: false, expected: setting{}},
		{name: "setboolflag", key: "bool-true", value: "", b: true, expected: setting{}},
		{name: "setboolflag", key: "bool-false", value: "", b: false, expected: setting{}},
	}


	for _, test := range tests {

		Convey("Setting a bool flag", t, func() {
			SetBoolFlag(test.key, test.value, test.b)

			Convey("Should result in the setting be set", func() {
				So(AppConfig.settings[test.key], ShouldResemble, test.expected)
			})

		})

	}

}
*/

func TestNotFoundErr(t *testing.T) {
	tests := []basic{
		basic{"notFoundErr test1", "setting", "setting not found", ""},
		basic{"notFoundErr test2", "grail", "grail not found", ""},
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
		basic{"notFoundErr test1", "dinosaur", "dinosaur: setting not found", ""},
		basic{"notFoundErr test2", "swallow type", "swallow type: setting not found", ""},
	}

	for _, test := range tests {

		Convey(test.name+"  given a string", t, func() {
			Convey("calling notFoundErr with it", func() {
				err := settingNotFoundErr(test.value)
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
