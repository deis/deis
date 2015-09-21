:title: Controller API v1.6
:description: The v1.6 REST API for Deis' Controller

Controller API v1.6
===================

This is the v1.6 REST API for the :ref:`Controller`.


What's New
----------

**New!** administrators no longer have to supply a password when changing another user's password.

**New!** administrators no longer have to supply a password when deleting another user.

**New!** ``?page_size`` query parameter for paginated requests to set the number of results per page.


Authentication
--------------


Register a New User
```````````````````

Example Request:

.. code-block:: console

    POST /v1/auth/register/ HTTP/1.1
    Host: deis.example.com
    Content-Type: application/json

    {
        "username": "test",
        "password": "opensesame",
        "email": "test@example.com"
    }

Optional Parameters:

.. code-block:: console

    {
        "first_name": "test",
        "last_name": "testerson"
    }

Example Response:

.. code-block:: console

    HTTP/1.1 201 CREATED
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1
    Content-Type: application/json

    {
        "id": 1,
        "last_login": "2014-10-19T22:01:00.601Z",
        "is_superuser": true,
        "username": "test",
        "first_name": "test",
        "last_name": "testerson",
        "email": "test@example.com",
        "is_staff": true,
        "is_active": true,
        "date_joined": "2014-10-19T22:01:00.601Z",
        "groups": [],
        "user_permissions": []
    }


Log in
``````

Example Request:

.. code-block:: console

    POST /v1/auth/login/ HTTP/1.1
    Host: deis.example.com
    Content-Type: application/json

    {"username": "test", "password": "opensesame"}

Example Response:

.. code-block:: console

    HTTP/1.1 200 OK
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1
    Content-Type: application/json

    {"token": "abc123"}


Cancel Account
``````````````

Example Request:

.. code-block:: console

    DELETE /v1/auth/cancel/ HTTP/1.1
    Host: deis.example.com
    Authorization: token abc123

Example Response:

.. code-block:: console

    HTTP/1.1 204 NO CONTENT
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1

Regenerate Token
````````````````

.. note::

    This command could require administrative privileges

Example Request:

.. code-block:: console

    POST /v1/auth/tokens/ HTTP/1.1
    Host: deis.example.com
    Authorization: token abc123

Optional Parameters:

.. code-block:: console

    {
        "username" : "test"
        "all" : "true"
    }

Example Response:

.. code-block:: console

    HTTP/1.1 200 OK
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1
    Content-Type: application/json

    {"token": "abc123"}

Change Password
```````````````

Example Request:

.. code-block:: console

    POST /v1/auth/passwd/ HTTP/1.1
    Host: deis.example.com
    Authorization: token abc123

    {
        "password": "foo",
        "new_password": "bar"
    }

Optional parameters:

.. code-block:: console

    {"username": "testuser"}

.. note::

    Using the ``username`` parameter requires administrative privileges and
    makes the ``password`` parameter optional.

Example Response:

.. code-block:: console

    HTTP/1.1 200 OK
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1


Applications
------------


List all Applications
`````````````````````

Example Request:

.. code-block:: console

    GET /v1/apps HTTP/1.1
    Host: deis.example.com
    Authorization: token abc123

Example Response:

.. code-block:: console

    HTTP/1.1 200 OK
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1
    Content-Type: application/json

    {
        "count": 1,
        "next": null,
        "previous": null,
        "results": [
            {
                "created": "2014-01-01T00:00:00UTC",
                "id": "example-go",
                "owner": "test",
                "structure": {},
                "updated": "2014-01-01T00:00:00UTC",
                "url": "example-go.example.com",
                "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
            }
        ]
    }


Create an Application
`````````````````````

Example Request:

.. code-block:: console

    POST /v1/apps/ HTTP/1.1
    Host: deis.example.com
    Content-Type: application/json
    Authorization: token abc123

Optional parameters:

.. code-block:: console

    {"id": "example-go"}


Example Response:

.. code-block:: console

    HTTP/1.1 201 CREATED
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1
    Content-Type: application/json

    {
        "created": "2014-01-01T00:00:00UTC",
        "id": "example-go",
        "owner": "test",
        "structure": {},
        "updated": "2014-01-01T00:00:00UTC",
        "url": "example-go.example.com",
        "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
    }


