# Docs Test
this container helps to test the documentation before you push to github

## How-to

### Using local docker

```
DEIS_ROOT=/home/myuser/deis-git-repository

$ docker build -t myuser/deis-docs .
$ docker run --rm -t -i -v $DEIS_ROOT:/deis -p 8000:8000 local/deis-docs
```

Point your browser to http://localhost:8000

### Using remote docker (e.g. boot2docker)

When using a remote docker you must commit your changes and set GIT_REPO so the remote host can use your files.
If you are using other branch than master, set the GIT_BRANCH too.

```
$ docker build -t myuser/deis-docs .
$ docker run --rm -t -i -e "GIT_REPO=https://github.com/deis/deis" -e "GIT_BRANCH=release" -p 8000:8000 myuser/deis-docs
```

Point your browser to http://remotehost:8000

### Using a HTTP Proxy

```
docker run -e "http_proxy=http://proxy:3128" -e "https_proxy=http://proxy:3128" --rm -t -i -v $DEIS_ROOT:/deis -p 8000:8000 myuser/deis-docs
```
