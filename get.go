package contour


// E versions, these return the error. Non-e versions are just wrapped calls to
// these functions with the error dropped.

// Get returns the setting Value as an interface{}.
func Get(k string) interface{} {
	tmp, _ := GetE(k)
	return tmp
}

// GetBool returns the setting Value as a bool.
func GetBool(k string) bool {
	tmp, _ := GetBoolE(k)
	return tmp
}

// GetInt returns the setting Value as an int.
func GetInt(k string) int {
	tmp, _ := GetIntE(k)
	return tmp
}

// GetString returns the setting Value as a string.
func GetString(k string) string {
	tmp, _ := GetStringE(k)
	return tmp
}

// GetInterface is a convenience wrapper function to Get
func GetInterface(k string) interface{} {	
	return Get(k)
}

// GetBoolFilterNames returns a list of filter names (flags).
func GetBoolFilterNames() []string {
	var names []string

	for k, setting := range appConfig.settings {
		if setting.IsFlag && setting.Type == "bool" {
			names = append(names, k)
		}
	}

	return names
}

// GetIntFilterNames returns a list of filter names (flags).
func GetIntFilterNames() []string {
	var names []string

	for k, setting := range appConfig.settings {
		if setting.IsFlag && setting.Type == "int" {
			names = append(names, k)
		}
	}

	return names
}

// GetStringFilterNames returns a list of filter names (flags).
func GetStringFilterNames() []string {
	var names []string

	for k, setting := range appConfig.settings {
		if setting.IsFlag && setting.Type == "string" {
			names = append(names, k)
		}
	}

	return names
}

// GetE returns the setting Value as an interface{}.
func GetE(k string) (interface{}, error) {
	_, ok := appConfig.settings[k]
	if !ok {
		return nil, notFoundErr(k)
	}

	return appConfig.settings[k].Value, nil
}

// GetBoolE returns the setting Value as a bool.
func GetBoolE(k string) (bool, error) {
	_, ok := appConfig.settings[k]
	if !ok {
		return false, notFoundErr(k)
	}

	return *appConfig.settings[k].Value.(*bool), nil
}

// GetIntE returns the setting Value as an int.
func GetIntE(k string) (int, error) {
	_, ok := appConfig.settings[k]
	if !ok {
		return 0, notFoundErr(k)
	}

	return *appConfig.settings[k].Value.(*int), nil
}

// GetStringE returns the setting Value as a string.
func GetStringE(k string) (string, error) {
	_, ok := appConfig.settings[k]
	if !ok {
		return "", notFoundErr(k)
	}

	return appConfig.settings[k].Value.(string), nil
}

// GetInterfaceE is a convenience wrapper function to Get
func GetInterfaceE(k string) (interface{}, error) {
	return GetE(k)
}
