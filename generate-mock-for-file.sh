#!/bin/sh
_cwd=$(pwd)
_base_dir=$(dirname "$0")

if [[ -f $1 ]]; then
    echo "Generate mocks for $1"
else
    echo "Usage mock-file.sh <source-file>"
    exit 1
fi

_filepath=$1
_package_folder=$(dirname $1)
_package=""
_filename="$(basename -s .go ${_filepath})"
_destination="mocks/${_filename}_mocks.go"

cd ${_base_dir}
echo "mockgen -source=${_filepath} -destination=${_destination} -package=mocks"
mockgen -source=${_filepath} -destination=${_destination} -package=mocks
cd ${_cwd}