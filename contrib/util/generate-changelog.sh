#!/usr/bin/env bash
# vim: et tw=0 sw=4 ts=4

usage() {
    echo "Usage: $0 <from> [to]"
    echo "The REPOROOT and TYPES variables may be used to customize"
    echo "what is scraped from the commit logs. Current settings are:"
    echo
    echo "   REPOROOT = $REPOROOT"
    echo "   TYPES = $TYPES"
    echo
}

# Print the information scraped from git
retrieve() {
    while read commit hash message
    do
        # Extract the subsystem where type(subsystem): message
        SUBSYSTEM=$(echo "$message" | cut -d'(' -f2 | cut -d')' -f1 | sed 's/\*/\(all\)/g')
        # Extract the message where type(subsystem): message
        MESSAGE=$(echo "$message" | awk -F ")" '{ print $2}' | sed 's/://' | cut -f2- -d' ')
        # Generate a link to the full legal commit on GitHub
        LINK="$REPOROOT/commit/$hash"
        # Echo all this in a way that makes the commit hash and message a link
        # to the commit in markdown
        echo ' -' "[\`$commit\`]($LINK)" "$SUBSYSTEM": "$MESSAGE"
    done < <(git --no-pager log --oneline --merges --oneline --format="%h %H %b" --grep="$1" "$FROM".."$TO")
    # Scrape the information from git
}

# Wrap feature type and show its relevant commits
subheading() {
    echo "#### $1"
    echo
    retrieve "$2"
    echo
}

main() {
    # if REPOROOT is already in the environment, don't overload it and use that
    if [ -z "$REPOROOT" ];
    then
        REPOROOT="https://github.com/deis/deis"
    fi

    # if TYPES is already in the environment, don't overload it and use that
    if [ -z "$TYPES" ];
    then
        # Based on https://github.com/deis/deis/blob/master/CONTRIBUTING.md
        # The format is in the form of $shortname;$longname $shortname;longname.
        TYPES="feat(;Features fix(;Fixes docs(;Documentation chore(;Maintenance"
    fi

    # Print usage summary if user didn't specify a beginning
    if [ -z "$1" ];
    then
        usage
        exit 1
    fi

    FROM=$1
    TO=${2:-"HEAD"}

    printf "### %s -> %s\n\n", "$FROM", "$TO"

    # Iterate over the types of messages specified
    for LEGALTYPE in $TYPES
    do
        SHORT=$(echo "$LEGALTYPE" | cut -f1 -d';')
        LONG=$(echo "$LEGALTYPE" | cut -f2 -d';')

        subheading "$LONG" "$SHORT"
    done
}

if (( SHLVL == 2 ))
then
    # If this is being run as a command
    main "$*"
    exit
else
    # Otherwise this is being sourced
    unset -f main
    unset -f usage
fi
