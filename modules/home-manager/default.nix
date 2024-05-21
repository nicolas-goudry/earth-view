# Test cases:
# --------------------------------------------------------------------------------------------------------------------------
# | enable | interval | autoStart | hmStartServices              | behavior                                                |
# | ------ | -------- | --------- | ---------------------------- | ------------------------------------------------------- |
# | false  | N/A      | N/A       | N/A                          | Nothing                                                 |
# | true   | null     | false     | suggest / false              | Start on login                                          |
# | true   | null     | false     | legacy  / true / sd-switch   | Start on login                                          |
# | true   | null     | true      | suggest / false              | Start on login + activation                             |
# | true   | null     | true      | legacy  / true / sd-switch   | Start on login + activation                             |
# | true   | "10s"    | false     | suggest / false              | Start on login + each 10s after manual activation       |
# | true   | "10s"    | false     | legacy  / true / sd-switch   | Start on login + activation + each 10s after activation |
# | true   | "10s"    | true      | suggest / false              | Start on login + activation + each 10s after activation |
# | true   | "10s"    | true      | legacy  / true / sd-switch   | Start on login + activation + each 10s after activation |
# --------------------------------------------------------------------------------------------------------------------------
{ config, lib, pkgs, ... }@args:

let
  cfg = config.services.earth-view;
  common = import ../common args;
  hmStartServices = config.systemd.user.startServices;
  startScript = common.mkStartScript "$HOME/${cfg.imageDirectory}";
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

        Install = {
          WantedBy = [
            "graphical-session.target"
          ] ++ (
            if (cfg.autoStart || hmStartServices == "legacy" || hmStartServices)
            then [ "default.target" ]
            else [ ]
          );
        };
      };
    }
    (lib.mkIf (cfg.interval != null) {
      systemd.user.timers.earth-view = {
        Unit = { Description = "Set random desktop background from Earth View"; };
        Timer = { OnUnitActiveSec = cfg.interval; };

        Install = {
          WantedBy = [
            "timers.target"
          ] ++ (
            if cfg.autoStart
            then [ "default.target" ]
            else [ ]
          );
        };
      };
    })
    (lib.mkIf (cfg.autoStart && (hmStartServices == "suggest" || !hmStartServices)) {
      home.activation.earth-view = lib.hm.dag.entryAfter [ "writeBoundary" ] (
        lib.concatStringsSep "\n" [
          ''
            #!${pkgs.bash}/bin/bash

            dryRun="''${DRY_RUN:-}"
          ''
          (if cfg.interval == null then "" else ''
            if test -n "$dryRun"; then
              ${pkgs.coreutils}/bin/echo "Would activate earth-view timer through systemctl"
            else
              ${pkgs.systemd}/bin/systemctl --user start earth-view.timer
            fi
          '')
          ''
            if test -n "$dryRun"; then
              ${pkgs.coreutils}/bin/echo "Would activate earth-view service through systemctl"
            else
              ${pkgs.systemd}/bin/systemctl --user start earth-view.service
            fi
          ''
        ]
      );
    })
  ]);
}
