
from django.core.management.base import BaseCommand
from django.core.management.base import CommandError


class Command(BaseCommand):

    help = "Launches the Deis command-line client."

    def handle(self, *args, **_kwargs):
        # pylint: disable=E0611,F0401
        try:
            from deis.client import DeisClient
        except ImportError:
            raise CommandError(
                'Please install the Deis command-line client.')
        else:
            return DeisClient().main()
