package contour

import (
	"testing"
)

func TestRegisterCfgFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		format   string
		err      string
	}{
		{"empty", "", "", "RegisterCfgFile expected a cfg filename: none received"},
		{"no extension", "cfg", "", "unable to determine cfg's config format: no extension"},
		{"toml", "cfg.toml", "toml", ""},
		{"yaml", "cfg.yaml", "yaml", ""},
		{"json", "cfg.json", "json", ""},
		{"xml", "cfg.xml", "xml", "unsupported cfg format: xml"},
		{"undefined", "cfg.bss", "bss", "unsupported cfg format: bss"},
	}
	for _, test := range tests {
		cfg := NewCfg(test.name)
		err := cfg.RegisterCfgFile(CfgFile, test.filename)
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
		fname, err := cfg.GetStringE(CfgFile)
		if err != nil {
			t.Errorf("RegisterCfgFilename %s: unexpected error retrieving filename, %q", test.name, err)
			continue
		}
		if fname != test.filename {
			t.Errorf("RegisterCfgFilename %s: expected %q got %q", test.name, test.filename, fname)
			continue
		}
		format, err := cfg.GetStringE(CfgFormat)
		if err != nil {
			t.Errorf("RegisterCfgFilename format %s: unexpected error retrieving ext, %q", test.name, err)
			continue
		}
		if format != test.format {
			t.Errorf("RegisterCfgFilename format %s: expected %q got %q", test.name, test.format, format)
		}
	}
}

