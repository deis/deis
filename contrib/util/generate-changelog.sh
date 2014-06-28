#!/bin/sh

usage() {
    echo "Usage: $0 <from> [to]"
}

retrieve() {
    git --no-pager log --oneline --no-merges --oneline --format=" - %h %s" --grep="$1" $FROM..$TO
}

subheading() {
    echo "#### $1\n"
    retrieve "$2"
    echo
}

FROM=$1
TO=${2:-"HEAD"}

if [ -z $1 ];
then
    usage
    exit 1
fi

echo "### $FROM -> $TO\n"

subheading "Features" "feat("
subheading "Fixes" "fix("
subheading "Documentation" "docs("
