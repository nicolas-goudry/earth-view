{ buildGoModule, lib, ... }:

buildGoModule rec {
  pname = "earth-view";
  version = "1.0.0";
  src = ./src;
  vendorHash = "sha256-eKeUhS2puz6ALb+cQKl7+DGvm9Cl+miZAHX0imf9wdg=";

  meta = {
    description = "List and download Google Earth View images";
    homepage = "https://github.com/nicolas-goudry/earth-view";
    license = lib.licenses.mit;
    maintainers = [ lib.maintainers.nicolas-goudry ];
    mainProgram = pname;
  };
}
