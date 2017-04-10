package contour

import (
	"testing"
)

func TestGetsE(t *testing.T) {
	tstSettings := newTestSettings()
	r, err := tstSettings.GetE("corebool")
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
	rb, err := tstSettings.BoolE("corebool")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err.Error())
	} else {
		if !rb {
			t.Errorf("Expected \"true\", got %t", rb)
		}
	}
	ri, err := tstSettings.IntE("coreint")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err.Error())
	} else {
		if ri != 42 {
			t.Errorf("Expected 42, got %d", ri)
		}
	}
	ri64, err := tstSettings.Int64E("coreint64")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err.Error())
	} else {
		if ri64 != int64(42) {
			t.Errorf("Expected 42, got %d", ri)
		}
	}
	rinf, err := tstSettings.InterfaceE("coreslice")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err.Error())
	} else {
		inf, ok := rinf.([]string)
		if !ok {
			t.Errorf("coreinterface: assertion to []string failed: %#v", rinf)
		} else {
			if len(inf) != 0 {
				t.Errorf("Expected slice to have 0 len; was %d: %v", len(inf), inf)
			}
		}
	}
	rs, err := tstSettings.StringE("corestring")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err.Error())
	} else {
		if rs != "a core string" {
			t.Errorf("Expected \"a core string\", got %q", rs)
		}
	}
}

func TestGets(t *testing.T) {
	r := Get("corebool")
	var b bool
	switch r.(type) {
	case bool:
		b = r.(bool)
	case *bool:
		b = *r.(*bool)
	}
	if !b {
		t.Errorf("Expected \"true\", got %v", r)
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
	rinf := Interface("coreslice")
	inf, ok := rinf.([]string)
	if !ok {
		t.Errorf("coreinterface: assertion of interface{} to []string failed: got %#v", rinf)
	} else {
		if len(inf) != 0 {
			t.Errorf("Expected len to be 0; got %d: %v", len(inf), inf)
		}
	}
	rs := String("corestring")
	if rs != "a core string" {
		t.Errorf("Expected \"a core string\", got %q", rs)
	}
}
