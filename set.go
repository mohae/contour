package contour
// Set contains all of contour's Set functions.Calling Set
// adds, or registers, the settings information to the AppConfig variable.
// The setting value, if there is one, is not saved to its environment 
// variable at this point.
//

// SetEnvs writes the current contents of AppConfig to their respective
// environmnet variables.

import (
	"errors"
	"os"
	"strconv"	
)

// SetEnvs goes through AppConfig and saves all of the settings to their
// environment variables.
func SetEnvs() error {
	var err error
	// For each setting
	for k, setting := range AppConfig.Settings {
		switch setting.Type {
		case "bool":
			err = os.Setenv(k, strconv.FormatBool(setting.Value.(bool)))
		case "int", "string":
			err = os.Setenv(k, setting.Value.(string))
		default:
			err = errors.New(k + "'s datatype, " + setting.Type + ", is not supported")
		}

		if err != nil {
			return err
		}

	}

	return nil

}

// setEnv set's the environment variable. No validation is done here, that's
// the callers responsibility.
func setEnv(k string, v interface{}) error {
	var err error

	// Update the setting with file's
	switch AppConfig.Settings[k].Type {
	case "int", "string":
		err = os.Setenv(k,v.(string))		
	case "bool":
		s := strconv.FormatBool(v.(bool))
		err = os.Setenv(k, s)
	default:
		err = errors.New(k + "'s datatype, " + AppConfig.Settings[k].Type + ", is not supported")
	}

	return err
}

// setEnvFromConfigFile goes through all the settings in the configFile and
// checks to see if the setting is updateable; saving those that are to their
// environment variable.
func setEnvFromConfigFile() error {
	var err error

	for k, v := range configFile {
		// Find the key in the settings
		_, ok := AppConfig.Settings[k]
		if !ok {
			// skip settings that don't already exist
			continue
		}

		// Skip if Immutable, IsCore, IsEnv since they aren't 
		//overridable by ConfigFile.
		if !canUpdate(k) {
			continue
		}

		err = setEnv(k, v)
		if err != nil {
			return err
		}

		// Update the setting with file's
		switch AppConfig.Settings[k].Type {
		case "string":
			err = UpdateString(k,v.(string))	
		case "bool":
			err = UpdateBool(k,v.(bool))	
		case "int":
			err = UpdateInt(k,v.(int))	
		default:
			return errors.New(k + "'s datatype, " + AppConfig.Settings[k].Type + ", is not supported")
		}

		if err != nil {
			return err
		}

	}

	return nil
}

// SetSetting
// TODO figure this out
func SetSetting(Type, k string, v interface{}, Code string, Immutable, IsCore, IsEnv, IsFlag bool) error {
	_, ok := AppConfig.Settings[k]
	if ok {
		return nil
	}

	AppConfig.Settings[k] = &setting{
		Type: Type,
		Value: v,
		Code: Code,
		Immutable: Immutable,
		IsCore: IsCore,
		IsEnv: IsEnv,
		IsFlag: IsFlag,
	}

	return nil
}
	
// SetFlagBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func SetFlagBool(k string, v bool, f string) error {
	err := SetSetting("bool", k, v, f, false, false, false, true)
	if err != nil {
		return err
	}

	s := strconv.FormatBool(v)
	return os.Setenv(k,s)
}

// SetFlagInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func SetFlagInt(k string, v int, f string) error {
	err := SetSetting("int", k, v, f, false, false, false, true)
	if err != nil {
		return err
	}

	return os.Setenv(k,string(v))
}

// SetFlagString adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func SetFlagString(k, v, f string) error {
	err := SetSetting("string", k, v, f, false, false, false, true)
	if err != nil {
		return err
	}

	return os.Setenv(k,v)
}

// SetImmutableBool adds the information to the AppsConfig struct, but does
// not save it to its environment variable.
func SetImmutableBool(k string, v bool) error {
	err := SetSetting("bool", k, v, "", true, false, false, false)
	if err != nil {
		return err
	}

	s := strconv.FormatBool(v)
	return os.Setenv(k,s)
}

// SetImmutableInt adds the information to the AppsConfig struct, but does
// not save it to its environment variable.
func SetImmutableInt(k string, v int) error {
	err := SetSetting("int", k, v, "", true, false, false, false)
	if err != nil {
		return err
	}

	return os.Setenv(k,string(v))
}

// SetImmutableString adds the information to the AppsConfig struct, but does
// not save it to its environment variable.
func SetImmutableString(k, v string) error {
	err := SetSetting("string", k, v, "", true, false, false, false)
	if err != nil {
		return err
	}

	return os.Setenv(k,v)
}
