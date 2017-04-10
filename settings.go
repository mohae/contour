package contour

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/kardianos/osext"
	"github.com/mohae/cjson"
	"gopkg.in/yaml.v2"
)

// Settings is a named group of settings and information related to that set.
//
// The name of the Settings is used for environment variable naming, if
// applicable.
//
// Configuration settings, settings that can be set from a configuration file,
// an environment variable, or a flag are registered with their default value.
// Settings will then update them with the anything it finds using the Set and
// ParseFlags methods. Configuration settings are updated according to the type
// that the are registered as; with higher precedence types being updatable
// from a lower precedence source, e.g. a Flag setting can be updated from a
// configuration file, an environment variable, or a flag, while a ConfFileVar
// setting can only be updated by a configuration file. The order of precedence
// is:
//    ConfFileVar
//    EnvVar
//    Flag
//
// In addition to configuration settings there are settings and Core settings.
// Core settings cannot be changed once they are set; any attempt to update a
// core setting will result in an error. Core settings are added using addCore
// methods. Regular settings cannot be updated by an external source, e.g. a
// configuration file, an environment variable, or a flag, but can be updated
// using an Update method.
//
// If the configuration filename was set using the SetConfFilename method or a
// setting was registered, an attempt will be made to laod setting information
// from a configuration file. If the configuration filename wasn't explicitly
// set, settings will build it using it's name and configured format, which
// defaults to JSON. The extension of the configuration file will be the
// format, i.i. JSON's extension will be 'json', TOML's extension will be
// 'toml', and YAML's extension will be 'yaml'.
//
// Settings will look for the configuration file according to the information
// it has and how it has been configured, e.g. look in additional paths
// provided to it, look in environment variables for paths, look in the
// executable directory, search the $PATH, etc. Settings defaults to returning
// an os.ErrNotExist os.PathError but it can be configured to not return an
// error if the configuration file is not found. New Settings are configured
// to search for the confiugration file in the $PATH.
//
// If any configuration settings were registered as a Flag or EnvVar setting,
// Settings will check the setting's corresponding environment variable. A
// setting's environment variable is a concatonation of the Settings' name and
// the setting's key, in UPPER_SNAKE_CASE, e.g. the environmnet variable name
// for a Settings with the name foo and a setting with the key bar would be
// FOO_BAR.
//
// Once a Settings has been updated from a configuration file and environment
// variables, configuration settings cannot be updated again from those
// sources, subsequent attempts will result in no error and nothing being done.
//
// If any configuration settings were registered as a Flag, ParseFlags should
// be called. Once the flags have been parsed, subsequent calls to ParseFlags
// will return an ErrFlagsParsed error. If settings is not configured to use
// flags, ParseFlags will return an ErrUSeFlagsFalse error.
//
// Settings are safe for concurrent use.
type Settings struct {
	name   string
	mu     sync.RWMutex
	format Format
	// if an attempt to load configuration from a file should error if the file
	// does not exist.
	errOnMissingConfFile bool
	// the name of the configuration filename variable: defaults to the
	// ConfFilenameSettingVarName constant.
	confFilenameVarName string
	// file is the name of the configuration file
	confFilename string
	// Encoding is what encoding scheme is used for this config.
	encoding string
	// Tracks the vars that are exposed to the configuration file. Only vars in
	// this map can be set from a configuration file.
	confFileVars map[string]struct{}
	// If settings should be updated from a configuration file.
	useConfFile bool
	// If the settings have been updated from configuration file.
	confFileVarsSet bool
	// confFilePaths is a list of paths in which to check for the configuration
	// file
	confFilePaths []string
	// confFilePathEnvVars is a list of environment variables to check for path
	// information in which to check for the configuration file.
	confFilePathEnvVars []string
	// Look in the working directory for the configuration file.
	checkWD bool
	// look in the executabledir for the configuration file
	checkExeDir bool
	// search the path env var, in addition to wd & executalbe dir, for the conf
	// file. The path env var is assumed to be NAMEPATH where name is
	// Settings.name
	searchPath bool
	// If settings should be loaded from environment variables. Environment
	// variable names will be in the form of NAME_VARNAME where NAME is this
	// Setting's name.
	useEnvVars bool
	// If the settings have been updated from environment variables.
	envVarsSet bool
	// flagset is the set of flags for arg parsing.
	flagSet *flag.FlagSet
	// If the settings should be updated from passed args.
	useFlags bool
	// If the settings have been updated from the passed args.
	flagsParsed bool
	// The map of variables that capture flag information.
	flagVars map[string]interface{}
	// Maps short flags to the long version
	shortFlags map[string]string
	// parsedFlags are flags that were passed and parsed. Short flags are
	// normalized to the flag name.
	parsedFlags []string
	// Settings contains a map of all the configuration settings for this
	// Setting and each setting's information, including current Value.
	settings map[string]setting
}

