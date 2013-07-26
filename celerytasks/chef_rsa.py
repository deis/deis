"""
This file as copied from pyChef at https://github.com/coderanger/pychef
"""

from ctypes import CDLL
from ctypes import byref
from ctypes import c_char_p
from ctypes import c_int
from ctypes import c_long
from ctypes import c_size_t
from ctypes import c_ulong
from ctypes import c_void_p
from ctypes import create_string_buffer
from ctypes import string_at
import sys


if sys.platform in ('win32', 'cygwin'):
    _eay = CDLL('libeay32.dll')
elif sys.platform == 'darwin':
    _eay = CDLL('libcrypto.dylib')
else:
    _eay = CDLL('libcrypto.so')

#unsigned long ERR_get_error(void);
ERR_get_error = _eay.ERR_get_error
ERR_get_error.argtypes = []
ERR_get_error.restype = c_ulong

#void ERR_error_string_n(unsigned long e, char *buf, size_t len);
ERR_error_string_n = _eay.ERR_error_string_n
ERR_error_string_n.argtypes = [c_ulong, c_char_p, c_size_t]
ERR_error_string_n.restype = None


class SSLError(Exception):
    """An error in OpenSSL."""

    def __init__(self, message, *args):
        message = message % args
        err = ERR_get_error()
        if err:
            message += ':'
        while err:
            buf = create_string_buffer(120)
            ERR_error_string_n(err, buf, 120)
            message += '\n%s' % string_at(buf, 119)
            err = ERR_get_error()
        super(SSLError, self).__init__(message)


#BIO *   BIO_new(BIO_METHOD *type);
BIO_new = _eay.BIO_new
BIO_new.argtypes = [c_void_p]
BIO_new.restype = c_void_p

# BIO *BIO_new_mem_buf(void *buf, int len);
BIO_new_mem_buf = _eay.BIO_new_mem_buf
BIO_new_mem_buf.argtypes = [c_void_p, c_int]
BIO_new_mem_buf.restype = c_void_p

#BIO_METHOD *BIO_s_mem(void);
BIO_s_mem = _eay.BIO_s_mem
BIO_s_mem.argtypes = []
BIO_s_mem.restype = c_void_p

#long    BIO_ctrl(BIO *bp,int cmd,long larg,void *parg);
BIO_ctrl = _eay.BIO_ctrl
BIO_ctrl.argtypes = [c_void_p, c_int, c_long, c_void_p]
BIO_ctrl.restype = c_long

#define BIO_CTRL_RESET          1  /* opt - rewind/zero etc */
BIO_CTRL_RESET = 1
##define BIO_CTRL_INFO           3  /* opt - extra tit-bits */
BIO_CTRL_INFO = 3


#define BIO_reset(b)            (int)BIO_ctrl(b,BIO_CTRL_RESET,0,NULL)
def BIO_reset(b):
    return BIO_ctrl(b, BIO_CTRL_RESET, 0, None)


##define BIO_get_mem_data(b,pp)  BIO_ctrl(b,BIO_CTRL_INFO,0,(char *)pp)
def BIO_get_mem_data(b, pp):
    return BIO_ctrl(b, BIO_CTRL_INFO, 0, pp)


# int    BIO_free(BIO *a)
BIO_free = _eay.BIO_free
BIO_free.argtypes = [c_void_p]
BIO_free.restype = c_int


def BIO_free_errcheck(result, func, arguments):
    if result == 0:
        raise SSLError('Unable to free BIO')
BIO_free.errcheck = BIO_free_errcheck

#RSA *PEM_read_bio_RSAPrivateKey(BIO *bp, RSA **x,
#                                        pem_password_cb *cb, void *u);
PEM_read_bio_RSAPrivateKey = _eay.PEM_read_bio_RSAPrivateKey
PEM_read_bio_RSAPrivateKey.argtypes = [c_void_p, c_void_p, c_void_p, c_void_p]
PEM_read_bio_RSAPrivateKey.restype = c_void_p

#RSA *PEM_read_bio_RSAPublicKey(BIO *bp, RSA **x,
#                                        pem_password_cb *cb, void *u);
PEM_read_bio_RSAPublicKey = _eay.PEM_read_bio_RSAPublicKey
PEM_read_bio_RSAPublicKey.argtypes = [c_void_p, c_void_p, c_void_p, c_void_p]
PEM_read_bio_RSAPublicKey.restype = c_void_p

#int PEM_write_bio_RSAPrivateKey(BIO *bp, RSA *x, const EVP_CIPHER *enc,
#                                        unsigned char *kstr, int klen,
#                                        pem_password_cb *cb, void *u);
PEM_write_bio_RSAPrivateKey = _eay.PEM_write_bio_RSAPrivateKey
PEM_write_bio_RSAPrivateKey.argtypes = [
    c_void_p, c_void_p, c_void_p, c_char_p, c_int, c_void_p, c_void_p]
