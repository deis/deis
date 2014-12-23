"""
Unit tests for the Deis api app.

Run the tests with "./manage.py test api"
"""

from __future__ import unicode_literals

import json

from django.contrib.auth.models import User
from django.test import TestCase
from rest_framework.authtoken.models import Token


CERT = """
-----BEGIN CERTIFICATE-----
MIICOzCCAaQCCQDSGuK9K1QGAjANBgkqhkiG9w0BAQsFADBiMQswCQYDVQQGEwJV
UzETMBEGA1UECAwKQ2FsaWZvcm5pYTEUMBIGA1UEBwwLTG9zIEFuZ2VsZXMxEDAO
BgNVBAoMB0RlaXMuaW8xFjAUBgNVBAMMDXRlc3QuY2VydC5jb20wHhcNMTQxMjA1
MDAwMTQ5WhcNMTUxMjA1MDAwMTQ5WjBiMQswCQYDVQQGEwJVUzETMBEGA1UECAwK
Q2FsaWZvcm5pYTEUMBIGA1UEBwwLTG9zIEFuZ2VsZXMxEDAOBgNVBAoMB0RlaXMu
aW8xFjAUBgNVBAMMDXRlc3QuY2VydC5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0A
MIGJAoGBALJwp3PEI4c0sV8h/J6v99iIdtw62ZdhDPSFUNd5mZ4l+6jFQ8M8HND4
kmbtsnWIoaX9Pry4wK0WzCIUqP+iAIFQAG4X5bTdWaPcmAz+5lBCxJRfYiOJpbZr
vk+Gdl6JJ/emyXRHI3MfaRV1Wdwcjp4eZ9RCzpLw/Dnkf+nOjEKNAgMBAAEwDQYJ
KoZIhvcNAQELBQADgYEAOeUDV1JNug8RW+l9tzSpM/cZ43QNJyUNW8aIDFmxSK4H
UMKZn5TPoi8JM6BC4G9CUwEbAcykWIKgs4x9Dl4ZnvEMx7C4ZMo8oOHUME16NhQQ
kckSxuWXfSJ9kHqT9ZPYMz9HklUGPwVlmRctgRMz0CJmsxg39mgGoEEa2Gk4g8I=
-----END CERTIFICATE-----
"""

CERT2 = """
-----BEGIN CERTIFICATE-----
MIICPTCCAaYCCQCho6FTN4hADjANBgkqhkiG9w0BAQsFADBjMQswCQYDVQQGEwJV
UzETMBEGA1UECAwKQ2FsaWZvcm5pYTEUMBIGA1UEBwwLTG9zIEFuZ2VsZXMxEDAO
BgNVBAoMB0RlaXMuaW8xFzAVBgNVBAMMDnRlc3QyLmNlcnQuY29tMB4XDTE0MTIw
OTA3NDAyOVoXDTE1MTIwOTA3NDAyOVowYzELMAkGA1UEBhMCVVMxEzARBgNVBAgM
CkNhbGlmb3JuaWExFDASBgNVBAcMC0xvcyBBbmdlbGVzMRAwDgYDVQQKDAdEZWlz
LmlvMRcwFQYDVQQDDA50ZXN0Mi5jZXJ0LmNvbTCBnzANBgkqhkiG9w0BAQEFAAOB
jQAwgYkCgYEAsnCnc8QjhzSxXyH8nq/32Ih23DrZl2EM9IVQ13mZniX7qMVDwzwc
0PiSZu2ydYihpf0+vLjArRbMIhSo/6IAgVAAbhfltN1Zo9yYDP7mUELElF9iI4ml
tmu+T4Z2Xokn96bJdEcjcx9pFXVZ3ByOnh5n1ELOkvD8OeR/6c6MQo0CAwEAATAN
BgkqhkiG9w0BAQsFAAOBgQBsBZupNw0zLn+YPLb10tfQXlXWJnRVTcLS4ILG/4F7
sV75TbeiP84FP4hN2hq2b1bOHVTexLtJFNMs6tsMOcVnMJoaE4w1m+FqSiXo+nYd
jhck0H5PxTdE+oODSGe7rMItoCxT2nYthyGeO6vYui7Qn4upS3HSRL7qyPnwL5oG
RQ==
-----END CERTIFICATE-----
"""

