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
	rb, err := GetBoolE("corebool")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err.Error())
	} else {
		if !rb {
			t.Errorf("Expected \"true\", got %t", rb)
		}
	}
	ri, err := GetIntE("coreint")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err.Error())
	} else {
		if ri != 42 {
			t.Errorf("Expected 42, got %d", ri)
		}
	}
	ri64, err := GetInt64E("coreint64")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err.Error())
	} else {
		if ri64 != int64(42) {
			t.Errorf("Expected 42, got %d", ri)
		}
	}
	rs, err := GetStringE("corestring")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err.Error())
	} else {
		if rs != "a core string" {
			t.Errorf("Expected \"a core string\", got %q", rs)
		}
	}
	rif, err := GetInterfaceE("corebool")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err.Error())
	} else {
		var b bool
		switch rif.(type) {
		case bool:
			b = rif.(bool)
		case *bool:
			b = *rif.(*bool)
		}
		if !b {
			t.Errorf("Expected \"true\", got %t", b)
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
	rb := GetBool("corebool")
	if !rb {
		t.Errorf("Expected true, got %t", rb)
	}
	ri := GetInt("coreint")
	if ri != 42 {
		t.Errorf("Expected 42, got %d", ri)
	}
	ri64 := GetInt64("coreint64")
	if ri64 != int64(42) {
		t.Errorf("Expected 42, got %d", ri)
	}
	rs := GetString("corestring")
	if rs != "a core string" {
		t.Errorf("Expected \"a core string\", got %q", rs)
	}
	rif := GetInterface("corebool")
	switch rif.(type) {
	case bool:
		b = rif.(bool)
	case *bool:
		b = *rif.(*bool)
	}
	if !b {
		t.Errorf("Expected true, got %t", b)
	}
}
