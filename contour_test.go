package contour

import (
	"bytes"
	"testing"

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

func TestSetAppCode(t *testing.T) {

	Convey("Given an excuse to use this convey", t, func() {

		Convey("calling SetAppCode with an empty value", func() {
			err := SetAppCode("")

			Convey("Should result in an error", func() {
				So(err.Error(), ShouldEqual, "code expected, none received")
			})

		})
		
		Convey("calling SetAppCode with a value", func() {
			err := SetAppCode("sup")
			
			Convey("Should not result in an error", func() {
				So(err, ShouldBeNil)
			})

			Convey("Should result in the appCode being set", func() {
				So(appCode, ShouldEqual, "sup")
			})

		})

	})

}

func TestSetConfigFile(t *testing.T) {
        
        Convey("Given an excuse to use this convey", t, func() {
                
                Convey("calling SetConfigFile with an empty value", func() {
                        err := SetConfigFile("")
                        
                        Convey("Should result in an error", func() {
                                So(err.Error(), ShouldEqual, "filename expected, none received")
                        })

                })      

                Convey("calling SetConfigFile with a value", func() {
                        err := SetConfigFile("sup")
                        
                        Convey("Should not result in an error", func() {
                                So(err, ShouldBeNil)
                        })

                        Convey("Should result in the configFilename being set", func() {
                                So(configFilename, ShouldEqual, "sup")
                        })

                })

        })

}

func TestSetConfigFormat(t *testing.T) {
        
        Convey("Given an excuse to use this convey", t, func() {
                
                Convey("calling SetConfigFormat with an empty value", func() {
                        err := SetConfigFormat("")
                        
                        Convey("Should result in an error", func() {
                                So(err.Error(), ShouldEqual, "config format was expected, none received")
                        })

                })      

                Convey("calling SetConfigFormat with json", func() {
			configFormat = ""
                        err := SetConfigFormat("json")
                        
                        Convey("Should not result in an error", func() {
                                So(err, ShouldBeNil)
                        })

                        Convey("Should result in the configFormat being set", func() {
                                So(configFormat, ShouldEqual, "json")
                        })

                })

                Convey("calling SetConfigFormat with toml", func() {
			configFormat = ""
                        err := SetConfigFormat("toml")
                        
                        Convey("Should not result in an error", func() {
                                So(err, ShouldBeNil)
                        })

                        Convey("Should result in the configFormat being set", func() {
                                So(configFormat, ShouldEqual, "toml")
                        })

                })

                Convey("calling SetConfigFormat with an invalid format: jso", func() {
			configFormat = ""
                        err := SetConfigFormat("jso")
                        
                        Convey("Should result in an error", func() {
                                So(err.Error(), ShouldEqual, "jso is not a supported configuration format")
                        })

                        Convey("Should result in the configFormat not being set", func() {
                                So(configFormat, ShouldEqual, "")
                        })

                })

        })

}

func TestsetConfigFormat(t *testing.T) {
        
        Convey("Given an excuse to use this convey", t, func() {
                
                Convey("calling setConfigFormat with the name not set", func() {
                        err := setConfigFormat()
                        
                        Convey("Should result in an error", func() {
                                So(err.Error(), ShouldEqual, "unable to determine config format, filename not set")
                        })

                })      

		Convey("setting the configFilename to a name without an extension", func() {
			configFilename = "config"
	
	                Convey("calling setConfigFormat", func() {
        	                err := setConfigFormat()
                        
	        	        Convey("Should result in an error", func() {
	                                So(err.Error(), ShouldEqual, "unable to determine config format, the configuration file config doesn't have an extension")
	                        })

        	        })

		})

		Convey("setting configFilename with config.toml", func() {
			configFilename = "config.toml"

		        Convey("calling setConfigFormat", func() {
		                err := setConfigFormat()
		                
		                Convey("Should not result in an error", func() {
		                        So(err, ShouldBeNil)
		                })

		                Convey("Should result in the configFormat being set", func() {
		                        So(configFormat, ShouldEqual, "toml")
		                })

		        })
		})

		Convey("setting configFilename with config.json", func() {
			configFilename = "config.json"

		        Convey("calling setConfigFormat", func() {
		                err := setConfigFormat()
		                
		                Convey("Should not result in an error", func() {
		                        So(err, ShouldBeNil)
		                })

		                Convey("Should result in the configFormat being set", func() {
		                        So(configFormat, ShouldEqual, "json")
		                })

		        })
		})

		Convey("setting configFilename with config.jso", func() {
			configFilename = "config.jso"

		        Convey("calling setConfigFormat", func() {
		                err := setConfigFormat()
		                
		                Convey("Should result in an error", func() {
		                        So(err, ShouldEqual, "jso is not a supported configuration format")
		                })

		                Convey("Should result in the configFormat not being set", func() {
		                        So(configFormat, ShouldEqual, "")
		                })

		        })

		})

        })

}

func TestIsSupportedFormat(t *testing.T) {
        
        Convey("Given an excuse to use this convey", t, func() {
                
                Convey("calling isSupportedFormat with an empty value", func() {
                        err := isSupportedFormat("")
                        
                        Convey("Should result in an error", func() {
                                So(err.Error(), ShouldEqual, " is not a supported configuration format")
                        })

                })      

                Convey("calling isSupportedFormat with json", func() {
			configFormat = ""
                        err := isSupportedFormat("json")
                        
                        Convey("Should not result in an error", func() {
                                So(err, ShouldBeNil)
                        })

                        Convey("Should result in the configFormat being set", func() {
                                So(configFormat, ShouldEqual, "json")
                        })

                })

                Convey("calling isSupportedFormat with toml", func() {
			configFormat = ""
                        err := isSupportedFormat("toml")
                        
                        Convey("Should not result in an error", func() {
                                So(err, ShouldBeNil)
                        })

                        Convey("Should result in the configFormat being set", func() {
                                So(configFormat, ShouldEqual, "toml")
                        })

                })

                Convey("calling isSupportedFormat with an invalid format: jso", func() {
			configFormat = ""
                        err := isSupportedFormat("jso")
                        
                        Convey("Should result in an error", func() {
                                So(err.Error(), ShouldEqual, "jso is not a supported configuration format")
                        })

                        Convey("Should result in the configFormat being set", func() {
                                So(configFormat, ShouldEqual, "")
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
