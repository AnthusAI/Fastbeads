{
  description = "beads (fbd) - An issue tracker designed for AI-supervised coding workflows";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachSystem
      [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ]
      (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
          fbdBase = pkgs.callPackage ./default.nix { inherit pkgs self; };
          # Wrap the base package with shell completions baked in
          fbd = pkgs.stdenv.mkDerivation {
            pname = "beads";
            version = fbdBase.version;

            phases = [ "installPhase" ];

            installPhase = ''
              mkdir -p $out/bin
              cp ${fbdBase}/bin/fbd $out/bin/fbd

              # Generate shell completions
              mkdir -p $out/share/fish/vendor_completions.d
              mkdir -p $out/share/bash-completion/completions
              mkdir -p $out/share/zsh/site-functions

              $out/bin/fbd completion fish > $out/share/fish/vendor_completions.d/fbd.fish
              $out/bin/fbd completion bash > $out/share/bash-completion/completions/fbd
              $out/bin/fbd completion zsh > $out/share/zsh/site-functions/_fbd
            '';

            meta = fbdBase.meta;
          };
        in
        {
          packages = {
            default = fbd;

            # Keep separate completion packages for users who only want specific shells
            fish-completions = pkgs.runCommand "fbd-fish-completions" { } ''
              mkdir -p $out/share/fish/vendor_completions.d
              ln -s ${fbd}/share/fish/vendor_completions.d/fbd.fish $out/share/fish/vendor_completions.d/fbd.fish
            '';

            bash-completions = pkgs.runCommand "fbd-bash-completions" { } ''
              mkdir -p $out/share/bash-completion/completions
              ln -s ${fbd}/share/bash-completion/completions/fbd $out/share/bash-completion/completions/fbd
            '';

            zsh-completions = pkgs.runCommand "fbd-zsh-completions" { } ''
              mkdir -p $out/share/zsh/site-functions
              ln -s ${fbd}/share/zsh/site-functions/_fbd $out/share/zsh/site-functions/_fbd
            '';
          };

          apps.default = {
            type = "app";
            program = "${self.packages.${system}.default}/bin/fbd";
          };

          devShells.default = pkgs.mkShell {
            buildInputs = with pkgs; [
              go
              git
              gopls
              gotools
              golangci-lint
              sqlite
            ];

            shellHook = ''
              echo "beads development shell"
              echo "Go version: $(go version)"
            '';
          };
        }
      );
}
