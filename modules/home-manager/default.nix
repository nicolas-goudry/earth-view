{ config, lib, pkgs, ... }:

with lib;

let
  inherit ((pkgs.callPackage ../.. { inherit pkgs; })) ev-fetcher;

  cfg = config.services.earth-view;

  fehFlags = concatStringsSep " "
    ([ "--bg-${cfg.display}" "--no-fehbg" ]
      ++ lib.optional (!cfg.enableXinerama) "--no-xinerama");

  # GNOME shell does not use X background (https://github.com/derf/feh/issues/225)
  # TODO: find a way to detect if GNOME is being used as we cannot use config attrset like we do in nixos module
  startScript = pkgs.writeScriptBin "start-earth-view" ''
    #!${pkgs.bash}/bin/bash

    file=$(${ev-fetcher}/bin/ev-fetcher $(${pkgs.coreutils}/bin/shuf -n1 $HOME/$1/.source) $HOME/$1)

    if test $? -ne 0; then
      echo "Error while fetching image"
      exit 1
    fi
    ${pkgs.glib}/bin/gsettings set org.gnome.desktop.background picture-uri file://$file || true
    ${pkgs.glib}/bin/gsettings set org.gnome.desktop.background picture-uri-dark file://$file || true
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