// New provides an initialized Settings named name.
func New(name string) *Settings {
	return &Settings{
		name:                 name,
		format:               format,
		errOnMissingConfFile: true,
		searchPath:           true,
		flagSet:              flag.NewFlagSet(name, flag.ContinueOnError),
		confFileVars:         map[string]struct{}{},
		flagVars:             map[string]interface{}{},
		shortFlags:           map[string]string{},
		settings:             map[string]setting{},
	}
}

// SetConfFilename sets the settings' configuration filename and configures
// settings to use a configuration file. If the filename is empty, an error is
// returned.
func (s *Settings) SetConfFilename(v string) error {
	if v == "" {
		return fmt.Errorf("configuration filename: set failed: no name provided")
	}
	// get the file's format from the extension
	f, err := formatFromFilename(v)
	if err != nil {
		return err
	}

	// store the key value being used as the configuration setting name by caller
	s.mu.Lock()
	defer s.mu.Unlock()
	s.confFilename = v
	s.useConfFile = true
	s.format = f
	return nil
}

// Set updates the settings' configuration from a configuration file and
// environment variables. This is only run once; subsequent calls will result
// in no changes. Only settings that are of type ConfFileVar or EnvVar will be
// affected. This does not handle flags.
//
// Once the standard settings has been set, updated, it will not update again;
// subsequent calls will result in nothing being done.
//
// All ConfFileVar, EnvVar, and Flag settings must be registered before calling
func (s *Settings) Set() error {
	// Set.
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.confFileVarsSet && s.envVarsSet {
		return nil
	}
	err := s.updateFromEnvVars()
	if err != nil {
		return fmt.Errorf("setting configuration from env failed: %s", err)
	}

	err = s.setFromConfFile()
	if err != nil {
		return fmt.Errorf("setting configuration from file failed: %s", err)
	}
	return nil
}

// SetFromEnvVars updates the settings' configuration from environment
// variables.
//
// Once the standard settings has been set, updated, from environment
// variables, they will not be updated again from environment variables;
// subsequent calls will result in nothing being done.
//
// A setting's env name is a concatonation of the settings' name, an underscore
// (_), and the setting's key, e.g. given a settings with the name 'foo', a
// setting whose key is 'bar' will be updateable with the environment variable
// FOO_BAR.
func (s *Settings) SetFromEnvVars() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.updateFromEnvVars()
}

func (s *Settings) updateFromEnvVars() error {
	if !s.useEnvVars || s.envVarsSet {
		return nil
	}
	var err error
	for k, v := range s.settings {
		if !v.IsEnvVar {
			continue
		}
		tmp := os.Getenv(s.EnvVarName(k))
		if tmp != "" {
			switch v.Type {
			case _bool:
				b, _ := strconv.ParseBool(tmp)
				err = s.updateBool(EnvVar, k, b)
			case _int:
				i, err := strconv.Atoi(tmp)
				if err != nil {
					return fmt.Errorf("getenv %s: %s", s.EnvVarName(k), err)
				}
				err = s.updateInt(EnvVar, k, i)
			case _int64:
				i, err := strconv.ParseInt(tmp, 10, 64)
				if err != nil {
					return fmt.Errorf("getenv %s: %s", s.EnvVarName(k), err)
				}
				err = s.updateInt64(EnvVar, k, i)
			case _string:
				err = s.updateString(EnvVar, k, tmp)
			default:
				return fmt.Errorf("%s: unsupported env variable type: %s", s.EnvVarName(k), v.Type)
			}
			if err != nil {
				return fmt.Errorf("get env %s: %s", s.EnvVarName(k), err)
			}
			// lock to check next setting, if there is one.
		}
	}
	// Rlock isn't sufficient for updating to close it and get a Lock() for update.
	s.envVarsSet = true
	return nil
}

