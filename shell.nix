{ pkgs ? import <nixpkgs> { } }:

let
  localPkgs = import ./. { inherit pkgs; };
in
pkgs.mkShell {
  nativeBuildInputs = [
    pkgs.go
    localPkgs.earth-view
  ];
}
