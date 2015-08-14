package builder

import (
	"time"

	"github.com/Masterminds/cookoo"
	"github.com/Masterminds/cookoo/fmt"
	"github.com/deis/deis/builder/confd"
	"github.com/deis/deis/builder/docker"
	"github.com/deis/deis/builder/env"
	"github.com/deis/deis/builder/etcd"
	"github.com/deis/deis/builder/git"
	"github.com/deis/deis/builder/sshd"
)

// routes builds the Cookoo registry.
//
// Esssentially this is a list of all of the things that Builder can do, broken
// down into a step-by-step list.
func routes(reg *cookoo.Registry) {

	// The "boot" route starts up the builder as a daemon process. Along the
	// way, it starts and configures multiple services, including etcd, confd,
	// and sshd.
	reg.AddRoute(cookoo.Route{
		Name: "boot",
		Help: "Boot the builder",
		Does: []cookoo.Task{

			// ENV: Make sure the environment is correct.
			cookoo.Cmd{
				Name: "vars",
				Fn:   env.Get,
				Using: []cookoo.Param{
					{Name: "HOST", DefaultValue: "127.0.0.1"},
					{Name: "ETCD_PORT", DefaultValue: "4001"},
					{Name: "ETCD_PATH", DefaultValue: "/deis/builder"},
					{Name: "ETCD_TTL", DefaultValue: "20"},
				},
			},
			cookoo.Cmd{ // This depends on others being processed first.
				Name: "vars2",
				Fn:   env.Get,
				Using: []cookoo.Param{
					{Name: "ETCD", DefaultValue: "http://$HOST:$ETCD_PORT"},
				},
			},

			// DOCKER: start up Docker and make sure it's running.
			// Then let it download the images while we keep going.
			cookoo.Cmd{
				Name: "docker",
				Fn:   docker.CreateClient,
				Using: []cookoo.Param{
					{Name: "url", DefaultValue: "unix:///var/run/docker.sock"},
				},
			},
			cookoo.Cmd{
				Name: "dockerclean",
				Fn:   docker.Cleanup,
			},
			cookoo.Cmd{
				Name: "dockerstart",
				Fn:   docker.Start,
			},
			cookoo.Cmd{
				Name: "waitfordocker",
				Fn:   docker.WaitForStart,
				Using: []cookoo.Param{
					{Name: "client", From: "cxt:docker"},
				},
			},
			/*
				cookoo.Cmd{
					Name: "slugbuilder",
					Fn:   docker.BuildImage,
					Using: []cookoo.Param{
						{Name: "tag", DefaultValue: "deis/slugbuilder"},
						{Name: "path", DefaultValue: "/usr/local/src/slugbuilder/"},
					},
				},
				cookoo.Cmd{
					Name: "slugrunner",
					Fn:   docker.BuildImage,
					Using: []cookoo.Param{
						{Name: "tag", DefaultValue: "deis/slugrunner"},
						{Name: "path", DefaultValue: "/usr/local/src/slugrunner/"},
					},
				},
			*/

			cookoo.Cmd{
				Name: "buildImages",
				Fn:   docker.ParallelBuild,
				Using: []cookoo.Param{
					{Name: "client", From: "cxt:docker"},
					{
						Name: "images",
						DefaultValue: []docker.BuildImg{
							{Path: "/usr/local/src/slugbuilder/", Tag: "deis/slugbuilder"},
							{Path: "/usr/local/src/slugrunner/", Tag: "deis/slugrunner"},
						},
					},
				},
			},

			// ETCD: Make sure Etcd is running, and do the initial population.
			cookoo.Cmd{
				Name:  "client",
				Fn:    etcd.CreateClient,
				Using: []cookoo.Param{{Name: "url", DefaultValue: "http://127.0.0.1:4001", From: "cxt:ETCD"}},
			},
			cookoo.Cmd{
				Name: "etcdup",
				Fn:   etcd.IsRunning,
				Using: []cookoo.Param{
					{Name: "client", From: "cxt:client"},
					{Name: "count", DefaultValue: 20},
				},
			},
			cookoo.Cmd{
				Name: "-",
				Fn:   Sleep,
				Using: []cookoo.Param{
					{Name: "duration", DefaultValue: 21 * time.Second},
					{Name: "message", DefaultValue: "Sleeping while etcd expires keys."},
				},
			},
			cookoo.Cmd{
				Name: "newdir",
				Fn:   fmt.Sprintf,
				Using: []cookoo.Param{
					{Name: "format", DefaultValue: "%s/users"},
					{Name: "0", From: "cxt:ETCD_PATH"},
				},
			},
			cookoo.Cmd{
				Name: "mkdir",
				Fn:   etcd.MakeDir,
				Using: []cookoo.Param{
					{Name: "path", From: "cxt:newdir"},
					{Name: "client", From: "cxt:client"},
				},
			},

			// SSHD: Create and configure host keys.
			cookoo.Cmd{
				Name: "installSshHostKeys",
				Fn:   etcd.StoreHostKeys,
				Using: []cookoo.Param{
					{Name: "client", From: "cxt:client"},
					{Name: "basepath", From: "cxt:ETCD_PATH"},
				},
			},
			cookoo.Cmd{
				Name: sshd.HostKeys,
				Fn:   sshd.ParseHostKeys,
			},
			cookoo.Cmd{
				Name: sshd.ServerConfig,
				Fn:   sshd.Configure,
			},

			// CONFD: Build out the templates, then start the Confd server.
			cookoo.Cmd{
				Name:  "once",
				Fn:    confd.RunOnce,
				Using: []cookoo.Param{{Name: "node", From: "cxt:ETCD"}},
			},
			cookoo.Cmd{
				Name:  "confd",
				Fn:    confd.Run,
				Using: []cookoo.Param{{Name: "node", From: "cxt:ETCD"}},
			},

			// Now we wait for Docker to finish downloading.
			cookoo.Cmd{
				Name: "dowloadImages",
				Fn:   docker.Wait,
				Using: []cookoo.Param{
					{Name: "wg", From: "cxt:buildImages"},
					{Name: "msg", DefaultValue: "Images downloaded"},
					{Name: "waiting", DefaultValue: "Downloading Docker images. This may take a long time. https://xkcd.com/303/"},
					{Name: "failures", From: "cxt:ParallelBuild.failN"},
				},
			},
			cookoo.Cmd{
				Name: "pushImages",
				Fn:   docker.Push,
				Using: []cookoo.Param{
					{Name: "tag", DefaultValue: "deis/slugrunner:latest"},
					{Name: "client", From: "cxt:client"},
				},
			},
			/*
				// Start SSHD.
				cookoo.Cmd{
					Name: "sshdstart",
					Fn:   sshd.Start,
				},
			*/

			// ETDCD: Now watch for events on etcd, and trigger a git check-repos for
			// each. For the most part, this runs in the background.
			cookoo.Cmd{
				Name: "Cleanup",
				Fn:   etcd.Watch,
				Using: []cookoo.Param{
					{Name: "client", From: "cxt:client"},
				},
			},
			// If there's an EXTERNAL_PORT, we publish info to etcd.
			cookoo.Cmd{
				Name: "externalport",
				Fn:   env.Get,
				Using: []cookoo.Param{
					{Name: "EXTERNAL_PORT", DefaultValue: ""},
				},
			},
			cookoo.Cmd{
				Name: "etcdupdate",
				Fn:   etcd.UpdateHostPort,
				Using: []cookoo.Param{
					{Name: "base", From: "cxt:ETCD_PATH"},
					{Name: "host", From: "cxt:HOST"},
					{Name: "port", From: "cxt:EXTERNAL_PORT"},
					{Name: "client", From: "cxt:client"},
					{Name: "sshdPid", From: "cxt:sshd"},
				},
			},

			// DAEMON: Finally, we wait around for a signal, and then cleanup.
			cookoo.Cmd{
				Name: "listen",
				Fn:   KillOnExit,
				Using: []cookoo.Param{
					{Name: "docker", From: "cxt:dockerstart"},
					{Name: "sshd", From: "cxt:sshdstart"},
				},
			},
		},
	})

	// This route is called during a user authentication for SSH.
	// The rough pattern is that we parse the local authorized keys file, and
	// then validate that the supplied user key matches an authorized key.
	//
	// This grants access to running git-receive, but does not grant access
	// to writing to the repo. That's handled by the sshReceive.
	reg.AddRoute(cookoo.Route{
		Name: "pubkeyAuth",
		Does: []cookoo.Task{
			// Parse the authorized keys file.
			// We do this every time because confd is constantly regenerating
			// the auth keys file.
			cookoo.Cmd{
				Name: "authorizedKeys",
				Fn:   sshd.ParseAuthorizedKeys,
				Using: []cookoo.Param{
					{Name: "path", DefaultValue: "/home/git/.ssh/authorized_keys"},
				},
			},

			// Auth against the keys
			cookoo.Cmd{
				Name: "authN",
				Fn:   sshd.AuthKey,
				Using: []cookoo.Param{
					{Name: "metadata", From: "cxt:metadata"},
					{Name: "key", From: "cxt:key"},
					{Name: "authorizedKeys", From: "cxt:authorizedKeys"},
				},
			},
		},
	})

	// This provides a very basic SSH ping.
	// Called by the sshd.Server
	reg.AddRoute(cookoo.Route{
		Name: "sshPing",
		Help: "Handles an ssh exec ping.",
		Does: []cookoo.Task{
			cookoo.Cmd{
				Name: "ping",
				Fn:   sshd.Ping,
				Using: []cookoo.Param{
					{Name: "request", From: "cxt:request"},
					{Name: "channel", From: "cxt:channel"},
				},
			},
		},
	})

	// This proxies a client session into a git receive.
	//
	// Called by the sshd.Server
	reg.AddRoute(cookoo.Route{
		Name: "sshGitReceive",
		Help: "Handle a git receive over an SSH connection.",
		Does: []cookoo.Task{
			// The Git receive handler needs the username. So we provide
			// it by looking up the name based on the key. When the
			// controller no longer requires username for SSH auth, we can
			// ditch this.
			cookoo.Cmd{
				Name: "fingerprint",
				Fn:   sshd.FingerprintKey,
				Using: []cookoo.Param{
					{Name: "key", From: "cxt:key"},
				},
			},
			cookoo.Cmd{
				Name: "username",
				Fn:   etcd.FindSSHUser,
				Using: []cookoo.Param{
					{Name: "client", From: "cxt:client"},
					{Name: "fingerprint", From: "cxt:fingerprint"},
				},
			},
			cookoo.Cmd{
				Name: "receive",
				Fn:   git.Receive,
				Using: []cookoo.Param{
					{Name: "request", From: "cxt:request"},
					{Name: "channel", From: "cxt:channel"},
					{Name: "operation", From: "cxt:operation"},
					{Name: "repoName", From: "cxt:repository"},
					{Name: "fingerprint", From: "cxt:fingerprint"},
					{Name: "permissions", From: "cxt:authN"},
					{Name: "user", From: "cxt:username"},
				},
			},
		},
	})
}
