#!/usr/bin/env bash

function echo_red {
  echo -e "\e[00;31m$1\e[00m"
}

# check for git
if ! which git > /dev/null; then
  echo_red 'Please install git and ensure it is in your $PATH.'
  exit 1
fi
