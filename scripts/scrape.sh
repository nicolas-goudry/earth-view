#!/usr/bin/env nix-shell
#! nix-shell -i bash --pure
#! nix-shell -p curl cacert jq
#! nix-shell -I nixpkgs=https://github.com/NixOS/nixpkgs/archive/057f9aecfb71c4437d2b27d3323df7f93c010b7e.tar.gz
# shellcheck shell=bash

OUT="$PWD/earthview-$(date +%y%m%d%H%M%S).json"
EARTH_VIEW="https://earthview.withgoogle.com"

declare -a urls=()

check_img() {
  # Only account for images not triggering an error (404) and return the final URL from redirect
  if URL=$(curl -fsL "$EARTH_VIEW/$1" -w "%{url_effective}" -o /dev/null); then
    echo "$URL"
  fi
}

echo "Scaping Earth View images... This could take a while (~10m)."

# Process N checks in parallel
N=10

# Append found images in bg process to bash array
while read -r img; do
  urls+=("$img");
done < <(
  # Loop over known indexes to find images URL
  for index in {1000..15000}; do
    # Wait if max parallel processes reached
    ((i=i%N)); ((i++==0)) && wait
    check_img "$index" &
  done
)

# Wait for latest results
wait

# Generate JSON from found image URLs
# - source is the original image URL
# - name is the last part of the original image URL appended with a '.jpg' extension
# - url is the image download URL (/download/<id>.jpg)
printf '%s\n' "${urls[@]}" | jq -R . | jq -s ". | map({ \
  source: ., \
  name: (. | capture(\"(?!.*/)(?<name>.+)\").name + \".jpg\"), \
  url: (\"$EARTH_VIEW/download/\" + (. | capture(\"(?<id>\\\d+)$\").id) + \".jpg\") \
}) | sort_by(.url)" > "$OUT"
echo "Results wrote to $OUT"
