{
  description = "golang devshell";
  nixConfig.bash-prompt = "[nix(golang)] ";
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";

  };

  outputs =
    { nixpkgs, ... }:
    {
      devShell.aarch64-darwin =
        let
          pkgs = import nixpkgs {
            system = "aarch64-darwin";
            config.allowUnfree = true;
          };
        in
        pkgs.mkShell {
          buildInputs = with pkgs; [
            go_1_24
            gopls
            gotools
            gofumpt
          ];

          CONSUL_HTTP_HOST = "common-consul-boe.bytedance.net";
          CONSUL_HTTP_PORT = "2280";
          ENV = "DEV";

          shellHook = ''
            for p in $NIX_PROFILES; do
              GOPATH="$p/share/go:$GOPATH"
            done
          '';
        };
    };
}