PEM_write_bio_RSAPrivateKey.restype = c_int

#int PEM_write_bio_RSAPublicKey(BIO *bp, RSA *x);
PEM_write_bio_RSAPublicKey = _eay.PEM_write_bio_RSAPublicKey
PEM_write_bio_RSAPublicKey.argtypes = [c_void_p, c_void_p]
PEM_write_bio_RSAPublicKey.restype = c_int

#int RSA_private_encrypt(int flen, unsigned char *from,
#    unsigned char *to, RSA *rsa,int padding);
RSA_private_encrypt = _eay.RSA_private_encrypt
RSA_private_encrypt.argtypes = [c_int, c_void_p, c_void_p, c_void_p, c_int]
RSA_private_encrypt.restype = c_int

#int RSA_public_decrypt(int flen, unsigned char *from,
#   unsigned char *to, RSA *rsa, int padding);
RSA_public_decrypt = _eay.RSA_public_decrypt
RSA_public_decrypt.argtypes = [c_int, c_void_p, c_void_p, c_void_p, c_int]
RSA_public_decrypt.restype = c_int

RSA_PKCS1_PADDING = 1
RSA_NO_PADDING = 3

# int RSA_size(const RSA *rsa);
RSA_size = _eay.RSA_size
RSA_size.argtypes = [c_void_p]
RSA_size.restype = c_int

#RSA *RSA_generate_key(int num, unsigned long e,
#    void (*callback)(int,int,void *), void *cb_arg);
RSA_generate_key = _eay.RSA_generate_key
RSA_generate_key.argtypes = [c_int, c_ulong, c_void_p, c_void_p]
RSA_generate_key.restype = c_void_p

##define RSA_F4  0x10001L
RSA_F4 = 0x10001

# void RSA_free(RSA *rsa);
RSA_free = _eay.RSA_free
RSA_free.argtypes = [c_void_p]


class Key(object):
    """An OpenSSL RSA key."""

    def __init__(self, fp=None):
        self.key = None
        self.public = False
        if not fp:
            return
        if isinstance(fp, basestring):
            if fp.startswith('-----'):
                # PEM formatted text
                self.raw = fp
            else:
                self.raw = open(fp, 'rb').read()
        else:
            self.raw = fp.read()
        self._load_key()

    def _load_key(self):
        if '\0' in self.raw:
            # Raw string has embedded nulls, treat it as binary data
            buf = create_string_buffer(self.raw, len(self.raw))
        else:
            buf = create_string_buffer(self.raw)

        bio = BIO_new_mem_buf(buf, len(buf))
        try:
            self.key = PEM_read_bio_RSAPrivateKey(bio, 0, 0, 0)
            if not self.key:
                BIO_reset(bio)
                self.public = True
                self.key = PEM_read_bio_RSAPublicKey(bio, 0, 0, 0)
            if not self.key:
                raise SSLError('Unable to load RSA key')
        finally:
            BIO_free(bio)

    @classmethod
    def generate(cls, size=1024, exp=RSA_F4):
        self = cls()
        self.key = RSA_generate_key(size, exp, None, None)
        return self

    def private_encrypt(self, value, padding=RSA_PKCS1_PADDING):
        if self.public:
            raise SSLError('private method cannot be used on a public key')
        buf = create_string_buffer(value, len(value))
        size = RSA_size(self.key)
        output = create_string_buffer(size)
        ret = RSA_private_encrypt(len(buf), buf, output, self.key, padding)
        if ret <= 0:
            raise SSLError('Unable to encrypt data')
        return output.raw[:ret]

    def public_decrypt(self, value, padding=RSA_PKCS1_PADDING):
        buf = create_string_buffer(value, len(value))
        size = RSA_size(self.key)
        output = create_string_buffer(size)
        ret = RSA_public_decrypt(len(buf), buf, output, self.key, padding)
        if ret <= 0:
            raise SSLError('Unable to decrypt data')
        return output.raw[:ret]

    def private_export(self):
        if self.public:
            raise SSLError('private method cannot be used on a public key')
        out = BIO_new(BIO_s_mem())
        PEM_write_bio_RSAPrivateKey(out, self.key, None, None, 0, None, None)
        buf = c_char_p()
        count = BIO_get_mem_data(out, byref(buf))
        pem = string_at(buf, count)
        BIO_free(out)
        return pem

    def public_export(self):
        out = BIO_new(BIO_s_mem())
        PEM_write_bio_RSAPublicKey(out, self.key)
        buf = c_char_p()
        count = BIO_get_mem_data(out, byref(buf))
        pem = string_at(buf, count)
        BIO_free(out)
        return pem

    def __del__(self):
        if self.key and RSA_free:
            RSA_free(self.key)
