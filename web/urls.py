
from __future__ import unicode_literals

from django.conf.urls import patterns
from django.conf.urls import url


urlpatterns = patterns(
    'web.views',
    url(r'^$', 'dashboard', name='dashboard'),
    url(r'^account/$', 'account', name='account'),
    url(r'^docs/$', 'docs', name='docs'),
    url(r'^formations/$', 'formations', name='formations'),
    url(r'^support/$', 'support', name='support'),
)
