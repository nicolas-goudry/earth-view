{ pkgs, ... }:

[
  {
    assertion = pkgs.stdenv.isLinux;
    message = "services.earth-view is only compatible with Linux systems";
  }
]
