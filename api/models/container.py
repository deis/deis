from __future__ import unicode_literals
import logging

from django.conf import settings
from django.db import models
from django.utils.encoding import python_2_unicode_compatible


from application import App, log_event
from base import UuidAuditedModel
from formation import Formation
from node import Node

logger = logging.getLogger(__name__)


class ContainerManager(models.Manager):

    def scale(self, app, structure, **kwargs):
        """Scale containers up or down to match requested."""
        requested_containers = structure.copy()
        formation = app.formation
        # increment new container nums off the most recent container
        all_containers = app.container_set.all().order_by('-created')
        container_num = 1 if not all_containers else all_containers[0].num + 1
        msg = 'Containers scaled ' + ' '.join(
            "{}={}".format(k, v) for k, v in requested_containers.items())
        # iterate and scale by container type (web, worker, etc)
        changed = False
        for container_type in requested_containers.keys():
            containers = list(app.container_set.filter(type=container_type).order_by('created'))
            requested = requested_containers.pop(container_type)
            diff = requested - len(containers)
            if diff == 0:
                continue
            changed = True
            while diff < 0:
                # get the next node with the most containers
                node = Node.objects.next_runtime_node(
                    formation, container_type, reverse=True)
                # delete a container attached to that node
                for c in containers:
                    if node == c.node:
                        containers.remove(c)
                        c.delete()
                        diff += 1
                        break
            while diff > 0:
                # get the next node with the fewest containers
                node = Node.objects.next_runtime_node(formation, container_type)
                port = Node.objects.next_runtime_port(formation)
                c = Container.objects.create(owner=app.owner,
                                             formation=formation,
                                             node=node,
                                             app=app,
                                             type=container_type,
                                             num=container_num,
                                             port=port)
                containers.append(c)
                container_num += 1
                diff -= 1
        log_event(app, msg)
        return changed

    def balance(self, formation, **kwargs):
        runtime_nodes = formation.node_set.filter(layer__runtime=True).order_by('created')
        all_containers = self.filter(formation=formation).order_by('-created')
        # get the next container number (e.g. web.19)
        container_num = 1 if not all_containers else all_containers[0].num + 1
        changed = False
        app = None
        # iterate by unique container type
        for container_type in set([c.type for c in all_containers]):
            # map node container counts => { 2: [b3, b4], 3: [ b1, b2 ] }
            n_map = {}
            for node in runtime_nodes:
                ct = len(node.container_set.filter(type=container_type))
                n_map.setdefault(ct, []).append(node)
            # loop until diff between min and max is 1 or 0
            while max(n_map.keys()) - min(n_map.keys()) > 1:
                # get the most over-utilized node
                n_max = max(n_map.keys())
                n_over = n_map[n_max].pop(0)
                if len(n_map[n_max]) == 0:
                    del n_map[n_max]
                # get the most under-utilized node
                n_min = min(n_map.keys())
                n_under = n_map[n_min].pop(0)
                if len(n_map[n_min]) == 0:
                    del n_map[n_min]
                # delete the oldest container from the most over-utilized node
                c = n_over.container_set.filter(type=container_type).order_by('created')[0]
                app = c.app  # pull ref to app for recreating the container
                c.delete()
                # create a container on the most under-utilized node
                self.create(owner=formation.owner,
                            formation=formation,
                            app=app,
                            type=container_type,
                            num=container_num,
                            node=n_under,
                            port=Node.objects.next_runtime_port(formation))
                container_num += 1
                # update the n_map accordingly
                for n in (n_over, n_under):
                    ct = len(n.container_set.filter(type=container_type))
                    n_map.setdefault(ct, []).append(n)
                changed = True
        if app:
            log_event(app, 'Containers balanced')
        return changed


@python_2_unicode_compatible
class Container(UuidAuditedModel):
    """
    Docker container used to securely host an application process.
    """

    objects = ContainerManager()

    owner = models.ForeignKey(settings.AUTH_USER_MODEL)
    formation = models.ForeignKey(Formation)
    node = models.ForeignKey(Node)
    app = models.ForeignKey(App)
    type = models.CharField(max_length=128)
    num = models.PositiveIntegerField()
    port = models.PositiveIntegerField()

    # TODO: add celery beat tasks for monitoring node health
    status = models.CharField(max_length=64, default='up')

    def short_name(self):
        return "{}.{}".format(self.type, self.num)
    short_name.short_description = 'Name'

    def __str__(self):
        return "{0} {1}".format(self.formation.id, self.short_name())

    class Meta:
        app_label = 'api'
        get_latest_by = '-created'
        ordering = ['created']
        unique_together = (('app', 'type', 'num'),
                           ('formation', 'port'))