// SetFromConfFile updates the settings' configuration from the configuration
// file. If the configuration filename was not set using the SetConfFilename
// method, settings will look for the configuration file using the settings'
// name and it's format: name.format.
//
// The format, in full, is used as the extension, i.e. JSON's extension will be
// 'json', TOML's extension will be 'toml', and YAML's extension will be
// 'yaml'.
//
// Once the settings has been set, updated from the configuration file, they
// will not be updated again from the configuration file; subsequent calls will
// result in nothing being done.
//
// The settings will look for the configuration file according to how it's been
// configured.
//     filename
//     confFilePaths + filename
//     confFileEnvVars + filename (each env var may have multiple path elements)
//     working directory + filename
//     executable directory + filename
//     $PATH element + filename (PATH may have multiple path elements)
//
// Any of the above elements that are either empty or false are skipped.
//
// If the file cannot be found, an os.PathError with an os.ErrNotExist and
// a list of all paths checked is returned.
func (s *Settings) SetFromConfFile() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.setFromConfFile()
}

// setFromFile set's Conf, Env, and Flag settings from the information found in
// the configuration file, if there is one. This assumes the caller already
// holds the lock.
// TODO:
//	add a confFileRequired flag
func (s *Settings) setFromConfFile() error {
	if !s.useConfFile {
		return nil
	}

	if s.confFilename == "" { // if it wasn't explicitly set, create the name
		s.confFilename = s.name + "." + s.format.String()
	}

	b, err := s.readConfFile(s.confFilename)
	if err != nil {
		return err
	}

	cnf, err := unmarshalConfBytes(s.format, b)
	if err != nil {
		return fmt.Errorf("%s: %s", s.confFilename, err)
	}

	// if nothing was returned and no error, nothing to do
	if cnf == nil {
		return nil
	}
	// Go through settings and update setting values.
	for k, v := range cnf.(map[string]interface{}) {
		// otherwise update the setting
		err = s.update(ConfFileVar, k, v)
		if err != nil {
			return fmt.Errorf("update setting: %s", err)
		}
	}
	s.confFileVarsSet = true
	return nil
}

// readConfFile reads the configuration file n.
func (s *Settings) readConfFile(n string) (b []byte, err error) {
	b, err = ioutil.ReadFile(n)
	if err == nil {
		return b, nil
	}

	// get the filename part of n
	fname := filepath.Base(n)

	var (
		ps   []string
		errS string // accumulates info for the error if nothing is found
	)

	if len(s.confFilePaths) > 0 {
		b, err = s.checkPaths(fname, s.confFilePaths)
		if err == nil { // was found
			return b, nil
		}
		errS = "; " + strings.Join(s.confFilePaths, "; ")
	}

	if len(s.confFilePathEnvVars) > 0 {
		for _, v := range s.confFilePathEnvVars {
			p := os.Getenv(v)
			tmp := getEnvVarPaths(p)
			if len(tmp) == 0 {
				continue
			}
			ps = append(ps, tmp...)
		}
		b, err = s.checkPaths(fname, ps)
		if err == nil {
			return b, nil
		}
		ps = ps[0:0]
		errS += "; $" + strings.Join(s.confFilePathEnvVars, "; $")
	}

	if s.checkWD {
		d, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("load conf file %s: get wd: %s", n, err)
		}
		ps = append(ps, d)
		b, err = s.checkPaths(fname, ps)
		if err == nil {
			return b, nil
		}
		ps = ps[0:0]
		errS += "; " + d
	}

	if s.checkExeDir {
		d, err := osext.ExecutableFolder()
		if err != nil {
			return nil, fmt.Errorf("load conf file %s: get wd: %s", n, err)
		}
		ps = append(ps, d)
		b, err = s.checkPaths(fname, ps)
		if err == nil {
			return b, nil
		}
		ps = ps[0:0]
		errS += "; " + d
	}

	// search the path, if applicable
	if s.searchPath {
		v := os.Getenv("PATH")
		ps = getEnvVarPaths(v)
		b, err = s.checkPaths(fname, ps)
		if err == nil {
			return b, nil
		}
		errS += "; $PATH"
	}

	if len(errS) >= 2 {
		errS = fmt.Sprintf("%s: %s", n, errS[2:len(errS)])
	} else {
		errS = n
	}

	return nil, &os.PathError{Op: "open file", Path: errS, Err: os.ErrNotExist}
}

