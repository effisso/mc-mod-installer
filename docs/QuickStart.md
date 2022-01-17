# Quick Start Guide

The following sections demonstrate how to get started with the tool if followed in order.

## Installing the Fabric Mod Framework

Clients connecting to YAMS must install the Fabric mod framework on their machine. 

The install is very easy, follow the guide on [Fabric's website](https://fabricmc.net/). Once complete, come back here and resume.

## Downloading and Installing the Mod Tool

To help maintain mods and remove the burden of manually updating them, this tool was created to automate the installation of mods. Head over to the [Releases page](https://github.com/effisso/mc-mod-installer/releases) and download the latest release version of the tool.

Unzip the contained application executable into your preferred location on your machine. A few ideas:

* Somewhere on the computer's PATH
* Directly in your Minecraft installation folder
* Use a combination of both of the above by adding the Minecraft install location to your PATH to get the best of both worlds

## Verifying the Tool Works

Open a new command prompt/terminal instance. If your executable is NOT somewhere in the PATH, [navigate the terminal's working directory](https://www.minitool.com/news/how-to-change-directory-in-cmd.html) to the folder where the executable was unzipped.

Type `mcmods version` into the terminal and hit enter. Verify that some version string is printed (e.g. `0.1.23`). That number should match the number of the release downloaded earlier.

## Configuring the Minecraft Install Location

The tool makes a best guess at where Minecraft is installed based on the operating system. You may need to change the path used by the tool if your Minecraft installation is customized. To see what path the tool is attempting to use, enter the command `mcmods mcpath`. If the output of the command matches where Minecraft is installed on the machine, no further action is needed.

However, if the path needs to be updated, use the set flag on the command followed by the full path to the installation (use quotes around the path if it contains spaces):

```mcmods mcpath --set C:\custom\path\to\.minecraft```

You can use `mcmods mcpath` again to verify that the configuration was updated.

## Installing Mods

Now the tool is ready to install mods! Before configuring any client-only mods on your machine, this version of the install will only install all of the required and recommended server mods. Further documentation about adding client-only mods is available in the [Adding Custom Mods doc](https://github.com/effisso/mc-mod-installer/tree/main/docs/AddingCustomMods.md); and information about excluding some optional server mods from the install can be found in the [Installing Mods doc](https://github.com/effisso/mc-mod-installer/tree/main/docs/InstallingMods.md).

Simply use the command `mcmods install` and wait for it to finish.

## Connecting to YAMS

There is a pinned message in the Discord server with the IP/port information. Use that information to connect to the server.
