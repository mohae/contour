contour
=======
Contour helps shape your application by supporting application settings, cfg settings, env settings, and flags with minimal extra work.

## About
Contour attempts to be a simple to use configuration package whose behavior can be modified to suit your application's need.  Contour can work with or without configuration files, tolerate missing config files, and use environment settings.

By default, the contour global cfg, a pointer to which can be obtained via the contour.AppCfg() call, is configured to check environment variables. If a cfg file is successfully set, via RegisterCfgFile(), it will set itself to use a config file.  If any flags are registered, via RegisterFlag*() funcs, flags will be filtered out of passed flags.

For flags, contour supports short flags and will accept either '-' or '--' for both flags and short flags; it is inconsistent with POSIX.

### Supported configuration sources
Contour supports various sources for configuration information:

* Application defaults
* Configuration file formats
  * JSON**
  * TOML
* Environment variables
* CLI flags.

Other than Application defaults, none of the others are required.  Contour can be configured to not error on the absence of a configuration file, if it is set to use one. 

#### JSON support
Standard compliant JSON, as defined in RFC 4627.  It also supports commented json. Even though this will offend some people, comments are useful in configuration files.  Before unmarshaling, these comments are stripped out of the JSON so that the JSON is RFC 4627 compliant.

Both line comments and block comments are supported. Line comments can start with either `//` or `#` and terminate at the end of the current line. Block comments start with `/*` and end with `*/`.  If those characters are found within a key or value, they are ignored.

I use `cjsn`, commented JSON, as the file extension. Since there is no such thing as commented JSON, any file extension, other than `json` or `jsn` will work. To avoid confusion using the standard JSON extensions for JSON files with comments is not recommended.

## Easy to use:
### Import `contour`
To use in a basic application, import the package:

	import "github.com/mohae/contour"

### Register settings
Configuration variables must be registered for Contour to recognize them. Registering a setting lets Contour know what that setting's datatype is, if it can be modified, and if so, by what. 

    contour.RegisterBool(key, value)
    contour.RegisterBoolFlag("log", "l", false, "false", "enable/disable logging")
	
### Initialize the configuration with `InitCfg()`
Once the settings have been registered, the Cfg needs to be initialized. This will merge the cfg file's values, if applicable, and the environement settings, if applicable, with the application defaults. Once initialized, settings can only be updated according to the rules applied to it's setting type.

    err := contour.InitCfg()
	
### Set the Cfg from args: `FilterArgs()`
If CLI flags are used, the command-line args need to be filtered:

    args, err := contour.FilterArgs(args)
	
Any args that are left over, after filtering, are returned. No need to create filter variables.

### Using the Configuration Information
To get the value of the setting, call `Cfg.Get()`, which returns the value as an ``interface{}`. If you want the value returned as its datatype, call that datatype's Get function, e.g. `Cfg.GetBool(name)` for a `bool` setting.

## Working with `Cfg`
`Cfg` is the contour struct for configuration setttings. Contour has a global config that is available using helper funcs. A pointer to the contour global cfg can be obtained by calling `contour.AppCfg()`.

A local Cfg can be obtained by calling `contour.NewCfg(name)`, where `name` is the name of the cfg.GetString.

### Setting a configuration file
A configuration file can be set via the `RegisterCfgFile(fileName)` function. Before setting the Cfg's config filename, if the cfg is set to use env variables, the Cfg will check the `CFGNAME_CFG_FILE` environment variable and see if it is not empty. If the config file environment variable is not empty, it will be used as the application config file, instead of the application default value. 

This allows for the application default cfg file value to be overridden.

One a config file has been registered, Cfg.useCfg will be set to `true`.

## Order of precedence
Except for the configuration file setting, a Cfg applies setting values in the following order, as is applicable:

* Application defaults: these are the values that are used when a setting is registered.
* cfg file: any cfg or flag settings found in the cfg file are updated with the cfg file value.
* env variables: any cfg of flag settings found in the env are updated with the env variable value.
* command-line: any flag settings that are found in the command-line args are updated with the passed value.

### Modifying a Cfg's behavior
#### `Cfg.useEnv`
A Cfg's `useEnv` variable tells the cfg whether or not environment variables should be used. If the Cfg supports environment variables, the environment will be checked for all cfg and flag settings. 

Contour generates a setting's environment variable name by prefixing the setting name with the Cfg's name and an underscore,  `_`, separator. The resulting environment setting name will also be in all UPPERCASE, regardless of the original casing.
    
Set whether or not to use environment variables:
    SetUseEnv(bool)
	myCfg.SetUseEnv(bool)
	
Get the current value of `useEnv`:
    b := UseEnv()
	b := myCfg.UseEnv()
	
#### `Cfg.useCfg`
A Cfg's `useCfg` variable tells the cfg whether or not to check for a configuration file. If it is `true`, the value of the `cfg_file` setting will be used as the location of the configuration file. Any cfg and flag settings which are found in the config file will be updated with the file's values. 

The setting names in the cfg file will be the same as the setting's for which they apply.

The `useCfg` flag is set automatically when `RegisterCfgFile()` is called.

To get the current value of `useCfg`:
    b := UseCfg()
    b := myCfg.UseCfg()

## Cfg settings
A contour setting is the basic datatype for a setting in contour. It's values are largely determined by how they are registered. Depending on the setting type, their values may, or may not, be updateable. 

Setting names must be unique and their short code must also be unique.

### supported datatypes
Currently, only the following datatypes are supported:
	* bool
	* int
	* int64
	* interface{}
	* string

For method sets that include datatype support, e.g. `Get*`, the version without a datatype, e.g. `Get()`, will return an interface{}.

### Registering settings
To use a setting, it must be first registered. Registering the setting ensures that everything gets set properly for a given setting type. What a setting is registered as will affect what operations are allowed on it after registration, e.g. __core__ setting values cannot be changed after registration.

All registration functions are in the form of `Register{datatype}{setting-type}`, with the __setting-type__ being optional.

#### Core settings
Core settings are settings, which once registered, cannot be modified: 
	contour.RegisterStringCore("corestring", "this is a core setting and cannot be changed after registration")
	
Cfg settings are settings that are exposed as configuration settings. Only cfg and flag settings can be in config files. If the Cfg supports environment variables, cfg settings are also exposed as environment variables. Once a cfg setting has been registered, they can only be updated using `Update` functions.
	contour.RegisterBoolCfg("cfgbool", false)
	
Flag settings are settings that are available as command-line flags. Flags are also exposed as configuration file settings and environment variables. Flags can be updated using the `Update` functions until `FilterArgs()` has been done. Once the args have been filtered, any changes to their value must be done with the `Override()` function.
	contour.RegisterStringFlag("flagstring", "this is a flag setting and can be used as a command-line flag")

Regular settings are settings that are not exposed in the config file, enviornment variables, or command-line args. Unlike Core settings, regular settings are modifiable.
	contour.RegisterInt("int", 42)