func (s *Settings) checkPaths(fname string, paths []string) (b []byte, err error) {
	for _, v := range paths {
		tmp := filepath.Join(v, fname)
		f, err := os.Open(tmp)
		if err == nil {
			defer f.Close()
			return ioutil.ReadAll(f)
		}
	}
	return nil, os.ErrNotExist
}

// ErrOnMissingConfFile returns if settings is configured to return an error
// if the configuration file cannot be located.
func (s *Settings) ErrOnMissingConfFile() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.errOnMissingConfFile
}

// SetErrOnMissingConfFile sets if settings should return an error if the
// confiugration file cannot be located.
func (s *Settings) SetErrOnMissingConfFile(b bool) {
	s.mu.Lock()
	s.errOnMissingConfFile = b
	s.mu.Unlock()
}

// ConfFilename returns the settings' configuration filename.
func (s *Settings) ConfFilename() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.confFilename
}

// ConfFilePaths sets the paths that settings should check when looking for the
// configuration file. The paths will be checked in the order provided.
func (s *Settings) ConfFilePaths(paths []string) {
	s.mu.Lock()
	s.confFilePaths = paths
	s.mu.Unlock()
}

// ConfFilePathEnvVars sets the names of the environment variables that have
// paths that settings should check when looking for the configuration file.
// The environment variables will be checked in the order provided. The
// environment variables may contain multiple paths.
func (s *Settings) ConfFilePathEnvVars(envVars []string) {
	s.mu.Lock()
	s.confFilePathEnvVars = envVars
	s.mu.Unlock()
}

// CheckWD returns if settings should check the working directory for the
// configuration file.
func (s *Settings) CheckWD() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.checkWD
}

// SetCheckWD sets if settings should check the working directory for the
// configuration file.
func (s *Settings) SetCheckWD(b bool) {
	s.mu.Lock()
	s.checkWD = b
	s.mu.Unlock()
}

// CheckExeDir returns if Settings should check the executable directory for the
// configuration file.
func (s *Settings) CheckExeDir() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.checkExeDir
}

// SetCheckExeDir sets if Settings should check the executable directory for
// the configuration file.
func (s *Settings) SetCheckExeDir(b bool) {
	s.mu.Lock()
	s.checkExeDir = b
	s.mu.Unlock()
}

// SearchPath sets if settings should use the user's PATH environment variable
// to check for the configuratiom file.
func (s *Settings) SearchPath(b bool) {
	s.mu.Lock()
	s.searchPath = b
	s.mu.Unlock()
}

// GetFormat returns the format to use if the ConfFilename hasn't been
// explicitly set.
func (s *Settings) GetFormat() Format {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.format
}

// SetFormat sets the format to use for the conffile if the ConfFilename hasn't
// been explicitly set.
func (s *Settings) SetFormat(f Format) {
	s.mu.Lock()
	s.format = f
	s.mu.Unlock()
}

// GetFormatString returns the format to use if the ConfFilename hasn't been
// explicitly set as a string
func (s *Settings) GetFormatString() string {
	s.mu.RLock()
	s.mu.RUnlock()
	return s.format.String()
}

// SetFormatString sets the format, using a string, to use for the
// configuration file if the ConfFilename hasn't been explicitly set. If the
// string isn't parsable to a supported Format, an UnsupportedFormatError will
// be returned.
func (s *Settings) SetFormatString(v string) error {
	f, err := ParseFormat(v)
	if err != nil {
		return err
	}
	s.SetFormat(f)
	return nil
}

// UseConfFile returns if settings will update its configuration settings from
// a configuration file.
func (s *Settings) UseConfFile() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.useConfFile
}

// UseEnvVars returns if settings will update its configuration settings from
// environment variables.
func (s *Settings) UseEnvVars() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.useEnvVars
}

// IsSet returns if settings' configuration settings have been set from all of
// its configured sources.
func (s *Settings) IsSet() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.useConfFile && !s.confFileVarsSet {
		return false
	}
	if s.useEnvVars && !s.envVarsSet {
		return false
	}
	if s.useFlags && !s.flagsParsed {
		return false
	}
	// Everything has been updated (is set) according to Settings' configuration.
	return true
}

