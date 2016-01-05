:title: Platform monitoring
:description: Configuring platform monitoring.

.. _platform_monitoring:

Platform monitoring
===================

While Deis itself doesn't have a built-in monitoring system, Deis components and deployed
applications alike run entirely within Docker containers. This means that monitoring tools and
services which support Docker containers should work with Deis. A few tools and monitoring services
which support Docker integrations are detailed below.

Tools
-----

cadvisor
~~~~~~~~

Google's Container Advisor (`cadvisor`_) runs inside a Docker container and shows memory and CPU
usage for all containers running on the host. To run cAdvisor:

.. code-block:: console

    sudo docker run \
    --volume=/:/rootfs:ro \
    --volume=/var/run:/var/run:rw \
    --volume=/sys:/sys:ro \
    --volume=/var/lib/docker/:/var/lib/docker:ro \
    --publish=8080:8080 \
    --detach=true \
    --name=cadvisor \
    google/cadvisor:latest

To run cAdvisor on all hosts in the cluster, you can submit and start a fleet service:

.. code-block:: console

    [Unit]
    Description=Google Container Advisor
    Requires=docker.socket
    After=docker.socket

    [Service]
    ExecStartPre=/bin/sh -c "docker history google/cadvisor:latest >/dev/null 2>&1 || docker pull google/cadvisor:latest"
    ExecStartPre=/bin/sh -c "docker inspect cadvisor >/dev/null 2>&1 && docker rm -f cadvisor || true"
    ExecStart=/usr/bin/docker run --volume=/:/rootfs:ro --volume=/var/run:/var/run:rw --volume=/sys:/sys:ro --volume=/var/lib/docker/:/var/lib/docker:ro --publish=8080:8080 --name=cadvisor google/cadvisor:latest
    ExecStopPost=-/usr/bin/docker rm -f cadvisor
    Restart=on-failure
    RestartSec=5

    [Install]
    WantedBy=multi-user.target

    [X-Fleet]
    Global=true

Save the file as ``cadvisor.service``. Load and start the service with
``fleetctl load cadvisor.service && fleetctl start cadvisor.service``.

The web interface will be accessible at port 8080 on each host.

In addition to starting a cAdvisor instance on each CoreOS host, there's also a project called
`heapster`_ from the Google Cloud Platform team, which seems to be a cluster-aware cAdvisor.

Monitoring services
-------------------

These are a few monitoring services which are known to provide Docker integrations.
Additions to this reference guide are much appreciated!

Datadog
~~~~~~~

