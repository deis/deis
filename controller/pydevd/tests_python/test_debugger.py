'''
    The idea is that we record the commands sent to the debugger and reproduce them from this script
    (so, this works as the client, which spawns the debugger as a separate process and communicates
    to it as if it was run from the outside)
    
    Note that it's a python script but it'll spawn a process to run as jython and as python.
'''
JYTHON_JAR_LOCATION = None
JAVA_LOCATION = None


import unittest 
port = 13336

def UpdatePort():
    global port
    port += 1

import os
def NormFile(filename):
    try:
        rPath = os.path.realpath #@UndefinedVariable
    except:
        # jython does not support os.path.realpath
        # realpath is a no-op on systems without islink support
        rPath = os.path.abspath   
    return os.path.normcase(rPath(filename))

PYDEVD_FILE = NormFile('../pydevd.py')
import sys
sys.path.append(os.path.dirname(PYDEVD_FILE))

SHOW_WRITES_AND_READS = False
SHOW_RESULT_STR = False
SHOW_OTHER_DEBUG_INFO = False


import subprocess
import socket
import threading
import time

#=======================================================================================================================
# ReaderThread
#=======================================================================================================================
class ReaderThread(threading.Thread):
    
    def __init__(self, sock):
        threading.Thread.__init__(self)
        self.setDaemon(True)
        self.sock = sock
        self.lastReceived = None
        
    def run(self):
        try:
            buf = ''
            while True:
                l = self.sock.recv(1024)
                buf += l
                
                if '\n' in buf:
                    self.lastReceived = buf
                    buf = ''
                    
                if SHOW_WRITES_AND_READS:
                    print 'Test Reader Thread Received %s' % self.lastReceived.strip()
        except:
            pass #ok, finished it
    
    def DoKill(self):
        self.sock.close()
    