// SetUsage sets settings' Usage func.
func (s *Settings) SetUsage(f func()) {
	s.mu.Lock()
	s.flagSet.Usage = f
	s.mu.Unlock()
}

// Name returns settings' name.
func (s *Settings) Name() string {
	return s.name
}

// IsCoreE returns if setting k is a Core setting. A SettingNotFoundErr will be
// returned if k doesn't exist in settings.
func (s *Settings) IsCoreE(k string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isCore(k)
}

func (s *Settings) isCore(k string) (bool, error) {
	val, ok := s.settings[k]
	if !ok {
		return false, SettingNotFoundError{settingType: Core, k: k}
	}
	return val.IsCore, nil
}

// IsCore returns if setting k is a Core setting. False will be returned if k
// doesn't exist in settings.
func (s *Settings) IsCore(k string) bool {
	b, _ := s.IsCoreE(k)
	return b
}

// IsConfFileVarE returns if setting k is a ConfFileVar setting. A
// SettingNotFoundErr will be returned if k doesn't exist in settings.
func (s *Settings) IsConfFileVarE(k string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isConfFileVar(k)
}

func (s *Settings) isConfFileVar(k string) (bool, error) {
	val, ok := s.settings[k]
	if !ok {
		return false, SettingNotFoundError{settingType: ConfFileVar, k: k}
	}
	return val.IsConfFileVar, nil
}

// IsConfFileVar returns if setting k is a ConfFileVar setting. False will be
// returned if k doesn't exist in settings.
func (s *Settings) IsConfFileVar(k string) bool {
	b, _ := s.IsConfFileVarE(k)
	return b
}

// IsEnvVarE returns if setting k is an EnvVar setting. A SettingNotFoundErr
// will be returned if k doesn't exist in settings.
func (s *Settings) IsEnvVarE(k string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isEnvVar(k)
}

func (s *Settings) isEnvVar(k string) (bool, error) {
	val, ok := s.settings[k]
	if !ok {
		return false, SettingNotFoundError{settingType: EnvVar, k: k}
	}
	return val.IsEnvVar, nil
}

// IsEnvVar returns if setting k is an EnvVar setting. False will be returned
// if k doesn't exist in settings.
func (s *Settings) IsEnvVar(k string) bool {
	b, _ := s.IsEnvVarE(k)
	return b
}

// IsFlagE returns if setting k is a Flag setting.  A SettingNotFoundErr will
// be returned if k doesn't exist in settings.
func (s *Settings) IsFlagE(k string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isFlag(k)
}

func (s *Settings) isFlag(k string) (bool, error) {
	val, ok := s.settings[k]
	if !ok {
		return false, SettingNotFoundError{settingType: Flag, k: k}
	}
	return val.IsFlag, nil
}

// IsFlag returns if setting k is a Flag setting. False will be returned if k
// doesn't exist in settings.
func (s *Settings) IsFlag(k string) bool {
	b, _ := s.IsFlagE(k)
	return b
}

// Exists returns if settings has a setting k.
func (s *Settings) Exists(k string) bool {
	s.mu.RLock()
	s.mu.RUnlock()
	return s.exists(k)
}

func (s *Settings) exists(k string) bool {
	_, err := s.get(k)
	if err == nil {
		return true
	}
	return false
}

// EnvVarName returns the environment variable name for k. This will be
// NAME_K, where K is k and NAME is settings' name.
func (s *Settings) EnvVarName(k string) string {
	return strings.ToUpper(fmt.Sprintf("%s_%s", s.name, k))
}

// formatFromFilename gets the format from the passed filename.  An error will
// be returned if either the format isn't supported or the extension doesn't
// exist.  If the passed string has multiple dots, the last dot is assumed to
// be the extension.
func formatFromFilename(s string) (Format, error) {
	if s == "" {
		return Unsupported, fmt.Errorf("no configuration filename")
	}
	parts := strings.Split(s, ".")
	format := ""
	// case 0 has already been evaluated
	switch len(parts) {
	case 1:
		return Unsupported, fmt.Errorf("unable to determine %s's format: no extension", strings.TrimSpace(s))
	case 2:
		format = parts[1]
	default:
		// assume its the last part
		format = parts[len(parts)-1]
	}

	return ParseFormat(format)
}