The `Datadog`_ cloud monitoring service provides a monitor agent which runs on the host and provides
metrics for all Docker containers (which is functionally similar to cAdvisor's implementation).
See `this blog post`_ for details. The `Datadog agent`_ for Docker can be run on a single host as
follows:

.. code-block:: console

    docker run -d --privileged --name dd-agent -h `hostname` -v /var/run/docker.sock:/var/run/docker.sock -v /proc/mounts:/host/proc/mounts:ro -v /sys/fs/cgroup/:/host/sys/fs/cgroup:ro -e API_KEY=YOUR_REAL_API_KEY datadog/docker-dd-agent

Be sure to substitute ``YOUR_REAL_API_KEY`` for your Datadog API key.

To run Datadog for the entire cluster, you can submit and start a fleet service (again, substitute ``YOUR_REAL_API_KEY``):

.. code-block:: console

    [Unit]
    Description=Datadog
    Requires=docker.socket
    After=docker.socket

    [Service]
    ExecStartPre=/bin/sh -c "docker history datadog/docker-dd-agent:latest >/dev/null || docker pull datadog/docker-dd-agent:latest"
    ExecStart=/usr/bin/docker run --privileged --name dd-agent -h %H -v /var/run/docker.sock:/var/run/docker.sock -v /proc/mounts:/host/proc/mounts:ro -v /sys/fs/cgroup/:/host/sys/fs/cgroup:ro -e API_KEY=YOUR_REAL_API_KEY datadog/docker-dd-agent

    [Install]
    WantedBy=multi-user.target

    [X-Fleet]
    Global=true

Save the file as ``datadog.service``. Load and start the service with
``fleetctl load datadog.service && fleetctl start datadog.service``.

Shortly thereafter, you should start to see metrics from your Deis cluster appear in your Datadog dashboard.

New Relic
~~~~~~~~~

The `New Relic`_ monitoring service's agent will run on the CoreOS host and report metrics to New Relic.

The Deis community's own Johannes Würbach has developed a fleet service for New Relic in his
`newrelic-sysmond`_ repository.

SPM Performance Monitoring
~~~~~~~~~~~~~~~~~~~~~~~~~~

`SPM for Docker`_ provides Monitoring, Anomaly Detection and Alerting. Detailed metrics and events are collected for the host machine and running containers (e.g. cpu, memory, IO, network, limits). The SPM agent runs in a Docker container.

`Create a new SPM App`_ of type "Docker" and copy the Application Token 

Set the Application Token (SPM_TOKEN) via etcd to make it available in the cluster.

.. code-block:: console

    etcdctl set /SPM_TOKEN YOUR_SPM_APP_TOKEN

SPM Docker Agent can be run using Docker as follows (assuming the SPM_TOKEN is set via etcdctl):

.. code-block:: console

    docker run -d --name sematext-agent -e SPM_TOKEN=`etcdctl get SPM_TOKEN` -e HOSTNAME=$HOSTNAME -v /var/run/docker.sock:/var/run/docker.sock sematext/sematext-agent-docker
    

To activate SPM Docker Agent for the entire cluster submit this unit file to fleet

.. code-block:: console

    [Unit]
    Description=SPM Docker Agent
    After=docker.service
    Requires=docker.service

    [Service]
    TimeoutStartSec=0
    EnvironmentFile=/etc/environment
    Restart=always
    RestartSec=30s
    ExecStartPre=-/usr/bin/docker kill sematext-agent
    ExecStartPre=-/usr/bin/docker rm sematext-agent
    ExecStartPre=/usr/bin/docker pull sematext/sematext-agent-docker:latest
    ExecStart=/bin/sh -c 'set -ex; /usr/bin/docker run --name sematext-agent -e SPM_TOKEN=$(etcdctl get SPM_TOKEN) -e HOSTNAME=$HOSTNAME -v /var/run/docker.sock:/var/run/docker.sock sematext/sematext-agent-docker'
    ExecStop=/usr/bin/docker stop sematext-agent

    [Install]
    WantedBy=multi-user.target
    
    [X-Fleet]
    Global=true

Save the file as ``sematext-agent.service``. 

.. code-block:: console

    wget https://raw.githubusercontent.com/sematext/sematext-agent-docker/master/coreos/sematext-agent.service

Load and start the service with

.. code-block:: console

    fleetctl load sematext-agent.service && fleetctl start sematext-agent.service 

After one minute, you should see metrics in SPM. 

Documentation, source code and support information is available here:
`https://github.com/sematext/sematext-agent-docker`_. 



.. _`cadvisor`: https://github.com/google/cadvisor
.. _`Datadog`: https://www.datadoghq.com
.. _`Datadog agent`: https://github.com/DataDog/docker-dd-agent
.. _`heapster`: https://github.com/GoogleCloudPlatform/heapster/blob/master/clusters/coreos/README.md
.. _`this blog post`: https://www.datadoghq.com/2014/06/monitor-docker-datadog/
.. _`New Relic`: http://newrelic.com/
.. _`newrelic-sysmond`: https://github.com/johanneswuerbach/newrelic-sysmond-service
.. _`SPM for Docker`: http://sematext.com/spm/integrations/docker-monitoring.html
.. _`Create a new SPM App`: https://apps.sematext.com/spm-reports/registerApplication.do
.. _`https://github.com/sematext/sematext-agent-docker`: https://github.com/sematext/sematext-agent-docker
