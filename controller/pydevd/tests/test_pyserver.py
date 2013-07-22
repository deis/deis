'''
@author Fabio Zadrozny 
'''
import sys
import os

#make it as if we were executing from the directory above this one (so that we can use pycompletionserver
#without the need for it being in the pythonpath)
sys.argv[0] = os.path.dirname(sys.argv[0]) 
#twice the dirname to get the previous level from this file.
sys.path.insert(1, os.path.join(os.path.dirname(sys.argv[0])))

IS_PYTHON_3K = 0
if sys.platform.find('java') == -1:
    
    
    try:
        import inspect
        import pycompletionserver
        import socket
        try:
            from urllib import quote_plus, unquote_plus
            def send(s, msg):
                s.send(msg)
        except ImportError:
            IS_PYTHON_3K = 1
            from urllib.parse import quote_plus, unquote_plus #Python 3.0
            def send(s, msg):
                s.send(bytearray(msg, 'utf-8'))
    except ImportError:
        pass #Not available in jython
    
    import unittest
    
    class Test(unittest.TestCase):
    
        def setUp(self):
            unittest.TestCase.setUp(self)
    
        def tearDown(self):
            unittest.TestCase.tearDown(self)
        
        def testMessage(self):
            t = pycompletionserver.T(0, 0)
            
            l = []
            l.append(('Def', 'description'  , 'args'))
            l.append(('Def1', 'description1', 'args1'))
            l.append(('Def2', 'description2', 'args2'))
            
            msg = t.processor.formatCompletionMessage(None, l)
            self.assertEquals('@@COMPLETIONS(None,(Def,description,args),(Def1,description1,args1),(Def2,description2,args2))END@@', msg)
            
            l = []
            l.append(('Def', 'desc,,r,,i()ption', ''))
            l.append(('Def(1', 'descriptio(n1', ''))
            l.append(('De,f)2', 'de,s,c,ription2', ''))
            msg = t.processor.formatCompletionMessage(None, l)
            self.assertEquals('@@COMPLETIONS(None,(Def,desc%2C%2Cr%2C%2Ci%28%29ption, ),(Def%281,descriptio%28n1, ),(De%2Cf%292,de%2Cs%2Cc%2Cription2, ))END@@', msg)
    
        def createConnections(self, p1=50002, p2=50003):
            '''
            Creates the connections needed for testing.
            '''
            t = pycompletionserver.T(p1, p2)
            
            t.start()
    
            sToWrite = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            sToWrite.connect((pycompletionserver.HOST, p1))
            
            sToRead = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            sToRead.bind((pycompletionserver.HOST, p2))
            sToRead.listen(1) #socket to receive messages.
    
            connToRead, addr = sToRead.accept()
    
            return t, sToWrite, sToRead, connToRead, addr
            
    
        def readMsg(self):
            finish = False
            msg = ''
            while finish == False:
                m = self.connToRead.recv(1024 * 4)
                if IS_PYTHON_3K:
                    m = m.decode('utf-8')
                if m.startswith('@@PROCESSING'):
                    sys.stdout.write('Status msg: %s\n' % (msg,))
                else:
                    msg += m
    
                if msg.find('END@@') != -1:
                    finish = True
    
            return msg
    
        def testCompletionSocketsAndMessages(self):
            t, sToWrite, sToRead, self.connToRead, addr = self.createConnections()
            
            try:
                #now that we have the connections all set up, check the code completion messages.
                msg = quote_plus('math')
                send(sToWrite, '@@IMPORTS:%sEND@@' % msg) #math completions
                completions = self.readMsg()
                #print_ unquote_plus(completions)
                
                #math is a builtin and because of that, it starts with None as a file
                start = '@@COMPLETIONS(None,(__doc__,'
                start_2 = '@@COMPLETIONS(None,(__name__,'
                self.assert_(completions.startswith(start) or completions.startswith(start_2), '%s DOESNT START WITH %s' % (completions, (start, start_2)))
        
                self.assert_('@@COMPLETIONS' in completions)
                self.assert_('END@@' in completions)
    
    
                #now, test i
                msg = quote_plus('__builtin__.list')
                send(sToWrite, "@@IMPORTS:%s\nEND@@" % msg)
                found = self.readMsg()
                self.assert_('sort' in found, 'Could not find sort in: %s' % (found,))
    
                #now, test search
                msg = quote_plus('inspect.ismodule')
                send(sToWrite, '@@SEARCH%sEND@@' % msg) #math completions
                found = self.readMsg()
                self.assert_('inspect.py' in found)
                self.assert_('33' in found or '34' in found or '51' in found or '50' in found, 'Could not find 33, 34, 50 or 51 in %s' % found)
    
                #now, test search
                msg = quote_plus('inspect.CO_NEWLOCALS')
                send(sToWrite, '@@SEARCH%sEND@@' % msg) #math completions
                found = self.readMsg()
                self.assert_('inspect.py' in found)
                self.assert_('CO_NEWLOCALS' in found)
                
                #now, test search
                msg = quote_plus('inspect.BlockFinder.tokeneater')
                send(sToWrite, '@@SEARCH%sEND@@' % msg) 
                found = self.readMsg()
                self.assert_('inspect.py' in found)
    #            self.assert_('CO_NEWLOCALS' in found)
    
            #reload modules test
    #        send(sToWrite, '@@RELOAD_MODULES_END@@')
    #        ok = self.readMsg()
    #        self.assertEquals('@@MSG_OK_END@@' , ok)
    #        this test is not executed because it breaks our current enviroment.
            
            
            finally:
                try:
                    sys.stdout.write('succedded...sending kill msg\n')
                    self.sendKillMsg(sToWrite)
                    
            
    #                while not hasattr(t, 'ended'):
    #                    pass #wait until it receives the message and quits.
            
                        
                    sToRead.close()
                    sToWrite.close()
                    self.connToRead.close()
                except:
                    pass
            
        def sendKillMsg(self, socket):
            socket.send(pycompletionserver.MSG_KILL_SERVER)

        
if __name__ == '__main__':
    if sys.platform.find('java') == -1:
        unittest.main()
    else:
        sys.stdout.write('Not running python tests in platform: %s\n' % (sys.platform,))