// unmarshalConfBytes accepts bytes and unmarshals them using the correct
// format. Either the unmarshaled data or an error is returned.
//
// Supported formats:
//   json
//   toml
//   yaml
func unmarshalConfBytes(f Format, buff []byte) (interface{}, error) {
	var ret interface{}
	switch f {
	case JSON:
		err := cjson.Unmarshal(buff, &ret)
		if err != nil {
			return nil, err
		}
		return ret, nil
	case TOML:
		_, err := toml.Decode(string(buff), &ret)
		if err != nil {
			return nil, err
		}
		return ret, nil
	case YAML:
		err := yaml.Unmarshal(buff, &ret)
		if err != nil {
			return nil, err
		}
		return ret, nil
	default:
		return nil, UnsupportedFormatError{f.String()}
	}
}

// get the value of an env var that is assumed to have path info; split it
// into its path elements and expand them.
func getEnvVarPaths(s string) []string {
	if s == "" {
		return nil
	}

	p := strings.Split(s, string(os.PathListSeparator))
	for i := range p {
		p[i] = os.ExpandEnv(p[i])
	}
	return p
}

// SetConfFilename sets the standard settings' configuration filename and
// configures settings to use a configuration file. If the filename is empty,
// an error is returned.
func SetConfFilename(v string) error { return std.SetConfFilename(v) }

// Set updates the standard settings' configuration from a configuration file
// and environment variables. This is only run once; subsequent calls will
// result in no changes. Only settings that are of type ConfFileVar or EnvVar
// will be affected. This does not handle flags.
//
// Once the standard settings has been set, updated, it will not update again;
// subsequent calls will result in nothing being done.
//
// All ConfFileVar, EnvVar, and Flag settings must be registered before calling
// Set.
func Set() error { return std.Set() }

// SetFromEnvVars updates the standard settings' configuration from environment
// variables.
//
// Once the standard settings has been set, updated, from environment variables
// they, will not be updated again from environment variables; subsequent calls
// will result in nothing being done.
//
// A setting's env name is a concatonation of the settings' name, an underscore
// (_), and the setting's key, e.g. given a settings with the name 'foo', a
// setting whose key is 'bar' will be updateable with the environment variable
// FOO_BAR.
func SetFromEnvVars() error { return std.SetFromEnvVars() }

// SetFromConfFile updates the settings' configuration from the configuration
// file. If the configuration filename was not set using the SetConfFilename
// method, settings will look for the configuration file using the settings'
// name and it's format: name.format.
//
// The format, in full, is used as the extension, i.e. JSON's extension will be
// 'json', TOML's extension will be 'toml', and YAML's extension will be
// 'yaml'.
//
// Once the settings has been set, updated from the configuration file, they
// will not be updated again from the configuration file; subsequent calls will
// result in nothing being done.
//
// The settings will look for the configuration file according to how it's been
// configured.
//     filename
//     confFilePaths + filename
//     confFileEnvVars + filename (each env var may have multiple path elements)
//     working directory + filename
//     executable directory + filename
//     $PATH element + filename (PATH may have multiple path elements)
//
// Any of the above elements that are either empty or false are skipped.
//
// If the file cannot be found, an os.PathError with an os.ErrNotExist and
// a list of all paths checked is returned.
func SetFromConfFile() error {
	return std.SetFromConfFile()
}

// ErrOnMissingConfFile returns if the standard settings is configured to
// return an error if the configuration file cannot be located.
func ErrOnMissingConfFile() bool { return std.ErrOnMissingConfFile() }

// SetErrOnMissingConfFile sets if the standard settings should return an error
// if the confiugration file cannot be located.
func SetErrOnMissingConfFile(b bool) { std.SetErrOnMissingConfFile(b) }

// ConfFilename returns the standard settings' configuration filename.
func ConfFilename() string { return std.ConfFilename() }

// ConfFilePaths sets the paths that the standard settings should check when
// looking for the configuration file. The paths will be checked in the order
// provided.
func ConfFilePaths(paths []string) {
	std.ConfFilePaths(paths)
}

// ConfFilePathEnvVars sets the names of the environment variables that have
// paths that the standard settings should check when looking for the
// configuration file. The environment variables will be checked in the order
// provided. The environment variables may contain multiple paths.
func ConfFilePathEnvVars(envVars []string) {
	std.ConfFilePathEnvVars(envVars)
}

