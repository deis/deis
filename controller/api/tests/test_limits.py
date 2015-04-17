import unittest

from api.serializers import MEMLIMIT_MATCH


class TestLimits(unittest.TestCase):
    """Tests the regex for unit format used by "deis limits:set --memory=<limit>".
    """

    def test_memlimit_regex(self):
        self.assertTrue(MEMLIMIT_MATCH.match("20MB"))
        self.assertFalse(MEMLIMIT_MATCH.match("20MK"))
        self.assertTrue(MEMLIMIT_MATCH.match("20gb"))
        self.assertFalse(MEMLIMIT_MATCH.match("20gK"))
