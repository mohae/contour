package contour

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetsE(t *testing.T) {
	Convey("Given a configuration with settings", t, func() {
		Convey("Getting an interface", func() {
			rif, err :=  testConfig.GetE("corebool")
			Convey("Should not result in an error", func() {
				So(err, ShouldBeNil)
			})
			Convey("Should result in a value", func() {
				So(rif.(bool), ShouldEqual, true)
			})
		})
 
		Convey("Getting a bool", func() {
			rb, err :=  testConfig.GetBoolE("corebool")
			Convey("Should not result in an error", func() {
				So(err, ShouldBeNil)
			})
			Convey("Should result in a value", func() {
				So(rb, ShouldEqual, true)
			})
		}) 

		Convey("Getting an int", func() {
			ri, err :=  testConfig.GetIntE("coreint")
			Convey("Should not result in an error", func() {
				So(err, ShouldBeNil)
			})
			Convey("Should result in a value", func() {
				So(ri, ShouldEqual, 42)
			})
		})

		Convey("Getting a string", func() {
			rb, err :=  testConfig.GetStringE("corestring")
			Convey("Should not result in an error", func() {
				So(err, ShouldBeNil)
			})
			Convey("Should result in a value", func() {
				So(rb, ShouldEqual, "a core string")
			})
		}) 

		Convey("Getting an interface", func() {
			rinter, err :=  testConfig.GetInterfaceE("corebool")
			Convey("Should not result in an error", func() {
				So(err, ShouldBeNil)
			})
			Convey("Should result in a value", func() {
				So(rinter.(bool), ShouldEqual, true)
			})
		})
	})
}

func TestGets(t *testing.T) {
	Convey("Given a configuration with settings", t, func() {
		Convey("Getting an interface", func() {
			rif :=  testConfig.Get("corebool")
			Convey("Should result in a value", func() {
				So(rif.(bool), ShouldEqual, true)
			})
		})
 
		Convey("Getting a bool", func() {
			rb :=  testConfig.GetBool("corebool")
			Convey("Should result in a value", func() {
				So(rb, ShouldEqual, true)
			})
		}) 

		Convey("Getting an int", func() {
			ri :=  testConfig.GetInt("coreint")
			Convey("Should result in a value", func() {
				So(ri, ShouldEqual, 42)
			})
		})

		Convey("Getting a string", func() {
			rb :=  testConfig.GetString("corestring")
			Convey("Should result in a value", func() {
				So(rb, ShouldEqual, "a core string")
			})
		}) 

		Convey("Getting an interface", func() {
			rinter :=  testConfig.GetInterface("corebool")
			Convey("Should result in a value", func() {
				So(rinter.(bool), ShouldEqual, true)
			})
		})
	})
}

func TestGetFilterNames(t * testing.T) {
	Convey("Given a configuration with settings", t, func() {
		Convey("Getting a list of Bool Filters", func() {
			boolFilters := testConfig.GetBoolFilterNames()
			Convey("Should result in a list of bool filters for this config", func() {
				So(toString.Get(boolFilters), ShouldEqual, "[\"flagbool\"]")
			})
		})

		Convey("Getting a list of int Filters", func() {
			intFilters := testConfig.GetIntFilterNames()
			Convey("Should result in a list of bool filters for this config", func() {
				So(toString.Get(intFilters), ShouldEqual, "[\"flagint\"]")
			})
		})

		Convey("Getting a list of string Filters", func() {
			stringFilters := testConfig.GetStringFilterNames()
			Convey("Should result in a list of string filters for this config", func() {
				So(toString.Get(stringFilters), ShouldEqual, "[\"flagstring\"]")
			})
		})

	})
}
