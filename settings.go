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

// Settings is a named group of settings and information related to that set,
// e.g. configuration file, if applicable, if it should emit an error when the
// configuration file is missing, flags parsed, etc. This is safe for
// concurrent use.
//
// The name of the Settings is used for environment variable naming, if
// applicable.
//
// If the configuration filename is set, an attempt will be made to laod
// setting information from the configuration file. Where to look for the file
// is configurable. Depending on the information provided and booleans set,
// Settings will look for the configuration file in those locations until the
// file is either found or all paths to check have been exhausted, which will
// result in an os.PathError containing an os.ErrNotExist being returned.
// Settings will look for the configuration file in the following order:
//
//    ConfFilePaths: in the order provided
//    ConfFilePathEnvVars: in the order provided
//    Working Directory: if set to true
//    Application Directory: if set to true
//    PATH environment variables: if set to search the $PATH
//
// The PathEnvVars will be checked using the values provided, they will not be
// prefixed with Settings name.
//
// There will always be a package global Settings whose name is the
// application's name.
type Settings struct {
	name string
	mu   sync.RWMutex
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
		errOnMissingConfFile: true,
		searchPath:           true,
		flagSet:              flag.NewFlagSet(name, flag.ContinueOnError),
		confFileVars:         map[string]struct{}{},
		flagVars:             map[string]interface{}{},
		shortFlags:           map[string]string{},
		settings:             map[string]setting{},
	}
}

// SetConfFilename sets the Settings' configuration filename and configures
// settings to use a configuration file. If the filename is empty, an error is
// returned.
func (s *Settings) SetConfFilename(v string) error {
	if v == "" {
		return fmt.Errorf("set configuration filename failed: no name provided")
	}

	// store the key value being used as the configuration setting name by caller
	s.mu.Lock()
	defer s.mu.Unlock()
	s.confFilename = v
	s.useConfFile = true
	return nil
}

// Set sets the registered settings according to Settings' configuration: it
// can be updated using a configuration file and/or environment variables; in
// that order of precedence. This is only run once; subsequent calls will
// result in no changes. Only settings that are of type ConfFileVar or EnvVar
// will be affected. This does not handle flags.
func (s *Settings) Set() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// if this has already been set from env vars and config, don't do it again.
	// TODO: decide if this should be handled differently to allow for reload.
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

// SetFromEnvVars sets the settings that are of type EnvVars from env vars if
// the Settings is set to use env vars. If any settings are registered as an
// EnnVars settings, the use env vars flag will be set to true. This can be
// overridden.
//
// Once a Settings has been set from environment variables they will not be
// updated again on subsequent calls.
//
// A setting's env name is a concatonation of the setting's name, an underscore
// (_), and the Settings' name, e.g. a Settings with the name 'foo' and a
// setting whose name is 'bar' will result in 'FOO_BAR'.
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

// SetFromConfFile set's the Conf, Env, and Flag settings from the information
// found in the configuration file if there is one. If Settings is not set to
// use a configuration file, if the configuration filename is not set, or if it
// has already been set, nothing is done and no error is returned.
//
// Settings may look for the configuration file according to how it's been
// configured if the file isn't found using the provided ConfFilename. When
// looking for the configuration file, Settings will extract the filename
// from the provided ConfFilename, if the information includes a path, and
// look for the configuration file according to its configuration:
//
// confFilePaths + filename
// confFileEnvVars + filename (each env var may have multiple path elements)
// working directory + filename
// executable directory + filename
// $PATH element + filename (PATH may have multiple path elements)
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
	if !s.useConfFile || s.confFilename == "" {
		return nil
	}
	// get the file's format from the extension
	f, err := formatFromFilename(s.confFilename)
	if err != nil {
		return err
	}

	b, err := s.readConfFile(s.confFilename)
	if err != nil {
		return err
	}

	cnf, err := unmarshalConfBytes(f, b)
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

	return nil, &os.PathError{Op: "open file", Path: fmt.Sprintf("%s: %s", n, errS[2:len(errS)]), Err: os.ErrNotExist}
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

// ErrOnMissingConfFile returns whether a missing config file should result in
// an error. This only applies when useConf == true
func (s *Settings) ErrOnMissingConfFile() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.errOnMissingConfFile
}

// SetErrOnMissingConfFile returns whether a missing config file should result
// in an error. This only applies when useFile == true
func (s *Settings) SetErrOnMissingConfFile(b bool) {
	s.mu.Lock()
	s.errOnMissingConfFile = b
	s.mu.Unlock()
}

// ConfFilePaths sets the paths that should be checked when looking for the
// configuration file. The paths will be checked in the order provided.
func (s *Settings) ConfFilePaths(paths []string) {
	s.mu.Lock()
	s.confFilePaths = paths
	s.mu.Unlock()
}

// ConfFilePathEnvVars sets the names of the environment variables that have
// paths that should be checked when looking for the configuration file. The
// environment variables will be checked in the order provided.
func ConfFilePathEnvVars(envVars []string) {
	settings.ConfFilePathEnvVars(envVars)
}

