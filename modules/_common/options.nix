{ lib, ... }:

{
  enable = lib.mkEnableOption "" // {
    description = ''
      Whether to enable Earth View service.

      Note, if you have set up a custom desktop manager session,
      then the session configuration must have the `bgSupport`
      option set to `true` or the background image set by this
      module may be overwritten.
    '';
  };

  interval = lib.mkOption {
    type = with lib.types; nullOr str;
    default = null;
    example = "1h";
    description = ''
      The duration between changing background image. Set to
      `null` to only set background when logging in. Should be
      formatted as a duration understood by systemd.
    '';
  };

  imageDirectory = lib.mkOption {
    type = lib.types.str;
    default = ".earth-view";
    example = "backgrounds";
    description = ''
      The directory to which background images should be
      downloaded, relative to HOME.
    '';
  };

  display = lib.mkOption {
    type = lib.types.enum [ "center" "fill" "max" "scale" "tile" ];
    default = "fill";
    description = ''
      Display background images according to this option.

      Note that this option has no effect on GNOME shell desktops.
    '';
  };

  enableXinerama = lib.mkOption {
    type = lib.types.bool;
    default = true;
    description = ''
      Will place a separate image per screen when enabled,
      otherwise a single image will be stretched across all
      screens.

      Note that this option has no effect on GNOME shell desktops.
    '';
  };

  autoStart = lib.mkOption {
    type = lib.types.bool;
    default = false;
    description = ''
      Whether to start the service automatically, along with its
      timer when `interval` is set.
    '';
  };
}
