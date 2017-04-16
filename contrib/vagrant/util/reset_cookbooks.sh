#!/bin/sh

knife cookbook bulk delete -y "w*"
cd ../..
berks install && berks upload