// ConfFilePathEnvVars sets the names of the environment variables that have
// paths that should be checked when looking for the configuration file. The
// environment variables will be checked in the order provided.
func (s *Settings) ConfFilePathEnvVars(envVars []string) {
	s.mu.Lock()
	s.confFilePathEnvVars = envVars
	s.mu.Unlock()
}

// CheckWD: if Settings should check the working directory for the
// configuration file.
func (s *Settings) CheckWD() {
	s.mu.RLock()
	s.checkWD = true
	defer s.mu.RUnlock()
}

// CheckExeDir: if Settings should check the executable directory for the
// configuration file.
func (s *Settings) CheckExeDir() {
	s.mu.RLock()
	s.checkExeDir = true
	defer s.mu.RUnlock()

}

// SearchPath: if Settings should use the user's PATH environment variable to
// check for the configuratiom file.
func (s *Settings) SearchPath(b bool) {
	s.mu.Lock()
	s.searchPath = b
	s.mu.Unlock()
}

// UseConfFile returns if ConfFileVar settings are to be updated from a
// configuration file.
func (s *Settings) UseConfFile() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.useConfFile
}

// UseEnvVars returns whether or not environment variables are used.
func (s *Settings) UseEnvVars() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.useEnvVars
}

// IsSet returns if the Settings has been set from all of its configured
// inputs, as applicable: env vars, configuration file, and flags.
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

// SetUsage sets the Usage func.
func (s *Settings) SetUsage(f func()) {
	s.mu.Lock()
	s.flagSet.Usage = f
	s.mu.Unlock()
}

// Name returns the Settings' name.
func (s *Settings) Name() string {
	return s.name
}

// IsCoreE returns if setting k is a Core setting. If setting k doesn't exist,
// a SettingNotFoundErr will be returned.
func (s *Settings) IsCoreE(k string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isCore(k)
}

func (s *Settings) isCore(k string) (bool, error) {
	val, ok := s.settings[k]
	if !ok {
		return false, SettingNotFoundErr{settingType: Core, k: k}
	}
	return val.IsCore, nil
}

// IsCore returns if setting k is a Core setting. If setting k doesn't exist,
// a false will be returned.
func (s *Settings) IsCore(k string) bool {
	b, _ := s.IsCoreE(k)
	return b
}

// IsConfFileVarE returns if setting k is a ConfFileVar setting. If setting k
// doesn't exist, a SettingNotFoundErr will be returned.
func (s *Settings) IsConfFileVarE(k string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isConfFileVar(k)
}

func (s *Settings) isConfFileVar(k string) (bool, error) {
	val, ok := s.settings[k]
	if !ok {
		return false, SettingNotFoundErr{settingType: ConfFileVar, k: k}
	}
	return val.IsConfFileVar, nil
}

// IsConfFileVar returns if setting k is a ConfFileVar setting. If setting k
// doesn't exist, a false will be returned.
func (s *Settings) IsConfFileVar(k string) bool {
	b, _ := s.IsConfFileVarE(k)
	return b
}

// IsEnvVarE returns if setting k is an EnvVar setting. If setting k doesn't
// exist, a SettingNotFoundErr will be returned.
func (s *Settings) IsEnvVarE(k string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isEnvVar(k)
}

func (s *Settings) isEnvVar(k string) (bool, error) {
	val, ok := s.settings[k]
	if !ok {
		return false, SettingNotFoundErr{settingType: EnvVar, k: k}
	}
	return val.IsEnvVar, nil
}

// IsEnvVar returns if setting k is an EnvVar setting. If setting k doesn't
// exist, a false will be returned.
func (s *Settings) IsEnvVar(k string) bool {
	b, _ := s.IsEnvVarE(k)
	return b
}

// IsFlagE returns if setting k is a Flag setting. If setting k doesn't exist,
// a SettingNotFoundErr will be returned.
func (s *Settings) IsFlagE(k string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isFlag(k)
}

func (s *Settings) isFlag(k string) (bool, error) {
	val, ok := s.settings[k]
	if !ok {
		return false, SettingNotFoundErr{settingType: Flag, k: k}
	}
	return val.IsFlag, nil
}

// IsFlag returns if setting k is a Flag setting. If setting k doesn't exist,
// a false will be returned.
func (s *Settings) IsFlag(k string) bool {
	b, _ := s.IsFlagE(k)
	return b
}

// EnvVarName returns the environment variable name for k. This will be
// NAME_K, where K is k and NAME is the name of Settings. For the pkg global
// Settings, this will be the executable name.
func (s *Settings) EnvVarName(k string) string {
	return strings.ToUpper(fmt.Sprintf("%s_%s", s.name, k))
}

// Exists returns if setting k exists.
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
		return nil, UnsupportedFormatErr{f.String()}
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

// SetConfFilename sets the Settings' configuration filename and configures
// settings to use a configuration file. If the filename is empty, an error is
// returned.
func SetConfFilename(v string) error { return settings.SetConfFilename(v) }

// Set sets the registered settings according to Settings' configuration: it
// can be updated using a configuration file and/or environment variables; in
// that order of precedence. This is only run once; subsequent calls will
// result in no changes. Only settings that are of type ConfFileVar or EnvVar
// will be affected. This does not handle flags.
func Set() error { return settings.Set() }

