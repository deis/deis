#!/usr/bin/env bash

THIS_DIR=$(cd $(dirname $0); pwd) # absolute path
CONTRIB_DIR=$(dirname $THIS_DIR)

# Use the first command-line argument as the user-data path
USER_DATA=${1:-$CONTRIB_DIR/coreos/user-data}
# Use the second command-line argument as DEIS_NUM_INSTANCES
NUM_INSTANCES=${2:-DEIS_NUM_INSTANCES}

function parse_yaml {
   local prefix=$2
   local s='[[:space:]]*' w='[a-zA-Z0-9_]*' fs=$(echo @|tr @ '\034')
   sed -ne "s|^\($s\):|\1|" \
        -e "s|^\($s\)\($w\)$s:$s[\"']\(.*\)[\"']$s\$|\1$fs\2$fs\3|p" \
        -e "s|^\($s\)\($w\)$s:$s\(.*\)$s\$|\1$fs\2$fs\3|p"  $1 |
   awk -F$fs '{
      indent = length($1)/2;
      vname[indent] = $2;
      for (i in vname) {if (i > indent) {delete vname[i]}}
      if (length($3) > 0) {
         vn=""; for (i=0; i<indent; i++) {vn=(vn)(vname[i])("_")}
         printf("%s%s%s=\"%s\"\n", "'$prefix'",vn, $2, $3);
      }
   }'
}

if [[ $NUM_INSTANCES -ne 1 ]] ; then
    parse_yaml $USER_DATA | grep -q "#DISCOVERY_URL"
    if [[ $? -ne 1 ]]; then
        echo "No etcd discovery URL set in $USER_DATA"
        exit 1
    fi
fi