KEY = """
-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQCycKdzxCOHNLFfIfyer/fYiHbcOtmXYQz0hVDXeZmeJfuoxUPD
PBzQ+JJm7bJ1iKGl/T68uMCtFswiFKj/ogCBUABuF+W03Vmj3JgM/uZQQsSUX2Ij
iaW2a75PhnZeiSf3psl0RyNzH2kVdVncHI6eHmfUQs6S8Pw55H/pzoxCjQIDAQAB
AoGANQXInFvB+uEre4tL15OOYCdcumA6XAMYqGgc94pInXfH6gSD+DWaknXqeu9S
wh4RepNf2xBDIKvPiKj+9scawtLrh4yksDXezb5c+rIVktW+dsiMVR59HAIpF7KX
nvA0w6FDeTz2xz6cYEFZJNHVqmNEEEnik7lHwcvVMv6eIwECQQDgE1t+tFfLYCaP
G6D69aYnCebDZBoEL7aikt1o3pdPYjaGP00lp/djIXDrlQEXtmp9PUqSPyQs4urC
ZTnHoUYhAkEAy9zYOqVFeruUAM+D7TiByUMY/yYpj1E1+2ytP86aI1mKd39Qbc5n
ZNbkvtv9ZTJF+AA9oRAv562ULhR/jEeW7QJBAKutqRg2zF1B2ckjff9JXnfimi9x
7ozukZuVspW6lWt48BWDQnRrcJs+7+lPTHsChCxYXV4Xinvpj7xJGi/dXIECQGxk
ylu0UJMHdZRQwha5uthmYr4XbnWTep5qlFue4Hn3PBZ5jSw1WOhXEl0g30SVTHqm
th4TW0VWF7nAkGjoD6kCQBSQdxtrRyQKtFd1SvjDsuNnPZZrBsM61Bd4Y0ppyn5C
era8PE+kBg7keazqwoKOFVY/1FMBrur6g2FBh7FwF/o=
-----END RSA PRIVATE KEY-----
"""

CERT_CN = 'test.cert.com'
CERT2_CN = 'test2.cert.com'


