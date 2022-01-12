# Adding Custom Mods

In addition to installing the server mods on a client machine, the tool can also be used to install custom client-only mods. The first step is adding a definition of the mod to the installer. Run the `mcmods add` command. There will be a series of prompts asking for input about the new mod. The collected information is described below:

* **Friendly Name** - human-readable name for this mod
* **CLI Name** - the short, yet concise name for the mod; ; lowercase and hyphens only
* **Description** - a description of the mod beyond what's implied by the friendly name; optional
* **Homepage/Wiki URL** - the main webpage for information about this mod
* **Package Download URL** - the HTTP URL for the version of the package to install

Once all prompts have been answered, the new mod configuration is saved to the local configuration file. To install the mod, just install all client-only mods with the command `mcmods install --client-only`. Without the --force flag, only mods not currently installed with the latest version will be downloaded.