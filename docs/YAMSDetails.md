# YAMS Details

Hosted on: [Apex Hosting](https://apexminecrafthosting.com/)
Mod Loader: [Fabric](https://fabricmc.net/)

## Server Configuration

* Obviously it's a modded server. We use Fabric and only mods which seem to be actively maintained.
* In addition to the 1.18 Vanilla terrain generation rewrite, we use Terralith 2.0 to enhance terrain generation and add new biomes, foliage, and natural features (no new blocks).
* Spawn Chunks have been entirely disabled to discorage building all farms right at spawn.
    * A mod was added which allows crafting an item which loads specific chunks. These loaders include random ticking, which is even better (but more expensive) than the vanilla spawn chunks.
    * By default, there is a max of 256 chunks that can be perpetually loaded, so be considerate of the aggregate when making new farms. This number can be changed, but not without broad performance implications.
    * In general, try to build farms using the lowest amount of chuncks necessary which require updates. It adds an extra challenge to farm-building :)
* An auto-crafting table has been added to the game. A big use for this is combining item outputs from farms into a denser block form.
* The server has an in-browser, interactive map mod hosted on the same IP with a different port. See the pinned message in the #general channel of the CDP Discord Server for the IP and port
* Carpet Mod is used to provide many small vanilla-friendly, reversable tweaks to the server. At the time of writing, here's what we use (and the [full list of options with descriptions](https://github.com/gnembon/fabric-carpet/wiki/Current-Available-Settings)):
    * fastRedstoneDust true
    * leadFix true
    * antiCheatDisabled true
    * persistentParrots true
    * xpNoCooldown true
    * allowSpawningOfflinePlayers true
    * chainStone true
    * commandPlayer true
    * commandScript ops
    * ctrlQCraftingFix true
    * desertShrubs true
    * flippinCactus true
    * huskSpawningInTemples true
    * lagFreeSpawning true
    * lightningKillsDropsFix true
    * optimizedTNT true
    * renewableCoral true
    * renewableSponges true
    * shulkerSpawningInEndCities true
    * silverFishDropGravel true
    * spawnChunksSize 0
    * stackableShulkerBoxes true
    * And a few from [carpet-extra](https://github.com/gnembon/carpet-extra)
        * dispensersFeedAnimals true
        * renewableSand true
