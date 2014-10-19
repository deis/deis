"""
View classes for presenting Deis web pages.
"""

from django.contrib.auth.decorators import login_required
from django.shortcuts import render

from api.models import App
from deis import __version__


@login_required
def account(request):
    """Return the user's account web page."""
    return render(request, 'web/account.html', {
        'page': 'account',
    })


@login_required
def dashboard(request):
    """Return the user's dashboard web page."""
    apps = App.objects.filter(owner=request.user)
    return render(request, 'web/dashboard.html', {
        'page': 'dashboard',
        'apps': apps,
        'version': __version__,
    })


@login_required
def apps(request):
    """Return the user's apps web page."""
    apps = App.objects.filter(owner=request.user)
    return render(request, 'web/apps.html', {
        'page': 'apps',
        'apps': apps,
    })


@login_required
def support(request):
    """Return the support ticket system home page."""
    return render(request, 'web/support.html', {
        'page': 'support',
    })