// CheckWD returns if settings should check the working directory for the
// configuration file.
func CheckWD() bool { return std.CheckWD() }

// SetCheckWD sets if settings should check the working directory for the
// configuration file.
func SetCheckWD(b bool) { std.SetCheckWD(b) }

// CheckExeDir returns if Settings should check the executable directory for the
// configuration file.
func CheckExeDir() bool { return std.CheckExeDir() }

// SetCheckExeDir sets if the standard settings should check the executable
// directory for the configuration file.
func SetCheckExeDir(b bool) { std.SetCheckExeDir(b) }

// SearchPath sets if the standard settings should use the user's PATH
// environment variable to check for the configuratiom file.
func SearchPath(b bool) { std.SearchPath(b) }

// GetFormat returns the standard settings' format to use if the ConfFilename
// hasn't been explicitly set.
func GetFormat() Format {
	return std.GetFormat()
}

// SetFormat sets the standard settings' format to use for the conffile if the
// ConfFilename hasn't been explicitly set.
func SetFormat(f Format) {
	std.SetFormat(f)
}

// GetFormatString returns the standard settings' format to use if the
// ConfFilename hasn't been explicitly set as a string
func GetFormatString() string {
	return std.GetFormatString()
}

// SetFormatString sets the standard settings' format, using a string, to use
// for the configuration file if the ConfFilename hasn't been explicitly set.
// If the string isn't parsable to a supported Format, an
// UnsupportedFormatError will be returned.
func SetFormatString(v string) error {
	return std.SetFormatString(v)
}

// UseConfFile returns if the standard settings will update its configuration
// settings from a configuration file.
func UseConfFile() bool { return std.UseConfFile() }

// UseEnvVars returns if the standard settings will update its configuration
// settings from environment variables.
func UseEnvVars() bool { return std.useEnvVars }

// IsSet returns if the standard settings' configuration settings have been set
// from all of its configured sources.
func IsSet() bool { return std.IsSet() }

// SetUsage sets the standard settings' Usage func.
func SetUsage(f func()) { std.SetUsage(f) }

// Name returns the standard settings' name.
func Name() string { return std.Name() }

// IsCoreE returns if setting k is a Core setting. A SettingNotFoundErr will be
// returned if k doesn't exist in the standard settings.
func IsCoreE(k string) (bool, error) { return std.IsCoreE(k) }

// IsCore returns if setting k is a Core setting. False will be returned if k
// doesn't exist in settings.
func IsCore(k string) bool { return std.IsCore(k) }

// IsConfFileVarE returns if setting k is a ConfFileVar setting. A
// SettingNotFoundErr will be returned if k doesn't exist in the standard
// settings.
func IsConfFileVarE(k string) (bool, error) { return std.IsConfFileVarE(k) }

// IsConfFileVar returns if setting k is a ConfFileVar setting. A false will be
// returned if k doesn't exist in the standard settings.
func IsConfFileVar(k string) bool { return std.IsConfFileVar(k) }

// IsEnvVarE returns if setting k is an EnvVar setting. A SettingNotFoundErr
// will be returned if k doesn't exist in the standard settings.
func IsEnvVarE(k string) (bool, error) { return std.IsEnvVarE(k) }

// IsEnvVar returns if setting k is an EnvVar setting. A false will be returned
// if k doesn't exist in the standard settings.
func IsEnvVar(k string) bool { return std.IsEnvVar(k) }

// IsFlagE returns if setting k is a Flag setting. A SettingNotFoundErr will be
// returned if k doesn't exist in the standard settings.
func IsFlagE(k string) (bool, error) { return std.IsFlagE(k) }

// IsFlag returns if setting k is a Flag setting. A false will be returned if k
// doesn't exist in the standard settings.
func IsFlag(k string) bool { return std.IsFlag(k) }

// Exists returns if setting k exists. A false will be be returned if k doesn't
// exist in the standard settings.
func Exists(k string) bool { return std.Exists(k) }

// EnvVarName returns the environment variable name for k. This will be
// NAME_K, where K is k and NAME is the standard settings' name (executable
// name).
func EnvVarName(k string) string { return std.EnvVarName(k) }
