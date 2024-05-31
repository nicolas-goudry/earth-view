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
  inherit (common) gcScript;

  cfg = config.services.earth-view;
  common = import ../_common args;
  startScript = common.mkStartScript "/etc/earth-view/.source";
in
{
  meta.maintainers = [ lib.maintainers.nicolas-goudry ];
  options.services.earth-view = common.options;

  config = lib.mkIf cfg.enable (lib.mkMerge [
    {
      inherit (common) assertions;

      # Copy source of truth to /etc
      environment.etc."earth-view/.source".source = ../../earth-view.json;

      # Define service
      systemd.user.services.earth-view = {
        unitConfig = {
          Description = "Set random desktop background from Earth View";
          After = [ "graphical-session.target" ];
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
    # Define timer if interval is defined
    (lib.mkIf (cfg.interval != null) {
      systemd.user.timers.earth-view = {
        unitConfig.Description = "Set random desktop background from Earth View";
        timerConfig.OnUnitActiveSec = cfg.interval;
        wantedBy = [ "timers.target" ];
      };
    })
    # Define garbage collector service if enabled
    (lib.mkIf cfg.gc.enable {
      systemd.user.services.earth-view-gc = {
        unitConfig = {
          Description = "Garbage collect Earth View images";
          After = [ "earth-view.service" ];
          PartOf = [ "earth-view.service" ];
        };

        serviceConfig = {
          Type = "oneshot";
          IOSchedulingClass = "idle";
          ExecStart = "${gcScript}/bin/gc";
        };
      };
    })
    # Make garbage collector run after main service if enabled without interval
    (lib.mkIf (cfg.gc.enable && (cfg.gc.interval == null)) {
      systemd.user.services.earth-view-gc.wantedBy = [ "earth-view.service" ];
    })
    # Define garbage collector timer if enabled with interval
    (lib.mkIf (cfg.gc.enable && cfg.gc.interval != null) {
      systemd.user.timers.earth-view-gc = {
        unitConfig.Description = "Garbage collect Earth View images";
        timerConfig.OnUnitActiveSec = cfg.gc.interval;
        wantedBy = [ "timers.target" ];
      };
    })
    # Define activation script if autoStart or garbage collection with interval are enabled
    (lib.mkIf (cfg.autoStart || (cfg.gc.enable && cfg.gc.interval != null)) {
      system.userActivationScripts.earthViewAutoStart.text = ''
          #!${pkgs.bash}/bin/bash
        ''
        # Start service timer if interval is set and autoStart is enabled
        + (lib.optionalString (cfg.autoStart && cfg.interval != null) ''
          ${pkgs.systemd}/bin/systemctl --user start earth-view.timer
        '')
        # Start main service if autoStart is enabled
        + (lib.optionalString cfg.autoStart ''
          ${pkgs.systemd}/bin/systemctl --user start earth-view.service
        '')
        # Start garbage collector service if enabled with interval
        + (lib.optionalString (cfg.gc.enable && cfg.gc.interval != null) ''
          ${pkgs.systemd}/bin/systemctl --user start earth-view-gc.service
        '');
    })
  ]);
}
