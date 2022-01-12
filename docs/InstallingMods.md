# Installing Mods on a Minecraft Client Machine

These instructions assume the Mod Framework and Mod Tool are already installed. For instructions on those topics, see the [Quick Start](https://github.com/effisso/mc-mod-installer/tree/main/docs/QuickStart.md)

`mcmods install --help` is a good resource for a quick overview/refresher of the install command. This document dives a little deeper into some of the specifics.

The install command is used for initially installing as well as updating mods. The mods it can install fall into two categories (the latter, with a few sub-categories):

* **Client Mods** - custom mods that *only* get installed on the client machine; the server doesn't know or care about these
* **Server Mods** - mods that are installed on the server; the sub-categories are referred to as "server groups"
    * **required** - mods that are necessary for the client to have to connect to the server
    * **server-only** - mods that should *not* be installed on the client machine, even though they are on the server
    * **performance** - optional mods that can be installed for performance boosts and many types of optimizations in Minecraft
    * **optional** - other mods installed on the server that are not necessary on the client, but can be installed, and recommended for client-side responsiveness to server events

To add client-only mods to the tool's configuration, see [Adding Client Mods](https://github.com/effisso/mc-mod-installer/tree/main/docs/AddingClientMods.md).

## Command Arguments

Brief descriptions of the arguments that can be used to customize the install

* **--client-only** - only install the custom mods defined on this client machine
* **--full-server** - do not use for a client install
* **--force** - force the mod to be redownloaded and written to disk, even if the latest already exists on this machine
* **--x-group** - exclude one or server groups by providing the group names after the flag; comma-separated
* **--x-mod** - exclude any mod (client or server) from being installed by providing its CLI name following the flag; comma-separated

## Sample Commands

* `mcmods install` installs all client mods, and all server mods, except the group server-only
* `mcmods install --client-only` installs all custom client-only mods; no server mods
* `mcmods install --client-only --x-mod somemod` excludes a mod from the client-only install
* `mcmods install --x-group performance,optional` exclude the performance and optional server groups

**NOTE**: For all install commands that don't explicitly speciy the `--full-server` flag, the `server-only` group is always automatically excluded.
