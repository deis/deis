"""
Unit tests for the Deis CLI providers commands.

Run these tests with "python -m unittest client.tests.test_providers"
or with "./manage.py test client.ProvidersTest".
"""

from __future__ import unicode_literals
from unittest import TestCase

import pexpect

from .utils import DEIS
from .utils import purge
from .utils import register


class ProvidersTest(TestCase):

    @classmethod
    def setUpClass(cls):
        cls.username, cls.password = register()

    @classmethod
    def tearDownClass(cls):
        purge(cls.username, cls.password)

    def test_seeded(self):
        """Test that our autotest user has some providers auto-seeded."""
        child = pexpect.spawn("{} providers".format(DEIS))
        child.expect(".* => .*")
        child.expect(pexpect.EOF)
