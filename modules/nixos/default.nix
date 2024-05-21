{ config, lib, pkgs, ... }:

let
  inherit ((pkgs.callPackage ../.. { inherit pkgs; })) earth-view;

  cfg = config.services.earth-view;

  fehFlags = lib.concatStringsSep " "
    ([ "--bg-${cfg.display}" "--no-fehbg" ]
      ++ lib.optional (!cfg.enableXinerama) "--no-xinerama");

  startScript = pkgs.writeScriptBin "start-earth-view" ''
    #!${pkgs.bash}/bin/bash

    mkdir -p $HOME/$1
    file=$(${earth-view}/bin/earth-view fetch random -i /etc/earth-view/.source -o $HOME/$1)

    if test $? -ne 0; then
      ${pkgs.coreutils}/bin/echo "Error while fetching image"
      exit 1
    fi

    if test "$XDG_CURRENT_DESKTOP" = "GNOME"; then
      ${pkgs.coreutils}/bin/echo "GNOME detected, use gsettings"
      ${pkgs.glib}/bin/gsettings set org.gnome.desktop.background picture-uri file://$file
      ${pkgs.glib}/bin/gsettings set org.gnome.desktop.background picture-uri-dark file://$file
      exit 0
    fi

    if test "$XDG_CURRENT_DESKTOP" = "KDE"; then
      ${pkgs.coreutils}/bin/echo "KDE detected, use plasma-apply-wallpaperimage"
      ${pkgs.libsForQt5.plasma-workspace}/bin/plasma-apply-wallpaperimage $file
      exit 0
    fi

    ${pkgs.coreutils}/bin/echo "Could not detect environment, use feh"
    ${pkgs.feh}/bin/feh ${fehFlags} $file
  '';
in
{
  meta.maintainers = [ lib.maintainers.nicolas-goudry ];

  options = {
    services.earth-view = {
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
    };
  };

  config = lib.mkIf cfg.enable (lib.mkMerge [
    {
      assertions = [
        {
          assertion = pkgs.stdenv.isLinux;
          message = "services.earth-view is only compatible with Linux systems";
        }
      ];

      environment.etc."earth-view/.source".source = ../../earth-view.json;

      systemd.user.services.earth-view = {
        unitConfig = {
          Description = "Set random desktop background from Earth View";
          After = [ "graphical-session-pre.target" ];
          PartOf = [ "graphical-session.target" ];
        };

        serviceConfig = {
          Type = "oneshot";
          IOSchedulingClass = "idle";
          ExecStart = "${startScript}/bin/start-earth-view ${cfg.imageDirectory}";
        };

        wantedBy = [ "graphical-session.target" ];
      };
    }
    (lib.mkIf (cfg.interval != null) {
      systemd.user.timers.earth-view = {
        unitConfig = { Description = "Set random desktop background from Earth View"; };
        timerConfig = { OnUnitActiveSec = cfg.interval; };
        wantedBy = [ "timers.target" ];
      };
    })
  ]);
}
