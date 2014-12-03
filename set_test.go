package contour

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSetBoolE(t *testing.T) {
	Convey("Given a new config", t, func() {
		c := NewCfg("setboole")

		Convey("Setting a boolean flag", func() {
			err := c.SetFlagBoolE("flagBoolKey", "b", true, true, "")
			Convey("Should not error", func() {
				So(err, ShouldBeNil)
			})

			Convey("And the setting should exist", func() {
				s, ok := c.settings["flagBoolKey"]

				Convey("And the key should exist", func() {
					So(ok, ShouldEqual, true)
				})

				Convey("And the flag should be set", func() {
					So(s.Value.(bool), ShouldEqual, true)
				})

				Convey("And the code should be set", func() {
					So(s.Short, ShouldEqual, "b")
				})

				Convey("And the type should be bool", func() {
					So(s.Type, ShouldEqual, "bool")
				})

				Convey("And it should not be flagged as a core setting", func() {
					So(s.IsCore, ShouldEqual, false)
				})

				Convey("And it should be flagged as a config", func() {
					So(s.IsCfg, ShouldEqual, true)
				})

				Convey("And it should be flagged as a flag", func() {
					So(s.IsFlag, ShouldEqual, true)
				})
			})
		})
	})
}

func TestSetIntE(t *testing.T) {
	Convey("Given a test config", t, func() {
		c := NewCfg("setinte")

		Convey("Setting an int flag", func() {
			err := c.SetFlagIntE("flagIntKey", "i", 42, 0, "")
			Convey("Should not error", func() {
				So(err, ShouldBeNil)
			})

			Convey("And the setting should exist", func() {
				s, ok := c.settings["flagIntKey"]

				Convey("And the key should exist", func() {
					So(ok, ShouldEqual, true)
				})

				Convey("And the flag should be set", func() {
					So(s.Value.(int), ShouldEqual, 42)
				})

				Convey("And the code should be set", func() {
					So(s.Short, ShouldEqual, "i")
				})

				Convey("And the type should be int", func() {
					So(s.Type, ShouldEqual, "int")
				})

				Convey("And it should not be flagged as a core setting", func() {
					So(s.IsCore, ShouldEqual, false)
				})

				Convey("And it should be flagged as a config", func() {
					So(s.IsCfg, ShouldEqual, true)
				})

				Convey("And it should be flagged as a flag", func() {
					So(s.IsFlag, ShouldEqual, true)
				})
			})
		})
	})
}

func TestSetStringE(t *testing.T) {
	Convey("Given a test config", t, func() {
		c := NewCfg("setstringe")

		Convey("Setting a string flag", func() {
			err := c.SetFlagStringE("flagStringKey", "s", "marvin", "", "")
			Convey("Should not error", func() {
				So(err, ShouldBeNil)
			})

			Convey("And the setting should exist", func() {
				s, ok := c.settings["flagStringKey"]

				Convey("And the key should exist", func() {
					So(ok, ShouldEqual, true)
				})

				Convey("And the flag should be set", func() {
					So(s.Value.(string), ShouldEqual, "marvin")
				})

				Convey("And the code should be set", func() {
					So(s.Short, ShouldEqual, "s")
				})

				Convey("And the type should be string", func() {
					So(s.Type, ShouldEqual, "string")
				})

				Convey("And it should not be flagged as a core setting", func() {
					So(s.IsCore, ShouldEqual, false)
				})

				Convey("And it should be flagged as a config", func() {
					So(s.IsCfg, ShouldEqual, true)
				})

				Convey("And it should be flagged as a flag", func() {
					So(s.IsFlag, ShouldEqual, true)
				})
			})
		})
	})
}

func TestSetBool(t *testing.T) {
	Convey("Given a new config", t, func() {
		c := NewCfg("setbool")

		Convey("Setting a boolean flag", func() {
			c.SetFlagBool("flagBoolKey", "b", true, true, "")

			Convey("And the setting should exist", func() {
				s, ok := c.settings["flagBoolKey"]

				Convey("And the key should exist", func() {
					So(ok, ShouldEqual, true)
				})

				Convey("And the flag should be set", func() {
					So(s.Value.(bool), ShouldEqual, true)
				})

				Convey("And the code should be set", func() {
					So(s.Short, ShouldEqual, "b")
				})

				Convey("And the type should be bool", func() {
					So(s.Type, ShouldEqual, "bool")
				})

				Convey("And it should not be flagged as a core setting", func() {
					So(s.IsCore, ShouldEqual, false)
				})

				Convey("And it should be flagged as a config", func() {
					So(s.IsCfg, ShouldEqual, true)
				})

				Convey("And it should be flagged as a flag", func() {
					So(s.IsFlag, ShouldEqual, true)
				})
			})
		})
	})
}

func TestSetInt(t *testing.T) {
	Convey("Given a test config", t, func() {
		c := NewCfg("setint")

		Convey("Setting an int flag", func() {
			c.SetFlagInt("flagIntKey", "i", 42, 42, "")

			Convey("And the setting should exist", func() {
				s, ok := c.settings["flagIntKey"]

				Convey("And the key should exist", func() {
					So(ok, ShouldEqual, true)
				})

				Convey("And the flag should be set", func() {
					So(s.Value.(int), ShouldEqual, 42)
				})

				Convey("And the code should be set", func() {
					So(s.Short, ShouldEqual, "i")
				})

				Convey("And the type should be int", func() {
					So(s.Type, ShouldEqual, "int")
				})

				Convey("And it should not be flagged as a core setting", func() {
					So(s.IsCore, ShouldEqual, false)
				})

				Convey("And it should be flagged as a config", func() {
					So(s.IsCfg, ShouldEqual, true)
				})

				Convey("And it should be flagged as a flag", func() {
					So(s.IsFlag, ShouldEqual, true)
				})
			})
		})
	})
}

func TestSetString(t *testing.T) {
	Convey("Given a test config", t, func() {
		c := NewCfg("setstring")

		Convey("Setting a string flag", func() {
			c.SetFlagStringE("flagStringKey", "s", "marvin", "", "")

			Convey("And the setting should exist", func() {
				s, ok := c.settings["flagStringKey"]

				Convey("And the key should exist", func() {
					So(ok, ShouldEqual, true)
				})

				Convey("And the flag should be set", func() {
					So(s.Value.(string), ShouldEqual, "marvin")
				})

				Convey("And the code should be set", func() {
					So(s.Short, ShouldEqual, "s")
				})

				Convey("And the type should be string", func() {
					So(s.Type, ShouldEqual, "string")
				})

				Convey("And it should not be flagged as a core setting", func() {
					So(s.IsCore, ShouldEqual, false)
				})

				Convey("And it should be flagged as a config", func() {
					So(s.IsCfg, ShouldEqual, true)
				})

				Convey("And it should be flagged as a flag", func() {
					So(s.IsFlag, ShouldEqual, true)
				})
			})
		})
	})
}