Destroy an Application
``````````````````````

Example Request:

.. code-block:: console

    DELETE /v1/apps/example-go/ HTTP/1.1
    Host: deis.example.com
    Authorization: token abc123

Example Response:

.. code-block:: console

    HTTP/1.1 204 NO CONTENT
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1


List Application Details
````````````````````````

Example Request:

.. code-block:: console

    GET /v1/apps/example-go/ HTTP/1.1
    Host: deis.example.com
    Authorization: token abc123

Example Response:

.. code-block:: console

    HTTP/1.1 200 OK
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1
    Content-Type: application/json

    {
        "created": "2014-01-01T00:00:00UTC",
        "id": "example-go",
        "owner": "test",
        "structure": {},
        "updated": "2014-01-01T00:00:00UTC",
        "url": "example-go.example.com",
        "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
    }


Retrieve Application Logs
`````````````````````````

Example Request:

.. code-block:: console

    GET /v1/apps/example-go/logs/ HTTP/1.1
    Host: deis.example.com
    Authorization: token abc123

Optional URL Query Parameters:

.. code-block:: console

    ?log_lines=

Example Response:

.. code-block:: console

    HTTP/1.1 200 OK
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1
    Content-Type: text/plain

    "16:51:14 deis[api]: test created initial release\n"


Run one-off Commands
````````````````````

.. code-block:: console

    POST /v1/apps/example-go/run/ HTTP/1.1
    Host: deis.example.com
    Content-Type: application/json
    Authorization: token abc123

    {"command": "echo hi"}

Example Response:

.. code-block:: console

    HTTP/1.1 200 OK
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1
    Content-Type: application/json

    [0, "hi\n"]


Certificates
------------


List all Certificates
`````````````````````

Example Request:

.. code-block:: console

    GET /v1/certs HTTP/1.1
    Host: deis.example.com
    Authorization: token abc123

Example Response:

.. code-block:: console

    HTTP/1.1 200 OK
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1
    Content-Type: application/json

    {
        "count": 1,
        "next": null,
        "previous": null,
        "results": [
            {
                "common_name": "test.example.com",
                "expires": "2014-01-01T00:00:00UTC"
            }
        ]
    }


List Certificate Details
````````````````````````

Example Request:

.. code-block:: console

    GET /v1/certs/test.example.com HTTP/1.1
    Host: deis.example.com
    Authorization: token abc123

Example Response:

.. code-block:: console

    HTTP/1.1 200 OK
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1
    Content-Type: application/json

    {
        "updated": "2014-01-01T00:00:00UTC",
        "created": "2014-01-01T00:00:00UTC",
        "expires": "2015-01-01T00:00:00UTC",
        "common_name": "test.example.com",
        "owner": "test",
        "id": 1
    }


Create Certificate
``````````````````

Example Request:

.. code-block:: console

    POST /v1/certs/ HTTP/1.1
    Host: deis.example.com
    Content-Type: application/json
    Authorization: token abc123

    {
        "certificate": "-----BEGIN CERTIFICATE-----",
        "key": "-----BEGIN RSA PRIVATE KEY-----"
    }

Optional Parameters:

.. code-block:: console

    {
        "common_name": "test.example.com"
    }


Example Response:

.. code-block:: console

    HTTP/1.1 201 CREATED
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1
    Content-Type: application/json

    {
        "updated": "2014-01-01T00:00:00UTC",
        "created": "2014-01-01T00:00:00UTC",
        "expires": "2015-01-01T00:00:00UTC",
        "common_name": "test.example.com",
        "owner": "test",
        "id": 1
    }


Destroy a Certificate
`````````````````````

Example Request:

.. code-block:: console

    DELETE /v1/certs/test.example.com HTTP/1.1
    Host: deis.example.com
    Authorization: token abc123

Example Response:

.. code-block:: console

    HTTP/1.1 204 NO CONTENT
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1


Containers
----------


List all Containers
```````````````````

Example Request:

