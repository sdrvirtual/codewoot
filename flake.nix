{
  description = "A Nix-flake-based Go development environment";

  inputs.nixpkgs.url = "https://flakehub.com/f/NixOS/nixpkgs/0.1"; # unstable Nixpkgs

  outputs =
    { self, ... }@inputs:

    let
      goVersion = 24; # Change this to update the whole stack

      supportedSystems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
      forEachSupportedSystem =
        f:
        inputs.nixpkgs.lib.genAttrs supportedSystems (
          system:
          f {
            pkgs = import inputs.nixpkgs {
              inherit system;
              overlays = [ inputs.self.overlays.default ];
            };
          }
        );
    in
      {
        overlays.default = final: prev: {
          go = final."go_1_${toString goVersion}";
        };

        devShells = forEachSupportedSystem (
          { pkgs }:
          {
            default = pkgs.mkShellNoCC {
              ANTHROPIC_API_KEY = "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";
              OPENAI_API_KEY = "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";
              packages = with pkgs; [
                # go (version is specified by overlay)
                go

                # goimports, godoc, etc.
                gotools

                # https://github.com/golangci/golangci-lint
                golangci-lint
                gopls
		sqlc
		goose
                aider-chat-with-playwright
		            delve
              ];
            };
          }
        );
      };
}
