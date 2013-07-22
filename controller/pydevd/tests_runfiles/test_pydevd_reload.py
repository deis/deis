'''
Changed the doctest frome the xreload to actual unittest.
'''

import sys
import os.path

import sys
IS_JYTHON = sys.platform.find('java') != -1

sys.path.append(os.path.split(os.path.split(__file__)[0])[0])

if sys.version_info[0] == 2 and sys.version_info[1] <= 4:
    SAMPLE_CODE = """
class C:
    attr = 42
    def foo(self):
        return 42
    
    def bar(cls):
        return 42, 42
    
    def stomp():
        return 42, 42, 42
"""
else:
    SAMPLE_CODE = """
class C:
    attr = 42
    def foo(self):
        return 42
    @classmethod
    def bar(cls):
        return 42, 42
    @staticmethod
    def stomp():
        return 42, 42, 42
"""

import shutil
from pydevd_reload import xreload
import tempfile

tempdir = None
save_path = None
import unittest

class Test(unittest.TestCase):
    

    def setUp(self, nused=None):
        global tempdir, save_path
        tempdir = tempfile.mktemp()
        print(tempdir)
        os.makedirs(tempdir)
        save_path = list(sys.path)
        sys.path.append(tempdir)
    
    
    def tearDown(self, unused=None):
        global tempdir, save_path
        if save_path is not None:
            sys.path = save_path
            save_path = None
        if tempdir is not None:
            shutil.rmtree(tempdir)
            tempdir = None
            
    
    def make_mod(self, name="x", repl=None, subst=None):
        assert tempdir
        fn = os.path.join(tempdir, name + ".py")
        f = open(fn, "w")
        sample = SAMPLE_CODE
        if repl is not None and subst is not None:
            sample = sample.replace(repl, subst)
        try:
            f.write(sample)
        finally:
            f.close()

        
    def testMet1(self):
        self.make_mod()
        import x #@UnresolvedImport -- this is the module we created at runtime.
        from x import C as Foo #@UnresolvedImport
        C = x.C
        Cfoo = C.foo
        Cbar = C.bar
        Cstomp = C.stomp
        b = C()
        bfoo = b.foo
        in_list = [C]
        self.assertEqual(b.foo(), 42)
        self.assertEqual(bfoo(), 42)
        self.assertEqual(Cfoo(b), 42)
        self.assertEqual(Cbar(), (42, 42))
        self.assertEqual(Cstomp(), (42, 42, 42))
        self.assertEqual(in_list[0].attr, 42)
        self.assertEqual(Foo.attr, 42)
        self.make_mod(repl="42", subst="24")
        xreload(x)
        self.assertEqual(b.foo(), 24)
        self.assertEqual(bfoo(), 24)
        self.assertEqual(Cfoo(b), 24)
        self.assertEqual(Cbar(), (24, 24))
        self.assertEqual(Cstomp(), (24, 24, 24))
        self.assertEqual(in_list[0].attr, 24)
        self.assertEqual(Foo.attr, 24)
        
        
#=======================================================================================================================
# main
#=======================================================================================================================
if __name__ == '__main__':
    #this is so that we can run it frem the jython tests -- because we don't actually have an __main__ module
    #(so, it won't try importing the __main__ module)
    if not IS_JYTHON: #Doesn't really work in Jython 
        unittest.TextTestRunner().run(unittest.makeSuite(Test))
