package contour

import (
	"testing"
)

func TestRegisterCfgFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		err      string
	}{
		{"empty", "", "cannot register configuration file: no name provided"},
		{"no extension", "cfg", "unable to determine cfg's format: no extension"},
		{"toml", "cfg.toml", ""},
		{"yaml", "cfg.yaml", ""},
		{"json", "cfg.json", ""},
		{"xml", "cfg.xml", "xml: unsupported configuration format"},
		{"undefined", "cfg.bss", "bss: unsupported configuration format"},
	}
	for _, test := range tests {
		cfg := New(test.name)
		err := cfg.RegisterCfgFile("cfg_file", test.filename)
		if err != nil {
			if test.err == "" {
				t.Errorf("RegisterCfgFilename %s: unexpected error: %q", test.name, err)
				goto cont
			}
			if test.err != err.Error() {
				t.Errorf("RegisterCfgFilename %s: expected error %q got %q", test.name, test.err, err.Error())
			}
		cont:
			continue
		}
		fname, err := cfg.StringE("cfg_file")
		if err != nil {
			t.Errorf("RegisterCfgFilename %s: unexpected error retrieving filename, %q", test.name, err)
			continue
		}
		if fname != test.filename {
			t.Errorf("RegisterCfgFilename %s: expected %q got %q", test.name, test.filename, fname)
			continue
		}
	}
}

