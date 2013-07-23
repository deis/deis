import pydev_runfiles_xml_rpc
import time
from _pytest import runner #@UnresolvedImport
from _pytest import unittest as pytest_unittest #@UnresolvedImport
from py._code import code #@UnresolvedImport
from pydevd_file_utils import _NormFile
import os


#=======================================================================================================================
# _CollectTestsFromUnittestCase
#=======================================================================================================================
class _CollectTestsFromUnittestCase:
    
    def __init__(self, found_methods_starting, unittest_case):
        self.found_methods_starting = found_methods_starting
        self.unittest_case = unittest_case

    
    def __call__(self):
        for name in self.found_methods_starting:
            yield pytest_unittest.TestCaseFunction(name, parent=self.unittest_case)
            

#=======================================================================================================================
# PydevPlugin
#=======================================================================================================================
class PydevPlugin:
    
    def __init__(self, py_test_accept_filter):
        self.py_test_accept_filter = py_test_accept_filter
        self._original_pytest_collect_makeitem = pytest_unittest.pytest_pycollect_makeitem
        pytest_unittest.pytest_pycollect_makeitem = self.__pytest_pycollect_makeitem
        self._using_xdist = False
    
    def reportCond(self, cond, filename, test, captured_output, error_contents, delta):
        '''
        @param filename: 'D:\\src\\mod1\\hello.py'
        @param test: 'TestCase.testMet1'
        @param cond: fail, error, ok
        '''
        time_str = '%.2f' % (delta,)
        pydev_runfiles_xml_rpc.notifyTest(cond, captured_output, error_contents, filename, test, time_str)
        
        
    def __pytest_pycollect_makeitem(self, collector, name, obj):
        if not self.py_test_accept_filter:
            return self._original_pytest_collect_makeitem(collector, name, obj)
            
        f = _NormFile(collector.fspath.strpath)
        
        if f not in self.py_test_accept_filter:
            return
        
        tests = self.py_test_accept_filter[f]
        found_methods_starting = []
        for test in tests:
            
            if test == name:
                #Direct match of the test (just go on with the default loading)
                return self._original_pytest_collect_makeitem(collector, name, obj)

            
            if test.startswith(name+'.'):
                found_methods_starting.append(test[len(name)+1:])
        else:
            if not found_methods_starting:
                return
            
        #Ok, we found some method starting with the test name, let's gather those
        #and load them.
        unittest_case = self._original_pytest_collect_makeitem(collector, name, obj)
        if unittest_case is None:
            return

        unittest_case.collect = _CollectTestsFromUnittestCase(
            found_methods_starting, unittest_case)
        return unittest_case
        
        

    def _MockFileRepresentation(self):
        code.ReprFileLocation._original_toterminal = code.ReprFileLocation.toterminal
        
        def toterminal(self, tw):
            # filename and lineno output for each entry,
            # using an output format that most editors understand
            msg = self.message
            i = msg.find("\n")
            if i != -1:
                msg = msg[:i]

            tw.line('File "%s", line %s\n%s' %(os.path.abspath(self.path), self.lineno, msg))
            
        code.ReprFileLocation.toterminal = toterminal


    def _UninstallMockFileRepresentation(self):
        code.ReprFileLocation.toterminal = code.ReprFileLocation._original_toterminal #@UndefinedVariable


    def pytest_cmdline_main(self, config):
        if hasattr(config.option, 'numprocesses'):
            if config.option.numprocesses:
                self._using_xdist = True
                pydev_runfiles_xml_rpc.notifyTestRunFinished('Unable to show results (py.test xdist plugin not compatible with PyUnit view)')


    def pytest_runtestloop(self, session):
        if self._using_xdist:
            #Yes, we don't have the hooks we'd need to show the results in the pyunit view...
            #Maybe the plugin maintainer may be able to provide these additional hooks?
            return None
        
        #This mock will make all file representations to be printed as Pydev expects, 
        #so that hyperlinks are properly created in errors. Note that we don't unmock it!
        self._MockFileRepresentation()
        
        #Based on the default run test loop: _pytest.session.pytest_runtestloop
        #but getting the times we need, reporting the number of tests found and notifying as each
        #test is run.
        
        start_total = time.time()
        try:
            pydev_runfiles_xml_rpc.notifyTestsCollected(len(session.session.items))
            
            if session.config.option.collectonly:
                return True
            
            for item in session.session.items:
                
                filename = item.fspath.strpath
                test = item.location[2]
                start = time.time()
                
                pydev_runfiles_xml_rpc.notifyStartTest(filename, test)
                
                #Don't use this hook because we need the actual reports.
                #item.config.hook.pytest_runtest_protocol(item=item)
                reports = runner.runtestprotocol(item)
                delta = time.time() - start
                
                captured_output = ''
                error_contents = ''
                
                
                status = 'ok'
                for r in reports:
                    if r.outcome not in ('passed', 'skipped'):
                        #It has only passed, skipped and failed (no error), so, let's consider error if not on call.
                        if r.when == 'setup':
                            if status == 'ok':
                                status = 'error'
                            
                        elif r.when == 'teardown':
                            if status == 'ok':
                                status = 'error'
                            
                        else:
                            #any error in the call (not in setup or teardown) is considered a regular failure.
                            status = 'fail'
                        
                    if hasattr(r, 'longrepr') and r.longrepr:
                        rep = r.longrepr
                        if hasattr(rep, 'reprcrash'):
                            reprcrash = rep.reprcrash
                            error_contents += str(reprcrash)
                            error_contents += '\n'
                            
                        if hasattr(rep, 'reprtraceback'):
                            error_contents += str(rep.reprtraceback)
                            
                        if hasattr(rep, 'sections'):
                            for name, content, sep in rep.sections:
                                error_contents += sep * 40 
                                error_contents += name 
                                error_contents += sep * 40 
                                error_contents += '\n'
                                error_contents += content 
                                error_contents += '\n'
                
                self.reportCond(status, filename, test, captured_output, error_contents, delta)
                
                if session.shouldstop:
                    raise session.Interrupted(session.shouldstop)
        finally:
            pydev_runfiles_xml_rpc.notifyTestRunFinished('Finished in: %.2f secs.' % (time.time() - start_total,))
        return True
            
