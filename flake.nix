{
  description = "Personal Website for Conner Ohnesorge";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    systems.url = "github:nix-systems/default";
    flake-utils = {
      url = "github:numtide/flake-utils";
      inputs.systems.follows = "systems";
    };
  };

  outputs = inputs @ {flake-utils, ...}:
    flake-utils.lib.eachDefaultSystem (system: let
      overlay = final: prev: {final.go = prev.go_1_24;};
      pkgs = import inputs.nixpkgs {
        inherit system;
        overlays = [
          overlay
        ];
        config.allowUnfree = true;
      };
      buildWithSpecificGo = pkg: pkg.override {buildGoModule = pkgs.buildGo124Module;};
    in {
      devShell = let
        scripts = {
          dx = {
            exec = ''$EDITOR $REPO_ROOT/flake.nix'';
            description = "Edit flake.nix";
          };
          gx = {
            exec = "$EDITOR $REPO_ROOT/go.mod";
            description = "Edit go.mod";
          };
          clean = {
            exec = ''${pkgs.git}/bin/git clean -fdx'';
            description = "Clean Project";
          };
          reset-db = {
            exec = ''
              rm ./master.db
              rm ./master.db-shm
              rm ./master.db-wal
            '';
            description = "Reset the database";
          };
          tests = {
            exec = ''${pkgs.go}/bin/go test -v ./...'';
            description = "Run all go tests";
          };
          lint = {
            exec = ''
              export REPO_ROOT=$(git rev-parse --show-toplevel)

              ${pkgs.golangci-lint}/bin/golangci-lint run
              ${pkgs.statix}/bin/statix check $REPO_ROOT/flake.nix
              ${pkgs.deadnix}/bin/deadnix $REPO_ROOT/flake.nix
            '';
            description = "Run Linting Steps.";
          };
          format = {
            exec = ''
              cd $(git rev-parse --show-toplevel)
              ${pkgs.go}/bin/go fmt ./...
              ${pkgs.git}/bin/git ls-files \
                --others \
                --exclude-standard \
                --cached \
                -- '*.js' '*.ts' '*.css' '*.md' '*.json' \
                | xargs prettier --write
              ${pkgs.golines}/bin/golines \
                -l \
                -w \
                --max-len=80 \
                --shorten-comments \
                --ignored-dirs=.direnv .
              cd -
            '';
            description = "Format code files";
          };
          generate-js = {
            exec = ''
              ${pkgs.bun}/bin/bun build \
                  $REPO_ROOT/index.js \
                  --minify \
                  --minify-syntax \
                  --minify-whitespace  \
                  --minify-identifiers \
                  --outdir $REPO_ROOT/cmd/conneroh/_static/dist/ &
            '';
            description = "Generate JS files";
          };
          run = {
            exec = ''
              cd $REPO_ROOT && air
            '';
            description = "Run the application with air for hot reloading";
          };
        };
        scriptPackages =
          pkgs.lib.mapAttrs
          (name: script: pkgs.writeShellScriptBin name script.exec)
          scripts;
      in
        pkgs.mkShell {
          shellHook = ''
            export REPO_ROOT=$(git rev-parse --show-toplevel)
            export CGO_CFLAGS="-O2"

            echo "Available commands:"
            ${pkgs.lib.concatStringsSep "\n" (
              pkgs.lib.mapAttrsToList (
                name: script: ''echo "  ${name} - ${script.description}"''
              )
              scripts
            )}

            echo "Git Status:"
            ${pkgs.git}/bin/git status
          '';
          packages = with pkgs;
            [
              alejandra # Nix
              nixd
              statix
              deadnix

              go_1_24 # Go Tools
              air
              templ
              golangci-lint
              (buildWithSpecificGo revive)
              (buildWithSpecificGo gopls)
              (buildWithSpecificGo templ)
              (buildWithSpecificGo golines)
              (buildWithSpecificGo golangci-lint-langserver)
              (buildWithSpecificGo gomarkdoc)
              (buildWithSpecificGo gotests)
              (buildWithSpecificGo gotools)
              (buildWithSpecificGo reftools)
              pprof
              graphviz
              goreleaser

              tailwindcss # Web
              tailwindcss-language-server
              bun
              nodePackages.typescript-language-server
              nodePackages.prettier
              svgcleaner
              sqlite-web

              openssl.dev
            ]
            ++ (with pkgs;
              lib.optionals stdenv.isDarwin [
                libiconv
              ])
            ++ (with pkgs;
              lib.optionals stdenv.isLinux [
                chromium # Chromium browser
                xorg.libXcomposite # X11 Composite extension - needed by browsers
                xorg.libXdamage # X11 Damage extension - needed by browsers
                xorg.libXfixes # X11 Fixes extension - needed by browsers
                xorg.libXrandr # X11 RandR extension - needed by browsers
                xorg.libX11 # X11 client-side library
                xorg.libxcb # X11 C Bindings library
                mesa # OpenGL implementation
                alsa-lib # Audio library
                nss # Network Security Services
                nspr # NetScape Portable Runtime
                pango # Text layout and rendering
              ])
            ++ builtins.attrValues scriptPackages;
        };

      packages = {};
    });
}
