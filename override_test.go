package contour

import (
	"testing"
)

func TestOverride(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		value       interface{}
		expectedErr string
	}{
		{"corebool", "corebool", true, "corebool is not a flag: only flags can be overridden"},
		{"coreint", "coreint", 42, "coreint is not a flag: only flags can be overridden"},
		{"corestring", "corestring", "beeblebrox", "corestring is not a flag: only flags can be overridden"},
		{"cfgbool", "cfgbool", true, "cfgbool is not a flag: only flags can be overridden"},
		{"cfgint", "cfgint", 43, "cfgint is not a flag: only flags can be overridden"},
		{"cfgstring", "cfgstring", "frood", "cfgstring is not a flag: only flags can be overridden"},
		{"flagbool", "flagbool", true, ""},
		{"flagint", "flagint", 41, ""},
		{"flagstring", "flagstring", "towel", ""},
		{"bool", "bool", true, "bool is not a flag: only flags can be overridden"},
		{"int", "int", 3, "int is not a flag: only flags can be overridden"},
		{"string", "string", "don't panic", "string is not a flag: only flags can be overridden"},
	}
	testCfg := newTestSettings()
	for i, test := range tests {
		err := testCfg.Override(test.key, test.value)
		if err != nil {
			if err.Error() != test.expectedErr {
				t.Errorf("%d: expected error to be %q, got %q", i, test.expectedErr, err)
			}
			continue
		}
		if test.expectedErr != "" {
			t.Errorf("%d: expected error to be %q: got none", i, test.expectedErr)
			continue
		}
		v := testCfg.Get(test.key)
		if v != test.value {
			t.Errorf("%d: expected %v got %v", i, test.value, v)
		}
	}
}
