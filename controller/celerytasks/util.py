import StringIO
import paramiko
import select
import socket
import time


def connect_ssh(username, hostname, port, key):
    key_f = StringIO.StringIO(key)
    pkey = paramiko.RSAKey.from_private_key(key_f)
    ssh = paramiko.SSHClient()
    ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
    for _ in range(10):
        try:
            ssh.connect(hostname, username=username, pkey=pkey)
            break
        except paramiko.AuthenticationException as e:
            raise paramiko.AuthenticationException(e)
        except socket.error:
            time.sleep(3)
    else:
        raise RuntimeError('SSH Connection Error')
    return ssh


def exec_ssh(ssh, command):
    tran = ssh.get_transport()
    chan = tran.open_session()
    # NOTE: pty breaks line ordering on commands like apt-get
    #chan.get_pty(term='vt100', width=300, height=50)
    chan.exec_command(command)
    output = read_from_ssh(chan)
    exit_status = chan.recv_exit_status()
    return output, exit_status


def read_from_ssh(chan):
    output = ''
    while True:
        r, w, e = select.select([chan], [], [], 10)  # @UnusedVariable
        if r:
            got_data = False
            if chan.recv_ready():
                data = r[0].recv(4096)
                if data:
                    got_data = True
                    output += data
                    #print("stdout => ", data)
            if chan.recv_stderr_ready():
                data = r[0].recv_stderr(4096)
                if data:
                    got_data = True
                    output += data
                    #print("stderr => ", data)
            if not got_data:
                return output
