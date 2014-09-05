package contour

import (
	"os"
	"strconv"
)

// UpdateBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateBool(k string, v bool) error {
	AppConfig.Settings[k].Value = v
	s := strconv.FormatBool(v)
	err := os.Setenv(k,s)
	return err 
}

// UpdateInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateInt(k string, v int) error {
	AppConfig.Settings[k].Value = v
	return os.Setenv(k,string(v))
}

// UpdateString adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateString(k, v string) error {
	AppConfig.Settings[k].Value = v
	return os.Setenv(k,v)

}

