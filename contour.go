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
	}
	return false
}

func ParseFormatE(s string) (Format, error) {
	ls := strings.ToLower(s)
	switch ls {
	case "json", "jsn", "cjsn", "cjson":
		return JSON, nil
	case "toml", "tml":
		return TOML, nil
	case "yaml", "yml":
		return YAML, nil
	}
	return Unsupported, UnsupportedFormatErr{s}
}

func ParseFormat(s string) Format {
	f, _ := ParseFormatE(s)
	return f
}

// These settings are in order of precedence. Each setting type can be set by
// any of the types with higher precedence if contour is configured to use that
// type.
const (
	// Core settings are immutable once set.
	Core SettingType = iota + 1
	// Env settings can be set from environment variables
	Env
	// File settings can be set from a configuration file.
	File
	// Flag settings can be set from flags.
	Flag
)

// SettingType is type of setting
type SettingType int

func (t SettingType) String() string {
	switch t {
	case Core:
		return "core"
	case File:
		return "file"
	case Env:
		return "env"
	case Flag:
		return "flag"
	default:
		return "unknown"
	}
}

// appCfg: contour's global config; contour config functinos operate on this.
var appCfg *Cfg

func init() {
	appCfg = NewCfg(app)
}

// NotFoundErr occurs when the value was not found.
type NotFoundErr struct {
	v string
}

func (e NotFoundErr) Error() string {
	return fmt.Sprintf("%s: not found", e.v)
}

// SettingNotFoundErr occurs when a setting isn't found.
type SettingNotFoundErr struct {
	settingType SettingType
	name        string
}

func (e SettingNotFoundErr) Error() string {
	if e.settingType <= 0 {
		return fmt.Sprintf("%s: setting not found", e.name)
	}
	return fmt.Sprintf("%s: %s setting not found", e.name, e.settingType)
}

// UnsupportedFormatErr occurs when the string cannot be matched to a
// supported configuration format.
type UnsupportedFormatErr struct {
	v string
}

func (e UnsupportedFormatErr) Error() string {
	return fmt.Sprintf("%s: unsupported configuration format", e.v)
}
