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
import boto.route53
import docopt
import json
import subprocess
import time


def parse_args():
    return docopt.docopt(__doc__)

def create_cname(args):
    conn = boto.route53.connect_to_region(args['--region'])
    zone = conn.get_zone(args['<zone>'])
    status = zone.add_cname(args['<name>'], args['<value>'], ttl=args['--ttl'])
    print("waiting for record to sync: {}".format(status))
    while status.update() != "INSYNC":
        time.sleep(2)
    print(status)

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

    # calculate zone from provided name
    zone = '.'.join(args['<name>'].split('.')[-2:])
    args['<zone>'] = zone

    # add a * to the provided name
    args['<name>'] = u'\\052.'+unicode(args['<name>'])

    if args['<action>'] == 'create':
        create_cname(args)
    elif args['<action>'] == 'delete':
        delete_cname(args)
    else:
        raise NotImplementedError
