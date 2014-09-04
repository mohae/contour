package contour

import (
	"bytes"
	_"os"
	"testing"

	"github.com/mohae/customjson"
	utils "github.com/mohae/utilitybelt"
	. "github.com/smartystreets/goconvey/convey"
)

type basic struct {
	name string
	value string
	expected string
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

var tomlResults = map[string]interface {}{
	"appVar1":true,
  	"appVar2":false,
	"appVar3":42,
	"appVar4":"zip",
	"appVar5":[]string{"less","sass","scss"},
	"logging":map[string]interface{}{
		"Logging":true,
		"LogConfig":"test/test.toml",
		"LogFileLevel":"debug",
		"LogStdoutLevel":"error",
	},
}

var jsonResults = map[string]interface {}{
	"appVar1":true,
  	"appVar2":false,
	"appVar3":42,
	"appVar4":"zip",
	"appVar5":[]string{"less","sass","scss"},
	"logging":map[string]interface{}{
		"logging":true,
		"logconfig":"test/test.toml",
		"logfilelevel":"debug",
		"logstdoutlevel":"error",
	},
}

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

func TestGetConfigFormat(t *testing.T) {
	tests := []basic {
		{"an empty configfilename", "", "", "a config filename was expected, none received"},
		{"a configfilename without an extension", "config", "", "unable to determine config format, the configuration file, config, doesn't have an extension"},
		{"a configfilename with an invalid extension", "config.bmp", "", "bmp is an unsupported format for configuration files"},
		{"a configfilename with a json extension", "config.json", "json", ""},
		{"a configfilename with a toml extension", "config.toml", "toml", ""},
	}

	for _, test := range tests {
		Convey("Given " + test.name + " Test", t, func() {
		
			Convey("Getting the config format", func() {
				format, err := getConfigFormat(test.value)
				checkTestReturn(test, format, err)
			})
		})
	}

}

func TestIsSupportedFormat(t *testing.T) {
	tests := []basic{
		{"empty format test", "", "false", ""},
		{"invalid format test", "bmp", "false", ""},
		{"json format test", "json", "true", ""},
		{"toml", "toml", "true", ""},
	}

	for _, test := range tests {
		Convey("Given some supported format tests", t, func() {
		
			Convey(test.name, func() {
				is := utils.BoolToString(isSupportedFormat(test.value))
			
				Convey("Should result in " + test.expected, func() {
					So(is, ShouldEqual, test.expected)
				})
			})

		})

	}

}

func TestLoadEnvs(t *testing.T) {


}

// Only testing failure for now
func TestLoadConfigFile(t *testing.T) {

	Convey("Given an empty config filename", t, func() {
		AppConfig.Settings[EnvConfigFilename] = &setting{Value: ""}
		AppConfig.Settings[EnvConfigFormat] = &setting{Value: ""}

		Convey("loading the config file", func() {		
			err := loadConfigFile()
			
			Convey("Should not result in an error", func() {
				So(err, ShouldBeNil)
			})
		
		})

	})

	Convey("Given an invalid config filename", t, func() {
		AppConfig.Settings[EnvConfigFilename] = &setting{Value: "holygrail"}
		AppConfig.Settings[EnvConfigFormat] = &setting{Value: ""}
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
				err := MarshalFormatReader("json", r)

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
				err := MarshalFormatReader("toml", r)

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

func TestCanUpdate(
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
					So(AppConfig.Settings[test.key], ShouldResemble, test.expected)
				})
				
			})

		}

	})

}

/*
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
				So(AppConfig.Settings[test.key], ShouldResemble, test.expected)
			})

		})

	}

}
*/
