from __future__ import unicode_literals

import json

from django.contrib.auth.models import User
from django.test import TestCase
from rest_framework.authtoken.models import Token

from api.models import App, Certificate


class CertificateTest(TestCase):

    """Tests creation of domain SSL certificates"""

    fixtures = ['tests.json']

    def setUp(self):
        self.user = User.objects.get(username='autotest')
        self.token = Token.objects.get(user=self.user).key
        self.user2 = User.objects.get(username='autotest2')
        self.token2 = Token.objects.get(user=self.user).key
        self.url = '/v1/certs'
        self.app = App.objects.create(owner=self.user, id='test-app')
        self.key = """-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAwyLIwjpUQkAmh/z6JvQMAtvNu/dBuCt+R8cnQMEw4VglglMw
YKAm2ZXA03LYWk5EO52YaDZKPAqjng+m4k+B0ble5XG4vFRTlBhln0cR3UAYlm7Z
tZp/6JR1STwph+9520DUsndPTOO3ApcMMuap5yLRYApHfOwbyoiaCCUuaE/XyZn6
FN/9Zj1V1IMcdu8//HtM2vbDkZ5yJbUzDSqUInHXZUp7OF+kUKwem0CN+SGk20ue
AQ3Loxg9RWcXwA8keZ1StQsNzmRzQTP/XGEvaTrJKgBkk9GHnxmC00L9q+zb0BaH
aZ5KXCKbwf0mCqOZngHuKTvKpD62TvPz46xE4wIDAQABAoIBABr5HO0UKP97ZJgZ
lO57f4mJnpej5vaxNGRxl/Bwg/QyPgUUwLQqjxQ2ig/waQ2akf33m9CT6JECG3nG
yhewS86UpBRtMs79jQwEj0+EAGkn6f4pVniu4Y1hsBCue0MqDBsNjBkbOt/y/iIi
hPIoRkYH3w86fIU9Ed5eIYSMtyx91wpGBwwpCh4ztfQ5jbBMZ0F5J+EnvzC41x2K
1o0bN6pr51epQBuyHz3SNAX0ce67f0jLhPSDl76nzsQsHem7rTPY4ZFTsRZE7lW0
lSA0S0z/sGpdoo1g0qvzg6T73/x8g0pdtf0N2ckbbafMvX1lba86Su9/KDRpS0RK
dymBkFkCgYEA6VQfKG2lZ1vEPq5JUQ8be1KbqzSEfvyqXd3Cb7iFcVVP3kNCRk6m
O04NJYUxDuF1LpWemGt5UCUUdLxcGFTYDW6gAKyfTuve87PPVvuHNsnJcJWW77aV
+yDhXgYUy9fCLMxtZwTwCCrqXEUtSgK5hvlwa8bYL/dE7YGhOa2ap/8CgYEA1hil
ezP8REe+Z+M8tSt2hoZsxrBuso2pZRAxuMqiO0/trA3d0w2M51vSm1/NxM2JpW2y
SPtE9CbngyGeHNdC/SvEkHOZxKacimoD2LUjAcVA+5r+shK+ssMqnniy9Qh13AGg
Pj3ba9j10T3zzAhItefpIu5E+swhqs1xmhTQwx0CgYBYVzY4y1K9kFv702cE3rBr
/7nal1a28ZjbUzPjsrwrTb6gi1yTXAHKIGIP257YYHpKefGDCeXzdyaIkCxaNf1b
EJBZ0QG8EsfmAyU0bKUkFEBFdQ2hksK0Qx2wyKKlDvqAlaGySIdMwFrdNn/QLrnp
pZVv6Og/OOKK/fJ58QXGJwKBgDOsmzRTZc3tKw3UEPEBXog1pceHChDalEoqUHXz
opiCQDFI34NzP9EPnpOV2gpoOZLOGTv4ObpcMYC6+ninlCmbCMR8wl5ugFYAJJGH
lr10qKyRymucjp6C8KRzKW5u7lN9qPmc4Hr1UM+CDnfuf+433VNrAwctgerBz2uL
HqAZAoGAYbrDiueIFxHDrkCkefSyAn4Wlo6KhPSUiSqvM9k5gBWZedcvJrjbvCmW
K1NefGc57cAb906Lwa3MpUmKEA5IYTGsO87iAFnDMcuu+w6RwiwV/DNY8xB6dtuz
r8G+so0UVAch6q1OBBSBaKC1Vn3fzT72zvS7/e5BZ0p5KrqCIZg=
-----END RSA PRIVATE KEY-----"""
        self.autotest_example_com_cert = """-----BEGIN CERTIFICATE-----
MIID3jCCAsYCCQDg75CmAL+avjANBgkqhkiG9w0BAQUFADCBsDELMAkGA1UEBhMC
Q0ExGTAXBgNVBAgTEEJyaXRpc2gtQ29sdW1iaWExEjAQBgNVBAcTCVZhbmNvdXZl
cjEtMCsGA1UEChMkRmlzaHdvcmtzIERldmVsb3BtZW50IGFuZCBDb25zdWx0aW5n
MR0wGwYDVQQDExRhdXRvdGVzdC5leGFtcGxlLmNvbTEkMCIGCSqGSIb3DQEJARYV
bWF0dGhld2ZAZmlzaHdvcmtzLmlvMB4XDTE1MDMwNjE3MTQyN1oXDTE2MDMwNTE3
MTQyN1owgbAxCzAJBgNVBAYTAkNBMRkwFwYDVQQIExBCcml0aXNoLUNvbHVtYmlh
MRIwEAYDVQQHEwlWYW5jb3V2ZXIxLTArBgNVBAoTJEZpc2h3b3JrcyBEZXZlbG9w
bWVudCBhbmQgQ29uc3VsdGluZzEdMBsGA1UEAxMUYXV0b3Rlc3QuZXhhbXBsZS5j
b20xJDAiBgkqhkiG9w0BCQEWFW1hdHRoZXdmQGZpc2h3b3Jrcy5pbzCCASIwDQYJ
KoZIhvcNAQEBBQADggEPADCCAQoCggEBAMMiyMI6VEJAJof8+ib0DALbzbv3Qbgr
fkfHJ0DBMOFYJYJTMGCgJtmVwNNy2FpORDudmGg2SjwKo54PpuJPgdG5XuVxuLxU
U5QYZZ9HEd1AGJZu2bWaf+iUdUk8KYfvedtA1LJ3T0zjtwKXDDLmqeci0WAKR3zs
G8qImgglLmhP18mZ+hTf/WY9VdSDHHbvP/x7TNr2w5GeciW1Mw0qlCJx12VKezhf
pFCsHptAjfkhpNtLngENy6MYPUVnF8APJHmdUrULDc5kc0Ez/1xhL2k6ySoAZJPR
h58ZgtNC/avs29AWh2meSlwim8H9JgqjmZ4B7ik7yqQ+tk7z8+OsROMCAwEAATAN
BgkqhkiG9w0BAQUFAAOCAQEAwYpXB8z4aOBedyHikbtVjDs1k0LEtWRAX/RXQY4I
BAYTnO+eGs/p7o+e3LGrIt/pX8kJ0RgD7TLITUJCZ69KkG9GzZaJ/CgQgqEa4Goh
JCI5u5a5nkTE6zZgAkkvpbA3Mj6WXGkGk7QEiO1e6e3y0jIBhDo1piD+DIppMWwM
OI0/r46FDlPHnm+y7UmTx+GZB4RAxnFaJE5L76w63oIPaRc/zkhS49AYiSmlawxj
thejiQz0ThCMBw7QMpVOiSvYAlQG0ATsRYwdTDqENIWKlerOLCSuxmbqe8XeDKhq
0ExzRJX9L9CjFIx9k+fIebIJWdv4Y4YUEtbLVmkKeghVJA==
-----END CERTIFICATE-----"""

    def test_create_certificate_with_domain(self):
        """Tests creating a certificate."""
        body = {'certificate': self.autotest_example_com_cert, 'key': self.key}
        response = self.client.post(self.url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)

    def test_create_wildcard_certificate(self):
        """Tests creating a wildcard certificate, which should be disabled."""
        body = {'certificate': self.autotest_example_com_cert,
                'key': self.key,
                'common_name': '*.example.com'}
        response = self.client.post(self.url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 400)
        self.assertEqual(json.loads(response.content),
                         {'common_name': ['Wildcard certificates are not supported']})

    def test_create_certificate_with_different_common_name(self):
        """
        In some cases such as with SAN certificates, the certificate can cover more
        than a single domain. In that case, we want to be able to specify the common
        name for the certificate/key.
        """
        body = {'certificate': self.autotest_example_com_cert,
                'key': self.key,
                'common_name': 'foo.example.com'}
        response = self.client.post(self.url, json.dumps(body), content_type='application/json',
                                    HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 201)
        self.assertEqual(response.data['common_name'], 'foo.example.com')

    def test_get_certificate_screens_data(self):
        """
        When a user retrieves a certificate, only the common name and expiry date should be
        displayed.
        """
        body = {'certificate': self.autotest_example_com_cert, 'key': self.key}
        self.client.post(self.url, json.dumps(body), content_type='application/json',
                         HTTP_AUTHORIZATION='token {}'.format(self.token))
        response = self.client.get('{}/{}'.format(self.url, 'autotest.example.com'),
                                   HTTP_AUTHORIZATION='token {}'.format(self.token))
        expected = {'common_name': 'autotest.example.com',
                    'expires': '2016-03-05T17:14:27UTC'}
        for key, value in expected.items():
            self.assertEqual(response.data[key], value)

    def test_certficate_denied_requests(self):
        """Disallow put/patch requests"""
        response = self.client.put(self.url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 405)
        response = self.client.patch(self.url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 405)

    def test_delete_certificate(self):
        """Destroying a certificate should generate a 204 response"""
        Certificate.objects.create(owner=self.user,
                                   common_name='autotest.example.com',
                                   certificate=self.autotest_example_com_cert)
        url = '/v1/certs/autotest.example.com'
        response = self.client.delete(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 204)
        # deleting a wildcard cert should work too (even though they're unsupported right now)
        # https://github.com/deis/deis/issues/3533
        Certificate.objects.create(owner=self.user,
                                   common_name='*.example.com',
                                   certificate=self.autotest_example_com_cert)
        url = '/v1/certs/*.example.com'
        response = self.client.delete(url, HTTP_AUTHORIZATION='token {}'.format(self.token))
        self.assertEqual(response.status_code, 204)
