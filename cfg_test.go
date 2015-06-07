package contour

import (
	"fmt"
	"os"
	"testing"
)

func TestCfg(t *testing.T) {
	c := NewCfg("test")
	if c == nil {
		t.Errorf("New test cfg was nil")
	} else {
		if c.name != "test" {
			t.Errorf("Expected test got %s", c.name)
		}
		if c.UseEnv() != true {
			t.Errorf("Expected true got %v", c.UseEnv())
		}
	}
	a := AppCfg()
	if a == nil {
		t.Errorf("New test cfg was nil")
	} else {
		if a.name != "app" {
			t.Errorf("Expected app got %s", a.name)
		}
		if c.UseEnv() != true {
			t.Errorf("Expected true got %v", c.UseEnv())
		}
	}

}

func TestLoadEnv(t *testing.T) {
	tests := []struct {
		name        string
		envValue    string
		typ         string
		origValue   interface{}
		expValue    interface{}
		expectedErr string
	}{
		{"tcfgbool", "true", "cfgbool", false, true, ""},
		{"tcfgint", "99", "cfgint", 42, 99, ""},
		{"tcfgstring", "bar", "cfgstring", "foo", "bar", ""},
		{"tflagbool", "true", "fkagbool", false, true, ""},
		{"tflagint", "88", "flagint", 42, 88, ""},
		{"tflagstring", "biz", "flagstring", "fiz", "biz", ""},
	}
	tCfg := NewCfg("test")
	for _, test := range tests {
		switch test.typ {
		case "cfgbool":
			tCfg.RegisterBoolCfg(test.name, test.origValue.(bool))
		case "cfgint":
			tCfg.RegisterIntCfg(test.name, test.origValue.(int))
		case "cfgstring":
			tCfg.RegisterStringCfg(test.name, test.origValue.(string))
		case "flagbool":
			tCfg.RegisterBoolFlag(test.name, "", test.origValue.(bool), "", "")
		case "flagint":
			tCfg.RegisterIntFlag(test.name, "", test.origValue.(int), "", "")
		case "flagstring":
			tCfg.RegisterStringFlag(test.name, "", test.origValue.(string), "", "")
		}
		os.Setenv(fmt.Sprintf("%s_%s", tCfg.Name(), test.name), test.envValue)
	}
	tCfg.UpdateFromEnv()
	for _, test := range tests {
		tmp := tCfg.Get(test.name)
		switch test.typ {
		case "cfgbool", "flagbool":
			if test.expValue != tmp.(bool) {
				t.Errorf("expected %v, got %v", test.expValue, tmp)
			}
		case "cfgint", "flagint":
			if test.expValue != tmp.(int) {
				t.Errorf("expected %v, got %v", test.expValue, tmp)
			}
		case "cfgstring", "flagstring":
			if test.expValue != tmp.(string) {
				t.Errorf("expected %v, got %v", test.expValue, tmp)
			}
		}
	}
}

func TestCfgBools(t *testing.T) {
	bTests := []struct {
		val      bool
		expected bool
	}{
		{true, true},
		{false, false},
		{true, true},
	}
	tCfg := NewCfg("test")
	for _, test := range bTests {
		tCfg.SetErrOnMissingCfg(test.val)
		b := tCfg.ErrOnMissingCfg()
		if b != test.expected {
			t.Errorf("ErrOnMissingCfg:  expected %v, got %v", test.expected, b)
		}
		tCfg.SetSearchPath(test.val)
		b = tCfg.SearchPath()
		if b != test.expected {
			t.Errorf("SearchPath:  expected %v, got %v", test.expected, b)
		}
		tCfg.SetUseCfg(test.val)
		b = tCfg.UseCfg()
		if b != test.expected {
			t.Errorf("SetUseCfgFile:  expected %v, got %v", test.expected, b)
		}
		tCfg.SetUseEnv(test.val)
		b = tCfg.UseEnv()
		if b != test.expected {
			t.Errorf("SetUseEnv:  expected %v, got %v", test.expected, b)
		}
	}

}

//func TestSetFromCfg(t *testing.T) {
//
//}

