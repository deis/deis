'''
This module was created to get information available in the interpreter, such as libraries,
paths, etc.

what is what:
sys.builtin_module_names: contains the builtin modules embeeded in python (rigth now, we specify all manually).
sys.prefix: A string giving the site-specific directory prefix where the platform independent Python files are installed

format is something as 
EXECUTABLE:python.exe|libs@compiled_dlls$builtin_mods

all internal are separated by |
'''
import sys

try:
    import os.path
    def fullyNormalizePath(path):
        '''fixes the path so that the format of the path really reflects the directories in the system
        '''
        return os.path.normpath(path)
    join = os.path.join
except: #ImportError or AttributeError.
    #See: http://stackoverflow.com/questions/10254353/error-while-installing-jython-for-pydev
    def fullyNormalizePath(path):
        '''fixes the path so that the format of the path really reflects the directories in the system
        '''
        return path

    def join(a, b):
        if a.endswith('/') or a.endswith('\\'):
            return a + b
        return a + '/' + b


IS_PYTHON_3K = 0

try:
    if sys.version_info[0] == 3:
        IS_PYTHON_3K = 1
except:
    #That's OK, not all versions of python have sys.version_info
    pass

try:
    #Just check if False and True are defined (depends on version, not whether it's jython/python)
    False
    True
except:
    exec ('True, False = 1,0') #An exec is used so that python 3k does not give a syntax error

import time

if sys.platform == "cygwin":

    try:
        import ctypes #use from the system if available
    except ImportError:
        sys.path.append(join(sys.path[0], 'third_party/wrapped_for_pydev'))
        import ctypes

    def nativePath(path):
        MAX_PATH = 512  # On cygwin NT, its 260 lately, but just need BIG ENOUGH buffer
        '''Get the native form of the path, like c:\\Foo for /cygdrive/c/Foo'''

        retval = ctypes.create_string_buffer(MAX_PATH)
        path = fullyNormalizePath(path)
        ctypes.cdll.cygwin1.cygwin_conv_to_win32_path(path, retval) #@UndefinedVariable
        return retval.value

else:
    def nativePath(path):
        return fullyNormalizePath(path)



def getfilesystemencoding():
    try:
        ret = sys.getfilesystemencoding()
        if not ret:
            raise RuntimeError('Unable to get encoding.')
        return ret
    except:
        #Only available from 2.3 onwards.
        if sys.platform == 'win32':
            return 'mbcs'
        return 'utf-8'

file_system_encoding = getfilesystemencoding()

def tounicode(s):
    if hasattr(s, 'decode'):
        #Depending on the platform variant we may have decode on string or not.
        return s.decode(file_system_encoding)
    return s

def toutf8(s):
    if hasattr(s, 'encode'):
        return s.encode('utf-8')
    return s


if __name__ == '__main__':
    try:
        #just give some time to get the reading threads attached (just in case)
        time.sleep(0.1)
    except:
        pass

    try:
        executable = nativePath(sys.executable)
    except:
        executable = sys.executable

    if sys.platform == "cygwin" and not executable.endswith('.exe'):
        executable += '.exe'


    try:
        major = str(sys.version_info[0])
        minor = str(sys.version_info[1])
    except AttributeError:
        #older versions of python don't have version_info
        import string
        s = string.split(sys.version, ' ')[0]
        s = string.split(s, '.')
        major = s[0]
        minor = s[1]

    s = tounicode('%s.%s') % (tounicode(major), tounicode(minor))

    contents = [tounicode('<xml>')]
    contents.append(tounicode('<version>%s</version>') % (tounicode(s),))

    contents.append(tounicode('<executable>%s</executable>') % tounicode(executable))

    #this is the new implementation to get the system folders 
    #(still need to check if it works in linux)
    #(previously, we were getting the executable dir, but that is not always correct...)
    prefix = tounicode(nativePath(sys.prefix))
    #print_ 'prefix is', prefix


    result = []

    path_used = sys.path
    try:
        path_used = path_used[:] #Use a copy.
    except:
        pass #just ignore it...

    for p in path_used:
        p = tounicode(nativePath(p))

        try:
            import string #to be compatible with older versions
            if string.find(p, prefix) == 0: #was startswith
                result.append((p, True))
            else:
                result.append((p, False))
        except (ImportError, AttributeError):
            #python 3k also does not have it
            #jython may not have it (depending on how are things configured)
            if p.startswith(prefix): #was startswith
                result.append((p, True))
            else:
                result.append((p, False))

    for p, b in result:
        if b:
            contents.append(tounicode('<lib path="ins">%s</lib>') % (p,))
        else:
            contents.append(tounicode('<lib path="out">%s</lib>') % (p,))

    #no compiled libs
    #nor forced libs

    for builtinMod in sys.builtin_module_names:
        contents.append(tounicode('<forced_lib>%s</forced_lib>') % tounicode(builtinMod))


    contents.append(tounicode('</xml>'))
    unic = tounicode('\n').join(contents)
    inutf8 = toutf8(unic)
    if IS_PYTHON_3K:
        #This is the 'official' way of writing binary output in Py3K (see: http://bugs.python.org/issue4571)
        sys.stdout.buffer.write(inutf8)
    else:
        sys.stdout.write(inutf8)

    try:
        sys.stdout.flush()
        sys.stderr.flush()
        #and give some time to let it read things (just in case)
        time.sleep(0.1)
    except:
        pass

    raise RuntimeError('Ok, this is so that it shows the output (ugly hack for some platforms, so that it releases the output).')
