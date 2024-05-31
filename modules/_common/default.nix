args:

{
  mkStartScript = import ./start.nix args;
  options = import ./options.nix args;
}
