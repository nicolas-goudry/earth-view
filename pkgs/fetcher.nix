{ buildGoModule, lib, ... }:

buildGoModule rec {
  pname = "ev-fetcher";
  version = "0.0.1";
  src = ../src;

  vendorHash = null;
  subPackages = [ "./fetcher" ];

  postInstall = ''
    mv $out/bin/fetcher $out/bin/${pname}
  '';

  meta = {
    description = "Earth View image fetcher";
    homepage = "https://github.com/nicolas-goudry/earth-view";
    license = lib.licenses.mit;
    maintainers = [ lib.maintainers.nicolas-goudry ];
    mainProgram = pname;
  };
}
