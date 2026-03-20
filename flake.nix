{
  description = "serve - serve a directory over HTTP with a QR code";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  };

  outputs =
    { self, nixpkgs }:
    let
      pkgs = nixpkgs.legacyPackages.x86_64-linux;
    in
    {
      packages.x86_64-linux.default = pkgs.buildGoModule {
        pname = "serve";
        version = self.shortRev or "dirty";
        vendorHash = "sha256-MHil/ykrrhjpfmx5m611bfv0nqpExpr0TiWiVZ8laag=";
        src = ./.;
      };

      overlays.default = final: prev: {
        serve = self.packages.${final.system}.default;
      };
    };
}
