#!/usr/bin/env bash
set -eo pipefail

herokuish buildpack install https://github.com/ddollar/heroku-buildpack-multi.git         6e79094
herokuish buildpack install https://github.com/heroku/heroku-buildpack-ruby.git           v129
herokuish buildpack install https://github.com/heroku/heroku-buildpack-nodejs.git         v64
herokuish buildpack install https://github.com/heroku/heroku-buildpack-java.git           v32
herokuish buildpack install https://github.com/heroku/heroku-buildpack-gradle.git         24a8ebe
herokuish buildpack install https://github.com/heroku/heroku-buildpack-grails.git         1ef927d
herokuish buildpack install https://github.com/heroku/heroku-buildpack-play.git           9c137b4
herokuish buildpack install https://github.com/heroku/heroku-buildpack-python.git         v53
herokuish buildpack install https://github.com/heroku/heroku-buildpack-php.git            v50
herokuish buildpack install https://github.com/heroku/heroku-buildpack-clojure.git        v63
herokuish buildpack install https://github.com/heroku/heroku-buildpack-scala.git          v43
herokuish buildpack install https://github.com/heroku/heroku-buildpack-go.git             0e20030

