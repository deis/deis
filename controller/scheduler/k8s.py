import copy
import httplib
import json
import random
import re
import string
import time

from django.conf import settings
from docker import Client
from .states import JobState
from . import AbstractSchedulerClient


POD_TEMPLATE = '''{
  "kind": "Pod",
  "apiVersion": "$version",
  "metadata": {
    "name": "$id"
  },
  "spec": {
    "containers": [
      {
        "name": "$id",
        "image": "$image"
      }
    ],
    "restartPolicy":"Never"
  }
}'''

RC_TEMPLATE = '''{
   "kind":"ReplicationController",
   "apiVersion":"$version",
   "metadata":{
      "name":"$name",
      "labels":{
         "name":"$id"
      }
   },
   "spec":{
      "replicas":$num,
      "selector":{
         "name":"$id",
         "version":"$appversion",
         "type":"$type"
      },
      "template":{
         "metadata":{
            "labels":{
               "name":"$id",
               "version":"$appversion",
               "type":"$type"
            }
         },
         "spec":{
            "containers":[
               {
                  "name":"$containername",
                  "image":"$image"
               }
            ]
         }
      }
   }
}'''

SERVICE_TEMPLATE = '''{
   "kind":"Service",
   "apiVersion":"$version",
   "metadata":{
      "name":"$name",
      "labels":{
         "name":"$label"
      }
   },
   "spec":{
      "ports": [
        {
          "port":80,
          "targetPort":$port,
          "protocol":"TCP"
        }
      ],
      "selector":{
         "name":"$label",
         "type":"$type"
      }
   }
}'''

POD_DELETE = '''{
}'''


RETRIES = 3
MATCH = re.compile(
    r'(?P<app>[a-z0-9-]+)_?(?P<version>v[0-9]+)?\.?(?P<c_type>[a-z-_]+)')


