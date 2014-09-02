package contour

import (
	"bytes"
	"testing"
	"os"

	"github.com/mohae/customjson"
	. "github.com/smartystreets/goconvey/convey"
)

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

func TestGetConfigFormat(t *testing.T) {
	Convey("Given a GetconfigFormatTest", t, func() {

		Convey("Given an empty config format environment variable", func() {
			os.Setenv(EnvConfigFormat, "")
			Convey("Given an empty configfilename", func() {
				os.Setenv(EnvConfigFilename, "")
				
				Convey("Getting the config format", func() {
					_, err := getConfigFormat()

					Convey("Should result in an error", func() {
						So(err.Error(), ShouldEqual, "unable to determine config format, filename not set")
					}) 

				})

			})

			Convey("Given a filename without an extension", func() {
				os.Setenv(EnvConfigFilename, "config")
				
				Convey("Getting the config format", func() {
					_, err := getConfigFormat()

					Convey("Should result in an error", func() {
						So(err.Error(), ShouldEqual, "unable to determine config format, the configuration file config doesn't have an extension")
					}) 

				})

			})

			Convey("Given a filename with an invalid format extension", func() {
				os.Setenv(EnvConfigFormat, "config.bmp")
				
				Convey("Getting the config format", func() {
					_, err := getConfigFormat()

					Convey("Should result in an error", func() {
						So(err.Error(), ShouldEqual, "unable to determine config format, the configuration file config doesn't have an extension")
					}) 

				})

			})

			Convey("Given a filename with a json extension", func() {				

				Convey("Getting the config format", func() {
					os.Setenv(EnvConfigFormat, "json")
					res, err := getConfigFormat()

					Convey("Should not error", func() {
						So(err, ShouldBeNil)
					}) 

					Convey("Should result in getting json as an extension", func() {
						So(res, ShouldEqual, "json")
					})

				})

			})

			Convey("Given a filename with a json extension", func() {
				
				Convey("Getting the config format", func() {
					os.Setenv(EnvConfigFormat, "toml")
					res, err := getConfigFormat()

					Convey("Should not error", func() {
						So(err, ShouldBeNil)
					}) 

					Convey("Should result in getting json as an extension", func() {
						So(res, ShouldEqual, "toml")
					})

				})

			})
		})

		Convey("Given setting the environment variable with json", func() {
			os.Setenv(EnvConfigFormat, "json")
			
			Convey("Calling getConfigFormat", func() {
				res, err := getConfigFormat()
				
				Convey("Should not error", func() {
					So(err, ShouldBeNil)
				})

				Convey("and the format should equal json", func() {
					So(res, ShouldEqual, "json")
				})

			})

		})


		Convey("Given setting the environment variable with toml", func() {
			os.Setenv(EnvConfigFormat, "toml")
			
			Convey("Calling getConfigFormat", func() {
				res, err := getConfigFormat()
				
				Convey("Should not error", func() {
					So(err, ShouldBeNil)
				})

				Convey("and the format should equal toml", func() {
					So(res, ShouldEqual, "toml")
				})

			})

		})

		Convey("Given setting the environment variable with an unsupported format", func() {
			os.Setenv(EnvConfigFormat, "png")
			
			Convey("Calling getConfigFormat", func() {
				_, err := getConfigFormat()
				
				Convey("Should error", func() {
					So(err.Error(), ShouldEqual, "unable to determine config format, the configuration file config doesn't have an extension")
				})

			})

		})

	})

}

func TestIsSupportedFormat(t *testing.T) {
	Convey("Given some supported format tests", t, func() {
		
		Convey("Checking to see if json is supported", func() {
			is := isSupportedFormat("json") 
			
			Convey("Should result in true", func() {
				So(is, ShouldEqual, true)
			})
		})

		Convey("Checking to see if toml is supported", func() {
			is := isSupportedFormat("toml") 
			
			Convey("Should result in true", func() {
				So(is, ShouldEqual, true)
			})
		})

		Convey("Checking to see if tnt is supported", func() {
			is := isSupportedFormat("tnt")  
			
			Convey("Should result in false", func() {
				So(is, ShouldEqual, false)
			})
		})

	})	

}


// Only testing failure for now
func TestLoadConfigFile(t *testing.T) {
	Convey("Given an unset config filename", t, func() {

		Convey("loading the config file", func() {		
			os.Setenv(EnvConfigFilename, "")
			err := LoadConfigFile()
			
			Convey("Should not result in an error", func() {
				So(err, ShouldBeNil)
			})
		
		})

	})

	Convey("Given an invalid config filename", t, func() {
		os.Setenv(EnvConfigFilename, "holygrail")
		Convey("loading the config file", func() {		
			err := LoadConfigFile()
			
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
