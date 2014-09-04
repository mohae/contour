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
func RegisterSetting( Type string, k string, v interface{}, Code string, Immutable, IsCore, IsEnv, IsFlag bool) {
	var update bool
	_, ok := AppConfig.Settings[k]
	if ok {

		// Core settings can't be re-registered.
		if AppConfig.Settings[k].IsCore {
			return
		}
	
		// Read-only settings that have bee set can't be re-registered.
		if AppConfig.Settings[k].Immutable {

			if AppConfig.Settings[k].Value != nil {
				return
			}

			update = true
			
		}

	}

	if update {
		AppConfig.Settings[k].Type = Type
		AppConfig.Settings[k].Value = v
		AppConfig.Settings[k].Code = Code
		AppConfig.Settings[k].Immutable = Immutable
		AppConfig.Settings[k].IsCore = IsCore
		AppConfig.Settings[k].IsEnv = IsEnv
		AppConfig.Settings[k].IsFlag = IsFlag
		return  
	}

	AppConfig.Settings[k] = &setting{
		Type: Type,
		Value: v,
		Code: Code,
		Immutable: Immutable,
		IsCore: IsCore,
		IsEnv: IsEnv,
		IsFlag: IsFlag}
}
	
// RegisterCoreBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable
func RegisterCoreBool(k string, v bool) {
	RegisterSetting("bool", k, v, "", true, true, false, false)
	return
}

// RegisterCoreInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable
func RegisterCoreInt(k string, v int) {
	RegisterSetting("int", k, v, "", true, true, false, false)
	return
}

// RegisterCoreString adds the information to the AppsConfig struct, but does not
// save it to its environment variable
func RegisterCoreString(k, v string) {
	RegisterSetting("string", k, v, "", true, true, false, false)
	return
}

// RegisterCoreFlagBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable
func RegisterCoreFlagBool(k string, v bool, f string) {
	RegisterSetting("bool", k, v, f, true, true, false, true)
	return
}

// RegisterCoreInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable
func RegisterCoreFlagInt(k string, v int, f string) {
	RegisterSetting("int", k, v, f, true, true, false, true)
	return
}

// RegisterCoreString adds the information to the AppsConfig struct, but does not
// save it to its environment variable
func RegisterCoreFlagString(k, v, f string) {
	RegisterSetting("string", k, v, f, true, true, false, true)
	return
}


// RegisterImmutableBool adds the information to the AppsConfig struct, but
// does not save it to its environment variable.
func RegisterImmutableBool(k string, v bool) {
	RegisterSetting("bool", k, v, "", true,  false, false, false)
	return
}

// RegisterImmutableInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func RegisterImmutableInt(k string, v int) {
	RegisterSetting("int", k, v, "", true,  false, false, false)
	return
}


// RegisterROString adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func RegisterImmutableString(k, v string) {
 	RegisterSetting("string", k, v, "", true,  false, false, false)
	return
}

// RegisterImmutableFlagBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func RegisterImmutableFlagBool(k string, v bool, f string) {
	RegisterSetting("bool", k, v, f, true, false, false, true)
	return
}

// RegisterImmutableFlagInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func RegisterImmutableFlagInt(k string, v int, f string) {
	RegisterSetting("int", k, v, f, true, false, false, true)
	return
}

// RegisterImmutableFlagString adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func RegisterImmutableFlagString(k, v, f string) {
 	RegisterSetting("string", k, v, f, true, false, false, true)
	return
}

// RegisterBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func RegisterBool(k string, v bool) {
	RegisterSetting("bool", k, v, "", false, false, false, false)
	return
}

// RegisterInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func RegisterInt(k string, v int) {
	RegisterSetting("int", k, v, "", false, false, false, false)
	return
}

// RegisterString adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func RegisterString(k, v string) {
	RegisterSetting("string", k, v, "", false, false, false, false)
	return
}

// RegisterBoolFlag adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func RegisterBoolFlag(k string, v bool, f string) {
	RegisterSetting("bool", k, v, f, false, false, false, true)
	return
}

// RegisterIntFlag adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func RegisterIntFlag(k string, v int, f string) {
	RegisterSetting("int", k, v, f, false, false, false, true)
	return
}

// RegisterStringFlag adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func RegisterStringFlag(k, v, f string) {
	RegisterSetting(k, v, "string", f, false, false, false, true)
	return
}

