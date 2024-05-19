{ buildGoModule, lib, ... }:

buildGoModule rec {
  pname = "ev-scraper";
  version = "0.0.1";
  src = ../src;

  vendorHash = null;
  subPackages = [ "./scraper" ];

  postInstall = ''
    mv $out/bin/scraper $out/bin/${pname}
  '';

  meta = {
    description = "Earth View image URLs scraper";
    homepage = "https://github.com/nicolas-goudry/earth-view";
    license = lib.licenses.mit;
    maintainers = [ lib.maintainers.nicolas-goudry ];
    mainProgram = pname;
  };
}
