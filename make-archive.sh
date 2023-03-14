#!/bin/bash

function die {
    echo -e "$*" > /dev/stdout
    echo "$0 VERSION" > /dev/stderr
    exit 1
}

version=$1

[ -z "$version" ] && die 'Version is not specified'

name="exais-explorer.preview.$version"
archive="${name}.tar"

test -f "$archive" && rm "$archive"
test -f "${archive}.gz" && rm "${archive}.gz"

echo
echo Creating $archive
echo

cd build
mkdir "$name"
cp nikto-explorer config_client.json  config_server.json $name
cp ../docker-compose.yml ../QUICK_START.md $name
mkdir $name/client
cp -a client/public client/scripts client/src client/webpack client/.eslintignore client/.eslintrc.js \
    client/README.md client/package-lock.json client/package.json client/postcss.config.js \
    client/size-plugin.json client/tailwind.config.js client/yarn.lock \
    $name/client

tar cvfz ../${archive}.gz $name
rm -rf $name
cd ..
