"""
Unit tests for the Deis CLI auth commands.

Run these tests with "python -m unittest client.tests.test_misc"
or with "./manage.py test client.HelpTest client.VersionTest".
"""

from __future__ import unicode_literals
from unittest import TestCase

import pexpect

from client.deis import __version__
from .utils import DEIS


class HelpTest(TestCase):
    """Test that the client can document its own behavior."""

    def test_deis(self):
        """Test that the `deis` command on its own returns usage."""
        child = pexpect.spawn(DEIS)
        child.expect('Usage: deis <command> \[<args>...\]')

    def test_help(self):
        """Test that the client reports its help message."""
        child = pexpect.spawn('{} --help'.format(DEIS))
        child.expect('The Deis command-line client.*to an application\.')
        child = pexpect.spawn('{} -h'.format(DEIS))
        child.expect('The Deis command-line client.*to an application\.')
        child = pexpect.spawn('{} help'.format(DEIS))
        child.expect('The Deis command-line client.*to an application\.')


class VersionTest(TestCase):
    """Test that the client can report its version string."""

    def test_version(self):
        """Test that the client reports its version string."""
        child = pexpect.spawn('{} --version'.format(DEIS))
        child.expect("Deis CLI {}".format(__version__))
