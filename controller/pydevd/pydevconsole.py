try:
    from code import InteractiveConsole
except ImportError:
    from pydevconsole_code_for_ironpython import InteractiveConsole 

import os
import sys

try:
    False
    True
except NameError: # version < 2.3 -- didn't have the True/False builtins
    import __builtin__
    setattr(__builtin__, 'True', 1) #Python 3.0 does not accept __builtin__.True = 1 in its syntax
    setattr(__builtin__, 'False', 0)

from pydev_console_utils import BaseStdIn, StdIn, BaseInterpreterInterface

        
try:
    class ExecState:
        FIRST_CALL = True
        PYDEV_CONSOLE_RUN_IN_UI = False #Defines if we should run commands in the UI thread.
        
    from org.python.pydev.core.uiutils import RunInUiThread #@UnresolvedImport
    from java.lang import Runnable #@UnresolvedImport
    class Command(Runnable):
    
        def __init__(self, interpreter, line):
            self.interpreter = interpreter
            self.line = line
            
        def run(self):
            if ExecState.FIRST_CALL:
                ExecState.FIRST_CALL = False
                sys.stdout.write('\nYou are now in a console within Eclipse.\nUse it with care as it can halt the VM.\n')
                sys.stdout.write('Typing a line with "PYDEV_CONSOLE_TOGGLE_RUN_IN_UI"\nwill start executing all the commands in the UI thread.\n\n')
                
            if self.line == 'PYDEV_CONSOLE_TOGGLE_RUN_IN_UI':
                ExecState.PYDEV_CONSOLE_RUN_IN_UI = not ExecState.PYDEV_CONSOLE_RUN_IN_UI
                if ExecState.PYDEV_CONSOLE_RUN_IN_UI:
                    sys.stdout.write('Running commands in UI mode. WARNING: using sys.stdin (i.e.: calling raw_input()) WILL HALT ECLIPSE.\n')
                else:
                    sys.stdout.write('No longer running commands in UI mode.\n')
                self.more = False
            else:
                self.more = self.interpreter.push(self.line)
                
            
    def Sync(runnable):
        if ExecState.PYDEV_CONSOLE_RUN_IN_UI:
            return RunInUiThread.sync(runnable)
        else:
            return runnable.run()
        
except:
    #If things are not there, define a way in which there's no 'real' sync, only the default execution.
    class Command:
    
        def __init__(self, interpreter, line):
            self.interpreter = interpreter
            self.line = line
            
        def run(self):
            self.more = self.interpreter.push(self.line)
        
    def Sync(runnable):
        runnable.run()


try:
    try:
        execfile #Not in Py3k
    except NameError:
        from pydev_imports import execfile
        import builtins #@UnresolvedImport -- only Py3K
        builtins.execfile = execfile
        
except:
    pass


#=======================================================================================================================
# InterpreterInterface
#=======================================================================================================================
class InterpreterInterface(BaseInterpreterInterface):
    '''
        The methods in this class should be registered in the xml-rpc server.
    '''
    
    def __init__(self, host, client_port):
        self.client_port = client_port
        self.host = host
        try:
            import pydevd #@UnresolvedImport
            if pydevd.GetGlobalDebugger() is None:
                raise RuntimeError() #Work as if the debugger does not exist as it's not connected.
        except:
            self.namespace = globals()
        else:
            #Adapted from the code in pydevd
            #patch provided by: Scott Schlesier - when script is run, it does not 
            #pretend pydevconsole is not the main module, and
            #convince the file to be debugged that it was loaded as main
            sys.modules['pydevconsole'] = sys.modules['__main__']
            sys.modules['pydevconsole'].__name__ = 'pydevconsole'            
            
            from imp import new_module
            m = new_module('__main__')
            sys.modules['__main__'] = m
            ns = m.__dict__
            try:
                ns['__builtins__'] = __builtins__
            except NameError:
                pass #Not there on Jython...
            self.namespace = ns
        self.interpreter = InteractiveConsole(self.namespace)
        self._input_error_printed = False

        
    def doAddExec(self, line):
        command = Command(self.interpreter, line)
        Sync(command)
        return command.more


    def getNamespace(self):
        return self.namespace
    
        
    def getCompletions(self, text, act_tok):
        try:
            from _pydev_completer import Completer
            completer = Completer(self.namespace, None)
            return completer.complete(act_tok)
        except:
            import traceback;traceback.print_exc()
            return []
        
    
    def close(self):
        sys.exit(0)
        
    