class CertTest(TestCase):

    """Tests application certificate management"""

    fixtures = ['tests.json']

    def setUp(self):
        # app
        self.user = User.objects.get(username='autotest')
        self.token = Token.objects.get(user=self.user).key
        url = '/v1/apps'
        response = self.client.post(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        self.app_id = response.data['id']  # noqa
        # app2
        self.user2 = User.objects.get(username='autotest2')
        self.token2 = Token.objects.get(user=self.user2).key
        url = '/v1/apps'
        response = self.client.post(url, HTTP_AUTHORIZATION='token {}'.format(self.token2))
        self.assertEqual(response.status_code, 201)
        self.app_id2 = response.data['id']  # noqa

    def test_manage_cert(self):
        """
        Test all CRUD operations
        """
        # update non-existing
        url = '/v1/apps/{}/certs/{}'.format(self.app_id, CERT_CN)
        body = {'cert': CERT, 'key': KEY}
        response = self.client.put(url, json.dumps(body), content_type='application/json',
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 404)
        # create bad cert
        url = '/v1/apps/{}/certs'.format(self.app_id)
        body = {'cert': 'NOT_A_VALID_CERT', 'key': KEY}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 400)
        # create
        url = '/v1/apps/{}/certs'.format(self.app_id)
        body = {'cert': CERT, 'key': KEY}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        # create already exists
        self.assertEqual(response.data['common_name'], CERT_CN)
        url = '/v1/apps/{}/certs'.format(self.app_id)
        body = {'cert': CERT, 'key': KEY}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 400)
        # list
        url = '/v1/apps/{}/certs'.format(self.app_id)
        response = self.client.get(url,
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['results'][0]['common_name'], CERT_CN)
        # update
        url = '/v1/apps/{}/certs/{}'.format(self.app_id, CERT_CN)
        body = {'cert': CERT, 'key': KEY}
        response = self.client.put(url, json.dumps(body), content_type='application/json',
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        # update wrong cn
        url = '/v1/apps/{}/certs/{}'.format(self.app_id, CERT_CN)
        body = {'cert': CERT2, 'key': KEY}
        response = self.client.put(url, json.dumps(body), content_type='application/json',
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 400)
        # delete
        url = '/v1/apps/{}/certs/{}'.format(self.app_id, CERT_CN)
        response = self.client.delete(url,
                                      HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 204)

    def test_manage_cert_invalid_app(self):
        """
        Attempting to manage certs for non-existent apps should 404
        """
        # create
        url = '/v1/apps/{}/certs'.format("poop")
        body = {'cert': CERT, 'key': KEY}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 404)
        # list
        url = '/v1/apps/{}/certs'.format("poop")
        response = self.client.get(url,
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 404)
        # update
        url = '/v1/apps/{}/certs/{}'.format(self.app_id, CERT_CN)
        body = {'cert': CERT, 'key': KEY}
        response = self.client.put(url, json.dumps(body), content_type='application/json',
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 404)
        # delete
        url = '/v1/apps/{}/certs/{}'.format("poop", "test.cert.com")
        response = self.client.delete(url,
                                      HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 404)

    def test_non_admin_owner_can_mange_certs(self):
        """
        If a non-admin can manage certs if they are owner.
        """
        # create
        url = '/v1/apps/{}/certs'.format(self.app_id2)
        body = {'cert': CERT, 'key': KEY}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token2))
        self.assertEqual(response.status_code, 201)
        # list
        url = '/v1/apps/{}/certs'.format(self.app_id2)
        response = self.client.get(url,
                                   HTTP_AUTHORIZATION='token {}'.format(self.token2))
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['results'][0]['common_name'], CERT_CN)
        # update
        url = '/v1/apps/{}/certs/{}'.format(self.app_id2, CERT_CN)
        body = {'cert': CERT, 'key': KEY}
        response = self.client.put(url, json.dumps(body), content_type='application/json',
                                   HTTP_AUTHORIZATION='token {}'.format(self.token2))
        self.assertEqual(response.status_code, 200)
        # delete
        url = '/v1/apps/{}/certs/{}'.format(self.app_id2, CERT_CN)
        response = self.client.delete(url,
                                      HTTP_AUTHORIZATION='token {}'.format(self.token2))
        self.assertEqual(response.status_code, 204)

    def test_non_admin_non_owner_can_not_mange_certs(self):
        """
        If a non-admin tries to manage cert for an app they don't own, deny.
        """
        # create
        url = '/v1/apps/{}/certs'.format(self.app_id)
        body = {'cert': CERT, 'key': KEY}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token2))
        self.assertEqual(response.status_code, 403)
        # list
        url = '/v1/apps/{}/certs'.format(self.app_id)
        response = self.client.get(url,
                                   HTTP_AUTHORIZATION='token {}'.format(self.token2))
        self.assertEqual(response.status_code, 403)
        # update
        url = '/v1/apps/{}/certs/{}'.format(self.app_id, CERT_CN)
        body = {'cert': CERT, 'key': KEY}
        response = self.client.put(url, json.dumps(body), content_type='application/json',
                                   HTTP_AUTHORIZATION='token {}'.format(self.token2))
        self.assertEqual(response.status_code, 403)
        # delete
        url = '/v1/apps/{}/certs/{}'.format(self.app_id, CERT_CN)
        response = self.client.delete(url,
                                      HTTP_AUTHORIZATION='token {}'.format(self.token2))
        self.assertEqual(response.status_code, 403)

    def test_admin_can_mange_certs_for_other_apps(self):
        """
        If a non-admin user creates an app, an administrator should be able
        manage certs for it.
        """
        # create
        url = '/v1/apps/{}/certs'.format(self.app_id2)
        body = {'cert': CERT, 'key': KEY}
        response = self.client.post(url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        # list
        url = '/v1/apps/{}/certs'.format(self.app_id2)
        response = self.client.get(url,
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        self.assertEqual(response.data['results'][0]['common_name'], CERT_CN)
        # update
        url = '/v1/apps/{}/certs/{}'.format(self.app_id2, CERT_CN)
        body = {'cert': CERT, 'key': KEY}
        response = self.client.put(url, json.dumps(body), content_type='application/json',
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 200)
        # delete
        url = '/v1/apps/{}/certs/{}'.format(self.app_id2, CERT_CN)
        response = self.client.delete(url,
                                      HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 204)
