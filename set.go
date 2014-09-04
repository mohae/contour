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
	
	utils "github.com/mohae/utilitybelt"
)

func SetEnvs() error {
	var err error
	// For each setting
	for k, setting := range AppConfig.Settings {
		switch setting.Type {
		case "bool":
			err = os.Setenv(k, utils.BoolToString(setting.Value.(bool)))
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
	case "bool", "int", "string":
		err = os.Setenv(k,v.(string))		
	default:
		err = errors.New(k + "'s datatype, " + AppConfig.Settings[k].Type + ", is not supported")
	}

	return err

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
	


// SetBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func SetBool(k string, v bool) error {
	err := SetSetting("bool", k, v, "", false, false, false, false)
	if err != nil {
		return err
	}

	s := utils.BoolToString(v)
	err = os.Setenv(k,s)
	return err 
}

// SetInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func SetInt(k string, v int) error {
	err := SetSetting("int", k, v, "", false, false, false, false)
	if err != nil {
		return err
	}

	err = os.Setenv(k,string(v))
	return err 
}

// SetString adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func SetString(k, v string) error {
	err := SetSetting("string", k, v, "", false, false, false, false)
	if err != nil {
		return err
	}

	err = os.Setenv(k,v)
	return err 
}

// SetFlagBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func SetFlagBool(k string, v bool, f string) error {
	err := SetSetting("bool", k, v, f, false, false, false, true)
	if err != nil {
		return err
	}

	s := utils.BoolToString(v)
	err = os.Setenv(k,s)
	return err 
}

// SetFlagInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func SetFlagInt(k string, v int, f string) error {
	err := SetSetting("int", k, v, f, false, false, false, true)
	if err != nil {
		return err
	}

	err = os.Setenv(k,string(v))
	return err 
}

// SetFlagString adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func SetFlagString(k, v, f string) error {
	err := SetSetting("string", k, v, f, false, false, false, true)
	if err != nil {
		return err
	}

	err = os.Setenv(k,v)
	return err 
}

// SetImmutableBool adds the information to the AppsConfig struct, but does
// not save it to its environment variable.
func SetImmutableBool(k string, v bool) error {
	err := SetSetting("bool", k, v, "", true, false, false, false)
	if err != nil {
		return err
	}

	s := utils.BoolToString(v)
	err = os.Setenv(k,s)
	return err 
}

// SetImmutableInt adds the information to the AppsConfig struct, but does
// not save it to its environment variable.
func SetImmutableInt(k string, v int) error {
	err := SetSetting("int", k, v, "", true, false, false, false)
	if err != nil {
		return err
	}

	err = os.Setenv(k,string(v))
	return err 
}

// SetImmutableString adds the information to the AppsConfig struct, but does
// not save it to its environment variable.
func SetImmutableString(k, v string) error {
	err := SetSetting("string", k, v, "", true, false, false, false)
	if err != nil {
		return err
	}

	err = os.Setenv(k,v)
	return err 
}


