import cStringIO
import base64
import copy
import json
import httplib
import time
import re
import string
import os
from django.conf import settings
from .states import JobState
from docker import Client
import etcd

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
         "version":"$appversion"
      },
      "template":{
         "metadata":{
            "labels":{
               "name":"$id",
               "version":"$appversion"
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
      "name":"$label",
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
         "name":"$label"
      }
   }
}'''

POD_DELETE = '''{
}'''


RETRIES = 3
MATCH = re.compile(
    r'(?P<app>[a-z0-9-]+)_?(?P<version>v[0-9]+)?\.?(?P<c_type>[a-z-_]+)')

class KubeHTTPClient():

    def __init__(self, target, auth, options, pkey):
        self.target = settings.K8S_MASTER
        self.port = "8080"
        self.registry = settings.REGISTRY_HOST+":"+settings.REGISTRY_PORT
        self.apiversion = "v1"
        self.conn = httplib.HTTPConnection(self.target+":"+self.port)
        #self.container_state = ""

    def _get_old_rc(self,name):
        con_app = httplib.HTTPConnection(self.target+":"+self.port)
        con_app.request('GET','/api/'+self.apiversion+'/namespaces/default/replicationcontrollers')
        resp = con_app.getresponse()
        data = resp.read()
        reason = resp.reason
        status = resp.status
        con_app.close()
        if not 200 <= status <= 299:
            errmsg = "Failed to get Replication Controllers: {} {} - {}".format(
                status, reason, data)
            raise RuntimeError(errmsg)
        parsed_json =  json.loads(data)
        exists = False
        prev_rc = []
        for rc in parsed_json['items']:
            if 'name' in rc['metadata']['labels'] and name == rc['metadata']['labels']['name']:
                exists = True
                prev_rc = rc
                break
        if exists :
            return prev_rc
        else:
            return 0


    def _get_rc_status(self,name):
      conn_rc = httplib.HTTPConnection(self.target+":"+self.port)
      conn_rc.request('GET','/api/'+self.apiversion+'/'+'namespaces/default/replicationcontrollers/'+name)
      resp = conn_rc.getresponse()
      status = resp.status
      conn_rc.close()
      return status

    def _get_rc_(self,name):
      conn_rc_resver = httplib.HTTPConnection(self.target+":"+self.port)
      conn_rc_resver.request('GET','/api/'+self.apiversion+'/'+'namespaces/default/replicationcontrollers/'+name)
      resp = conn_rc_resver.getresponse()
      data = resp.read()
      reason = resp.reason
      status = resp.status
      conn_rc_resver.close()
      if not 200 <= status <= 299:
          errmsg = "Failed to get Replication Controller:{} {} {} - {}".format(
              name,status, reason, data)
          raise RuntimeError(errmsg)
      parsed_json =  json.loads(data)
      return parsed_json

    def deploy(self, name, image, command, **kwargs):
        app_name = kwargs.get('aname', {})
        old_rc = self._get_old_rc(app_name)
        new_rc = self._create_rc(name, image, command, **kwargs)
        desired = int(old_rc["spec"]["replicas"])
        old_rc_name = old_rc["metadata"]["name"]
        new_rc_name = new_rc["metadata"]["name"]
        try:
            count = 1
            while desired >= count:
                new_rc = self._scale_app(new_rc_name,count)
                old_rc = self._scale_app(old_rc_name,desired-count)
                count += 1
        except Exception as e:
            if count > 1:
                self._scale_app(new_rc["metadata"]["name"],0)
            self._delete_rc(new_rc["metadata"]["name"])
            self._scale_app(old_rc["metadata"]["name"],desired)
            err = '{} (deploy): {}'.format(name, e)
            raise RuntimeError(err)
        self._delete_rc(old_rc_name)

    def _scale_rc(self,rc):
        name= rc['metadata']['name']
        num = rc["spec"]["replicas"]
        headers = {'Content-Type': 'application/json'}
        conn_scalepod = httplib.HTTPConnection(self.target+":"+self.port)
        conn_scalepod.request('PUT', '/api/'+self.apiversion+'/namespaces/default/replicationcontrollers/'+name,
                          headers=headers,body=json.dumps(rc))
        resp = conn_scalepod.getresponse()
        data = resp.read()
        reason = resp.reason
        status = resp.status
        conn_scalepod.close()
        if not 200 <= status <= 299:
            errmsg = "Failed to scale Replication Controller:{} {} {} - {}".format(
                name,rc, reason, data)
            raise RuntimeError(errmsg)
        for _ in xrange(120):
            count = 0
            status,data,reason = self._get_pods()
            parsed_json =  json.loads(data)
            for pod in parsed_json['items']:
                if pod['metadata']['generateName'] == name+'-' and pod['status']['phase'] == 'Running':
                    count += 1
            if count == num:
                break
            time.sleep(1)

    def _scale_app(self,name,num):
        js_template = self._get_rc_(name)
        js_template["spec"]["replicas"]= num
        self._scale_rc(js_template)

    def scale(self, name, image, command, **kwargs):
      rc_name = name.split(".")[0]
      rc_name = rc_name.replace("_","-")
      if not 200 <= self._get_rc_status(rc_name) <= 299 :
          self.create(name, image, command, **kwargs)
          return
      name = name.split(".")[0]
      name = name.replace("_","-")
      num = kwargs.get('num', {})
      self._scale_app(name,num)

    def _create_rc(self, name, image, command, **kwargs):
        container_fullname = name
        app_name = kwargs.get('aname', {})
        app_type = name.split(".")[1]
        container_name = app_name+"-"+app_type
        name = name.split(".")[0]
        name = name.replace("_","-")
        args= command.split()

        num = kwargs.get('num', {})
        l = {}
        l["name"]= name
        l["id"]= app_name
        l["appversion"] = kwargs.get('version', {})
        l["version"]=self.apiversion
        l["image"]=self.registry+"/"+image
        l['num'] =  num
        l['containername'] = container_name
        template=string.Template(RC_TEMPLATE).substitute(l)
        js_template = json.loads(template)
        js_template["spec"]["template"]["spec"]["containers"][0]['args'] = args
        loc = locals().copy()
        loc.update(re.match(MATCH, container_fullname).groupdict())
        mem = kwargs.get('memory', {}).get(loc['c_type'])
        cpu = kwargs.get('cpu', {}).get(loc['c_type'])
        if mem or cpu :
            js_template["spec"]["template"]["spec"]["containers"][0]["resources"] = {"limits":{}}
        if mem:
            if mem[-2:-1].isalpha() and mem[-1].isalpha():
                mem = mem[:-1]
            mem = mem+"i"
            js_template["spec"]["template"]["spec"]["containers"][0]["resources"]["limits"]["memory"] = mem
        if cpu:
            cpu = float(cpu)/1024
            js_template["spec"]["template"]["spec"]["containers"][0]["resources"]["limits"]["cpu"] = cpu
        headers = {'Content-Type': 'application/json'}
        conn_rc = httplib.HTTPConnection(self.target+":"+self.port)
        conn_rc.request('POST', '/api/'+self.apiversion+'/namespaces/default/replicationcontrollers',
                  headers=headers, body=json.dumps(js_template))
        resp = conn_rc.getresponse()
        data = resp.read()
        reason = resp.reason
        status = resp.status
        conn_rc.close()
        if not 200 <= status <= 299:
            errmsg = "Failed to create Replication Controller:{} {} {} - {}".format(
                name,status, reason, data)
            raise RuntimeError(errmsg)
        create = False
        for _ in xrange(30):
            if not create and self._get_rc_status(name) == 404 :
                time.sleep(1)
                continue
            create = True
            rc = self._get_rc_(name)
            if "observedGeneration" in rc["status"] and rc["metadata"]["generation"] == rc["status"]["observedGeneration"]:
                break
            time.sleep(1)
        return json.loads(data)

    def create(self, name, image, command, **kwargs):
        #self.container_state = "create"
        self._create_rc(name, image, command, **kwargs)
        name = name.split(".")[0]
        name = name.replace("_","-")
        app_name = kwargs.get('aname', {})
        self._create_service(name , app_name)

    def _create_service(self,name,app_name):
      actual_pod = {}
      for _ in xrange(300):
          status,data,reason = self._get_pods()
          parsed_json =  json.loads(data)
          for pod in parsed_json['items']:
              if 'generateName' in pod['metadata'] and pod['metadata']['generateName'] == name+'-':
                  actual_pod = pod
                  break
          if actual_pod and actual_pod['status']['phase'] == 'Running':
              break
          time.sleep(1)

      container_id = actual_pod['status']['containerStatuses'][0]['containerID'].split("//")[1]
      ip = actual_pod['status']['hostIP']
      docker_cli = Client("tcp://{}:2375".format(ip),timeout=1200, version='1.17')
      port = int(docker_cli.inspect_container(container_id)['Config']['ExposedPorts'].keys()[0].split("/")[0])
      l = {}
      l["version"]=self.apiversion
      l["label"] =app_name
      l["port"] = port
      template=string.Template(SERVICE_TEMPLATE).substitute(l)
      headers = {'Content-Type': 'application/json'}
      conn_serv = httplib.HTTPConnection(self.target+":"+self.port)
      conn_serv.request('POST', '/api/'+self.apiversion+'/namespaces/default/services',
                  headers=headers, body=copy.deepcopy(template))
      resp = conn_serv.getresponse()
      data = resp.read()
      reason = resp.reason
      status = resp.status
      conn_serv.close()
      if not 200 <= status <= 299:
          errmsg = "Failed to create Service:{} {} {} - {}".format(
              name,status, reason, data)
          raise RuntimeError(errmsg)

    def start(self, name):
        """
        Start a container
        """
        #self.container_state = "start"
        return

    def stop(self, name):
        """
        Stop a container
        """
        return

    def _delete_rc(self,name):
        headers = {'Content-Type': 'application/json'}
        con_dest = httplib.HTTPConnection(self.target+":"+self.port)
        con_dest.request('DELETE','/api/'+self.apiversion+'/namespaces/default/replicationcontrollers/'+name,headers=headers,body=POD_DELETE)
        resp = con_dest.getresponse()
        reason = resp.reason
        status = resp.status
        data = resp.read()
        con_dest.close()
        if not 200 <= status <= 299:
            errmsg = "Failed to delete Replication Controller:{} {} {} - {}".format(
                name,status, reason, data)
            raise RuntimeError(errmsg)

    def destroy(self, name):
        """
        Destroy a container
        """
        name = name.split(".")[0]
        name = name.replace("_","-")

        appname = self._get_rc_(name)["metadata"]["labels"]["name"]

        headers = {'Content-Type': 'application/json'}
        con_dest = httplib.HTTPConnection(self.target+":"+self.port)
        con_dest.request('DELETE','/api/'+self.apiversion+'/namespaces/default/replicationcontrollers/'+name,headers=headers,body=POD_DELETE)
        resp = con_dest.getresponse()
        reason = resp.reason
        status = resp.status
        data = resp.read()
        con_dest.close()
        if status == 404:
            return
        if not 200 <= status <= 299:
            errmsg = "Failed to delete Replication Controller:{} {} {} - {}".format(
                name,status, reason, data)
            raise RuntimeError(errmsg)

        con_serv = httplib.HTTPConnection(self.target+":"+self.port)
        con_serv.request('DELETE','/api/'+self.apiversion+'/namespaces/default/services/'+appname)
        resp = con_serv.getresponse()
        reason = resp.reason
        status = resp.status
        data = resp.read()
        con_serv.close()
        if not 200 <= status <= 299:
            errmsg = "Failed to delete service:{} {} {} - {}".format(
                name,status, reason, data)
            raise RuntimeError(errmsg)

        status,data,reason = self._get_pods()
        parsed_json =  json.loads(data)
        for pod in parsed_json['items']:
            if 'generateName' in pod['metadata'] and pod['metadata']['generateName'] == name+'-':
                self._delete_pod(pod['metadata']['name'])


    def _get_pod(self,name):
        conn_pod = httplib.HTTPConnection(self.target+":"+self.port)
        conn_pod.request('GET','/api/'+self.apiversion+'/namespaces/default/pods/'+name)
        resp = conn_pod.getresponse()
        status = resp.status
        data =resp.read()
        reason = resp.reason
        conn_pod.close()
        return (status,data,reason)

    def _get_pods(self):
        con_get = httplib.HTTPConnection(self.target+":"+self.port)
        con_get.request('GET','/api/'+self.apiversion+'/namespaces/default/pods')
        resp = con_get.getresponse()
        reason = resp.reason
        status = resp.status
        data = resp.read()
        con_get.close()
        if not 200 <= status <= 299:
            errmsg = "Failed to get Pods: {} {} - {}".format(
                status, reason, data)
            raise RuntimeError(errmsg)
        return (status,data,reason)

    def _delete_pod(self,name):
        headers = {'Content-Type': 'application/json'}
        con_dest_pod = httplib.HTTPConnection(self.target+":"+self.port)
        con_dest_pod.request('DELETE','/api/'+self.apiversion+'/namespaces/default/pods/'+name,headers=headers,body=POD_DELETE)
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
            status,data,reason = self._get_pod(name)
            if status != 404:
                time.sleep(1)
                continue
            break
        if status != 404 :
            errmsg = "Failed to delete Pod: {} {} - {}".format(
                status, reason, data)
            raise RuntimeError(errmsg)

    def _pod_log(self,name):
        headers = {'Content-Type': 'application/json'}
        conn_log = httplib.HTTPConnection(self.target+":"+self.port)
        conn_log.request('GET', '/api/'+self.apiversion+'/namespaces/default/pods/'+name+'/log')
        resp = conn_log.getresponse()
        status = resp.status
        data = resp.read()
        reason = resp.reason
        conn_log.close()
        if not 200 <= status <= 299:
            errmsg = "Failed to get the log: {} {} - {}".format(
                status, reason, data)
            raise RuntimeError(errmsg)
        return (status,data,reason)

    def logs(self,name):
        name = name.split(".")[0]
        name = name.replace("_","-")
        status,data,reason = self._get_pods()
        parsed_json =  json.loads(data)
        log_data = ''
        for pod in parsed_json['items']:
            if name in pod['metadata']['generateName'] and pod['status']['phase'] == 'Running':
                status,data,reason = self._pod_log(pod['metadata']['name'])
                log_data += data
        return log_data

    def run(self, name, image, entrypoint, command):
        """
        Run a one-off command
        """
        name = name.split(".")[0]
        name = name.replace("_","-")
        l = {}
        l["id"]=name
        l["version"]=self.apiversion
        l["image"]=self.registry+"/"+image
        template=string.Template(POD_TEMPLATE).substitute(l)
        if command.startswith("-c "):
            args = command.split(' ' ,1)
            args[1] = args[1][1:-1]
        else :
            args = [command]
        js_template = json.loads(template)
        js_template['spec']['containers'][0]['command'] = [entrypoint]
        js_template['spec']['containers'][0]['args'] = args

        con_dest = httplib.HTTPConnection(self.target+":"+self.port)
        headers = {'Content-Type': 'application/json'}
        con_dest.request('POST', '/api/'+self.apiversion+'/namespaces/default/pods',
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
                status,data,reason = self._get_pod(name)
                if not 200 <= status <= 299:
                    time.sleep(1)
                    continue
                parsed_json =  json.loads(data)
                break
            if not 200 <= status <= 299:
                errmsg = "Failed to create a Pod: {} {} - {}".format(
                    status, reason, data)
                raise RuntimeError(errmsg)
            if parsed_json['status']['phase'] == 'Succeeded':
                status,data,reason = self._pod_log(name)
                self._delete_pod(name)
                return 0, data
            elif parsed_json['status']['phase'] == 'Failed':
                err_code = parsed_json['status']['containerStatuses'][0]['state']['terminated']['exitCode']
                self._delete_pod(name)
                return err_code,data
            time.sleep(1)
        return 0, data

    def _get_pod_state(self, name):
        try:
            name = name.split(".")[0]
            name = name.replace("_","-")
            for _ in xrange(120):
                status,data,reason = self._get_pods()
                parsed_json =  json.loads(data)
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
        try:
            return self._get_pod_state(name)
        except KeyError:
            return JobState.error
        except RuntimeError:
            return JobState.destroyed

    def attach(self, name):
        """
        Attach to a job's stdin, stdout and stderr
        """
        return StringIO(), StringIO(), StringIO()

SchedulerClient = KubeHTTPClient