#=======================================================================================================================
# AbstractWriterThread
#=======================================================================================================================
class AbstractWriterThread(threading.Thread):
    
    def __init__(self):
        threading.Thread.__init__(self)
        self.setDaemon(True)
        self.finishedOk = False
        
    def DoKill(self):
        if hasattr(self, 'readerThread'):
            #if it's not created, it's not there...
            self.readerThread.DoKill()
        self.sock.close()
        
    def Write(self, s):
        last = self.readerThread.lastReceived
        if SHOW_WRITES_AND_READS:
            print 'Test Writer Thread Written %s' % (s,)
        self.sock.send(s + '\n')
        time.sleep(0.2)
        
        i = 0
        while last == self.readerThread.lastReceived and i < 10:
            i += 1
            time.sleep(0.1)
        
    
    def StartSocket(self):
        if SHOW_WRITES_AND_READS:
            print 'StartSocket'
        
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.bind(('', port))
        s.listen(1)
        if SHOW_WRITES_AND_READS:
            print 'Waiting in socket.accept()'
        newSock, addr = s.accept()
        if SHOW_WRITES_AND_READS:
            print 'Test Writer Thread Socket:', newSock, addr
            
        readerThread = self.readerThread = ReaderThread(newSock)
        readerThread.start()
        self.sock = newSock
        
        self._sequence = -1
        #initial command is always the version
        self.WriteVersion()
    
    def NextSeq(self):
        self._sequence += 2
        return self._sequence
            
        
    def WaitForNewThread(self):
        i = 0
        #wait for hit breakpoint
        while not '<xml><thread name="' in self.readerThread.lastReceived or '<xml><thread name="pydevd.' in self.readerThread.lastReceived:
            i += 1
            time.sleep(1)
            if i >= 15:
                raise AssertionError('After %s seconds, a thread was not created.' % i)
        
        #we have something like <xml><thread name="MainThread" id="12103472" /></xml>
        splitted = self.readerThread.lastReceived.split('"')
        threadId = splitted[3]
        return threadId
        
    def WaitForBreakpointHit(self, reason='111', get_line=False): 
        '''
            108 is over
            109 is return
            111 is breakpoint
        '''
        i = 0
        #wait for hit breakpoint
        while not ('stop_reason="%s"' % reason) in self.readerThread.lastReceived:
            i += 1
            time.sleep(1)
            if i >= 10:
                raise AssertionError('After %s seconds, a break with reason: %s was not hit. Found: %s' % \
                    (i, reason, self.readerThread.lastReceived))
            
        #we have something like <xml><thread id="12152656" stop_reason="111"><frame id="12453120" ...
        splitted = self.readerThread.lastReceived.split('"')
        threadId = splitted[1]
        frameId = splitted[5]
        if get_line:
            return threadId, frameId, int(splitted[11])
            
        return threadId, frameId
        
    def WaitForVars(self, expected): 
        i = 0
        #wait for hit breakpoint
        while not expected in self.readerThread.lastReceived:
            i += 1
            time.sleep(1)
            if i >= 10:
                raise AssertionError('After %s seconds, the vars were not found. Last found:\n%s' % 
                    (i, self.readerThread.lastReceived))

        return True
        
    def WaitForMultipleVars(self, expected_vars): 
        i = 0
        #wait for hit breakpoint
        while True:
            for expected in expected_vars:
                if expected not in self.readerThread.lastReceived:
                    break #Break out of loop (and don't get to else)
            else:
                return True
            
            i += 1
            time.sleep(1)
            if i >= 10:
                raise AssertionError('After %s seconds, the vars were not found. Last found:\n%s' % 
                    (i, self.readerThread.lastReceived))

        return True
    
    def WriteMakeInitialRun(self):
        self.Write("101\t%s\t" % self.NextSeq())
        
    def WriteVersion(self):
        self.Write("501\t%s\t1.0" % self.NextSeq())
        
    def WriteAddBreakpoint(self, line, func):
        '''
            @param line: starts at 1
        '''
        if func is not None:
            self.Write("111\t%s\t%s\t%s\t**FUNC**%s\tNone" % (self.NextSeq(), self.TEST_FILE, line, func))
        else:
            self.Write("111\t%s\t%s\t%s\tNone" % (self.NextSeq(), self.TEST_FILE, line))
            
    def WriteRemoveBreakpoint(self, line):
        self.Write("112\t%s\t%s\t%s" % (self.NextSeq(), self.TEST_FILE, line))
        
    def WriteGetFrame(self, threadId, frameId):
        self.Write("114\t%s\t%s\t%s\tFRAME" % (self.NextSeq(), threadId, frameId))
        
    def WriteStepOver(self, threadId):
        self.Write("108\t%s\t%s" % (self.NextSeq(), threadId,))
        
    def WriteStepIn(self, threadId):
        self.Write("107\t%s\t%s" % (self.NextSeq(), threadId,))
        
    def WriteStepReturn(self, threadId):
        self.Write("109\t%s\t%s" % (self.NextSeq(), threadId,))

    def WriteSuspendThread(self, threadId):
        self.Write("105\t%s\t%s" % (self.NextSeq(), threadId,))
        
    def WriteRunThread(self, threadId):
        self.Write("106\t%s\t%s" % (self.NextSeq(), threadId,))
        
    def WriteKillThread(self, threadId):
        self.Write("104\t%s\t%s" % (self.NextSeq(), threadId,))

    def WriteDebugConsoleExpression(self, locator):
        self.Write("126\t%s\t%s"%(self.NextSeq(), locator))

#=======================================================================================================================
# WriterThreadCase14 - [Test Case]: Interactive Debug Console
#======================================================================================================================
class WriterThreadCase14(AbstractWriterThread):

    TEST_FILE = NormFile('_debugger_case14.py')

    def run(self):
        self.StartSocket()
        self.WriteAddBreakpoint(22, 'main')
        self.WriteMakeInitialRun()

        threadId, frameId, line = self.WaitForBreakpointHit('111', True)

        # Access some variable
        self.WriteDebugConsoleExpression("%s\t%s\tEVALUATE\tcarObj.color"%(threadId, frameId))
        self.WaitForMultipleVars(['<more>False</more>', '%27Black%27'])
        assert 7 == self._sequence, 'Expected 9. Had: %s' % self._sequence

        # Change some variable
        self.WriteDebugConsoleExpression("%s\t%s\tEVALUATE\tcarObj.color='Red'"%(threadId, frameId))
        self.WriteDebugConsoleExpression("%s\t%s\tEVALUATE\tcarObj.color"%(threadId, frameId))
        self.WaitForMultipleVars(['<more>False</more>', '%27Red%27'])
        assert 11 == self._sequence, 'Expected 13. Had: %s' % self._sequence

        # Iterate some loop
        self.WriteDebugConsoleExpression("%s\t%s\tEVALUATE\tfor i in range(3):"%(threadId, frameId))
        self.WaitForVars('<xml><more>True</more></xml>')
        self.WriteDebugConsoleExpression("%s\t%s\tEVALUATE\t    print i"%(threadId, frameId))
        self.WriteDebugConsoleExpression("%s\t%s\tEVALUATE\t"%(threadId, frameId))
        self.WaitForVars('<xml><more>False</more><output message="0"></output><output message="1"></output><output message="2"></output></xml>')
        assert 17 == self._sequence, 'Expected 19. Had: %s' % self._sequence
        
        self.WriteRunThread(threadId)
        self.finishedOk = True


