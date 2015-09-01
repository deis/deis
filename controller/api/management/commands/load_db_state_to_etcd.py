from django.core.management.base import BaseCommand

from api.models import Key, App, Domain, Certificate, Config


class Command(BaseCommand):
    """Management command for publishing Deis platform state from the database
    to etcd.
    """
    def handle(self, *args, **options):
        """Publishes Deis platform state from the database to etcd."""
        print "Publishing DB state to etcd..."
        for model in (Key, App, Domain, Certificate, Config):
            for obj in model.objects.all():
                obj.save()
        print "Done Publishing DB state to etcd."
