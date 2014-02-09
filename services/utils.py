import random
import string
import urllib


def gen_random_string(length=16):
    return ''.join([random.choice(string.ascii_uppercase + string.digits)
                    for x in range(length)]).lower()


def get_external_ip_address():
    return urllib.urlopen('http://ipv4.icanhazip.com/').read().rstrip()
