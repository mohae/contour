package contour

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfig(t *testing.T) {
	initConfigs()
	Convey("Calling NewConfig", t, func() {
		c := NewConfig("test")
		Convey("should result in an config with name 'test'", func() {
			So(c.Code(), ShouldEqual, "")
		})
	})

}
