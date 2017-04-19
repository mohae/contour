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
// Registering a configuration setting will result in settings being configured
// to use that setting type's source, along with any lower precedence sources
// during the Set process, e.g. registering a ConfFileVar will result in
// settings being configured to use a configuration file and registering a Flag
// will result in settings being configured to use a configuration file and
// check environment variables.
//
// After registration, settings can be set to ignore certain configuration
// sources using the SetUseConfFile and SetUseEnvVars methods. Flag parsing is
// always explicitly done by the caller with the ParseFlags method.
//
// The configuration file can be explicitly set, in which case the
// configuration file format will be inferred from the extension with unknown
// extensions resulting in an UnsupportedFormatError. If there are
// configuration settings registered, of any type, and the configuration file
// has not been set, it will be assumed to be SettingsName.Format where
// SettingsName is the name of the settings and Format is the configuration
// file format that settings is set to use, which defaults to JSON. The format
// can be set using the SetFormat method. This only needs to be done if the
// configuration file is not explictly set using the SetConfFilename method.
// The supported configuration formats are: JSON, TOML, and YAML. A settings
// can also be configured to search for the configuration file until it is
// found. Where it looks depends on how it has been configured and what
// additional information the settings has been provided:
//    configuration filename
//    paths set with SetConfFilePaths
//    paths extracted from env vars set by SetConfFilePathEnvVars*
//    working directory
//    executable directory
//    paths extracted from the PATH*
//
//    * the env vars may contain multiple paths; each path will be checked
//
// By default, a missing configuration file results in an os.PathError with
// a list of all paths that were checked along with an os.IsNotExist error. A
// settings can be set to not return an error when the configuration file
// cannot be found by using the SetErrOnMissingConfFile method.
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
//
// For convenience, there's a predefined 'standard' Settings, whose name is the
// executable's name.
package contour

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mohae/appname"
)

var (
	Exe    = appname.Get() // Exe is the name of the running executable.
	std    *Settings       // std: contour's global Settinngs set; contour functions operate on this.
	format = JSON          // format: the default format
)

func init() {
	std = New(Exe)
}

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
// UnsupportedFormatError if it can't be matched to a supported format. The
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
	return Unsupported, UnsupportedFormatError{s}
}

// ParseFilenameFormat takes a string that represents a filename and returns
// the files format based on its extension. If the filename either doesn't have
// an extension or the extension is not one of a supported file format an
// UnsupportedFormatError will be returned.
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

// DataTypeError occurs when the requested setting's data type is different
// than the type requested.
type DataTypeError struct {
	k   string
	is  string
	not dataType
}

func (e DataTypeError) Error() string {
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

var ErrNoSettingName = errors.New("no setting name provided")

// SettingExistsError occurs when a setting being Added or Registered already
// exists under the same name (k).
type SettingExistsError struct {
	typ SettingType
	k   string
}

func (e SettingExistsError) Error() string {
	// if typ is unknown or basic, don't include it in the o utput.
	if e.typ <= 1 {
		return fmt.Sprintf("%s: setting exists", e.k)
	}
	return fmt.Sprintf("%s: %s setting exists", e.k, e.typ)
}

// ShortFlagExistsError occurs when registering a flag whose short flag already
// exists/
type ShortFlagExistsError struct {
	k         string
	short     string
	shortName string
}

func (e ShortFlagExistsError) Error() string {
	return fmt.Sprintf("%s: short flag %q already exists for %q", e.k, e.short, e.shortName)
}

// SettingNotFoundError occurs when a setting isn't found.
type SettingNotFoundError struct {
	settingType SettingType
	k           string
}

func (e SettingNotFoundError) Error() string {
	if e.settingType <= 0 {
		return fmt.Sprintf("%s: setting not found", e.k)
	}
	return fmt.Sprintf("%s: %s setting not found", e.k, e.settingType)
}

// UnsupportedFormatError occurs when the string cannot be matched to a
// supported configuration format.
type UnsupportedFormatError struct {
	v string
}

func (e UnsupportedFormatError) Error() string {
	return fmt.Sprintf("%s: unsupported configuration format", e.v)
}

// NewUnsupportedFormatError returns an UnsupportedFormatError using the provided
// s.
func NewUnsupportedFormatError(s string) UnsupportedFormatError {
	return UnsupportedFormatError{s}
}

// PathsFromEnvVars returns a list of expanded paths found in the environemnt
// variable s. If nothing was found, or the environment variable was empty, a
// nil will be returned.
func PathsFromEnvVar(s string) []string {
	if s == "" {
		return nil
	}
	v := os.Getenv(s)
	if v == "" {
		return nil
	}
	paths := strings.Split(v, string(os.PathListSeparator))
	for i := range paths {
		paths[i] = os.ExpandEnv(paths[i])
	}
	return paths
}
