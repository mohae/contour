package contour

import (
	"testing"
)

func TestRegisterCfgSettings(t *testing.T) {
	tests := []struct {
		name          string
		typ           dataType
		value         interface{}
		expected      interface{}
		expectedErr   string
		checkValues   bool
		IsCore        bool
		IsConfFileVar bool
		IsEnvVar      bool
		IsFlag        bool
	}{
		{"", _bool, true, true, "registration failed: setting name was empty", false, false, false, false, false},
		{"cfgbool", _bool, true, true, "", true, false, true, true, false},
		{"cfgbool", _bool, false, true, "cfgbool: registration failed: setting exists", true, false, true, true, false},
		{"", _int, 42, 42, "registration failed: setting name was empty", false, false, false, false, false},
		{"cfgint", _int, 42, 42, "", true, false, true, true, false},
		{"cfgint", _int, 84, 42, "cfgint: registration failed: setting exists", true, false, true, true, false},
		{"", _int64, int64(42), int64(42), "registration failed: setting name was empty", false, false, false, false, false},
		{"cfgint64", _int64, int64(42), int64(42), "", true, false, true, true, false},
		{"cfgint64", _int64, int64(84), int64(42), "cfgint64: registration failed: setting exists", true, false, true, true, false},
		{"", _string, "bar", "bar", "registration failed: setting name was empty", false, false, false, false, false},
		{"cfgstring", _string, "bar", "bar", "", true, false, true, true, false},
		{"cfgstring", _string, "baz", "bar", "cfgstring: registration failed: setting exists", true, false, true, true, false},
	}
	tstSettings := New("test register")
	var err error
	for i, test := range tests {
		switch test.typ {
		case _bool:
			err = tstSettings.RegisterBoolConfFileVar(test.name, test.value.(bool))
		case _int:
			err = tstSettings.RegisterIntConfFileVar(test.name, test.value.(int))
		case _int64:
			err = tstSettings.RegisterInt64ConfFileVar(test.name, test.value.(int64))
		case _string:
			err = tstSettings.RegisterStringConfFileVar(test.name, test.value.(string))
		default:
			t.Errorf("%d: unsupported typ: %s", i, test.typ)
			continue
		}
		if err != nil {
			if err.Error() != test.expectedErr {
				t.Errorf("%d error: expected %s got %s", i, test.expectedErr, err.Error())
			}
		}
		if !test.checkValues {
			continue
		}
		if tstSettings.Get(test.name) != test.expected {
			t.Errorf("%d: expected %v got %v", i, test.expected, tstSettings.Get(test.name))
		}
		if tstSettings.IsCore(test.name) != test.IsCore {
			t.Errorf("%d expected IsCore to be %v, got %v", i, test.IsCore, tstSettings.IsCore(test.name))
		}
		if tstSettings.IsConfFileVar(test.name) != test.IsConfFileVar {
			t.Errorf("%d expected IsConfFileVar to be %v, got %v", i, test.IsConfFileVar, tstSettings.IsConfFileVar(test.name))
		}
		if tstSettings.IsEnvVar(test.name) != test.IsEnvVar {
			t.Errorf("%d expected IsEnvVar to be %v, got %v", i, test.IsEnvVar, tstSettings.IsEnvVar(test.name))
		}
		if tstSettings.IsFlag(test.name) != test.IsFlag {
			t.Errorf("%d expected IsFlag to be %v, got %v", i, test.IsFlag, tstSettings.IsFlag(test.name))
		}
	}
}

func TestRegisterFlagSettings(t *testing.T) {
	tests := []struct {
		name          string
		short         string
		typ           dataType
		value         interface{}
		expected      interface{}
		expectedErr   string
		checkValues   bool
		IsCore        bool
		IsConfFileVar bool
		IsEnvVar      bool
		IsFlag        bool
	}{
		{"", "", _bool, true, true, "registration failed: setting name was empty", false, false, false, false, false},
		{"flagbool", "b", _bool, true, true, "", true, false, true, true, true},
		{"flagbool", "", _bool, false, true, "flagbool: registration failed: setting exists", true, false, true, true, true},
		{"", "", _int, 42, 42, "registration failed: setting name was empty", false, false, false, false, false},
		{"flagint", "i", _int, 42, 42, "", true, false, true, true, true},
		{"flagint", "", _int, 84, 42, "flagint: registration failed: setting exists", true, false, true, true, true},
		{"", "", _int64, int64(42), int64(42), "registration failed: setting name was empty", false, false, false, false, false},
		{"flagint64", "6", _int64, int64(42), int64(42), "", true, false, true, true, true},
		{"flagint64", "", _int64, int64(84), int64(42), "flagint64: registration failed: setting exists", true, false, true, true, true},
		{"", "", _string, "bar", "bar", "registration failed: setting name was empty", false, false, false, false, false},
		{"flagstring", "s", _string, "bar", "bar", "", true, false, true, true, true},
		{"flagstring", "", _string, "baz", "bar", "flagstring: registration failed: setting exists", true, false, true, true, true},
	}
	tstSettings := New("test register")
	var err error
	for i, test := range tests {
		switch test.typ {
		case _bool:
			err = tstSettings.RegisterBoolFlag(test.name, test.short, test.value.(bool), "", "usage")
		case _int:
			err = tstSettings.RegisterIntFlag(test.name, test.short, test.value.(int), "", "usage")
		case _int64:
			err = tstSettings.RegisterInt64Flag(test.name, test.short, test.value.(int64), "", "usage")
		case _string:
			err = tstSettings.RegisterStringFlag(test.name, test.short, test.value.(string), "", "usage")
		default:
			t.Errorf("%d: unsupported typ: %s", i, test.typ)
			continue
		}
		if err != nil {
			if err.Error() != test.expectedErr {
				t.Errorf("%d error: expected %s got %s", i, test.expectedErr, err.Error())
			}
		}
		if !test.checkValues {
			continue
		}
		if tstSettings.Get(test.name) != test.expected {
			t.Errorf("%d: expected %v got %v", i, test.expected, tstSettings.Get(test.name))
		}
		if tstSettings.IsCore(test.name) != test.IsCore {
			t.Errorf("%d expected IsCore to be %v, got %v", i, test.IsCore, tstSettings.IsCore(test.name))
		}
		if tstSettings.IsConfFileVar(test.name) != test.IsConfFileVar {
			t.Errorf("%d expected IsConfFileVar to be %v, got %v", i, test.IsConfFileVar, tstSettings.IsConfFileVar(test.name))
		}
		if tstSettings.IsEnvVar(test.name) != test.IsEnvVar {
			t.Errorf("%d expected IsEnvVar to be %v, got %v", i, test.IsEnvVar, tstSettings.IsEnvVar(test.name))
		}
		if tstSettings.IsFlag(test.name) != test.IsFlag {
			t.Errorf("%d expected IsFlag to be %v, got %v", i, test.IsFlag, tstSettings.IsFlag(test.name))
		}
	}
}
