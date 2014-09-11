package contour

// Set contains all of contour's Override functions. Override can set both
// settings whose values were obtained from environment variables and regular
// settings that are also Flags. Override cannot set any configuration setting
// that is not a flag. Also, override cannot set any Immutable or IsCore
// settings.
//
// A common user for overrides is to set values obtained by flags.
import (
	"fmt"
	"os"
	"strconv"
)

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

	// write it to its environment variable
	var tmp string
	var err error

	switch appConfig.settings[k].Type {
	case "string":
		err = os.Setenv(k, *v.(*string))

	case "int":
		err = os.Setenv(k, string(*v.(*int)))

	case "bool":
		tmp = strconv.FormatBool(*v.(*bool))
		err = os.Setenv(k, tmp)

	default:
		err = fmt.Errorf("Unable to override setting %s: type is unsupported %s", k, appConfig.settings[k].Type)
	}

	if err != nil {
		return err
	}

	appConfig.settings[k].Value = v

	return nil
}
