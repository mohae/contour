package contour

import (
	"strconv"
	"testing"
)

func TestUpdateBools(t *testing.T) {
	bTests := []struct {
		key   string
		value bool
		err   string
	}{
		{"corebool", false, "cannot update \"corebool\": core settings cannot be updated"},
		{"corebool", true, "cannot update \"corebool\": core settings cannot be updated"},
		{"flagbool", false, ""},
		{"flagbool", true, ""},
		{"cfgbool", false, ""},
		{"cfgbool", true, ""},
		{"bool", false, ""},
		{"bool", true, ""},
	}
	testCfg := newTestCfg()
	testCfg.name = "update"
	for i, test := range bTests {
		err := testCfg.UpdateBoolE(test.key, test.value)
		if err != nil {
			if test.err != err.Error() {
				t.Errorf("%d: expected %q got %q", i, test.err, err)
			}
			continue
		}
		b, err := testCfg.GetBoolE(test.key)
		if err != nil {
			if test.err != err.Error() {
				t.Errorf("%d: expected %q got %q", i, test.err, err)
			}
			continue
		}
		if b != test.value {
			t.Errorf("%d: expected %t got %t", i, test.value, b)
		}
	}
}

func TestUpdateStrings(t *testing.T) {
	sTests := []struct {
		key   string
		value string
		err   string
	}{
		{"", "false", "cannot update \"\": not found"},
		{"corestring", "false", "cannot update \"corestring\": core settings cannot be updated"},
		{"corestring", "t", "cannot update \"corestring\": core settings cannot be updated"},
		{"flagstring", "false", ""},
		{"flagstring", "t", ""},
		{"cfgstring", "false", ""},
		{"cfgstring", "t", ""},
		{"string", "false", ""},
		{"string", "t", ""},
	}
	testCfg := newTestCfg()
	testCfg.name = "update"
	for i, test := range sTests {
		err := testCfg.UpdateStringE(test.key, test.value)
		if err != nil {
			if test.err != err.Error() {
				t.Errorf("%d: expected %q got %q", i, test.err, err)
			}
			continue
		}
		s, err := testCfg.GetStringE(test.key)
		if err != nil {
			if test.err != err.Error() {
				t.Errorf("%d: expected %q got %q", i, test.err, err)
			}
			continue
		}
		if s != test.value {
			t.Errorf("%d: expected %t got %s", i, test.value, s)
		}
	}
}

func TestUpdateInts(t *testing.T) {
	iTests := []struct {
		key   string
		value int
		err   string
	}{
		{"coreint", 42, "cannot update \"coreint\": core settings cannot be updated"},
		{"", 0, "cannot update \"\": not found"},
		{"flagint", 42, ""},
		{"cfgint", 42, ""},
		{"int", 42, ""},
	}
	testCfg := newTestCfg()
	testCfg.name = "update"
	for i, test := range iTests {
		err := testCfg.UpdateIntE(test.key, test.value)
		if err != nil {
			if test.err != err.Error() {
				t.Errorf("%ds: expected %q got %q", i, test.err, err)
			}
			continue
		}
		i, err := testCfg.GetIntE(test.key)
		if err != nil {
			if test.err != err.Error() {
				t.Errorf("%d: expected %q got %q", i, test.err, err)
			}
			continue
		}
		if i != test.value {
			t.Errorf("%d: expected %q got %q", i, test.value, strconv.Itoa(i))
		}
	}
}

func TestUpdateInt64s(t *testing.T) {
	i64Tests := []struct {
		key   string
		value int64
		err   string
	}{
		{"coreint", int64(42), "cannot update \"coreint\": core settings cannot be updated"},
		{"", int64(0), "cannot update \"\": not found"},
		{"flagint", int64(42), ""},
		{"cfgint", int64(42), ""},
		{"int", int64(42), ""},
	}
	testCfg := newTestCfg()
	testCfg.name = "update"
	for i, test := range i64Tests {
		err := testCfg.UpdateInt64E(test.key, test.value)
		if err != nil {
			if test.err != err.Error() {
				t.Errorf("%d: expected %q got %q", i, test.err, err)
			}
			continue
		}
		i64, err := testCfg.GetInt64E(test.key)
		if err != nil {
			if test.err != err.Error() {
				t.Errorf("%d: expected %q got %q", i, test.err, err)
			}
			continue
		}
		if i64 != test.value {
			t.Errorf("%d: expected %q got %q", i, test.value, strconv.Itoa(int(i64)))
		}
	}
}
