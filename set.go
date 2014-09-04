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
	


// SetImmutableBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func SetImmutableBool(k string, v bool) {
	SetSetting("bool", k, v, "", true, false, false, false)
	return 
}

// SetImmutableInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func SetImmutableInt(k string, v int) {
	SetSetting("int", k, v, "", true, false, false, false)
	return 
}

// SetImmutableString adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func SetImmutableString(k, v string) {
 	SetSetting("string", k, v, "", true, false, false, false)
	return 
}

// SetImmutableFlagBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func SetImmutableFlagBool(k string, v bool, f string) {
	SetSetting("bool", k, v, f, true, false, false, true)
	return 
}

// SetImmutableFlagInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func SetImmutableFlagInt(k string, v int, f string) {
	SetSetting("int", k, v, f, true, false, false, true)
	return 
}

// SetImmutableFlagString adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func SetImmutableFlagString(k, v, f string) {
 	SetSetting("string", k, v, f, true, false, false, true)
	return 
}

// SetBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func SetBool(k string, v bool) {
	SetSetting("bool", k, v, "", false, false, false, false)
	return 
}

// SetInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func SetInt(k string, v int) {
	SetSetting("int", k, v, "", false, false, false, false)
	return 
}

// SetString adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func SetString(k, v string) {
	SetSetting("string", k, v, "", false, false, false, false)
	return 
}

// SetFlagBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func SetFlagBool(k string, v bool, f string) {
	SetSetting("bool", k, v, f, false, false, false, true)
	return 
}

// SetFlagInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func SetFlagInt(k string, v int, f string) {
	SetSetting("int", k, v, f, false, false, false, true)
	return 
}

// SetFlagString adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func SetFlagString(k, v, f string) {
	SetSetting("string", k, v, f, false, false, false, true)
	return 
}

