{
  description = "Beyond all reason unit info";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-24.05";
  };

  outputs = {
    self,
    nixpkgs,
  }: let
    # to work with older version of flakes
    lastModifiedDate = self.lastModifiedDate or self.lastModified or "19700101";

    # Generate a user-friendly version number.
    version = builtins.substring 0 8 lastModifiedDate;

    # System types to support.
    supportedSystems = ["x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin"];

    # Helper function to generate an attrset '{ x86_64-linux = f "x86_64-linux"; ... }'.
    forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

    # Nixpkgs instantiated for supported system types.
    nixpkgsFor = forAllSystems (system: import nixpkgs {inherit system;});

    barRepoWithPkgs = pkgs:
      pkgs.fetchgit {
        url = "https://github.com/beyond-all-reason/Beyond-All-Reason.git";
        sparseCheckout = [
          "units"
          "language/en"
          "luaui/configs"
        ];
        hash = "sha256-XXlFn5Cz1VRn3Arx4VIS4exEigxILqbrl1s8ZcgqPKc=";
      };
  in {
    # Provide some binary packages for selected system types.
    packages = forAllSystems (system: let
      pkgs = nixpkgsFor.${system};
      barRepo = barRepoWithPkgs pkgs;
    in {
      bar-unit-info = pkgs.buildGoModule {
        pname = "bar-unit-info";
        inherit version;
        src = ./.;
        # vendorHash = pkgs.lib.fakeHash;
        vendorHash = "sha256-RmiEm0l/SKZmgAc5u6QeaNqg/UgcyeRSntjTPyj2ZMA=";

        # nativeBuildInputs = [barRepo];
        postConfigure = ''
          GAME_REPO=${barRepo} go generate ./...
        '';
      };
    });

    devShells = forAllSystems (system: let
      pkgs = nixpkgsFor.${system};
      barRepo = barRepoWithPkgs pkgs;
    in {
      default = pkgs.mkShell {
        buildInputs = with pkgs; [go just lua lua54Packages.dkjson jq barRepo];
        shellHook = ''
          export GAME_REPO=${barRepo}
          ln -s ${barRepo} bar-repo
        '';
      };
    });

    defaultPackage = forAllSystems (system: self.packages.${system}.bar-unit-info);
  };
}
