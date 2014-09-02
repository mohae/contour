package contour

import (
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

/*
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

*/

