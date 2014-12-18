package contour

import (
	"strconv"
	"testing"
)

func TestUpdates(t *testing.T) {
	bTests := []struct {
		name  string
		key   string
		value string
		err   string
	}{
		{"false", "corebool", "false", "config[update]: \"corebool\" is not updateable"},
		{"true", "corebool", "t", "config[update]: \"corebool\" is not updateable"},
		{"unset", "", "", "config[update]: \"\" is not updateable"},
		{"false", "flagbool", "false", ""},
		{"true", "flagbool", "t", ""},
		{"false", "cfgbool", "false", ""},
		{"true", "cfgbool", "t", ""},
		{"false", "bool", "false", ""},
		{"true", "bool", "t", ""},
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
		if b != test.value {
			t.Errorf("%s: expected %q got %q", test.name, test.value, b)
		}
	}

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
			t.Errorf("%s: expected %q got %q", test.name, test.value, s)
		}
	}

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
