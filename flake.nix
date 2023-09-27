{
  description = "Cloud API Adaptor for Confidential Containers";

  inputs = {
    nixpkgsWorking = {
      url = "github:katexochen/nixpkgs/working";
    };
    nixpkgsUnstable = {
      url = "github:nixos/nixpkgs/nixos-unstable";
    };
    flake-utils = {
      url = "github:numtide/flake-utils";
    };
  };

  outputs =
    { self
    , nixpkgsWorking
    , nixpkgsUnstable
    , flake-utils
    }:
    flake-utils.lib.eachDefaultSystem
      (system:
      let
        pkgsWorking = import nixpkgsWorking { inherit system; };
        pkgsUnstable = import nixpkgsUnstable { inherit system; };

        mkosiDev = pkgsWorking.mkosi;
        mkosiDevFull = pkgsWorking.mkosi-full;
      in
      {
        devShells = {
          # Shell for building a podvm image with mkosi.
          podvm-mkosi = pkgsUnstable.mkShell {
            nativeBuildInputs = with pkgsUnstable; [
              dnf5
              rpm
              btrfs-progs
              squashfsTools
              dosfstools
              mtools
              cryptsetup
              util-linux
            ] ++ [ mkosiDevFull ];
          };
        };

        formatter = nixpkgsUnstable.legacyPackages.${system}.nixpkgs-fmt;
      });
}
