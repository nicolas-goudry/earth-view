#!/usr/bin/env nix-shell
#! nix-shell -i bash --pure
#! nix-shell -p curl
#! nix-shell -I nixpkgs=https://github.com/NixOS/nixpkgs/archive/057f9aecfb71c4437d2b27d3323df7f93c010b7e.tar.gz
# shellcheck shell=bash

OUT="$PWD/earthview$(date +%y%m%d%H%M%S).json"
echo "[" > "$OUT"

for i in {1003..5000}; do
  echo -n "Trying $i... "
  if URL=$(curl -fsL "https://earthview.withgoogle.com/$i" -w "%{url_effective}" -o /dev/null); then
    echo "$(test "$i" -gt 1003 && echo "," || echo ""){\"url\":\"$URL\",\"download\":\"https://earthview.withgoogle.com/download/$i.jpg\"}" >> "$OUT"
    echo "yep!"
  else
    echo "nope."
  fi
done

echo "]" >> "$OUT"
