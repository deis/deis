"""
Unit tests for the Deis Chef module

Run the tests with "python -m unittest discover"
"""

from __future__ import unicode_literals

import json
import unittest

from celerytasks.chef import ChefAPI
from deis import settings


@unittest.skip('Need to set up TEST_CHEF_SERVER somehow.')
class ChefAPITest(unittest.TestCase):
    """Tests the client interface to Chef Server API."""

    def setUp(self):
        self.client = ChefAPI(
            settings.TEST_CHEF_SERVER_URL, settings.TEST_CHEF_CLIENT_NAME,
            settings.TEST_CHEF_CLIENT_KEY)

    def test_databag(self):
        databag_name = 'testing'
        ditem_name = 'item1'
        ditem_value = {'something': 1, 'else': 2}

        # delete the databag to make sure we are creating a new one
        resp, status = self.client.delete_databag(databag_name)

        resp, status = self.client.create_databag(databag_name)
        self.assertEqual(status, 201)
        self.assertTrue(resp)

        resp = self.client.create_databag_item(
            databag_name, ditem_name, ditem_value)
        self.assertEqual(status, 201)
        self.assertTrue(resp)

        resp, status = self.client.get_databag(databag_name)
        self.assertEqual(status, 200)
        resp, status = self.client.get_databag_item(databag_name, ditem_name)
        self.assertEqual(status, 200)

        ditem_value = json.loads(resp)
        ditem_value['newvalue'] = 'databag'
        resp, status = self.client.update_databag_item(
            databag_name, ditem_name, ditem_value)
        self.assertEqual(status, 200)
        resp, status = self.client.get_databag_item(databag_name, ditem_name)
        self.assertEqual(status, 200)
        self.assertTrue('newvalue' in json.loads(resp))
