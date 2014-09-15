contour
=======
__Under Development__: this includes the readme.

contour: application configuration supporting defaults, configuration files, and flags.

## About
Contour seeks to be a simpe to use configuration package that is flexible and powerful enough to meet more complex uses.

Contour can act like the main configuration by using its functions directly, or you can have any amount of your own custom configurations that your application interacts with directly.

Maps are the core data structure used, which is usefull when _n_ gets large.

To use in a basic application, import the package:

	import "github.com/mohae/contour"

Configuration variables must be registered for Contour to recognize them. Registering a setting lets Contour know what that setting's datatype is, if it can be modified, and if so, by what. 

Supported datatypes:
	* Bool
	* Int
	* Interface{}
	* String

Types of configuration settings:
	* __core__ cannot be updated once registered. These are typically application settings that you wish to make available but do not want anything to be able to overwrite its initial setting, regardless of where the new value is coming from.
	* __conf__ configuration settings are settings that may be set during registration and can be updated by a configuration file. __conf__ settings must be updated using the `Update*` functions.
	* __flag__ configuration settings may be set by the application, updated via configuration files, and are exposed as command-line flags. Registering a setting as a flag automatically makes that setting part of the flag filtering process while processing command-line arguments.
	* __setting__ is just a plain configuration setting. These are not core settings, are not exposed to the configuration file, and are not supported by flags. They are there for you to use within your application however you wish.

Register the configuration variables:

	contour.RegisterCoreString("name", "testapp")
	contour.RegisterConfString("logconfig", "log.cfg")
	contour.RegisterFlagBool("logging", false, "l")
	contour.RegisterSettingInt("ultimateanswer", 42)

Retrieving the set configuration:

	b, err := contour.GetBoolE("logging")
	lcfg := contour.GetString("logconfig")
	
	// Get() returns the interface{} for Value. GetInterface() is an alias
	// for Get()
	if := contour.Get("ultimateanswer")

Anytime `contour` functions are used, like above, the main `app` configuration is used. These functions are convenience functions that call the `app` config for you. Or you can get a pointer to the `app` config and work with it directly:

	cfg := contour.AppConfig()
	cfg.RegisterSettingString("greeting", "hello")

	// Returns "hello": contour.GetString() == cfg.GetString(), in this case.
	s := contour.GetString("greeting")

Or you can get a different configuration by requesting one and providing a key by which the configuration will be identified:

	logCfg := contour.NewConfig("log")

	logCfg.RegisterFlagBool("logging", true, "l")

	lCfg := logCfg.GetBool("logging")
	aCfg := countour.GetBool("logging")

	fmt.Println(lCfg)	// Prints true
	fmt.Println(aCfg)	// Prints false

### Precedence
Contour assu

## Notes:
Currently works but not necessarily safe. Even though each setting has a mutex, support for it has not been implemented. Once implemented, each setting will have RW lock support, allowing for multiple concurrent reads and safe writing.

## Wishlist:
_Environment Variable_ support. The original version of contour supported environment variables, but support for it was not implemented in the initial rewrite, to allow for focus on a cleaner precedence hierarchy and simpler management-though it can still get complicated, if you want. Given the usefulness of having configuration in environment variables for some deployment scenarios, it seems like a good idea to add it back in.

*Slice version*: maps are useful for large numbers, but they are overkill for a large number of applications. For these uses, a slice oriented version, which would be more performant, would be better. 