#=======================================================================================================================
# WriterThreadCase13
#======================================================================================================================
class WriterThreadCase13(AbstractWriterThread):
    
    TEST_FILE = NormFile('_debugger_case13.py')
    
    def run(self):
        self.StartSocket()
        self.WriteAddBreakpoint(35, 'main')
        self.Write("124\t%s\t%s" % (self.NextSeq(), "true;false;false;true"))
        self.WriteMakeInitialRun()
        threadId, frameId, line = self.WaitForBreakpointHit('111', True)

        self.WriteGetFrame(threadId, frameId)
        
        self.WriteStepIn(threadId)
        threadId, frameId, line = self.WaitForBreakpointHit('107', True)
        # Should go inside setter method
        assert line == 25, 'Expected return to be in line 25, was: %s' % line
        
        self.WriteStepIn(threadId)
        threadId, frameId, line = self.WaitForBreakpointHit('107', True)

        self.WriteStepIn(threadId)
        threadId, frameId, line = self.WaitForBreakpointHit('107', True)
        # Should go inside getter method
        assert line == 21, 'Expected return to be in line 21, was: %s' % line

        self.WriteStepIn(threadId)
        threadId, frameId, line = self.WaitForBreakpointHit('107', True)

        # Disable property tracing
        self.Write("124\t%s\t%s" % (self.NextSeq(), "true;true;true;true"))
        self.WriteStepIn(threadId)
        threadId, frameId, line = self.WaitForBreakpointHit('107', True)
        # Should Skip step into properties setter
        assert line == 39, 'Expected return to be in line 39, was: %s' % line

        # Enable property tracing
        self.Write("124\t%s\t%s" % (self.NextSeq(), "true;false;false;true"))
        self.WriteStepIn(threadId)
        threadId, frameId, line = self.WaitForBreakpointHit('107', True)
        # Should go inside getter method
        assert line == 8, 'Expected return to be in line 8, was: %s' % line

        self.WriteRunThread(threadId)
        
        self.finishedOk = True

#=======================================================================================================================
# WriterThreadCase12
#======================================================================================================================
class WriterThreadCase12(AbstractWriterThread):
    
    TEST_FILE = NormFile('_debugger_case10.py')
        
    def run(self):
        self.StartSocket()
        self.WriteAddBreakpoint(2, '') #Should not be hit: setting empty function (not None) should only hit global.
        self.WriteAddBreakpoint(6, 'Method1a')  
        self.WriteAddBreakpoint(11, 'Method2') 
        self.WriteMakeInitialRun()
        
        threadId, frameId, line = self.WaitForBreakpointHit('111', True)
        
        assert line == 11, 'Expected return to be in line 11, was: %s' % line
        
        self.WriteStepReturn(threadId)
        
        threadId, frameId, line = self.WaitForBreakpointHit('111', True) #not a return (it stopped in the other breakpoint)
        
        assert line == 6, 'Expected return to be in line 6, was: %s' % line
        
        self.WriteRunThread(threadId)

        assert 13 == self._sequence, 'Expected 13. Had: %s' % self._sequence
        
        self.finishedOk = True
        