class KubeHTTPClient(AbstractSchedulerClient):

    def __init__(self, target, auth, options, pkey):
        super(KubeHTTPClient, self).__init__(target, auth, options, pkey)
        self.target = settings.K8S_MASTER
        self.port = "8080"
        self.registry = settings.REGISTRY_HOST+":"+settings.REGISTRY_PORT
        self.apiversion = "v1"
        self.conn = httplib.HTTPConnection(self.target+":"+self.port)

    def _get_old_rc(self, name, app_type):
        con_app = httplib.HTTPConnection(self.target+":"+self.port)
        con_app.request('GET', '/api/'+self.apiversion +
                        '/namespaces/'+name+'/replicationcontrollers')
        resp = con_app.getresponse()
        data = resp.read()
        reason = resp.reason
        status = resp.status
        con_app.close()
        if not 200 <= status <= 299:
            errmsg = "Failed to get Replication Controllers: {} {} - {}".format(
                status, reason, data)
            raise RuntimeError(errmsg)
        parsed_json = json.loads(data)
        exists = False
        prev_rc = []
        for rc in parsed_json['items']:
            if('name' in rc['metadata']['labels'] and name == rc['metadata']['labels']['name'] and
               'type' in rc['spec']['selector'] and app_type == rc['spec']['selector']['type']):
                exists = True
                prev_rc = rc
                break
        if exists:
            return prev_rc
        else:
            return 0

    def _get_rc_status(self, name, namespace):
        conn_rc = httplib.HTTPConnection(self.target+":"+self.port)
        conn_rc.request('GET', '/api/'+self.apiversion+'/' +
                        'namespaces/'+namespace+'/replicationcontrollers/'+name)
        resp = conn_rc.getresponse()
        status = resp.status
        conn_rc.close()
        return status

    def _get_rc_(self, name, namespace):
        conn_rc_resver = httplib.HTTPConnection(self.target+":"+self.port)
        conn_rc_resver.request('GET', '/api/'+self.apiversion+'/' +
                               'namespaces/'+namespace+'/replicationcontrollers/'+name)
        resp = conn_rc_resver.getresponse()
        data = resp.read()
        reason = resp.reason
        status = resp.status
        conn_rc_resver.close()
        if not 200 <= status <= 299:
            errmsg = "Failed to get Replication Controller:{} {} {} - {}".format(
                name, status, reason, data)
            raise RuntimeError(errmsg)
        parsed_json = json.loads(data)
        return parsed_json

    def deploy(self, name, image, command, **kwargs):
        app_name = kwargs.get('aname', {})
        app_type = name.split(".")[1]
        old_rc = self._get_old_rc(app_name, app_type)
        new_rc = self._create_rc(name, image, command, **kwargs)
        desired = int(old_rc["spec"]["replicas"])
        old_rc_name = old_rc["metadata"]["name"]
        new_rc_name = new_rc["metadata"]["name"]
        try:
            count = 1
            while desired >= count:
                new_rc = self._scale_app(new_rc_name, count, app_name)
                old_rc = self._scale_app(old_rc_name, desired-count, app_name)
                count += 1
        except Exception as e:
            self._scale_app(new_rc["metadata"]["name"], 0, app_name)
            self._delete_rc(new_rc["metadata"]["name"], app_name)
            self._scale_app(old_rc["metadata"]["name"], desired, app_name)
            err = '{} (deploy): {}'.format(name, e)
            raise RuntimeError(err)
        self._delete_rc(old_rc_name, app_name)

    def _get_events(self, namespace):
        con_get = httplib.HTTPConnection(self.target+":"+self.port)
        con_get.request('GET', '/api/'+self.apiversion+'/namespaces/'+namespace+'/events')
        resp = con_get.getresponse()
        reason = resp.reason
        status = resp.status
        data = resp.read()
        con_get.close()
        if not 200 <= status <= 299:
            errmsg = "Failed to get events: {} {} - {}".format(
                status, reason, data)
            raise RuntimeError(errmsg)
        return (status, data, reason)

    def _get_schedule_status(self, name, num, namespace):
        pods = []
        for _ in xrange(120):
            count = 0
            pods = []
            status, data, reason = self._get_pods(namespace)
            parsed_json = json.loads(data)
            for pod in parsed_json['items']:
                if pod['metadata']['generateName'] == name+'-':
                    count += 1
                    pods.append(pod['metadata']['name'])
            if count == num:
                break
            time.sleep(1)
        for _ in xrange(120):
            count = 0
            status, data, reason = self._get_events(namespace)
            parsed_json = json.loads(data)
            for event in parsed_json['items']:
                if(event['involvedObject']['name'] in pods and
                   event['source']['component'] == 'scheduler'):
                    if event['reason'] == 'scheduled':
                        count += 1
                    else:
                        raise RuntimeError(event['message'])
            if count == num:
                break
            time.sleep(1)

    def _scale_rc(self, rc, namespace):
        name = rc['metadata']['name']
        num = rc["spec"]["replicas"]
        headers = {'Content-Type': 'application/json'}
        conn_scalepod = httplib.HTTPConnection(self.target+":"+self.port)
        conn_scalepod.request('PUT', '/api/'+self.apiversion+'/namespaces/'+namespace+'/' +
                              'replicationcontrollers/'+name, headers=headers, body=json.dumps(rc))
        resp = conn_scalepod.getresponse()
        data = resp.read()
        reason = resp.reason
        status = resp.status
        conn_scalepod.close()
        if not 200 <= status <= 299:
            errmsg = "Failed to scale Replication Controller:{} {} {} - {}".format(
                name, status, reason, data)
            raise RuntimeError(errmsg)
        resource_ver = rc['metadata']['resourceVersion']
        for _ in xrange(30):
            js_template = self._get_rc_(name, namespace)
            if js_template["metadata"]["resourceVersion"] != resource_ver:
                break
            time.sleep(1)
        self._get_schedule_status(name, num, namespace)
        for _ in xrange(120):
            count = 0
            status, data, reason = self._get_pods(namespace)
            parsed_json = json.loads(data)
            for pod in parsed_json['items']:
                if(pod['metadata']['generateName'] == name+'-' and
                   pod['status']['phase'] == 'Running'):
                    count += 1
            if count == num:
                break
            time.sleep(1)

    def _scale_app(self, name, num, namespace):
        js_template = self._get_rc_(name, namespace)
        js_template["spec"]["replicas"] = num
        self._scale_rc(js_template, namespace)

    def scale(self, name, image, command, **kwargs):
        app_name = kwargs.get('aname', {})
        rc_name = name.replace(".", "-")
        rc_name = rc_name.replace("_", "-")
        if not 200 <= self._get_rc_status(rc_name, app_name) <= 299:
            self.create(name, image, command, **kwargs)
            return
        name = name.replace(".", "-")
        name = name.replace("_", "-")
        num = kwargs.get('num', {})
        js_template = self._get_rc_(name, app_name)
        old_replicas = js_template["spec"]["replicas"]
        try:
            self._scale_app(name, num, app_name)
        except Exception as e:
            self._scale_app(name, old_replicas, app_name)
            err = '{} (Scale): {}'.format(name, e)
            raise RuntimeError(err)

    def _create_rc(self, name, image, command, **kwargs):
        container_fullname = name
        app_name = kwargs.get('aname', {})
        app_type = name.split(".")[1]
        container_name = app_name+"-"+app_type
        name = name.replace(".", "-")
        name = name.replace("_", "-")
        args = command.split()

        num = kwargs.get('num', {})
        l = {}
        l["name"] = name
        l["id"] = app_name
        l["appversion"] = kwargs.get('version', {})
        l["version"] = self.apiversion
        l["image"] = self.registry+"/"+image
        l['num'] = num
        l['containername'] = container_name
        l['type'] = app_type
        template = string.Template(RC_TEMPLATE).substitute(l)
        js_template = json.loads(template)
        containers = js_template["spec"]["template"]["spec"]["containers"]
        containers[0]['args'] = args
        loc = locals().copy()
        loc.update(re.match(MATCH, container_fullname).groupdict())
        mem = kwargs.get('memory', {}).get(loc['c_type'])
        cpu = kwargs.get('cpu', {}).get(loc['c_type'])
        if mem or cpu:
            containers[0]["resources"] = {"limits": {}}
        if mem:
            if mem[-2:-1].isalpha() and mem[-1].isalpha():
                mem = mem[:-1]
            mem = mem+"i"
            containers[0]["resources"]["limits"]["memory"] = mem
        if cpu:
            cpu = float(cpu)/1024
            containers[0]["resources"]["limits"]["cpu"] = cpu
        headers = {'Content-Type': 'application/json'}
        conn_rc = httplib.HTTPConnection(self.target+":"+self.port)
        conn_rc.request('POST', '/api/'+self.apiversion+'/namespaces/'+app_name+'/' +
                        'replicationcontrollers', headers=headers, body=json.dumps(js_template))
        resp = conn_rc.getresponse()
        data = resp.read()
        reason = resp.reason
        status = resp.status
        conn_rc.close()
        if not 200 <= status <= 299:
            errmsg = "Failed to create Replication Controller:{} {} {} - {}".format(
                name, status, reason, data)
            raise RuntimeError(errmsg)
        create = False
        for _ in xrange(30):
            if not create and self._get_rc_status(name, app_name) == 404:
                time.sleep(1)
                continue
            create = True
            rc = self._get_rc_(name, app_name)
            if ("observedGeneration" in rc["status"]
                    and rc["metadata"]["generation"] == rc["status"]["observedGeneration"]):
                break
            time.sleep(1)
        return json.loads(data)

    def create(self, name, image, command, **kwargs):
        """Create a container."""
        self._create_rc(name, image, command, **kwargs)
        app_type = name.split(".")[1]
        name = name.replace(".", "-")
        name = name.replace("_", "-")
        app_name = kwargs.get('aname', {})
        try:
            self._create_service(name, app_name, app_type)
        except Exception as e:
            self._scale_app(name, 0, app_name)
            self._delete_rc(name, app_name)
            err = '{} (create): {}'.format(name, e)
            raise RuntimeError(err)

    def _get_service(self, name, namespace):
        con_get = httplib.HTTPConnection(self.target+":"+self.port)
        con_get.request('GET', '/api/'+self.apiversion+'/namespaces/'+namespace+'/services/'+name)
        resp = con_get.getresponse()
        reason = resp.reason
        status = resp.status
        data = resp.read()
        con_get.close()
        if not 200 <= status <= 299:
            errmsg = "Failed to get Service: {} {} - {}".format(
                status, reason, data)
            raise RuntimeError(errmsg)
        return (status, data, reason)

    def _create_service(self, name, app_name, app_type):
        random.seed(app_name)
        app_id = random.randint(1, 100000)
        appname = "app-"+str(app_id)
        actual_pod = {}
        for _ in xrange(300):
            status, data, reason = self._get_pods(app_name)
            parsed_json = json.loads(data)
            for pod in parsed_json['items']:
                if('generateName' in pod['metadata'] and
                   pod['metadata']['generateName'] == name+'-'):
                    actual_pod = pod
                    break
            if actual_pod and actual_pod['status']['phase'] == 'Running':
                break
            time.sleep(1)
        container_id = actual_pod['status']['containerStatuses'][0]['containerID'].split("//")[1]
        ip = actual_pod['status']['hostIP']
        docker_cli = Client("tcp://{}:2375".format(ip), timeout=1200, version='1.17')
        container = docker_cli.inspect_container(container_id)
        port = int(container['Config']['ExposedPorts'].keys()[0].split("/")[0])
        l = {}
        l["version"] = self.apiversion
        l["label"] = app_name
        l["port"] = port
        l['type'] = app_type
        l["name"] = appname
        template = string.Template(SERVICE_TEMPLATE).substitute(l)
        headers = {'Content-Type': 'application/json'}
        conn_serv = httplib.HTTPConnection(self.target+":"+self.port)
        conn_serv.request('POST', '/api/'+self.apiversion+'/namespaces/'+app_name+'/services',
                          headers=headers, body=copy.deepcopy(template))
        resp = conn_serv.getresponse()
        data = resp.read()
        reason = resp.reason
        status = resp.status
        conn_serv.close()
        if status == 409:
            status, data, reason = self._get_service(appname, app_name)
            srv = json.loads(data)
            if srv['spec']['selector']['type'] == 'web':
                return
            srv['spec']['selector']['type'] = app_type
            srv['spec']['ports'][0]['targetPort'] = port
            headers = {'Content-Type': 'application/json'}
            conn_scalepod = httplib.HTTPConnection(self.target+":"+self.port)
            conn_scalepod.request('PUT', '/api/'+self.apiversion+'/namespaces/'+app_name+'/' +
                                  'services/'+appname, headers=headers, body=json.dumps(srv))
            resp = conn_scalepod.getresponse()
            data = resp.read()
            reason = resp.reason
            status = resp.status
            conn_scalepod.close()
            if not 200 <= status <= 299:
                errmsg = "Failed to update the Service:{} {} {} - {}".format(
                    name, status, reason, data)
                raise RuntimeError(errmsg)
        elif not 200 <= status <= 299:
            errmsg = "Failed to create Service:{} {} {} - {}".format(
                     name, status, reason, data)
            raise RuntimeError(errmsg)

    def start(self, name):
        """Start a container."""
        pass

    def stop(self, name):
        """Stop a container."""
        pass

    def _delete_rc(self, name, namespace):
        headers = {'Content-Type': 'application/json'}
        con_dest = httplib.HTTPConnection(self.target+":"+self.port)
        con_dest.request('DELETE', '/api/'+self.apiversion+'/namespaces/'+namespace+'/' +
                         'replicationcontrollers/'+name, headers=headers, body=POD_DELETE)
        resp = con_dest.getresponse()
        reason = resp.reason
        status = resp.status
        data = resp.read()
        con_dest.close()
        if not 200 <= status <= 299:
            errmsg = "Failed to delete Replication Controller:{} {} {} - {}".format(
                name, status, reason, data)
            raise RuntimeError(errmsg)

    def destroy(self, name):
        """Destroy a container."""
        appname = name.split("_")[0]
        name = name.split(".")
        name = name[0]+'-'+name[1]
        name = name.replace("_", "-")

        headers = {'Content-Type': 'application/json'}
        con_dest = httplib.HTTPConnection(self.target+":"+self.port)
        con_dest.request('DELETE', '/api/'+self.apiversion+'/namespaces/'+appname+'/' +
                         'replicationcontrollers/'+name, headers=headers, body=POD_DELETE)
        resp = con_dest.getresponse()
        reason = resp.reason
        status = resp.status
        data = resp.read()
        con_dest.close()
        if status == 404:
            return
        if not 200 <= status <= 299:
            errmsg = "Failed to delete Replication Controller:{} {} {} - {}".format(
                name, status, reason, data)
            raise RuntimeError(errmsg)

        random.seed(appname)
        app_id = random.randint(1, 100000)
        app_name = "app-"+str(app_id)
        con_serv = httplib.HTTPConnection(self.target+":"+self.port)
        con_serv.request('DELETE', '/api/'+self.apiversion +
                         '/namespaces/'+appname+'/services/'+app_name)
        resp = con_serv.getresponse()
        reason = resp.reason
        status = resp.status
        data = resp.read()
        con_serv.close()
        if status != 404 and not 200 <= status <= 299:
            errmsg = "Failed to delete service:{} {} {} - {}".format(
                name, status, reason, data)
            raise RuntimeError(errmsg)

        status, data, reason = self._get_pods(appname)
        parsed_json = json.loads(data)
        for pod in parsed_json['items']:
            if 'generateName' in pod['metadata'] and pod['metadata']['generateName'] == name+'-':
                self._delete_pod(pod['metadata']['name'], appname)
        con_ns = httplib.HTTPConnection(self.target+":"+self.port)
        con_ns.request('DELETE', '/api/'+self.apiversion+'/namespaces/'+appname)
        resp = con_ns.getresponse()
        reason = resp.reason
        status = resp.status
        data = resp.read()
        con_ns.close()
        if not 200 <= status <= 299:
            errmsg = "Failed to delete namespace:{} {} {} - {}".format(
                appname, status, reason, data)
            raise RuntimeError(errmsg)

    def _get_pod(self, name, namespace):
        conn_pod = httplib.HTTPConnection(self.target+":"+self.port)
        conn_pod.request('GET', '/api/'+self.apiversion+'/namespaces/'+namespace+'/pods/'+name)
        resp = conn_pod.getresponse()
        status = resp.status
        data = resp.read()
        reason = resp.reason
        conn_pod.close()
        return (status, data, reason)

    def _get_pods(self, namespace):
        con_get = httplib.HTTPConnection(self.target+":"+self.port)
        con_get.request('GET', '/api/'+self.apiversion+'/namespaces/'+namespace+'/pods')
        resp = con_get.getresponse()
        reason = resp.reason
        status = resp.status
        data = resp.read()
        con_get.close()
        if not 200 <= status <= 299:
            errmsg = "Failed to get Pods: {} {} - {}".format(
                status, reason, data)
            raise RuntimeError(errmsg)
        return (status, data, reason)

    def _delete_pod(self, name, namespace):
        headers = {'Content-Type': 'application/json'}
        con_dest_pod = httplib.HTTPConnection(self.target+":"+self.port)
        con_dest_pod.request('DELETE', '/api/'+self.apiversion+'/namespaces/' +
                             namespace+'/pods/'+name, headers=headers, body=POD_DELETE)
        resp = con_dest_pod.getresponse()
        reason = resp.reason
        status = resp.status
        data = resp.read()
        con_dest_pod.close()
        if not 200 <= status <= 299:
            errmsg = "Failed to delete Pod: {} {} - {}".format(
                status, reason, data)
            raise RuntimeError(errmsg)
        for _ in xrange(5):
            status, data, reason = self._get_pod(name, namespace)
            if status != 404:
                time.sleep(1)
                continue
            break
        if status != 404:
            errmsg = "Failed to delete Pod: {} {} - {}".format(
                status, reason, data)
            raise RuntimeError(errmsg)

    def _pod_log(self, name, namespace):
        conn_log = httplib.HTTPConnection(self.target+":"+self.port)
        conn_log.request('GET', '/api/'+self.apiversion+'/namespaces/' +
                         namespace+'/pods/'+name+'/log')
        resp = conn_log.getresponse()
        status = resp.status
        data = resp.read()
        reason = resp.reason
        conn_log.close()
        if not 200 <= status <= 299:
            errmsg = "Failed to get the log: {} {} - {}".format(
                status, reason, data)
            raise RuntimeError(errmsg)
        return (status, data, reason)

    def logs(self, name):
        appname = name.split("_")[0]
        name = name.replace(".", "-")
        name = name.replace("_", "-")
        status, data, reason = self._get_pods(appname)
        parsed_json = json.loads(data)
        log_data = ''
        for pod in parsed_json['items']:
            if name in pod['metadata']['generateName'] and pod['status']['phase'] == 'Running':
                status, data, reason = self._pod_log(pod['metadata']['name'], appname)
                log_data += data
        return log_data

    def run(self, name, image, entrypoint, command):
        """Run a one-off command."""
        appname = name.split("_")[0]
        name = name.replace(".", "-")
        name = name.replace("_", "-")
        l = {}
        l["id"] = name
        l["version"] = self.apiversion
        l["image"] = self.registry+"/"+image
        template = string.Template(POD_TEMPLATE).substitute(l)
        if command.startswith("-c "):
            args = command.split(' ', 1)
            args[1] = args[1][1:-1]
        else:
            args = [command[1:-1]]
        js_template = json.loads(template)
        js_template['spec']['containers'][0]['command'] = [entrypoint]
        js_template['spec']['containers'][0]['args'] = args

        con_dest = httplib.HTTPConnection(self.target+":"+self.port)
        headers = {'Content-Type': 'application/json'}
        con_dest.request('POST', '/api/'+self.apiversion+'/namespaces/'+appname+'/pods',
                         headers=headers, body=json.dumps(js_template))
        resp = con_dest.getresponse()
        data = resp.read()
        status = resp.status
        reason = resp.reason
        con_dest.close()
        if not 200 <= status <= 299:
            errmsg = "Failed to create a Pod: {} {} - {}".format(
                status, reason, data)
            raise RuntimeError(errmsg)
        while(1):
            parsed_json = {}
            status = 404
            reason = ''
            data = ''
            for _ in xrange(5):
                status, data, reason = self._get_pod(name, appname)
                if not 200 <= status <= 299:
                    time.sleep(1)
                    continue
                parsed_json = json.loads(data)
                break
            if not 200 <= status <= 299:
                errmsg = "Failed to create a Pod: {} {} - {}".format(
                    status, reason, data)
                raise RuntimeError(errmsg)
            if parsed_json['status']['phase'] == 'Succeeded':
                status, data, reason = self._pod_log(name, appname)
                self._delete_pod(name, appname)
                return 0, data
            elif parsed_json['status']['phase'] == 'Failed':
                pod_state = parsed_json['status']['containerStatuses'][0]['state']
                err_code = pod_state['terminated']['exitCode']
                self._delete_pod(name, appname)
                return err_code, data
            time.sleep(1)
        return 0, data

    def _get_pod_state(self, name):
        try:
            appname = name.split("_")[0]
            name = name.split(".")
            name = name[0]+'-'+name[1]
            name = name.replace("_", "-")
            for _ in xrange(120):
                status, data, reason = self._get_pods(appname)
                parsed_json = json.loads(data)
                for pod in parsed_json['items']:
                    if pod['metadata']['generateName'] == name+'-':
                        actual_pod = pod
                        break
                if actual_pod and actual_pod['status']['phase'] == 'Running':
                    return JobState.up
                time.sleep(1)
            return JobState.destroyed
        except:
            return JobState.destroyed

    def state(self, name):
        """Display the given job's running state."""
        try:
            return self._get_pod_state(name)
        except KeyError:
            return JobState.error
        except RuntimeError:
            return JobState.destroyed

SchedulerClient = KubeHTTPClient
