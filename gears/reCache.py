from elasticsearch import Elasticsearch
import json

# Listens on key expirations

def fetch_and_load(sha, index):
    try:
        rec = execute("HMGET", "es-proxy:recache:{}:{}".format(index, sha), "QUERY", "TTL")
        execute("HSET", "LOG", "QUERY", rec[0], "TTL", rec[1])
        es = Elasticsearch(hosts="es")
        res = es.search(index=index, body=rec[0])
        execute("SETEX", "es-proxy:query:{}:{}".format(index, sha), rec[1], json.dumps(res))
        execute("INCR", "es-proxy:internal_stats:recache_success")
    except Exception as res:
        execute("INCR", "es-proxy:internal_stats:recache_errors")
        execute('SET', "LAST_ERROR", res)

def runIt(x):
    if x['event'] == "expired":
        info = x['key'].split(":")
        res = fetch_and_load(info[3], info[2])

gb = GearsBuilder(
    reader='KeysReader',
    defaultArg='es-proxy:query:*',
    desc="Automatically add and remove all users from a set")

gb.map(runIt)  
gb.register('es-proxy:query:*')