#=======================================================================================================================
# WriterThreadCase11
#======================================================================================================================
class WriterThreadCase11(AbstractWriterThread):
    
    TEST_FILE = NormFile('_debugger_case10.py')
        
    def run(self):
        self.StartSocket()
        self.WriteAddBreakpoint(2, 'Method1') 
        self.WriteMakeInitialRun()
        
        threadId, frameId = self.WaitForBreakpointHit('111')
        
        self.WriteStepOver(threadId)
        
        threadId, frameId, line = self.WaitForBreakpointHit('108', True)
        
        assert line == 3, 'Expected return to be in line 3, was: %s' % line
        
        self.WriteStepOver(threadId)
        
        threadId, frameId, line = self.WaitForBreakpointHit('108', True)
        
        assert line == 11, 'Expected return to be in line 11, was: %s' % line
        
        self.WriteStepOver(threadId)
        
        threadId, frameId, line = self.WaitForBreakpointHit('108', True)
        
        assert line == 12, 'Expected return to be in line 12, was: %s' % line
        
        self.WriteRunThread(threadId)

        assert 13 == self._sequence, 'Expected 13. Had: %s' % self._sequence
        
        self.finishedOk = True
        



#=======================================================================================================================
# WriterThreadCase10
#======================================================================================================================
class WriterThreadCase10(AbstractWriterThread):
    
    TEST_FILE = NormFile('_debugger_case10.py')
        
    def run(self):
        self.StartSocket()
        self.WriteAddBreakpoint(2, 'None') #None or Method should make hit. 
        self.WriteMakeInitialRun()
        
        threadId, frameId = self.WaitForBreakpointHit('111')
        
        self.WriteStepReturn(threadId)
        
        threadId, frameId, line = self.WaitForBreakpointHit('109', True)
        
        assert line == 11, 'Expected return to be in line 11, was: %s' % line
        
        self.WriteStepOver(threadId)
        
        threadId, frameId, line = self.WaitForBreakpointHit('108', True)
        
        assert line == 12, 'Expected return to be in line 12, was: %s' % line
        
        self.WriteRunThread(threadId)

        assert 11 == self._sequence, 'Expected 11. Had: %s' % self._sequence
        
        self.finishedOk = True
        


#=======================================================================================================================
# WriterThreadCase9
#======================================================================================================================
class WriterThreadCase9(AbstractWriterThread):
    
    TEST_FILE = NormFile('_debugger_case89.py')
        
    def run(self):
        self.StartSocket()
        self.WriteAddBreakpoint(10, 'Method3') 
        self.WriteMakeInitialRun()
        
        threadId, frameId = self.WaitForBreakpointHit('111')
        
        self.WriteStepOver(threadId)
        
        threadId, frameId, line = self.WaitForBreakpointHit('108', True)
        
        assert line == 11, 'Expected return to be in line 11, was: %s' % line
        
        self.WriteStepOver(threadId)
        
        threadId, frameId, line = self.WaitForBreakpointHit('108', True)
        
        assert line == 12, 'Expected return to be in line 12, was: %s' % line
        
        self.WriteRunThread(threadId)

        assert 11 == self._sequence, 'Expected 11. Had: %s' % self._sequence
        
        self.finishedOk = True
        

#=======================================================================================================================
# WriterThreadCase8
#======================================================================================================================
class WriterThreadCase8(AbstractWriterThread):
    
    TEST_FILE = NormFile('_debugger_case89.py')
        
    def run(self):
        self.StartSocket()
        self.WriteAddBreakpoint(10, 'Method3') 
        self.WriteMakeInitialRun()
        
        threadId, frameId = self.WaitForBreakpointHit('111')
        
        self.WriteStepReturn(threadId)
        
        threadId, frameId, line = self.WaitForBreakpointHit('109', True)
        
        assert line == 15, 'Expected return to be in line 15, was: %s' % line
        
        self.WriteRunThread(threadId)

        assert 9 == self._sequence, 'Expected 9. Had: %s' % self._sequence
        
        self.finishedOk = True
        



#=======================================================================================================================
# WriterThreadCase7
#======================================================================================================================
class WriterThreadCase7(AbstractWriterThread):
    
    TEST_FILE = NormFile('_debugger_case7.py')
        
    def run(self):
        self.StartSocket()
        self.WriteAddBreakpoint(2, 'Call') 
        self.WriteMakeInitialRun()
        
        threadId, frameId = self.WaitForBreakpointHit('111')
        
        self.WriteGetFrame(threadId, frameId)
        
        self.WaitForVars('<xml></xml>') #no vars at this point
        
        self.WriteStepOver(threadId)
        
        self.WriteGetFrame(threadId, frameId)
        
        self.WaitForVars('<xml><var name="variable_for_test_1" type="int" value="int%253A 10" />%0A</xml>')
        
        self.WriteStepOver(threadId)
        
        self.WriteGetFrame(threadId, frameId)
        
        self.WaitForVars('<xml><var name="variable_for_test_1" type="int" value="int%253A 10" />%0A<var name="variable_for_test_2" type="int" value="int%253A 20" />%0A</xml>')
        
        self.WriteRunThread(threadId)

        assert 17 == self._sequence, 'Expected 17. Had: %s' % self._sequence
        
        self.finishedOk = True
        


