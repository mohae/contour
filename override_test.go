package contour

import (
	"testing"
)

func TestOverride(t *testing.T) {
	tests := []struct {
		key         string
		value       interface{}
		expectedErr string
	}{
		{"corebool", true, "corebool is not a flag: only flags can be overridden"},
		{"coreint", 42, "coreint is not a flag: only flags can be overridden"},
		{"corestring", "beeblebrox", "corestring is not a flag: only flags can be overridden"},
		{"cfgbool", true, "cfgbool is not a flag: only flags can be overridden"},
		{"cfgint", 43, "cfgint is not a flag: only flags can be overridden"},
		{"cfgstring", "frood", "cfgstring is not a flag: only flags can be overridden"},
		{"flagbool", true, ""},
		{"flagint", 41, ""},
		{"flagstring", "towel", ""},
		{"bool", true, "bool is not a flag: only flags can be overridden"},
		{"int", 3, "int is not a flag: only flags can be overridden"},
		{"string", "don't panic", "string is not a flag: only flags can be overridden"},
	}
	testCfg := newTestSettings()
	for _, test := range tests {
		err := testCfg.Override(test.key, test.value)
		if err != nil {
			if err.Error() != test.expectedErr {
				t.Errorf("%s: expected error to be %q, got %q", test.key, test.expectedErr, err)
			}
			continue
		}
		if test.expectedErr != "" {
			t.Errorf("%s: expected error to be %q: got none", test.key, test.expectedErr)
			continue
		}
		v := testCfg.Get(test.key)
		if v != test.value {
			t.Errorf("%s: expected %v got %v", test.key, test.value, v)
		}
	}
}
