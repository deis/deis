"""
URL routing patterns for the Deis project.

This is the "master" urls.py which then includes the urls.py files of
installed apps.
"""

from __future__ import unicode_literals

from django.conf.urls import patterns, include, url
from django.contrib import admin


admin.autodiscover()


# pylint: disable=C0103
urlpatterns = patterns(
    '',
    url(r'^accounts/', include('allauth.urls')),
    url(r'^admin/doc/', include('django.contrib.admindocs.urls')),
    url(r'^admin/', include(admin.site.urls)),
    url(r'^api/', include('api.urls')),
    url(r'^', include('web.urls')),
)