#=======================================================================================================================
# WriterThreadCase6
#=======================================================================================================================
class WriterThreadCase6(AbstractWriterThread):
    
    TEST_FILE = NormFile('_debugger_case56.py')
        
    def run(self):
        self.StartSocket()
        self.WriteAddBreakpoint(2, 'Call2') 
        self.WriteMakeInitialRun()
        
        threadId, frameId = self.WaitForBreakpointHit()
        
        self.WriteGetFrame(threadId, frameId)
        
        self.WriteStepReturn(threadId)
        
        threadId, frameId, line = self.WaitForBreakpointHit('109', True)

        assert line == 8, 'Expecting it to go to line 8. Went to: %s' % line
        
        self.WriteStepIn(threadId)
        
        threadId, frameId, line = self.WaitForBreakpointHit('107', True)
        
        #goes to line 4 in jython (function declaration line)
        assert line in (4, 5), 'Expecting it to go to line 4 or 5. Went to: %s' % line
        
        self.WriteRunThread(threadId)

        assert 13 == self._sequence, 'Expected 15. Had: %s' % self._sequence
        
        self.finishedOk = True

#=======================================================================================================================
# WriterThreadCase5
#=======================================================================================================================
class WriterThreadCase5(AbstractWriterThread):
    
    TEST_FILE = NormFile('_debugger_case56.py')
        
    def run(self):
        self.StartSocket()
        self.WriteAddBreakpoint(2, 'Call2') 
        self.WriteMakeInitialRun()
        
        threadId, frameId = self.WaitForBreakpointHit()
        
        self.WriteGetFrame(threadId, frameId)
        
        self.WriteRemoveBreakpoint(2)
        
        self.WriteStepReturn(threadId)
        
        threadId, frameId, line = self.WaitForBreakpointHit('109', True)

        assert line == 8, 'Expecting it to go to line 8. Went to: %s' % line
        
        self.WriteStepIn(threadId)
        
        threadId, frameId, line = self.WaitForBreakpointHit('107', True)
        
        #goes to line 4 in jython (function declaration line)
        assert line in (4, 5), 'Expecting it to go to line 4 or 5. Went to: %s' % line
        
        self.WriteRunThread(threadId)

        assert 15 == self._sequence, 'Expected 15. Had: %s' % self._sequence
        
        self.finishedOk = True


#=======================================================================================================================
# WriterThreadCase4
#=======================================================================================================================
class WriterThreadCase4(AbstractWriterThread):
    
    TEST_FILE = NormFile('_debugger_case4.py')
        
    def run(self):
        self.StartSocket()
        self.WriteMakeInitialRun()
        
        threadId = self.WaitForNewThread()
        
        self.WriteSuspendThread(threadId)

        time.sleep(4) #wait for time enough for the test to finish if it wasn't suspended
        
        self.WriteRunThread(threadId)
        
        self.finishedOk = True


#=======================================================================================================================
# WriterThreadCase3
#=======================================================================================================================
class WriterThreadCase3(AbstractWriterThread):
    
    TEST_FILE = NormFile('_debugger_case3.py')
        
    def run(self):
        self.StartSocket()
        self.WriteMakeInitialRun()
        time.sleep(1)
        self.WriteAddBreakpoint(4, '') 
        self.WriteAddBreakpoint(5, 'FuncNotAvailable') #Check that it doesn't get hit in the global when a function is available
        
        threadId, frameId = self.WaitForBreakpointHit()
        
        self.WriteGetFrame(threadId, frameId)
        
        self.WriteRunThread(threadId)
        
        threadId, frameId = self.WaitForBreakpointHit()
        
        self.WriteGetFrame(threadId, frameId)
        
        self.WriteRemoveBreakpoint(4)
        
        self.WriteRunThread(threadId)
        
        assert 17 == self._sequence, 'Expected 17. Had: %s' % self._sequence
        
        self.finishedOk = True

