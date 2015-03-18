package contour

import (
	"testing"
)

func TestGetsE(t *testing.T) {
	r, err := testCfg.GetE("corebool")
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

	rb, err := testCfg.GetBoolE("corebool")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err.Error())
	} else {
		if !rb {
			t.Errorf("Expected \"true\", got %t", rb)
		}
	}

	ri, err := testCfg.GetIntE("coreint")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err.Error())
	} else {
		if ri != 42 {
			t.Errorf("Expected 42, got %d", ri)
		}
	}

	rs, err := testCfg.GetStringE("corestring")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err.Error())
	} else {
		if rs != "a core string" {
			t.Errorf("Expected \"a core string\", got %q", rs)
		}
	}

	rif, err := testCfg.GetInterfaceE("corebool")
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
	r := testCfg.Get("corebool")
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

	rb := testCfg.GetBool("corebool")
	if !rb {
		t.Errorf("Expected true, got %t", rb)
	}

	ri := testCfg.GetInt("coreint")
	if ri != 42 {
		t.Errorf("Expected 42, got %d", ri)
	}

	rs := testCfg.GetString("corestring")
	if rs != "a core string" {
		t.Errorf("Expected \"a core string\", got %q", rs)
	}

	rif := testCfg.GetInterface("corebool")
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

func TestGetFilterNames(t *testing.T) {
	boolFilters := testCfg.GetBoolFilterNames()
	if toString.Get(boolFilters) != "[\"flagbool\"]" {
		t.Errorf("Expected [\"flagbool\"], got %s", toString.Get(boolFilters))
	}

	intFilters := testCfg.GetIntFilterNames()
	if toString.Get(intFilters) != "[\"flagint\"]" {
		t.Errorf("Expected [\"flagint\"], got %s", toString.Get(intFilters))
	}

	stringFilters := testCfg.GetStringFilterNames()
	if toString.Get(stringFilters) != "[\"flagstring\"]" {
		t.Errorf("Expected [\"flagstring\"], got %s", toString.Get(stringFilters))
	}
}
