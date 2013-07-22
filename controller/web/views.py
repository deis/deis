
from django.contrib.auth.decorators import login_required
from django.shortcuts import render

from api.models import Formation


@login_required
def account(request):
    """Return the user's account web page."""
    return render(request, 'web/account.html', {
        'current_page': 'account',
    })


@login_required
def dashboard(request):
    """Return the user's dashboard web page."""
    formations = Formation.objects.filter(owner=request.user)
    return render(request, 'web/dashboard.html', {
        'current_page': 'dashboard',
        'formations': formations,
    })


@login_required
def formations(request):
    """Return the user's formations web page."""
    formations = Formation.objects.filter(owner=request.user)
    return render(request, 'web/formations.html', {
        'current_page': 'formations',
        'formations': formations
    })


@login_required
def docs(request):
    """Return the documentation index."""
    return render(request, 'web/docs.html', {
        'current_page': 'docs',
    })


@login_required
def support(request):
    """Return the support ticket system home page."""
    return render(request, 'web/support.html', {
        'current_page': 'support',
    })
