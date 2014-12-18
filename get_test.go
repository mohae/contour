package contour

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetsE(t *testing.T) {
	Convey("Given a cfguration with settings", t, func() {
		Convey("Getting an interface", func() {
			rif, err := testCfg.GetE("corebool")
			Convey("Should not result in an error", func() {
				So(err, ShouldBeNil)
			})
			Convey("Should result in a value", func() {
				So(rif.(string), ShouldEqual, "true")
			})
		})

		Convey("Getting a bool", func() {
			rb, err := testCfg.GetBoolE("corebool")
			Convey("Should not result in an error", func() {
				So(err, ShouldBeNil)
			})
			Convey("Should result in a value", func() {
				So(rb, ShouldEqual, "true")
			})
		})

		Convey("Getting an int", func() {
			ri, err := testCfg.GetIntE("coreint")
			Convey("Should not result in an error", func() {
				So(err, ShouldBeNil)
			})
			Convey("Should result in a value", func() {
				So(ri, ShouldEqual, 42)
			})
		})

		Convey("Getting a string", func() {
			rb, err := testCfg.GetStringE("corestring")
			Convey("Should not result in an error", func() {
				So(err, ShouldBeNil)
			})
			Convey("Should result in a value", func() {
				So(rb, ShouldEqual, "a core string")
			})
		})

		Convey("Getting an interface", func() {
			rinter, err := testCfg.GetInterfaceE("corebool")
			Convey("Should not result in an error", func() {
				So(err, ShouldBeNil)
			})
			Convey("Should result in a value", func() {
				So(rinter.(string), ShouldEqual, "true")
			})
		})
	})
}

func TestGets(t *testing.T) {
	Convey("Given a cfguration with settings", t, func() {
		Convey("Getting an interface", func() {
			rif := testCfg.Get("corebool")
			Convey("Should result in a value", func() {
				So(rif.(string), ShouldEqual, "true")
			})
		})

		Convey("Getting a bool", func() {
			rb := testCfg.GetBool("corebool")
			Convey("Should result in a value", func() {
				So(rb, ShouldEqual, "true")
			})
		})

		Convey("Getting an int", func() {
			ri := testCfg.GetInt("coreint")
			Convey("Should result in a value", func() {
				So(ri, ShouldEqual, 42)
			})
		})

		Convey("Getting a string", func() {
			rb := testCfg.GetString("corestring")
			Convey("Should result in a value", func() {
				So(rb, ShouldEqual, "a core string")
			})
		})

		Convey("Getting an interface", func() {
			rinter := testCfg.GetInterface("corebool")
			Convey("Should result in a value", func() {
				So(rinter.(string), ShouldEqual, "true")
			})
		})
	})
}

func TestGetFilterNames(t *testing.T) {
	Convey("Given a cfguration with settings", t, func() {
		Convey("Getting a list of Bool Filters", func() {
			boolFilters := testCfg.GetBoolFilterNames()
			Convey("Should result in a list of bool filters for this cfg", func() {
				So(toString.Get(boolFilters), ShouldEqual, "[\"flagbool\"]")
			})
		})

		Convey("Getting a list of int Filters", func() {
			intFilters := testCfg.GetIntFilterNames()
			Convey("Should result in a list of bool filters for this cfg", func() {
				So(toString.Get(intFilters), ShouldEqual, "[\"flagint\"]")
			})
		})

		Convey("Getting a list of string Filters", func() {
			stringFilters := testCfg.GetStringFilterNames()
			Convey("Should result in a list of string filters for this cfg", func() {
				So(toString.Get(stringFilters), ShouldEqual, "[\"flagstring\"]")
			})
		})

	})
}
