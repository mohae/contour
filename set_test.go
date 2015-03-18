package contour

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSetBoolE(t *testing.T) {
	c := NewCfg("setboole")
	err := c.SetFlagBoolE("flagBoolKey", "b", true, "true", "")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err.Error())
		return
	}

	s, ok := c.settings["flagBoolKey"]
	if !ok {
		t.Errorf("Expected 'ok' to be true, got %t", ok)
	} else {
		var b bool
		switch s.Value.(type) {
		case bool:
			b = s.Value.(bool)
		case *bool:
			b = *s.Value.(*bool)
		}
		if !b {
			t.Errorf("Expected \"true\", got %t", b)
		}

		if s.Short != "b" {
			t.Errorf("Expected short to be \"b\", got %q", s.Short)
		}

		if s.Type != "bool" {
			t.Errorf("Expected type to be \"bool\", got %q", s.Type)
		}

		if s.IsCore {
			t.Errorf("Expected IsCore to be false, got %t", s.IsCore)
		}

		if !s.IsCfg {
			t.Errorf("Expected IsCfg to be true, got %t", s.IsCfg)
		}

		if !s.IsFlag {
			t.Errorf("Expected IsCfg to be true, got %t", s.IsFlag)
		}
	}
}

func TestSetIntE(t *testing.T) {
	Convey("Given a test config", t, func() {
		c := NewCfg("setinte")

		Convey("Setting an int flag", func() {
			err := c.SetFlagIntE("flagIntKey", "i", 42, "42", "")
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
	c := NewCfg("setbool")
	c.SetFlagBool("flagBoolKey", "b", true, "true", "")

	s, ok := c.settings["flagBoolKey"]
	if !ok {
		t.Errorf("Expected 'ok' to be true, got %t", ok)
	} else {
		var b bool
		switch s.Value.(type) {
		case bool:
			b = s.Value.(bool)
		case *bool:
			b = *s.Value.(*bool)
		}
		if !b {
			t.Errorf("Expected \"true\", got %t", b)
		}

		if s.Short != "b" {
			t.Errorf("Expected short to be \"b\", got %q", s.Short)
		}

		if s.Type != "bool" {
			t.Errorf("Expected type to be \"bool\", got %q", s.Type)
		}

		if s.IsCore {
			t.Errorf("Expected IsCore to be false, got %t", s.IsCore)
		}

		if !s.IsCfg {
			t.Errorf("Expected IsCfg to be true, got %t", s.IsCfg)
		}

		if !s.IsFlag {
			t.Errorf("Expected IsCfg to be true, got %t", s.IsFlag)
		}
	}
}

func TestSetInt(t *testing.T) {
	Convey("Given a test config", t, func() {
		c := NewCfg("setint")

		Convey("Setting an int flag", func() {
			c.SetFlagInt("flagIntKey", "i", 42, "42", "")

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
