{ config, lib, pkgs, ... }:

let
  inherit ((pkgs.callPackage ../../. { inherit pkgs; })) earth-view;

  cfg = config.services.earth-view;

  fehFlags = lib.concatStringsSep " "
    ([ "--bg-${cfg.display}" "--no-fehbg" ]
      ++ lib.optional (!cfg.enableXinerama) "--no-xinerama");
in
source:
pkgs.writeScriptBin "start" ''
  #!${pkgs.bash}/bin/bash

  outdir="$HOME/${cfg.imageDirectory}"

  ${pkgs.coreutils}/bin/mkdir -p $outdir
  file=$(${earth-view}/bin/earth-view fetch random -i ${source} -o $outdir)

  if test $? -ne 0; then
    ${pkgs.coreutils}/bin/echo "Error while fetching image"
    exit 1
  fi

  if test "$XDG_CURRENT_DESKTOP" = "GNOME"; then
    ${pkgs.coreutils}/bin/echo "GNOME detected, use gsettings"
    ${pkgs.glib}/bin/gsettings set org.gnome.desktop.background picture-uri file://$file
    ${pkgs.glib}/bin/gsettings set org.gnome.desktop.background picture-uri-dark file://$file
  elif test "$XDG_CURRENT_DESKTOP" = "KDE"; then
    ${pkgs.coreutils}/bin/echo "KDE detected, use plasma-apply-wallpaperimage"
    ${pkgs.kdePackages.plasma-workspace}/bin/plasma-apply-wallpaperimage $file
  else
    ${pkgs.coreutils}/bin/echo "Could not detect environment, use feh"
    ${pkgs.feh}/bin/feh ${fehFlags} $file
  fi

  ${pkgs.coreutils}/bin/ln -fs $file $outdir/.current
''
