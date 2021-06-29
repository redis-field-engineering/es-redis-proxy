from elasticsearch import Elasticsearch
import json

# Takes
# [ SHA, Index, Query, TTL ]

# Returns JSON

def fetch_and_load(sha, index, query, ttl):
    try:
        es = Elasticsearch(hosts="es")
        res = es.search(index=index, body=query)
        execute("SETEX", "{}:{}".format(index, sha), ttl, json.dumps(res))
        return(res, 0)
    except Exception as res:
        return(res.error, 1)

def runIt(x):
    w = execute("GET", "{}:{}".format(x[2],x[1]))
    if w :
        t = execute("TTL", "{}:{}".format(x[2],x[1]))
        j = json.loads(w)
        out = {"result": j, "ttl": t, "exit_code": 0}
    else:
        res, exit_code = fetch_and_load(x[1], x[2], x[3], x[4])
        out = {"result": res, "ttl": x[4], "exit_code": exit_code}

    return(json.dumps(out))

gb = GB('CommandReader', desc="Query Upstream ES if not already in cache")
gb.map(runIt)
gb.register(trigger='es-search')

