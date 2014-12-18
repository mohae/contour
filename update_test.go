package contour

import (
	"testing"
)

func TestUpdates(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
		err   string
	}{
		{"false", "corebool", "false", "corebool is not updateable"},
		{"true", "corebool", "t", "corebool is not updateable"},
		{"unset", "", "", ""},
		{"false", "flagbool", "false", ""},
		{"true", "flagbool", "t", ""},
		{"false", "confbool", "false", ""},
		{"true", "confbool", "t", ""},
		{"false", "settingbool", "false", ""},
		{"true", "settingbool", "t", ""},
	}

	cfg := testCfg
	for _, test := range tests {
		err := cfg.UpdateBoolE(test.key, test.value)
		if err != nil {
			if test.err == "" {
				t.Errorf("%s: unexpected error %q", test.name, err)
				goto cont
			}
			if test.err != err.Error() {
				t.Errorf("%s: expected %q got %q", test.name, test.err, err)
			}
		cont:
			continue
		}
		b, err := cfg.GetBoolE(test.key)
		if err != nil {
			if test.err == "" {
				t.Errorf("%s: unexpected error %q", test.name, err)
			} else {
				if test.err != err.Error() {
					t.Errorf("%s: expected %q got %q", test.name, test.err, err)
				}
			}
			continue
		}
		if b != test.value {
			t.Errorf("%s: expected %q got %q", test.name, test.value, b)
		}
	}
}
