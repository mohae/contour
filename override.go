package contour

// Set contains all of contour's Override functions. Override can set both
// settings whose values were obtained from environment variables and regular
// settings that are also Flags. Override cannot set any configuration setting
// that is not a flag. Also, override cannot set any Immutable or IsCore
// settings.
//
// A common use for overrides is to set values obtained by flags.

func Override(k string, v interface{}) error {
	if v == nil {
		return nil
	}
	// If it can't be overriden, return it.
	// This is currently a silent fail.
	// TODO:
	//	log failure
	//	return error instead of silently failing?
	if appConfig.settings[k].IsCore || appConfig.settings[k].Immutable || !appConfig.settings[k].IsFlag {
		return nil
	}

	// Write to environment variable
	err := appConfig.Setenv(k, v)
	if err != nil {
		return err
	}

	appConfig.settings[k].Value = v

	return nil
}
