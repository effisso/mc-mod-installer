# Printing Tool Info

## Tool Version

`mcmods version`

## Minecraft Install Path

This is the path the tool assumes exists and has Minecraft installed in it.

`mcmods mcpath`

To update the path the tool uses, use the `--set` flag

`mcmods mcpath --set C:\path\to\.minecraft`

## List Valid Server Groups

Mods on the server fall into one of these four groups:

* required
* server-only
* performance
* optional

The names above are used verbatim in the tool for commands referencing the groups. To print these group names in the terminal:

`mcmods list groups`

## Describe Server Group

Further descriptions of each group can be printed by calling `mcmods describe group x` where `x` is the name above (see also [Installing Mods](https://github.com/effisso/mc-mod-installer/tree/main/docs/InstallingMods.md)).

## List Mods

To list various sets of mods, use the `mcmods list mods` command. Mods are displayed by their CLI name, which is the name used to refer to the mod in other operations with this tool. Here are the flags/filters for this command:

* `--installed` - show installed mods
* `--not-installed` - show mods that are *not* installed
* `--client` - show mods that are client-only
* `--server` - show that come from the server
* `--group` - show mods from the provided server group name

### Examples

All of the filters can be used in conjunction. Below are several examples of the command and the explanations of what they do.

* `mcmods list mods --client` - shows all client-only mods which are defined in the tool's local configuration
* `mcmods list mods --client --installed` - shows all client-only mods which are currently installed on the machine
* `mcmods list mods --group required --not-installed` - shows all mods required by the server which are not installed on the machine
* `mcmods list mods --server` - shows all mods on the server

## Describe Mod

Additional details about each mod are available by calling `mcmods describe mod x` where `x` is the mod CLI name.

## Describe Mod Install

Print the details of a mod's installation metadata: `mcmods describe installation x` where `x` is the mod CLI name.

## Visiting a Mod's homepage/wiki

To make learning accessing mod documentation easier, the tool can be used to quickly visit the main informational webpage about each mod. Just use the command `mcmods visit x` where `x` is the mod CLI name.

## Accessing Tool Documentation

`mcmods docs` will open up the main doc folder in a browser.