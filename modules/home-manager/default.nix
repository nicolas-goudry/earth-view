{ config, lib, pkgs, ... }@args:

let
  cfg = config.services.earth-view;
  opts = import ../options.nix { inherit lib; };
  startScript = (import ../script.nix args) "$HOME/${cfg.imageDirectory}";
in
{
  meta.maintainers = [ lib.maintainers.nicolas-goudry ];
  options.services.earth-view = opts;

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

        Install = { WantedBy = [ "graphical-session.target" ]; };
      };
    }
    (lib.mkIf (cfg.interval != null) {
      systemd.user.timers.earth-view = {
        Unit = { Description = "Set random desktop background from Earth View"; };
        Timer = { OnUnitActiveSec = cfg.interval; };
        Install = { WantedBy = [ "timers.target" ]; };
      };
    })
  ]);
}
