#!/usr/bin/env python
"""
The route53-wildcard utility is used to add and remove wildcard
DNS entries from an AWS Route53 zonefile.

Usage: route53-wildcard.py <action> <name> <value> [options]

Options:

  --zone=<zone>      name of zone to use [defaults to parsing name]
  --region=<region>  AWS region to use [default: us-west-2].
  --ttl=<ttl>        record TTL to use [default: 60].

Examples:

./route53-wildcard.py create mycluster.gabrtv.io deis-elb.blah.amazonaws.com
./route53-wildcard.py delete mycluster.gabrtv.io deis-elb.blah.amazonaws.com
"""

from __future__ import print_function

import boto.route53
import docopt
import socket
import time
import uuid


def parse_args():
    return docopt.docopt(__doc__)


def create_cname(args):
    conn = boto.route53.connect_to_region(args['--region'])
    zone = conn.get_zone(args['<zone>'])
    name = args['<name>']
    status = zone.add_cname(name, args['<value>'], ttl=args['--ttl'])
    print("waiting for record to sync: {}".format(status))
    while status.update() != "INSYNC":
        time.sleep(2)
    print(status)
    print('waiting for wildcard domain to become available...', end='')
    # AWS docs say it can take up to 30 minutes for route53 changes to happen, although
    # it seems to be almost immediate.
    for i in xrange(120):
        try:
            random_hostname = str(uuid.uuid4())[:8]
            if socket.gethostbyname("{}.{}".format(random_hostname, name)):
                print('ok')
                break
        except socket.gaierror:
            time.sleep(15)


def delete_cname(args, block=False):
    conn = boto.route53.connect_to_region(args['--region'])
    zone = conn.get_zone(args['<zone>'])
    status = zone.delete_cname(args['<name>'])
    if block:
        print("waiting for record to sync: {}".format(status))
        while status.update() != "INSYNC":
            time.sleep(2)
        print(status)

if __name__ == '__main__':
    args = parse_args()

    if args['--zone'] == None:
        zone = '.'.join(args['<name>'].split('.')[-2:])
        args['<zone>'] = zone
    else:
        args['<zone>'] = args['--zone']

    # add a * to the provided name
    args['<name>'] = u'\\052.'+unicode(args['<name>'])

    if args['<action>'] == 'create':
        create_cname(args)
    elif args['<action>'] == 'delete':
        delete_cname(args)
    else:
        raise NotImplementedError