try:
    from pydev_ipython_console import InterpreterInterface
except:
    sys.stderr.write('PyDev console: using default backend (IPython not available).\n')
    pass #IPython not available, proceed as usual.
    
#=======================================================================================================================
# _DoExit
#=======================================================================================================================
def _DoExit(*args):
    '''
        We have to override the exit because calling sys.exit will only actually exit the main thread,
        and as we're in a Xml-rpc server, that won't work.
    '''
    
    try:
        import java.lang.System
        java.lang.System.exit(1)
    except ImportError:
        if len(args) == 1:
            os._exit(args[0])
        else:
            os._exit(0)
    
    
#=======================================================================================================================
# StartServer
#=======================================================================================================================
def StartServer(host, port, client_port):
    #replace exit (see comments on method)
    #note that this does not work in jython!!! (sys method can't be replaced).
    sys.exit = _DoExit
    
    from pydev_imports import SimpleXMLRPCServer
    try:
        interpreter = InterpreterInterface(host, client_port)
        server = SimpleXMLRPCServer((host, port), logRequests=False)
    except:
        sys.stderr.write('Error starting server with host: %s, port: %s, client_port: %s\n' % (host, port, client_port))
        raise

    
    if True:
        #Functions for basic protocol
        server.register_function(interpreter.addExec)
        server.register_function(interpreter.getCompletions)
        server.register_function(interpreter.getDescription)
        server.register_function(interpreter.close)

        #Functions so that the console can work as a debugger (i.e.: variables view, expressions...)
        server.register_function(interpreter.connectToDebugger)
        server.register_function(interpreter.postCommand)
        server.register_function(interpreter.hello)

        server.serve_forever()
        
    else:
        #This is still not finished -- that's why the if True is there :)
        from pydev_imports import Queue
        queue_requests_received = Queue.Queue() #@UndefinedVariable
        queue_return_computed = Queue.Queue() #@UndefinedVariable
        
        def addExec(line):
            queue_requests_received.put(('addExec', line))
            return queue_return_computed.get(block=True)
        
        def getCompletions(text):
            queue_requests_received.put(('getCompletions', text))
            return queue_return_computed.get(block=True)
        
        def getDescription(text):
            queue_requests_received.put(('getDescription', text))
            return queue_return_computed.get(block=True)
        
        def close():
            queue_requests_received.put(('close', None))
            return queue_return_computed.get(block=True)
            
        server.register_function(addExec)
        server.register_function(getCompletions)
        server.register_function(getDescription)
        server.register_function(close)
        try:
            import PyQt4.QtGui #We can only start the PyQt4 loop if we actually have access to it.
        except ImportError:
            print('Unable to process gui events (PyQt4.QtGui not imported)')
            server.serve_forever()
        else:
            import threading
            class PydevHandleRequestsThread(threading.Thread):
                
                def run(self):
                    while 1:
                        #This is done on a thread (so, it may be blocking or not blocking, it doesn't matter)
                        #anyways, the request will be put on a queue and the return will be gotten from another
                        #one -- and those queues are shared with the main thread.
                        server.handle_request()
                    
            app = PyQt4.QtGui.QApplication([])
            def serve_forever():
                """Handle one request at a time until doomsday."""
                while 1:
                    try:
                        try:
                            func, param = queue_requests_received.get(block=True,timeout=1.0/20.0) #20 loops/second
                            attr = getattr(interpreter, func)
                            if param is not None:
                                queue_return_computed.put(attr(param))
                            else:
                                queue_return_computed.put(attr())
                        except Queue.Empty: #@UndefinedVariable
                            pass
                        
                        PyQt4.QtGui.qApp.processEvents()
                    except:
                        import traceback;traceback.print_exc()
                    
            PydevHandleRequestsThread().start()
            serve_forever()


    
#=======================================================================================================================
# main
#=======================================================================================================================
if __name__ == '__main__':
    sys.stdin = BaseStdIn()
    port, client_port = sys.argv[1:3]
    import pydev_localhost
    StartServer(pydev_localhost.get_localhost(), int(port), int(client_port))
    
