package contour
// Register contains all of contour's Register functions.Calling Register
// adds, or registers, the settings information to the AppConfig variable.
// The setting value, if there is one, is not saved to its environment 
// variable at this point.
//
// This allows for 
//
// These should be called at app startup to register all configuration
// settings that the application uses.

// RegisterSetting checks to see if the entry already exists and adds the
// new setting if it does not.
func RegisterSetting(k string, v interface{}, Type, Code, string, IsFlag, IsRO, IsEnv, IsCore bool) error {
	_, ok := AppConfig.Settings[k]
	if ok {
		return
	}

	AppConfig.Settings[k] = &setting{Value: v, Type: Type, Code: Code, IsFlag: IsFlag, IsRO: IsRO, IsEnv: IsEnv, IsCore: IsCore}
}
	
// RegisterCoreBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable
func RegisterCoreBool(k string, v bool) {
	return RegisterSetting(k, v, "bool", "",false, true, false, true)
}

// RegisterCoreInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable
func RegisterCoreInt(k string, v int) {
	return RegisterSetting(k, v, "int", "", false, true, false, true)
}

// RegisterCoreString adds the information to the AppsConfig struct, but does not
// save it to its environment variable
func RegisterCoreString(k, v string) {
	return RegisterSetting(k, v, "string", "", false, true, false, true)
}

// RegisterCoreFlagBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable
func RegisterCoreFlagBool(k string, v bool, f string) {
	return RegisterSetting(k, v, "bool", f, true, true, false, true)
}

// RegisterCoreInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable
func RegisterCoreFlagInt(k string, v int, f string) {
	return RegisterSetting(k, v, "int", f, true, true, false, true)
}

// RegisterCoreString adds the information to the AppsConfig struct, but does not
// save it to its environment variable
func RegisterCoreFlagString(k, v, f string) {
	return RegisterSetting(k, v, "string", f, true, true, false, true)
}


// RegisterROBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func RegisterROBool(k string, v bool) {
	return RegisterSetting(k, v, "bool", "", false, true, false, false)
}

// RegisterReadOnlyBool is an alias for RegisterROBool
func RegisterReadOnlyBool(k string, v bool) {
	return RegisterROBool(k, v)
}

// RegisterROInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func RegisterROInt(k string, v int) {
	return RegisterSetting(k, v, "int", "", false, true, false, false)
}

// RegisterReadOnlyInt is an alias for RegisterROInt
func RegisterReadOnlyInt(k string, v int) {
	return RegisterROInt(k, v)
}

// RegisterROString adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func RegisterROString(k, v string) error {
 	return RegisterSetting(k, v, "string", "", false, true, false, false)
}

// RegisterReadOnlyString is an alias for RegisterROString
func RegisterReadOnlyString(k, v string) {
	return RegisterROString(k, v)
}

// RegisterROFlagBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func RegisterROFlagBool(k string, v bool, f string) {
	return RegisterSetting(k, v, "bool", f, true, true, false, false)
}

// RegisterReadOnlyFlagBool is an alias for RegisterROBool
func RegisterReadOnlyFlagBool(k string, v bool, f string) {
	return RegisterROFlagBool(k, v, f)
}

// RegisterROFlagInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func RegisterROFlagInt(k string, v int, f string) {
	return RegisterSetting(k, v, "int", f, true, true, false, false)
}

// RegisterReadOnlyFlagInt is an alias for RegisterROInt
func RegisterReadOnlyFlagInt(k string, v int, f string) {
	return RegisterROInt(k, v. f)
}

// RegisterROFlagString adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func RegisterROFlagString(k, v, f string) error {
 	return RegisterSetting(k, v, "string", f, true, true, false, false)
}

// RegisterReadOnlyFlagString is an alias for RegisterROString
func RegisterReadOnlyFlagString(k, v, f string) {
	return RegisterROString(k, v, f)
}

// RegisterBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func RegisterBool(k string, v bool) error {
	return RegisterSetting(k, v, "", true, false, false, false)
}

// RegisterInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func RegisterInt(k string, v int) error {
	return RegisterSetting(k, v, "", true, false, false, false)
}

// RegisterString adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func RegisterString(k, v string) error {
	return RegisterSetting(k, v, "", true, false, false, false)
}

// RegisterFlagBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func RegisterFlagBool(k string, v bool, f string) error {
	return RegisterSetting(k, v, "bool", f, true, false, false, false)
}

// RegisterFlagInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func RegisterFlagInt(k string, v int, f string) error {
	return RegisterSetting(k, v, "int", f, true, false, false, false)
}

// RegisterFlagString adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func RegisterFlagString(k, v, f string) error {
	return RegisterSetting(k, v, f, true, false, false, false)
}