func TestCfgProcessed(t *testing.T) {
	tests := []struct {
		useCfg       bool
		cfgSet       bool
		useEnv       bool
		envSet       bool
		useFlags     bool
		argsFiltered bool
		expected     bool
	}{
		// 0
		{false, false, false, false, false, false, true},
		{false, false, false, false, false, true, true},
		{false, false, false, false, true, false, false},
		{false, false, false, false, true, true, true},
		{false, false, false, true, false, false, true},
		// 5
		{false, false, false, true, false, true, true},
		{false, false, false, true, true, false, false},
		{false, false, false, true, true, true, true},
		{false, false, true, false, false, false, false},
		{false, false, true, false, false, true, false},
		// 10
		{false, false, true, false, true, false, false},
		{false, false, true, false, true, true, false},
		{false, false, true, true, false, false, true},
		{false, false, true, true, true, false, false},
		{false, false, true, true, false, true, true},
		// 15
		{false, false, true, true, true, true, true},
		{false, true, false, false, false, false, true},
		{false, true, false, false, false, true, true},
		{false, true, false, false, true, false, false},
		{false, true, false, false, true, true, true},
		{false, true, false, true, false, false, true},
		// 20
		{false, true, false, true, false, true, true},
		{false, true, false, true, true, false, false},
		{false, true, false, true, true, true, true},
		{false, true, true, false, false, false, false},
		{false, true, true, false, false, true, false},
		// 25
		{false, true, true, false, true, false, false},
		{false, true, true, false, true, true, false},
		{false, true, true, true, false, false, true},
		{false, true, true, true, true, false, false},
		{false, true, true, true, false, true, true},
		// 30
		{false, true, true, true, true, true, true},
		{true, false, false, false, false, false, false},
		{true, false, false, false, false, true, false},
		{true, false, false, false, true, false, false},
		{true, false, false, false, true, true, false},
		// 35
		{true, false, false, true, false, false, false},
		{true, false, false, true, false, true, false},
		{true, false, false, true, true, false, false},
		{true, false, false, true, true, true, false},
		{true, false, true, false, false, false, false},
		// 40
		{true, false, true, false, false, true, false},
		{true, false, true, false, true, false, false},
		{true, false, true, false, true, true, false},
		{true, false, true, true, false, false, false},
		{true, false, true, true, true, false, false},
		// 45
		{true, false, true, true, false, true, false},
		{true, false, true, true, true, true, false},
		{true, true, false, false, false, false, true},
		{true, true, false, false, false, true, true},
		{true, true, false, false, true, false, false},
		// 50
		{true, true, false, false, true, true, true},
		{true, true, false, true, false, false, true},
		{true, true, false, true, false, true, true},
		{true, true, false, true, true, false, false},
		{true, true, false, true, true, true, true},
		// 55
		{true, true, true, false, false, false, false},
		{true, true, true, false, false, true, false},
		{true, true, true, false, true, false, false},
		{true, true, true, false, true, true, false},
		{true, true, true, true, false, false, true},
		// 60
		{true, true, true, true, true, false, false},
		{true, true, true, true, false, true, true},
		{true, true, true, true, true, true, true},
	}
	tCfg := NewCfg("test")
	for i, test := range tests {
		tCfg.SetUseCfg(test.useCfg)
		tCfg.cfgSet = test.cfgSet
		tCfg.envSet = test.envSet
		tCfg.SetUseEnv(test.useEnv)
		tCfg.useFlags = test.useFlags
		tCfg.argsFiltered = test.argsFiltered
		b := tCfg.CfgProcessed()
		if b != test.expected {
			t.Errorf("%d expected %v, got %v", i, test.expected, b)
		}
	}
}

func TestCanUpdate(t *testing.T) {
	tests := []struct {
		name         string
		argsFiltered bool
		expected     bool
		err          string
	}{
		// 0
		{"corebool", false, false, "cannot update \"corebool\": core settings cannot be updated"},
		{"x-corebool", false, false, "cannot update \"x-corebool\": setting not found"},
		{"coreint", false, false, "cannot update \"coreint\": core settings cannot be updated"},
		{"x-coreint", false, false, "cannot update \"x-coreint\": setting not found"},
		{"coreint64", false, false, "cannot update \"coreint64\": core settings cannot be updated"},
		// 5
		{"x-coreint64", false, false, "cannot update \"x-coreint64\": setting not found"},
		{"corestring", false, false, "cannot update \"corestring\": core settings cannot be updated"},
		{"x-corestring", false, false, "cannot update \"x-corestring\": setting not found"},
		{"cfgbool", false, true, ""},
		{"x-cfgbool", false, false, "cannot update \"x-cfgbool\": setting not found"},
		// 10
		{"cfgint", false, true, ""},
		{"x-cfgint", false, false, "cannot update \"x-cfgint\": setting not found"},
		{"cfgint64", false, true, ""},
		{"x-cfgint64", false, false, "cannot update \"x-cfgint64\": setting not found"},
		{"cfgstring", false, true, ""},
		// 15
		{"x-cfgstring", false, false, "cannot update \"x-cfgstring\": setting not found"},
		{"flagbool", false, true, ""},
		{"x-flagbool", false, false, "cannot update \"x-flagbool\": setting not found"},
		{"flagint", false, true, ""},
		{"x-flagint", false, false, "cannot update \"x-flagint\": setting not found"},
		// 20
		{"flagint64", false, true, ""},
		{"x-flagint64", false, false, "cannot update \"x-flagint64\": setting not found"},
		{"flagstring", false, true, ""},
		{"x-flagstring", false, false, "cannot update \"x-flagstring\": setting not found"},
		{"flagbool", true, false, "cannot update \"flagbool\": flag settings cannot be updated after arg filtering"},
		// 25
		{"x-flagbool", true, false, "cannot update \"x-flagbool\": setting not found"},
		{"flagint", true, false, "cannot update \"flagint\": flag settings cannot be updated after arg filtering"},
		{"x-flagint", true, false, "cannot update \"x-flagint\": setting not found"},
		{"flagint64", true, false, "cannot update \"flagint64\": flag settings cannot be updated after arg filtering"},
		{"x-flagint64", true, false, "cannot update \"x-flagint64\": setting not found"},
		// 30
		{"flagstring", true, false, "cannot update \"flagstring\": flag settings cannot be updated after arg filtering"},
		{"x-flagstring", true, false, "cannot update \"x-flagstring\": setting not found"},
		{"bool", false, true, ""},
		{"x-bool", false, false, "cannot update \"x-bool\": setting not found"},
		{"int", false, true, ""},
		// 35
		{"x-int", false, false, "cannot update \"x-int\": setting not found"},
		{"int64", false, true, ""},
		{"x-int64", false, false, "cannot update \"x-int64\": setting not found"},
		{"string", false, true, ""},
		{"x-string", false, false, "cannot update \"x-string\": setting not found"},
	}
	tCfg := newTestCfg()
	for i, test := range tests {
		tCfg.argsFiltered = test.argsFiltered
		b, err := tCfg.canUpdate(test.name)
		if err != nil {
			if err.Error() != test.err {
				t.Errorf("%d expected %q got %q", i, test.err, err.Error())
			}
			if b {
				t.Errorf("%d: expected returned value to be false on an error, it was not false", i)
			}
			continue
		}
		if b != test.expected {
			t.Errorf("%d: expected %v got %v", i, test.expected, b)
		}
	}
}