.. code-block:: console

    GET /v1/apps/example-go/containers/ HTTP/1.1
    Host: deis.example.com
    Authorization: token abc123

Example Response:

.. code-block:: console

    HTTP/1.1 200 OK
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1
    Content-Type: application/json

    {
        "count": 1,
        "next": null,
        "previous": null,
        "results": [
            {
                "owner": "test",
                "app": "example-go",
                "release": "v2",
                "created": "2014-01-01T00:00:00UTC",
                "updated": "2014-01-01T00:00:00UTC",
                "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
                "type": "web",
                "num": 1,
                "state": "up"
            }
        ]
    }


List all Containers by Type
```````````````````````````

Example Request:

.. code-block:: console

    GET /v1/apps/example-go/containers/web/ HTTP/1.1
    Host: deis.example.com
    Authorization: token abc123

Example Response:

.. code-block:: console

    HTTP/1.1 200 OK
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1
    Content-Type: application/json

    {
        "count": 1,
        "next": null,
        "previous": null,
        "results": [
            {
                "owner": "test",
                "app": "example-go",
                "release": "v2",
                "created": "2014-01-01T00:00:00UTC",
                "updated": "2014-01-01T00:00:00UTC",
                "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
                "type": "web",
                "num": 1,
                "state": "up"
            }
        ]
    }


Restart All Containers
``````````````````````

Example Request:

.. code-block:: console

    POST /v1/apps/example-go/containers/restart/ HTTP/1.1
    Host: deis.example.com
    Authorization: token abc123

Example Response:

.. code-block:: console

    HTTP/1.1 200 OK
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1
    Content-Type: application/json

    [
        {
            "owner": "test",
            "app": "example-go",
            "release": "v2",
            "created": "2014-01-01T00:00:00UTC",
            "updated": "2014-01-01T00:00:00UTC",
            "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
            "type": "web",
            "num": 1,
            "state": "up"
        }
    ]


Restart Containers by Type
``````````````````````````

Example Request:

.. code-block:: console

    POST /v1/apps/example-go/containers/web/restart/ HTTP/1.1
    Host: deis.example.com
    Authorization: token abc123

Example Response:

.. code-block:: console

    HTTP/1.1 200 OK
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1
    Content-Type: application/json

    [
        {
            "owner": "test",
            "app": "example-go",
            "release": "v2",
            "created": "2014-01-01T00:00:00UTC",
            "updated": "2014-01-01T00:00:00UTC",
            "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
            "type": "web",
            "num": 1,
            "state": "up"
        }
    ]


Restart Containers by Type and Number
`````````````````````````````````````

Example Request:

.. code-block:: console

    POST /v1/apps/example-go/containers/web/1/restart/ HTTP/1.1
    Host: deis.example.com
    Authorization: token abc123

Example Response:

.. code-block:: console

    HTTP/1.1 200 OK
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1
    Content-Type: application/json

    [
        {
            "owner": "test",
            "app": "example-go",
            "release": "v2",
            "created": "2014-01-01T00:00:00UTC",
            "updated": "2014-01-01T00:00:00UTC",
            "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
            "type": "web",
            "num": 1,
            "state": "up"
        }
    ]


Scale Containers
````````````````

Example Request:

.. code-block:: console

    POST /v1/apps/example-go/scale/ HTTP/1.1
    Host: deis.example.com
    Content-Type: application/json
    Authorization: token abc123

    {"web": 3}

Example Response:

.. code-block:: console

    HTTP/1.1 204 NO CONTENT
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1


Configuration
-------------


List Application Configuration
``````````````````````````````

Example Request:

.. code-block:: console

    GET /v1/apps/example-go/config/ HTTP/1.1
    Host: deis.example.com
    Authorization: token abc123

Example Response:

.. code-block:: console

    HTTP/1.1 200 OK
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1
    Content-Type: application/json

    {
        "owner": "test",
        "app": "example-go",
        "values": {
          "PLATFORM": "deis"
        },
        "memory": {},
        "cpu": {},
        "tags": {},
        "created": "2014-01-01T00:00:00UTC",
        "updated": "2014-01-01T00:00:00UTC",
        "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
    }


