Whether built for evaluation or to host production applications, when managing a
small Deis cluster (three to five nodes), it is reasonable to accept the
platform's default behavior wherein the Control Plane, Data Plane, and Router
Mesh are not isolated from one another. (See :ref:`architecture`.) This means
Control Plane components such as the :ref:`controller` or :ref:`database` will
be eligible to run on any node, as will the Router Mesh and the Data Plane
components such as :ref:`logspout`, :ref:`publisher`, and deployed applications.

In larger clusters however, nodes are more easily thought of as a commodity.
Operators may scale clusters out to meet demand or in to conserve resources. In
such cases, it is beneficial to isolate the Control Plane, which has no
significant need to scale (and optionally, the Router Mesh) to a small, fixed
number of nodes that are exempt from such scaling events.  This eliminates the
possibility that Control Plane components running on a decommissioned node will
experience downtime as they are rescheduled.  Additionally, this reserves the
resources of a large (and possibly dynamic) pool of nodes for the workloads that
are most likely to scale-- applications.
