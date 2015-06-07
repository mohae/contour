package contour

import (
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type iTest struct {
	name        string
	key         string
	value       interface{}
	expected    string
	expectedErr string
}

func NewITests() []iTest {
	tests := make([]iTest, 12)
	tests[0] = iTest{"corebool", "corebool", true, "", "corebool is not a flag: only flags can be overridden"}
	tests[1] = iTest{"coreint", "coreint", 42, "", "coreint is not a flag: only flags can be overridden"}
	tests[2] = iTest{"corestring", "corestring", "beeblebrox", "", "corestring is not a flag: only flags can be overridden"}
	tests[3] = iTest{"cfgbool", "cfgbool", true, "", "cfgbool is not a flag: only flags can be overridden"}
	tests[4] = iTest{"cfgint", "cfgint", 43, "", "cfgint is not a flag: only flags can be overridden"}
	tests[5] = iTest{"cfgstring", "cfgstring", "frood", "", "cfgstring is not a flag: only flags can be overridden"}
	tests[6] = iTest{"flagbool", "flagbool", true, "true", ""}
	tests[7] = iTest{"flagint", "flagint", 41, "41", ""}
	tests[8] = iTest{"flagstring", "flagstring", "towel", "towel", ""}
	tests[9] = iTest{"bool", "bool", true, "", "bool is not a flag: only flags can be overridden"}
	tests[10] = iTest{"int", "int", 3, "", "int is not a flag: only flags can be overridden"}
	tests[11] = iTest{"string", "string", "don't panic", "dont' panic", "string is not a flag: only flags can be overridden"}
	return tests
}

func TestOverride(t *testing.T) {
	testCfg := newTestCfg()
	tests := NewITests()
	for _, test := range tests {
		Convey("Given a test and a testCfg", t, func() {
			Convey(test.name+": Overridding "+test.key, func() {
				err := testCfg.Override(test.key, test.value)
				if err != nil {
					if test.expectedErr != "" {
						Convey("Should result in an error", func() {
							So(err.Error(), ShouldEqual, test.expectedErr)
						})
					} else {
						Convey("Should not error", func() {
							So(err, ShouldBeNil)
						})
					}
				} else {
					if test.expectedErr != "" {
						Convey("Should not result in an error", func() {
							So(err.Error(), ShouldEqual, test.expectedErr)
						})
					} else {
						Convey("Should result in expected", func() {
							i := testCfg.Get(test.key)
							switch i.(type) {
							case bool:
								So(strconv.FormatBool(i.(bool)), ShouldEqual, test.expected)
							case int:
								So(strconv.Itoa(i.(int)), ShouldEqual, test.expected)
							case string:
								So(i.(string), ShouldEqual, test.expected)
							}
						})
					}
				}
			})
		})
	}
}
