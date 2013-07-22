# base
default.deis.dir = '/opt/deis'
default.deis.username = 'deis'
default.deis.group = 'deis'
default.deis.log_dir = '/var/log/deis'

# chroot
default.deis.chroot.cache = '/var/cache/debootstrap'
default.deis.runtime.suite = 'lucid'
default.deis.runtime.arch = 'amd64'

# docker
default.deis.go.url = 'https://go.googlecode.com/files/go1.1.linux-amd64.tar.gz'
default.deis.docker.dir = '/opt/deis/docker'
default.deis.docker.url = 'https://s3.amazonaws.com/gabrtv-docker/docker-bind-mounts.tar.gz'

# runtime
default.deis.runtime.dir = '/opt/deis/runtime'
default.deis.runtime.chroot = '/opt/deis/runtime/chroot'
default.deis.runtime.suite = 'lucid'
default.deis.runtime.arch = 'amd64'
default.deis.runtime.image = 'lucid64-runtime'
default.deis.runtime.slug_root = '/opt/deis/runtime/slugs'

# gitosis
default.deis.gitosis.dir = '/opt/deis/gitosis'

# build
default.deis.build.dir = '/opt/deis/build'
default.deis.build.repository = 'https://github.com/opdemand/buildstep'
default.deis.build.image = 'deis/buildstep'
default.deis.build.slug_dir = '/opt/deis/build/slugs'

# database
default.deis.database.name = 'deis'
default.deis.database.user = 'deis'

# registry
default.deis.registry.dir = '/opt/deis/registry'
default.deis.registry.repository = 'https://github.com/dynport/docker-private-registry'

# server/api
default.deis.controller.dir = '/opt/deis/controller'
default.deis.controller.repository = 'git@github.com:opdemand/deis-controller.git'
default.deis.controller.secret_key = 'atotallysecretkey9876543210!'
default.deis.controller.debug = 'False'
default.deis.controller.workers = 4
default.deis.controller.worker_port = 8000
default.deis.controller.http_port = 80
default.deis.controller.https_port = 443

# client
default.deis.client.dir = '/opt/deis/client'
default.deis.client.repository = 'git@github.com:opdemand/deis-cli.git'
