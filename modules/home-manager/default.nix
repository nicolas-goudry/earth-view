# Test cases:
# -------------------------------------------------------------------------------------------
# | enable | interval | autoStart | behavior                                                |
# | ------ | -------- | --------- | --------------------------------------------------------|
# | false  | N/A      | N/A       | Nothing                                                 |
# | true   | null     | false     | Start on login                                          |
# | true   | null     | true      | Start on login + activation                             |
# | true   | "10s"    | false     | Start on login + each 10s after login                   |
# | true   | "10s"    | true      | Start on login + activation + each 10s after activation |
# -------------------------------------------------------------------------------------------
{ config, lib, pkgs, ... }@args:

let
  inherit (common) gcScript;

  cfg = config.services.earth-view;
  common = import ../_common args;
  hmStartServices = config.systemd.user.startServices == true
    || config.systemd.user.startServices == "legacy"
    || config.systemd.user.startServices == "sd-switch";
  startScript = common.mkStartScript "$HOME/${cfg.imageDirectory}/.source";
in
{
  meta.maintainers = [ lib.maintainers.nicolas-goudry ];
  options.services.earth-view = common.options;

  config = lib.mkIf cfg.enable (lib.mkMerge [
    {
      inherit (common) assertions;

      # Copy source of truth to image directory
      home.file."${cfg.imageDirectory}/.source".source = ../../earth-view.json;

      # Define service
      systemd.user.services.earth-view = {
        Unit = {
          Description = "Set random desktop background from Earth View";
          After = [ "graphical-session.target" ];
          PartOf = [ "graphical-session.target" ];
        };

        Service = {
          Type = "oneshot";
          IOSchedulingClass = "idle";
          ExecStart = "${startScript}/bin/start";
        };

        Install.WantedBy = [ "graphical-session.target" ];
      };
    }
    # Define timer if interval is defined
    (lib.mkIf (cfg.interval != null) {
      systemd.user.timers.earth-view = {
        Unit.Description = "Set random desktop background from Earth View";
        Timer.OnUnitActiveSec = cfg.interval;
        Install.WantedBy = [ "timers.target" "earth-view.service" ];
      };
    })
    # Define garbage collector service if enabled
    (lib.mkIf cfg.gc.enable {
      systemd.user.services.earth-view-gc = {
        Unit = {
          Description = "Garbage collect Earth View images";
          After = [ "earth-view.service" ];
          PartOf = [ "earth-view.service" ];
        };

        Service = {
          Type = "oneshot";
          IOSchedulingClass = "idle";
          ExecStart = "${gcScript}/bin/gc";
        };
      };
    })
    # Make garbage collector run after main service if enabled without interval
    (lib.mkIf (cfg.gc.enable && (cfg.gc.interval == null)) {
      systemd.user.services.earth-view-gc.Install.WantedBy = [ "earth-view.service" ];
    })
    # Define garbage collector timer if enabled with interval
    (lib.mkIf (cfg.gc.enable && cfg.gc.interval != null) {
      systemd.user.timers.earth-view-gc = {
        Unit.Description = "Garbage collect Earth View images";
        Timer.OnUnitActiveSec = cfg.gc.interval;
        Install.WantedBy = [ "timers.target" "earth-view-gc.service" ];
      };
    })
    # Define activation script if autoStart or garbage collection with interval are enabled
    # Skipped if home manager startServices is enabled
    (lib.mkIf (hmStartServices == false && (cfg.autoStart || (cfg.gc.enable && cfg.gc.interval != null))) {
      home.activation.earth-view = lib.hm.dag.entryAfter [ "writeBoundary" ] (
        ''
          #!${pkgs.bash}/bin/bash

          run() {
            if test -n "''${DRY_RUN:-}"; then
              ${pkgs.coreutils}/bin/echo "$@"
            else
              eval "$@"
            fi
          }

        ''
        # Start main service if autoStart is enabled
        + (lib.optionalString cfg.autoStart ''
          run ${pkgs.systemd}/bin/systemctl --user start earth-view.service
        '')
        # Start garbage collector service if enabled with interval
        + (lib.optionalString (cfg.gc.enable && cfg.gc.interval != null) ''
          run ${pkgs.systemd}/bin/systemctl --user start earth-view-gc.service
        '')
      );
    })
  ]);
}