Create new Config
`````````````````

Example Request:

.. code-block:: console

    POST /v1/apps/example-go/config/ HTTP/1.1
    Host: deis.example.com
    Content-Type: application/json
    Authorization: token abc123

    {"values": {"HELLO": "world", "PLATFORM": "deis"}}

Example Response:

.. code-block:: console

    HTTP/1.1 201 CREATED
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1
    Content-Type: application/json
    X-Deis-Release: 3

    {
        "owner": "test",
        "app": "example-go",
        "values": {
            "DEIS_APP": "example-go",
            "DEIS_RELEASE": "v3",
            "HELLO": "world",
            "PLATFORM": "deis"

        },
        "memory": {},
        "cpu": {},
        "tags": {},
        "created": "2014-01-01T00:00:00UTC",
        "updated": "2014-01-01T00:00:00UTC",
        "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
    }


Unset Config Variable
`````````````````````

Example Request:

.. code-block:: console

    POST /v1/apps/example-go/config/ HTTP/1.1
    Host: deis.example.com
    Content-Type: application/json
    Authorization: token abc123

    {"values": {"HELLO": null}}

Example Response:

.. code-block:: console

    HTTP/1.1 201 CREATED
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1
    Content-Type: application/json
    X-Deis-Release: 4

    {
        "owner": "test",
        "app": "example-go",
        "values": {
            "DEIS_APP": "example-go",
            "DEIS_RELEASE": "v4",
            "PLATFORM": "deis"
       },
        "memory": {},
        "cpu": {},
        "tags": {},
        "created": "2014-01-01T00:00:00UTC",
        "updated": "2014-01-01T00:00:00UTC",
        "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
    }


Domains
-------


List Application Domains
````````````````````````

Example Request:

.. code-block:: console

    GET /v1/apps/example-go/domains/ HTTP/1.1
    Host: deis.example.com
    Authorization: token abc123

Example Response:

.. code-block:: console

    HTTP/1.1 200 OK
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1
    Content-Type: application/json

    {
        "count": 1,
        "next": null,
        "previous": null,
        "results": [
            {
                "app": "example-go",
                "created": "2014-01-01T00:00:00UTC",
                "domain": "example.example.com",
                "owner": "test",
                "updated": "2014-01-01T00:00:00UTC"
            }
        ]
    }


Add Domain
``````````

Example Request:

.. code-block:: console

    POST /v1/apps/example-go/domains/ HTTP/1.1
    Host: deis.example.com
    Authorization: token abc123

    {'domain': 'example.example.com'}

Example Response:

.. code-block:: console

    HTTP/1.1 201 CREATED
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1
    Content-Type: application/json

    {
        "app": "example-go",
        "created": "2014-01-01T00:00:00UTC",
        "domain": "example.example.com",
        "owner": "test",
        "updated": "2014-01-01T00:00:00UTC"
    }



Remove Domain
`````````````

Example Request:

.. code-block:: console

    DELETE /v1/apps/example-go/domains/example.example.com HTTP/1.1
    Host: deis.example.com
    Authorization: token abc123

Example Response:

.. code-block:: console

    HTTP/1.1 204 NO CONTENT
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1


Builds
------


List Application Builds
```````````````````````

Example Request:

.. code-block:: console

    GET /v1/apps/example-go/builds/ HTTP/1.1
    Host: deis.example.com
    Authorization: token abc123

Example Response:

.. code-block:: console

    HTTP/1.1 200 OK
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1
    Content-Type: application/json

    {
        "count": 1,
        "next": null,
        "previous": null,
        "results": [
            {
                "app": "example-go",
                "created": "2014-01-01T00:00:00UTC",
                "dockerfile": "FROM deis/slugrunner RUN mkdir -p /app WORKDIR /app ENTRYPOINT [\"/runner/init\"] ADD slug.tgz /app ENV GIT_SHA 060da68f654e75fac06dbedd1995d5f8ad9084db",
                "image": "example-go",
                "owner": "test",
                "procfile": {
                    "web": "example-go"
                },
                "sha": "060da68f",
                "updated": "2014-01-01T00:00:00UTC",
                "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
            }
        ]
    }


