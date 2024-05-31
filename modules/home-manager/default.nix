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
  startScript = common.mkStartScript "$HOME/${cfg.imageDirectory}/.source";
in
{
  meta.maintainers = [ lib.maintainers.nicolas-goudry ];
  options.services.earth-view = common.options;

  config = lib.mkIf cfg.enable (lib.mkMerge [
    {
      inherit (common) assertions;

      home.file."${cfg.imageDirectory}/.source".source = ../../earth-view.json;

      systemd.user.services.earth-view = {
        Unit = {
          Description = "Set random desktop background from Earth View";
          After = [ "graphical-session-pre.target" ];
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
    (lib.mkIf (cfg.interval != null) {
      systemd.user.timers.earth-view = {
        Unit.Description = "Set random desktop background from Earth View";
        Timer.OnUnitActiveSec = cfg.interval;
        Install.WantedBy = [ "timers.target" ];
      };
    })
    (lib.mkIf (cfg.autoStart || (cfg.gc.enable && cfg.gc.interval != null)) {
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
        + (lib.optionalString (cfg.autoStart && cfg.interval != null) ''
          run ${pkgs.systemd}/bin/systemctl --user start earth-view.timer
        '')
        + (lib.optionalString cfg.autoStart ''
          run ${pkgs.systemd}/bin/systemctl --user start earth-view.service
        '')
        + (lib.optionalString (cfg.gc.enable && cfg.gc.interval != null) ''
          run ${pkgs.systemd}/bin/systemctl --user start earth-view-gc.service
        '')
      );
    })
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
    (lib.mkIf (cfg.gc.enable && (cfg.gc.interval == null)) {
      systemd.user.services.earth-view-gc.Install.WantedBy = [ "earth-view.service" ];
    })
    (lib.mkIf (cfg.gc.enable && cfg.gc.interval != null) {
      systemd.user.timers.earth-view-gc = {
        Unit.Description = "Garbage collect Earth View images";
        Timer.OnUnitActiveSec = cfg.gc.interval;
        Install.WantedBy = [ "timers.target" ];
      };
    })
  ]);
}
