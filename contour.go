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

// ParseFormat takes a string and returns the Format it represents or an
// UnsupportedFormatErr if it can't be matched to a format. The string is
// normalized to lower case before matching.
func ParseFormat(s string) (Format, error) {
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

const (
	_interface dataType = iota + 1
	_bool
	_int
	_int64
	_string
)

// dataType is the setting's data type.
type dataType int

func (t dataType) String() string {
	switch t {
	case _string:
		return "string"
	case _int:
		return "int"
	case _int64:
		return "int64"
	case _bool:
		return "bool"
	case _interface:
		return "interface{}"
	}
	return "unknown data type"
}

func parseDataType(s string) (dataType, error) {
	v := strings.ToLower(s)
	switch v {
	case "string":
		return _string, nil
	case "int":
		return _int, nil
	case "int64":
		return _int64, nil
	case "bool":
		return _bool, nil
	}
	return 0, fmt.Errorf("%s: not a supported data type", s)
}

// DataTypeErr occurs when the requested setting's data type is different than
// the type requested.
type DataTypeErr struct {
	name string
	is   string
	not  dataType
}

func (e DataTypeErr) Error() string {
	return fmt.Sprintf("%s is %s, not %s", e.name, e.is, e.not)
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

// settings: contour's global Settinngs set; contour functions operate on this.
var settings *Settings

func init() {
	settings = New(app)
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
