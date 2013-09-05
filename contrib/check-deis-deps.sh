#!/usr/bin/env bash

# check for git
if ! which git > /dev/null; then
  echo 'Please install git and ensure it is in your $PATH.'
  exit
fi

# check for RubyGems and friends
if ! which ruby > /dev/null; then
  echo 'Please install ruby and ensure it is in your $PATH.'
  exit
fi
if ! which gem > /dev/null; then
  echo 'Please install RubyGems and ensure "gem" is in your $PATH.'
  exit
fi
if ! which bundle > /dev/null; then
  echo 'Please install the bundler ruby gem and ensure "bundle" is in your $PATH.'
  exit
fi
bundles=`bundle list | egrep 'berkshelf|chef|foodcritic|knife-' | wc -l`
if ! [ $bundles -ge 4 ]; then
  echo 'Please run "bundle install" for required ruby gems.'
  exit
fi

# check for working knife
if ! which knife > /dev/null; then
  echo 'Please install a knife-<provider> ruby gem and ensure "knife" is in your $PATH.'
  exit
fi
if ! knife client list > /dev/null; then
  echo 'Please ensure the knife.rb file is set up correctly for your Chef account.'
  exit
fi
