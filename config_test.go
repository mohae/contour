package contour

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAppConfig(t *testing.T) {
	initConfigs()
	Convey("Given some config tests", t, func() { 
		Convey("calling AppConfig", func() {
			c := AppConfig()
			Convey("should result in an appConfig", func() {
				So(toString.Get(c), ShouldEqual, toString.Get(testConfigs[app]))
			})
		})

		Convey("calling AppConfig", func() {
			c := AppConfig()
			Convey("should result in a config object", func() {
				So(toString.Get(c), ShouldEqual, toString.Get(testConfigs[app]))
			})

		})

	})

}

func TestConfig(t *testing.T) {
	initConfigs()
	Convey("Given an empty Config", t, func() { 
		Convey("calling Config", func() {
			_, err := Config("test")
			Convey("should result in an error", func() {
				So(err.Error(), ShouldEqual, "test config was requested; it does not exist")
			})
		})
	})

	Convey("Given an configs with an config", t, func() { 
		Convey("calling AppConfig", func() {
			NewConfig("test")
			c, err := Config("test")
			Convey("should not result in an error", func() {
				So(err, ShouldBeNil)
			})
			Convey("should result in a config object", func() {
				So(toString.Get(c), ShouldEqual, toString.Get(testConfigs[app]))
			})

		})

	})

}

// Things seem wonky with NewConfig. It appears to be returning that it already exists.
// Even though it doesn't. Returning the err processing to include the returned c, which
// is an empty *config, according to spec, allows the tests to pass.
//
// Running this code in a separate example code shows that it is working properly,
// i.e. creating a new *config, registering it in configs, and returning it w/o error.
func TestNewconfig(t *testing.T) {
	initConfigs()
	Convey("Given an empty Config", t, func() { 
		Convey("calling NewConfig", func() {
			c, err := NewConfig("test")
			Convey("should not result in an error", func() {
				So(err, ShouldBeNil)
			})
			Convey("should result in a config object", func() {
				So(toString.Get(c), ShouldEqual, toString.Get(testConfigs[app]))
			})
		})

		Convey("calling NewConfig", func() {
			c, err := NewConfig("test2")
			Convey("should not result in an error", func() {
				So(err, ShouldBeNil)
			})
			Convey("should result in a config object", func() {
				So(toString.Get(c), ShouldEqual, toString.Get(testConfigs[app]))
			})
		})

	})

	Convey("Given an Config with a 'test' config", t, func() { 
		Convey("calling NewConfig", func() {
			_, err := NewConfig("test")
			Convey("should result in an error", func() {
				So(err.Error(), ShouldEqual, "unable to create a new config for test, it already exists")
			})
		})
	})

}

