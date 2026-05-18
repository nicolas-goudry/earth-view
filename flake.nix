{
  description = "Set background wallpaper to a random image from Google Earth View.";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";

    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    {
      self,
      nixpkgs,
      treefmt-nix,
      ...
    }:
    let
      systems = [
        "x86_64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
        "aarch64-linux"
      ];
      eachSystem = f: nixpkgs.lib.genAttrs systems (system: f nixpkgs.legacyPackages.${system});
      treefmtEval = eachSystem (pkgs: treefmt-nix.lib.evalModule pkgs ./treefmt.nix);
    in
    {
      # nix flake check
      checks = eachSystem (
        pkgs:
        let
          checks = {
            formatting = treefmtEval.${pkgs.stdenv.hostPlatform.system}.config.build.check self;
          };
        in
        checks
        // {
          all = pkgs.runCommand "all-checks" { buildInputs = builtins.attrValues checks; } "touch $out";
        }
      );

      # nix fmt
      formatter = eachSystem (pkgs: treefmtEval.${pkgs.stdenv.hostPlatform.system}.config.build.wrapper);

      # Home Manager module
      homeManagerModules.earth-view = import ./modules/home-manager;

      # NixOS module
      nixosModules = {
        earth-view = import ./modules/nixos;
        default = self.nixosModules.earth-view;
      };

      # Packages
      packages = eachSystem (
        pkgs:
        import ./. {
          inherit pkgs;
        }
      );

      # Development environment with packages used by the module available in PATH
      devShells = eachSystem (pkgs: {
        default = pkgs.callPackage ./shell.nix { };
      });
    };
}
