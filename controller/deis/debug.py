
from pydevd import pydevd
pydevd.settrace('box.techverity.com', stdoutToServer=True, stderrToServer=True, port=33000, suspend=False,
                trace_only_current_thread=False)

from .settings import *  # @UnusedWildImport