Create Application Build
````````````````````````

Example Request:

.. code-block:: console

    POST /v1/apps/example-go/builds/ HTTP/1.1
    Host: deis.example.com
    Content-Type: application/json
    Authorization: token abc123

    {"image": "deis/example-go:latest"}

Optional Parameters:

.. code-block:: console

    {
        "procfile": {
          "web": "./cmd"
        }
    }

Example Response:

.. code-block:: console

    HTTP/1.1 201 CREATED
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1
    Content-Type: application/json
    X-Deis-Release: 4

    {
        "app": "example-go",
        "created": "2014-01-01T00:00:00UTC",
        "dockerfile": "",
        "image": "deis/example-go:latest",
        "owner": "test",
        "procfile": {},
        "sha": "",
        "updated": "2014-01-01T00:00:00UTC",
        "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
    }


Releases
--------


List Application Releases
`````````````````````````

Example Request:

.. code-block:: console

    GET /v1/apps/example-go/releases/ HTTP/1.1
    Host: deis.example.com
    Authorization: token abc123

Example Response:

.. code-block:: console

    HTTP/1.1 200 OK
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1
    Content-Type: application/json

    {
        "count": 3,
        "next": null,
        "previous": null,
        "results": [
            {
                "app": "example-go",
                "build": "202d8e4b-600e-4425-a85c-ffc7ea607f61",
                "config": "ed637ceb-5d32-44bd-9406-d326a777a513",
                "created": "2014-01-01T00:00:00UTC",
                "owner": "test",
                "summary": "test changed nothing",
                "updated": "2014-01-01T00:00:00UTC",
                "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
                "version": 3
            },
            {
                "app": "example-go",
                "build": "202d8e4b-600e-4425-a85c-ffc7ea607f61",
                "config": "95bd6dea-1685-4f78-a03d-fd7270b058d1",
                "created": "2014-01-01T00:00:00UTC",
                "owner": "test",
                "summary": "test deployed 060da68",
                "updated": "2014-01-01T00:00:00UTC",
                "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
                "version": 2
            },
            {
                "app": "example-go",
                "build": null,
                "config": "95bd6dea-1685-4f78-a03d-fd7270b058d1",
                "created": "2014-01-01T00:00:00UTC",
                "owner": "test",
                "summary": "test created initial release",
                "updated": "2014-01-01T00:00:00UTC",
                "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
                "version": 1
            }
        ]
    }


List Release Details
````````````````````

Example Request:

.. code-block:: console

    GET /v1/apps/example-go/releases/v1/ HTTP/1.1
    Host: deis.example.com
    Authorization: token abc123

Example Response:

.. code-block:: console

    HTTP/1.1 200 OK
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1
    Content-Type: application/json

    {
        "app": "example-go",
        "build": null,
        "config": "95bd6dea-1685-4f78-a03d-fd7270b058d1",
        "created": "2014-01-01T00:00:00UTC",
        "owner": "test",
        "summary": "test created initial release",
        "updated": "2014-01-01T00:00:00UTC",
        "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
        "version": 1
    }


Rollback Release
````````````````

Example Request:

.. code-block:: console

    POST /v1/apps/example-go/releases/rollback/ HTTP/1.1
    Host: deis.example.com
    Content-Type: application/json
    Authorization: token abc123

    {"version": 1}

Example Response:

.. code-block:: console

    HTTP/1.1 201 CREATED
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1
    Content-Type: application/json

    {"version": 5}


Keys
----


List Keys
`````````

Example Request:

.. code-block:: console

    GET /v1/keys/ HTTP/1.1
    Host: deis.example.com
    Authorization: token abc123

Example Response:

.. code-block:: console

    HTTP/1.1 201 CREATED
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1
    Content-Type: application/json

    {
        "count": 1,
        "next": null,
        "previous": null,
        "results": [
            {
                "created": "2014-01-01T00:00:00UTC",
                "id": "test@example.com",
                "owner": "test",
                "public": "ssh-rsa <...>",
                "updated": "2014-01-01T00:00:00UTC",
                "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
            }
        ]
    }


