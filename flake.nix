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
  in {
    # Provide some binary packages for selected system types.
    packages = forAllSystems (system: let
      pkgs = nixpkgsFor.${system};
    in {
      bar-unit-info = pkgs.buildGoModule {
        pname = "bar-unit-info";
        inherit version;
        # In 'nix develop', we don't need a copy of the source tree
        # in the Nix store.
        src = ./.;

        # This hash locks the dependencies of this package. It is
        # necessary because of how Go requires network access to resolve
        # VCS.  See https://www.tweag.io/blog/2021-03-04-gomod2nix/ for
        # details. Normally one can build with a fake hash and rely on native Go
        # mechanisms to tell you what the hash should be or determine what
        # it should be "out-of-band" with other tooling (eg. gomod2nix).
        # To begin with it is recommended to set this, but one must
        # remember to bump this hash when your dependencies change.
        # vendorHash = pkgs.lib.fakeHash;

        vendorHash = "sha256-5s7PvRc9S3nE2XZd21tUchSXhBkyxluqQIdlh6NWMeA=";
      };
    });

    # Add dependencies that are only needed for development
    devShells = forAllSystems (system: let
      pkgs = nixpkgsFor.${system};
      barRepo = pkgs.fetchgit {
        url = "https://github.com/beyond-all-reason/Beyond-All-Reason.git";
        sparseCheckout = [
          "units"
          "language/en"
          "luaui/configs"
        ];
        hash = "sha256-Uw8X4tXrXub4AVwRhDT9KZqrBFgkOvxeDPD6XsD7hvQ=";
      };
    in {
      default = pkgs.mkShell {
        buildInputs = with pkgs; [go just lua lua54Packages.dkjson jq barRepo];
        shellHook = ''
          export GAME_REPO=${barRepo}
          ln -s ${barRepo} bar-repo
        '';
      };
    });

    # The default package for 'nix build'. This makes sense if the
    # flake provides only one package or there is a clear "main"
    # package.
    defaultPackage = forAllSystems (system: self.packages.${system}.bar-unit-info);
  };
}
