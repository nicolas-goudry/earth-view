# Changelog

## v2.2.1 - 09/09/2025

🐞 **Bug Fixes**

- Update plasma-workspace package set (https://github.com/nicolas-goudry/earth-view/pull/15)

## v2.2.0 - 01/06/2024

🐞 **Bug Fixes**

- Early run (on `graphical-session-pre.target`) failing to detect desktop
- 🏠 Prevent auto start if Home Manager's `systemd.user.startServices` is enabled

🚀 **Features**

- Keep track of currently set background
- Add automatic garbage collection through `gc` attribute set option

✨ **Polish**

- Rename `common` directory to `_common`
- Rename `script.nix` to `start.nix`
- Make assertions common to both modules
- 🏠 Improve activation script
- Use systemd dependencies to start timers

## v2.1.1 - 21/05/2024

🐞 **Bug Fixes**

- 🏠 Invalid source file given to start script

## v2.1.0 - 21/05/2024

🚀 **Features**

- Full rewrite of Go module (https://github.com/nicolas-goudry/earth-view/pull/8)
- Add `autoStart` option (https://github.com/nicolas-goudry/earth-view/pull/12)

✨ **Polish**

- Avoid Nix `with` usage (https://github.com/nicolas-goudry/earth-view/pull/9)
- Deduplicate code (https://github.com/nicolas-goudry/earth-view/pull/11)

## v2.0.0 - 19/05/2024

🐞 **Bug Fixes**

- 🏠 Invalid module path (https://github.com/nicolas-goudry/earth-view/pull/7)

## v1.4.2 - 28/04/2024

🐞 **Bug Fixes**

- 📦 Downgrade Go to v1.21.9 (NixOS 23.11 compatibility)

## v1.4.1 - 27/04/2024

🐞 **Bug Fixes**

- ❄️ Create `imageDirectory` if it does not exist

## v1.4.0 - 27/04/2024

🚀 **Features**

- Add KDE support
- ❄️ Use user-managed systemd service

✨ **Polish**

- ❄️ Sync service definition with Home Manager module
- 🏠 Remove useless assertion

## v1.3.0 - 27/04/2024

🚀 **Features**

- 🏠 Add GNOME detection

🐞 **Bug Fixes**

- 🏠 Set `imageDirectory` relative to `$HOME` directory

✨ **Polish**

- ❄️ Only apply `extraGSettingsOverrides` on GNOME
- ❄️ Update config files names
- ❄️ Create current image dummy file only on GNOME

## v1.2.0 - 26/04/2024

🚀 **Features**

- Support offline mode (fail service)

## v1.1.0 - 26/04/2024

🚀 **Features**

- Avoid downloading images if they already exist

## v1.0.0 - 26/04/2024

🚀 **Features**

- Full rewrite to Flake exporting:
  - ❄️ NixOS module
  - 🏠 Home Manager module
  - Go packages
  - Development environment

✨ **Polish**

- Replace Bash scraper script by Go package
- Update source of truth file format and name

## v0.1.0 - 03/12/2023

🚀 **Features**

- Add initial module
- Add Bash scraper script
- Add source of truth up-to-date as of Dec 3, 2023
