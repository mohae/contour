package contour

import (
	"strconv"
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

func TestRegisterCoreSettings(t *testing.T) {
	cfg := NewCfg("test register")
	cfg.RegisterBoolCore("corebool", true)
	b, err := cfg.GetBoolE("corebool")
	if err != nil {
		t.Errorf("%s corebool: unexpected error %q", cfg.name, err)
	} else {
		if b != true {
			t.Errorf("%s corebool: expected \"true\" got %q", cfg.name, strconv.FormatBool(b))
		}
	}

	cfg.RegisterIntCore("coreint", 42)
	i, err := cfg.GetIntE("coreint")
	if err != nil {
		t.Errorf("%s coreint: unexpected error %q", cfg.name, err)
	} else {
		if i != 42 {
			t.Errorf("%s coreint: expected %d got %d", cfg.name, "42", strconv.Itoa(i))
		}
	}

	cfg.RegisterInt64Core("coreint64", int64(42))
	i64, err := cfg.GetInt64E("coreint64")
	if err != nil {
		t.Errorf("%s coreint64: unexpected error %q", cfg.name, err)
	} else {
		if i64 != int64(42) {
			t.Errorf("%s coreint: expected %d got %d", cfg.name, "42", strconv.Itoa(int(i)))
		}
	}

	cfg.RegisterStringCore("corestring", "value")
	s, err := cfg.GetStringE("corestring")
	if err != nil {
		t.Errorf("%s corestring: unexpected error %q", cfg.name, err)
	} else {
		if s != "value" {
			t.Errorf("%s corestring: expected \"value\" got %q", cfg.name, s)
		}
	}

}

func TestRegisterSettings(t *testing.T) {
	cfg := NewCfg("test register")
	cfg.RegisterBool("settingbool", true)
	b, err := cfg.GetBoolE("settingbool")
	if err != nil {
		t.Errorf("%s settingbool: unexpected error %q", cfg.name, err)
	}

	cfg.RegisterInt("settingint", 42)
	i, err := cfg.GetIntE("settingint")
	if err != nil {
		t.Errorf("%s settingint: unexpected error %q", cfg.name, err)
	} else {
		if i != 42 {
			t.Errorf("%s settingint: expected %d got %d", cfg.name, "42", strconv.Itoa(i))
		}
	}

	cfg.RegisterInt64("settingint64", int64(42))
	i64, err := cfg.GetInt64E("settingint64")
	if err != nil {
		t.Errorf("%s settingint64: unexpected error %q", cfg.name, err)
	} else {
		if i64 != int64(42) {
			t.Errorf("%s settingint64: expected %d got %d", cfg.name, "42", strconv.Itoa(int(i64)))
		}
	}

	cfg.RegisterString("settingstring", "value")
	s, err := cfg.GetStringE("settingstring")
	if err != nil {
		t.Errorf("%s settingstring: unexpected error %q", cfg.name, err)
	} else {
		if s != "value" {
			t.Errorf("%s settingstring: expected \"value\" got %q", cfg.name, b)
		}
	}

}

func TestRegisterConfSettings(t *testing.T) {
	cfg := NewCfg("test register Conf")
	cfg.RegisterBoolCfg("settingbool", "rancher_test_setting", true)
	b, err := cfg.GetBoolE("settingbool")
	if err != nil {
		t.Errorf("%s settingbool: unexpected error %q", cfg.name, err)
	} else {
		if !b {
			t.Error("Expected true got false")
		}
	}

	cfg.RegisterIntCfg("settingint", "rancher_test_int", 42)
	i, err := cfg.GetIntE("settingint")
	if err != nil {
		t.Errorf("%s settingint: unexpected error %q", cfg.name, err)
	} else {
		if i != 42 {
			t.Errorf("%s settingint: expected %d got %d", cfg.name, "42", strconv.Itoa(i))
		}
	}

	cfg.RegisterInt64Cfg("settingint64", "rancher_test_int64", int64(42))
	i64, err := cfg.GetInt64E("settingint64")
	if err != nil {
		t.Errorf("%s settingint64: unexpected error %q", cfg.name, err)
	} else {
		if i64 != int64(42) {
			t.Errorf("%s settingint64: expected %d got %d", cfg.name, "42", strconv.Itoa(int(i64)))
		}
	}

	cfg.RegisterStringCfg("settingstring", "rancher_test_string", "value")
	s, err := cfg.GetStringE("settingstring")
	if err != nil {
		t.Errorf("%s settingstring: unexpected error %q", cfg.name, err)
	} else {
		if s != "value" {
			t.Errorf("%s settingstring: expected \"value\" got %q", cfg.name, b)
		}
	}

}

func TestRegisterFlagSettings(t *testing.T) {
	cfg := NewCfg("test register flag")
	cfg.RegisterBoolFlag("settingbool", "b", "rancher_test_bool_flag", true, "true", "usage")
	b, err := cfg.GetBoolE("settingbool")
	if err != nil {
		t.Errorf("%s settingbool: unexpected error %q", cfg.name, err)
	} else {
		if !b {
			t.Errorf("%s settingbool: expected true got false", cfg.name)
		}
	}

	cfg.RegisterIntFlag("settingint", "i", "rancher_test_int_flag", 42, "42", "usage")
	i, err := cfg.GetIntE("settingint")
	if err != nil {
		t.Errorf("%s settingint: unexpected error %q", cfg.name, err)
	} else {
		if i != 42 {
			t.Errorf("%s settingint: expected %d got %d", cfg.name, "42", strconv.Itoa(i))
		}
	}

	cfg.RegisterInt64Flag("settingint64", "x", "rancher_test_int64_flag", int64(42), "42", "usage")
	i64, err := cfg.GetInt64E("settingint64")
	if err != nil {
		t.Errorf("%s settingint64: unexpected error %q", cfg.name, err)
	} else {
		if i64 != int64(42) {
			t.Errorf("%s settingint64: expected %d got %d", cfg.name, "42", strconv.Itoa(int(i64)))
		}
	}

	cfg.RegisterStringFlag("settingstring", "s", "rancher_test_string_flag", "value", "value", "usage")
	s, err := cfg.GetStringE("settingstring")
	if err != nil {
		t.Errorf("%s settingstring: unexpected error %q", cfg.name, err)
	} else {
		if s != "value" {
			t.Errorf("%s settingstring: expected \"value\" got %q", cfg.name, b)
		}
	}
}
