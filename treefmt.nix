_:

{
  projectRootFile = "flake.nix";

  programs = {
    nixfmt.enable = true;
    goimports.enable = true;
    golines.enable = true;

    prettier = {
      enable = true;

      excludes = [ "earth-view.json" ];
    };

    shellcheck = {
      enable = true;

      excludes = [ ".envrc" ];
    };
  };
}
