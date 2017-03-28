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
		err             error
	}{
		// 0
		{"corebool", Basic, false, false, false, false, CoreUpdateErr{k: "corebool"}},
		{"corebool", Basic, true, false, false, false, CoreUpdateErr{k: "corebool"}},
		{"corebool", Basic, false, true, false, false, CoreUpdateErr{k: "corebool"}},
		{"corebool", Basic, true, true, false, false, CoreUpdateErr{k: "corebool"}},
		{"corebool", Basic, false, false, true, false, CoreUpdateErr{k: "corebool"}},
		// 5
		{"corebool", Basic, false, true, true, false, CoreUpdateErr{k: "corebool"}},
		{"corebool", Basic, true, false, true, false, CoreUpdateErr{k: "corebool"}},
		{"corebool", Basic, true, true, true, false, CoreUpdateErr{k: "corebool"}},
		{"x-corebool", Basic, false, false, false, false, SettingNotFoundErr{name: "x-corebool"}},
		{"corebool", Core, false, false, false, false, CoreUpdateErr{k: "corebool"}},
		// 10
		{"corebool", Core, true, false, false, false, CoreUpdateErr{k: "corebool"}},
		{"corebool", Core, false, true, false, false, CoreUpdateErr{k: "corebool"}},
		{"corebool", Core, true, true, false, false, CoreUpdateErr{k: "corebool"}},
		{"corebool", Core, false, false, true, false, CoreUpdateErr{k: "corebool"}},
		{"corebool", Core, false, true, true, false, CoreUpdateErr{k: "corebool"}},
		// 15
		{"corebool", Core, true, false, true, false, CoreUpdateErr{k: "corebool"}},
		{"corebool", Core, true, true, true, false, CoreUpdateErr{k: "corebool"}},
		{"x-corebool", Core, false, false, false, false, SettingNotFoundErr{name: "x-corebool"}},
		{"corebool", ConfFileVar, false, false, false, false, CoreUpdateErr{k: "corebool"}},
		{"corebool", ConfFileVar, true, false, false, false, CoreUpdateErr{k: "corebool"}},
		// 20
		{"corebool", ConfFileVar, false, true, false, false, CoreUpdateErr{k: "corebool"}},
		{"corebool", ConfFileVar, true, true, false, false, CoreUpdateErr{k: "corebool"}},
		{"corebool", ConfFileVar, false, false, true, false, CoreUpdateErr{k: "corebool"}},
		{"corebool", ConfFileVar, false, true, true, false, CoreUpdateErr{k: "corebool"}},
		{"corebool", ConfFileVar, true, false, true, false, CoreUpdateErr{k: "corebool"}},
		// 25
		{"corebool", ConfFileVar, true, true, true, false, CoreUpdateErr{k: "corebool"}},
		{"x-corebool", ConfFileVar, false, false, false, false, SettingNotFoundErr{name: "x-corebool"}},
		{"corebool", Env, false, false, false, false, CoreUpdateErr{k: "corebool"}},
		{"corebool", Env, true, false, false, false, CoreUpdateErr{k: "corebool"}},
		{"corebool", Env, false, true, false, false, CoreUpdateErr{k: "corebool"}},
		// 30
		{"corebool", Env, true, true, false, false, CoreUpdateErr{k: "corebool"}},
		{"corebool", Env, false, false, true, false, CoreUpdateErr{k: "corebool"}},
		{"corebool", Env, false, true, true, false, CoreUpdateErr{k: "corebool"}},
		{"corebool", Env, true, false, true, false, CoreUpdateErr{k: "corebool"}},
		{"corebool", Env, true, true, true, false, CoreUpdateErr{k: "corebool"}},
		// 35
		{"x-corebool", Env, false, false, false, false, SettingNotFoundErr{name: "x-corebool"}},
		{"corebool", Flag, false, false, false, false, CoreUpdateErr{k: "corebool"}},
		{"corebool", Flag, true, false, false, false, CoreUpdateErr{k: "corebool"}},
		{"corebool", Flag, false, true, false, false, CoreUpdateErr{k: "corebool"}},
		{"corebool", Flag, true, true, false, false, CoreUpdateErr{k: "corebool"}},
		// 40
		{"corebool", Flag, false, false, true, false, CoreUpdateErr{k: "corebool"}},
		{"corebool", Flag, false, true, true, false, CoreUpdateErr{k: "corebool"}},
		{"corebool", Flag, true, false, true, false, CoreUpdateErr{k: "corebool"}},
		{"corebool", Flag, true, true, true, false, CoreUpdateErr{k: "corebool"}},
		{"x-corebool", Flag, false, false, false, false, SettingNotFoundErr{name: "x-corebool"}},
		// 45
		{"coreint", Basic, false, false, false, false, CoreUpdateErr{k: "coreint"}},
		{"coreint", Basic, true, false, false, false, CoreUpdateErr{k: "coreint"}},
		{"coreint", Basic, false, true, false, false, CoreUpdateErr{k: "coreint"}},
		{"coreint", Basic, true, true, false, false, CoreUpdateErr{k: "coreint"}},
		{"coreint", Basic, false, false, true, false, CoreUpdateErr{k: "coreint"}},
		// 50
		{"coreint", Basic, false, true, true, false, CoreUpdateErr{k: "coreint"}},
		{"coreint", Basic, true, false, true, false, CoreUpdateErr{k: "coreint"}},
		{"coreint", Basic, true, true, true, false, CoreUpdateErr{k: "coreint"}},
		{"x-coreint", Basic, false, false, false, false, SettingNotFoundErr{name: "x-coreint"}},
		{"coreint", Core, false, false, false, false, CoreUpdateErr{k: "coreint"}},
		// 55
		{"coreint", Core, true, false, false, false, CoreUpdateErr{k: "coreint"}},
		{"coreint", Core, false, true, false, false, CoreUpdateErr{k: "coreint"}},
		{"coreint", Core, true, true, false, false, CoreUpdateErr{k: "coreint"}},
		{"coreint", Core, false, false, true, false, CoreUpdateErr{k: "coreint"}},
		{"coreint", Core, false, true, true, false, CoreUpdateErr{k: "coreint"}},
		// 60
		{"coreint", Core, true, false, true, false, CoreUpdateErr{k: "coreint"}},
		{"coreint", Core, true, true, true, false, CoreUpdateErr{k: "coreint"}},
		{"x-coreint", Core, false, false, false, false, SettingNotFoundErr{name: "x-coreint"}},
		{"coreint", ConfFileVar, false, false, false, false, CoreUpdateErr{k: "coreint"}},
		{"coreint", ConfFileVar, true, false, false, false, CoreUpdateErr{k: "coreint"}},
		// 65
		{"coreint", ConfFileVar, false, true, false, false, CoreUpdateErr{k: "coreint"}},
		{"coreint", ConfFileVar, true, true, false, false, CoreUpdateErr{k: "coreint"}},
		{"coreint", ConfFileVar, false, false, true, false, CoreUpdateErr{k: "coreint"}},
		{"coreint", ConfFileVar, false, true, true, false, CoreUpdateErr{k: "coreint"}},
		{"coreint", ConfFileVar, true, false, true, false, CoreUpdateErr{k: "coreint"}},
		// 70
		{"coreint", ConfFileVar, true, true, true, false, CoreUpdateErr{k: "coreint"}},
		{"x-coreint", ConfFileVar, false, false, false, false, SettingNotFoundErr{name: "x-coreint"}},
		{"coreint", Env, false, false, false, false, CoreUpdateErr{k: "coreint"}},
		{"coreint", Env, true, false, false, false, CoreUpdateErr{k: "coreint"}},
		{"coreint", Env, false, true, false, false, CoreUpdateErr{k: "coreint"}},
		// 75
		{"coreint", Env, true, true, false, false, CoreUpdateErr{k: "coreint"}},
		{"coreint", Env, false, false, true, false, CoreUpdateErr{k: "coreint"}},
		{"coreint", Env, false, true, true, false, CoreUpdateErr{k: "coreint"}},
		{"coreint", Env, true, false, true, false, CoreUpdateErr{k: "coreint"}},
		{"coreint", Env, true, true, true, false, CoreUpdateErr{k: "coreint"}},
		// 80
		{"x-coreint", Env, false, false, false, false, SettingNotFoundErr{name: "x-coreint"}},
		{"coreint", Flag, false, false, false, false, CoreUpdateErr{k: "coreint"}},
		{"coreint", Flag, true, false, false, false, CoreUpdateErr{k: "coreint"}},
		{"coreint", Flag, false, true, false, false, CoreUpdateErr{k: "coreint"}},
		{"coreint", Flag, true, true, false, false, CoreUpdateErr{k: "coreint"}},
		// 85
		{"coreint", Flag, false, false, true, false, CoreUpdateErr{k: "coreint"}},
		{"coreint", Flag, false, true, true, false, CoreUpdateErr{k: "coreint"}},
		{"coreint", Flag, true, false, true, false, CoreUpdateErr{k: "coreint"}},
		{"coreint", Flag, true, true, true, false, CoreUpdateErr{k: "coreint"}},
		{"x-coreint", Flag, false, false, false, false, SettingNotFoundErr{name: "x-coreint"}},
		// 90
		{"coreint64", Basic, false, false, false, false, CoreUpdateErr{k: "coreint64"}},
		{"coreint64", Basic, true, false, false, false, CoreUpdateErr{k: "coreint64"}},
		{"coreint64", Basic, false, true, false, false, CoreUpdateErr{k: "coreint64"}},
		{"coreint64", Basic, true, true, false, false, CoreUpdateErr{k: "coreint64"}},
		{"coreint64", Basic, false, false, true, false, CoreUpdateErr{k: "coreint64"}},
		// 95
		{"coreint64", Basic, false, true, true, false, CoreUpdateErr{k: "coreint64"}},
		{"coreint64", Basic, true, false, true, false, CoreUpdateErr{k: "coreint64"}},
		{"coreint64", Basic, true, true, true, false, CoreUpdateErr{k: "coreint64"}},
		{"x-coreint64", Basic, false, false, false, false, SettingNotFoundErr{name: "x-coreint64"}},
		{"coreint64", Core, false, false, false, false, CoreUpdateErr{k: "coreint64"}},
		// 100
		{"coreint64", Core, true, false, false, false, CoreUpdateErr{k: "coreint64"}},
		{"coreint64", Core, false, true, false, false, CoreUpdateErr{k: "coreint64"}},
		{"coreint64", Core, true, true, false, false, CoreUpdateErr{k: "coreint64"}},
		{"coreint64", Core, false, false, true, false, CoreUpdateErr{k: "coreint64"}},
		{"coreint64", Core, false, true, true, false, CoreUpdateErr{k: "coreint64"}},
		// 105
		{"coreint64", Core, true, false, true, false, CoreUpdateErr{k: "coreint64"}},
		{"coreint64", Core, true, true, true, false, CoreUpdateErr{k: "coreint64"}},
		{"x-coreint64", Core, false, false, false, false, SettingNotFoundErr{name: "x-coreint64"}},
		{"coreint64", ConfFileVar, false, false, false, false, CoreUpdateErr{k: "coreint64"}},
		{"coreint64", ConfFileVar, true, false, false, false, CoreUpdateErr{k: "coreint64"}},
		// 110
		{"coreint64", ConfFileVar, false, true, false, false, CoreUpdateErr{k: "coreint64"}},
		{"coreint64", ConfFileVar, true, true, false, false, CoreUpdateErr{k: "coreint64"}},
		{"coreint64", ConfFileVar, false, false, true, false, CoreUpdateErr{k: "coreint64"}},
		{"coreint64", ConfFileVar, false, true, true, false, CoreUpdateErr{k: "coreint64"}},
		{"coreint64", ConfFileVar, true, false, true, false, CoreUpdateErr{k: "coreint64"}},
		// 115
		{"coreint64", ConfFileVar, true, true, true, false, CoreUpdateErr{k: "coreint64"}},
		{"x-coreint64", ConfFileVar, false, false, false, false, SettingNotFoundErr{name: "x-coreint64"}},
		{"coreint64", Env, false, false, false, false, CoreUpdateErr{k: "coreint64"}},
		{"coreint64", Env, true, false, false, false, CoreUpdateErr{k: "coreint64"}},
		{"coreint64", Env, false, true, false, false, CoreUpdateErr{k: "coreint64"}},
		// 120
		{"coreint64", Env, true, true, false, false, CoreUpdateErr{k: "coreint64"}},
		{"coreint64", Env, false, false, true, false, CoreUpdateErr{k: "coreint64"}},
		{"coreint64", Env, false, true, true, false, CoreUpdateErr{k: "coreint64"}},
		{"coreint64", Env, true, false, true, false, CoreUpdateErr{k: "coreint64"}},
		{"coreint64", Env, true, true, true, false, CoreUpdateErr{k: "coreint64"}},
		// 125
		{"x-coreint64", Env, false, false, false, false, SettingNotFoundErr{name: "x-coreint64"}},
		{"coreint64", Flag, false, false, false, false, CoreUpdateErr{k: "coreint64"}},
		{"coreint64", Flag, true, false, false, false, CoreUpdateErr{k: "coreint64"}},
		{"coreint64", Flag, false, true, false, false, CoreUpdateErr{k: "coreint64"}},
		{"coreint64", Flag, true, true, false, false, CoreUpdateErr{k: "coreint64"}},
		// 130
		{"coreint64", Flag, false, false, true, false, CoreUpdateErr{k: "coreint64"}},
		{"coreint64", Flag, false, true, true, false, CoreUpdateErr{k: "coreint64"}},
		{"coreint64", Flag, true, false, true, false, CoreUpdateErr{k: "coreint64"}},
		{"coreint64", Flag, true, true, true, false, CoreUpdateErr{k: "coreint64"}},
		{"x-coreint64", Flag, false, false, false, false, SettingNotFoundErr{name: "x-coreint64"}},
		// 135
		{"corestring", Basic, false, false, false, false, CoreUpdateErr{k: "corestring"}},
		{"corestring", Basic, true, false, false, false, CoreUpdateErr{k: "corestring"}},
		{"corestring", Basic, false, true, false, false, CoreUpdateErr{k: "corestring"}},
		{"corestring", Basic, true, true, false, false, CoreUpdateErr{k: "corestring"}},
		{"corestring", Basic, false, false, true, false, CoreUpdateErr{k: "corestring"}},
		// 140
		{"corestring", Basic, false, true, true, false, CoreUpdateErr{k: "corestring"}},
		{"corestring", Basic, true, false, true, false, CoreUpdateErr{k: "corestring"}},
		{"corestring", Basic, true, true, true, false, CoreUpdateErr{k: "corestring"}},
		{"x-corestring", Basic, false, false, false, false, SettingNotFoundErr{name: "x-corestring"}},
		{"corestring", Core, false, false, false, false, CoreUpdateErr{k: "corestring"}},
		// 145
		{"corestring", Core, true, false, false, false, CoreUpdateErr{k: "corestring"}},
		{"corestring", Core, false, true, false, false, CoreUpdateErr{k: "corestring"}},
		{"corestring", Core, true, true, false, false, CoreUpdateErr{k: "corestring"}},
		{"corestring", Core, false, false, true, false, CoreUpdateErr{k: "corestring"}},
		{"corestring", Core, false, true, true, false, CoreUpdateErr{k: "corestring"}},
		// 150
		{"corestring", Core, true, false, true, false, CoreUpdateErr{k: "corestring"}},
		{"corestring", Core, true, true, true, false, CoreUpdateErr{k: "corestring"}},
		{"x-corestring", Core, false, false, false, false, SettingNotFoundErr{name: "x-corestring"}},
		{"corestring", ConfFileVar, false, false, false, false, CoreUpdateErr{k: "corestring"}},
		{"corestring", ConfFileVar, true, false, false, false, CoreUpdateErr{k: "corestring"}},
		// 155
		{"corestring", ConfFileVar, false, true, false, false, CoreUpdateErr{k: "corestring"}},
		{"corestring", ConfFileVar, true, true, false, false, CoreUpdateErr{k: "corestring"}},
		{"corestring", ConfFileVar, false, false, true, false, CoreUpdateErr{k: "corestring"}},
		{"corestring", ConfFileVar, false, true, true, false, CoreUpdateErr{k: "corestring"}},
		{"corestring", ConfFileVar, true, false, true, false, CoreUpdateErr{k: "corestring"}},
		// 160
		{"corestring", ConfFileVar, true, true, true, false, CoreUpdateErr{k: "corestring"}},
		{"x-corestring", ConfFileVar, false, false, false, false, SettingNotFoundErr{name: "x-corestring"}},
		{"corestring", Env, false, false, false, false, CoreUpdateErr{k: "corestring"}},
		{"corestring", Env, true, false, false, false, CoreUpdateErr{k: "corestring"}},
		{"corestring", Env, false, true, false, false, CoreUpdateErr{k: "corestring"}},
		// 165
		{"corestring", Env, true, true, false, false, CoreUpdateErr{k: "corestring"}},
		{"corestring", Env, false, false, true, false, CoreUpdateErr{k: "corestring"}},
		{"corestring", Env, false, true, true, false, CoreUpdateErr{k: "corestring"}},
		{"corestring", Env, true, false, true, false, CoreUpdateErr{k: "corestring"}},
		{"corestring", Env, true, true, true, false, CoreUpdateErr{k: "corestring"}},
		// 170
		{"x-corestring", Env, false, false, false, false, SettingNotFoundErr{name: "x-corestring"}},
		{"corestring", Flag, false, false, false, false, CoreUpdateErr{k: "corestring"}},
		{"corestring", Flag, true, false, false, false, CoreUpdateErr{k: "corestring"}},
		{"corestring", Flag, false, true, false, false, CoreUpdateErr{k: "corestring"}},
		{"corestring", Flag, true, true, false, false, CoreUpdateErr{k: "corestring"}},
		// 175
		{"corestring", Flag, false, false, true, false, CoreUpdateErr{k: "corestring"}},
		{"corestring", Flag, false, true, true, false, CoreUpdateErr{k: "corestring"}},
		{"corestring", Flag, true, false, true, false, CoreUpdateErr{k: "corestring"}},
		{"corestring", Flag, true, true, true, false, CoreUpdateErr{k: "corestring"}},
		{"x-corestring", Flag, false, false, false, false, SettingNotFoundErr{name: "x-corestring"}},
		// 180
		{"cfgbool", Basic, false, false, false, false, BasicUpdateErr{typ: ConfFileVar, k: "cfgbool"}},
		{"cfgbool", Basic, true, false, false, false, BasicUpdateErr{typ: ConfFileVar, k: "cfgbool"}},
		{"cfgbool", Basic, false, true, false, false, BasicUpdateErr{typ: ConfFileVar, k: "cfgbool"}},
		{"cfgbool", Basic, true, true, false, false, BasicUpdateErr{typ: ConfFileVar, k: "cfgbool"}},
		{"cfgbool", Basic, false, false, true, false, BasicUpdateErr{typ: ConfFileVar, k: "cfgbool"}},
		// 185
		{"cfgbool", Basic, false, true, true, false, BasicUpdateErr{typ: ConfFileVar, k: "cfgbool"}},
		{"cfgbool", Basic, true, false, true, false, BasicUpdateErr{typ: ConfFileVar, k: "cfgbool"}},
		{"cfgbool", Basic, true, true, true, false, BasicUpdateErr{typ: ConfFileVar, k: "cfgbool"}},
		{"x-cfgbool", Basic, false, false, false, false, SettingNotFoundErr{name: "x-cfgbool"}},
		{"cfgbool", Core, false, false, false, false, UpdateErr{typ: Core, k: "cfgbool", slug: "invalid update type"}},
		// 190
		{"cfgbool", Core, true, false, false, false, UpdateErr{typ: Core, k: "cfgbool", slug: "invalid update type"}},
		{"cfgbool", Core, false, true, false, false, UpdateErr{typ: Core, k: "cfgbool", slug: "invalid update type"}},
		{"cfgbool", Core, true, true, false, false, UpdateErr{typ: Core, k: "cfgbool", slug: "invalid update type"}},
		{"cfgbool", Core, false, false, true, false, UpdateErr{typ: Core, k: "cfgbool", slug: "invalid update type"}},
		{"cfgbool", Core, false, true, true, false, UpdateErr{typ: Core, k: "cfgbool", slug: "invalid update type"}},
		// 195
		{"cfgbool", Core, true, false, true, false, UpdateErr{typ: Core, k: "cfgbool", slug: "invalid update type"}},
		{"cfgbool", Core, true, true, true, false, UpdateErr{typ: Core, k: "cfgbool", slug: "invalid update type"}},
		{"x-cfgbool", Core, false, false, false, false, SettingNotFoundErr{name: "x-cfgbool"}},
		{"cfgbool", ConfFileVar, false, false, false, true, nil},
		{"cfgbool", ConfFileVar, true, false, false, false, UpdateErr{typ: ConfFileVar, k: "cfgbool", slug: "already set from the configuration file"}},
		// 200
		{"cfgbool", ConfFileVar, false, true, false, false, UpdateErr{typ: ConfFileVar, k: "cfgbool", slug: "already set from env vars"}},
		{"cfgbool", ConfFileVar, true, true, false, false, UpdateErr{typ: ConfFileVar, k: "cfgbool", slug: "already set from env vars"}},
		{"cfgbool", ConfFileVar, false, false, true, false, UpdateErr{typ: ConfFileVar, k: "cfgbool", slug: "already set from flags"}},
		{"cfgbool", ConfFileVar, false, true, true, false, UpdateErr{typ: ConfFileVar, k: "cfgbool", slug: "already set from flags"}},
		{"cfgbool", ConfFileVar, true, false, true, false, UpdateErr{typ: ConfFileVar, k: "cfgbool", slug: "already set from flags"}},
		// 205
		{"cfgbool", ConfFileVar, true, true, true, false, UpdateErr{typ: ConfFileVar, k: "cfgbool", slug: "already set from flags"}},
		{"x-cfgbool", ConfFileVar, false, false, false, false, SettingNotFoundErr{name: "x-cfgbool"}},
		{"cfgbool", Env, false, false, false, false, UpdateErr{typ: Env, k: "cfgbool", slug: "is not an env var"}},
		{"cfgbool", Env, true, false, false, false, UpdateErr{typ: Env, k: "cfgbool", slug: "is not an env var"}},
		{"cfgbool", Env, false, true, false, false, UpdateErr{typ: Env, k: "cfgbool", slug: "is not an env var"}},
		// 210
		{"cfgbool", Env, true, true, false, false, UpdateErr{typ: Env, k: "cfgbool", slug: "is not an env var"}},
		{"cfgbool", Env, false, false, true, false, UpdateErr{typ: Env, k: "cfgbool", slug: "is not an env var"}},
		{"cfgbool", Env, false, true, true, false, UpdateErr{typ: Env, k: "cfgbool", slug: "is not an env var"}},
		{"cfgbool", Env, true, false, true, false, UpdateErr{typ: Env, k: "cfgbool", slug: "is not an env var"}},
		{"cfgbool", Env, true, true, true, false, UpdateErr{typ: Env, k: "cfgbool", slug: "is not an env var"}},
		// 215
		{"x-cfgbool", Env, false, false, false, false, SettingNotFoundErr{name: "x-cfgbool"}},
		{"cfgbool", Flag, false, false, false, false, UpdateErr{typ: Flag, k: "cfgbool", slug: "is not a flag"}},
		{"cfgbool", Flag, true, false, false, false, UpdateErr{typ: Flag, k: "cfgbool", slug: "is not a flag"}},
		{"cfgbool", Flag, false, true, false, false, UpdateErr{typ: Flag, k: "cfgbool", slug: "is not a flag"}},
		{"cfgbool", Flag, true, true, false, false, UpdateErr{typ: Flag, k: "cfgbool", slug: "is not a flag"}},
		// 220
		{"cfgbool", Flag, false, false, true, false, UpdateErr{typ: Flag, k: "cfgbool", slug: "is not a flag"}},
		{"cfgbool", Flag, false, true, true, false, UpdateErr{typ: Flag, k: "cfgbool", slug: "is not a flag"}},
		{"cfgbool", Flag, true, false, true, false, UpdateErr{typ: Flag, k: "cfgbool", slug: "is not a flag"}},
		{"cfgbool", Flag, true, true, true, false, UpdateErr{typ: Flag, k: "cfgbool", slug: "is not a flag"}},
		{"x-cfgbool", Flag, false, false, false, false, SettingNotFoundErr{name: "x-cfgbool"}},
		// 225
		{"cfgint", Basic, false, false, false, false, BasicUpdateErr{typ: ConfFileVar, k: "cfgint"}},
		{"cfgint", Basic, true, false, false, false, BasicUpdateErr{typ: ConfFileVar, k: "cfgint"}},
		{"cfgint", Basic, false, true, false, false, BasicUpdateErr{typ: ConfFileVar, k: "cfgint"}},
		{"cfgint", Basic, true, true, false, false, BasicUpdateErr{typ: ConfFileVar, k: "cfgint"}},
		{"cfgint", Basic, false, false, true, false, BasicUpdateErr{typ: ConfFileVar, k: "cfgint"}},
		// 230
		{"cfgint", Basic, false, true, true, false, BasicUpdateErr{typ: ConfFileVar, k: "cfgint"}},
		{"cfgint", Basic, true, false, true, false, BasicUpdateErr{typ: ConfFileVar, k: "cfgint"}},
		{"cfgint", Basic, true, true, true, false, BasicUpdateErr{typ: ConfFileVar, k: "cfgint"}},
		{"x-cfgint", Basic, false, false, false, false, SettingNotFoundErr{name: "x-cfgint"}},
		{"cfgint", Core, false, false, false, false, UpdateErr{typ: Core, k: "cfgint", slug: "invalid update type"}},
		// 235
		{"cfgint", Core, true, false, false, false, UpdateErr{typ: Core, k: "cfgint", slug: "invalid update type"}},
		{"cfgint", Core, false, true, false, false, UpdateErr{typ: Core, k: "cfgint", slug: "invalid update type"}},
		{"cfgint", Core, true, true, false, false, UpdateErr{typ: Core, k: "cfgint", slug: "invalid update type"}},
		{"cfgint", Core, false, false, true, false, UpdateErr{typ: Core, k: "cfgint", slug: "invalid update type"}},
		{"cfgint", Core, false, true, true, false, UpdateErr{typ: Core, k: "cfgint", slug: "invalid update type"}},
		// 240
		{"cfgint", Core, true, false, true, false, UpdateErr{typ: Core, k: "cfgint", slug: "invalid update type"}},
		{"cfgint", Core, true, true, true, false, UpdateErr{typ: Core, k: "cfgint", slug: "invalid update type"}},
		{"x-cfgint", Core, false, false, false, false, SettingNotFoundErr{name: "x-cfgint"}},
		{"cfgint", ConfFileVar, false, false, false, true, nil},
		{"cfgint", ConfFileVar, true, false, false, false, UpdateErr{typ: ConfFileVar, k: "cfgint", slug: "already set from the configuration file"}},
		// 245
		{"cfgint", ConfFileVar, false, true, false, false, UpdateErr{typ: ConfFileVar, k: "cfgint", slug: "already set from env vars"}},
		{"cfgint", ConfFileVar, true, true, false, false, UpdateErr{typ: ConfFileVar, k: "cfgint", slug: "already set from env vars"}},
		{"cfgint", ConfFileVar, false, false, true, false, UpdateErr{typ: ConfFileVar, k: "cfgint", slug: "already set from flags"}},
		{"cfgint", ConfFileVar, false, true, true, false, UpdateErr{typ: ConfFileVar, k: "cfgint", slug: "already set from flags"}},
		{"cfgint", ConfFileVar, true, false, true, false, UpdateErr{typ: ConfFileVar, k: "cfgint", slug: "already set from flags"}},
		// 250
		{"cfgint", ConfFileVar, true, true, true, false, UpdateErr{typ: ConfFileVar, k: "cfgint", slug: "already set from flags"}},
		{"x-cfgint", ConfFileVar, false, false, false, false, SettingNotFoundErr{name: "x-cfgint"}},
		{"cfgint", Env, false, false, false, false, UpdateErr{typ: Env, k: "cfgint", slug: "is not an env var"}},
		{"cfgint", Env, true, false, false, false, UpdateErr{typ: Env, k: "cfgint", slug: "is not an env var"}},
		{"cfgint", Env, false, true, false, false, UpdateErr{typ: Env, k: "cfgint", slug: "is not an env var"}},
		// 255
		{"cfgint", Env, true, true, false, false, UpdateErr{typ: Env, k: "cfgint", slug: "is not an env var"}},
		{"cfgint", Env, false, false, true, false, UpdateErr{typ: Env, k: "cfgint", slug: "is not an env var"}},
		{"cfgint", Env, false, true, true, false, UpdateErr{typ: Env, k: "cfgint", slug: "is not an env var"}},
		{"cfgint", Env, true, false, true, false, UpdateErr{typ: Env, k: "cfgint", slug: "is not an env var"}},
		{"cfgint", Env, true, true, true, false, UpdateErr{typ: Env, k: "cfgint", slug: "is not an env var"}},
		// 260
		{"x-cfgint", Env, false, false, false, false, SettingNotFoundErr{name: "x-cfgint"}},
		{"cfgint", Flag, false, false, false, false, UpdateErr{typ: Flag, k: "cfgint", slug: "is not a flag"}},
		{"cfgint", Flag, true, false, false, false, UpdateErr{typ: Flag, k: "cfgint", slug: "is not a flag"}},
		{"cfgint", Flag, false, true, false, false, UpdateErr{typ: Flag, k: "cfgint", slug: "is not a flag"}},
		{"cfgint", Flag, true, true, false, false, UpdateErr{typ: Flag, k: "cfgint", slug: "is not a flag"}},
		// 265
		{"cfgint", Flag, false, false, true, false, UpdateErr{typ: Flag, k: "cfgint", slug: "is not a flag"}},
		{"cfgint", Flag, false, true, true, false, UpdateErr{typ: Flag, k: "cfgint", slug: "is not a flag"}},
		{"cfgint", Flag, true, false, true, false, UpdateErr{typ: Flag, k: "cfgint", slug: "is not a flag"}},
		{"cfgint", Flag, true, true, true, false, UpdateErr{typ: Flag, k: "cfgint", slug: "is not a flag"}},
		{"x-cfgint", Flag, false, false, false, false, SettingNotFoundErr{name: "x-cfgint"}},
		// 270
		{"cfgint64", Basic, false, false, false, false, BasicUpdateErr{typ: ConfFileVar, k: "cfgint64"}},
		{"cfgint64", Basic, true, false, false, false, BasicUpdateErr{typ: ConfFileVar, k: "cfgint64"}},
		{"cfgint64", Basic, false, true, false, false, BasicUpdateErr{typ: ConfFileVar, k: "cfgint64"}},
		{"cfgint64", Basic, true, true, false, false, BasicUpdateErr{typ: ConfFileVar, k: "cfgint64"}},
		{"cfgint64", Basic, false, false, true, false, BasicUpdateErr{typ: ConfFileVar, k: "cfgint64"}},
		// 275
		{"cfgint64", Basic, false, true, true, false, BasicUpdateErr{typ: ConfFileVar, k: "cfgint64"}},
		{"cfgint64", Basic, true, false, true, false, BasicUpdateErr{typ: ConfFileVar, k: "cfgint64"}},
		{"cfgint64", Basic, true, true, true, false, BasicUpdateErr{typ: ConfFileVar, k: "cfgint64"}},
		{"x-cfgint64", Basic, false, false, false, false, SettingNotFoundErr{name: "x-cfgint64"}},
		{"cfgint64", Core, false, false, false, false, UpdateErr{typ: Core, k: "cfgint64", slug: "invalid update type"}},
		// 280
		{"cfgint64", Core, true, false, false, false, UpdateErr{typ: Core, k: "cfgint64", slug: "invalid update type"}},
		{"cfgint64", Core, false, true, false, false, UpdateErr{typ: Core, k: "cfgint64", slug: "invalid update type"}},
		{"cfgint64", Core, true, true, false, false, UpdateErr{typ: Core, k: "cfgint64", slug: "invalid update type"}},
		{"cfgint64", Core, false, false, true, false, UpdateErr{typ: Core, k: "cfgint64", slug: "invalid update type"}},
		{"cfgint64", Core, false, true, true, false, UpdateErr{typ: Core, k: "cfgint64", slug: "invalid update type"}},
		// 285
		{"cfgint64", Core, true, false, true, false, UpdateErr{typ: Core, k: "cfgint64", slug: "invalid update type"}},
		{"cfgint64", Core, true, true, true, false, UpdateErr{typ: Core, k: "cfgint64", slug: "invalid update type"}},
		{"x-cfgint64", Core, false, false, false, false, SettingNotFoundErr{name: "x-cfgint64"}},
		{"cfgint64", ConfFileVar, false, false, false, true, nil},
		{"cfgint64", ConfFileVar, true, false, false, false, UpdateErr{typ: ConfFileVar, k: "cfgint64", slug: "already set from the configuration file"}},
		// 290
		{"cfgint64", ConfFileVar, false, true, false, false, UpdateErr{typ: ConfFileVar, k: "cfgint64", slug: "already set from env vars"}},
		{"cfgint64", ConfFileVar, true, true, false, false, UpdateErr{typ: ConfFileVar, k: "cfgint64", slug: "already set from env vars"}},
		{"cfgint64", ConfFileVar, false, false, true, false, UpdateErr{typ: ConfFileVar, k: "cfgint64", slug: "already set from flags"}},
		{"cfgint64", ConfFileVar, false, true, true, false, UpdateErr{typ: ConfFileVar, k: "cfgint64", slug: "already set from flags"}},
		{"cfgint64", ConfFileVar, true, false, true, false, UpdateErr{typ: ConfFileVar, k: "cfgint64", slug: "already set from flags"}},
		// 295
		{"cfgint64", ConfFileVar, true, true, true, false, UpdateErr{typ: ConfFileVar, k: "cfgint64", slug: "already set from flags"}},
		{"x-cfgint64", ConfFileVar, false, false, false, false, SettingNotFoundErr{name: "x-cfgint64"}},
		{"cfgint64", Env, false, false, false, false, UpdateErr{typ: Env, k: "cfgint64", slug: "is not an env var"}},
		{"cfgint64", Env, true, false, false, false, UpdateErr{typ: Env, k: "cfgint64", slug: "is not an env var"}},
		{"cfgint64", Env, false, true, false, false, UpdateErr{typ: Env, k: "cfgint64", slug: "is not an env var"}},
		// 300
		{"cfgint64", Env, true, true, false, false, UpdateErr{typ: Env, k: "cfgint64", slug: "is not an env var"}},
		{"cfgint64", Env, false, false, true, false, UpdateErr{typ: Env, k: "cfgint64", slug: "is not an env var"}},
		{"cfgint64", Env, false, true, true, false, UpdateErr{typ: Env, k: "cfgint64", slug: "is not an env var"}},
		{"cfgint64", Env, true, false, true, false, UpdateErr{typ: Env, k: "cfgint64", slug: "is not an env var"}},
		{"cfgint64", Env, true, true, true, false, UpdateErr{typ: Env, k: "cfgint64", slug: "is not an env var"}},
		// 305
		{"x-cfgint64", Env, false, false, false, false, SettingNotFoundErr{name: "x-cfgint64"}},
		{"cfgint64", Flag, false, false, false, false, UpdateErr{typ: Flag, k: "cfgint64", slug: "is not a flag"}},
		{"cfgint64", Flag, true, false, false, false, UpdateErr{typ: Flag, k: "cfgint64", slug: "is not a flag"}},
		{"cfgint64", Flag, false, true, false, false, UpdateErr{typ: Flag, k: "cfgint64", slug: "is not a flag"}},
		{"cfgint64", Flag, true, true, false, false, UpdateErr{typ: Flag, k: "cfgint64", slug: "is not a flag"}},
		// 310
		{"cfgint64", Flag, false, false, true, false, UpdateErr{typ: Flag, k: "cfgint64", slug: "is not a flag"}},
		{"cfgint64", Flag, false, true, true, false, UpdateErr{typ: Flag, k: "cfgint64", slug: "is not a flag"}},
		{"cfgint64", Flag, true, false, true, false, UpdateErr{typ: Flag, k: "cfgint64", slug: "is not a flag"}},
		{"cfgint64", Flag, true, true, true, false, UpdateErr{typ: Flag, k: "cfgint64", slug: "is not a flag"}},
		{"x-cfgint64", Flag, false, false, false, false, SettingNotFoundErr{name: "x-cfgint64"}},
		// 315
		{"cfgstring", Basic, false, false, false, false, BasicUpdateErr{typ: ConfFileVar, k: "cfgstring"}},
		{"cfgstring", Basic, true, false, false, false, BasicUpdateErr{typ: ConfFileVar, k: "cfgstring"}},
		{"cfgstring", Basic, false, true, false, false, BasicUpdateErr{typ: ConfFileVar, k: "cfgstring"}},
		{"cfgstring", Basic, true, true, false, false, BasicUpdateErr{typ: ConfFileVar, k: "cfgstring"}},
		{"cfgstring", Basic, false, false, true, false, BasicUpdateErr{typ: ConfFileVar, k: "cfgstring"}},
		// 320
		{"cfgstring", Basic, false, true, true, false, BasicUpdateErr{typ: ConfFileVar, k: "cfgstring"}},
		{"cfgstring", Basic, true, false, true, false, BasicUpdateErr{typ: ConfFileVar, k: "cfgstring"}},
		{"cfgstring", Basic, true, true, true, false, BasicUpdateErr{typ: ConfFileVar, k: "cfgstring"}},
		{"x-cfgstring", Basic, false, false, false, false, SettingNotFoundErr{name: "x-cfgstring"}},
		{"cfgstring", Core, false, false, false, false, UpdateErr{typ: Core, k: "cfgstring", slug: "invalid update type"}},
		// 325
		{"cfgstring", Core, true, false, false, false, UpdateErr{typ: Core, k: "cfgstring", slug: "invalid update type"}},
		{"cfgstring", Core, false, true, false, false, UpdateErr{typ: Core, k: "cfgstring", slug: "invalid update type"}},
		{"cfgstring", Core, true, true, false, false, UpdateErr{typ: Core, k: "cfgstring", slug: "invalid update type"}},
		{"cfgstring", Core, false, false, true, false, UpdateErr{typ: Core, k: "cfgstring", slug: "invalid update type"}},
		{"cfgstring", Core, false, true, true, false, UpdateErr{typ: Core, k: "cfgstring", slug: "invalid update type"}},
		// 330
		{"cfgstring", Core, true, false, true, false, UpdateErr{typ: Core, k: "cfgstring", slug: "invalid update type"}},
		{"cfgstring", Core, true, true, true, false, UpdateErr{typ: Core, k: "cfgstring", slug: "invalid update type"}},
		{"x-cfgstring", Core, false, false, false, false, SettingNotFoundErr{name: "x-cfgstring"}},
		{"cfgstring", ConfFileVar, false, false, false, true, nil},
		{"cfgstring", ConfFileVar, true, false, false, false, UpdateErr{typ: ConfFileVar, k: "cfgstring", slug: "already set from the configuration file"}},
		// 335
		{"cfgstring", ConfFileVar, false, true, false, false, UpdateErr{typ: ConfFileVar, k: "cfgstring", slug: "already set from env vars"}},
		{"cfgstring", ConfFileVar, true, true, false, false, UpdateErr{typ: ConfFileVar, k: "cfgstring", slug: "already set from env vars"}},
		{"cfgstring", ConfFileVar, false, false, true, false, UpdateErr{typ: ConfFileVar, k: "cfgstring", slug: "already set from flags"}},
		{"cfgstring", ConfFileVar, false, true, true, false, UpdateErr{typ: ConfFileVar, k: "cfgstring", slug: "already set from flags"}},
		{"cfgstring", ConfFileVar, true, false, true, false, UpdateErr{typ: ConfFileVar, k: "cfgstring", slug: "already set from flags"}},
		// 340
		{"cfgstring", ConfFileVar, true, true, true, false, UpdateErr{typ: ConfFileVar, k: "cfgstring", slug: "already set from flags"}},
		{"x-cfgstring", ConfFileVar, false, false, false, false, SettingNotFoundErr{name: "x-cfgstring"}},
		{"cfgstring", Env, false, false, false, false, UpdateErr{typ: Env, k: "cfgstring", slug: "is not an env var"}},
		{"cfgstring", Env, true, false, false, false, UpdateErr{typ: Env, k: "cfgstring", slug: "is not an env var"}},
		{"cfgstring", Env, false, true, false, false, UpdateErr{typ: Env, k: "cfgstring", slug: "is not an env var"}},
		// 345
		{"cfgstring", Env, true, true, false, false, UpdateErr{typ: Env, k: "cfgstring", slug: "is not an env var"}},
		{"cfgstring", Env, false, false, true, false, UpdateErr{typ: Env, k: "cfgstring", slug: "is not an env var"}},
		{"cfgstring", Env, false, true, true, false, UpdateErr{typ: Env, k: "cfgstring", slug: "is not an env var"}},
		{"cfgstring", Env, true, false, true, false, UpdateErr{typ: Env, k: "cfgstring", slug: "is not an env var"}},
		{"cfgstring", Env, true, true, true, false, UpdateErr{typ: Env, k: "cfgstring", slug: "is not an env var"}},
		// 350
		{"x-cfgstring", Env, false, false, false, false, SettingNotFoundErr{name: "x-cfgstring"}},
		{"cfgstring", Flag, false, false, false, false, UpdateErr{typ: Flag, k: "cfgstring", slug: "is not a flag"}},
		{"cfgstring", Flag, true, false, false, false, UpdateErr{typ: Flag, k: "cfgstring", slug: "is not a flag"}},
		{"cfgstring", Flag, false, true, false, false, UpdateErr{typ: Flag, k: "cfgstring", slug: "is not a flag"}},
		{"cfgstring", Flag, true, true, false, false, UpdateErr{typ: Flag, k: "cfgstring", slug: "is not a flag"}},
		// 355
		{"cfgstring", Flag, false, false, true, false, UpdateErr{typ: Flag, k: "cfgstring", slug: "is not a flag"}},
		{"cfgstring", Flag, false, true, true, false, UpdateErr{typ: Flag, k: "cfgstring", slug: "is not a flag"}},
		{"cfgstring", Flag, true, false, true, false, UpdateErr{typ: Flag, k: "cfgstring", slug: "is not a flag"}},
		{"cfgstring", Flag, true, true, true, false, UpdateErr{typ: Flag, k: "cfgstring", slug: "is not a flag"}},
		{"x-cfgstring", Flag, false, false, false, false, SettingNotFoundErr{name: "x-cfgstring"}},

		// 360
		{"envbool", Basic, false, false, false, false, BasicUpdateErr{typ: Env, k: "envbool"}},
		{"envbool", Basic, true, false, false, false, BasicUpdateErr{typ: Env, k: "envbool"}},
		{"envbool", Basic, false, true, false, false, BasicUpdateErr{typ: Env, k: "envbool"}},
		{"envbool", Basic, true, true, false, false, BasicUpdateErr{typ: Env, k: "envbool"}},
		{"envbool", Basic, false, false, true, false, BasicUpdateErr{typ: Env, k: "envbool"}},
		// 365
		{"envbool", Basic, false, true, true, false, BasicUpdateErr{typ: Env, k: "envbool"}},
		{"envbool", Basic, true, false, true, false, BasicUpdateErr{typ: Env, k: "envbool"}},
		{"envbool", Basic, true, true, true, false, BasicUpdateErr{typ: Env, k: "envbool"}},
		{"x-envbool", Basic, false, false, false, false, SettingNotFoundErr{name: "x-envbool"}},

		{"envbool", Core, false, false, false, false, UpdateErr{typ: Core, k: "envbool", slug: "invalid update type"}},
		// 370
		{"envbool", Core, true, false, false, false, UpdateErr{typ: Core, k: "envbool", slug: "invalid update type"}},
		{"envbool", Core, false, true, false, false, UpdateErr{typ: Core, k: "envbool", slug: "invalid update type"}},
		{"envbool", Core, true, true, false, false, UpdateErr{typ: Core, k: "envbool", slug: "invalid update type"}},
		{"envbool", Core, false, false, true, false, UpdateErr{typ: Core, k: "envbool", slug: "invalid update type"}},
		{"envbool", Core, false, true, true, false, UpdateErr{typ: Core, k: "envbool", slug: "invalid update type"}},
		// 375
		{"envbool", Core, true, false, true, false, UpdateErr{typ: Core, k: "envbool", slug: "invalid update type"}},
		{"envbool", Core, true, true, true, false, UpdateErr{typ: Core, k: "envbool", slug: "invalid update type"}},
		{"x-envbool", Core, false, false, false, false, SettingNotFoundErr{name: "x-envbool"}},
		{"envbool", ConfFileVar, false, false, false, true, nil},
		{"envbool", ConfFileVar, true, false, false, false, UpdateErr{typ: ConfFileVar, k: "envbool", slug: "already set from the configuration file"}},
		// 380
		{"envbool", ConfFileVar, false, true, false, false, UpdateErr{typ: ConfFileVar, k: "envbool", slug: "already set from env vars"}},
		{"envbool", ConfFileVar, true, true, false, false, UpdateErr{typ: ConfFileVar, k: "envbool", slug: "already set from env vars"}},
		{"envbool", ConfFileVar, false, false, true, false, UpdateErr{typ: ConfFileVar, k: "envbool", slug: "already set from flags"}},
		{"envbool", ConfFileVar, false, true, true, false, UpdateErr{typ: ConfFileVar, k: "envbool", slug: "already set from flags"}},
		{"envbool", ConfFileVar, true, false, true, false, UpdateErr{typ: ConfFileVar, k: "envbool", slug: "already set from flags"}},
		// 385
		{"envbool", ConfFileVar, true, true, true, false, UpdateErr{typ: ConfFileVar, k: "envbool", slug: "already set from flags"}},
		{"x-envbool", ConfFileVar, false, false, false, false, SettingNotFoundErr{name: "x-envbool"}},
		{"envbool", Env, false, false, false, true, nil},
		{"envbool", Env, true, false, false, true, nil},
		{"envbool", Env, false, true, false, false, UpdateErr{typ: Env, k: "envbool", slug: "already set from env vars"}},
		// 390
		{"envbool", Env, true, true, false, false, UpdateErr{typ: Env, k: "envbool", slug: "already set from env vars"}},
		{"envbool", Env, false, false, true, false, UpdateErr{typ: Env, k: "envbool", slug: "already set from flags"}},
		{"envbool", Env, false, true, true, false, UpdateErr{typ: Env, k: "envbool", slug: "already set from flags"}},
		{"envbool", Env, true, false, true, false, UpdateErr{typ: Env, k: "envbool", slug: "already set from flags"}},
		{"envbool", Env, true, true, true, false, UpdateErr{typ: Env, k: "envbool", slug: "already set from flags"}},
		// 395
		{"x-envbool", Env, false, false, false, false, SettingNotFoundErr{name: "x-envbool"}},
		{"envbool", Flag, false, false, false, false, UpdateErr{typ: Flag, k: "envbool", slug: "is not a flag"}},
		{"envbool", Flag, true, false, false, false, UpdateErr{typ: Flag, k: "envbool", slug: "is not a flag"}},
		{"envbool", Flag, false, true, false, false, UpdateErr{typ: Flag, k: "envbool", slug: "is not a flag"}},
		{"envbool", Flag, true, true, false, false, UpdateErr{typ: Flag, k: "envbool", slug: "is not a flag"}},
		// 400
		{"envbool", Flag, false, false, true, false, UpdateErr{typ: Flag, k: "envbool", slug: "is not a flag"}},
		{"envbool", Flag, false, true, true, false, UpdateErr{typ: Flag, k: "envbool", slug: "is not a flag"}},
		{"envbool", Flag, true, false, true, false, UpdateErr{typ: Flag, k: "envbool", slug: "is not a flag"}},
		{"envbool", Flag, true, true, true, false, UpdateErr{typ: Flag, k: "envbool", slug: "is not a flag"}},
		{"x-envbool", Flag, false, false, false, false, SettingNotFoundErr{name: "x-envbool"}},
		// 405
		{"envint", Basic, false, false, false, false, BasicUpdateErr{typ: Env, k: "envint"}},
		{"envint", Basic, true, false, false, false, BasicUpdateErr{typ: Env, k: "envint"}},
		{"envint", Basic, false, true, false, false, BasicUpdateErr{typ: Env, k: "envint"}},
		{"envint", Basic, true, true, false, false, BasicUpdateErr{typ: Env, k: "envint"}},
		{"envint", Basic, false, false, true, false, BasicUpdateErr{typ: Env, k: "envint"}},
		// 410
		{"envint", Basic, false, true, true, false, BasicUpdateErr{typ: Env, k: "envint"}},
		{"envint", Basic, true, false, true, false, BasicUpdateErr{typ: Env, k: "envint"}},
		{"envint", Basic, true, true, true, false, BasicUpdateErr{typ: Env, k: "envint"}},
		{"x-envint", Basic, false, false, false, false, SettingNotFoundErr{name: "x-envint"}},
		{"envint", Core, false, false, false, false, UpdateErr{typ: Core, k: "envint", slug: "invalid update type"}},
		// 415
		{"envint", Core, true, false, false, false, UpdateErr{typ: Core, k: "envint", slug: "invalid update type"}},
		{"envint", Core, false, true, false, false, UpdateErr{typ: Core, k: "envint", slug: "invalid update type"}},
		{"envint", Core, true, true, false, false, UpdateErr{typ: Core, k: "envint", slug: "invalid update type"}},
		{"envint", Core, false, false, true, false, UpdateErr{typ: Core, k: "envint", slug: "invalid update type"}},
		{"envint", Core, false, true, true, false, UpdateErr{typ: Core, k: "envint", slug: "invalid update type"}},
		// 420
		{"envint", Core, true, false, true, false, UpdateErr{typ: Core, k: "envint", slug: "invalid update type"}},
		{"envint", Core, true, true, true, false, UpdateErr{typ: Core, k: "envint", slug: "invalid update type"}},
		{"x-envint", Core, false, false, false, false, SettingNotFoundErr{name: "x-envint"}},
		{"envint", ConfFileVar, false, false, false, true, nil},
		{"envint", ConfFileVar, true, false, false, false, UpdateErr{typ: ConfFileVar, k: "envint", slug: "already set from the configuration file"}},
		// 425
		{"envint", ConfFileVar, false, true, false, false, UpdateErr{typ: ConfFileVar, k: "envint", slug: "already set from env vars"}},
		{"envint", ConfFileVar, true, true, false, false, UpdateErr{typ: ConfFileVar, k: "envint", slug: "already set from env vars"}},
		{"envint", ConfFileVar, false, false, true, false, UpdateErr{typ: ConfFileVar, k: "envint", slug: "already set from flags"}},
		{"envint", ConfFileVar, false, true, true, false, UpdateErr{typ: ConfFileVar, k: "envint", slug: "already set from flags"}},
		{"envint", ConfFileVar, true, false, true, false, UpdateErr{typ: ConfFileVar, k: "envint", slug: "already set from flags"}},
		// 430
		{"envint", ConfFileVar, true, true, true, false, UpdateErr{typ: ConfFileVar, k: "envint", slug: "already set from flags"}},
		{"x-envint", ConfFileVar, false, false, false, false, SettingNotFoundErr{name: "x-envint"}},
		{"envint", Env, false, false, false, true, nil},
		{"envint", Env, true, false, false, true, nil},
		{"envint", Env, false, true, false, false, UpdateErr{typ: Env, k: "envint", slug: "already set from env vars"}},
		// 435
		{"envint", Env, true, true, false, false, UpdateErr{typ: Env, k: "envint", slug: "already set from env vars"}},
		{"envint", Env, false, false, true, false, UpdateErr{typ: Env, k: "envint", slug: "already set from flags"}},
		{"envint", Env, false, true, true, false, UpdateErr{typ: Env, k: "envint", slug: "already set from flags"}},
		{"envint", Env, true, false, true, false, UpdateErr{typ: Env, k: "envint", slug: "already set from flags"}},
		{"envint", Env, true, true, true, false, UpdateErr{typ: Env, k: "envint", slug: "already set from flags"}},
		// 440
		{"x-envint", Env, false, false, false, false, SettingNotFoundErr{name: "x-envint"}},
		{"envint", Flag, false, false, false, false, UpdateErr{typ: Flag, k: "envint", slug: "is not a flag"}},
		{"envint", Flag, true, false, false, false, UpdateErr{typ: Flag, k: "envint", slug: "is not a flag"}},
		{"envint", Flag, false, true, false, false, UpdateErr{typ: Flag, k: "envint", slug: "is not a flag"}},
		{"envint", Flag, true, true, false, false, UpdateErr{typ: Flag, k: "envint", slug: "is not a flag"}},
		// 445
		{"envint", Flag, false, false, true, false, UpdateErr{typ: Flag, k: "envint", slug: "is not a flag"}},
		{"envint", Flag, false, true, true, false, UpdateErr{typ: Flag, k: "envint", slug: "is not a flag"}},
		{"envint", Flag, true, false, true, false, UpdateErr{typ: Flag, k: "envint", slug: "is not a flag"}},
		{"envint", Flag, true, true, true, false, UpdateErr{typ: Flag, k: "envint", slug: "is not a flag"}},
		{"x-envint", Flag, false, false, false, false, SettingNotFoundErr{name: "x-envint"}},
		// 450
		{"envint64", Basic, false, false, false, false, BasicUpdateErr{typ: Env, k: "envint64"}},
		{"envint64", Basic, true, false, false, false, BasicUpdateErr{typ: Env, k: "envint64"}},
		{"envint64", Basic, true, false, false, false, BasicUpdateErr{typ: Env, k: "envint64"}},
		{"envint64", Basic, false, true, false, false, BasicUpdateErr{typ: Env, k: "envint64"}},
		{"envint64", Basic, true, true, false, false, BasicUpdateErr{typ: Env, k: "envint64"}},
		// 455
		{"envint64", Basic, false, false, true, false, BasicUpdateErr{typ: Env, k: "envint64"}},
		{"envint64", Basic, false, true, true, false, BasicUpdateErr{typ: Env, k: "envint64"}},
		{"envint64", Basic, true, false, true, false, BasicUpdateErr{typ: Env, k: "envint64"}},
		{"envint64", Basic, true, true, true, false, BasicUpdateErr{typ: Env, k: "envint64"}},
		{"x-envint64", Basic, false, false, false, false, SettingNotFoundErr{name: "x-envint64"}},
		// 460
		{"envint64", Core, false, false, false, false, UpdateErr{typ: Core, k: "envint64", slug: "invalid update type"}},
		{"envint64", Core, true, false, false, false, UpdateErr{typ: Core, k: "envint64", slug: "invalid update type"}},
		{"envint64", Core, false, true, false, false, UpdateErr{typ: Core, k: "envint64", slug: "invalid update type"}},
		{"envint64", Core, true, true, false, false, UpdateErr{typ: Core, k: "envint64", slug: "invalid update type"}},
		{"envint64", Core, false, false, true, false, UpdateErr{typ: Core, k: "envint64", slug: "invalid update type"}},
		// 465
		{"envint64", Core, false, true, true, false, UpdateErr{typ: Core, k: "envint64", slug: "invalid update type"}},
		{"envint64", Core, true, false, true, false, UpdateErr{typ: Core, k: "envint64", slug: "invalid update type"}},
		{"envint64", Core, true, true, true, false, UpdateErr{typ: Core, k: "envint64", slug: "invalid update type"}},
		{"x-envint64", Core, false, false, false, false, SettingNotFoundErr{name: "x-envint64"}},
		{"envint64", ConfFileVar, false, false, false, true, nil},
		// 470
		{"envint64", ConfFileVar, true, false, false, false, UpdateErr{typ: ConfFileVar, k: "envint64", slug: "already set from the configuration file"}},
		{"envint64", ConfFileVar, false, true, false, false, UpdateErr{typ: ConfFileVar, k: "envint64", slug: "already set from env vars"}},
		{"envint64", ConfFileVar, true, true, false, false, UpdateErr{typ: ConfFileVar, k: "envint64", slug: "already set from env vars"}},
		{"envint64", ConfFileVar, false, false, true, false, UpdateErr{typ: ConfFileVar, k: "envint64", slug: "already set from flags"}},
		{"envint64", ConfFileVar, false, true, true, false, UpdateErr{typ: ConfFileVar, k: "envint64", slug: "already set from flags"}},
		// 475
		{"envint64", ConfFileVar, true, false, true, false, UpdateErr{typ: ConfFileVar, k: "envint64", slug: "already set from flags"}},
		{"envint64", ConfFileVar, true, true, true, false, UpdateErr{typ: ConfFileVar, k: "envint64", slug: "already set from flags"}},
		{"x-envint64", ConfFileVar, false, false, false, false, SettingNotFoundErr{name: "x-envint64"}},
		{"envint64", Env, false, false, false, true, nil},
		{"envint64", Env, true, false, false, true, nil},
		// 480
		{"envint64", Env, false, true, false, false, UpdateErr{typ: Env, k: "envint64", slug: "already set from env vars"}},
		{"envint64", Env, true, true, false, false, UpdateErr{typ: Env, k: "envint64", slug: "already set from env vars"}},
		{"envint64", Env, false, false, true, false, UpdateErr{typ: Env, k: "envint64", slug: "already set from flags"}},
		{"envint64", Env, false, true, true, false, UpdateErr{typ: Env, k: "envint64", slug: "already set from flags"}},
		{"envint64", Env, true, false, true, false, UpdateErr{typ: Env, k: "envint64", slug: "already set from flags"}},
		// 485
		{"envint64", Env, true, true, true, false, UpdateErr{typ: Env, k: "envint64", slug: "already set from flags"}},
		{"x-envint64", Env, false, false, false, false, SettingNotFoundErr{name: "x-envint64"}},
		{"envint64", Flag, false, false, false, false, UpdateErr{typ: Flag, k: "envint64", slug: "is not a flag"}},
		{"envint64", Flag, true, false, false, false, UpdateErr{typ: Flag, k: "envint64", slug: "is not a flag"}},
		{"envint64", Flag, false, true, false, false, UpdateErr{typ: Flag, k: "envint64", slug: "is not a flag"}},
		// 490
		{"envint64", Flag, true, true, false, false, UpdateErr{typ: Flag, k: "envint64", slug: "is not a flag"}},
		{"envint64", Flag, false, false, true, false, UpdateErr{typ: Flag, k: "envint64", slug: "is not a flag"}},
		{"envint64", Flag, false, true, true, false, UpdateErr{typ: Flag, k: "envint64", slug: "is not a flag"}},
		{"envint64", Flag, true, false, true, false, UpdateErr{typ: Flag, k: "envint64", slug: "is not a flag"}},
		{"envint64", Flag, true, true, true, false, UpdateErr{typ: Flag, k: "envint64", slug: "is not a flag"}},
		// 495
		{"x-envint64", Flag, false, false, false, false, SettingNotFoundErr{name: "x-envint64"}},
		{"envstring", Basic, false, false, false, false, BasicUpdateErr{typ: Env, k: "envstring"}},
		{"envstring", Basic, true, false, false, false, BasicUpdateErr{typ: Env, k: "envstring"}},
		{"envstring", Basic, false, true, false, false, BasicUpdateErr{typ: Env, k: "envstring"}},
		{"envstring", Basic, true, true, false, false, BasicUpdateErr{typ: Env, k: "envstring"}},
		// 500
		{"envstring", Basic, false, false, true, false, BasicUpdateErr{typ: Env, k: "envstring"}},
		{"envstring", Basic, false, true, true, false, BasicUpdateErr{typ: Env, k: "envstring"}},
		{"envstring", Basic, true, false, true, false, BasicUpdateErr{typ: Env, k: "envstring"}},
		{"envstring", Basic, true, true, true, false, BasicUpdateErr{typ: Env, k: "envstring"}},
		{"x-envstring", Basic, false, false, false, false, SettingNotFoundErr{name: "x-envstring"}},
		// 505
		{"envstring", Core, false, false, false, false, UpdateErr{typ: Core, k: "envstring", slug: "invalid update type"}},
		{"envstring", Core, true, false, false, false, UpdateErr{typ: Core, k: "envstring", slug: "invalid update type"}},
		{"envstring", Core, false, true, false, false, UpdateErr{typ: Core, k: "envstring", slug: "invalid update type"}},
		{"envstring", Core, true, true, false, false, UpdateErr{typ: Core, k: "envstring", slug: "invalid update type"}},
		{"envstring", Core, false, false, true, false, UpdateErr{typ: Core, k: "envstring", slug: "invalid update type"}},
		// 510
		{"envstring", Core, false, true, true, false, UpdateErr{typ: Core, k: "envstring", slug: "invalid update type"}},
		{"envstring", Core, true, false, true, false, UpdateErr{typ: Core, k: "envstring", slug: "invalid update type"}},
		{"envstring", Core, true, true, true, false, UpdateErr{typ: Core, k: "envstring", slug: "invalid update type"}},
		{"x-envstring", Core, false, false, false, false, SettingNotFoundErr{name: "x-envstring"}},
		{"envstring", ConfFileVar, false, false, false, true, nil},
		// 515
		{"envstring", ConfFileVar, true, false, false, false, UpdateErr{typ: ConfFileVar, k: "envstring", slug: "already set from the configuration file"}},
		{"envstring", ConfFileVar, false, true, false, false, UpdateErr{typ: ConfFileVar, k: "envstring", slug: "already set from env vars"}},
		{"envstring", ConfFileVar, true, true, false, false, UpdateErr{typ: ConfFileVar, k: "envstring", slug: "already set from env vars"}},
		{"envstring", ConfFileVar, false, false, true, false, UpdateErr{typ: ConfFileVar, k: "envstring", slug: "already set from flags"}},
		{"envstring", ConfFileVar, false, true, true, false, UpdateErr{typ: ConfFileVar, k: "envstring", slug: "already set from flags"}},
		// 520
		{"envstring", ConfFileVar, true, false, true, false, UpdateErr{typ: ConfFileVar, k: "envstring", slug: "already set from flags"}},
		{"envstring", ConfFileVar, true, true, true, false, UpdateErr{typ: ConfFileVar, k: "envstring", slug: "already set from flags"}},
		{"x-envstring", ConfFileVar, false, false, false, false, SettingNotFoundErr{name: "x-envstring"}},
		{"envstring", Env, false, false, false, true, nil},
		{"envstring", Env, true, false, false, true, nil},
		// 525
		{"envstring", Env, false, true, false, false, UpdateErr{typ: Env, k: "envstring", slug: "already set from env vars"}},
		{"envstring", Env, true, true, false, false, UpdateErr{typ: Env, k: "envstring", slug: "already set from env vars"}},
		{"envstring", Env, false, false, true, false, UpdateErr{typ: Env, k: "envstring", slug: "already set from flags"}},
		{"envstring", Env, false, true, true, false, UpdateErr{typ: Env, k: "envstring", slug: "already set from flags"}},
		{"envstring", Env, true, false, true, false, UpdateErr{typ: Env, k: "envstring", slug: "already set from flags"}},
		// 530
		{"envstring", Env, true, true, true, false, UpdateErr{typ: Env, k: "envstring", slug: "already set from flags"}},
		{"x-envstring", Env, false, false, false, false, SettingNotFoundErr{name: "x-envstring"}},
		{"envstring", Flag, false, false, false, false, UpdateErr{typ: Flag, k: "envstring", slug: "is not a flag"}},
		{"envstring", Flag, true, false, false, false, UpdateErr{typ: Flag, k: "envstring", slug: "is not a flag"}},
		{"envstring", Flag, false, true, false, false, UpdateErr{typ: Flag, k: "envstring", slug: "is not a flag"}},
		// 535
		{"envstring", Flag, true, true, false, false, UpdateErr{typ: Flag, k: "envstring", slug: "is not a flag"}},
		{"envstring", Flag, false, false, true, false, UpdateErr{typ: Flag, k: "envstring", slug: "is not a flag"}},
		{"envstring", Flag, false, true, true, false, UpdateErr{typ: Flag, k: "envstring", slug: "is not a flag"}},
		{"envstring", Flag, true, false, true, false, UpdateErr{typ: Flag, k: "envstring", slug: "is not a flag"}},
		{"envstring", Flag, true, true, true, false, UpdateErr{typ: Flag, k: "envstring", slug: "is not a flag"}},
		// 540
		{"x-envstring", Flag, false, false, false, false, SettingNotFoundErr{name: "x-envstring"}},
		{"flagbool", Basic, false, false, false, false, BasicUpdateErr{typ: Flag, k: "flagbool"}},
		{"flagbool", Basic, true, false, false, false, BasicUpdateErr{typ: Flag, k: "flagbool"}},
		{"flagbool", Basic, false, true, false, false, BasicUpdateErr{typ: Flag, k: "flagbool"}},
		{"flagbool", Basic, true, true, false, false, BasicUpdateErr{typ: Flag, k: "flagbool"}},
		// 545
		{"flagbool", Basic, false, false, true, false, BasicUpdateErr{typ: Flag, k: "flagbool"}},
		{"flagbool", Basic, false, true, true, false, BasicUpdateErr{typ: Flag, k: "flagbool"}},
		{"flagbool", Basic, true, false, true, false, BasicUpdateErr{typ: Flag, k: "flagbool"}},
		{"flagbool", Basic, true, true, true, false, BasicUpdateErr{typ: Flag, k: "flagbool"}},
		{"x-flagbool", Basic, false, false, false, false, SettingNotFoundErr{name: "x-flagbool"}},
		// 550
		{"flagbool", Core, false, false, false, false, UpdateErr{typ: Core, k: "flagbool", slug: "invalid update type"}},
		{"flagbool", Core, true, false, false, false, UpdateErr{typ: Core, k: "flagbool", slug: "invalid update type"}},
		{"flagbool", Core, false, true, false, false, UpdateErr{typ: Core, k: "flagbool", slug: "invalid update type"}},
		{"flagbool", Core, true, true, false, false, UpdateErr{typ: Core, k: "flagbool", slug: "invalid update type"}},
		{"flagbool", Core, false, false, true, false, UpdateErr{typ: Core, k: "flagbool", slug: "invalid update type"}},
		// 555
		{"flagbool", Core, false, true, true, false, UpdateErr{typ: Core, k: "flagbool", slug: "invalid update type"}},
		{"flagbool", Core, true, false, true, false, UpdateErr{typ: Core, k: "flagbool", slug: "invalid update type"}},
		{"flagbool", Core, true, true, true, false, UpdateErr{typ: Core, k: "flagbool", slug: "invalid update type"}},
		{"x-flagbool", Core, false, false, false, false, SettingNotFoundErr{name: "x-flagbool"}},
		{"flagbool", ConfFileVar, false, false, false, true, nil},
		// 560
		{"flagbool", ConfFileVar, true, false, false, false, UpdateErr{typ: ConfFileVar, k: "flagbool", slug: "already set from the configuration file"}},
		{"flagbool", ConfFileVar, false, true, false, false, UpdateErr{typ: ConfFileVar, k: "flagbool", slug: "already set from env vars"}},
		{"flagbool", ConfFileVar, true, true, false, false, UpdateErr{typ: ConfFileVar, k: "flagbool", slug: "already set from env vars"}},
		{"flagbool", ConfFileVar, false, false, true, false, UpdateErr{typ: ConfFileVar, k: "flagbool", slug: "already set from flags"}},
		{"flagbool", ConfFileVar, false, true, true, false, UpdateErr{typ: ConfFileVar, k: "flagbool", slug: "already set from flags"}},
		// 565
		{"flagbool", ConfFileVar, true, false, true, false, UpdateErr{typ: ConfFileVar, k: "flagbool", slug: "already set from flags"}},
		{"flagbool", ConfFileVar, true, true, true, false, UpdateErr{typ: ConfFileVar, k: "flagbool", slug: "already set from flags"}},
		{"x-flagbool", ConfFileVar, false, false, false, false, SettingNotFoundErr{name: "x-flagbool"}},
		{"flagbool", Env, false, false, false, true, nil},
		{"flagbool", Env, true, false, false, true, nil},
		// 570
		{"flagbool", Env, false, true, false, false, UpdateErr{typ: Env, k: "flagbool", slug: "already set from env vars"}},
		{"flagbool", Env, true, true, false, false, UpdateErr{typ: Env, k: "flagbool", slug: "already set from env vars"}},
		{"flagbool", Env, false, false, true, false, UpdateErr{typ: Env, k: "flagbool", slug: "already set from flags"}},
		{"flagbool", Env, false, true, true, false, UpdateErr{typ: Env, k: "flagbool", slug: "already set from flags"}},
		{"flagbool", Env, true, false, true, false, UpdateErr{typ: Env, k: "flagbool", slug: "already set from flags"}},
		// 575
		{"flagbool", Env, true, true, true, false, UpdateErr{typ: Env, k: "flagbool", slug: "already set from flags"}},
		{"x-flagbool", Env, false, false, false, false, SettingNotFoundErr{name: "x-flagbool"}},
		{"flagbool", Flag, false, false, false, true, nil},
		{"flagbool", Flag, true, false, false, true, nil},
		{"flagbool", Flag, false, true, false, true, nil},
		// 580
		{"flagbool", Flag, true, true, false, true, nil},
		{"flagbool", Flag, false, false, true, false, UpdateErr{typ: Flag, k: "flagbool", slug: "already set from flags"}},
		{"flagbool", Flag, false, true, true, false, UpdateErr{typ: Flag, k: "flagbool", slug: "already set from flags"}},
		{"flagbool", Flag, true, false, true, false, UpdateErr{typ: Flag, k: "flagbool", slug: "already set from flags"}},
		{"flagbool", Flag, true, true, true, false, UpdateErr{typ: Flag, k: "flagbool", slug: "already set from flags"}},
		// 585
		{"x-flagbool", Flag, false, false, false, false, SettingNotFoundErr{name: "x-flagbool"}},
		{"flagint", Basic, false, false, false, false, BasicUpdateErr{typ: Flag, k: "flagint"}},
		{"flagint", Basic, true, false, false, false, BasicUpdateErr{typ: Flag, k: "flagint"}},
		{"flagint", Basic, true, false, false, false, BasicUpdateErr{typ: Flag, k: "flagint"}},
		{"flagint", Basic, false, true, false, false, BasicUpdateErr{typ: Flag, k: "flagint"}},
		// 590
		{"flagint", Basic, true, true, false, false, BasicUpdateErr{typ: Flag, k: "flagint"}},
		{"flagint", Basic, false, false, true, false, BasicUpdateErr{typ: Flag, k: "flagint"}},
		{"flagint", Basic, false, true, true, false, BasicUpdateErr{typ: Flag, k: "flagint"}},
		{"flagint", Basic, true, false, true, false, BasicUpdateErr{typ: Flag, k: "flagint"}},
		{"flagint", Basic, true, true, true, false, BasicUpdateErr{typ: Flag, k: "flagint"}},
		// 595
		{"x-flagint", Basic, false, false, false, false, SettingNotFoundErr{name: "x-flagint"}},
		{"flagint", Core, false, false, false, false, UpdateErr{typ: Core, k: "flagint", slug: "invalid update type"}},
		{"flagint", Core, true, false, false, false, UpdateErr{typ: Core, k: "flagint", slug: "invalid update type"}},
		{"flagint", Core, false, true, false, false, UpdateErr{typ: Core, k: "flagint", slug: "invalid update type"}},
		{"flagint", Core, true, true, false, false, UpdateErr{typ: Core, k: "flagint", slug: "invalid update type"}},
		// 600
		{"flagint", Core, false, false, true, false, UpdateErr{typ: Core, k: "flagint", slug: "invalid update type"}},
		{"flagint", Core, false, true, true, false, UpdateErr{typ: Core, k: "flagint", slug: "invalid update type"}},
		{"flagint", Core, true, false, true, false, UpdateErr{typ: Core, k: "flagint", slug: "invalid update type"}},
		{"flagint", Core, true, true, true, false, UpdateErr{typ: Core, k: "flagint", slug: "invalid update type"}},
		{"x-flagint", Core, false, false, false, false, SettingNotFoundErr{name: "x-flagint"}},
		// 605
		{"flagint", ConfFileVar, false, false, false, true, nil},
		{"flagint", ConfFileVar, true, false, false, false, UpdateErr{typ: ConfFileVar, k: "flagint", slug: "already set from the configuration file"}},
		{"flagint", ConfFileVar, false, true, false, false, UpdateErr{typ: ConfFileVar, k: "flagint", slug: "already set from env vars"}},
		{"flagint", ConfFileVar, true, true, false, false, UpdateErr{typ: ConfFileVar, k: "flagint", slug: "already set from env vars"}},
		{"flagint", ConfFileVar, false, false, true, false, UpdateErr{typ: ConfFileVar, k: "flagint", slug: "already set from flags"}},
		// 610
		{"flagint", ConfFileVar, false, true, true, false, UpdateErr{typ: ConfFileVar, k: "flagint", slug: "already set from flags"}},
		{"flagint", ConfFileVar, true, false, true, false, UpdateErr{typ: ConfFileVar, k: "flagint", slug: "already set from flags"}},
		{"flagint", ConfFileVar, true, true, true, false, UpdateErr{typ: Flag, k: "flagint", slug: "already set from flags"}},
		{"x-flagint", ConfFileVar, false, false, false, false, SettingNotFoundErr{name: "x-flagint"}},
		{"flagint", Env, false, false, false, true, nil},
		// 615
		{"flagint", Env, true, false, false, true, nil},
		{"flagint", Env, false, true, false, false, UpdateErr{typ: Env, k: "flagint", slug: "already set from env vars"}},
		{"flagint", Env, true, true, false, false, UpdateErr{typ: Env, k: "flagint", slug: "already set from env vars"}},
		{"flagint", Env, false, false, true, false, UpdateErr{typ: Env, k: "flagint", slug: "already set from flags"}},
		{"flagint", Env, false, true, true, false, UpdateErr{typ: Env, k: "flagint", slug: "already set from flags"}},
		// 620
		{"flagint", Env, true, false, true, false, UpdateErr{typ: Env, k: "flagint", slug: "already set from flags"}},
		{"flagint", Env, true, true, true, false, UpdateErr{typ: Env, k: "flagint", slug: "already set from flags"}},
		{"x-flagint", Env, false, false, false, false, SettingNotFoundErr{name: "x-flagint"}},
		{"flagint", Flag, false, false, false, true, nil},
		{"flagint", Flag, true, false, false, true, nil},
		// 625
		{"flagint", Flag, false, true, false, true, nil},
		{"flagint", Flag, true, true, false, true, nil},
		{"flagint", Flag, false, false, true, false, UpdateErr{typ: Flag, k: "flagint", slug: "already set from flags"}},
		{"flagint", Flag, false, true, true, false, UpdateErr{typ: Flag, k: "flagint", slug: "already set from flags"}},
		{"flagint", Flag, true, false, true, false, UpdateErr{typ: Flag, k: "flagint", slug: "already set from flags"}},
		// 630
		{"flagint", Flag, true, true, true, false, UpdateErr{typ: Flag, k: "flagint", slug: "already set from flags"}},
		{"x-flagint", Flag, false, false, false, false, SettingNotFoundErr{name: "x-flagint"}},
		{"flagint64", Basic, false, false, false, false, BasicUpdateErr{typ: Flag, k: "flagint64"}},
		{"flagint64", Basic, true, false, false, false, BasicUpdateErr{typ: Flag, k: "flagint64"}},
		{"flagint64", Basic, true, false, false, false, BasicUpdateErr{typ: Flag, k: "flagint64"}},
		// 635
		{"flagint64", Basic, false, true, false, false, BasicUpdateErr{typ: Flag, k: "flagint64"}},
		{"flagint64", Basic, true, true, false, false, BasicUpdateErr{typ: Flag, k: "flagint64"}},
		{"flagint64", Basic, false, false, true, false, BasicUpdateErr{typ: Flag, k: "flagint64"}},
		{"flagint64", Basic, false, true, true, false, BasicUpdateErr{typ: Flag, k: "flagint64"}},
		{"flagint64", Basic, true, false, true, false, BasicUpdateErr{typ: Flag, k: "flagint64"}},
		// 640
		{"flagint64", Basic, true, true, true, false, BasicUpdateErr{typ: Flag, k: "flagint64"}},
		{"x-flagint64", Basic, false, false, false, false, SettingNotFoundErr{name: "x-flagint64"}},
		{"flagint64", Core, false, false, false, false, UpdateErr{typ: Core, k: "flagint64", slug: "invalid update type"}},
		{"flagint64", Core, true, false, false, false, UpdateErr{typ: Core, k: "flagint64", slug: "invalid update type"}},
		{"flagint64", Core, false, true, false, false, UpdateErr{typ: Core, k: "flagint64", slug: "invalid update type"}},
		// 645
		{"flagint64", Core, true, true, false, false, UpdateErr{typ: Core, k: "flagint64", slug: "invalid update type"}},
		{"flagint64", Core, false, false, true, false, UpdateErr{typ: Core, k: "flagint64", slug: "invalid update type"}},
		{"flagint64", Core, false, true, true, false, UpdateErr{typ: Core, k: "flagint64", slug: "invalid update type"}},
		{"flagint64", Core, true, false, true, false, UpdateErr{typ: Core, k: "flagint64", slug: "invalid update type"}},
		{"flagint64", Core, true, true, true, false, UpdateErr{typ: Core, k: "flagint64", slug: "invalid update type"}},
		// 650
		{"x-flagint64", Core, false, false, false, false, SettingNotFoundErr{name: "x-flagint64"}},
		{"flagint64", ConfFileVar, false, false, false, true, nil},
		{"flagint64", ConfFileVar, true, false, false, false, UpdateErr{typ: ConfFileVar, k: "flagint64", slug: "already set from the configuration file"}},
		{"flagint64", ConfFileVar, false, true, false, false, UpdateErr{typ: ConfFileVar, k: "flagint64", slug: "already set from env vars"}},
		{"flagint64", ConfFileVar, true, true, false, false, UpdateErr{typ: ConfFileVar, k: "flagint64", slug: "already set from env vars"}},
		// 655
		{"flagint64", ConfFileVar, false, false, true, false, UpdateErr{typ: ConfFileVar, k: "flagint64", slug: "already set from flags"}},
		{"flagint64", ConfFileVar, false, true, true, false, UpdateErr{typ: ConfFileVar, k: "flagint64", slug: "already set from flags"}},
		{"flagint64", ConfFileVar, true, false, true, false, UpdateErr{typ: ConfFileVar, k: "flagint64", slug: "already set from flags"}},
		{"flagint64", ConfFileVar, true, true, true, false, UpdateErr{typ: ConfFileVar, k: "flagint64", slug: "already set from flags"}},
		{"x-flagint64", ConfFileVar, false, false, false, false, SettingNotFoundErr{name: "x-flagint64"}},
		// 660
		{"flagint64", Env, false, false, false, true, nil},
		{"flagint64", Env, true, false, false, true, nil},
		{"flagint64", Env, false, true, false, false, UpdateErr{typ: Env, k: "flagint64", slug: "already set from env vars"}},
		{"flagint64", Env, true, true, false, false, UpdateErr{typ: Env, k: "flagint64", slug: "already set from env vars"}},
		{"flagint64", Env, false, false, true, false, UpdateErr{typ: Env, k: "flagint64", slug: "already set from flags"}},
		// 665
		{"flagint64", Env, false, true, true, false, UpdateErr{typ: Env, k: "flagint64", slug: "already set from flags"}},
		{"flagint64", Env, true, false, true, false, UpdateErr{typ: Env, k: "flagint64", slug: "already set from flags"}},
		{"flagint64", Env, true, true, true, false, UpdateErr{typ: Env, k: "flagint64", slug: "already set from flags"}},
		{"x-flagint64", Env, false, false, false, false, SettingNotFoundErr{name: "x-flagint64"}},
		{"flagint64", Flag, false, false, false, true, nil},
		// 670
		{"flagint64", Flag, true, false, false, true, nil},
		{"flagint64", Flag, false, true, false, true, nil},
		{"flagint64", Flag, true, true, false, true, nil},
		{"flagint64", Flag, false, false, true, false, UpdateErr{typ: Flag, k: "flagint64", slug: "already set from flags"}},
		{"flagint64", Flag, false, true, true, false, UpdateErr{typ: Flag, k: "flagint64", slug: "already set from flags"}},
		// 675
		{"flagint64", Flag, true, false, true, false, UpdateErr{typ: Flag, k: "flagint64", slug: "already set from flags"}},
		{"flagint64", Flag, true, true, true, false, UpdateErr{typ: Flag, k: "flagint64", slug: "already set from flags"}},
		{"x-flagint64", Flag, false, false, false, false, SettingNotFoundErr{name: "x-flagint64"}},
		{"flagstring", Basic, false, false, false, false, BasicUpdateErr{typ: Flag, k: "flagstring"}},
		{"flagstring", Basic, true, false, false, false, BasicUpdateErr{typ: Flag, k: "flagstring"}},
		// 680
		{"flagstring", Basic, true, false, false, false, BasicUpdateErr{typ: Flag, k: "flagstring"}},
		{"flagstring", Basic, false, true, false, false, BasicUpdateErr{typ: Flag, k: "flagstring"}},
		{"flagstring", Basic, true, true, false, false, BasicUpdateErr{typ: Flag, k: "flagstring"}},
		{"flagstring", Basic, false, false, true, false, BasicUpdateErr{typ: Flag, k: "flagstring"}},
		{"flagstring", Basic, false, true, true, false, BasicUpdateErr{typ: Flag, k: "flagstring"}},
		// 685
		{"flagstring", Basic, true, false, true, false, BasicUpdateErr{typ: Flag, k: "flagstring"}},
		{"flagstring", Basic, true, true, true, false, BasicUpdateErr{typ: Flag, k: "flagstring"}},
		{"x-flagstring", Basic, false, false, false, false, SettingNotFoundErr{name: "x-flagstring"}},
		{"flagstring", Core, false, false, false, false, UpdateErr{typ: Core, k: "flagstring", slug: "invalid update type"}},
		{"flagstring", Core, true, false, false, false, UpdateErr{typ: Core, k: "flagstring", slug: "invalid update type"}},
		// 690
		{"flagstring", Core, false, true, false, false, UpdateErr{typ: Core, k: "flagstring", slug: "invalid update type"}},
		{"flagstring", Core, true, true, false, false, UpdateErr{typ: Core, k: "flagstring", slug: "invalid update type"}},
		{"flagstring", Core, false, false, true, false, UpdateErr{typ: Core, k: "flagstring", slug: "invalid update type"}},
		{"flagstring", Core, false, true, true, false, UpdateErr{typ: Core, k: "flagstring", slug: "invalid update type"}},
		{"flagstring", Core, true, false, true, false, UpdateErr{typ: Core, k: "flagstring", slug: "invalid update type"}},
		// 695
		{"flagstring", Core, true, true, true, false, UpdateErr{typ: Core, k: "flagstring", slug: "invalid update type"}},
		{"x-flagstring", Core, false, false, false, false, SettingNotFoundErr{name: "x-flagstring"}},
		{"flagstring", ConfFileVar, false, false, false, true, nil},
		{"flagstring", ConfFileVar, true, false, false, false, UpdateErr{typ: ConfFileVar, k: "flagstring", slug: "already set from the configuration file"}},
		{"flagstring", ConfFileVar, false, true, false, false, UpdateErr{typ: ConfFileVar, k: "flagstring", slug: "already set from env vars"}},
		// 700
		{"flagstring", ConfFileVar, true, true, false, false, UpdateErr{typ: ConfFileVar, k: "flagstring", slug: "already set from env vars"}},
		{"flagstring", ConfFileVar, false, false, true, false, UpdateErr{typ: ConfFileVar, k: "flagstring", slug: "already set from flags"}},
		{"flagstring", ConfFileVar, false, true, true, false, UpdateErr{typ: ConfFileVar, k: "flagstring", slug: "already set from flags"}},
		{"flagstring", ConfFileVar, true, false, true, false, UpdateErr{typ: ConfFileVar, k: "flagstring", slug: "already set from flags"}},
		{"flagstring", ConfFileVar, true, true, true, false, UpdateErr{typ: ConfFileVar, k: "flagstring", slug: "already set from flags"}},
		// 705
		{"x-flagstring", ConfFileVar, false, false, false, false, SettingNotFoundErr{name: "x-flagstring"}},
		{"flagstring", Env, false, false, false, true, nil},
		{"flagstring", Env, true, false, false, true, nil},
		{"flagstring", Env, false, true, false, false, UpdateErr{typ: Env, k: "flagstring", slug: "already set from env vars"}},
		{"flagstring", Env, true, true, false, false, UpdateErr{typ: Env, k: "flagstring", slug: "already set from env vars"}},
		// 710
		{"flagstring", Env, false, false, true, false, UpdateErr{typ: Env, k: "flagstring", slug: "already set from flags"}},
		{"flagstring", Env, false, true, true, false, UpdateErr{typ: Env, k: "flagstring", slug: "already set from flags"}},
		{"flagstring", Env, true, false, true, false, UpdateErr{typ: Env, k: "flagstring", slug: "already set from flags"}},
		{"flagstring", Env, true, true, true, false, UpdateErr{typ: Env, k: "flagstring", slug: "already set from flags"}},
		{"x-flagstring", Env, false, false, false, false, SettingNotFoundErr{name: "x-flagstring"}},
		// 715
		{"flagstring", Flag, false, false, false, true, nil},
		{"flagstring", Flag, true, false, false, true, nil},
		{"flagstring", Flag, false, true, false, true, nil},
		{"flagstring", Flag, true, true, false, true, nil},
		{"flagstring", Flag, false, false, true, false, UpdateErr{typ: Flag, k: "flagstring", slug: "already set from flags"}},
		// 720
		{"flagstring", Flag, false, true, true, false, UpdateErr{typ: Flag, k: "flagstring", slug: "already set from flags"}},
		{"flagstring", Flag, true, false, true, false, UpdateErr{typ: Flag, k: "flagstring", slug: "already set from flags"}},
		{"flagstring", Flag, true, true, true, false, UpdateErr{typ: Flag, k: "flagstring", slug: "already set from flags"}},
		{"x-flagstring", Flag, false, false, false, false, SettingNotFoundErr{name: "x-flagstring"}},
		{"bool", Basic, false, false, false, true, nil},

		// 725
		{"bool", Basic, true, false, false, true, nil},
		{"bool", Basic, false, true, false, true, nil},
		{"bool", Basic, true, true, false, true, nil},
		{"bool", Basic, false, false, true, true, nil},
		{"bool", Basic, false, true, true, true, nil},
		// 730
		{"bool", Basic, true, false, true, true, nil},
		{"bool", Basic, true, true, true, true, nil},
		{"x-bool", Basic, false, false, false, false, SettingNotFoundErr{name: "x-bool"}},
		{"bool", Core, false, false, false, false, UpdateErr{typ: Core, k: "bool", slug: "invalid update type"}},
		{"bool", Core, true, false, false, false, UpdateErr{typ: Core, k: "bool", slug: "invalid update type"}},
		// 735
		{"bool", Core, false, true, false, false, UpdateErr{typ: Core, k: "bool", slug: "invalid update type"}},
		{"bool", Core, true, true, false, false, UpdateErr{typ: Core, k: "bool", slug: "invalid update type"}},
		{"bool", Core, false, false, true, false, UpdateErr{typ: Core, k: "bool", slug: "invalid update type"}},
		{"bool", Core, false, true, true, false, UpdateErr{typ: Core, k: "bool", slug: "invalid update type"}},
		{"bool", Core, true, false, true, false, UpdateErr{typ: Core, k: "bool", slug: "invalid update type"}},
		// 740
		{"bool", Core, true, true, true, false, UpdateErr{typ: Core, k: "bool", slug: "invalid update type"}},
		{"x-bool", Core, false, false, false, false, SettingNotFoundErr{name: "x-bool"}},
		{"bool", ConfFileVar, false, false, false, false, UpdateErr{typ: ConfFileVar, k: "bool", slug: fmt.Sprintf("is not a %s", ConfFileVar)}},
		{"bool", ConfFileVar, true, false, false, false, UpdateErr{typ: ConfFileVar, k: "bool", slug: fmt.Sprintf("is not a %s", ConfFileVar)}},
		{"bool", ConfFileVar, false, true, false, false, UpdateErr{typ: ConfFileVar, k: "bool", slug: fmt.Sprintf("is not a %s", ConfFileVar)}},
		// 745
		{"bool", ConfFileVar, true, true, false, false, UpdateErr{typ: ConfFileVar, k: "bool", slug: fmt.Sprintf("is not a %s", ConfFileVar)}},
		{"bool", ConfFileVar, false, true, true, false, UpdateErr{typ: ConfFileVar, k: "bool", slug: fmt.Sprintf("is not a %s", ConfFileVar)}},
		{"bool", ConfFileVar, false, false, true, false, UpdateErr{typ: ConfFileVar, k: "bool", slug: fmt.Sprintf("is not a %s", ConfFileVar)}},
		{"bool", ConfFileVar, true, false, true, false, UpdateErr{typ: ConfFileVar, k: "bool", slug: fmt.Sprintf("is not a %s", ConfFileVar)}},
		{"bool", ConfFileVar, true, true, true, false, UpdateErr{typ: ConfFileVar, k: "bool", slug: fmt.Sprintf("is not a %s", ConfFileVar)}},
		// 750
		{"x-bool", ConfFileVar, false, false, false, false, SettingNotFoundErr{name: "x-bool"}},
		{"bool", Env, false, false, false, false, UpdateErr{typ: ConfFileVar, k: "bool", slug: fmt.Sprintf("is not an %s", Env)}},
		{"bool", Env, true, false, false, false, UpdateErr{typ: ConfFileVar, k: "bool", slug: fmt.Sprintf("is not an %s", Env)}},
		{"bool", Env, false, true, false, false, UpdateErr{typ: ConfFileVar, k: "bool", slug: fmt.Sprintf("is not an %s", Env)}},
		{"bool", Env, true, true, false, false, UpdateErr{typ: ConfFileVar, k: "bool", slug: fmt.Sprintf("is not an %s", Env)}},
		// 755
		{"bool", Env, false, false, true, false, UpdateErr{typ: ConfFileVar, k: "bool", slug: fmt.Sprintf("is not an %s", Env)}},
		{"bool", Env, false, true, true, false, UpdateErr{typ: ConfFileVar, k: "bool", slug: fmt.Sprintf("is not an %s", Env)}},
		{"bool", Env, true, false, true, false, UpdateErr{typ: ConfFileVar, k: "bool", slug: fmt.Sprintf("is not an %s", Env)}},
		{"bool", Env, true, true, true, false, UpdateErr{typ: ConfFileVar, k: "bool", slug: fmt.Sprintf("is not an %s", Env)}},
		{"x-bool", Env, false, false, false, false, SettingNotFoundErr{name: "x-bool"}},
		// 760
		{"bool", Flag, false, false, false, false, UpdateErr{typ: ConfFileVar, k: "bool", slug: fmt.Sprintf("is not a %s", Flag)}},
		{"bool", Flag, true, false, false, false, UpdateErr{typ: ConfFileVar, k: "bool", slug: fmt.Sprintf("is not a %s", Flag)}},
		{"bool", Flag, false, true, false, false, UpdateErr{typ: ConfFileVar, k: "bool", slug: fmt.Sprintf("is not a %s", Flag)}},
		{"bool", Flag, true, true, false, false, UpdateErr{typ: ConfFileVar, k: "bool", slug: fmt.Sprintf("is not a %s", Flag)}},
		{"bool", Flag, false, false, true, false, UpdateErr{typ: ConfFileVar, k: "bool", slug: fmt.Sprintf("is not a %s", Flag)}},
		// 765
		{"bool", Flag, false, true, true, false, UpdateErr{typ: ConfFileVar, k: "bool", slug: fmt.Sprintf("is not a %s", Flag)}},
		{"bool", Flag, true, false, true, false, UpdateErr{typ: ConfFileVar, k: "bool", slug: fmt.Sprintf("is not a %s", Flag)}},
		{"bool", Flag, true, true, true, false, UpdateErr{typ: ConfFileVar, k: "bool", slug: fmt.Sprintf("is not a %s", Flag)}},
		{"x-bool", Flag, false, false, false, false, SettingNotFoundErr{name: "x-bool"}},
		{"int", Basic, false, false, false, true, nil},
		// 770
		{"int", Basic, true, false, false, true, nil},
		{"int", Basic, false, true, false, true, nil},
		{"int", Basic, true, true, false, true, nil},
		{"int", Basic, false, false, true, true, nil},
		{"int", Basic, false, true, true, true, nil},
		// 775
		{"int", Basic, true, false, true, true, nil},
		{"int", Basic, true, true, true, true, nil},
		{"x-int", Basic, false, false, false, false, SettingNotFoundErr{name: "x-int"}},
		{"int", Core, false, false, false, false, UpdateErr{typ: Core, k: "int", slug: "invalid update type"}},
		{"int", Core, true, false, false, false, UpdateErr{typ: Core, k: "int", slug: "invalid update type"}},
		// 780
		{"int", Core, false, true, false, false, UpdateErr{typ: Core, k: "int", slug: "invalid update type"}},
		{"int", Core, true, true, false, false, UpdateErr{typ: Core, k: "int", slug: "invalid update type"}},
		{"int", Core, false, false, true, false, UpdateErr{typ: Core, k: "int", slug: "invalid update type"}},
		{"int", Core, false, true, true, false, UpdateErr{typ: Core, k: "int", slug: "invalid update type"}},
		{"int", Core, true, false, true, false, UpdateErr{typ: Core, k: "int", slug: "invalid update type"}},
		// 785
		{"int", Core, true, true, true, false, UpdateErr{typ: Core, k: "int", slug: "invalid update type"}},
		{"x-int", Core, false, false, false, false, SettingNotFoundErr{name: "x-int"}},
		{"int", ConfFileVar, false, false, false, false, UpdateErr{typ: ConfFileVar, k: "int", slug: fmt.Sprintf("is not a %s", ConfFileVar)}},
		{"int", ConfFileVar, true, false, false, false, UpdateErr{typ: ConfFileVar, k: "int", slug: fmt.Sprintf("is not a %s", ConfFileVar)}},
		{"int", ConfFileVar, false, true, false, false, UpdateErr{typ: ConfFileVar, k: "int", slug: fmt.Sprintf("is not a %s", ConfFileVar)}},
		// 790
		{"int", ConfFileVar, true, true, false, false, UpdateErr{typ: ConfFileVar, k: "int", slug: fmt.Sprintf("is not a %s", ConfFileVar)}},
		{"int", ConfFileVar, false, false, true, false, UpdateErr{typ: ConfFileVar, k: "int", slug: fmt.Sprintf("is not a %s", ConfFileVar)}},
		{"int", ConfFileVar, false, true, true, false, UpdateErr{typ: ConfFileVar, k: "int", slug: fmt.Sprintf("is not a %s", ConfFileVar)}},
		{"int", ConfFileVar, true, false, true, false, UpdateErr{typ: ConfFileVar, k: "int", slug: fmt.Sprintf("is not a %s", ConfFileVar)}},
		// 795
		{"int", ConfFileVar, true, true, true, false, UpdateErr{typ: ConfFileVar, k: "int", slug: fmt.Sprintf("is not a %s", ConfFileVar)}},
		{"x-int", ConfFileVar, false, false, false, false, SettingNotFoundErr{name: "x-int"}},
		{"int", Env, false, false, false, false, UpdateErr{typ: ConfFileVar, k: "int", slug: fmt.Sprintf("is not an %s", Env)}},
		{"int", Env, true, false, false, false, UpdateErr{typ: ConfFileVar, k: "int", slug: fmt.Sprintf("is not an %s", Env)}},
		{"int", Env, false, true, false, false, UpdateErr{typ: ConfFileVar, k: "int", slug: fmt.Sprintf("is not an %s", Env)}},
		// 800
		{"int", Env, true, true, false, false, UpdateErr{typ: ConfFileVar, k: "int", slug: fmt.Sprintf("is not an %s", Env)}},
		{"int", Env, false, false, true, false, UpdateErr{typ: ConfFileVar, k: "int", slug: fmt.Sprintf("is not an %s", Env)}},
		{"int", Env, false, true, true, false, UpdateErr{typ: ConfFileVar, k: "int", slug: fmt.Sprintf("is not an %s", Env)}},
		{"int", Env, true, false, true, false, UpdateErr{typ: ConfFileVar, k: "int", slug: fmt.Sprintf("is not an %s", Env)}},
		{"int", Env, true, true, true, false, UpdateErr{typ: ConfFileVar, k: "int", slug: fmt.Sprintf("is not an %s", Env)}},
		{"x-int", Env, false, false, false, false, SettingNotFoundErr{name: "x-int"}},
		// 805
		{"int", Flag, false, false, false, false, UpdateErr{typ: ConfFileVar, k: "int", slug: fmt.Sprintf("is not a %s", Flag)}},
		{"int", Flag, true, false, false, false, UpdateErr{typ: ConfFileVar, k: "int", slug: fmt.Sprintf("is not a %s", Flag)}},
		{"int", Flag, false, true, false, false, UpdateErr{typ: ConfFileVar, k: "int", slug: fmt.Sprintf("is not a %s", Flag)}},
		{"int", Flag, true, true, false, false, UpdateErr{typ: ConfFileVar, k: "int", slug: fmt.Sprintf("is not a %s", Flag)}},
		{"int", Flag, false, false, true, false, UpdateErr{typ: ConfFileVar, k: "int", slug: fmt.Sprintf("is not a %s", Flag)}},
		// 810
		{"int", Flag, false, true, true, false, UpdateErr{typ: ConfFileVar, k: "int", slug: fmt.Sprintf("is not a %s", Flag)}},
		{"int", Flag, true, false, true, false, UpdateErr{typ: ConfFileVar, k: "int", slug: fmt.Sprintf("is not a %s", Flag)}},
		{"int", Flag, true, true, true, false, UpdateErr{typ: ConfFileVar, k: "int", slug: fmt.Sprintf("is not a %s", Flag)}},
		{"x-int", Flag, false, false, false, false, SettingNotFoundErr{name: "x-int"}},
		{"int64", Basic, false, false, false, true, nil},
		// 815
		{"int64", Basic, true, false, false, true, nil},
		{"int64", Basic, false, true, false, true, nil},
		{"int64", Basic, true, true, false, true, nil},
		{"int64", Basic, false, false, true, true, nil},
		{"int64", Basic, false, true, true, true, nil},
		// 820
		{"int64", Basic, true, false, true, true, nil},
		{"int64", Basic, true, true, true, true, nil},
		{"x-int64", Basic, false, false, false, false, SettingNotFoundErr{name: "x-int64"}},
		{"int64", Core, false, false, false, false, UpdateErr{typ: Core, k: "int64", slug: "invalid update type"}},
		{"int64", Core, true, false, false, false, UpdateErr{typ: Core, k: "int64", slug: "invalid update type"}},
		// 825
		{"int64", Core, false, true, false, false, UpdateErr{typ: Core, k: "int64", slug: "invalid update type"}},
		{"int64", Core, true, true, false, false, UpdateErr{typ: Core, k: "int64", slug: "invalid update type"}},
		{"int64", Core, false, false, true, false, UpdateErr{typ: Core, k: "int64", slug: "invalid update type"}},
		{"int64", Core, false, true, true, false, UpdateErr{typ: Core, k: "int64", slug: "invalid update type"}},
		{"int64", Core, true, false, true, false, UpdateErr{typ: Core, k: "int64", slug: "invalid update type"}},
		// 830
		{"int64", Core, true, true, true, false, UpdateErr{typ: Core, k: "int64", slug: "invalid update type"}},
		{"x-int64", Core, false, false, false, false, SettingNotFoundErr{name: "x-int64"}},
		{"int64", ConfFileVar, false, false, false, false, UpdateErr{typ: ConfFileVar, k: "int64", slug: fmt.Sprintf("is not a %s", ConfFileVar)}},
		{"int64", ConfFileVar, true, false, false, false, UpdateErr{typ: ConfFileVar, k: "int64", slug: fmt.Sprintf("is not a %s", ConfFileVar)}},
		{"int64", ConfFileVar, false, true, false, false, UpdateErr{typ: ConfFileVar, k: "int64", slug: fmt.Sprintf("is not a %s", ConfFileVar)}},
		// 835
		{"int64", ConfFileVar, true, true, false, false, UpdateErr{typ: ConfFileVar, k: "int64", slug: fmt.Sprintf("is not a %s", ConfFileVar)}},
		{"int64", ConfFileVar, false, false, true, false, UpdateErr{typ: ConfFileVar, k: "int64", slug: fmt.Sprintf("is not a %s", ConfFileVar)}},
		{"int64", ConfFileVar, false, true, true, false, UpdateErr{typ: ConfFileVar, k: "int64", slug: fmt.Sprintf("is not a %s", ConfFileVar)}},
		{"int64", ConfFileVar, true, false, true, false, UpdateErr{typ: ConfFileVar, k: "int64", slug: fmt.Sprintf("is not a %s", ConfFileVar)}},
		{"int64", ConfFileVar, true, true, true, false, UpdateErr{typ: ConfFileVar, k: "int64", slug: fmt.Sprintf("is not a %s", ConfFileVar)}},
		// 840
		{"x-int64", ConfFileVar, false, false, false, false, SettingNotFoundErr{name: "x-int64"}},
		{"int64", Env, false, false, false, false, UpdateErr{typ: ConfFileVar, k: "int64", slug: fmt.Sprintf("is not an %s", Env)}},
		{"int64", Env, true, false, false, false, UpdateErr{typ: ConfFileVar, k: "int64", slug: fmt.Sprintf("is not an %s", Env)}},
		{"int64", Env, false, true, false, false, UpdateErr{typ: ConfFileVar, k: "int64", slug: fmt.Sprintf("is not an %s", Env)}},
		{"int64", Env, true, true, false, false, UpdateErr{typ: ConfFileVar, k: "int64", slug: fmt.Sprintf("is not an %s", Env)}},
		// 845
		{"int64", Env, false, false, true, false, UpdateErr{typ: ConfFileVar, k: "int64", slug: fmt.Sprintf("is not an %s", Env)}},
		{"int64", Env, false, true, true, false, UpdateErr{typ: ConfFileVar, k: "int64", slug: fmt.Sprintf("is not an %s", Env)}},
		{"int64", Env, true, false, true, false, UpdateErr{typ: ConfFileVar, k: "int64", slug: fmt.Sprintf("is not an %s", Env)}},
		{"int64", Env, true, true, true, false, UpdateErr{typ: ConfFileVar, k: "int64", slug: fmt.Sprintf("is not an %s", Env)}},
		{"x-int64", Env, false, false, false, false, SettingNotFoundErr{name: "x-int64"}},
		// 850
		{"int64", Flag, false, false, false, false, UpdateErr{typ: ConfFileVar, k: "int64", slug: fmt.Sprintf("is not a %s", Flag)}},
		{"int64", Flag, true, false, false, false, UpdateErr{typ: ConfFileVar, k: "int64", slug: fmt.Sprintf("is not a %s", Flag)}},
		{"int64", Flag, false, true, false, false, UpdateErr{typ: ConfFileVar, k: "int64", slug: fmt.Sprintf("is not a %s", Flag)}},
		{"int64", Flag, true, true, false, false, UpdateErr{typ: ConfFileVar, k: "int64", slug: fmt.Sprintf("is not a %s", Flag)}},
		{"int64", Flag, false, false, true, false, UpdateErr{typ: ConfFileVar, k: "int64", slug: fmt.Sprintf("is not a %s", Flag)}},
		// 855
		{"int64", Flag, false, true, true, false, UpdateErr{typ: ConfFileVar, k: "int64", slug: fmt.Sprintf("is not a %s", Flag)}},
		{"int64", Flag, true, false, true, false, UpdateErr{typ: ConfFileVar, k: "int64", slug: fmt.Sprintf("is not a %s", Flag)}},
		{"int64", Flag, true, true, true, false, UpdateErr{typ: ConfFileVar, k: "int64", slug: fmt.Sprintf("is not a %s", Flag)}},
		{"x-int64", Flag, false, false, false, false, SettingNotFoundErr{name: "x-int64"}},
		{"string", Basic, false, false, false, true, nil},
		// 860
		{"string", Basic, true, false, false, true, nil},
		{"string", Basic, false, true, false, true, nil},
		{"string", Basic, true, true, false, true, nil},
		{"string", Basic, false, false, true, true, nil},
		{"string", Basic, false, true, true, true, nil},
		// 865
		{"string", Basic, true, false, true, true, nil},
		{"string", Basic, true, true, true, true, nil},
		{"x-string", Basic, false, false, false, false, SettingNotFoundErr{name: "x-string"}},
		{"string", Core, false, false, false, false, UpdateErr{typ: Core, k: "string", slug: "invalid update type"}},
		{"string", Core, true, false, false, false, UpdateErr{typ: Core, k: "string", slug: "invalid update type"}},
		// 870
		{"string", Core, false, true, false, false, UpdateErr{typ: Core, k: "string", slug: "invalid update type"}},
		{"string", Core, true, true, false, false, UpdateErr{typ: Core, k: "string", slug: "invalid update type"}},
		{"string", Core, false, false, true, false, UpdateErr{typ: Core, k: "string", slug: "invalid update type"}},
		{"string", Core, false, true, true, false, UpdateErr{typ: Core, k: "string", slug: "invalid update type"}},
		{"string", Core, true, false, true, false, UpdateErr{typ: Core, k: "string", slug: "invalid update type"}},
		// 875
		{"string", Core, true, true, true, false, UpdateErr{typ: Core, k: "string", slug: "invalid update type"}},
		{"x-string", Core, false, false, false, false, SettingNotFoundErr{name: "x-string"}},
		{"string", ConfFileVar, false, false, false, false, UpdateErr{typ: ConfFileVar, k: "string", slug: fmt.Sprintf("is not a %s", ConfFileVar)}},
		{"string", ConfFileVar, true, false, false, false, UpdateErr{typ: ConfFileVar, k: "string", slug: fmt.Sprintf("is not a %s", ConfFileVar)}},
		{"string", ConfFileVar, false, true, false, false, UpdateErr{typ: ConfFileVar, k: "string", slug: fmt.Sprintf("is not a %s", ConfFileVar)}},
		// 880
		{"string", ConfFileVar, true, true, false, false, UpdateErr{typ: ConfFileVar, k: "string", slug: fmt.Sprintf("is not a %s", ConfFileVar)}},
		{"string", ConfFileVar, false, false, true, false, UpdateErr{typ: ConfFileVar, k: "string", slug: fmt.Sprintf("is not a %s", ConfFileVar)}},
		{"string", ConfFileVar, false, true, true, false, UpdateErr{typ: ConfFileVar, k: "string", slug: fmt.Sprintf("is not a %s", ConfFileVar)}},
		{"string", ConfFileVar, true, false, true, false, UpdateErr{typ: ConfFileVar, k: "string", slug: fmt.Sprintf("is not a %s", ConfFileVar)}},
		{"string", ConfFileVar, true, true, true, false, UpdateErr{typ: ConfFileVar, k: "string", slug: fmt.Sprintf("is not a %s", ConfFileVar)}},
		// 885
		{"x-string", ConfFileVar, false, false, false, false, SettingNotFoundErr{name: "x-string"}},
		{"string", Env, false, false, false, false, UpdateErr{typ: ConfFileVar, k: "string", slug: fmt.Sprintf("is not an %s", Env)}},
		{"string", Env, true, false, false, false, UpdateErr{typ: ConfFileVar, k: "string", slug: fmt.Sprintf("is not an %s", Env)}},
		{"string", Env, false, true, false, false, UpdateErr{typ: ConfFileVar, k: "string", slug: fmt.Sprintf("is not an %s", Env)}},
		{"string", Env, true, true, false, false, UpdateErr{typ: ConfFileVar, k: "string", slug: fmt.Sprintf("is not an %s", Env)}},
		// 890
		{"string", Env, false, false, true, false, UpdateErr{typ: ConfFileVar, k: "string", slug: fmt.Sprintf("is not an %s", Env)}},
		{"string", Env, false, true, true, false, UpdateErr{typ: ConfFileVar, k: "string", slug: fmt.Sprintf("is not an %s", Env)}},
		{"string", Env, true, false, true, false, UpdateErr{typ: ConfFileVar, k: "string", slug: fmt.Sprintf("is not an %s", Env)}},
		{"string", Env, true, true, true, false, UpdateErr{typ: ConfFileVar, k: "string", slug: fmt.Sprintf("is not an %s", Env)}},
		{"x-string", Env, false, false, false, false, SettingNotFoundErr{name: "x-string"}},
		// 895
		{"string", Flag, false, false, false, false, UpdateErr{typ: ConfFileVar, k: "string", slug: fmt.Sprintf("is not a %s", Flag)}},
		{"string", Flag, true, false, false, false, UpdateErr{typ: ConfFileVar, k: "string", slug: fmt.Sprintf("is not a %s", Flag)}},
		{"string", Flag, false, true, false, false, UpdateErr{typ: ConfFileVar, k: "string", slug: fmt.Sprintf("is not a %s", Flag)}},
		{"string", Flag, true, true, false, false, UpdateErr{typ: ConfFileVar, k: "string", slug: fmt.Sprintf("is not a %s", Flag)}},
		{"string", Flag, false, false, true, false, UpdateErr{typ: ConfFileVar, k: "string", slug: fmt.Sprintf("is not a %s", Flag)}},
		// 900
		{"string", Flag, false, true, true, false, UpdateErr{typ: ConfFileVar, k: "string", slug: fmt.Sprintf("is not a %s", Flag)}},
		{"string", Flag, true, false, true, false, UpdateErr{typ: ConfFileVar, k: "string", slug: fmt.Sprintf("is not a %s", Flag)}},
		{"string", Flag, true, true, true, false, UpdateErr{typ: ConfFileVar, k: "string", slug: fmt.Sprintf("is not a %s", Flag)}},
		{"x-string", Flag, false, false, false, false, SettingNotFoundErr{name: "x-string"}},
	}
	appCfg := newTestSettings()
	for i, test := range tests {
		appCfg.flagsParsed = test.flagsParsed
		appCfg.confFileVarsSet = test.confFileVarsSet
		appCfg.envSet = test.envVarsSet
		appCfg.flagsParsed = test.flagsParsed
		b, err := appCfg.canUpdate(test.typ, test.name)
		if err != nil {
			if test.err == nil {
				t.Errorf("%d: %s: %s: got %s; wanted no error", i, test.name, test.typ, err)
			} else if err.Error() != test.err.Error() {
				t.Errorf("%d: %s: %s: got %s; want %s", i, test.name, test.typ, err, test.err)
			}
		} else {
			if test.err != nil {
				t.Errorf("%d: %s: %s: got no error; wanted %s", i, test.name, test.typ, test.err)
			}
		}

		if b != test.expected {
			t.Errorf("%d: %s:%s: got %v; want %v", i, test.name, test.typ, b, test.expected)
		}
	}
}
