//
package contour

import (
	"bytes"
	"encoding/json"
	// "encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	// "os"
	// "strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

const app = "app"

const (
	Unsupported Format = iota
	JSON
	TOML
	YAML
	XML
)

type Format int

func (f Format) String() string {
	switch f {
	case JSON:
		return "json"
	case TOML:
		return "toml"
	case YAML:
		return "yaml"
	case XML:
		return "xml"
	default:
		return "unsupported"
	}
}

// isSupported checks to see if the passed string represents a supported config
// format.
func (f Format) isSupported() bool {
	switch f {
	case YAML:
		return true
	case JSON:
		return true
	case TOML:
		return true
	case XML:
		return false
	}
	return false
}

func ParseFormatE(s string) (Format, error) {
	ls := strings.ToLower(s)
	switch ls {
	case "json":
		return JSON, nil
	case "toml":
		return TOML, nil
	case "yaml", "yml":
		return YAML, nil
	case "xml":
		return XML, nil
	}
	return Unsupported, unsupportedFormatErr(ls)
}

func ParseFormat(s string) Format {
	f, _ := ParseFormatE(s)
	return f
}

// appCfg: contour's global config; contour config functinos operate on this.
var appCfg *Cfg

// Contour variable names for the pre-configured core setting names that it
// comes with. These are public and are directly settable if you wish to use
// your own values, just set them before doing anything with Contour.
var (
	CfgFile   string = "cfgfile"
	CfgFormat string = "cfgformat"
)

func init() {
	initCfgs()
}

// initConfigs initializes the configs var. This can be called to reset it in
// testing too.
func initCfgs() {
	appCfg = NewCfg(app)
}

// formatFromFilename gets the format from the passed filename.  An error will
// be returned if either the format isn't supported or the extension doesn't
// exist.  If the passed string has multiple dots, the last dot is assumed to
// be the extension.
func formatFromFilename(s string) (Format, error) {
	if s == "" {
		return Unsupported, fmt.Errorf("no config filename")
	}
	parts := strings.Split(s, ".")
	format := ""
	// case 0 has already been evaluated
	switch len(parts) {
	case 1:
		return Unsupported, fmt.Errorf("unable to determine %s's config format: no extension", strings.TrimSpace(s))
	case 2:
		format = parts[1]
	default:
		// assume its the last part
		format = parts[len(parts)-1]
	}
	f := ParseFormat(format)
	if !f.isSupported() {
		return Unsupported, unsupportedFormatErr(format)
	}
	return f, nil
}

// readCfgFile reads the configFile and returns the resulting slice. The entire
// contents of the file are read at once.
func readCfgFile(n string) ([]byte, error) {
	cfg, err := ioutil.ReadFile(n)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// unmarshalFormatReader accepts an io.Reader and unmarshals it using the
// correct format.
//
// Supported formats:
//   json
//   toml
// TODO
//   add YAML support
//   add HCL support
func unmarshalFormatReader(f Format, r io.Reader) (interface{}, error) {
	b := new(bytes.Buffer)
	b.ReadFrom(r)
	var ret interface{}
	switch f {
	case JSON:
		err := json.Unmarshal(b.Bytes(), &ret)
		if err != nil {
			return nil, err
		}
	case TOML:
		_, err := toml.Decode(b.String(), &ret)
		if err != nil {
			return nil, err
		}
	default:
		err := unsupportedFormatErr(f.String())
		return nil, err
	}
	return ret, nil
}

// notFoundErr returns a standadized notFoundErr.
func notFoundErr(k string) error {
	return fmt.Errorf("not found: %s", k)
}

// settingNotFoundErr standadized settingNotFoundErr.
func settingNotFoundErr(k string) error {
	return fmt.Errorf("setting not found: %s", k)
}

// unsupportedFormatErr standadized the unsupportedFormatErr.
func unsupportedFormatErr(k string) error {
	return fmt.Errorf("unsupported cfg format: %s", k)
}
