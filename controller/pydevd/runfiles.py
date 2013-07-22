
def main():
    import sys
    
    #Separate the nose params and the pydev params.
    pydev_params = []
    other_test_framework_params = []
    found_other_test_framework_param = None
    
    NOSE_PARAMS = '--nose-params' 
    PY_TEST_PARAMS = '--py-test-params'
    
    for arg in sys.argv[1:]:
        if not found_other_test_framework_param and arg != NOSE_PARAMS and arg != PY_TEST_PARAMS:
            pydev_params.append(arg)
            
        else:
            if not found_other_test_framework_param:
                found_other_test_framework_param = arg
            else:
                other_test_framework_params.append(arg)
                
    
    #Here we'll run either with nose or with the pydev_runfiles.
    import pydev_runfiles
    import pydev_runfiles_xml_rpc
    import pydevd_constants
    from pydevd_file_utils import _NormFile
    
    DEBUG = 0
    if DEBUG:
        sys.stdout.write('Received parameters: %s\n' % (sys.argv,))
        sys.stdout.write('Params for pydev: %s\n' % (pydev_params,))
        if found_other_test_framework_param:
            sys.stdout.write('Params for test framework: %s, %s\n' % (found_other_test_framework_param, other_test_framework_params))
        
    try:
        configuration = pydev_runfiles.parse_cmdline([sys.argv[0]] + pydev_params)
    except:
        sys.stderr.write('Command line received: %s\n' % (sys.argv,))
        raise
    pydev_runfiles_xml_rpc.InitializeServer(configuration.port) #Note that if the port is None, a Null server will be initialized.

    NOSE_FRAMEWORK = 1
    PY_TEST_FRAMEWORK = 2
    try:
        if found_other_test_framework_param:
            test_framework = 0 #Default (pydev)
            if found_other_test_framework_param == NOSE_PARAMS:
                import nose
                test_framework = NOSE_FRAMEWORK
                
            elif found_other_test_framework_param == PY_TEST_PARAMS:
                import pytest
                test_framework = PY_TEST_FRAMEWORK
                
            else:
                raise ImportError()
                
        else:
            raise ImportError()
        
    except ImportError:
        if found_other_test_framework_param:
            sys.stderr.write('Warning: Could not import the test runner: %s. Running with the default pydev unittest runner instead.\n' % (
                found_other_test_framework_param,))
            
        test_framework = 0
        
    #Clear any exception that may be there so that clients don't see it.
    #See: https://sourceforge.net/tracker/?func=detail&aid=3408057&group_id=85796&atid=577329
    if hasattr(sys, 'exc_clear'):
        sys.exc_clear()
    
    if test_framework == 0:
        
        pydev_runfiles.main(configuration)
        
    else:
        #We'll convert the parameters to what nose or py.test expects.
        #The supported parameters are: 
        #runfiles.py  --config-file|-t|--tests <Test.test1,Test2>  dirs|files --nose-params xxx yyy zzz
        #(all after --nose-params should be passed directly to nose)
        
        #In java: 
        #--tests = Constants.ATTR_UNITTEST_TESTS
        #--config-file = Constants.ATTR_UNITTEST_CONFIGURATION_FILE
        
        
        #The only thing actually handled here are the tests that we want to run, which we'll 
        #handle and pass as what the test framework expects.

        py_test_accept_filter = {}
        files_to_tests = configuration.files_to_tests
            
        if files_to_tests:
            #Handling through the file contents (file where each line is a test)
            files_or_dirs = []
            for file, tests in files_to_tests.items():
                if test_framework == NOSE_FRAMEWORK:
                    for test in tests:
                        files_or_dirs.append(file+':'+test)
                        
                elif test_framework == PY_TEST_FRAMEWORK:
                    file = _NormFile(file)
                    py_test_accept_filter[file] = tests
                    files_or_dirs.append(file)
                
                else:
                    raise AssertionError('Cannot handle test framework: %s at this point.' % (test_framework,))
                        
        else:
            if configuration.tests:
                #Tests passed (works together with the files_or_dirs)
                files_or_dirs = []
                for file in configuration.files_or_dirs:
                    if test_framework == NOSE_FRAMEWORK:
                        for t in configuration.tests:
                            files_or_dirs.append(file+':'+t)
                            
                    elif test_framework == PY_TEST_FRAMEWORK:
                        file = _NormFile(file)
                        py_test_accept_filter[file] = configuration.tests
                        files_or_dirs.append(file)
                        
                    else:
                        raise AssertionError('Cannot handle test framework: %s at this point.' % (test_framework,))
            else:
                #Only files or dirs passed (let it do the test-loading based on those paths)
                files_or_dirs = configuration.files_or_dirs
                
        argv = other_test_framework_params + files_or_dirs
        

        if test_framework == NOSE_FRAMEWORK:
            #Nose usage: http://somethingaboutorange.com/mrl/projects/nose/0.11.2/usage.html
            #show_stdout_option = ['-s']
            #processes_option = ['--processes=2']
            argv.insert(0, sys.argv[0])
            if DEBUG:
                sys.stdout.write('Final test framework args: %s\n' % (argv[1:],))
                
            import pydev_runfiles_nose
            PYDEV_NOSE_PLUGIN_SINGLETON = pydev_runfiles_nose.StartPydevNosePluginSingleton(configuration)
            argv.append('--with-pydevplugin')
            nose.run(argv=argv, addplugins=[PYDEV_NOSE_PLUGIN_SINGLETON])

        elif test_framework == PY_TEST_FRAMEWORK:
            if DEBUG:
                sys.stdout.write('Final test framework args: %s\n' % (argv,))
                
            from pydev_runfiles_pytest import PydevPlugin
            pydev_plugin = PydevPlugin(py_test_accept_filter)
            
            pytest.main(argv, plugins=[pydev_plugin])
        
        else:
            raise AssertionError('Cannot handle test framework: %s at this point.' % (test_framework,))
        
        
if __name__ == '__main__':
    try:
        main()
    finally:
        try:
            #The server is not a daemon thread, so, we have to ask for it to be killed!
            import pydev_runfiles_xml_rpc
            pydev_runfiles_xml_rpc.forceServerKill()
        except:
            pass #Ignore any errors here
