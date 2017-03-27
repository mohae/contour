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
		{"", false, ": setting not found"},
		{"corebool", false, "corebool: core settings cannot be updated"},
		{"corebool", true, "corebool: core settings cannot be updated"},
		{"flagbool", false, ""},
		{"flagbool", true, ""},
		{"cfgbool", false, ""},
		{"cfgbool", true, ""},
		{"bool", false, ""},
		{"bool", true, ""},
	}
	tstSettings := newTestSettings()
	for i, test := range bTests {
		err := tstSettings.UpdateBoolE(test.key, test.value)
		if err != nil {
			if test.err != err.Error() {
				t.Errorf("%d: expected %q got %q", i, test.err, err.Error())
			}
			continue
		}
		b, err := tstSettings.BoolE(test.key)
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
		tstSettings.UpdateBool(test.key, false)
		b = Bool(test.key)
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
		{"", 0, ": setting not found"},
		{"coreint", 42, "coreint: core settings cannot be updated"},
		{"flagint", 42, ""},
		{"cfgint", 42, ""},
		{"int", 42, ""},
	}
	tstSettings := newTestSettings()
	for i, test := range iTests {
		err := tstSettings.UpdateIntE(test.key, test.value)
		if err != nil {
			if test.err != err.Error() {
				t.Errorf("%ds: expected %q got %q", i, test.err, err)
			}
			continue
		}
		i, err := tstSettings.IntE(test.key)
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
		tstSettings.UpdateInt(test.key, test.value+10)
		i = tstSettings.Int(test.key)
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
		{"coreint64", int64(42), "coreint64: core settings cannot be updated"},
		{"", int64(0), ": setting not found"},
		{"flagint64", int64(42), ""},
		{"cfgint64", int64(42), ""},
		{"int", int64(42), ""},
	}
	tstSettings := newTestSettings()
	for i, test := range i64Tests {
		err := tstSettings.UpdateInt64E(test.key, test.value)
		if err != nil {
			if test.err != err.Error() {
				t.Errorf("%d: expected %q got %q", i, test.err, err.Error())
			}
			continue
		}
		i64, err := tstSettings.Int64E(test.key)
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
		tstSettings.UpdateInt64(test.key, test.value+int64(10))
		i64 = tstSettings.Int64(test.key)
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
		{"", "false", ": setting not found"},
		{"corestring", "false", "corestring: core settings cannot be updated"},
		{"corestring", "t", "corestring: core settings cannot be updated"},
		{"flagstring", "false", ""},
		{"flagstring", "t", ""},
		{"cfgstring", "false", ""},
		{"cfgstring", "t", ""},
		{"string", "false", ""},
		{"string", "t", ""},
	}
	tstSettings := newTestSettings()
	for i, test := range sTests {
		err := tstSettings.UpdateStringE(test.key, test.value)
		if err != nil {
			if test.err != err.Error() {
				t.Errorf("%d: expected %q got %q", i, test.err, err.Error())
			}
			continue
		}
		s, err := tstSettings.StringE(test.key)
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
		tstSettings.UpdateString(test.key, fmt.Sprintf("%s %s", test.value, test.value))
		s = tstSettings.String(test.key)
		if s != fmt.Sprintf("%s %s", test.value, test.value) {
			t.Errorf("%d: expected %v got %v", i, fmt.Sprintf("%s %s", test.value, test.value), s)
		}
	}
}

