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
		{"", _bool, true, true, "no setting name provided", false, false, false, false, false},
		{"xx_cfgbool", _bool, true, true, "", true, false, true, false, false},
		{"xx_cfgbool", _bool, false, true, "xx_cfgbool: configuration file var setting exists", true, false, true, false, false},
		{"", _int, 42, 42, "no setting name provided", false, false, false, false, false},
		{"xx_cfgint", _int, 42, 42, "", true, false, true, false, false},

		{"xx_cfgint", _int, 84, 42, "xx_cfgint: configuration file var setting exists", true, false, true, false, false},
		{"", _int64, int64(42), int64(42), "no setting name provided", false, false, false, false, false},
		{"xx_cfgint64", _int64, int64(42), int64(42), "", true, false, true, false, false},
		{"xx_cfgint64", _int64, int64(84), int64(42), "xx_cfgint64: configuration file var setting exists", true, false, true, false, false},
		{"", _interface, 42, 42, "no setting name provided", false, false, false, false, false},

		{"xx_cfginterface", _interface, 42, 42, "", true, false, true, false, false},
		{"xx_cfginterface", _interface, 84, 42, "xx_cfginterface: configuration file var setting exists", true, false, true, false, false},
		{"", _string, "bar", "bar", "no setting name provided", false, false, false, false, false},
		{"xx_cfgstring", _string, "bar", "bar", "", true, false, true, false, false},
		{"xx_cfgstring", _string, "baz", "bar", "xx_cfgstring: configuration file var setting exists", true, false, true, false, false},
	}
	var err error
	for i, test := range tests {
		switch test.typ {
		case _bool:
			err = RegisterBoolConfFileVar(test.name, test.value.(bool))
		case _int:
			err = RegisterIntConfFileVar(test.name, test.value.(int))
		case _int64:
			err = RegisterInt64ConfFileVar(test.name, test.value.(int64))
		case _string:
			err = RegisterStringConfFileVar(test.name, test.value.(string))
		case _interface:
			err = RegisterInterfaceConfFileVar(test.name, test.value)
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
		if Get(test.name) != test.expected {
			t.Errorf("%d: expected %v got %v", i, test.expected, Get(test.name))
		}
		if IsCore(test.name) != test.IsCore {
			t.Errorf("%d: expected IsCore to be %v, got %v", i, test.IsCore, IsCore(test.name))
		}
		if IsConfFileVar(test.name) != test.IsConfFileVar {
			t.Errorf("%d: expected IsConfFileVar to be %v, got %v", i, test.IsConfFileVar, IsConfFileVar(test.name))
		}
		if IsEnvVar(test.name) != test.IsEnvVar {
			t.Errorf("%d: expected IsEnvVar to be %v, got %v", i, test.IsEnvVar, IsEnvVar(test.name))
		}
		if IsFlag(test.name) != test.IsFlag {
			t.Errorf("%d: expected IsFlag to be %v, got %v", i, test.IsFlag, IsFlag(test.name))
		}
		if !UseConfFile() {
			t.Errorf("%d: useConfFile: got %v; want true", i, UseConfFile())
		}
		if UseEnvVars() {
			t.Errorf("%d: useEnvVars: got %v; want false", i, UseEnvVars())
		}
		if UseFlags() {
			t.Errorf("%d: useFlags: got %v; want false", i, UseFlags())
		}
	}
}

