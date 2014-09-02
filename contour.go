//
package contour

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"strings"

	"github.com/BurntSushi/toml"
)

// appCode is your short code for your application, if you use one. This is 
// used to prefix the environment variable name. This can be left empty.
var appCode string

// configFilename is the name of the configuration file.
var configFilename string

// configFormat is the format for the configuration file
var configFormat string

// configFile holds the contents of the configuration file
var configFile map[string]interface{} = make(map[string]interface{})


// setAppCode sets the application code.
func SetAppCode(code string) error {
	if code == "" {
		return errors.New("code expected, none received")
	}
	
	appCode = code

	return nil
}

// setConfigFile sets the configuration filename
func SetConfigFile(name string) error {
	if name == "" {
		return errors.New("filename expected, none received")
	}

	configFilename = name

	return nil
}

// SetConfigFormat exposes the configFormat variable. Use this to explicitely
// set the format type. setConfigFormat will not override this value.
// The set format must be a supported config file format.
func SetConfigFormat(s string) error {
	if s == "" {
		return errors.New("config format was expected, none received")
	}

	err := isSupportedFormat(s)
	if err != nil {
		return err
	}

	configFormat = s

	return nil
}

// setConfigFormat parses the configFilename to determine the format being
// used. If the format cannot be determined, or is not supported, an error
// is returned
func setConfigFormat() error {
	// If the format is already set, we don't override the setting.
	// A nil is returned because this is not an error.
	if configFormat != "" {
		return nil
	}	
	
	parts := strings.Split(configFilename, ".")

	switch len(parts) {
	case 0:
		return errors.New("unable to determine config format, filename not set")
	case 1:
		return errors.New("unable to determine config format, the configuration file " + configFilename + " doesn't have an extension")
	case 2:
		configFormat = parts[1]
	default:
		// assume its the last part
		configFormat = parts[len(parts) - 1]
	}

	return nil

}

// isSupportedFormat checks to see if the passed string represents a supported
// config format. If it is, it returns a nil, otherwise an error.
func isSupportedFormat(s string) error {
        switch s {
        case "json", "toml":
                configFormat = s
        default:
                err := errors.New(s + " is not a supported configuration format")
                return  err
        }

        return nil
}

// LoadConfigFile() is the entry point for reading the configuration file.
func LoadConfigFile() error {
	if configFilename == "" {
		return errors.New("config filename not set")
	}

	fBytes, err := readConfigFile()
	if err != nil {
		return err
	}

	err = MarshalFormatReader(configFormat,bytes.NewReader(fBytes)) 
	if err != nil {
		return err
	}

	return nil
}

// readConfigFile reads the configFile
func readConfigFile() ([]byte, error) {
	cfg, err := ioutil.ReadFile(configFilename)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// MarshalFormatReader 
func MarshalFormatReader(t string, r io.Reader) error {
	b := new(bytes.Buffer)
	b.ReadFrom(r)

	switch t{
	case "json":
		err := json.Unmarshal(b.Bytes(), &configFile)
		if err != nil {
			return err
		}

	case "toml":
		_, err := toml.Decode(b.String(), &configFile)
		if err != nil {
			return err
		}

	}
	return nil
}
