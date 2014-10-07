import logging

from restkit.contrib.wsgi_proxy import HostProxy
from webob import Request

logger = logging.getLogger("api")

# TODO: hardcoded
proxy = HostProxy("http://172.17.8.100:5000")

def forward_to_registry(environ, start_response):
  req = Request(environ)
  resp = req.get_response(proxy)
  return resp(environ, start_response)
