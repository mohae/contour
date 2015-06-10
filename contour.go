//
package contour

import (
	"fmt"
	"strings"
)

const app = "app"

const (
	Unsupported Format = iota
	JSON
	TOML
	YAML
	XML
)

type Format int

func (f Format) String() string {
	switch f {
	case JSON:
		return "json"
	case TOML:
		return "toml"
	case YAML:
		return "yaml"
	case XML:
		return "xml"
	default:
		return "unsupported"
	}
}

// isSupported checks to see if the passed string represents a supported config
// format.
func (f Format) isSupported() bool {
	switch f {
	case YAML:
		return true
	case JSON:
		return true
	case TOML:
		return true
	case XML:
		return false
	}
	return false
}

func ParseFormatE(s string) (Format, error) {
	ls := strings.ToLower(s)
	switch ls {
	case "json", "jsn":
		return JSON, nil
	case "toml", "tml":
		return TOML, nil
	case "yaml", "yml":
		return YAML, nil
	case "xml":
		return XML, nil
	}
	return Unsupported, unsupportedFormatErr(ls)
}

func ParseFormat(s string) Format {
	f, _ := ParseFormatE(s)
	return f
}

// appCfg: contour's global config; contour config functinos operate on this.
var appCfg *Cfg

func init() {
	appCfg = NewCfg(app)
}

// notFoundErr returns a standadized notFoundErr.
func notFoundErr(k string) error {
	return fmt.Errorf("not found: %s", k)
}

// settingNotFoundErr standadized settingNotFoundErr.
func settingNotFoundErr(k string) error {
	return fmt.Errorf("setting not found: %s", k)
}

// unsupportedFormatErr standadized the unsupportedFormatErr.
func unsupportedFormatErr(k string) error {
	return fmt.Errorf("unsupported cfg format: %s", k)
}
