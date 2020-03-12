# Using the version command 
You can use `hpecli --help` to find out about the usage of the command. You can get the version of the hpecli with `hpecli version`. You can specify a `verbose` flag to get more details. 

<img src="hpecli-version-new.gif" align="center">

# Using the update command
If a newer version is available, which we check when you invoke `hpecli version`, we can use `hpecli update` to download the latest version.

<img src="hpecli-update-new.gif" align="center">

# Using oneview commands
There is a plugin for HPE OneView which allows to control one or more HPE OneView appliance using the `hpecli oneview` command set. Run `hpe oneview --help` to find out the supported verbs. Then we can use `hpecli oneview login`to connect to an HPE OneView instance then `hpecli oneview get servers` or `hpecli oneview get enclosures` to retrieve information from HPE OneView. When done we can use `hpecli oneview logout` to terminate the session.

<img src="hpecli-oneview-new.gif" align="center">

# Using context
We cache the latest context used for a given plugin. But you can switch to a different context using `hpe oneview context --host <anotherhost>`. You can also find out which context is currently active with `hpecli context`.

<img src="hpecli-context-new.gif" align="center">

# Using iLO/Redfish commands
There is a plugin for HPE iLO/Redfish which allows to control one or more HPE iLO BMC using the `hpecli ilo` command set. Run `hpe ilo --help` to find out the supported verbs. Then we can use `hpecli ilo login`to connect to an HPE iLO instance then `hpecli ilo get serviceroot` to retrieve root information about this iLO/Redfish instance. When done we can use `hpecli ilo logout` to terminate the session. You can use `hpecli ilo context` to find out which iLO you are opererating and `hpecli ilo context --host <anotherhost>` to switch context.

<img src="hpecli-ilo-new.gif" align="center">