Add Key to User
```````````````

Example Request:

.. code-block:: console

    POST /v1/keys/ HTTP/1.1
    Host: deis.example.com
    Authorization: token abc123

    {
        "id": "example",
        "public": "ssh-rsa <...>"
    }

Example Response:

.. code-block:: console

    HTTP/1.1 201 CREATED
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1
    Content-Type: application/json

    {
        "created": "2014-01-01T00:00:00UTC",
        "id": "example",
        "owner": "example",
        "public": "ssh-rsa <...>",
        "updated": "2014-01-01T00:00:00UTC",
        "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
    }


Remove Key from User
````````````````````

Example Request:

.. code-block:: console

    DELETE /v1/keys/example HTTP/1.1
    Host: deis.example.com
    Authorization: token abc123

Example Response:

.. code-block:: console

    HTTP/1.1 204 NO CONTENT
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1


Permissions
-----------


List Application Permissions
````````````````````````````

.. note::

    This does not include the app owner.

Example Request:

.. code-block:: console

    GET /v1/apps/example-go/perms/ HTTP/1.1
    Host: deis.example.com
    Authorization: token abc123

Example Response:

.. code-block:: console

    HTTP/1.1 200 OK
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1
    Content-Type: application/json

    {
        "users": [
            "test",
            "foo"
        ]
    }


Create Application Permission
`````````````````````````````

Example Request:

.. code-block:: console

    POST /v1/apps/example-go/perms/ HTTP/1.1
    Host: deis.example.com
    Authorization: token abc123

    {"username": "example"}

Example Response:

.. code-block:: console

    HTTP/1.1 201 CREATED
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1


Remove Application Permission
`````````````````````````````

Example Request:

.. code-block:: console

    DELETE /v1/apps/example-go/perms/example HTTP/1.1
    Host: deis.example.com
    Authorization: token abc123

Example Response:

.. code-block:: console

    HTTP/1.1 204 NO CONTENT
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1

List Administrators
```````````````````

Example Request:

.. code-block:: console

    GET /v1/admin/perms/ HTTP/1.1
    Host: deis.example.com
    Authorization: token abc123

Example Response:

.. code-block:: console

    HTTP/1.1 200 OK
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1
    Content-Type: application/json

    {
        "count": 2,
        "next": null
        "previous": null,
        "results": [
            {
                "username": "test",
                "is_superuser": true
            },
            {
                "username": "foo",
                "is_superuser": true
            }
        ]
    }


Grant User Administrative Privileges
````````````````````````````````````

.. note::

    This command requires administrative privileges

Example Request:

.. code-block:: console

    POST /v1/admin/perms HTTP/1.1
    Host: deis.example.com
    Authorization: token abc123

    {"username": "example"}

Example Response:

.. code-block:: console

    HTTP/1.1 201 CREATED
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1

Remove User's Administrative Privileges
```````````````````````````````````````

.. note::

    This command requires administrative privileges

Example Request:

.. code-block:: console

    DELETE /v1/admin/perms/example HTTP/1.1
    Host: deis.example.com
    Authorization: token abc123

Example Response:

.. code-block:: console

    HTTP/1.1 204 NO CONTENT
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1

Users
-----

List all users
``````````````

.. note::

    This command requires administrative privileges

Example Request:

.. code-block:: console

    GET /v1/users HTTP/1.1
    Host: deis.example.com
    Authorization: token abc123

Example Response:

.. code-block:: console

    HTTP/1.1 200 OK
    DEIS_API_VERSION: 1.6
    DEIS_PLATFORM_VERSION: 1.10.1
    Content-Type: application/json

    {
        "count": 1,
        "next": null,
        "previous": null,
        "results": [
            {
                "id": 1,
                "last_login": "2014-10-19T22:01:00.601Z",
                "is_superuser": true,
                "username": "test",
                "first_name": "test",
                "last_name": "testerson",
                "email": "test@example.com",
                "is_staff": true,
                "is_active": true,
                "date_joined": "2014-10-19T22:01:00.601Z",
                "groups": [],
                "user_permissions": []
            }
        ]
    }
