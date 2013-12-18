#!/opt/deis/controller/venv/bin/python
from yaml.error import YAMLError
import argparse
import getpass
import json
import os
import subprocess
import sys
import uuid
import yaml

SLUG_DIR = os.environ['SLUG_DIR']
CONTROLLER_DIR = os.environ['CONTROLLER_DIR']

def parse_args():
    desc = """
Process a git push by running it through the buildpack process

Note this script must be run as the `git` user.
"""
    parser = argparse.ArgumentParser(description=desc,
                                     formatter_class=argparse.RawDescriptionHelpFormatter)
    parser.add_argument('src', action='store',
                        help='path to source repository')
    parser.add_argument('user', action='store',
                        help='name of user who started build',
                        default='system')
    # check execution environment
    if getpass.getuser() != 'git':
        sys.stderr.write('This script must be run as the "git" user\n')
        parser.print_help()
        sys.exit(1)
    args = parser.parse_args()
    args.src = os.path.abspath(args.src)
    # store app ID
    args.app = '/'.join(args.src.split(os.path.sep)[-1:]).replace('.git', '')
    return args


def puts_step(s):
    sys.stdout.write("-----> %s\n" % s)
    sys.stdout.flush()


def puts_line():
    sys.stdout.write("\n")
    sys.stdout.flush()


def puts(s):
    sys.stdout.write("       " + s)
    sys.stdout.flush()


def exit_on_error(error_code, msg):
    sys.stderr.write(msg)
    sys.stderr.write('\n')
    sys.stderr.flush()
    sys.exit(error_code)


if __name__ == '__main__':
    args = parse_args()
    # get sha of master
    try:
        with open(os.path.join(args.src, 'refs/heads/master')) as f:
            sha = f.read().strip('\n')
    except IOError:
        exit_on_error(2, 'Could not read the repository SHA--is the refspec correct?')
    # prepare for buildpack run
    slug_path = os.path.join(SLUG_DIR, '{0}-{1}.tar.gz'.format(args.app, sha))
    # create cache dir
    cache_dir = os.path.join(args.src, 'cache')
    if not os.path.exists(cache_dir):
        os.mkdir(cache_dir)
    # build/compile
    cmd = "git archive master | docker run -i -a stdin" \
          " -v {cache_dir}:/tmp/cache:rw " \
          " -v /opt/deis/build/packs:/tmp/buildpacks:rw " \
          " deis/slugbuilder"
    cmd = cmd.format(**locals())
    p = subprocess.Popen(cmd, cwd=args.src, shell=True, stdout=subprocess.PIPE)
    rc = p.wait()
    if rc != 0:
        exit_on_error(rc, 'Could not execute build, leaving current release in place')
    container = p.stdout.read().strip('\n')
    # attach to container and wait for build
    cmd = 'docker attach {container}'.format(**locals())
    p = subprocess.Popen(cmd, cwd=args.src, shell=True)
    rc = p.wait()
    if rc != 0:
        exit_on_error(rc, 'Build failed, leaving current release in place')
    # extract slug
    cmd = 'docker cp {container}:/tmp/slug.tgz .'.format(**locals())
    p = subprocess.Popen(cmd, cwd=args.src, shell=True)
    rc = p.wait()
    if rc != 0:
        exit_on_error(rc, 'Could not extract slug from container')
    os.rename(os.path.join(args.src, 'slug.tgz'), slug_path)
    # extract procfile
    cmd = 'tar xfO {slug_path} ./Procfile'.format(**locals())
    p = subprocess.Popen(cmd, cwd=args.src, shell=True, stdout=subprocess.PIPE)
    rc = p.wait()
    if rc != 0:
        exit_on_error(rc, 'Could not extract Procfile from container')
    try:
        procfile = yaml.safe_load(p.stdout.read())
    except YAMLError as e:
        exit_on_error(1, 'Invalid Procfile format: {0}'.format(e))
    # extract release
    cmd = 'tar xfO {slug_path} ./.release'.format(**locals())
    p = subprocess.Popen(cmd, cwd=args.src, shell=True, stdout=subprocess.PIPE)
    rc = p.wait()
    if rc != 0:
        exit_on_error(rc, 'Could not extract Release from container')
    try:
        release = yaml.safe_load(p.stdout.read())
    except YAMLError as e:
        exit_on_error(1, 'Invalid Release format: {0}'.format(e))
    # remove the container
    cmd = 'docker rm {container}'.format(**locals())
    p = subprocess.Popen(cmd, cwd=args.src, shell=True, stdout=subprocess.PIPE)
    rc = p.wait()
    if rc != 0:
        exit_on_error(rc, 'Could not remove build container')
    # calculate checksum
    p = subprocess.Popen(['sha256sum', slug_path], stdout=subprocess.PIPE)
    rc = p.wait()
    if rc != 0:
        exit_on_error(rc, 'Could not calculate SHA of slug')
    checksum = p.stdout.read().split(' ')[0]
    # prepare the push-hook
    push = {'username': args.user, 'app': args.app, 'sha': sha, 'checksum': checksum,
            'config': release.get('config_vars', {})}
    # TODO: why can't we run this with `sudo -u deis`?
    output = subprocess.check_output(
        ['sudo', '-u', 'deis', "{}/bin/pre-push-hook".format(CONTROLLER_DIR)])
    data = json.loads(output)
    ip = data['domain']
    push['url'] = "http://{ip}/slugs/{args.app}-{sha}.tar.gz".format(**locals())
    push['procfile'] = procfile
    # calculate slug size
    push['size'] = os.stat(slug_path).st_size
    puts_line()
    # run stage
    sys.stdout.write("       " + "Launching... ")
    sys.stdout.flush()
    p = subprocess.Popen(['sudo', '-u', 'deis', "{}/bin/push-hook".format(CONTROLLER_DIR)],
                         stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    stdout, stderr = p.communicate(json.dumps(push))
    rc = p.wait()
    if rc != 0:
        raise RuntimeError('Build error {0}'.format(stderr))
    databag = json.loads(stdout)
    sys.stdout.write("done, v{}\n".format(databag['release']['version']))
    sys.stdout.flush()
    puts_line()
    puts_step("{args.app} deployed to Deis".format(**locals()))
    domains = databag.get('domains', [])
    if domains:
        for domain in domains:
            puts("http://{domain}\n".format(**locals()))
    else:
        puts('No proxy nodes found for this formation.\n')
    puts_line()
    puts('To learn more, use `deis help` or visit http://deis.io')
    puts_line()
    puts_line()
