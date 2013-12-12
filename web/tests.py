"""
Unit tests for the Deis web app.

Run the tests with "./manage.py test web"
"""

from __future__ import unicode_literals

from django.test import TestCase


class WebViewsTest(TestCase):

    fixtures = ['test_web.json']

    def setUp(self):
        self.client.login(username='autotest-1', password='password')

    def test_account(self):
        response = self.client.get('/account/')
        self.assertContains(response, '<title>Deis | Account</title>', html=True)
        self.assertContains(response, 'autotest-1')
        self.assertContains(response, '<img src="//www.gravatar.com/avatar')
        self.assertContains(
            response, '<form method="post" action="/accounts/logout/">')

    def test_dashboard(self):
        response = self.client.get('/')
        self.assertContains(response, '<title>Deis | Dashboard</title>', html=True)
        self.assertContains(
            response,
            r'You have <a href="/formations/">one formation</a> and <a href="/apps/">one app</a>.')

    def test_formations(self):
        response = self.client.get('/formations/')
        self.assertContains(response, '<title>Deis | Formations</title>', html=True)
        self.assertContains(response, '<h1>One Formation</h1>')
        self.assertContains(response, '<h3>autotest-1</h3>')
        self.assertContains(response, '<dt>Owned by</dt>')
        self.assertContains(response, '<dd>autotest-1</dd>')

    def test_apps(self):
        response = self.client.get('/apps/')
        self.assertContains(response, '<title>Deis | Apps</title>', html=True)
        self.assertContains(response, '<h1>One App</h1>')
        self.assertContains(response, '<h3>autotest-1-app</h3>')

    def test_support(self):
        response = self.client.get('/support/')
        self.assertContains(response, '<title>Deis | Support</title>', html=True)
        self.assertContains(response, '<div class="forkImage">')
        self.assertContains(response, '<h2>IRC</h2>')
        self.assertContains(response, '<h2>GitHub</h2>')
