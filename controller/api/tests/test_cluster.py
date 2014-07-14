"""
Unit tests for the Deis api app.

Run the tests with "./manage.py test api"
"""

from __future__ import unicode_literals

import json

from django.test import TestCase
from django.test.utils import override_settings


@override_settings(CELERY_ALWAYS_EAGER=True)
class ClusterTest(TestCase):

    """Tests cluster management"""

    fixtures = ['tests.json']

    def setUp(self):
        self.assertTrue(
            self.client.login(username='autotest', password='password'))

    def test_cluster(self):
        """
        Test that an administrator can create, read, update and delete a cluster
        """
        url = '/api/clusters'
        options = {'key': 'val'}
        body = {'id': 'autotest', 'domain': 'autotest.local', 'type': 'mock',
                'hosts': 'host1,host2', 'auth': 'base64string', 'options': options}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        cluster_id = response.data['id']  # noqa
        self.assertIn('owner', response.data)
        self.assertIn('id', response.data)
        self.assertIn('domain', response.data)
        self.assertIn('hosts', response.data)
        self.assertIn('auth', response.data)
        self.assertIn('options', response.data)
        self.assertEqual(response.data['hosts'], 'host1,host2')
        self.assertEqual(json.loads(response.data['options']), {'key': 'val'})
        response = self.client.get('/api/clusters')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(len(response.data['results']), 1)
        # ensure we can delete the cluster with an app
        # see https://github.com/deis/deis/issues/927
        url = '/api/apps'
        body = {'cluster': 'autotest'}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 201)
        url = '/api/clusters/{cluster_id}'.format(**locals())
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
        new_hosts, new_options = 'host2,host3', {'key': 'val2'}
        body = {'hosts': new_hosts, 'options': new_options}
        response = self.client.patch(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['hosts'], new_hosts)
        self.assertEqual(json.loads(response.data['options']), new_options)
        response = self.client.delete(url)
        self.assertEqual(response.status_code, 204)

    def test_cluster_perms_denied(self):
        """
        Test that a user cannot make changes to a cluster
        """
        url = '/api/clusters'
        options = {'key': 'val'}
        self.client.login(username='autotest2', password='password')
        body = {'id': 'autotest2', 'domain': 'autotest.local', 'type': 'mock',
                'hosts': 'host1,host2', 'auth': 'base64string', 'options': options}
        response = self.client.post(url, json.dumps(body), content_type='application/json')
        self.assertEqual(response.status_code, 403)

    def test_user_read_access(self):
        """
        Test that a user has read-only access to the clusters endpoint
        """
        self.client.login(username='autotest2', password='password')
        url = '/api/clusters'
        response = self.client.get(url)
        self.assertEqual(response.status_code, 200)
