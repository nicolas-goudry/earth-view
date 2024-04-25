{ pkgs ? import <nixpkgs> { } }:

let
  localPkgs = import ./. { inherit pkgs; };
in
pkgs.mkShell {
  nativeBuildInputs = with pkgs; [
    go
    localPkgs.ev-scraper
    localPkgs.ev-fetcher
  ];
}