#=======================================================================================================================
# WriterThreadCase2
#=======================================================================================================================
class WriterThreadCase2(AbstractWriterThread):
    
    TEST_FILE = NormFile('_debugger_case2.py')
        
    def run(self):
        self.StartSocket()
        self.WriteAddBreakpoint(3, 'Call4') #seq = 3
        self.WriteMakeInitialRun()
        
        threadId, frameId = self.WaitForBreakpointHit()
        
        self.WriteGetFrame(threadId, frameId)
        
        self.WriteAddBreakpoint(14, 'Call2')
        
        self.WriteRunThread(threadId)
        
        threadId, frameId = self.WaitForBreakpointHit()
        
        self.WriteGetFrame(threadId, frameId)
        
        self.WriteRunThread(threadId)
        
        assert 15 == self._sequence, 'Expected 15. Had: %s' % self._sequence
        
        self.finishedOk = True
        
#=======================================================================================================================
# WriterThreadCase1
#=======================================================================================================================
class WriterThreadCase1(AbstractWriterThread):
    
    TEST_FILE = NormFile('_debugger_case1.py')
        
    def run(self):
        self.StartSocket()
        self.WriteAddBreakpoint(6, 'SetUp')
        self.WriteMakeInitialRun()
        
        threadId, frameId = self.WaitForBreakpointHit()
        
        self.WriteGetFrame(threadId, frameId)

        self.WriteStepOver(threadId)
        
        self.WriteGetFrame(threadId, frameId)
        
        self.WriteRunThread(threadId)
        
        assert 13 == self._sequence, 'Expected 13. Had: %s' % self._sequence
        
        self.finishedOk = True
        
#=======================================================================================================================
# Test
#=======================================================================================================================
class Test(unittest.TestCase):
    
    def CheckCase(self, writerThreadClass, run_as_python=True):
        UpdatePort()
        writerThread = writerThreadClass()
        writerThread.start()
        
        import pydev_localhost
        localhost = pydev_localhost.get_localhost()
        if run_as_python:
            args = [
                'python',
                PYDEVD_FILE,
                '--DEBUG_RECORD_SOCKET_READS',
                '--client',
                localhost,
                '--port',
                str(port),
                '--file',
                writerThread.TEST_FILE,
            ]
            
        else:
            #run as jython
            args = [
                JAVA_LOCATION,
                '-classpath',
                JYTHON_JAR_LOCATION,
                'org.python.util.jython',
                PYDEVD_FILE,
                '--DEBUG_RECORD_SOCKET_READS',
                '--client',
                localhost,
                '--port',
                str(port),
                '--file',
                writerThread.TEST_FILE,
            ]
        
        if SHOW_OTHER_DEBUG_INFO:
            print 'executing', ' '.join(args)
            
        process = subprocess.Popen(args, stdout=subprocess.PIPE, stderr=subprocess.STDOUT, cwd=os.path.dirname(PYDEVD_FILE))
        class ProcessReadThread(threading.Thread):
            def run(self):
                self.resultStr = None
                self.resultStr = process.stdout.read()
                process.stdout.close()
                
            def DoKill(self):
                process.stdout.close()
                
        processReadThread = ProcessReadThread()
        processReadThread.setDaemon(True)
        processReadThread.start()
        if SHOW_OTHER_DEBUG_INFO:
            print 'Both processes started'
        
        #polls can fail (because the process may finish and the thread still not -- so, we give it some more chances to
        #finish successfully).
        pools_failed = 0
        while writerThread.isAlive():
            if process.poll() is not None:
                pools_failed += 1
            time.sleep(.2)
            if pools_failed == 10:
                break
        
        if process.poll() is None:
            for i in range(10):
                if processReadThread.resultStr is None:
                    time.sleep(.5)
                else:
                    break
            else:
                writerThread.DoKill()
        
        else:
            if process.poll() < 0:
                self.fail("The other process exited with error code: " + str(process.poll()) + " result:" + processReadThread.resultStr)
                    
        
        if SHOW_RESULT_STR:
            print processReadThread.resultStr
            
        if processReadThread.resultStr is None:
            self.fail("The other process may still be running -- and didn't give any output")
            
        if 'TEST SUCEEDED' not in processReadThread.resultStr:
            self.fail(processReadThread.resultStr)
            
        if not writerThread.finishedOk:
            self.fail("The thread that was doing the tests didn't finish successfully. Output: %s" % processReadThread.resultStr)
            

            
    def testCase1(self):
        self.CheckCase(WriterThreadCase1)
        
    def testCase2(self):
        self.CheckCase(WriterThreadCase2)
        
    def testCase3(self):
        self.CheckCase(WriterThreadCase3)
        
    def testCase4(self):
        self.CheckCase(WriterThreadCase4)
            
    def testCase5(self):
        self.CheckCase(WriterThreadCase5)
            
    def testCase6(self):
        self.CheckCase(WriterThreadCase6)
        
    def testCase7(self):
        self.CheckCase(WriterThreadCase7)
        
    def testCase8(self):
        self.CheckCase(WriterThreadCase8)
        
    def testCase9(self):
        self.CheckCase(WriterThreadCase9)
        
    def testCase10(self):
        self.CheckCase(WriterThreadCase10)
        
    def testCase11(self):
        self.CheckCase(WriterThreadCase11)
        
    def testCase12(self):
        self.CheckCase(WriterThreadCase12)
        
    def testCase13(self):
        self.CheckCase(WriterThreadCase13)

    def testCase14(self):
        self.CheckCase(WriterThreadCase14)
            
    def testCase1a(self):
        self.CheckCase(WriterThreadCase1, False)
        
    def testCase2a(self):
        self.CheckCase(WriterThreadCase2, False)
        
    def testCase3a(self):
        self.CheckCase(WriterThreadCase3, False)
        
    def testCase4a(self):
        self.CheckCase(WriterThreadCase4, False)
        
    def testCase5a(self):
        self.CheckCase(WriterThreadCase5, False)
        
    def testCase6a(self):
        self.CheckCase(WriterThreadCase6, False)
        
    def testCase7a(self):
        self.CheckCase(WriterThreadCase7, False)
        
    def testCase8a(self):
        self.CheckCase(WriterThreadCase8, False)
        
    def testCase9a(self):
        self.CheckCase(WriterThreadCase9, False)
        
    def testCase10a(self):
        self.CheckCase(WriterThreadCase10, False)
        
    def testCase11a(self):
        self.CheckCase(WriterThreadCase11, False)
        
    def testCase12a(self):
        self.CheckCase(WriterThreadCase12, False)

    def testCase14a(self):
        self.CheckCase(WriterThreadCase14, False)

