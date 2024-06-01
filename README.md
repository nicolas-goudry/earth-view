# 🌎 earth-view

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white) ![NixOS](https://img.shields.io/badge/NIXOS-5277C3.svg?style=for-the-badge&logo=NixOS&logoColor=white) ![Home Manager](https://img.shields.io/badge/home%20manager-EC733B?style=for-the-badge) [![FlakeHub](https://img.shields.io/endpoint?url=https%3A%2F%2Fflakehub.com%2Ff%2Fnicolas-goudry%2Fearth-view%2Fbadge&style=for-the-badge)](https://flakehub.com/flake/nicolas-goudry/earth-view)

Randomly set desktop background from 2600+ images sourced from [Google Earth View](https://earthview.withgoogle.com).

Currently supporting:

- X desktops
- GNOME on Wayland or X11
- KDE on Wayland or X11

> Wayland compositors support will come soon.

**Disclaimer:** this project is not affiliated with Google in any way. Images are protected by copyright and are licensed only for use as wallpapers.

## 📥 Installation

### ❄️ NixOS

#### Flakes (recommended)

```nix
{
  # Use unstable to get latest updates
  inputs.earth-view.url = "github:nicolas-goudry/earth-view";

  # Pin to a given version
  #inputs.earth-view.url = "github:nicolas-goudry/earth-view/v1.4.1";

  # Use FlakeHub
  #inputs.earth-view.url = "https://flakehub.com/f/nicolas-goudry/earth-view/*.tar.gz";

  # You may optionally follow your nixpkgs input
  #inputs.earth-view.inputs.nixpkgs.follows = "nixpkgs";

  outputs = { self, nixpkgs, earth-view }: {
    nixosConfigurations.yourhostname = nixpkgs.lib.nixosSystem {
      # Customize to your system
      system = "x86_64-linux";
      modules = [
        ./configuration.nix
        earth-view.nixosModules.earth-view
      ];
    };
  };
}
```

#### `fetchTarball`

```nix
{ lib, ... }:

{
  imports = let
    # Replace this with an actual commit or tag
    rev = "<replace>";
  in [
    "${builtins.fetchTarball {
      url = "https://github.com/nicolas-goudry/earth-view/archive/${rev}.tar.gz";
      # Replace this with an actual hash
      sha256 = lib.fakeHash;
    }}/modules/nixos"
  ];
}
```

### 🏠 Home Manager

#### Flakes: NixOS system-wide home-manager configuration

```nix
{
  # Use unstable to get latest updates
  inputs.earth-view.url = "github:nicolas-goudry/earth-view";

  # Pin to a given version
  #inputs.earth-view.url = "github:nicolas-goudry/earth-view/v1.4.1";

  # Use FlakeHub
  #inputs.earth-view.url = "https://flakehub.com/f/nicolas-goudry/earth-view/*.tar.gz";

  # You may optionally follow your nixpkgs input
  #inputs.earth-view.inputs.nixpkgs.follows = "nixpkgs";

  outputs = { self, nixpkgs, home-manager, earth-view }: {
    nixosConfigurations.yourhostname = nixpkgs.lib.nixosSystem {
      # Customize to your system
      system = "x86_64-linux";
      modules = [
        ./configuration.nix
        home-manager.nixosModules.home-manager {
          home-manager.sharedModules = [
            inputs.earth-view.homeManagerModules.earth-view
          ];
        }
      ];
    };
  };
}
```

#### Flakes: Configuration via home.nix

```nix
{ inputs, ... }:

{
  imports = [
    inputs.earth-view.homeManagerModules.earth-view
  ];
}
```

#### `fetchTarball`: Configuration via home.nix

```nix
{ lib, ... }:

{
  imports = let
    # Replace this with an actual commit or tag
    rev = "<replace>";
  in [
    "${builtins.fetchTarball {
      url = "https://github.com/nicolas-goudry/earth-view/archive/${rev}.tar.gz";
      # Replace this with an actual hash
      sha256 = lib.fakeHash;
    }}/modules/home-manager"
  ];
}
```

## 🧑‍💻 Usage

```nix
{
  # Default values
  services.earth-view = {
    enable = false;
    interval = null;
    imageDirectory = ".earth-view";
    display = "fill";
    enableXinerama = true;
    autoStart = false;
    autoUpscale = false;
    gc = {
      enable = false;
      keep = 10;
      interval = null;
      sizeThreshold = "0";
    };
  };
}
```

> [!TIP]
> If [`autoStart`](#autostart) is set to `false`, the service will only be started on the next login. To set the background after installation you have to manually start the main service:
>
> ```shell
> systemctl --user start earth-view
> ```
>
> This command can also be used to manually switch the background.
>
> Same goes for the timer, if [`autoStart`](#autostart) is set to `false`, it will only be started on the next login. To launch the timer after installation you have to manually start it:
>
> ```shell
> systemctl --user start earth-view.timer
> ```

### `enable`

Whether to enable Earth View service.

Note, if you are using NixOS and have set up a custom desktop manager session for Home Manager, then the session configuration must have the `bgSupport` option set to `true` or the background image set by this module may be overwritten.

### `interval`

The duration between changing background image. Set to `null` to only set background when logging in. Should be formatted as a [duration understood by systemd](https://www.freedesktop.org/software/systemd/man/latest/systemd.time.html#Parsing%20Time%20Spans).

### `imageDirectory`

The directory to which background images should be downloaded, relative to `$HOME`.

### `display`

Display background images according to this option. See [`feh` documentation](https://man.archlinux.org/man/feh.1.en#BACKGROUND_SETTING) for details.

> [!NOTE]
> This option has no effect neither on GNOME nor KDE.

### `enableXinerama`

Will place a separate image per screen when enabled, otherwise a single image will be stretched across all screens.

> [!NOTE]
> This option has no effect neither on GNOME nor KDE.

### `autoStart`

Whether to start the service automatically, along with its timer when [`interval`](#interval) is set.

> [!NOTE]
> This feature relies on activation scripts from NixOS (`system.userActivationScripts`) and Home Manager (`home.activation`).
>
> If you are using Home Manager, you may want to use [`systemd.user.startServices`](https://nix-community.github.io/home-manager/options.xhtml#opt-systemd.user.startServices) instead.

### `autoUpscale`

All Google Earth View images have a resolution of 1200×1800px. Because of this low-ish size, they may seem blurry or even ugly on very high resolutions as well as ultra-wide monitors.

This option aims to fix this by upscaling images to a 4× scale factor after they are downloaded. The upscaling process uses [Real-ESRGAN-ncnn-vulkan](https://github.com/xinntao/Real-ESRGAN-ncnn-vulkan) under the hood. Results are mostly very good, but some images are still a bit blurry (on a 3440×1440 ultra-wide monitor).

> [!CAUTION]
> **A GPU must be available for this option to work.** The module will try to upscale the image and resort to use the original image in case of failure.
>
> When the option is enabled, it is important to keep in mind the following things:
> * the service will use much more compute resources than when upscaling is disabled
> * the service execution time will greatly increase
> * the upscaled images will take much more disk space than their original counterparts
>   * greatest size observed was ~55MB for an original of ~900KB
>   * size of upscaled images does not seem to be strictly related to their original size (an image of ~1MB got upscaled to 40MB)
>
> If you still choose to enable the option, it is recommended to:
> * set a high enough [`interval`](#interval) to avoid disturbing other processes
> * enable the [`gc`](#gcenable) option if you care about disk space (after all, the whole upscaled collection is believed to weigh more than 100GB 🐘)

### `gc.enable`

Whether to enable automatic garbage collection.

### `gc.keep`

The number of images to keep from being garbage collected. Only the most recent images will be kept. The current background will **never** be deleted.

### `gc.interval`

The duration between garbage collection runs. Set to `null` to run garbage collection along with the main service. Should be formatted as a [duration understood by systemd](https://www.freedesktop.org/software/systemd/man/latest/systemd.time.html#Parsing%20Time%20Spans).

### `gc.sizeThreshold`

Garbage collect images only if collection size exceeds this threshold. Should be formatted as a string [understood by `du`'s size option](https://man.archlinux.org/man/du.1.en).

Examples:
* `"10M"` or `"10MiB"`: deletes images when collection exceeds 10MiB (power of 1024)
* `"1GB"`: deletes images when collection exceeds 1GB (power of 1000)

## 🧐 How it works

### Source of truth

All discovered images URLs from Earth View are saved in [`earth-view.json`](./earth-view.json), which is the source of truth of this Nix module.

To create this file, we use a [custom CLI application](./src/main.go) written in Go which provides commands to list and download images from Google Earth View. In particular, we use the `list` command which scrapes the Earth View static assets in order to find valid identifiers for images. If you want to use the CLI locally:

```shell
# Use go (go must be available in your path! Hint: nix shell 'nixpkgs#go_1_22')
cd src
go run .

# Use a devshell
nix develop # ...or nix-shell
earth-view

# Run via nix
nix run '.#earth-view'

# Build it
nix build '.#earth-view' # ...or nix-build -A earth-view
./result/bin/earth-view
```

### Image selection

To select an image, the `fetch random` command is used to select a random image identifier from the source of truth, download it and save it to the `imageDirectory` directory.

### systemd

Both modules use a systemd user-managed unit, along with a timer when [`interval`](#interval) is specified.

The service executes a Bash script which uses the Go module described in the previous section to fetch the image and then set the desktop background accordingly. Read further for more details.

### Some background

Setting the background depends on the desktop manager in use. We detect the current desktop environment with the `XDG_CURRENT_DESKTOP` environment variable and set the background with the right program:

- GNOME on Wayland or X11: `gsettings`
- KDE on Wayland or X11: `plasma-apply-wallpaperimage`
- X: [`feh`](https://github.com/derf/feh)

Why not only use `feh`, would you ask? Well, as of today it does not support setting the GNOME background image. [And it may not ever support it](https://github.com/derf/feh/issues/225). It seems that it also does not work with KDE. And obviously it does not work with Wayland compositors.

## 🎩 Acknowledgments

This module is heavily based on the [`random-background` service](https://github.com/nix-community/home-manager/blob/9f9e277b60a6e6915ad3a129e06861044b50fdf2/modules/services/random-background.nix) of [`home-manager`](https://github.com/nix-community/home-manager), by [rycee](https://github.com/rycee).

The idea to create this module was triggered by the [Google Earth Wallpaper Gnome extension](https://extensions.gnome.org/extension/1295/google-earth-wallpaper/), which is not anymore compatible with Gnome 45 or Gnome 46.

Last but not least, this would not be possible without the great [Google Earth View](https://earthview.withgoogle.com) website (if it is down, use the [Wayback Machine](https://web.archive.org/web/20240118081629/http://earthview.withgoogle.com/) to browse it) and the [Chrome extension](https://chromewebstore.google.com/detail/earth-view-from-google-ea/bhloflhklmhfpedakmangadcdofhnnoh).

## 📝 TODO

- [ ] 🏗 Setup Github Actions to update the image URLs source file
- [ ] ✨ Add support for all Wayland compositors with [`swaybg`](https://github.com/swaywm/swaybg)
- [ ] 📡 Make sure network is up to avoid image download failure
