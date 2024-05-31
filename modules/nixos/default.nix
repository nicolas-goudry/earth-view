# Test cases:
# -------------------------------------------------------------------------------------------
# | enable | interval | autoStart | behavior                                                |
# | ------ | -------- | --------- | ------------------------------------------------------- |
# | false  | N/A      | N/A       | Nothing                                                 |
# | true   | null     | false     | Start on login                                          |
# | true   | null     | true      | Start on login + activation                             |
# | true   | "10s"    | false     | Start on login + each 10s after login                   |
# | true   | "10s"    | true      | Start on login + activation + each 10s after activation |
# -------------------------------------------------------------------------------------------
{ config, lib, pkgs, ... }@args:

let
  cfg = config.services.earth-view;
  common = import ../_common args;
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
        unitConfig.Description = "Set random desktop background from Earth View";
        timerConfig.OnUnitActiveSec = cfg.interval;
        # Dependency to earth-view.service is required here so that the timer starts when the service starts
        # On HM module, the timer is automatically started along with the service (why? how? dunno)
        wantedBy = [ "timers.target" "earth-view.service" ];
      };
    })
    (lib.mkIf cfg.autoStart {
      system.userActivationScripts.earthViewAutoStart.text = "${pkgs.systemd}/bin/systemctl --user start earth-view.service";
    })
  ]);
}
