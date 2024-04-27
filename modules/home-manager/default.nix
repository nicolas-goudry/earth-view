{ config, lib, pkgs, ... }:

with lib;

let
  inherit ((pkgs.callPackage ../.. { inherit pkgs; })) ev-fetcher;

  cfg = config.services.earth-view;

  fehFlags = concatStringsSep " "
    ([ "--bg-${cfg.display}" "--no-fehbg" ]
      ++ lib.optional (!cfg.enableXinerama) "--no-xinerama");

  singleSwayBg = pkgs.writeScriptBin "single-swaybg" ''
    #!${pkgs.bash}/bin/bash

    pidfile="$1/.swaybg_pid"
    if ! test -e "$pidfile"; then
      ${pkgs.coreutils}/bin/echo > $pidfile
    fi
    current_pid=$(${pkgs.coreutils}/bin/cat $pidfile)

    ${pkgs.coreutils}/bin/nohup ${pkgs.swaybg}/bin/swaybg -i $2 -m fill & ${pkgs.coreutils}/bin/echo $! > $pidfile

    sleep 3
    if test -n "$current_pid"; then
      ${pkgs.coreutils}/bin/kill $current_pid
    fi
  '';

  # GNOME shell does not use X background (https://github.com/derf/feh/issues/225)
  startScript = pkgs.writeScriptBin "start-earth-view" ''
    #!${pkgs.bash}/bin/bash

    file=$(${ev-fetcher}/bin/ev-fetcher $(${pkgs.coreutils}/bin/shuf -n1 $HOME/$1/.source) $HOME/$1)

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

    if test "$XDG_SESSION_TYPE" = "wayland"; then
      ${pkgs.coreutils}/bin/echo "Wayland detected, use swaybg"
      ${singleSwayBg}/bin/single-swaybg $HOME/$1 $file
      exit 0
    fi

    ${pkgs.coreutils}/bin/echo "Could not detect environment, use feh"
    ${pkgs.feh}/bin/feh ${fehFlags} $file
  '';
in
{
  meta.maintainers = [ maintainers.nicolas-goudry ];

  options = {
    services.earth-view = {
      enable = mkEnableOption "" // {
        description = ''
          Whether to enable Earth View service.

          Note, if you are using NixOS and have set up a custom
          desktop manager session, then the session configuration must
          have the `bgSupport` option set to `true` or the background
          image set by this module may be overwritten.
        '';
      };

      interval = mkOption {
        type = types.nullOr types.str;
        default = null;
        example = "1h";
        description = ''
          The duration between changing background image. Set to
          `null` to only set background when logging in. Should be
          formatted as a duration understood by systemd.
        '';
      };

      imageDirectory = mkOption {
        type = types.str;
        default = ".earth-view";
        example = "backgrounds";
        description = ''
          The directory to which background images should be
          downloaded, relative to HOME.
        '';
      };

      display = mkOption {
        type = types.enum [ "center" "fill" "max" "scale" "tile" ];
        default = "fill";
        description = ''
          Display background images according to this option.

          Note that this option has no effect on GNOME shell desktops.
        '';
      };

      enableXinerama = mkOption {
        type = types.bool;
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

  config = mkIf cfg.enable (mkMerge [
    {
      assertions = [
        {
          assertion = pkgs.stdenv.isLinux;
          message = "services.earth-view is only compatible with Linux systems";
        }
        {
          assertion = config.home.homeDirectory != null;
          message = "'home.homeDirectory' must be defined to the user home directory";
        }
      ];

      home.file."${cfg.imageDirectory}/.source".source = ../../_earthview.txt;

      systemd.user.services.earth-view = {
        Unit = {
          Description = "Set random desktop background from Earth View";
          After = [ "graphical-session-pre.target" ];
          PartOf = [ "graphical-session.target" ];
        };

        Service = {
          Type = "oneshot";
          IOSchedulingClass = "idle";
          ExecStart = "${startScript}/bin/start-earth-view ${cfg.imageDirectory}";
        };

        Install = { WantedBy = [ "graphical-session.target" ]; };
      };
    }
    (mkIf (cfg.interval != null) {
      systemd.user.timers.earth-view = {
        Unit = { Description = "Set random desktop background from Earth View"; };
        Timer = { OnUnitActiveSec = cfg.interval; };
        Install = { WantedBy = [ "timers.target" ]; };
      };
    })
  ]);
}
