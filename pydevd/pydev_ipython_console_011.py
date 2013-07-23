from IPython.frontend.terminal.interactiveshell import TerminalInteractiveShell
from IPython.core.inputsplitter import IPythonInputSplitter
from IPython.utils import io
import sys
import codeop, re
original_stdout = sys.stdout
original_stderr = sys.stderr


#=======================================================================================================================
# _showtraceback
#=======================================================================================================================
def _showtraceback(*args, **kwargs):
    import traceback;traceback.print_exc()
    
    
    
#=======================================================================================================================
# PyDevFrontEnd
#=======================================================================================================================
class PyDevFrontEnd:

    def __init__(self, *args, **kwargs):        
        #Initialization based on: from IPython.testing.globalipapp import start_ipython
        
        self._curr_exec_line = 0
        # Store certain global objects that IPython modifies
        _displayhook = sys.displayhook
        _excepthook = sys.excepthook
    
        # Create and initialize our IPython instance.
        shell = TerminalInteractiveShell.instance()
        # Create an intput splitter to handle input separation
        self.input_splitter = IPythonInputSplitter()

        shell.showtraceback = _showtraceback
        # IPython is ready, now clean up some global state...
        
        # Deactivate the various python system hooks added by ipython for
        # interactive convenience so we don't confuse the doctest system
        sys.displayhook = _displayhook
        sys.excepthook = _excepthook
    
        # So that ipython magics and aliases can be doctested (they work by making
        # a call into a global _ip object).  Also make the top-level get_ipython
        # now return this without recursively calling here again.
        get_ipython = shell.get_ipython
        try:
            import __builtin__
        except:
            import builtins as __builtin__
        __builtin__._ip = shell
        __builtin__.get_ipython = get_ipython
        
        # We want to print to stdout/err as usual.
        io.stdout = original_stdout
        io.stderr = original_stderr

        self.ipython = shell

    def complete(self, string):
        return self.ipython.complete(None, line=string)

    def is_complete(self, string):
        return  not self.input_splitter.push_accepts_more()

    def getNamespace(self):
        return self.ipython.user_ns

    def addExec(self, line):
        self.input_splitter.push(line)
        if self.is_complete(line):
            self.ipython.run_cell(self.input_splitter.source_reset())
            return False
        else:
            return True
