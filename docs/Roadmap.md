# Roadmap

This document captures intended/potential additions to the tool. Nothing is guaranteed, just some ideas.

* Update specific fields for mod definitions
    * Needed for at least latest download URL
* Uninstall/Remove for mods/definitions
    * Specifically client-only mods
* Dry run (at least for intstall)
* Secure password prompt instead of command line arg for FTP password

Potential Tech Debt Problems:

* general architecture of cmd files
    * bad design; needs dependency injection, per-call cmd instantiation
    * unit tests all call through RootCmd
        * not *quite* unit tests, but still valid tests
        * tests the topmost user-facing layer of the application for each cmd, including cmd layering and persistent flags/funcs
        * incurs additional mocking overhead, but this is handled by an common func (plus other per-test necessities)
        * might not need to change
