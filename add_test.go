package contour

import "testing"

func TestAddCoreSettings(t *testing.T) {
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
		{"corebool", _bool, true, true, "", true, true, false, false, false},
		{"corebool", _bool, true, true, "corebool: core setting exists", true, true, false, false, false},
		{"", _int, 42, 42, "no setting name provided", false, false, false, false, false},
		{"coreint", _int, 42, 42, "", true, true, false, false, false},
		{"coreint", _int, 84, 42, "coreint: core setting exists", true, true, false, false, false},
		{"", _int64, int64(42), int64(42), "no setting name provided", false, false, false, false, false},
		{"coreint64", _int64, int64(42), int64(42), "", true, true, false, false, false},
		{"coreint64", _int64, int64(84), int64(42), "coreint64: core setting exists", true, true, false, false, false},
		{"", _string, "bar", "bar", "no setting name provided", false, false, false, false, false},
		{"corestring", _string, "bar", "bar", "", true, true, false, false, false},
		{"corestring", _string, "baz", "bar", "corestring: core setting exists", true, true, false, false, false},
	}
	tstSettings := New("test register")
	var err error
	for i, test := range tests {
		switch test.typ {
		case _bool:
			err = tstSettings.AddBoolCore(test.name, test.value.(bool))
		case _int:
			err = tstSettings.AddIntCore(test.name, test.value.(int))
		case _int64:
			err = tstSettings.AddInt64Core(test.name, test.value.(int64))
		case _string:
			err = tstSettings.AddStringCore(test.name, test.value.(string))
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

func TestAddSettings(t *testing.T) {
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
		{"bool", _bool, true, true, "", true, false, false, false, false},
		{"bool", _bool, true, true, "bool: setting exists", true, false, false, false, false},
		{"", _int, 42, 42, "no setting name provided", false, false, false, false, false},
		{"int", _int, 42, 42, "", true, false, false, false, false},
		{"int", _int, 84, 42, "int: setting exists", true, false, false, false, false},
		{"", _int64, int64(42), int64(42), "no setting name provided", false, false, false, false, false},
		{"int64", _int64, int64(42), int64(42), "", true, false, false, false, false},
		{"int64", _int64, int64(84), int64(42), "int64: setting exists", true, false, false, false, false},
		{"", _string, "bar", "bar", "no setting name provided", false, false, false, false, false},
		{"string", _string, "bar", "bar", "", true, false, false, false, false},
		{"string", _string, "baz", "bar", "string: setting exists", true, false, false, false, false},
	}
	tstSettings := New("test add")
	var err error
	for i, test := range tests {
		switch test.typ {
		case _bool:
			err = tstSettings.AddBool(test.name, test.value.(bool))
		case _int:
			err = tstSettings.AddInt(test.name, test.value.(int))
		case _int64:
			err = tstSettings.AddInt64(test.name, test.value.(int64))
		case _string:
			err = tstSettings.AddString(test.name, test.value.(string))
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
