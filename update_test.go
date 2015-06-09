package contour

import (
	"fmt"
	"strconv"
	"testing"
)

func TestUpdateBools(t *testing.T) {
	bTests := []struct {
		key   string
		value bool
		err   string
	}{
		{"", false, "cannot update \"\": setting not found"},
		{"corebool", false, "cannot update \"corebool\": core settings cannot be updated"},
		{"corebool", true, "cannot update \"corebool\": core settings cannot be updated"},
		{"flagbool", false, ""},
		{"flagbool", true, ""},
		{"cfgbool", false, ""},
		{"cfgbool", true, ""},
		{"bool", false, ""},
		{"bool", true, ""},
	}
	appCfg = newTestCfg()
	for i, test := range bTests {
		err := UpdateBoolE(test.key, test.value)
		if err != nil {
			if test.err != err.Error() {
				t.Errorf("%d: expected %q got %q", i, test.err, err.Error())
			}
			continue
		}
		b, err := GetBoolE(test.key)
		if err != nil {
			if test.err != err.Error() {
				t.Errorf("%d: expected %q got %q", i, test.err, err.Error())
			}
			continue
		}
		if b != test.value {
			t.Errorf("%d: expected %v got %v", i, test.value, b)
		}
		// Non-E
		UpdateBool(test.key, false)
		b = GetBool(test.key)
		if b != false {
			t.Errorf("%d: expected false got %v", i, b)
		}
	}
}

func TestUpdateInts(t *testing.T) {
	iTests := []struct {
		key   string
		value int
		err   string
	}{
		{"", 0, "cannot update \"\": setting not found"},
		{"coreint", 42, "cannot update \"coreint\": core settings cannot be updated"},
		{"flagint", 42, ""},
		{"cfgint", 42, ""},
		{"int", 42, ""},
	}
	appCfg = newTestCfg()
	for i, test := range iTests {
		err := UpdateIntE(test.key, test.value)
		if err != nil {
			if test.err != err.Error() {
				t.Errorf("%ds: expected %q got %q", i, test.err, err)
			}
			continue
		}
		i, err := GetIntE(test.key)
		if err != nil {
			if test.err != err.Error() {
				t.Errorf("%d: expected %q got %q", i, test.err, err)
			}
			continue
		}
		if i != test.value {
			t.Errorf("%d: expected %q got %q", i, test.value, strconv.Itoa(i))
		}
		// Non-e
		UpdateInt(test.key, test.value+10)
		i = GetInt(test.key)
		if i != test.value+10 {
			t.Errorf("%d: expected %v got %v", i, test.value+10, i)
		}
	}
}

func TestUpdateInt64s(t *testing.T) {
	i64Tests := []struct {
		key   string
		value int64
		err   string
	}{
		{"coreint64", int64(42), "cannot update \"coreint64\": core settings cannot be updated"},
		{"", int64(0), "cannot update \"\": setting not found"},
		{"flagint64", int64(42), ""},
		{"cfgint64", int64(42), ""},
		{"int", int64(42), ""},
	}
	appCfg = newTestCfg()
	for i, test := range i64Tests {
		err := UpdateInt64E(test.key, test.value)
		if err != nil {
			if test.err != err.Error() {
				t.Errorf("%d: expected %q got %q", i, test.err, err.Error())
			}
			continue
		}
		i64, err := GetInt64E(test.key)
		if err != nil {
			if test.err != err.Error() {
				t.Errorf("%d: expected %q got %q", i, test.err, err.Error())
			}
			continue
		}
		if i64 != test.value {
			t.Errorf("%d: expected %v got %v", i, test.value, i64)
		}
		// Non-e
		UpdateInt64(test.key, test.value+int64(10))
		i64 = GetInt64(test.key)
		if i64 != test.value+int64(10) {
			t.Errorf("%d: expected %v got %v", i, test.value+int64(10), i64)
		}
	}
}

func TestUpdateStrings(t *testing.T) {
	sTests := []struct {
		key   string
		value string
		err   string
	}{
		{"", "false", "cannot update \"\": setting not found"},
		{"corestring", "false", "cannot update \"corestring\": core settings cannot be updated"},
		{"corestring", "t", "cannot update \"corestring\": core settings cannot be updated"},
		{"flagstring", "false", ""},
		{"flagstring", "t", ""},
		{"cfgstring", "false", ""},
		{"cfgstring", "t", ""},
		{"string", "false", ""},
		{"string", "t", ""},
	}
	appCfg = newTestCfg()
	for i, test := range sTests {
		err := UpdateStringE(test.key, test.value)
		if err != nil {
			if test.err != err.Error() {
				t.Errorf("%d: expected %q got %q", i, test.err, err.Error())
			}
			continue
		}
		s, err := GetStringE(test.key)
		if err != nil {
			if test.err != err.Error() {
				t.Errorf("%d: expected %q got %q", i, test.err, err.Error())
			}
			continue
		}
		if s != test.value {
			t.Errorf("%d: expected %s got %s", i, test.value, s)
		}
		// Non-e
		UpdateString(test.key, fmt.Sprintf("%s %s", test.value, test.value))
		s = GetString(test.key)
		if s != fmt.Sprintf("%s %s", test.value, test.value) {
			t.Errorf("%d: expected %v got %v", i, fmt.Sprintf("%s %s", test.value, test.value), s)
		}
	}
}
