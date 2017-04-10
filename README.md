contour
=======
[![GoDoc](https://godoc.org/github.com/mohae/contour?status.svg)](https://godoc.org/github.com/mohae/contour)[![Build Status](https://travis-ci.org/mohae/contour.png)](https://travis-ci.org/mohae/contour)

Contour is a configuration handling package supporting files (JSON, YAML, TOML), environment variables, and flags.

## About
Contour attempts to be a simple to use package that can be used for key value pairs or configuration information and whose behavior can be modified to suit your application's need.  Contour can work with or without a configuration file, use environment variables, and supports flags.

Each key value pair is saved as a setting, which also contains information about that pair that enables Contour to figure out what it can do with each setting. Groups of settings are stored in a Settings struct, which also contains information about that group, including its name, and how it should behave. For convenience, Contour provides a standard Settings with its name set to the executable name. All Contour functions operate on this Settings. For your own Settings, use NewSettings.

Settings can be used concurrently.

## Setting
A setting is a key value pair.

### Application Setting
An application setting is any setting that cannot be modified by a configuration file, environment variable, or flag. A Core setting cannot be updated once added. Any attempt to modify a Core setting will result in an error. A regular application setting can only be updated using an Update method; a configuration file, environment variable, or flag cannot modify an application setting, they are not exposed outside of the application. These settings are added to a setting via Add methods. Settings, whose values can be modified, are changed using an Update method.

### Configuration Settings
Configuration settings are any settings that can be modified by a configuration file, environment variable, or flag, in that order of precedence. These settings are Registered using Register methods. They have a default value that can be set by one or more of the above methods, depending on what is allowed to update them.

The update rules for a configuration setting is at a per setting level. Setting Foo can be registered as a ConfFileVar, which means it can only be updated by a configuration file, while Bar can be registered as a Flag, which means it can be updated by a configuration file, environment variable, or flag. What is allowed up update a configuration setting can be set per setting, e.g. flag Biz can be set to be updateable by a configuration file or a flag but not an environment variable.

Flags can be registered with a short flag, or alias.

### Configuration files
Contour supports various formats for configuration files:

* JSON (default)
* TOML
* YAML

If a configuration file has not been explicitly set, the settings will use its name as the filename and the format it has been set to use as the file's extension. Contour will search for the configuration file using the filename, using any additional paths and environment variables it has been given along with the working directory, executable directory, and $PATH. The search behavior is fully configurable.

Currently, only the top level keys of the configuration file are parsed with keys whose values are not `bool`, `int`, `int64`, or a `string` being saved as an interface{}.

If the configuration file is optional, settings can be set to not emit an error when it can't find it.

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
If flags are used, the command-line args need to be parsed for flags. The standard logger uses `os.Args[1:]`:

    args, err := contour.ParseFlags()

 All other settings must have the args passed:

 		s := NewSettings("foo")
		args, eerr := s.ParseFlags(args)

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
