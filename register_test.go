package contour

import (
	"strconv"
	"testing"
)

func TestRegisterCfgFilename(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		format   string
		err      string
	}{
		{"empty", "", "", "A config filename was expected, none received"},
		{"no extension", "cfg", "", "unable to determine config format, the configuration file, cfg, doesn't have an extension"},
		{"toml", "cfg.toml", "toml", ""},
		{"yaml", "cfg.yaml", "yaml", ""},
		{"json", "cfg.json", "json", ""},
		{"xml", "cfg.xml", "xml", "unsupported configuration file format: xml"},
		{"undefined", "cfg.bss", "bss", "unsupported configuration file format: bss"},
	}

	for _, test := range tests {
		cfg := NewCfg(test.name)
		err := cfg.RegisterCfgFilename(CfgFilename, test.filename)
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
		fname, err := cfg.GetStringE(CfgFilename)
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
		if format != test.filename {
			t.Errorf("RegisterCfgFilename format %s: expected %q got %q", test.name, test.format, format)
		}
	}
}

func TestRegisterCoreSettings(t *testing.T) {
	cfg := NewCfg("test register")
	cfg.RegisterBoolCore("corebool", "true")
	b, err := cfg.GetBoolE("corebool")
	if err != nil {
		t.Errorf("%s corebool: unexpected error %q", cfg.name, err)
	} else {
		if b != "true" {
			t.Errorf("%s corebool: expected \"true\" got %q", cfg.name, b)
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
	cfg.RegisterBool("settingbool", "true")
	b, err := cfg.GetBoolE("settingbool")
	if err != nil {
		t.Errorf("%s settingbool: unexpected error %q", cfg.name, err)
	} else {
		bl, err := strconv.ParseBool(b)
		if err != nil {
			t.Errorf("%s settingbool: unexpected error %q", cfg.name, err)

		} else {
			if !bl {
				t.Errorf("%s settingbool: expected true got %q", cfg.name, b)
			}
		}
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
	cfg.RegisterBoolConf("settingbool", "true")
	b, err := cfg.GetBoolE("settingbool")
	if err != nil {
		t.Errorf("%s settingbool: unexpected error %q", cfg.name, err)
	} else {
		bl, err := strconv.ParseBool(b)
		if err != nil {
			t.Errorf("%s settingbool: unexpected error %q", cfg.name, err)

		} else {
			if !bl {
				t.Errorf("%s settingbool: expected true got %q", cfg.name, b)
			}
		}
	}

	cfg.RegisterIntConf("settingint", 42)
	i, err := cfg.GetIntE("settingint")
	if err != nil {
		t.Errorf("%s settingint: unexpected error %q", cfg.name, err)
	} else {
		if i != 42 {
			t.Errorf("%s settingint: expected %d got %d", cfg.name, "42", strconv.Itoa(i))
		}
	}

	cfg.RegisterInt64Conf("settingint64", int64(42))
	i64, err := cfg.GetInt64E("settingint64")
	if err != nil {
		t.Errorf("%s settingint64: unexpected error %q", cfg.name, err)
	} else {
		if i64 != int64(42) {
			t.Errorf("%s settingint64: expected %d got %d", cfg.name, "42", strconv.Itoa(int(i64)))
		}
	}

	cfg.RegisterStringConf("settingstring", "value")
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
	cfg.RegisterBoolFlag("settingbool", "b", "true", "true", "usage")
	b, err := cfg.GetBoolE("settingbool")
	if err != nil {
		t.Errorf("%s settingbool: unexpected error %q", cfg.name, err)
	} else {
		bl, err := strconv.ParseBool(b)
		if err != nil {
			t.Errorf("%s settingbool: unexpected error %q", cfg.name, err)

		} else {
			if !bl {
				t.Errorf("%s settingbool: expected true got %q", cfg.name, b)
			}
		}
	}

	cfg.RegisterIntFlag("settingint", "i", 42, "42", "usage")
	i, err := cfg.GetIntE("settingint")
	if err != nil {
		t.Errorf("%s settingint: unexpected error %q", cfg.name, err)
	} else {
		if i != 42 {
			t.Errorf("%s settingint: expected %d got %d", cfg.name, "42", strconv.Itoa(i))
		}
	}

	cfg.RegisterInt64Flag("settingint64", "x", int64(42), "42", "usage")
	i64, err := cfg.GetInt64E("settingint64")
	if err != nil {
		t.Errorf("%s settingint64: unexpected error %q", cfg.name, err)
	} else {
		if i64 != int64(42) {
			t.Errorf("%s settingint64: expected %d got %d", cfg.name, "42", strconv.Itoa(int(i64)))
		}
	}

	cfg.RegisterStringFlag("settingstring", "s", "value", "value", "usage")
	s, err := cfg.GetStringE("settingstring")
	if err != nil {
		t.Errorf("%s settingstring: unexpected error %q", cfg.name, err)
	} else {
		if s != "value" {
			t.Errorf("%s settingstring: expected \"value\" got %q", cfg.name, b)
		}
	}
}
