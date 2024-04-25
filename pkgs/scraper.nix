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

  meta = with lib; {
    description = "Earth View image URLs scraper";
    homepage = "https://github.com/nicolas-goudry/earth-view";
    license = licenses.mit;
    maintainers = with maintainers; [ nicolas-goudry ];
    mainProgram = pname;
  };
}
