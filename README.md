contour
=======
[![GoDoc](https://godoc.org/github.com/mohae/contour?status.svg)](https://godoc.org/github.com/mohae/contour)[![Build Status](https://travis-ci.org/mohae/contour.png)](https://travis-ci.org/mohae/contour)

Contour is a package for storing settings. These settings may be configuration settings or application settings.

## About
Contour attempts to be a simple to use package that can be used for key value pairs or configuration information and whose behavior can be modified to suit your application's need.  Contour can work with or without configuration files, tolerate missing config files, use environment settings, and supports flags.

Each key value pair is saved as a setting, which also contains information about that pair that enables Contour to figure out what it can do with each setting. Groups of settings are stored in a Settings struct, which also contains information about that group, including its name, and how it should behave. Contour provides a package global Settings with its name set to the executable name. All Contour functions operate on this Settings. For your own Settings, use NewSettings.

All Settings operations are thread-safe.

## Settings

### Application Settings
Application settings are any setting that cannot be modified by a configuration file, environment variable, or flag; they either are not updateable, Core settings, or only updateable by calling an Update function, everything else. These settings are Added to Contour via Add functions. Settings, whose values can be modified, are changed using an Update function.

### Configuration Settings
Configuration settings are any settings that can be modified by a configuration file, environment variable, or a flag, in that order of precedence. These settings are Registered using Register functions. They have a default value that can be set by one or more of the above methods, depending on what is allowed to update them.

The update rules for a configuration setting is at a per setting level. Setting Foo can be registered as a ConfigFileVar, which means it can only be updated by a configuration file, while Bar can be registered as a Flag, which means it can be updated by a configuration file, environment variable, or a flag. Custom update flags can be set: e.g. Biz can be set to be updateable by a configuration file or a flag but not an environment variable.

Flags can be registered with a short flag, or alias.

A Settings can be configured to search the PATH for the configuration file. It can also be set to not error when the configuration file is missing.

Contour supports various formats for configuration files:

* JSON
* TOML
* YAML


## Easy to use:
### Import `contour`
To use in a basic application, import the package:

	import "github.com/mohae/contour"

### Register settings
Configuration variables must be registered for Contour to recognize them. Registering a setting lets Contour know what that setting's datatype is, if it can be modified, and if so, by what.

    contour.AddString(key, value)
    contour.RegisterStringConfFileVar("foo", "bar")
    contour.RegisterBoolFlag("log", "l", false, "false", "enable/disable logging")

### Initialize the configuration
Once all settings have been registered, `Set` needs to be run to update the settings with all available configuration file settings and environment variables.

    err := contour.Set()

### Parse flags
If flags are used, the command-line args need to be parsed for flags:

    args, err := contour.ParseFlags(args)

Any args that are left over, after parsing, are returned. A list of all flags parsed is maintained and either the whole list can be retrieved:

    flgs := contour.Visited()

or a specific flag can be checked:

    was := contour.WasVisited("foo")

### supported datatypes
Currently, only the following datatypes are supported:
	* bool
	* int
	* int64
	* interface{}
	* string
