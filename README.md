# 🌎 earth-view

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white) ![NixOS](https://img.shields.io/badge/NIXOS-5277C3.svg?style=for-the-badge&logo=NixOS&logoColor=white) ![Home Manager](https://img.shields.io/badge/home%20manager-EC733B?style=for-the-badge)

Randomly set desktop background from 2000+ images sourced from [Google Earth View](https://earthview.withgoogle.com).

Currently supporting:

- X desktops
- GNOME on Wayland or X11

> Wayland compositors support will come soon. [Track issue](https://github.com/nicolas-goudry/earth-view/issues/2).

## 📥 Installation

### ❄️ NixOS

#### Flakes (recommended)

```nix
{
  # Use unstable to get latest updates
  inputs.earth-view.url = "github:nicolas-goudry/earth-view";
  # Pin to a given revision
  #inputs.earth-view.url = "github:nicolas-goudry/earth-view/8c193eeb245cf4b5394f6441d31728775657a80a";
  # Optional, follow your nixpkgs input
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
  # Pin to a given revision
  #inputs.earth-view.url = "github:nicolas-goudry/earth-view/8c193eeb245cf4b5394f6441d31728775657a80a";
  # Optional, follow your nixpkgs input
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
  services.earth-view = {
    enable = true;
    interval = "1h";
    imageDirectory = ".earth-view"; # Home Manager only
    display = "fill";
    enableXinerama = true;
  }
}
```

> [!TIP]
> Currently, the systemd service is not automatically started. To manually start it, you can use the following commands after applying your configuration:
>
> ```shell
> # NixOS module
> sudo systemctl start earth-view.service
> sudo systemctl start earth-view.timer
>
> # Home Manager module
> systemctl --user start eath-view.service
> systemctl --user start eath-view.timer
> ```

> [!WARNING]
> When using the NixOS module with GNOME and a custom background image was already set, you have to reset it in order for the service to work:
>
> ```shell
> gsettings reset org.gnome.desktop.background picture-uri
> gsettings reset org.gnome.desktop.background picture-uri-dark
> ```
>
> You may also have to logout and login again to see the background applied.

### `enable`

Whether to enable Earth View service.

Note, if you are using NixOS and have set up a custom desktop manager session for Home Manager, then the session configuration must have the `bgSupport` option set to `true` or the background image set by this module may be overwritten.

### `interval`

The duration between changing background image. Set to `null` to only set background when logging in. Should be formatted as a [duration understood by systemd](https://www.freedesktop.org/software/systemd/man/latest/systemd.time.html#Parsing%20Time%20Spans).

### `imageDirectory`

The directory to which background images should be downloaded, relative to `$HOME`. `home.homeDirectory` must be set by your Home Manager configuration.

> [!IMPORTANT]
> This option is only available in the Home Manager module, since with the NixOS module we use systemd via the system manager and therefore cannot access the user home directory.
>
> Images will be stored in `/etc/earth-view`.

### `display`

Display background images according to this option. See [`feh` documentation](https://man.archlinux.org/man/feh.1.en#BACKGROUND_SETTING) for details.

> [!NOTE]
> This option has no effect on GNOME shell desktops.

### `enableXinerama`

Will place a separate image per screen when enabled, otherwise a single image will be stretched across all screens.

> [!NOTE]
> This option has no effect on GNOME shell desktops.

## 🧐 How it works

### Source of truth

All discovered images URLs from Earth View are saved in [`_earthview.txt`](./_earthview.txt), which is the source of truth of this module.

To create this file, we use a small [Go module](./src/scraper/main.go) which scrapes the Earth View static assets in order to find valid images URLs. If you want to use it locally:

```shell
# Use go
go run ./scraper

# Use a devshell
nix develop # ...or nix-shell
ev-scraper

# Run via nix
nix run '.#ev-scraper'

# Build it
nix build '.#ev-scraper' # ...or nix-build -A ev-scraper
./result/bin/ev-scraper
```

### Image selection

To select an image, a random line from the source of truth is read and a [Go module](./src/fetcher/main.go) is used to download it and save it to a given location. The need for a Go module comes from the fact that Earth View exposes images as JSON object with a `dataUri` key containing the base64 encoded image. It also contribute to reduce Bash usage.

The images are downloaded at different locations given the module used:

- for NixOS, the path is `/etc/earth-view`, it is not configurable
- for Home Manager, the path defined by the `imageDirectory` option is used

### systemd

Both modules use a systemd unit, along with a timer when `interval` is specified. The NixOS module uses a system-wide service, while the Home Manager module uses a user-managed service.

These services execute a Bash script which uses the Go module described in the previous section to fetch the image and then set the desktop background accordingly. Read further for more details.

### Some background

Setting the background is handled differently by NixOS and Home Manager modules and also depends on the desktop manager used. The programs used by the module are the following:

- GNOME on Wayland or X11: `gsettings`
- X: [`feh`](https://github.com/derf/feh)

Why not only use `feh`, would you ask? Well, as of today it does not support setting the GNOME background image. [And it may not ever support it](https://github.com/derf/feh/issues/225). It seems that it also does not work with KDE. And obviously it does not work with Wayland compositors.

#### Home Manager

We detect the current desktop environment with the `XDG_CURRENT_DESKTOP` environment variable and set the background with the right program. If we cannot detect the desktop we rely on `feh`.

#### NixOS

**GNOME:**

Things are a little bit hacky, but it works™️.

We must deal with `systemd` being run as `root` as well as NixOS immutability:

- we cannot use `gsettings` as it will not have any effect on the current user configuration
- we cannot interact directly with `dconf`, neither via the command line nor by messing with `/etc/dconf/db/local.d`

Instead, the module does the following:

- write a dummy file at `/etc/earth-view/current`
- define `extraGSettingsOverrides` to set the GNOME background to this file
- force link the background image to this file

**Other desktops:**

Only `feh` is used.

## 🎩 Acknowledgments

This module is heavily based on the [`random-background` service](https://github.com/nix-community/home-manager/blob/9f9e277b60a6e6915ad3a129e06861044b50fdf2/modules/services/random-background.nix) of [`home-manager`](https://github.com/nix-community/home-manager), by [rycee](https://github.com/rycee).

The idea to create this module was triggered by the [Google Earth Wallpaper Gnome extension](https://extensions.gnome.org/extension/1295/google-earth-wallpaper/), which is not anymore compatible with Gnome 45 or Gnome 46.

Last but not least, this would not be possible without the great [Google Earth View](https://earthview.withgoogle.com) website and the [Chrome extension](https://chromewebstore.google.com/detail/earth-view-from-google-ea/bhloflhklmhfpedakmangadcdofhnnoh).

## 📝 TODO

- [ ] 🏠 Find a way to detect if GNOME is being used as we cannot use config attrset [like we do in NixOS module](./modules/nixos/default.nix#L9)
- [ ] 🏗 Setup Github Actions to update the image URLs source file
- [ ] ✨ Add support for all Wayland compositors with [`swaybg`](https://github.com/swaywm/swaybg)
- [ ] ✨ Add `autoStart` option to enable and start the systemd services
- [ ] 🧹 Add `autoGC` option to enable images garbage collection