func TestRegisterEnvVarSettings(t *testing.T) {
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
		{"", _bool, true, true, "no setting name provided", false, false, false, false, false},
		{"envbool", _bool, true, true, "", true, false, true, true, false},
		{"envbool", _bool, false, true, "envbool: env var setting exists", true, false, true, true, false},
		{"", _int, 42, 42, "no setting name provided", false, false, false, false, false},
		{"envint", _int, 42, 42, "", true, false, true, true, false},

		{"envint", _int, 84, 42, "envint: env var setting exists", true, false, true, true, false},
		{"", _int64, int64(42), int64(42), "no setting name provided", false, false, false, false, false},
		{"envint64", _int64, int64(42), int64(42), "", true, false, true, true, false},
		{"envint64", _int64, int64(84), int64(42), "envint64: env var setting exists", true, false, true, true, false},
		{"", _string, "bar", "bar", "no setting name provided", false, false, false, false, false},

		{"envstring", _string, "bar", "bar", "", true, false, true, true, false},
		{"envstring", _string, "baz", "bar", "envstring: env var setting exists", true, false, true, true, false},
	}
	tstSettings := New("test register")
	var err error
	for i, test := range tests {
		switch test.typ {
		case _bool:
			err = tstSettings.RegisterBoolEnvVar(test.name, test.value.(bool))
		case _int:
			err = tstSettings.RegisterIntEnvVar(test.name, test.value.(int))
		case _int64:
			err = tstSettings.RegisterInt64EnvVar(test.name, test.value.(int64))
		case _string:
			err = tstSettings.RegisterStringEnvVar(test.name, test.value.(string))
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
		if !tstSettings.UseConfFile() {
			t.Errorf("%d: useConfFile: got %v; want true", i, tstSettings.UseConfFile())
		}
		if !tstSettings.UseEnvVars() {
			t.Errorf("%d: useEnvVars: got %v; want true", i, tstSettings.UseEnvVars())
		}
		if tstSettings.UseFlags() {
			t.Errorf("%d: useFlags: got %v; want false", i, tstSettings.UseFlags())
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
		{"", "", _bool, true, true, "no setting name provided", false, false, false, false, false},
		{"flagbool", "b", _bool, true, true, "", true, false, true, true, true},
		{"flagbool", "", _bool, false, true, "flagbool: flag setting exists", true, false, true, true, true},
		{"flagboolz", "b", _bool, true, true, "flagboolz: short flag \"b\" already exists for \"flagbool\"", true, false, true, true, true},
		{"", "", _int, 42, 42, "no setting name provided", false, false, false, false, false},

		{"flagint", "i", _int, 42, 42, "", true, false, true, true, true},
		{"flagint", "", _int, 84, 42, "flagint: flag setting exists", true, false, true, true, true},
		{"flagintz", "i", _int, 42, 42, "flagintz: short flag \"i\" already exists for \"flagint\"", true, false, true, true, true},
		{"", "", _int64, int64(42), int64(42), "no setting name provided", false, false, false, false, false},
		{"flagint64", "6", _int64, int64(42), int64(42), "", true, false, true, true, true},

		{"flagint64", "", _int64, int64(84), int64(42), "flagint64: flag setting exists", true, false, true, true, true},
		{"flagint64z", "6", _int64, int64(42), int64(42), "flagint64z: short flag \"6\" already exists for \"flagint64\"", true, false, true, true, true},
		{"", "", _string, "bar", "bar", "no setting name provided", false, false, false, false, false},
		{"flagstring", "s", _string, "bar", "bar", "", true, false, true, true, true},
		{"flagstring", "", _string, "baz", "bar", "flagstring: flag setting exists", true, false, true, true, true},

		{"flagstringz", "s", _string, "bar", "bar", "flagstringz: short flag \"s\" already exists for \"flagstring\"", true, false, true, true, true},
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
			// don't do any additional comparisions for short flag errors
			if _, ok := err.(ShortFlagExistsError); ok {
				continue
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
		if !tstSettings.UseConfFile() {
			t.Errorf("%d: useConfFile: got %v; want true", i, tstSettings.UseConfFile())
		}
		if !tstSettings.UseEnvVars() {
			t.Errorf("%d: useEnvVars: got %v; want true", i, tstSettings.UseEnvVars())
		}
		if !tstSettings.UseFlags() {
			t.Errorf("%d: useFlags: got %v; want true", i, tstSettings.UseFlags())
		}
	}
}
