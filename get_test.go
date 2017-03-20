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
	tstSettings := newTestSettings()
	r := tstSettings.Get("corebool")
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
	rb := tstSettings.Bool("corebool")
	if !rb {
		t.Errorf("Expected true, got %t", rb)
	}
	ri := tstSettings.Int("coreint")
	if ri != 42 {
		t.Errorf("Expected 42, got %d", ri)
	}
	ri64 := tstSettings.Int64("coreint64")
	if ri64 != int64(42) {
		t.Errorf("Expected 42, got %d", ri)
	}
	rs := tstSettings.String("corestring")
	if rs != "a core string" {
		t.Errorf("Expected \"a core string\", got %q", rs)
	}
}