func TestCanUpdate(t *testing.T) {
	tests := []struct {
		name            string
		typ             SettingType
		confFileVarsSet bool
		envVarsSet      bool
		flagsParsed     bool
		expected        bool
	}{
		// 0
		{"corebool", Basic, false, false, false, false},
		{"corebool", Basic, true, false, false, false},
		{"corebool", Basic, false, true, false, false},
		{"corebool", Basic, true, true, false, false},
		{"corebool", Basic, false, false, true, false},
		// 5
		{"corebool", Basic, false, true, true, false},
		{"corebool", Basic, true, false, true, false},
		{"corebool", Basic, true, true, true, false},
		{"x-corebool", Basic, false, false, false, false},
		{"corebool", Core, false, false, false, false},
		// 10
		{"corebool", Core, true, false, false, false},
		{"corebool", Core, false, true, false, false},
		{"corebool", Core, true, true, false, false},
		{"corebool", Core, false, false, true, false},
		{"corebool", Core, false, true, true, false},
		// 15
		{"corebool", Core, true, false, true, false},
		{"corebool", Core, true, true, true, false},
		{"x-corebool", Core, false, false, false, false},
		{"corebool", ConfFileVar, false, false, false, false},
		{"corebool", ConfFileVar, true, false, false, false},
		// 20
		{"corebool", ConfFileVar, false, true, false, false},
		{"corebool", ConfFileVar, true, true, false, false},
		{"corebool", ConfFileVar, false, false, true, false},
		{"corebool", ConfFileVar, false, true, true, false},
		{"corebool", ConfFileVar, true, false, true, false},
		// 25
		{"corebool", ConfFileVar, true, true, true, false},
		{"x-corebool", ConfFileVar, false, false, false, false},
		{"corebool", Env, false, false, false, false},
		{"corebool", Env, true, false, false, false},
		{"corebool", Env, false, true, false, false},
		// 30
		{"corebool", Env, true, true, false, false},
		{"corebool", Env, false, false, true, false},
		{"corebool", Env, false, true, true, false},
		{"corebool", Env, true, false, true, false},
		{"corebool", Env, true, true, true, false},
		// 35
		{"x-corebool", Env, false, false, false, false},
		{"corebool", Flag, false, false, false, false},
		{"corebool", Flag, true, false, false, false},
		{"corebool", Flag, false, true, false, false},
		{"corebool", Flag, true, true, false, false},
		// 40
		{"corebool", Flag, false, false, true, false},
		{"corebool", Flag, false, true, true, false},
		{"corebool", Flag, true, false, true, false},
		{"corebool", Flag, true, true, true, false},
		{"x-corebool", Flag, false, false, false, false},
		// 45
		{"coreint", Basic, false, false, false, false},
		{"coreint", Basic, true, false, false, false},
		{"coreint", Basic, false, true, false, false},
		{"coreint", Basic, true, true, false, false},
		{"coreint", Basic, false, false, true, false},
		// 50
		{"coreint", Basic, false, true, true, false},
		{"coreint", Basic, true, false, true, false},
		{"coreint", Basic, true, true, true, false},
		{"x-coreint", Basic, false, false, false, false},
		{"coreint", Core, false, false, false, false},
		// 55
		{"coreint", Core, true, false, false, false},
		{"coreint", Core, false, true, false, false},
		{"coreint", Core, true, true, false, false},
		{"coreint", Core, false, false, true, false},
		{"coreint", Core, false, true, true, false},
		// 60
		{"coreint", Core, true, false, true, false},
		{"coreint", Core, true, true, true, false},
		{"x-coreint", Core, false, false, false, false},
		{"coreint", ConfFileVar, false, false, false, false},
		{"coreint", ConfFileVar, true, false, false, false},
		// 65
		{"coreint", ConfFileVar, false, true, false, false},
		{"coreint", ConfFileVar, true, true, false, false},
		{"coreint", ConfFileVar, false, false, true, false},
		{"coreint", ConfFileVar, false, true, true, false},
		{"coreint", ConfFileVar, true, false, true, false},
		// 70
		{"coreint", ConfFileVar, true, true, true, false},
		{"x-coreint", ConfFileVar, false, false, false, false},
		{"coreint", Env, false, false, false, false},
		{"coreint", Env, true, false, false, false},
		{"coreint", Env, false, true, false, false},
		// 75
		{"coreint", Env, true, true, false, false},
		{"coreint", Env, false, false, true, false},
		{"coreint", Env, false, true, true, false},
		{"coreint", Env, true, false, true, false},
		{"coreint", Env, true, true, true, false},
		// 80
		{"x-coreint", Env, false, false, false, false},
		{"coreint", Flag, false, false, false, false},
		{"coreint", Flag, true, false, false, false},
		{"coreint", Flag, false, true, false, false},
		{"coreint", Flag, true, true, false, false},
		// 85
		{"coreint", Flag, false, false, true, false},
		{"coreint", Flag, false, true, true, false},
		{"coreint", Flag, true, false, true, false},
		{"coreint", Flag, true, true, true, false},
		{"x-coreint", Flag, false, false, false, false},
		// 90
		{"coreint64", Basic, false, false, false, false},
		{"coreint64", Basic, true, false, false, false},
		{"coreint64", Basic, false, true, false, false},
		{"coreint64", Basic, true, true, false, false},
		{"coreint64", Basic, false, false, true, false},
		// 95
		{"coreint64", Basic, false, true, true, false},
		{"coreint64", Basic, true, false, true, false},
		{"coreint64", Basic, true, true, true, false},
		{"x-coreint64", Basic, false, false, false, false},
		{"coreint64", Core, false, false, false, false},
		// 100
		{"coreint64", Core, true, false, false, false},
		{"coreint64", Core, false, true, false, false},
		{"coreint64", Core, true, true, false, false},
		{"coreint64", Core, false, false, true, false},
		{"coreint64", Core, false, true, true, false},
		// 105
		{"coreint64", Core, true, false, true, false},
		{"coreint64", Core, true, true, true, false},
		{"x-coreint64", Core, false, false, false, false},
		{"coreint64", ConfFileVar, false, false, false, false},
		{"coreint64", ConfFileVar, true, false, false, false},
		// 110
		{"coreint64", ConfFileVar, false, true, false, false},
		{"coreint64", ConfFileVar, true, true, false, false},
		{"coreint64", ConfFileVar, false, false, true, false},
		{"coreint64", ConfFileVar, false, true, true, false},
		{"coreint64", ConfFileVar, true, false, true, false},
		// 115
		{"coreint64", ConfFileVar, true, true, true, false},
		{"x-coreint64", ConfFileVar, false, false, false, false},
		{"coreint64", Env, false, false, false, false},
		{"coreint64", Env, true, false, false, false},
		{"coreint64", Env, false, true, false, false},
		// 120
		{"coreint64", Env, true, true, false, false},
		{"coreint64", Env, false, false, true, false},
		{"coreint64", Env, false, true, true, false},
		{"coreint64", Env, true, false, true, false},
		{"coreint64", Env, true, true, true, false},
		// 125
		{"x-coreint64", Env, false, false, false, false},
		{"coreint64", Flag, false, false, false, false},
		{"coreint64", Flag, true, false, false, false},
		{"coreint64", Flag, false, true, false, false},
		{"coreint64", Flag, true, true, false, false},
		// 130
		{"coreint64", Flag, false, false, true, false},
		{"coreint64", Flag, false, true, true, false},
		{"coreint64", Flag, true, false, true, false},
		{"coreint64", Flag, true, true, true, false},
		{"x-coreint64", Flag, false, false, false, false},
		// 135
		{"corestring", Basic, false, false, false, false},
		{"corestring", Basic, true, false, false, false},
		{"corestring", Basic, false, true, false, false},
		{"corestring", Basic, true, true, false, false},
		{"corestring", Basic, false, false, true, false},
		// 140
		{"corestring", Basic, false, true, true, false},
		{"corestring", Basic, true, false, true, false},
		{"corestring", Basic, true, true, true, false},
		{"x-corestring", Basic, false, false, false, false},
		{"corestring", Core, false, false, false, false},
		// 145
		{"corestring", Core, true, false, false, false},
		{"corestring", Core, false, true, false, false},
		{"corestring", Core, true, true, false, false},
		{"corestring", Core, false, false, true, false},
		{"corestring", Core, false, true, true, false},
		// 150
		{"corestring", Core, true, false, true, false},
		{"corestring", Core, true, true, true, false},
		{"x-corestring", Core, false, false, false, false},
		{"corestring", ConfFileVar, false, false, false, false},
		{"corestring", ConfFileVar, true, false, false, false},
		// 155
		{"corestring", ConfFileVar, false, true, false, false},
		{"corestring", ConfFileVar, true, true, false, false},
		{"corestring", ConfFileVar, false, false, true, false},
		{"corestring", ConfFileVar, false, true, true, false},
		{"corestring", ConfFileVar, true, false, true, false},
		// 160
		{"corestring", ConfFileVar, true, true, true, false},
		{"x-corestring", ConfFileVar, false, false, false, false},
		{"corestring", Env, false, false, false, false},
		{"corestring", Env, true, false, false, false},
		{"corestring", Env, false, true, false, false},
		// 165
		{"corestring", Env, true, true, false, false},
		{"corestring", Env, false, false, true, false},
		{"corestring", Env, false, true, true, false},
		{"corestring", Env, true, false, true, false},
		{"corestring", Env, true, true, true, false},
		// 170
		{"x-corestring", Env, false, false, false, false},
		{"corestring", Flag, false, false, false, false},
		{"corestring", Flag, true, false, false, false},
		{"corestring", Flag, false, true, false, false},
		{"corestring", Flag, true, true, false, false},
		// 175
		{"corestring", Flag, false, false, true, false},
		{"corestring", Flag, false, true, true, false},
		{"corestring", Flag, true, false, true, false},
		{"corestring", Flag, true, true, true, false},
		{"x-corestring", Flag, false, false, false, false},
		// 180
		{"cfgbool", Basic, false, false, false, false},
		{"cfgbool", Basic, true, false, false, false},
		{"cfgbool", Basic, false, true, false, false},
		{"cfgbool", Basic, true, true, false, false},
		{"cfgbool", Basic, false, false, true, false},
		// 185
		{"cfgbool", Basic, false, true, true, false},
		{"cfgbool", Basic, true, false, true, false},
		{"cfgbool", Basic, true, true, true, false},
		{"x-cfgbool", Basic, false, false, false, false},
		{"cfgbool", Core, false, false, false, false},
		// 190
		{"cfgbool", Core, true, false, false, false},
		{"cfgbool", Core, false, true, false, false},
		{"cfgbool", Core, true, true, false, false},
		{"cfgbool", Core, false, false, true, false},
		{"cfgbool", Core, false, true, true, false},
		// 195
		{"cfgbool", Core, true, false, true, false},
		{"cfgbool", Core, true, true, true, false},
		{"x-cfgbool", Core, false, false, false, false},
		{"cfgbool", ConfFileVar, false, false, false, true},
		{"cfgbool", ConfFileVar, true, false, false, false},
		// 200
		{"cfgbool", ConfFileVar, false, true, false, false},
		{"cfgbool", ConfFileVar, true, true, false, false},
		{"cfgbool", ConfFileVar, false, false, true, false},
		{"cfgbool", ConfFileVar, false, true, true, false},
		{"cfgbool", ConfFileVar, true, false, true, false},
		// 205
		{"cfgbool", ConfFileVar, true, true, true, false},
		{"x-cfgbool", ConfFileVar, false, false, false, false},
		{"cfgbool", Env, false, false, false, true},
		{"cfgbool", Env, true, false, false, true},
		{"cfgbool", Env, false, true, false, false},
		// 210
		{"cfgbool", Env, true, true, false, false},
		{"cfgbool", Env, false, false, true, false},
		{"cfgbool", Env, false, true, true, false},
		{"cfgbool", Env, true, false, true, false},
		{"cfgbool", Env, true, true, true, false},
		// 215
		{"x-cfgbool", Env, false, false, false, false},
		{"cfgbool", Flag, false, false, false, false},
		{"cfgbool", Flag, true, false, false, false},
		{"cfgbool", Flag, false, true, false, false},
		{"cfgbool", Flag, true, true, false, false},
		// 220
		{"cfgbool11", Flag, false, false, true, false},
		{"cfgbool", Flag, false, true, true, false},
		{"cfgbool", Flag, true, false, true, false},
		{"cfgbool", Flag, true, true, true, false},
		{"x-cfgbool", Flag, false, false, false, false},
		// 225
		{"cfgint", Basic, false, false, false, false},
		{"cfgint", Basic, true, false, false, false},
		{"cfgint", Basic, false, true, false, false},
		{"cfgint", Basic, true, true, false, false},
		{"cfgint", Basic, false, false, true, false},
		// 230
		{"cfgint", Basic, false, true, true, false},
		{"cfgint", Basic, true, false, true, false},
		{"cfgint", Basic, true, true, true, false},
		{"x-cfgint", Basic, false, false, false, false},
		{"cfgint", Core, false, false, false, false},
		// 235
		{"cfgint", Core, true, false, false, false},
		{"cfgint", Core, false, true, false, false},
		{"cfgint", Core, true, true, false, false},
		{"cfgint", Core, false, false, true, false},
		{"cfgint", Core, false, true, true, false},
		// 240
		{"cfgint", Core, true, false, true, false},
		{"cfgbool", Core, true, true, true, false},
		{"x-cfgint", Core, false, false, false, false},
		{"cfgint", ConfFileVar, false, false, false, true},
		{"cfgint", ConfFileVar, true, false, false, false},
		// 245
		{"cfgint", ConfFileVar, false, true, false, false},
		{"cfgint", ConfFileVar, true, true, false, false},
		{"cfgint", ConfFileVar, false, false, true, false},
		{"cfgint", ConfFileVar, false, true, true, false},
		{"cfgint", ConfFileVar, true, false, true, false},
		// 250
		{"cfgint", ConfFileVar, true, true, true, false},
		{"x-cfgint", ConfFileVar, false, false, false, false},
		{"cfgint", Env, false, false, false, true},
		{"cfgint", Env, true, false, false, true},
		{"cfgint", Env, false, true, false, false},
		// 255
		{"cfgint", Env, true, true, false, false},
		{"cfgint", Env, false, false, true, false},
		{"cfgint", Env, false, true, true, false},
		{"cfgint", Env, true, false, true, false},
		{"cfgint", Env, true, true, true, false},
		// 260
		{"x-cfgint", Env, false, false, false, false},
		{"cfgint", Flag, false, false, false, false},
		{"cfgint", Flag, true, false, false, false},
		{"cfgint", Flag, false, true, false, false},
		{"cfgint", Flag, true, true, false, false},
		// 265
		{"cfgint", Flag, false, false, true, false},
		{"cfgint", Flag, false, true, true, false},
		{"cfgint", Flag, true, false, true, false},
		{"cfgint", Flag, true, true, true, false},
		{"x-cfgint", Flag, false, false, false, false},
		// 270
		{"cfgint64", Basic, false, false, false, false},
		{"cfgint64", Basic, true, false, false, false},
		{"cfgint64", Basic, false, true, false, false},
		{"cfgint64", Basic, true, true, false, false},
		{"cfgint64", Basic, false, false, true, false},
		// 275
		{"cfgint64", Basic, false, true, true, false},
		{"cfgint64", Basic, true, false, true, false},
		{"cfgint64", Basic, true, true, true, false},
		{"x-cfgint64", Basic, false, false, false, false},
		{"cfgint64", Core, false, false, false, false},
		// 280
		{"cfgint64", Core, true, false, false, false},
		{"cfgint64", Core, false, true, false, false},
		{"cfgint64", Core, true, true, false, false},
		{"cfgint64", Core, false, false, true, false},
		{"cfgint64", Core, false, true, true, false},
		// 285
		{"cfgint64", Core, true, false, true, false},
		{"cfgint64", Core, true, true, true, false},
		{"x-cfgint64", Core, false, false, false, false},
		{"cfgint64", ConfFileVar, false, false, false, true},
		{"cfgint64", ConfFileVar, true, false, false, false},
		// 290
		{"cfgint64", ConfFileVar, false, true, false, false},
		{"cfgint64", ConfFileVar, true, true, false, false},
		{"cfgint64", ConfFileVar, false, false, true, false},
		{"cfgint64", ConfFileVar, false, true, true, false},
		{"cfgint64", ConfFileVar, true, false, true, false},
		// 295
		{"cfgint64", ConfFileVar, true, true, true, false},
		{"x-cfgint64", ConfFileVar, false, false, false, false},
		{"cfgint64", Env, false, false, false, true},
		{"cfgint64", Env, true, false, false, true},
		{"cfgint64", Env, false, true, false, false},
		// 300
		{"cfgint64", Env, true, true, false, false},
		{"cfgint64", Env, false, false, true, false},
		{"cfgint64", Env, false, true, true, false},
		{"cfgint64", Env, true, false, true, false},
		{"cfgint64", Env, true, true, true, false},
		// 305
		{"x-cfgint64", Env, false, false, false, false},
		{"cfgint64", Flag, false, false, false, false},
		{"cfgint64", Flag, true, false, false, false},
		{"cfgint64", Flag, false, true, false, false},
		{"cfgint64", Flag, true, true, false, false},
		// 310
		{"cfgint64", Flag, false, false, true, false},
		{"cfgint64", Flag, false, true, true, false},
		{"cfgint64", Flag, true, false, true, false},
		{"cfgint64", Flag, true, true, true, false},
		{"x-cfgint64", Flag, false, false, false, false},
		// 315
		{"cfgstring", Basic, false, false, false, false},
		{"cfgstring", Basic, true, false, false, false},
		{"cfgstring", Basic, false, true, false, false},
		{"cfgstring", Basic, true, true, false, false},
		{"cfgstring", Basic, false, false, true, false},
		// 320
		{"cfgstring", Basic, false, true, true, false},
		{"cfgstring", Basic, true, false, true, false},
		{"cfgstring", Basic, true, true, true, false},
		{"x-cfgstring", Basic, false, false, false, false},
		{"cfgstring", Core, false, false, false, false},
		// 325
		{"cfgstring", Core, true, false, false, false},
		{"cfgstring", Core, false, true, false, false},
		{"cfgstring", Core, true, true, false, false},
		{"cfgstring", Core, false, false, true, false},
		{"cfgstring", Core, false, true, true, false},
		// 330
		{"cfgstring", Core, true, false, true, false},
		{"cfgstring", Core, true, true, true, false},
		{"x-cfgstring", Core, false, false, false, false},
		{"cfgstring", ConfFileVar, false, false, false, true},
		{"cfgstring", ConfFileVar, true, false, false, false},
		// 335
		{"cfgstring", ConfFileVar, false, true, false, false},
		{"cfgstring", ConfFileVar, true, true, false, false},
		{"cfgstring", ConfFileVar, false, false, true, false},
		{"cfgstring", ConfFileVar, false, true, true, false},
		{"cfgstring", ConfFileVar, true, false, true, false},
		// 340
		{"cfgstring", ConfFileVar, true, true, true, false},
		{"x-cfgstring", ConfFileVar, false, false, false, false},
		{"cfgstring", Env, false, false, false, true},
		{"cfgstring", Env, true, false, false, true},
		{"cfgstring", Env, false, true, false, false},
		// 345
		{"cfgstring", Env, true, true, false, false},
		{"cfgstring", Env, false, false, true, false},
		{"cfgstring", Env, false, true, true, false},
		{"cfgstring", Env, true, false, true, false},
		{"cfgstring", Env, true, true, true, false},
		// 350
		{"x-cfgstring", Env, false, false, false, false},
		{"cfgstring", Flag, false, false, false, false},
		{"cfgstring", Flag, true, false, false, false},
		{"cfgstring", Flag, false, true, false, false},
		{"cfgstring", Flag, true, true, false, false},
		// 355
		{"cfgstring", Flag, false, false, true, false},
		{"cfgstring", Flag, false, true, true, false},
		{"cfgstring", Flag, true, false, true, false},
		{"cfgstring", Flag, true, true, true, false},
		{"x-cfgstring", Flag, false, false, false, false},
		// 360
		{"envbool", Basic, false, false, false, false},
		{"envbool", Basic, true, false, false, false},
		{"envbool", Basic, false, true, false, false},
		{"envbool", Basic, true, true, false, false},
		{"envbool", Basic, false, false, true, false},
		// 365
		{"envbool", Basic, false, true, true, false},
		{"envbool", Basic, true, false, true, false},
		{"envbool", Basic, true, true, true, false},
		{"x-envbool", Basic, false, false, false, false},
		{"envbool", Core, false, false, false, false},
		// 370
		{"envbool", Core, true, false, false, false},
		{"envbool", Core, false, true, false, false},
		{"envbool", Core, true, true, false, false},
		{"envbool", Core, false, false, true, false},
		{"envbool", Core, false, true, true, false},
		// 375
		{"envbool", Core, true, false, true, false},
		{"envbool", Core, true, true, true, false},
		{"x-envbool", Core, false, false, false, false},
		{"envbool", ConfFileVar, false, false, false, false},
		{"envbool", ConfFileVar, true, false, false, false},
		// 380
		{"envbool", ConfFileVar, false, true, false, false},
		{"envbool", ConfFileVar, true, true, false, false},
		{"envbool", ConfFileVar, false, false, true, false},
		{"envbool", ConfFileVar, false, true, true, false},
		{"envbool", ConfFileVar, true, false, true, false},
		// 385
		{"envbool", ConfFileVar, true, true, true, false},
		{"x-envbool", ConfFileVar, false, false, false, false},
		{"envbool", Env, false, false, false, false},
		{"envbool", Env, true, false, false, false},
		{"envbool", Env, false, true, false, false},
		// 390
		{"envbool", Env, true, true, false, false},
		{"envbool", Env, false, false, true, false},
		{"envbool", Env, false, true, true, false},
		{"envbool", Env, true, false, true, false},
		{"envbool", Env, true, true, true, false},
		// 395
		{"x-envbool", Env, false, false, false, false},
		{"envbool", Flag, false, false, false, false},
		{"envbool", Flag, true, false, false, false},
		{"envbool", Flag, false, true, false, false},
		{"envbool", Flag, true, true, false, false},
		// 400
		{"envbool", Flag, false, false, true, false},
		{"envbool", Flag, false, true, true, false},
		{"envbool", Flag, true, false, true, false},
		{"envbool", Flag, true, true, true, false},
		{"x-envbool", Flag, false, false, false, false},
		// 405
		{"envint", Basic, false, false, false, false},
		{"envint", Basic, true, false, false, false},
		{"envint", Basic, false, true, false, false},
		{"envint", Basic, true, true, false, false},
		{"envint", Basic, false, false, true, false},
		// 410
		{"envint", Basic, false, true, true, false},
		{"envint", Basic, true, false, true, false},
		{"envint", Basic, true, true, true, false},
		{"x-envint", Basic, false, false, false, false},
		{"envbool", Core, false, false, false, false},
		// 415
		{"envbool", Core, true, false, false, false},
		{"envbool", Core, false, true, false, false},
		{"envbool", Core, true, true, false, false},
		{"envbool", Core, false, false, true, false},
		{"envbool", Core, false, true, true, false},
		// 420
		{"envbool", Core, true, false, true, false},
		{"envbool", Core, true, true, true, false},
		{"x-envbool", Core, false, false, false, false},
		{"envbool", ConfFileVar, false, false, false, false},
		{"envbool", ConfFileVar, true, false, false, false},
		// 425
		{"envbool", ConfFileVar, false, true, false, false},
		{"envbool", ConfFileVar, true, true, false, false},
		{"envbool", ConfFileVar, false, false, true, false},
		{"envbool", ConfFileVar, false, true, true, false},
		{"envbool", ConfFileVar, true, false, true, false},
		// 430
		{"envbool", ConfFileVar, true, true, true, false},
		{"x-envbool", ConfFileVar, false, false, false, false},
		{"envbool", Env, false, false, false, false},
		{"envbool", Env, true, false, false, false},
		{"envbool", Env, false, true, false, false},
		// 435
		{"envbool", Env, true, true, false, false},
		{"envbool", Env, false, false, true, false},
		{"envbool", Env, false, true, true, false},
		{"envbool", Env, true, false, true, false},
		{"envbool", Env, true, true, true, false},
		// 440
		{"x-envbool", Env, false, false, false, false},
		{"envbool", Flag, false, false, false, false},
		{"envbool", Flag, true, false, false, false},
		{"envbool", Flag, false, true, false, false},
		{"envbool", Flag, true, true, false, false},
		// 445
		{"envbool", Flag, false, false, true, false},
		{"envbool", Flag, false, true, true, false},
		{"envbool", Flag, true, false, true, false},
		{"envbool", Flag, true, true, true, false},
		{"x-envbool", Flag, false, false, false, false},
		// 450
		{"envint64", Basic, false, false, false, false},
		{"envint64", Basic, true, false, false, false},
		{"envint64", Basic, true, false, false, false},
		{"envint64", Basic, false, true, false, false},
		{"envint64", Basic, true, true, false, false},
		// 455
		{"envint64", Basic, false, false, true, false},
		{"envint64", Basic, false, true, true, false},
		{"envint64", Basic, true, false, true, false},
		{"envint64", Basic, true, true, true, false},
		{"x-envint64", Basic, false, false, false, false},
		// 460
		{"envint64", Core, false, false, false, false},
		{"envint64", Core, true, false, false, false},
		{"envint64", Core, false, true, false, false},
		{"envint64", Core, true, true, false, false},
		{"envint64", Core, false, false, true, false},
		// 465
		{"envint64", Core, false, true, true, false},
		{"envint64", Core, true, false, true, false},
		{"envint64", Core, true, true, true, false},
		{"x-envint64", Core, false, false, false, false},
		{"envint64", ConfFileVar, false, false, false, false},
		// 470
		{"envint64", ConfFileVar, true, false, false, false},
		{"envint64", ConfFileVar, false, true, false, false},
		{"envint64", ConfFileVar, true, true, false, false},
		{"envint64", ConfFileVar, false, false, true, false},
		{"envint64", ConfFileVar, false, true, true, false},
		// 475
		{"envint64", ConfFileVar, true, false, true, false},
		{"envint64", ConfFileVar, true, true, true, false},
		{"x-envint64", ConfFileVar, false, false, false, false},
		{"envint64", Env, false, false, false, false},
		{"envint64", Env, true, false, false, false},
		// 480
		{"envint64", Env, false, true, false, false},
		{"envint64", Env, true, true, false, false},
		{"envint64", Env, false, false, true, false},
		{"envint64", Env, false, true, true, false},
		{"envint64", Env, true, false, true, false},
		// 485
		{"envint64", Env, true, true, true, false},
		{"x-envint64", Env, false, false, false, false},
		{"envint64", Flag, false, false, false, false},
		{"envint64", Flag, true, false, false, false},
		{"envint64", Flag, false, true, false, false},
		// 490
		{"envint64", Flag, true, true, false, false},
		{"envint64", Flag, false, false, true, false},
		{"envint64", Flag, false, true, true, false},
		{"envint64", Flag, true, false, true, false},
		{"envint64", Flag, true, true, true, false},
		// 495
		{"x-envint64", Flag, false, false, false, false},
		{"envstring", Basic, false, false, false, false},
		{"envstring", Basic, true, false, false, false},
		{"envstring", Basic, false, true, false, false},
		{"envstring", Basic, true, true, false, false},
		// 500
		{"envstring", Basic, false, false, true, false},
		{"envstring", Basic, false, true, true, false},
		{"envstring", Basic, true, false, true, false},
		{"envstring", Basic, true, true, true, false},
		{"x-envstring", Basic, false, false, false, false},
		// 505
		{"envstring", Core, false, false, false, false},
		{"envstring", Core, true, false, false, false},
		{"envstring", Core, false, true, false, false},
		{"envstring", Core, true, true, false, false},
		{"envstring", Core, false, false, true, false},
		// 510
		{"envstring", Core, false, true, true, false},
		{"envstring", Core, true, false, true, false},
		{"envstring", Core, true, true, true, false},
		{"x-envstring", Core, false, false, false, false},
		{"envstring", ConfFileVar, false, false, false, false},
		// 515
		{"envstring", ConfFileVar, true, false, false, false},
		{"envstring", ConfFileVar, false, true, false, false},
		{"envstring", ConfFileVar, true, true, false, false},
		{"envstring", ConfFileVar, false, false, true, false},
		{"envstring", ConfFileVar, false, true, true, false},
		// 520
		{"envstring", ConfFileVar, true, false, true, false},
		{"envstring", ConfFileVar, true, true, true, false},
		{"x-envstring", ConfFileVar, false, false, false, false},
		{"envstring", Env, false, false, false, false},
		{"envstring", Env, true, false, false, false},
		// 525
		{"envstring", Env, false, true, false, false},
		{"envstring", Env, true, true, false, false},
		{"envstring", Env, false, false, true, false},
		{"envstring", Env, false, true, true, false},
		{"envstring", Env, true, false, true, false},
		// 530
		{"envstring", Env, true, true, true, false},
		{"x-envstring", Env, false, false, false, false},
		{"envstring", Flag, false, false, false, false},
		{"envstring", Flag, true, false, false, false},
		{"envstring", Flag, false, true, false, false},
		// 535
		{"envstring", Flag, true, true, false, false},
		{"envstring", Flag, false, false, true, false},
		{"envstring", Flag, false, true, true, false},
		{"envstring", Flag, true, false, true, false},
		{"envstring", Flag, true, true, true, false},
		// 540
		{"x-envstring", Flag, false, false, false, false},
		{"flagbool", Basic, false, false, false, false},
		{"flagbool", Basic, true, false, false, false},
		{"flagbool", Basic, false, true, false, false},
		{"flagbool", Basic, true, true, false, false},
		// 545
		{"flagbool", Basic, false, false, true, false},
		{"flagbool", Basic, false, true, true, false},
		{"flagbool", Basic, true, false, true, false},
		{"flagbool", Basic, true, true, true, false},
		{"x-flagbool", Basic, false, false, false, false},
		// 550
		{"flagbool", Core, false, false, false, false},
		{"flagbool", Core, true, false, false, false},
		{"flagbool", Core, false, true, false, false},
		{"flagbool", Core, true, true, false, false},
		{"flagbool", Core, false, false, true, false},
		// 555
		{"flagbool", Core, false, true, true, false},
		{"flagbool", Core, true, false, true, false},
		{"flagbool", Core, true, true, true, false},
		{"x-flagbool", Core, false, false, false, false},
		{"flagbool", ConfFileVar, false, false, false, true},
		// 560
		{"flagbool", ConfFileVar, true, false, false, false},
		{"flagbool", ConfFileVar, false, true, false, false},
		{"flagbool", ConfFileVar, true, true, false, false},
		{"flagbool", ConfFileVar, false, false, true, false},
		{"flagbool", ConfFileVar, false, true, true, false},
		// 565
		{"flagbool", ConfFileVar, true, false, true, false},
		{"flagbool", ConfFileVar, true, true, true, false},
		{"x-flagbool", ConfFileVar, false, false, false, false},
		{"flagbool", Env, false, false, false, true},
		{"flagbool", Env, true, false, false, true},
		// 570
		{"flagbool", Env, false, true, false, false},
		{"flagbool", Env, true, true, false, false},
		{"flagbool", Env, false, false, true, false},
		{"flagbool", Env, false, true, true, false},
		{"flagbool", Env, true, false, true, false},
		// 575
		{"flagbool", Env, true, true, true, false},
		{"x-flagbool", Env, false, false, false, false},
		{"flagbool", Flag, false, false, false, true},
		{"flagbool", Flag, true, false, false, true},
		{"flagbool", Flag, false, true, false, true},
		// 580
		{"flagbool", Flag, true, true, false, true},
		{"flagbool", Flag, false, false, true, false},
		{"flagbool", Flag, false, true, true, false},
		{"flagbool", Flag, true, false, true, false},
		{"flagbool", Flag, true, true, true, false},
		// 585
		{"x-flagbool", Flag, false, false, false, false},
		{"flagint", Basic, false, false, false, false},
		{"flagint", Basic, true, false, false, false},
		{"flagint", Basic, true, false, false, false},
		{"flagint", Basic, false, true, false, false},
		// 590
		{"flagint", Basic, true, true, false, false},
		{"flagint", Basic, false, false, true, false},
		{"flagint", Basic, false, true, true, false},
		{"flagint", Basic, true, false, true, false},
		{"flagint", Basic, true, true, true, false},
		// 595
		{"x-flagint", Basic, false, false, false, false},
		{"flagbool", Core, false, false, false, false},
		{"flagbool", Core, true, false, false, false},
		{"flagbool", Core, false, true, false, false},
		{"flagbool", Core, true, true, false, false},
		// 600
		{"flagbool", Core, false, false, true, false},
		{"flagbool", Core, false, true, true, false},
		{"flagbool", Core, true, false, true, false},
		{"flagbool", Core, true, true, true, false},
		{"x-flagbool", Core, false, false, false, false},
		// 605
		{"flagbool", ConfFileVar, false, false, false, true},
		{"flagbool", ConfFileVar, true, false, false, false},
		{"flagbool", ConfFileVar, false, true, false, false},
		{"flagbool", ConfFileVar, true, true, false, false},
		{"flagbool", ConfFileVar, false, false, true, false},
		// 610
		{"flagbool", ConfFileVar, false, true, true, false},
		{"flagbool", ConfFileVar, true, false, true, false},
		{"flagbool", ConfFileVar, true, true, true, false},
		{"x-flagbool", ConfFileVar, false, false, false, false},
		{"flagbool", Env, false, false, false, true},
		// 615
		{"flagbool", Env, true, false, false, true},
		{"flagbool", Env, false, true, false, false},
		{"flagbool", Env, true, true, false, false},
		{"flagbool", Env, false, false, true, false},
		{"flagbool", Env, false, true, true, false},
		// 620
		{"flagbool", Env, true, false, true, false},
		{"flagbool", Env, true, true, true, false},
		{"x-flagbool", Env, false, false, false, false},
		{"flagbool", Flag, false, false, false, true},
		{"flagbool", Flag, true, false, false, true},
		// 625
		{"flagbool", Flag, false, true, false, true},
		{"flagbool", Flag, true, true, false, true},
		{"flagbool", Flag, false, false, true, false},
		{"flagbool", Flag, false, true, true, false},
		{"flagbool", Flag, true, false, true, false},
		// 630
		{"flagbool", Flag, true, true, true, false},
		{"x-flagbool", Flag, false, false, false, false},
		{"flagint64", Basic, false, false, false, false},
		{"flagint64", Basic, true, false, false, false},
		{"flagint64", Basic, true, false, false, false},
		// 635
		{"flagint64", Basic, false, true, false, false},
		{"flagint64", Basic, true, true, false, false},
		{"flagint64", Basic, false, false, true, false},
		{"flagint64", Basic, false, true, true, false},
		{"flagint64", Basic, true, false, true, false},
		// 640
		{"flagint64", Basic, true, true, true, false},
		{"x-flagint64", Basic, false, false, false, false},
		{"flagint64", Core, false, false, false, false},
		{"flagint64", Core, true, false, false, false},
		{"flagint64", Core, false, true, false, false},
		// 645
		{"flagint64", Core, true, true, false, false},
		{"flagint64", Core, false, false, true, false},
		{"flagint64", Core, false, true, true, false},
		{"flagint64", Core, true, false, true, false},
		{"flagint64", Core, true, true, true, false},
		// 650
		{"x-flagint64", Core, false, false, false, false},
		{"flagint64", ConfFileVar, false, false, false, true},
		{"flagint64", ConfFileVar, true, false, false, false},
		{"flagint64", ConfFileVar, false, true, false, false},
		{"flagint64", ConfFileVar, true, true, false, false},
		// 655
		{"flagint64", ConfFileVar, false, false, true, false},
		{"flagint64", ConfFileVar, false, true, true, false},
		{"flagint64", ConfFileVar, true, false, true, false},
		{"flagint64", ConfFileVar, true, true, true, false},
		{"x-flagint64", ConfFileVar, false, false, false, false},
		// 660
		{"flagint64", Env, false, false, false, true},
		{"flagint64", Env, true, false, false, true},
		{"flagint64", Env, false, true, false, false},
		{"flagint64", Env, true, true, false, false},
		{"flagint64", Env, false, false, true, false},
		// 665
		{"flagint64", Env, false, true, true, false},
		{"flagint64", Env, true, false, true, false},
		{"flagint64", Env, true, true, true, false},
		{"x-flagint64", Env, false, false, false, false},
		{"flagint64", Flag, false, false, false, true},
		// 670
		{"flagint64", Flag, true, false, false, true},
		{"flagint64", Flag, false, true, false, true},
		{"flagint64", Flag, true, true, false, true},
		{"flagint64", Flag, false, false, true, false},
		{"flagint64", Flag, false, true, true, false},
		// 675
		{"flagint64", Flag, true, false, true, false},
		{"flagint64", Flag, true, true, true, false},
		{"x-flagint64", Flag, false, false, false, false},
		{"flagstring", Basic, false, false, false, false},
		{"flagstring", Basic, true, false, false, false},
		// 680
		{"flagstring", Basic, true, false, false, false},
		{"flagstring", Basic, false, true, false, false},
		{"flagstring", Basic, true, true, false, false},
		{"flagstring", Basic, false, false, true, false},
		{"flagstring", Basic, false, true, true, false},
		// 685
		{"flagstring", Basic, true, false, true, false},
		{"flagstring", Basic, true, true, true, false},
		{"x-flagstring", Basic, false, false, false, false},
		{"flagstring", Core, false, false, false, false},
		{"flagstring", Core, true, false, false, false},
		// 690
		{"flagstring", Core, false, true, false, false},
		{"flagstring", Core, true, true, false, false},
		{"flagstring", Core, false, false, true, false},
		{"flagstring", Core, false, true, true, false},
		{"flagstring", Core, true, false, true, false},
		// 695
		{"flagstring", Core, true, true, true, false},
		{"x-flagstring", Core, false, false, false, false},
		{"flagstring", ConfFileVar, false, false, false, true},
		{"flagstring", ConfFileVar, true, false, false, false},
		{"flagstring", ConfFileVar, false, true, false, false},
		// 700
		{"flagstring", ConfFileVar, true, true, false, false},
		{"flagstring", ConfFileVar, false, false, true, false},
		{"flagstring", ConfFileVar, false, true, true, false},
		{"flagstring", ConfFileVar, true, false, true, false},
		{"flagstring", ConfFileVar, true, true, true, false},
		// 705
		{"x-flagstring", ConfFileVar, false, false, false, false},
		{"flagstring", Env, false, false, false, true},
		{"flagstring", Env, true, false, false, true},
		{"flagstring", Env, false, true, false, false},
		{"flagstring", Env, true, true, false, false},
		// 710
		{"flagstring", Env, false, false, true, false},
		{"flagstring", Env, false, true, true, false},
		{"flagstring", Env, true, false, true, false},
		{"flagstring", Env, true, true, true, false},
		{"x-flagstring", Env, false, false, false, false},
		// 715
		{"flagstring", Flag, false, false, false, true},
		{"flagstring", Flag, true, false, false, true},
		{"flagstring", Flag, false, true, false, true},
		{"flagstring", Flag, true, true, false, true},
		{"flagstring", Flag, false, false, true, false},
		// 720
		{"flagstring", Flag, false, true, true, false},
		{"flagstring", Flag, true, false, true, false},
		{"flagstring", Flag, true, true, true, false},
		{"x-flagstring", Flag, false, false, false, false},
		{"bool", Basic, false, false, false, true},
		// 725
		{"bool", Basic, true, false, false, true},
		{"bool", Basic, false, true, false, true},
		{"bool", Basic, true, true, false, true},
		{"bool", Basic, false, false, true, true},
		{"bool", Basic, false, true, true, true},
		// 730
		{"bool", Basic, true, false, true, true},
		{"bool", Basic, true, true, true, true},
		{"x-bool", Basic, false, false, false, false},
		{"bool", Core, false, false, false, false},
		{"bool", Core, true, false, false, false},
		// 735
		{"bool", Core, false, true, false, false},
		{"bool", Core, true, true, false, false},
		{"bool", Core, false, false, true, false},
		{"bool", Core, false, true, true, false},
		{"bool", Core, true, false, true, false},
		// 740
		{"bool", Core, true, true, true, false},
		{"x-bool", Core, false, false, false, false},
		{"bool", ConfFileVar, false, false, false, false},
		{"bool", ConfFileVar, true, false, false, false},
		{"bool", ConfFileVar, false, true, false, false},
		// 745
		{"bool", ConfFileVar, true, true, false, false},
		{"bool", ConfFileVar, false, false, true, false},
		{"bool", ConfFileVar, false, true, true, false},
		{"bool", ConfFileVar, true, false, true, false},
		{"bool", ConfFileVar, true, true, true, false},
		// 750
		{"x-bool", ConfFileVar, false, false, false, false},
		{"bool", Env, false, false, false, false},
		{"bool", Env, true, false, false, false},
		{"bool", Env, false, true, false, false},
		{"bool", Env, true, true, false, false},
		// 755
		{"bool", Env, false, false, true, false},
		{"bool", Env, false, true, true, false},
		{"bool", Env, true, false, true, false},
		{"bool", Env, true, true, true, false},
		{"x-bool", Env, false, false, false, false},
		// 760
		{"bool", Flag, false, false, false, false},
		{"bool", Flag, true, false, false, false},
		{"bool", Flag, false, true, false, false},
		{"bool", Flag, true, true, false, false},
		{"bool", Flag, false, false, true, false},
		// 765
		{"bool", Flag, false, true, true, false},
		{"bool", Flag, true, false, true, false},
		{"bool", Flag, true, true, true, false},
		{"x-bool", Flag, false, false, false, false},
		{"int", Basic, false, false, false, true},
		// 770
		{"int", Basic, true, false, false, true},
		{"int", Basic, false, true, false, true},
		{"int", Basic, true, true, false, true},
		{"int", Basic, false, false, true, true},
		{"int", Basic, false, true, true, true},
		// 775
		{"int", Basic, true, false, true, true},
		{"int", Basic, true, true, true, true},
		{"x-int", Basic, false, false, false, false},
		{"int", Core, false, false, false, false},
		{"int", Core, true, false, false, false},
		// 780
		{"int", Core, false, true, false, false},
		{"int", Core, true, true, false, false},
		{"int", Core, false, false, true, false},
		{"int", Core, false, true, true, false},
		{"int", Core, true, false, true, false},
		// 785
		{"int", Core, true, true, true, false},
		{"x-int", Core, false, false, false, false},
		{"int", ConfFileVar, false, false, false, false},
		{"int", ConfFileVar, true, false, false, false},
		{"int", ConfFileVar, false, true, false, false},
		// 790
		{"int", ConfFileVar, true, true, false, false},
		{"int", ConfFileVar, false, false, true, false},
		{"int", ConfFileVar, false, true, true, false},
		{"int", ConfFileVar, true, false, true, false},
		// 795
		{"int", ConfFileVar, true, true, true, false},
		{"x-int", ConfFileVar, false, false, false, false},
		{"int", Env, false, false, false, false},
		{"int", Env, true, false, false, false},
		{"int", Env, false, true, false, false},
		// 800
		{"int", Env, true, true, false, false},
		{"int", Env, false, false, true, false},
		{"int", Env, false, true, true, false},
		{"int", Env, true, false, true, false},
		{"int", Env, true, true, true, false},
		{"x-int", Env, false, false, false, false},
		// 805
		{"int", Flag, false, false, false, false},
		{"int", Flag, true, false, false, false},
		{"int", Flag, false, true, false, false},
		{"int", Flag, true, true, false, false},
		{"int", Flag, false, false, true, false},
		// 810
		{"int", Flag, false, true, true, false},
		{"int", Flag, true, false, true, false},
		{"int", Flag, true, true, true, false},
		{"x-int", Flag, false, false, false, false},
		{"int64", Basic, false, false, false, true},
		// 815
		{"int64", Basic, true, false, false, true},
		{"int64", Basic, false, true, false, true},
		{"int64", Basic, true, true, false, true},
		{"int64", Basic, false, false, true, true},
		{"int64", Basic, false, true, true, true},
		// 820
		{"int64", Basic, true, false, true, true},
		{"int64", Basic, true, true, true, true},
		{"x-int64", Basic, false, false, false, false},
		{"int64", Core, false, false, false, false},
		{"int64", Core, true, false, false, false},
		// 825
		{"int64", Core, false, true, false, false},
		{"int64", Core, true, true, false, false},
		{"int64", Core, false, false, true, false},
		{"int64", Core, false, true, true, false},
		{"int64", Core, true, false, true, false},
		// 830
		{"int64", Core, true, true, true, false},
		{"x-int64", Core, false, false, false, false},
		{"int64", ConfFileVar, false, false, false, false},
		{"int64", ConfFileVar, true, false, false, false},
		{"int64", ConfFileVar, false, true, false, false},
		// 835
		{"int64", ConfFileVar, true, true, false, false},
		{"int64", ConfFileVar, false, false, true, false},
		{"int64", ConfFileVar, false, true, true, false},
		{"int64", ConfFileVar, true, false, true, false},
		{"int64", ConfFileVar, true, true, true, false},
		// 840
		{"x-int64", ConfFileVar, false, false, false, false},
		{"int64", Env, false, false, false, false},
		{"int64", Env, true, false, false, false},
		{"int64", Env, false, true, false, false},
		{"int64", Env, true, true, false, false},
		// 845
		{"int64", Env, false, false, true, false},
		{"int64", Env, false, true, true, false},
		{"int64", Env, true, false, true, false},
		{"int64", Env, true, true, true, false},
		{"x-int64", Env, false, false, false, false},
		// 850
		{"int64", Flag, false, false, false, false},
		{"int64", Flag, true, false, false, false},
		{"int64", Flag, false, true, false, false},
		{"int64", Flag, true, true, false, false},
		{"int64", Flag, false, false, true, false},
		// 855
		{"int64", Flag, false, true, true, false},
		{"int64", Flag, true, false, true, false},
		{"int64", Flag, true, true, true, false},
		{"x-int64", Flag, false, false, false, false},
		{"string", Basic, false, false, false, true},
		// 860
		{"string", Basic, true, false, false, true},
		{"string", Basic, false, true, false, true},
		{"string", Basic, true, true, false, true},
		{"string", Basic, false, false, true, true},
		{"string", Basic, false, true, true, true},
		// 865
		{"string", Basic, true, false, true, true},
		{"string", Basic, true, true, true, true},
		{"x-string", Basic, false, false, false, false},
		{"string", Core, false, false, false, false},
		{"string", Core, true, false, false, false},
		// 870
		{"string", Core, false, true, false, false},
		{"string", Core, true, true, false, false},
		{"string", Core, false, false, true, false},
		{"string", Core, false, true, true, false},
		{"string", Core, true, false, true, false},
		// 875
		{"string", Core, true, true, true, false},
		{"x-string", Core, false, false, false, false},
		{"string", ConfFileVar, false, false, false, false},
		{"string", ConfFileVar, true, false, false, false},
		{"string", ConfFileVar, false, true, false, false},
		// 880
		{"string", ConfFileVar, true, true, false, false},
		{"string", ConfFileVar, false, false, true, false},
		{"string", ConfFileVar, false, true, true, false},
		{"string", ConfFileVar, true, false, true, false},
		{"string", ConfFileVar, true, true, true, false},
		// 885
		{"x-string", ConfFileVar, false, false, false, false},
		{"string", Env, false, false, false, false},
		{"string", Env, true, false, false, false},
		{"string", Env, false, true, false, false},
		{"string", Env, true, true, false, false},
		// 890
		{"string", Env, false, false, true, false},
		{"string", Env, false, true, true, false},
		{"string", Env, true, false, true, false},
		{"string", Env, true, true, true, false},
		{"x-string", Env, false, false, false, false},
		// 895
		{"string", Flag, false, false, false, false},
		{"string", Flag, true, false, false, false},
		{"string", Flag, false, true, false, false},
		{"string", Flag, true, true, false, false},
		{"string", Flag, false, false, true, false},
		// 900
		{"string", Flag, false, true, true, false},
		{"string", Flag, true, false, true, false},
		{"string", Flag, true, true, true, false},
		{"x-string", Flag, false, false, false, false},
	}
	appCfg := newTestSettings()
	for i, test := range tests {
		appCfg.flagsParsed = test.flagsParsed
		appCfg.confFileVarsSet = test.confFileVarsSet
		appCfg.envSet = test.envVarsSet
		appCfg.flagsParsed = test.flagsParsed
		b := appCfg.canUpdate(test.typ, test.name)
		if b != test.expected {
			t.Errorf("%d: %s:%s: expected %v got %v", i, test.name, test.typ, test.expected, b)
		}
	}
}
