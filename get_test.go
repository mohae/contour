package contour

import (
	"testing"
)

func TestGetsE(t *testing.T) {
	appCfg = newTestCfg()
	r, err := GetE("corebool")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err.Error())
	} else {
		var b bool
		switch r.(type) {
		case bool:
			b = r.(bool)
		case *bool:
			b = *r.(*bool)
		}
		if !b {
			t.Errorf("Expected \"true\", got %t", b)
		}
	}
	rb, err := BoolE("corebool")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err.Error())
	} else {
		if !rb {
			t.Errorf("Expected \"true\", got %t", rb)
		}
	}
	ri, err := IntE("coreint")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err.Error())
	} else {
		if ri != 42 {
			t.Errorf("Expected 42, got %d", ri)
		}
	}
	ri64, err := Int64E("coreint64")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err.Error())
	} else {
		if ri64 != int64(42) {
			t.Errorf("Expected 42, got %d", ri)
		}
	}
	rs, err := StringE("corestring")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err.Error())
	} else {
		if rs != "a core string" {
			t.Errorf("Expected \"a core string\", got %q", rs)
		}
	}
}

func TestGets(t *testing.T) {
	appCfg = newTestCfg()
	r := Get("corebool")
	var b bool
	switch r.(type) {
	case bool:
		b = r.(bool)
	case *bool:
		b = *r.(*bool)
	}
	if !b {
		t.Errorf("Expected \"true\", got %t", r)
	}
	rb := Bool("corebool")
	if !rb {
		t.Errorf("Expected true, got %t", rb)
	}
	ri := Int("coreint")
	if ri != 42 {
		t.Errorf("Expected 42, got %d", ri)
	}
	ri64 := Int64("coreint64")
	if ri64 != int64(42) {
		t.Errorf("Expected 42, got %d", ri)
	}
	rs := String("corestring")
	if rs != "a core string" {
		t.Errorf("Expected \"a core string\", got %q", rs)
	}
}
