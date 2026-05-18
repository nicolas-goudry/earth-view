{ buildGoModule, lib, ... }:

buildGoModule (finalAttrs: {
  pname = "earth-view";
  version = "1.0.0";
  src = ./src;
  vendorHash = "sha256-kiYMJXsrRJxU2P6mxFtt0kZd5qu1Qbd3uIXjXFUyjZA=";
  doCheck = false;

  meta = {
    mainProgram = finalAttrs.pname;
    description = "List and download Google Earth View images";
    homepage = "https://github.com/nicolas-goudry/earth-view";
    license = lib.licenses.mit;
    maintainers = [ lib.maintainers.nicolas-goudry ];
  };
})
