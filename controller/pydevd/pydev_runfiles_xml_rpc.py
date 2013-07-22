from pydev_imports import xmlrpclib
from pydevd_constants import *
import traceback
import threading
try:
    from Queue import Queue
except:
    from queue import Queue

#This may happen in IronPython (in Python it shouldn't happen as there are
#'fast' replacements that are used in xmlrpclib.py)
import warnings
warnings.filterwarnings(
    'ignore', 'The xmllib module is obsolete.*', DeprecationWarning)



#=======================================================================================================================
# _ServerHolder
#=======================================================================================================================
class _ServerHolder:
    '''
    Helper so that we don't have to use a global here.
    '''
    SERVER = None


#=======================================================================================================================
# SetServer
#=======================================================================================================================
def SetServer(server):
    _ServerHolder.SERVER = server



#=======================================================================================================================
# ParallelNotification
#=======================================================================================================================
class ParallelNotification(object):
    
    def __init__(self, method, args):
        self.method = method
        self.args = args

    def ToTuple(self):
        return self.method, self.args
    
        
        
#=======================================================================================================================
# KillServer
#=======================================================================================================================
class KillServer(object):
    pass


#=======================================================================================================================
# ServerFacade
#=======================================================================================================================
class ServerFacade(object):
    
    
    def __init__(self, notifications_queue):
        self.notifications_queue = notifications_queue
    
    
    def notifyTestsCollected(self, *args):
        self.notifications_queue.put_nowait(ParallelNotification('notifyTestsCollected', args))
    
    def notifyConnected(self, *args):
        self.notifications_queue.put_nowait(ParallelNotification('notifyConnected', args))
    
    
    def notifyTestRunFinished(self, *args):
        self.notifications_queue.put_nowait(ParallelNotification('notifyTestRunFinished', args))
        
        
    def notifyStartTest(self, *args):
        self.notifications_queue.put_nowait(ParallelNotification('notifyStartTest', args))
        
        
    def notifyTest(self, *args):
        self.notifications_queue.put_nowait(ParallelNotification('notifyTest', args))





#=======================================================================================================================
# ServerComm
#=======================================================================================================================
class ServerComm(threading.Thread):
    

    
    def __init__(self, notifications_queue, port):
        threading.Thread.__init__(self)
        self.setDaemon(False) #Wait for all the notifications to be passed before exiting!
        self.finished = False
        self.notifications_queue = notifications_queue
        
        import pydev_localhost
        self.server = xmlrpclib.Server('http://%s:%s' % (pydev_localhost.get_localhost(), port))
        
    
    def run(self):
        while True:
            kill_found = False
            commands = []
            command = self.notifications_queue.get(block=True)
            if isinstance(command, KillServer):
                kill_found = True
            else:
                assert isinstance(command, ParallelNotification)
                commands.append(command.ToTuple())
                
            try:
                while True:
                    command = self.notifications_queue.get(block=False) #No block to create a batch.
                    if isinstance(command, KillServer):
                        kill_found = True
                    else:
                        assert isinstance(command, ParallelNotification)
                        commands.append(command.ToTuple())
            except:
                pass #That's OK, we're getting it until it becomes empty so that we notify multiple at once.


            if commands:
                try:
                    self.server.notifyCommands(commands)
                except:
                    traceback.print_exc()
            
            if kill_found:
                self.finished = True
                return



#=======================================================================================================================
# InitializeServer
#=======================================================================================================================
def InitializeServer(port):
    if _ServerHolder.SERVER is None:
        if port is not None:
            notifications_queue = Queue()
            _ServerHolder.SERVER = ServerFacade(notifications_queue)
            _ServerHolder.SERVER_COMM = ServerComm(notifications_queue, port)
            _ServerHolder.SERVER_COMM.start()
        else:
            #Create a null server, so that we keep the interface even without any connection.
            _ServerHolder.SERVER = Null()
            _ServerHolder.SERVER_COMM = Null()
        
    try:
        _ServerHolder.SERVER.notifyConnected()
    except:
        traceback.print_exc()

    
    
#=======================================================================================================================
# notifyTest
#=======================================================================================================================
def notifyTestsCollected(tests_count):
    assert tests_count is not None
    try:
        _ServerHolder.SERVER.notifyTestsCollected(tests_count)
    except:
        traceback.print_exc()
    
    
#=======================================================================================================================
# notifyStartTest
#=======================================================================================================================
def notifyStartTest(file, test):
    '''
    @param file: the tests file (c:/temp/test.py)
    @param test: the test ran (i.e.: TestCase.test1)
    '''
    assert file is not None
    if test is None:
        test = '' #Could happen if we have an import error importing module.
        
    try:
        _ServerHolder.SERVER.notifyStartTest(file, test)
    except:
        traceback.print_exc()

    
#=======================================================================================================================
# notifyTest
#=======================================================================================================================
def notifyTest(cond, captured_output, error_contents, file, test, time):
    '''
    @param cond: ok, fail, error
    @param captured_output: output captured from stdout
    @param captured_output: output captured from stderr
    @param file: the tests file (c:/temp/test.py)
    @param test: the test ran (i.e.: TestCase.test1)
    @param time: float with the number of seconds elapsed
    '''
    assert cond is not None
    assert captured_output is not None
    assert error_contents is not None
    assert file is not None
    if test is None:
        test = '' #Could happen if we have an import error importing module.
    assert time is not None
    try:
        _ServerHolder.SERVER.notifyTest(cond, captured_output, error_contents, file, test, time)
    except:
        traceback.print_exc()

#=======================================================================================================================
# notifyTestRunFinished
#=======================================================================================================================
def notifyTestRunFinished(total_time):
    assert total_time is not None
    try:
        _ServerHolder.SERVER.notifyTestRunFinished(total_time)
    except:
        traceback.print_exc()
    
    
#=======================================================================================================================
# forceServerKill
#=======================================================================================================================
def forceServerKill():
    _ServerHolder.SERVER_COMM.notifications_queue.put_nowait(KillServer())
