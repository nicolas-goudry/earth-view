{ config, lib, pkgs, ... }@args:

let
  cfg = config.services.earth-view;
  common = import ../common args;
  startScript = common.mkStartScript "/etc/earth-view/.source";
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
          ExecStart = "${startScript}/bin/start";
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
