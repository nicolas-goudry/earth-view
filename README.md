# earth-view

This repository holds a simple Nix module that registers a systemd service whose goal is to randomly pick an image from the [Google Earth View website](https://earthview.withgoogle.com) and set it as the Gnome desktop background image.

## Options

The module is exposed under `services.earthView` and accepts the following options:

| Name | Default value | Description |
| ---- | ------------- | ----------- |
| enable | None | Enable the service. |
| imageDirectory | `$HOME/.earth-view` | Directory where the images will be stored.<br/>_Will be created if it does not exist already._ |
| interval | `null` | Duration between background image updates.<br/>_If set to `null`, background will be updated upon login only._<br/>_The value should be formatted as a duration understood by systemd._ |

## Usage

```nix
# Import the module into your configuration
imports = [
  (builtins.fetchTarball {
    url = "https://github.com/nicolas-goudry/earth-view/archive/master.tar.gz";
    # You may have to build your system first and see it failing before updating this to the correct SHA
    sha256 = "sha256:1xs8hmr8g4fqblih0pk1sqccp1nfcwmmbbqy4a0vvjwkvl8rmczr";
  })
];

# Enable and configure the module
services.earthView = {
  enable = true;
  interval = "4h";
};
```

## Inner workings

### Source of truth

All discovered images from Earth View are saved in the [`earthview.json` file](./earthview.json), which is the source of truth of this module.

This file can be updated by running the [`scrape.sh` script](./scripts/scrape.sh), which uses `nix-shell` to create a reproducible interpreted script.\
The script is using `curl` and `jq` to scrape the Earth View website for known images index (_from 1000 to 15000_). It uses parallelization to speed up the process, but it can be quite long to run depending on the hardware capabilities of the host.\
_For the record, it takes about 11 minutes to run on a Dell XPS 15 7590 with 64GB of memory, an Intel i9 processor and a not-so-bad fiber connection to the internet._

> [!NOTE]
> The JSON file is up-to-date as of **December 3, 2023** and contains **2609** references to images.

### systemd

As mentioned earlier, the module will register a new systemd service (along with a timer service when `interval` is specified). This service will execute a bash script to download the latest version of the source of truth JSON, pick a random image and set it as the Gnome desktop background image (_for both light and dark color schemes_).

The content of the script can be viewed [here](./default.nix#L50-L76).

## Request for help

If you like this module and have knowledge of Nix, it would be greatly appreciated if you could help improve it !

The module is obviously working, but since I’m new to Nix I believe it could be improved with the following features:

* expose as a flake
* support other desktop managers
* less extensive usage of Bash
* support other image sources ?

## Acknowledgments

This module is heavily inspired by the [`random-background` service](https://github.com/nix-community/home-manager/blob/9f9e277b60a6e6915ad3a129e06861044b50fdf2/modules/services/random-background.nix) of [`home-manager`](https://github.com/nix-community/home-manager).

The idea to create this module was triggered by the [Google Earth Wallpaper Gnome extension](https://extensions.gnome.org/extension/1295/google-earth-wallpaper/), which is not anymore compatible with Gnome 45.