func TestRegisterSettings(t *testing.T) {
	tests := []struct {
		name        string
		typ         dataType
		value       interface{}
		expected    interface{}
		expectedErr string
		checkValues bool
		IsCore      bool
		IsCfg       bool
		IsEnv       bool
		IsFlag      bool
	}{
		{"", _bool, true, true, "cannot register an unnamed setting", false, false, false, false, false},
		{"bool", _bool, true, true, "", true, false, false, false, false},
		{"bool", _bool, true, true, "bool is already registered, cannot re-register settings", true, false, false, false, false},
		{"", _int, 42, 42, "cannot register an unnamed setting", false, false, false, false, false},
		{"int", _int, 42, 42, "", true, false, false, false, false},
		{"int", _int, 84, 42, "int is already registered, cannot re-register settings", true, false, false, false, false},
		{"", _int64, int64(42), int64(42), "cannot register an unnamed setting", false, false, false, false, false},
		{"int64", _int64, int64(42), int64(42), "", true, false, false, false, false},
		{"int64", _int64, int64(84), int64(42), "int64 is already registered, cannot re-register settings", true, false, false, false, false},
		{"", _string, "bar", "bar", "cannot register an unnamed setting", false, false, false, false, false},
		{"string", _string, "bar", "bar", "", true, false, false, false, false},
		{"string", _string, "baz", "bar", "string is already registered, cannot re-register settings", true, false, false, false, false},
	}
	cfg := New("test register")
	var err error
	for i, test := range tests {
		switch test.typ {
		case _bool:
			err = cfg.RegisterBoolE(test.name, test.value.(bool))
		case _int:
			err = cfg.RegisterIntE(test.name, test.value.(int))
		case _int64:
			err = cfg.RegisterInt64E(test.name, test.value.(int64))
		case _string:
			err = cfg.RegisterStringE(test.name, test.value.(string))
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
		if cfg.Get(test.name) != test.expected {
			t.Errorf("%d: expected %v got %v", i, test.expected, cfg.Get(test.name))
		}
		if cfg.IsCore(test.name) != test.IsCore {
			t.Errorf("%d expected IsCore to be %v, got %v", i, test.IsCore, cfg.IsCore(test.name))
		}
		if cfg.IsCfg(test.name) != test.IsCfg {
			t.Errorf("%d expected IsCfg to be %v, got %v", i, test.IsCfg, cfg.IsCfg(test.name))
		}
		if cfg.IsEnv(test.name) != test.IsEnv {
			t.Errorf("%d expected IsEnv to be %v, got %v", i, test.IsEnv, cfg.IsEnv(test.name))
		}
		if cfg.IsFlag(test.name) != test.IsFlag {
			t.Errorf("%d expected IsFlag to be %v, got %v", i, test.IsFlag, cfg.IsFlag(test.name))
		}
	}
}

func TestRegisterCoreSettings(t *testing.T) {
	tests := []struct {
		name        string
		typ         dataType
		value       interface{}
		expected    interface{}
		expectedErr string
		checkValues bool
		IsCore      bool
		IsCfg       bool
		IsEnv       bool
		IsFlag      bool
	}{
		{"", _bool, true, true, "cannot register an unnamed setting", false, false, false, false, false},
		{"corebool", _bool, true, true, "", true, true, false, false, false},
		{"corebool", _bool, true, true, "corebool is already registered, cannot re-register settings", true, true, false, false, false},
		{"", _int, 42, 42, "cannot register an unnamed setting", false, false, false, false, false},
		{"coreint", _int, 42, 42, "", true, true, false, false, false},
		{"coreint", _int, 84, 42, "coreint is already registered, cannot re-register settings", true, true, false, false, false},
		{"", _int64, int64(42), int64(42), "cannot register an unnamed setting", false, false, false, false, false},
		{"coreint64", _int64, int64(42), int64(42), "", true, true, false, false, false},
		{"coreint64", _int64, int64(84), int64(42), "coreint64 is already registered, cannot re-register settings", true, true, false, false, false},
		{"", _string, "bar", "bar", "cannot register an unnamed setting", false, false, false, false, false},
		{"corestring", _string, "bar", "bar", "", true, true, false, false, false},
		{"corestring", _string, "baz", "bar", "corestring is already registered, cannot re-register settings", true, true, false, false, false},
	}
	tstSettings := New("test register")
	var err error
	for i, test := range tests {
		switch test.typ {
		case _bool:
			err = tstSettings.RegisterBoolCoreE(test.name, test.value.(bool))
		case _int:
			err = tstSettings.RegisterIntCoreE(test.name, test.value.(int))
		case _int64:
			err = tstSettings.RegisterInt64CoreE(test.name, test.value.(int64))
		case _string:
			err = tstSettings.RegisterStringCoreE(test.name, test.value.(string))
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
			t.Errorf("%d: expected %v got %v", i, test.expected, Get(test.name))
		}
		if tstSettings.IsCore(test.name) != test.IsCore {
			t.Errorf("%d expected IsCore to be %v, got %v", i, test.IsCore, IsCore(test.name))
		}
		if tstSettings.IsCfg(test.name) != test.IsCfg {
			t.Errorf("%d expected IsCfg to be %v, got %v", i, test.IsCfg, IsCfg(test.name))
		}
		if tstSettings.IsEnv(test.name) != test.IsEnv {
			t.Errorf("%d expected IsEnv to be %v, got %v", i, test.IsEnv, IsEnv(test.name))
		}
		if tstSettings.IsFlag(test.name) != test.IsFlag {
			t.Errorf("%d expected IsFlag to be %v, got %v", i, test.IsFlag, IsFlag(test.name))
		}
	}
	// NonE
	tstSettings = New("test register")
	for i, test := range tests {
		// because we aren't checking errors, don't test empty names
		if test.name == "" {
			continue
		}
		switch test.typ {
		case _bool:
			tstSettings.RegisterBoolCoreE(test.name, test.value.(bool))
		case _int:
			tstSettings.RegisterIntCoreE(test.name, test.value.(int))
		case _int64:
			tstSettings.RegisterInt64CoreE(test.name, test.value.(int64))
		case _string:
			tstSettings.RegisterStringCoreE(test.name, test.value.(string))
		default:
			t.Errorf("%d: unsupported typ: %s", i, test.typ)
			continue
		}
		if !test.checkValues {
			continue
		}
		if tstSettings.Get(test.name) != test.expected {
			t.Errorf("%d: expected %v got %v", i, test.expected, Get(test.name))
		}
		if tstSettings.IsCore(test.name) != test.IsCore {
			t.Errorf("%d expected IsCore to be %v, got %v", i, test.IsCore, IsCore(test.name))
		}
		if tstSettings.IsCfg(test.name) != test.IsCfg {
			t.Errorf("%d expected IsCfg to be %v, got %v", i, test.IsCfg, IsCfg(test.name))
		}
		if tstSettings.IsEnv(test.name) != test.IsEnv {
			t.Errorf("%d expected IsEnv to be %v, got %v", i, test.IsEnv, IsEnv(test.name))
		}
		if tstSettings.IsFlag(test.name) != test.IsFlag {
			t.Errorf("%d expected IsFlag to be %v, got %v", i, test.IsFlag, IsFlag(test.name))
		}
	}
}

func TestRegisterCfgSettings(t *testing.T) {
	tests := []struct {
		name        string
		typ         dataType
		value       interface{}
		expected    interface{}
		expectedErr string
		checkValues bool
		IsCore      bool
		IsCfg       bool
		IsEnv       bool
		IsFlag      bool
	}{
		{"", _bool, true, true, "cannot register an unnamed setting", false, false, false, false, false},
		{"cfgbool", _bool, true, true, "", true, false, true, true, false},
		{"cfgbool", _bool, false, true, "cfgbool is already registered, cannot re-register settings", true, false, true, true, false},
		{"", _int, 42, 42, "cannot register an unnamed setting", false, false, false, false, false},
		{"cfgint", _int, 42, 42, "", true, false, true, true, false},
		{"cfgint", _int, 84, 42, "cfgint is already registered, cannot re-register settings", true, false, true, true, false},
		{"", _int64, int64(42), int64(42), "cannot register an unnamed setting", false, false, false, false, false},
		{"cfgint64", _int64, int64(42), int64(42), "", true, false, true, true, false},
		{"cfgint64", _int64, int64(84), int64(42), "cfgint64 is already registered, cannot re-register settings", true, false, true, true, false},
		{"", _string, "bar", "bar", "cannot register an unnamed setting", false, false, false, false, false},
		{"cfgstring", _string, "bar", "bar", "", true, false, true, true, false},
		{"cfgstring", _string, "baz", "bar", "cfgstring is already registered, cannot re-register settings", true, false, true, true, false},
	}
	tstSettings := New("test register")
	var err error
	for i, test := range tests {
		switch test.typ {
		case _bool:
			err = tstSettings.RegisterBoolCfgE(test.name, test.value.(bool))
		case _int:
			err = tstSettings.RegisterIntCfgE(test.name, test.value.(int))
		case _int64:
			err = tstSettings.RegisterInt64CfgE(test.name, test.value.(int64))
		case _string:
			err = tstSettings.RegisterStringCfgE(test.name, test.value.(string))
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
			t.Errorf("%d: expected %v got %v", i, test.expected, Get(test.name))
		}
		if tstSettings.IsCore(test.name) != test.IsCore {
			t.Errorf("%d expected IsCore to be %v, got %v", i, test.IsCore, IsCore(test.name))
		}
		if tstSettings.IsCfg(test.name) != test.IsCfg {
			t.Errorf("%d expected IsCfg to be %v, got %v", i, test.IsCfg, IsCfg(test.name))
		}
		if tstSettings.IsEnv(test.name) != test.IsEnv {
			t.Errorf("%d expected IsEnv to be %v, got %v", i, test.IsEnv, IsEnv(test.name))
		}
		if tstSettings.IsFlag(test.name) != test.IsFlag {
			t.Errorf("%d expected IsFlag to be %v, got %v", i, test.IsFlag, IsFlag(test.name))
		}
	}
	// Non-E
	tstSettings = New("test")
	for i, test := range tests {
		// skip empty names since we don't check errors
		if test.name == "" {
			continue
		}
		switch test.typ {
		case _bool:
			tstSettings.RegisterBoolCfg(test.name, test.value.(bool))
		case _int:
			tstSettings.RegisterIntCfg(test.name, test.value.(int))
		case _int64:
			tstSettings.RegisterInt64Cfg(test.name, test.value.(int64))
		case _string:
			tstSettings.RegisterStringCfg(test.name, test.value.(string))
		default:
			t.Errorf("%d: unsupported typ: %s", i, test.typ)
			continue
		}
		if !test.checkValues {
			continue
		}
		if tstSettings.Get(test.name) != test.expected {
			t.Errorf("%d: expected %v got %v", i, test.expected, Get(test.name))
		}
		if tstSettings.IsCore(test.name) != test.IsCore {
			t.Errorf("%d expected IsCore to be %v, got %v", i, test.IsCore, IsCore(test.name))
		}
		if tstSettings.IsCfg(test.name) != test.IsCfg {
			t.Errorf("%d expected IsCfg to be %v, got %v", i, test.IsCfg, IsCfg(test.name))
		}
		if tstSettings.IsEnv(test.name) != test.IsEnv {
			t.Errorf("%d expected IsEnv to be %v, got %v", i, test.IsEnv, IsEnv(test.name))
		}
		if tstSettings.IsFlag(test.name) != test.IsFlag {
			t.Errorf("%d expected IsFlag to be %v, got %v", i, test.IsFlag, IsFlag(test.name))
		}
	}

}

func TestRegisterFlagSettings(t *testing.T) {
	tests := []struct {
		name        string
		short       string
		typ         dataType
		value       interface{}
		expected    interface{}
		expectedErr string
		checkValues bool
		IsCore      bool
		IsCfg       bool
		IsEnv       bool
		IsFlag      bool
	}{
		{"", "", _bool, true, true, "cannot register an unnamed setting", false, false, false, false, false},
		{"flagbool", "b", _bool, true, true, "", true, false, true, true, true},
		{"flagbool", "", _bool, false, true, "flagbool is already registered, cannot re-register settings", true, false, true, true, true},
		{"", "", _int, 42, 42, "cannot register an unnamed setting", false, false, false, false, false},
		{"flagint", "i", _int, 42, 42, "", true, false, true, true, true},
		{"flagint", "", _int, 84, 42, "flagint is already registered, cannot re-register settings", true, false, true, true, true},
		{"", "", _int64, int64(42), int64(42), "cannot register an unnamed setting", false, false, false, false, false},
		{"flagint64", "6", _int64, int64(42), int64(42), "", true, false, true, true, true},
		{"flagint64", "", _int64, int64(84), int64(42), "flagint64 is already registered, cannot re-register settings", true, false, true, true, true},
		{"", "", _string, "bar", "bar", "cannot register an unnamed setting", false, false, false, false, false},
		{"flagstring", "s", _string, "bar", "bar", "", true, false, true, true, true},
		{"flagstring", "", _string, "baz", "bar", "flagstring is already registered, cannot re-register settings", true, false, true, true, true},
	}
	tstSettings := New("test register")
	var err error
	for i, test := range tests {
		switch test.typ {
		case _bool:
			err = tstSettings.RegisterBoolFlagE(test.name, test.short, test.value.(bool), "", "usage")
		case _int:
			err = tstSettings.RegisterIntFlagE(test.name, test.short, test.value.(int), "", "usage")
		case _int64:
			err = tstSettings.RegisterInt64FlagE(test.name, test.short, test.value.(int64), "", "usage")
		case _string:
			err = tstSettings.RegisterStringFlagE(test.name, test.short, test.value.(string), "", "usage")
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
			t.Errorf("%d: expected %v got %v", i, test.expected, Get(test.name))
		}
		if tstSettings.IsCore(test.name) != test.IsCore {
			t.Errorf("%d expected IsCore to be %v, got %v", i, test.IsCore, IsCore(test.name))
		}
		if tstSettings.IsCfg(test.name) != test.IsCfg {
			t.Errorf("%d expected IsCfg to be %v, got %v", i, test.IsCfg, IsCfg(test.name))
		}
		if tstSettings.IsEnv(test.name) != test.IsEnv {
			t.Errorf("%d expected IsEnv to be %v, got %v", i, test.IsEnv, IsEnv(test.name))
		}
		if tstSettings.IsFlag(test.name) != test.IsFlag {
			t.Errorf("%d expected IsFlag to be %v, got %v", i, test.IsFlag, IsFlag(test.name))
		}
	}
	// Non-E
	tstSettings = New("test register")
	for i, test := range tests {
		// since no error checking is being done, we skip empty names
		switch test.typ {
		case _bool:
			tstSettings.RegisterBoolFlag(test.name, test.short, test.value.(bool), "", "usage")
		case _int:
			tstSettings.RegisterIntFlag(test.name, test.short, test.value.(int), "", "usage")
		case _int64:
			tstSettings.RegisterInt64Flag(test.name, test.short, test.value.(int64), "", "usage")
		case _string:
			tstSettings.RegisterStringFlag(test.name, test.short, test.value.(string), "", "usage")
		default:
			t.Errorf("%d: unsupported typ: %s", i, test.typ)
			continue
		}
		if !test.checkValues {
			continue
		}
		if tstSettings.Get(test.name) != test.expected {
			t.Errorf("%d: expected %v got %v", i, test.expected, Get(test.name))
		}
		if tstSettings.IsCore(test.name) != test.IsCore {
			t.Errorf("%d expected IsCore to be %v, got %v", i, test.IsCore, IsCore(test.name))
		}
		if tstSettings.IsCfg(test.name) != test.IsCfg {
			t.Errorf("%d expected IsCfg to be %v, got %v", i, test.IsCfg, IsCfg(test.name))
		}
		if tstSettings.IsEnv(test.name) != test.IsEnv {
			t.Errorf("%d expected IsEnv to be %v, got %v", i, test.IsEnv, IsEnv(test.name))
		}
		if tstSettings.IsFlag(test.name) != test.IsFlag {
			t.Errorf("%d expected IsFlag to be %v, got %v", i, test.IsFlag, IsFlag(test.name))
		}
	}
}