func TestRegisterSettings(t *testing.T) {
	tests := []struct {
		name        string
		typ         string
		value       interface{}
		expected    interface{}
		expectedErr string
		checkValues bool
		IsCore      bool
		IsCfg       bool
		IsEnv       bool
		IsFlag      bool
	}{
		{"", "bool", true, true, "cannot register an unnamed setting", false, false, false, false, false},
		{"bool", "bool", true, true, "", true, false, false, false, false},
		{"bool", "bool", true, true, "bool is already registered, cannot re-register settings", true, false, false, false, false},
		{"", "int", 42, 42, "cannot register an unnamed setting", false, false, false, false, false},
		{"int", "int", 42, 42, "", true, false, false, false, false},
		{"int", "int", 84, 42, "int is already registered, cannot re-register settings", true, false, false, false, false},
		{"", "int64", int64(42), int64(42), "cannot register an unnamed setting", false, false, false, false, false},
		{"int64", "int64", int64(42), int64(42), "", true, false, false, false, false},
		{"int64", "int64", int64(84), int64(42), "int64 is already registered, cannot re-register settings", true, false, false, false, false},
		{"", "string", "bar", "bar", "cannot register an unnamed setting", false, false, false, false, false},
		{"string", "string", "bar", "bar", "", true, false, false, false, false},
		{"string", "string", "baz", "bar", "string is already registered, cannot re-register settings", true, false, false, false, false},
	}
	cfg := NewCfg("test register")
	var err error
	for i, test := range tests {
		switch test.typ {
		case "bool":
			err = cfg.RegisterBoolE(test.name, test.value.(bool))
		case "int":
			err = cfg.RegisterIntE(test.name, test.value.(int))
		case "int64":
			err = cfg.RegisterInt64E(test.name, test.value.(int64))
		case "string":
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
		typ         string
		value       interface{}
		expected    interface{}
		expectedErr string
		checkValues bool
		IsCore      bool
		IsCfg       bool
		IsEnv       bool
		IsFlag      bool
	}{
		{"", "bool", true, true, "cannot register an unnamed setting", false, false, false, false, false},
		{"corebool", "bool", true, true, "", true, true, false, false, false},
		{"corebool", "bool", true, true, "corebool is already registered, cannot re-register settings", true, true, false, false, false},
		{"", "int", 42, 42, "cannot register an unnamed setting", false, false, false, false, false},
		{"coreint", "int", 42, 42, "", true, true, false, false, false},
		{"coreint", "int", 84, 42, "coreint is already registered, cannot re-register settings", true, true, false, false, false},
		{"", "int64", int64(42), int64(42), "cannot register an unnamed setting", false, false, false, false, false},
		{"coreint64", "int64", int64(42), int64(42), "", true, true, false, false, false},
		{"coreint64", "int64", int64(84), int64(42), "coreint64 is already registered, cannot re-register settings", true, true, false, false, false},
		{"", "string", "bar", "bar", "cannot register an unnamed setting", false, false, false, false, false},
		{"corestring", "string", "bar", "bar", "", true, true, false, false, false},
		{"corestring", "string", "baz", "bar", "corestring is already registered, cannot re-register settings", true, true, false, false, false},
	}
	appCfg = NewCfg("test register")
	var err error
	for i, test := range tests {
		switch test.typ {
		case "bool":
			err = RegisterBoolCoreE(test.name, test.value.(bool))
		case "int":
			err = RegisterIntCoreE(test.name, test.value.(int))
		case "int64":
			err = RegisterInt64CoreE(test.name, test.value.(int64))
		case "string":
			err = RegisterStringCoreE(test.name, test.value.(string))
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
		if IsCfg(test.name) != test.IsCfg {
			t.Errorf("%d expected IsCfg to be %v, got %v", i, test.IsCfg, IsCfg(test.name))
		}
		if IsEnv(test.name) != test.IsEnv {
			t.Errorf("%d expected IsEnv to be %v, got %v", i, test.IsEnv, IsEnv(test.name))
		}
		if IsFlag(test.name) != test.IsFlag {
			t.Errorf("%d expected IsFlag to be %v, got %v", i, test.IsFlag, IsFlag(test.name))
		}
	}
	// NonE
	appCfg = NewCfg("test register")
	for i, test := range tests {
		// because we aren't checking errors, don't test empty names
		if test.name == "" {
			continue
		}
		switch test.typ {
		case "bool":
			RegisterBoolCoreE(test.name, test.value.(bool))
		case "int":
			RegisterIntCoreE(test.name, test.value.(int))
		case "int64":
			RegisterInt64CoreE(test.name, test.value.(int64))
		case "string":
			RegisterStringCoreE(test.name, test.value.(string))
		default:
			t.Errorf("%d: unsupported typ: %s", i, test.typ)
			continue
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
		if IsCfg(test.name) != test.IsCfg {
			t.Errorf("%d expected IsCfg to be %v, got %v", i, test.IsCfg, IsCfg(test.name))
		}
		if IsEnv(test.name) != test.IsEnv {
			t.Errorf("%d expected IsEnv to be %v, got %v", i, test.IsEnv, IsEnv(test.name))
		}
		if IsFlag(test.name) != test.IsFlag {
			t.Errorf("%d expected IsFlag to be %v, got %v", i, test.IsFlag, IsFlag(test.name))
		}
	}
}

func TestRegisterCfgSettings(t *testing.T) {
	tests := []struct {
		name        string
		typ         string
		value       interface{}
		expected    interface{}
		expectedErr string
		checkValues bool
		IsCore      bool
		IsCfg       bool
		IsEnv       bool
		IsFlag      bool
	}{
		{"", "bool", true, true, "cannot register an unnamed setting", false, false, false, false, false},
		{"cfgbool", "bool", true, true, "", true, false, true, true, false},
		{"cfgbool", "bool", false, true, "cfgbool is already registered, cannot re-register settings", true, false, true, true, false},
		{"", "int", 42, 42, "cannot register an unnamed setting", false, false, false, false, false},
		{"cfgint", "int", 42, 42, "", true, false, true, true, false},
		{"cfgint", "int", 84, 42, "cfgint is already registered, cannot re-register settings", true, false, true, true, false},
		{"", "int64", int64(42), int64(42), "cannot register an unnamed setting", false, false, false, false, false},
		{"cfgint64", "int64", int64(42), int64(42), "", true, false, true, true, false},
		{"cfgint64", "int64", int64(84), int64(42), "cfgint64 is already registered, cannot re-register settings", true, false, true, true, false},
		{"", "string", "bar", "bar", "cannot register an unnamed setting", false, false, false, false, false},
		{"cfgstring", "string", "bar", "bar", "", true, false, true, true, false},
		{"cfgstring", "string", "baz", "bar", "cfgstring is already registered, cannot re-register settings", true, false, true, true, false},
	}
	appCfg = NewCfg("test register")
	var err error
	for i, test := range tests {
		switch test.typ {
		case "bool":
			err = RegisterBoolCfgE(test.name, test.value.(bool))
		case "int":
			err = RegisterIntCfgE(test.name, test.value.(int))
		case "int64":
			err = RegisterInt64CfgE(test.name, test.value.(int64))
		case "string":
			err = RegisterStringCfgE(test.name, test.value.(string))
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
		if IsCfg(test.name) != test.IsCfg {
			t.Errorf("%d expected IsCfg to be %v, got %v", i, test.IsCfg, IsCfg(test.name))
		}
		if IsEnv(test.name) != test.IsEnv {
			t.Errorf("%d expected IsEnv to be %v, got %v", i, test.IsEnv, IsEnv(test.name))
		}
		if IsFlag(test.name) != test.IsFlag {
			t.Errorf("%d expected IsFlag to be %v, got %v", i, test.IsFlag, IsFlag(test.name))
		}
	}
	// Npn=E
	appCfg = NewCfg("test")
	for i, test := range tests {
		// skip empty names since we don't check errors
		if test.name == "" {
			continue
		}
		switch test.typ {
		case "bool":
			RegisterBoolCfg(test.name, test.value.(bool))
		case "int":
			RegisterIntCfg(test.name, test.value.(int))
		case "int64":
			RegisterInt64Cfg(test.name, test.value.(int64))
		case "string":
			RegisterStringCfg(test.name, test.value.(string))
		default:
			t.Errorf("%d: unsupported typ: %s", i, test.typ)
			continue
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
		if IsCfg(test.name) != test.IsCfg {
			t.Errorf("%d expected IsCfg to be %v, got %v", i, test.IsCfg, IsCfg(test.name))
		}
		if IsEnv(test.name) != test.IsEnv {
			t.Errorf("%d expected IsEnv to be %v, got %v", i, test.IsEnv, IsEnv(test.name))
		}
		if IsFlag(test.name) != test.IsFlag {
			t.Errorf("%d expected IsFlag to be %v, got %v", i, test.IsFlag, IsFlag(test.name))
		}
	}

}

func TestRegisterFlagSettings(t *testing.T) {
	tests := []struct {
		name        string
		short       string
		typ         string
		value       interface{}
		expected    interface{}
		expectedErr string
		checkValues bool
		IsCore      bool
		IsCfg       bool
		IsEnv       bool
		IsFlag      bool
	}{
		{"", "", "bool", true, true, "cannot register an unnamed setting", false, false, false, false, false},
		{"flagbool", "b", "bool", true, true, "", true, false, true, true, true},
		{"flagbool", "", "bool", false, true, "flagbool is already registered, cannot re-register settings", true, false, true, true, true},
		{"", "", "int", 42, 42, "cannot register an unnamed setting", false, false, false, false, false},
		{"flagint", "i", "int", 42, 42, "", true, false, true, true, true},
		{"flagint", "", "int", 84, 42, "flagint is already registered, cannot re-register settings", true, false, true, true, true},
		{"", "", "int64", int64(42), int64(42), "cannot register an unnamed setting", false, false, false, false, false},
		{"flagint64", "6", "int64", int64(42), int64(42), "", true, false, true, true, true},
		{"flagint64", "", "int64", int64(84), int64(42), "flagint64 is already registered, cannot re-register settings", true, false, true, true, true},
		{"", "", "string", "bar", "bar", "cannot register an unnamed setting", false, false, false, false, false},
		{"flagstring", "s", "string", "bar", "bar", "", true, false, true, true, true},
		{"flagstring", "", "string", "baz", "bar", "flagstring is already registered, cannot re-register settings", true, false, true, true, true},
	}
	appCfg = NewCfg("test register")
	var err error
	for i, test := range tests {
		switch test.typ {
		case "bool":
			err = RegisterBoolFlagE(test.name, test.short, test.value.(bool), "", "usage")
		case "int":
			err = RegisterIntFlagE(test.name, test.short, test.value.(int), "", "usage")
		case "int64":
			err = RegisterInt64FlagE(test.name, test.short, test.value.(int64), "", "usage")
		case "string":
			err = RegisterStringFlagE(test.name, test.short, test.value.(string), "", "usage")
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
		if IsCfg(test.name) != test.IsCfg {
			t.Errorf("%d expected IsCfg to be %v, got %v", i, test.IsCfg, IsCfg(test.name))
		}
		if IsEnv(test.name) != test.IsEnv {
			t.Errorf("%d expected IsEnv to be %v, got %v", i, test.IsEnv, IsEnv(test.name))
		}
		if IsFlag(test.name) != test.IsFlag {
			t.Errorf("%d expected IsFlag to be %v, got %v", i, test.IsFlag, IsFlag(test.name))
		}
	}
	// Non-E
	appCfg = NewCfg("test register")
	for i, test := range tests {
		// since no error checking is being done, we skip empty names
		switch test.typ {
		case "bool":
			RegisterBoolFlag(test.name, test.short, test.value.(bool), "", "usage")
		case "int":
			RegisterIntFlag(test.name, test.short, test.value.(int), "", "usage")
		case "int64":
			RegisterInt64Flag(test.name, test.short, test.value.(int64), "", "usage")
		case "string":
			RegisterStringFlag(test.name, test.short, test.value.(string), "", "usage")
		default:
			t.Errorf("%d: unsupported typ: %s", i, test.typ)
			continue
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
		if IsCfg(test.name) != test.IsCfg {
			t.Errorf("%d expected IsCfg to be %v, got %v", i, test.IsCfg, IsCfg(test.name))
		}
		if IsEnv(test.name) != test.IsEnv {
			t.Errorf("%d expected IsEnv to be %v, got %v", i, test.IsEnv, IsEnv(test.name))
		}
		if IsFlag(test.name) != test.IsFlag {
			t.Errorf("%d expected IsFlag to be %v, got %v", i, test.IsFlag, IsFlag(test.name))
		}
	}
}
