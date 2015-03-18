package contour

import (
	"strconv"
	"testing"
)

func TestUpdateBools(t *testing.T) {
	bTests := []struct {
		name     string
		key      string
		value    bool
		expected bool
		err      string
	}{
		{"false", "corebool", false, false, "config[update]: \"corebool\" is not updateable"},
		{"true", "corebool", true, true, "config[update]: \"corebool\" is not updateable"},
		{"false", "flagbool", false, false, ""},
		{"true", "flagbool", true, true, ""},
		{"false", "cfgbool", false, false, ""},
		{"true", "cfgbool", true, true, ""},
		{"false", "bool", false, false, ""},
		{"true", "bool", true, true, ""},
	}

	cfg := testCfg
	cfg.name = "update"
	for _, test := range bTests {
		err := cfg.UpdateBoolE(test.key, test.value)
		if err != nil {
			if test.err == "" {
				t.Errorf("%s: unexpected error %q", test.name, err)
				goto cont
			}
			if test.err != err.Error() {
				t.Errorf("%s: expected %q got %q", test.name, test.err, err)
			}
		cont:
			continue
		}
		b, err := cfg.GetBoolE(test.key)
		if err != nil {
			if test.err == "" {
				t.Errorf("%s: unexpected error %q", test.name, err)
			} else {
				if test.err != err.Error() {
					t.Errorf("%s: expected %q got %q", test.name, test.err, err)
				}
			}
			continue
		}
		if b != test.expected {
			t.Errorf("%s: expected %t got %t", test.name, test.value, b)
		}
	}
}
func TestUpdateStrings(t *testing.T) {
	sTests := []struct {
		name  string
		key   string
		value string
		err   string
	}{
		{"false", "corestring", "false", "config[update]: \"corestring\" is not updateable"},
		{"true", "corestring", "t", "config[update]: \"corestring\" is not updateable"},
		{"unset", "", "", "config[update]: \"\" is not updateable"},
		{"false", "flagstring", "false", ""},
		{"true", "flagstring", "t", ""},
		{"false", "cfgstring", "false", ""},
		{"true", "cfgstring", "t", ""},
		{"false", "string", "false", ""},
		{"true", "string", "t", ""},
	}

	cfg := testCfg
	cfg.name = "update"

	for _, test := range sTests {
		err := cfg.UpdateStringE(test.key, test.value)
		if err != nil {
			if test.err == "" {
				t.Errorf("%s: unexpected error %q", test.name, err)
			} else {
				if test.err != err.Error() {
					t.Errorf("%s: expected %q got %q", test.name, test.err, err)
				}
			}
			continue
		}
		s, err := cfg.GetStringE(test.key)
		if err != nil {
			if test.err == "" {
				t.Errorf("%s: unexpected error %q", test.name, err)
			} else {
				if test.err != err.Error() {
					t.Errorf("%s: expected %q got %q", test.name, test.err, err)
				}
			}
			continue
		}
		if s != test.value {
			t.Errorf("%s: expected %t got %s", test.name, test.value, s)
		}
	}
}

func TestUpdateInts(t *testing.T) {
	iTests := []struct {
		name  string
		key   string
		value int
		err   string
	}{
		{"42", "coreint", 42, "config[update]: \"coreint\" is not updateable"},
		{"unset", "", 0, "config[update]: \"\" is not updateable"},
		{"42", "flagint", 42, ""},
		{"42", "cfgint", 42, ""},
		{"42", "int", 42, ""},
	}

	cfg := testCfg
	cfg.name = "update"

	for _, test := range iTests {
		err := cfg.UpdateIntE(test.key, test.value)
		if err != nil {
			if test.err == "" {
				t.Errorf("%s: unexpected error %q", test.name, err)
			} else {
				if test.err != err.Error() {
					t.Errorf("%s: expected %q got %q", test.name, test.err, err)
				}
			}
			continue
		}
		i, err := cfg.GetIntE(test.key)
		if err != nil {
			if test.err == "" {
				t.Errorf("%s: unexpected error %q", test.name, err)
			} else {
				if test.err != err.Error() {
					t.Errorf("%s: expected %q got %q", test.name, test.err, err)
				}
			}
			continue
		}
		if i != test.value {
			t.Errorf("%s: expected %q got %q", test.name, test.value, strconv.Itoa(i))
		}
	}
}

func TestUpdateInt64s(t *testing.T) {
	i64Tests := []struct {
		name  string
		key   string
		value int64
		err   string
	}{
		{"42", "coreint", int64(42), "config[update]: \"coreint\" is not updateable"},
		{"unset", "", int64(0), "config[update]: \"\" is not updateable"},
		{"42", "flagint", int64(42), ""},
		{"42", "cfgint", int64(42), ""},
		{"42", "int", int64(42), ""},
	}

	cfg := testCfg
	cfg.name = "update"

	for _, test := range i64Tests {
		err := cfg.UpdateInt64E(test.key, test.value)
		if err != nil {
			if test.err == "" {
				t.Errorf("%s: unexpected error %q", test.name, err)
			} else {
				if test.err != err.Error() {
					t.Errorf("%s: expected %q got %q", test.name, test.err, err)
				}
			}
			continue
		}
		i64, err := cfg.GetInt64E(test.key)
		if err != nil {
			if test.err == "" {
				t.Errorf("%s: unexpected error %q", test.name, err)
			} else {
				if test.err != err.Error() {
					t.Errorf("%s: expected %q got %q", test.name, test.err, err)
				}
			}
			continue
		}
		if i64 != test.value {
			t.Errorf("%s: expected %q got %q", test.name, test.value, strconv.Itoa(int(i64)))
		}
	}
}
