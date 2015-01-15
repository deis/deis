Due to changes introduced in Docker 1.3.1 related to insecure Docker registries, the hosts running
Deis must be able to communicate via a private network in one of the RFC 1918 or RFC 6598 private
address spaces: ``10.0.0.0/8``, ``172.16.0.0/12``, ``192.168.0.0/16``, or ``100.64.0.0/10``.
