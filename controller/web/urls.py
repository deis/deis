"""
URL patterns and routing for the Deis web app.
"""

from __future__ import unicode_literals

from django.conf.urls import patterns
from django.conf.urls import url


urlpatterns = patterns(
    'web.views',
    url(r'^$', 'dashboard', name='dashboard'),
    url(r'^account/$', 'account', name='account'),
    url(r'^apps/$', 'apps', name='apps'),
    url(r'^support/$', 'support', name='support'),
)
