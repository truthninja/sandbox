import os
import time
import json
import subprocess
import multiprocessing

POD_MANIFEST = {"desiredState": {"manifest": {"version": "v1beta1", "id": "php", "containers": [{"image": "dockerfile/nginx", "name": "nginx", "livenessProbe": {"type": "http", "initialDelaySeconds": 30, "enabled": True, "httpGet": {"path": "/index.html", "port": "8081"}}, "ports": [{"containerPort": 80, "hostPort": 8081}]}]}}, "kind": "Pod", "labels": {"name": "foo"}, "id": "php", "apiVersion": "v1beta1"}

def get(state=None):
    cmd = './cluster/kubectl.sh get pods'
    if state:
        search = '| grep -i %s | wc -l' % state
        cmd = '%s %s' % (cmd, search)
    p = subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE, stdin=subprocess.PIPE)
    return p.communicate()[0]


def kubectl(start_id, finish_id, create=True):
    action = 'delete'
    if create:
        action = 'create'
    for i in range(start_id, finish_id):
        str_json = json.dumps(POD_MANIFEST)
        str_json = str_json.replace('8081', str(i)).replace('foo', 'foo%s'%i).replace('php', 'php%s'%i)
        p = subprocess.Popen(['./cluster/kubectl.sh', action, '-f', '-'], stdout=subprocess.PIPE, stdin=subprocess.PIPE)
        p.stdin.write(str_json)
        p.communicate()[0]

def pod_cycle(start_port, end_port, wait_limit):
    print 'creating pods '
    kubectl(start_port, end_port, create=True)
    running_pods = '0'
    while int(running_pods) < wait_limit:
        running_pods=get('Running')
        print 'Running pods %s' % running_pods
        time.sleep(1)

    print 'deleting pods '
    kubectl(start_port, end_port, create=False)
    while int(running_pods) > wait_limit:
        running_pods=get('Running')
        print 'Running pods %s' % running_pods
        time.sleep(1)

def pod_stress():
    proc_list = []
    for i in range(1, 10):
        port_base = 8000+i*10
        p = multiprocessing.Process(target=pod_cycle, args=(port_base, port_base+10, 5))
        p.start()
        proc_list.append(p)

    for p in proc_list:
        p.join()


pod_stress()


# pod_stress(8011, 8021, wait_limit=5)
# pod_stress(8022, 8023, wait_limit=5)
# pod_stress(8024, 8034, wait_limit=5)