// SetFromEnvVars sets the settings that are of type EnvVar from env vars if
// the Settings is set to use env vars. If any settings are registered as an
// EnvVar settings, the use env vars flag will be set to true. This can be
// overridden.
//
// Once a Settings has been set from environment variables they will not be
// updated again on subsequent calls.
//
// A setting's env var name is a concatonation of the setting's name, an
// underscore, (_), and the Settings' name, e.g. a Settings with the name
// 'foo' and a setting whose name is 'bar' will result in 'FOO_BAR'.
func SetFromEnvVars() error { return settings.SetFromEnvVars() }

// SetFromConfFile set's the Conf, Env, and Flag settings from the information
// found in the configuration file if there is one. If Settings is not set to
// use a configuration file, if the configuration filename is not set, or if it
// has already been set, nothing is done and no error is returned.
//
// Settings may look for the configuration file according to how it's been
// configured if the file isn't found using the provided ConfFilename. When
// looking for the configuration file, Settings will extract the filename
// from the provided ConfFilename, if the information includes a path, and
// look for the configuration file according to its configuration:
//
// confFilePaths + filename
// confFileEnvVars + filename (each env var may have multiple path elements)
// working directory + filename
// executable directory + filename
// $PATH element + filename (PATH may have multiple path elements)
//
// Any of the above elements that are either empty or false are skipped.
//
// If the file cannot be found, an os.PathError with an os.ErrNotExist and
// a list of all paths checked is returned.
func SetFromConfFile() error {
	return settings.SetFromConfFile()
}

// ErrOnMissingConfFile returns whether a missing config file should result in
// an error. This only applies when useConf == true
func ErrOnMissingConfFile() bool { return settings.ErrOnMissingConfFile() }

// SetErrOnMissingConfFile returns whether a missing config file should result
// in an error. This only applies when useFile == true
func SetErrOnMissingConfFile(b bool) { settings.SetErrOnMissingConfFile(b) }

// ConfFilePaths sets the paths that should be checked when looking for the
// configuration file. The paths will be checked in the order provided.
func ConfFilePaths(paths []string) {
	settings.ConfFilePaths(paths)
}

// CheckWD: if Settings should check the working directory for the
// configuration file.
func CheckWD() { settings.CheckWD() }

// CheckExeDir: if Settings should check the executable directory for the
// configuration file.
func CheckExeDir() { settings.CheckExeDir() }

// SearchPath: if Settings should use the user's PATH environment variable to
// check for the configuratiom file.
func SearchPath(b bool) { settings.SearchPath(b) }

// UseConfFile returns if ConfFileVar settings are to be updated from a
// configuration file.
func UseConfFile() bool { return settings.UseConfFile() }

// UseEnvVars returns whether or not environment variables are used.
func UseEnvVars() bool { return settings.useEnvVars }

// IsSet returns if the Settings has been set from all of its configured
// inputs, as applicable: env vars, configuration file, and flags.
func IsSet() bool { return settings.IsSet() }

// SetUsage sets the Usage func.
func SetUsage(f func()) { settings.SetUsage(f) }

// Name returns the Settings' name.
func Name() string { return settings.Name() }

// IsCoreE returns if setting k is a Core setting. If setting k doesn't exist,
// a SettingNotFoundErr will be returned.
func IsCoreE(k string) (bool, error) { return settings.IsCoreE(k) }

// IsCore returns if setting k is a Core setting. If setting k doesn't exist,
// a false will be returned.
func IsCore(k string) bool { return settings.IsCore(k) }

// IsConfFileVarE returns if setting k is a ConfFileVar setting. If setting k
// doesn't exist, a SettingNotFoundErr will be returned.
func IsConfFileVarE(k string) (bool, error) { return settings.IsConfFileVarE(k) }

// IsConfFileVar returns if setting k is a ConfFileVar setting. If setting k
// doesn't exist, a false will be returned.
func IsConfFileVar(k string) bool { return settings.IsConfFileVar(k) }

// IsEnvVarE returns if setting k is an EnvVar setting. If setting k doesn't
// exist, a SettingNotFoundErr will be returned.
func IsEnvVarE(k string) (bool, error) { return settings.IsEnvVarE(k) }

// IsEnvVar returns if setting k is an EnvVar setting. If setting k doesn't
// exist, a false will be returned.
func IsEnvVar(k string) bool { return settings.IsEnvVar(k) }

// IsFlagE returns if setting k is a Flag setting. If setting k doesn't exist,
// a SettingNotFoundErr will be returned.
func IsFlagE(k string) (bool, error) { return settings.IsFlagE(k) }

// IsFlag returns if setting k is a Flag setting. If setting k doesn't exist,
// a false will be returned.
func IsFlag(k string) bool { return settings.IsFlag(k) }

// EnvVarName returns the environment variable name for k. This will be
// NAME_K, where K is k and NAME is the name of Settings. For the pkg global
// Settings, this will be the executable name.
func EnvVarName(k string) string { return settings.EnvVarName(k) }

// Exists returns if setting k exists.
func Exists(k string) bool { return settings.Exists(k) }
