{ pkgs ? import <nixpkgs> { } }:

{
  earth-view = pkgs.callPackage ./package.nix { };
}
