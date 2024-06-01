{ config, pkgs, ... }:

let
  imgDir = config.services.earth-view.imageDirectory;
  cfg = config.services.earth-view.gc;
in
pkgs.writeScriptBin "gc" ''
  #!${pkgs.bash}/bin/bash

  outdir="$HOME/${imgDir}"

  if ! test -d $outdir; then
    ${pkgs.coreutils}/bin/echo "Image directory does not exist"
    exit 1
  fi

  if test $(${pkgs.findutils}/bin/find $outdir -type f | ${pkgs.coreutils}/bin/wc -l) -le ${toString cfg.keep}; then
    ${pkgs.coreutils}/bin/echo "Not enough candidates, skipping garbage collection"
    exit 0
  fi

  if test $(${pkgs.coreutils}/bin/du -t ${toString cfg.sizeThreshold} $outdir | ${pkgs.coreutils}/bin/wc -l) -eq 0; then
    ${pkgs.coreutils}/bin/echo "Collection has not reached size threshold, skipping garbage collection"
    exit 0
  fi

  ${pkgs.findutils}/bin/find $outdir -type f -printf '%Ts\t%h/%P\n' | \
    ${pkgs.coreutils}/bin/sort -n | \
    ${pkgs.coreutils}/bin/cut -f2 | \
    ${pkgs.gnugrep}/bin/grep -v $(${pkgs.coreutils}/bin/readlink $outdir/.current) | \
    ${pkgs.coreutils}/bin/head -n -${toString (cfg.keep - 1)} | \
    ${pkgs.findutils}/bin/xargs rm -f
  ${pkgs.findutils}/bin/find $outdir -xtype l -delete
''
