# Quick Start Guide

The following sections demonstrate how to get started with the tool if followed in order.

## Installing {{TBD MOD LOADER}} Mod Framework

Clients connecting to YAMS must use install a separate mod framework on their machine. This framework must match the server for mod compatibility. Without the proper framework or mods, the connection will not be allowed, or will be unstable.

TODO

## Downloading and Installing the Mod Tool

To help maintain mods and remove the burden of manually updating things, this tool was created to automate the installation of mods. Head over to the [Releases page](https://github.com/effisso/mc-mod-installer/releases) and download the latest version of the tool.

Unzip the contained application executable into your preferred location on your machine. A few suggestions:
* Somewhere on the computer's PATH, so the command can be run from anywhere
* Directly in your Minecraft installation folder, so it's easy to remember where it is
* Use a combination of both of the above by adding the Minecraft install location to your PATH to get the best of both worlds

## Verifying the Tool Works

Open a new command prompt/terminal instance. If your executable is NOT somewhere in the PATH, [navigate the terminal's working directory](https://www.minitool.com/news/how-to-change-directory-in-cmd.html) to the folder where the executable was unzipped.

Type `mcmods version` into the terminal and hit enter. Verify that some version string is printed (e.g. 0.1.23). That number should match the number of the release downloaded earlier.

## Configuring the Minecraft Install Location

The tool makes a best guess at where Minecraft is installed based on the operating system. You may need to change the path used by the tool if your Minecraft installation is customized. To see what path the tool is attempting to use, enter the command `mcmods mcpath`. If the output of the command matches where Minecraft is installed on the machine, no further action is needed.

However, if the path needs to be updated, use the set flag on the command followed by the full path to the installation (use quotes around the path if it contains spaces):

```mcmods mcpath --set C:\custom\path\to\.minecraft```

You can use `mcmods mcpath` again to verify that the configuration was updated.

## Installing Mods

Now the tool is ready to install mods! In order to play on the YAMS server, certain mods are required to be installed on the client This portion of the guide demonstrates a basic client install with no customizations or client-only mods. Further documentation about those advanced topics is available in the [docs folder](https://github.com/effisso/mc-mod-installer/tree/main/docs). Here, we will install all of the required mods, as well as some that offer performance boosts and bug fixes.

Simply type `mcmods install` and wait for it to finish.

That's it.

With that, the machine should have all the mods installed in the correct location for the mod framework; no additional configuration required. Fire up Minecraft and head on over to YAMS.