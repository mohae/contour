package contour

// Get returns the setting Value as an interface{}.
func Get(k string) (interface{}, error) {
	_, ok := AppConfig.Settings[k]
	if !ok {
		return nil, notFoundErr(k)
	}

	return AppConfig.Settings[k].Value, nil
}

// GetBool returns the setting Value as a bool.
func GetBool(k string) (bool, error) {
	_, ok := AppConfig.Settings[k]
	if !ok {
		return false, notFoundErr(k)
	}

	return AppConfig.Settings[k].Value.(bool), nil
}

// GetInt returns the setting Value as an int.
func GetInt(k string) (int, error) {
	_, ok := AppConfig.Settings[k]
	if !ok {
		return 0, notFoundErr(k)
	}

	return AppConfig.Settings[k].Value.(int), nil
}

// GetString returns the setting Value as a string.
func GetString(k string) (string, error) {
	_, ok := AppConfig.Settings[k]
	if !ok {
		return "", notFoundErr(k)
	}

	return AppConfig.Settings[k].Value.(string), nil
}

// GetInterface is a convenience wrapper function to Get
func GetInterface(k string) (interface{}, error) {
	return Get(k)
}

// GetBoolFilterNames returns a list of filter names (flags).
func GetBoolFilterNames() []string {
	var names []string

	for k, setting := range AppConfig.Settings {
		if setting.IsFlag && setting.Type == "bool" {
			names = append(names, k)
		}
	}

	return names
}

// GetIntFilterNames returns a list of filter names (flags).
func GetIntFilterNames() []string {
	var names []string

	for k, setting := range AppConfig.Settings {
		if setting.IsFlag && setting.Type == "int" {
			names = append(names, k)
		}
	}

	return names
}

// GetStringFilterNames returns a list of filter names (flags).
func GetStringFilterNames() []string {
	var names []string

	for k, setting := range AppConfig.Settings {
		if setting.IsFlag && setting.Type == "string" {
			names = append(names, k)
		}
	}

	return names
}
