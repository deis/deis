from __future__ import unicode_literals

import unittest

from celerytasks import azuresms
# from deis import settings


class AzureTest(unittest.TestCase):
    """Tests the client interface to Chef Server API."""

    @unittest.skip('Windows Azure is not yet supported.')
    def test_launch(self):
        l = azuresms.launch_node(None, None, None, None, None, None)
        print "L is ", l
