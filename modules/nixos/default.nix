{ config, lib, pkgs, ... }:

with lib;

let
  inherit ((pkgs.callPackage ../.. { inherit pkgs; })) ev-fetcher;

  cfg = config.services.earth-view;
  isGnome = config.services.xserver.desktopManager.gnome.enable;

  fehFlags = lib.concatStringsSep " "
    ([ "--bg-${cfg.display}" "--no-fehbg" ]
      ++ lib.optional (!cfg.enableXinerama) "--no-xinerama");

  # GNOME shell does not use X background (https://github.com/derf/feh/issues/225)
  startScript = pkgs.writeScriptBin "start-earth-view" (concatStringsSep "\n" ([
    ''
      #!${pkgs.bash}/bin/bash
      file=$(${ev-fetcher}/bin/ev-fetcher $(${pkgs.coreutils}/bin/shuf -n1 /etc/earth-view/source.txt) ''${1})
      if test $? -ne 0; then
        echo "Error while fetching image"
        exit 1
      fi
    ''
  ] ++ optional isGnome ''
    ${pkgs.coreutils}/bin/ln -sf $file /etc/earth-view/current
  '' ++ optional (!isGnome) ''
    ${pkgs.feh}/bin/feh ${fehFlags} $file
  ''));
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
      ];

      environment.etc = {
        "earth-view/source.txt".source = ../../_earthview.txt;
        "earth-view/current" = {
          text = "";
          mode = "0600";
        };
      };

      services.xserver.desktopManager.gnome.extraGSettingsOverrides = ''
        [org.gnome.desktop.background]
        picture-uri='file:///etc/earth-view/current'
        picture-uri-dark='file:///etc/earth-view/current'
      '';

      systemd.services.earth-view = {
        unitConfig = {
          Description = "Set random desktop background from Earth View";
          After = [ "graphical-session-pre.target" ];
          PartOf = [ "graphical-session.target" ];
        };

        serviceConfig = {
          Type = "oneshot";
          IOSchedulingClass = "idle";
          ExecStart = "${startScript}/bin/start-earth-view %E/earth-view";
        };

        wantedBy = [ "graphical-session.target" ];
      };
    }
    (mkIf (cfg.interval != null) {
      systemd.timers.earth-view = {
        unitConfig = { Description = "Set random desktop background from Earth View"; };
        timerConfig = { OnUnitActiveSec = cfg.interval; };
        wantedBy = [ "timers.target" ];
      };
    })
  ]);
}
