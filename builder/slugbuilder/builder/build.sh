#!/bin/bash
set -eo pipefail

slug_file=/tmp/slug.tgz
app_dir=/app
build_root=/tmp/build
cache_root=/tmp/cache
buildpack_root=/tmp/buildpacks

mkdir -p $app_dir
mkdir -p $cache_root
mkdir -p $buildpack_root
mkdir -p $build_root/.profile.d

function output_redirect() {
	if [[ "$slug_file" == "-" ]]; then
		cat - 1>&2
	else
		cat -
	fi
}

function echo_title() {
  echo $'\e[1G----->' $* | output_redirect
}

function echo_normal() {
  echo $'\e[1G      ' $* | output_redirect
}

function ensure_indent() {
  while read line; do
    if [[ "$line" == --* ]]; then
      echo $'\e[1G'$line | output_redirect
    else
      echo $'\e[1G      ' "$line" | output_redirect
    fi

  done
}

## Copy application code over
if [ -d "/tmp/app" ]; then
	cp -rf /tmp/app/. $app_dir
else
	cat | tar -xmC $app_dir
fi

# In heroku, there are two separate directories, and some
# buildpacks expect that.
cp -r $app_dir/. $build_root

# Grant the slugbuilder user access to all relevant
# directories, then run the rest as slugbuilder
chown -R slugbuilder:slugbuilder	$app_dir \
					$cache_root \
					$buildpack_root \
					$build_root
su slugbuilder

## Buildpack fixes

export APP_DIR="$app_dir"
export HOME="$app_dir"
export REQUEST_ID=$(openssl rand -base64 32)

## Buildpack detection

buildpacks=($buildpack_root/*)
selected_buildpack=

if [[ -n "$BUILDPACK_URL" ]]; then
	echo_title "Fetching custom buildpack"

	buildpack="$buildpack_root/custom"
	rm -fr "$buildpack"
	git clone --quiet --depth=1 "$BUILDPACK_URL" "$buildpack"
	selected_buildpack="$buildpack"
	buildpack_name=$($buildpack/bin/detect "$build_root") && selected_buildpack=$buildpack
else
    for buildpack in "${buildpacks[@]}"; do
    	buildpack_name=$($buildpack/bin/detect "$build_root") && selected_buildpack=$buildpack && break
    done
fi

if [[ -n "$selected_buildpack" ]]; then
	echo_title "$buildpack_name app detected"
	else
	echo_title "Unable to select a buildpack"
	exit 1
fi

## Buildpack compile

$selected_buildpack/bin/compile "$build_root" "$cache_root" | ensure_indent

$selected_buildpack/bin/release "$build_root" "$cache_root" > $build_root/.release

## Display process types

echo_title "Discovering process types"
if [[ -f "$build_root/Procfile" ]]; then
	types=$(ruby -e "require 'yaml';puts YAML.load_file('$build_root/Procfile').keys().join(', ')")
	echo_normal "Procfile declares types -> $types"
fi
default_types=""
if [[ -s "$build_root/.release" ]]; then
	default_types=$(ruby -e "require 'yaml';puts (YAML.load_file('$build_root/.release')['default_process_types'] || {}).keys().join(', ')")
	[[ $default_types ]] && echo_normal "Default process types for $buildpack_name -> $default_types"
fi


## Produce slug

if [[ -f "$build_root/.slugignore" ]]; then
	tar --exclude='.git' --use-compress-program=pigz -X "$build_root/.slugignore" -C $build_root -cf $slug_file . | cat
else
	tar --exclude='.git' --use-compress-program=pigz -C $build_root -cf $slug_file . | cat
fi

if [[ "$slug_file" != "-" ]]; then
	slug_size=$(du -Sh "$slug_file" | cut -f1)
	echo_title "Compiled slug size is $slug_size"
fi
