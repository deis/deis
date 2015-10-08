"""
Unit tests for the Deis web app.

Run the tests with "./manage.py test web"
"""

from __future__ import unicode_literals

from django.conf import settings
from django.template import Context
from django.template import Template
from django.template import TemplateSyntaxError
from django.test import TestCase


class WebViewsTest(TestCase):

    fixtures = ['test_web.json']

    @classmethod
    def setUpClass(cls):
        settings.WEB_ENABLED = True

    def setUp(self):
        self.client.login(username='autotest-1', password='password')

    def test_account(self):
        response = self.client.get('/account/')
        self.assertContains(response, '<title>Deis | Account</title>', html=True)
        self.assertContains(response, 'autotest-1')
        self.assertContains(response, '<img src="//www.gravatar.com/avatar')

    def test_dashboard(self):
        response = self.client.get('/')
        self.assertContains(response, '<title>Deis | Dashboard</title>', html=True)

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


class GravatarTagsTest(TestCase):

    def _render_template(self, t, ctx=None):
        """Test that the tag renders a gravatar URL."""
        tmpl = Template(t)
        return tmpl.render(Context(ctx)).strip()

    def test_render(self):
        tmpl = """\
{% load gravatar_tags %}
{% gravatar_url email %}
"""
        rendered = self._render_template(tmpl, {'email': 'github@deis.io'})
        self.assertEquals(
            rendered,
            r'//www.gravatar.com/avatar/058ff74579b6a8fa1e10ab98c990e945?s=24&d=mm')

    def test_render_syntax_error(self):
        """Test that the tag requires one argument."""
        tmpl = """
{% load gravatar_tags %}
{% gravatar_url %}
"""
        self.assertRaises(TemplateSyntaxError, self._render_template, tmpl)

    def test_render_context_error(self):
        """Test that an empty email returns an empty string."""
        tmpl = """
{% load gravatar_tags %}
{% gravatar_url email %}
"""
        rendered = self._render_template(tmpl, {})
        self.assertEquals(rendered, '')
