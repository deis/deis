"""
Unit tests for the Deis web app.

Run the tests with "./manage.py test web"
"""

from __future__ import unicode_literals

from django.test import TestCase

# pylint: disable=R0904


class SimpleTest(TestCase):

    def test_basic_addition(self):
        """
        Tests that 1 + 1 always equals 2.
        """
        self.assertEqual(1 + 1, 2)
