// Copyright 2016, Joel Scoble. All rights reserved.
// Licensed under the MIT License. See the included LICENSE file.

// Package contour: a package for storing settings. These settings may be
// configuration settings or application settings.
//
// Application settings are either Core or Basic settings. Core settings cannot
// be modified once Added with AddCore functions. Basic settings are set by
// Add functions and can be updated with Update functions. These settings are
// not exposed as configuration file variables, environment variables, or
// flags. They cannot be modified by any of them.
//
// Configuration settings are settings that are updateable by one or more of
// the following, depending on what type of configuration setting they are, in
// order of override precedence: configuration file variable, environment
// variable, and flag. Configuration settings are registered. They are
// registered with default values and are overridable according to their
// configuration setting type. For custom override properties, e.g. can be
// set by either a configuration file or a flag but not by an environment
// variable, use the Register function.
//
// The configuration file, and the format that it is in, can be specified. The
// supported formats are: JSON, TOML, and YAML. A Contour Settings can also
// be configured to search the PATH for the configuration file. For situations
// where the configuration file is optional, contour can be set to not generate
// an error when it cannot be found.
//
// Contour only saves the top level keys of configuration files as settings.
// For configuration file settings that are arrays, maps, or objects, their
// values will be saved as an interface{}.
//
// Environment variables are UPPER CASE and use a NAME_KEY as the variable
// name, where NAME is the name of the Settings, the executable name for
// the package global Settings, and KEY is the name, or key, of the setting.
//
// Flags can be registered with either a short flag or alias using the short
// parameter of Register Flag functions.
//
// All operations are thread-safe.
//
// The workflow for application configurations is:
//    Register all configuration settings.
//    Call the Set() function.
//    Call the ParseFlags() function.
//
// Non-configuration application settings, Core and Basic, can be added at
// anytime.
package contour

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mohae/appname"
)

// Exe: the name of the running executable.
var Exe = appname.Get()

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
	settings = New(Exe)
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
