package contour

// UpdateBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateBool(k string, v bool) error {
	appConfig.settings[k].Value = v
	return  appConfig.Setenv(k,v)
}

// UpdateInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateInt(k string, v int) error {
	appConfig.settings[k].Value = v
	return appConfig.Setenv(k,v)
}

// UpdateString adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateString(k, v string) error {
	appConfig.settings[k].Value = v
	return  appConfig.Setenv(k,v)

}