#This case requires decorators to work (which are not present on Jython 2.1), so, this test is just removed from the jython run.
#    def testCase13a(self):
#        self.CheckCase(WriterThreadCase13, False)

def GetLocationFromLine(line):
    loc = line.split('=')[1].strip()
    if loc.endswith(';'):
        loc = loc[:-1]
    if loc.endswith('"'):
        loc = loc[:-1]
    if loc.startswith('"'):
        loc = loc[1:]
    return loc
    
#=======================================================================================================================
# Main        
#=======================================================================================================================
if __name__ == '__main__':
    import platform
    sysname = platform.system().lower()
    test_dependent = os.path.join('../../../', 'org.python.pydev.core', 'tests', 'org', 'python', 'pydev', 'core', 'TestDependent.' + sysname + '.properties')
    f = open(test_dependent)
    try:
        for line in f.readlines():
            if 'JYTHON_JAR_LOCATION' in line:
                JYTHON_JAR_LOCATION = GetLocationFromLine(line)
                
            if 'JAVA_LOCATION' in line:
                JAVA_LOCATION = GetLocationFromLine(line)
    finally:
        f.close()
        
    assert JYTHON_JAR_LOCATION, 'JYTHON_JAR_LOCATION not found in %s' % (test_dependent,)
    assert JAVA_LOCATION, 'JAVA_LOCATION not found in %s' % (test_dependent,)
    assert os.path.exists(JYTHON_JAR_LOCATION), 'The location: %s is not valid' % (JYTHON_JAR_LOCATION,)
    assert os.path.exists(JAVA_LOCATION), 'The location: %s is not valid' % (JAVA_LOCATION,)

    suite = unittest.makeSuite(Test)
    
    suite = unittest.TestSuite()
    suite.addTest(Test('testCase14'))
#    suite.addTest(Test('testCase10a'))
    unittest.TextTestRunner(verbosity=3).run(suite)

