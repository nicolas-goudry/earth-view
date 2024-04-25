{ pkgs ? import <nixpkgs> { } }:

{
  ev-fetcher = pkgs.callPackage ./pkgs/fetcher.nix { };
  ev-scraper = pkgs.callPackage ./pkgs/scraper.nix { };
}
