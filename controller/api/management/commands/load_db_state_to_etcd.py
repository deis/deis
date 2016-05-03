from __future__ import print_function

from django.core.management.base import BaseCommand

from api.models import Key, App, Domain, Certificate


class Command(BaseCommand):
    """Management command for publishing Deis platform state from the database
    to etcd.
    """
    def handle(self, *args, **options):
        """Publishes Deis platform state from the database to etcd."""
        print("Publishing DB state to etcd...")
        for app in App.objects.all():
            app.save()
            app.config_set.latest().save()
        for model in (Key, Domain, Certificate):
            for obj in model.objects.all():
                obj.save()
        print("Done Publishing DB state to etcd.")
