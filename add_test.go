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
		{"x_corebool", _bool, true, true, "", true, true, false, false, false},
		{"x_corebool", _bool, true, true, "x_corebool: core setting exists", true, true, false, false, false},
		{"", _int, 42, 42, "no setting name provided", false, false, false, false, false},
		{"x_coreint", _int, 42, 42, "", true, true, false, false, false},
		{"x_coreint", _int, 84, 42, "x_coreint: core setting exists", true, true, false, false, false},
		{"", _int64, int64(42), int64(42), "no setting name provided", false, false, false, false, false},
		{"x_coreint64", _int64, int64(42), int64(42), "", true, true, false, false, false},
		{"x_coreint64", _int64, int64(84), int64(42), "x_coreint64: core setting exists", true, true, false, false, false},
		{"", _interface, 42, 42, "no setting name provided", false, false, false, false, false},
		{"x_coreinterface", _interface, 42, 42, "", true, true, false, false, false},
		{"x_coreinterface", _interface, 42, 42, "x_coreinterface: core setting exists", true, true, false, false, false},
		{"", _string, "bar", "bar", "no setting name provided", false, false, false, false, false},
		{"x_corestring", _string, "bar", "bar", "", true, true, false, false, false},
		{"x_corestring", _string, "baz", "bar", "x_corestring: core setting exists", true, true, false, false, false},
	}
	var err error
	for i, test := range tests {
		switch test.typ {
		case _bool:
			err = AddBoolCore(test.name, test.value.(bool))
		case _int:
			err = AddIntCore(test.name, test.value.(int))
		case _int64:
			err = AddInt64Core(test.name, test.value.(int64))
		case _string:
			err = AddStringCore(test.name, test.value.(string))
		case _interface:
			err = AddInterfaceCore(test.name, test.value)
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
			t.Errorf("%d expected IsCore to be %v, got %v", i, test.IsCore, IsCore(test.name))
		}
		if IsConfFileVar(test.name) != test.IsConfFileVar {
			t.Errorf("%d expected IsConfFileVar to be %v, got %v", i, test.IsConfFileVar, IsConfFileVar(test.name))
		}
		if IsEnvVar(test.name) != test.IsEnvVar {
			t.Errorf("%d expected IsEnvVar to be %v, got %v", i, test.IsEnvVar, IsEnvVar(test.name))
		}
		if IsFlag(test.name) != test.IsFlag {
			t.Errorf("%d expected IsFlag to be %v, got %v", i, test.IsFlag, IsFlag(test.name))
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
		{"x_bool", _bool, true, true, "", true, false, false, false, false},
		{"x_bool", _bool, true, true, "x_bool: setting exists", true, false, false, false, false},
		{"", _int, 42, 42, "no setting name provided", false, false, false, false, false},
		{"x_int", _int, 42, 42, "", true, false, false, false, false},
		{"x_int", _int, 84, 42, "x_int: setting exists", true, false, false, false, false},
		{"", _int64, int64(42), int64(42), "no setting name provided", false, false, false, false, false},
		{"x_int64", _int64, int64(42), int64(42), "", true, false, false, false, false},
		{"x_int64", _int64, int64(84), int64(42), "x_int64: setting exists", true, false, false, false, false},
		{"", _int, 42, 42, "no setting name provided", false, false, false, false, false},
		{"x_interface", _interface, 42, 42, "", true, false, false, false, false},
		{"x_interface", _interface, 84, 42, "x_interface: setting exists", true, false, false, false, false},
		{"", _string, "bar", "bar", "no setting name provided", false, false, false, false, false},
		{"x_string", _string, "bar", "bar", "", true, false, false, false, false},
		{"x_string", _string, "baz", "bar", "x_string: setting exists", true, false, false, false, false},
	}
	var err error
	for i, test := range tests {
		switch test.typ {
		case _bool:
			err = AddBool(test.name, test.value.(bool))
		case _int:
			err = AddInt(test.name, test.value.(int))
		case _int64:
			err = AddInt64(test.name, test.value.(int64))
		case _string:
			err = AddString(test.name, test.value.(string))
		case _interface:
			err = AddInterface(test.name, test.value)
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
			t.Errorf("%d expected IsCore to be %v, got %v", i, test.IsCore, IsCore(test.name))
		}
		if IsConfFileVar(test.name) != test.IsConfFileVar {
			t.Errorf("%d expected IsConfFileVar to be %v, got %v", i, test.IsConfFileVar, IsConfFileVar(test.name))
		}
		if IsEnvVar(test.name) != test.IsEnvVar {
			t.Errorf("%d expected IsEnvVar to be %v, got %v", i, test.IsEnvVar, IsEnvVar(test.name))
		}
		if IsFlag(test.name) != test.IsFlag {
			t.Errorf("%d expected IsFlag to be %v, got %v", i, test.IsFlag, IsFlag(test.name))
		}
	}
}
