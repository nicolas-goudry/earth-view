{ config, lib, pkgs, ... }:

with lib;

let
  cfg = config.services.earthView;
in {
  options = {
    services.earthView = {
      enable = mkEnableOption "" // {
        description = "Whether to enable Earth View.";
      };

      imageDirectory = mkOption {
        type = types.str;
        default = "$HOME/.earth-view";
        description = "The directory where images will be downloaded.";
      };

      interval = mkOption {
        default = null;
        type = types.nullOr types.str;
        example = "1h";
        description = ''
          The duration between changing background image, set to null to only set
          background when logging in. Should be formatted as a duration understood by
          systemd.
        '';
      };
    };
  };

  config = mkIf cfg.enable (mkMerge ([
    {
      assertions = [
        (hm.assertions.assertPlatform "services.earthView" pkgs platforms.linux)
      ];

      systemd.user.services.earth-view = {
        Unit = {
          Description = "Set random desktop background from Earth View";
          After = [ "graphical-session-pre.target" ];
          PartOf = [ "graphical-session.target" ];
        };

        Service = {
          Type = "oneshot";
          IOSchedulingClass = "idle";
          ExecStart = pkgs.writers.writeBash "earth-view" ''
            SOURCE=$(${pkgs.curl}/bin/curl -fsS https://raw.githubusercontent.com/nicolas-goudry/earth-view/master/earthview.json)

            if test -n "$SOURCE"; then
              COUNT=$(echo "$SOURCE" | ${pkgs.jq}/bin/jq -r '. | length')

              if test "$COUNT" -gt 0; then
                RND=$(("$RANDOM" % "$COUNT"))
                IMG=$(echo "$SOURCE" | ${pkgs.jq}/bin/jq -r '.['"$RND"'].download')
                DEST_ROOT="${cfg.imageDirectory}"

                if ! test "$IMG" == "null"; then
                  DEST_NAME="$DEST_ROOT/$(basename $IMG)"

                  mkdir -p "$DEST_ROOT"

                  if ! test -f "$DEST_NAME"; then
                    ${pkgs.curl}/bin/curl -fsSL "$IMG" -o "$DEST_NAME"
                  fi

                  if test -f "$DEST_NAME"; then
                    gsettings set org.gnome.desktop.background picture-uri-dark file://$DEST_NAME
                    gsettings set org.gnome.desktop.background picture-uri file://$DEST_NAME
                  fi
                fi
              fi
            fi
          '';
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
  ]));
}
