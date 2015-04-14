import unittest
import re


MEMLIMIT = re.compile(r'^(?P<mem>[0-9]+(MB|KB|GB|[BKMG]))$', re.IGNORECASE)


class TestLimits(unittest.TestCase):

    def test_upper(self):
        self.assertTrue(MEMLIMIT.match("20MB"))
        self.assertFalse(MEMLIMIT.match("20MK"))
        self.assertTrue(MEMLIMIT.match("20gb"))
        self.assertFalse(MEMLIMIT.match("20gK"))
