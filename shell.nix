
{ pkgs ? import <nixpkgs> {} }:
  pkgs.mkShell {
    nativeBuildInputs = with pkgs.buildPackages; [
        go
        air
        sqlite
        templ
        sqlitebrowser
    ];
}
