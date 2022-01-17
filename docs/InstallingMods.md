# Installing Mods on a Minecraft Client Machine

These instructions assume the Mod Framework and Mod Tool are already installed. For instructions on those topics, see the [Quick Start](https://github.com/effisso/mc-mod-installer/tree/main/docs/QuickStart.md)

`mcmods install --help` is a good resource for a quick overview/refresher of the install command. This document dives a little deeper into some of the specifics. Generally, filtering out any of the server mods should only be done if it's 1: optional, and 2: incompatible with a client-only mod you'd like to use.

**IMPORTANT: If you use a VPN, disconnect it while installing to avoid issues with CloudFlare when downloading mods.**"

The install command is used for initially installing as well as updating mods. The mods it can install fall into two categories (the latter, with a few sub-categories):

* **Client Mods** - custom mods that *only* get installed on the client machine; the server doesn't know or care about these
* **Server Mods** - mods that are installed on the server; the sub-categories are referred to as "server groups"
    * **required** - mods that are necessary for the client to have to connect to the server
    * **server-only** - mods that should *not* be installed on the client machine, even though they are on the server
    * **performance** - optional mods that can be installed for performance boosts and many types of optimizations in Minecraft
    * **optional** - other mods installed on the server that are not necessary on the client, but can be installed, and recommended for client-side responsiveness to server events

To add client-only mods to the tool's configuration, see [Adding Custom Mods](https://github.com/effisso/mc-mod-installer/tree/main/docs/AddingCustomMods.md). Client mods can be installed separately from the server mods (and incrementally over time), so there's no urgency to decide to use these or not.

First-time installations not wanting to replace any of the optional server mods should simply use `mcmods install` to install all the server mods, as well as any client-only mods. Further use cases for this command can be read about below.

## Command Arguments

Brief descriptions of the arguments that can be used to customize the install

* **--client-only** - only install the custom mods defined on this client machine; server mods are ignored
* **--full-server** - should only be used when deploying new server mods via FTP; do not use for a client install
* **--force** - force the mods to be downloaded, even if the latest package already exists locally
* **--x-group** - exclude one or server groups by providing the group names after the flag; comma-separated.
* **--x-mod** - exclude any mod (client or server) from being installed by providing its CLI name following the flag; comma-separated.

## Sample Commands

* `mcmods install` installs all client mods, and all server mods, except the group server-only
* `mcmods install --client-only` installs all custom client-only mods; no server mods
* `mcmods install --client-only --x-mod somemod` excludes a mod from the client-only install
* `mcmods install --x-group performance,optional` exclude the performance and optional server groups (i.e. only install the required mods)
* `mcmods install --force` forces all required, optional, and client-only mods to be redownloaded, even if the latest version exists according to the install config.

**NOTE**: For all install commands that don't explicitly speciy the `--full-server` flag, the `server-only` group is always automatically excluded.
