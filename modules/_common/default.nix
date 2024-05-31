args:

{
  options = import ./options.nix args;
  mkStartScript = import ./script.nix args;
}
