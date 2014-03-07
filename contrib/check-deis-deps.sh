#!/usr/bin/env bash

function echo_red {
  echo -e "\e[00;31m$1\e[00m"
}

# check for git
if ! which git > /dev/null; then
  echo_red 'Please install git and ensure it is in your $PATH.'
  exit 1
fi

# check for RubyGems and friends
if ! which ruby > /dev/null; then
  echo_red 'Please install ruby and ensure it is in your $PATH.'
  exit 1
fi
if ! which gem > /dev/null; then
  echo_red 'Please install RubyGems and ensure "gem" is in your $PATH.'
  exit 1
fi
if ! which bundle > /dev/null; then
  echo_red 'Please install the bundler ruby gem and ensure "bundle" is in your $PATH.'
  exit 1
fi
bundles=`bundle list | egrep 'berkshelf|chef|foodcritic|knife-' | wc -l`
if ! [ $bundles -ge 4 ]; then
  echo_red 'Please run "bundle install" for required ruby gems.'
  exit 1
fi
# check for working knife
if ! which knife > /dev/null; then
  echo_red 'Please install a knife-<provider> ruby gem and ensure "knife" is in your $PATH.'
  exit 1
fi
if ! bundle exec knife client list > /dev/null; then
  echo_red 'Please ensure the knife.rb file is set up correctly for your Chef account.'
  exit 1
fi
