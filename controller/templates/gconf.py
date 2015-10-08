bind = '0.0.0.0'
try:
    workers = int({{ if exists "/deis/controller/workers" }}{{ getv "/deis/controller/workers" }}{{ else }}"not set"{{end}})
    if workers < 1:
        raise ValueError()
except (NameError, ValueError):
    import multiprocessing
    try:
        workers = multiprocessing.cpu_count() * 2 + 1
    except NotImplementedError:
        workers = 8
timeout = 1200
pidfile = '/tmp/gunicorn.pid'
loglevel = 'info'
errorlog = '-'
accesslog = '-'
access_log_format = '%(h)s "%(r)s" %(s)s %(b)s "%(a)s"'


def worker_int(worker):
    """Print a stack trace when a worker receives a SIGINT or SIGQUIT signal."""
    worker.log.warning('worker terminated')
    import traceback
    traceback.print_stack()


def worker_abort(worker):
    """Print a stack trace when a worker receives a SIGABRT signal, generally on timeout."""
    worker.log.warning('worker aborted')
    import traceback
    traceback.print_stack()
