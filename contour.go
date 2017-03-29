// Copyright 2016, Joel Scoble. All rights reserved.
// Licensed under the MIT License. See the included LICENSE file.

// Package contour: a package for settings.
//
// Contour supports application settings, loading settings from a configuration
// file, environment variables, and flags. Where a setting can be set from is
// configurable at the per setting level. In addition to setting's
// overridability being configurable, the Settings behavior is also
// configurable.
//
// Application settings can either be core settings, not updateable once set,
// or they can be settings that can only be updated within the application
// using the Update methods; these settings are not exposed to configuration
// files, environment variables, or flags and cannot be modified by them.
//
// For flags, short flag aliases are also supported.
//
// The package global Settings uses the application's name as its name.
//
// All operations are thread-safe.
package contour

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mohae/appname"
)

var app = appname.Get()

const (
	// Unsupported configuration encoding format.
	Unsupported Format = iota
	// JSON encoding format.
	JSON
	// TOML encoding format.
	TOML
	// YAML encoding format
	YAML
)

// Format is the type of esupported encoding for configuration files.
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
// UnsupportedFormatErr if it can't be matched to a supported format. The
// string is normalized to lower case before matching.
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

// ParseFilenameFormat takes a string that represents a filename and returns
// the files format based on its extension. If the filename either doesn't have
// an extension or the extension is not one of a supported file format an
// UnsupportedFormatErr will be returned.
func ParseFilenameFormat(s string) (Format, error) {
	ext := strings.TrimPrefix(filepath.Ext(s), ".")
	return ParseFormat(ext)
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

func parseDataType(s string) dataType {
	v := strings.ToLower(s)
	switch v {
	case "string":
		return _string
	case "int":
		return _int
	case "int64":
		return _int64
	case "bool":
		return _bool
	}
	// everything else is an interface{}, the user of the setting will be
	// expected to know what it is.
	return _interface
}

// DataTypeErr occurs when the requested setting's data type is different than
// the type requested.
type DataTypeErr struct {
	k   string
	is  string
	not dataType
}

func (e DataTypeErr) Error() string {
	return fmt.Sprintf("%s is %s, not %s", e.k, e.is, e.not)
}

// These settings are in order of precedence. Each setting type can be set by
// any of the types with higher precedence if contour is configured to use that
// type.
const (
	// Basic settings are settings that are none of the below. These are often
	// referred to as application settings: settings that can only be updated
	// within an application and not by configuration files, environment
	// variables, or flags. These settings do not have to be registered.
	Basic SettingType = iota + 1
	// Core settings are immutable once registered.
	Core
	// ConFileVar settings can be set from a configuration file.
	ConfFileVar
	// EnvVar settings can be set from a configuration file and an environment
	// variable; unless it has been explicitly set to not be updateable from a
	// configuration file.
	EnvVar
	// Flag settings can be set from a configuration file, an environment
	// variable, and a flag; unless it has been explicitly set to not be
	// updateable from either a configuration file or an environment variable.
	Flag
)

// SettingType is type of setting.
type SettingType int

func (t SettingType) String() string {
	switch t {
	case Basic:
		return "basic"
	case Core:
		return "core"
	case ConfFileVar:
		return "configuration file var"
	case EnvVar:
		return "env var"
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

var ErrNoSettingName = errors.New("no setting name provided")

// SettingExistsErr occurs when a setting being Added or Registered already
// exists under the same name (k).
type SettingExistsErr struct {
	typ SettingType
	k   string
}

func (e SettingExistsErr) Error() string {
	// if typ is unknown or basic, don't include it in the o utput.
	if e.typ <= 1 {
		return fmt.Sprintf("%s: setting exists", e.k)
	}
	return fmt.Sprintf("%s: %s setting exists", e.k, e.typ)
}

// ShortFlagExistsErr occurs when registering a flag whose short flag already
// exists/
type ShortFlagExistsErr struct {
	k         string
	short     string
	shortName string
}

func (e ShortFlagExistsErr) Error() string {
	return fmt.Sprintf("%s: short flag %q already exists for %q", e.k, e.short, e.shortName)
}

// SettingNotFoundErr occurs when a setting isn't found.
type SettingNotFoundErr struct {
	settingType SettingType
	k           string
}

func (e SettingNotFoundErr) Error() string {
	if e.settingType <= 0 {
		return fmt.Sprintf("%s: setting not found", e.k)
	}
	return fmt.Sprintf("%s: %s setting not found", e.k, e.settingType)
}

// UnsupportedFormatErr occurs when the string cannot be matched to a
// supported configuration format.
type UnsupportedFormatErr struct {
	v string
}

func (e UnsupportedFormatErr) Error() string {
	return fmt.Sprintf("%s: unsupported configuration format", e.v)
}
