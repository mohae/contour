package contour
// Set contains all of contour's Set functions.Calling Set
// adds, or registers, the settings information to the AppConfig variable.
// The setting value, if there is one, is not saved to its environment 
// variable at this point.
//
// This allows for 
//
// These should be called at app startup to register all configuration
// settings that the application uses.

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

