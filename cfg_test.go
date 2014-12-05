package contour

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfig(t *testing.T) {
	initCfgs()
	Convey("Calling NewConfig", t, func() {
		c := NewCfg("test")
		Convey("should result in an config with name 'test'", func() {
			So(c.name, ShouldEqual, "test")
		})
	})

}
