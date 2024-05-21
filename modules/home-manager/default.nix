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
  cfg = config.services.earth-view;
  common = import ../common args;
  startScript = common.mkStartScript "$HOME/${cfg.imageDirectory}/.source";
in
{
  meta.maintainers = [ lib.maintainers.nicolas-goudry ];
  options.services.earth-view = common.options;

  config = lib.mkIf cfg.enable (lib.mkMerge [
    {
      assertions = [
        {
          assertion = pkgs.stdenv.isLinux;
          message = "services.earth-view is only compatible with Linux systems";
        }
      ];

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
    (lib.mkIf cfg.autoStart {
      home.activation.earth-view = lib.hm.dag.entryAfter [ "writeBoundary" ] ''
        #!${pkgs.bash}/bin/bash

        run() {
          if test -n "''${DRY_RUN:-}"; then
            ${pkgs.coreutils}/bin/echo "$@"
          else
            eval "$@"
          fi
        }

        if test -n "${toString cfg.interval}"; then
          run ${pkgs.systemd}/bin/systemctl --user start earth-view.timer
        fi

        run ${pkgs.systemd}/bin/systemctl --user start earth-view.service
      '';
    })
  ]);
}
