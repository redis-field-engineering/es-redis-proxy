#!/usr/bin/env python

from jsondiff import diff
import requests

# don't do this in production
import warnings
warnings.filterwarnings("ignore")

import json

from elasticsearch import Elasticsearch

for x in ['{"query": {"match_all": {}}}','{"query":{"query_string":{"query":"ZINC OR GOOG","default_field":"symbol"}}}']:

    # Query elasticsearch
    es = Elasticsearch()
    esres = es.search(index="instruments", body=x)
    
    
    # Query the proxy
    proxy = Elasticsearch([
        {'host': 'localhost', 'port': 8080},
    ])
    proxyres = proxy.search(index="instruments", body=x)
    
    d = diff(esres, proxyres)

    # delet the took as it's going to be different with the cache version
    if 'took' in d: del d['took']

    print("Query: ", x)
    
    if len(d) > 0:
        print("\t", d)
    else:
        print("\tNo differences")
    

    { "query": { "query_string": { "query": "MSBHF OR GOOG", "default_field": "instrument" } } }

print("Testing to make sure it is cached")

r = requests.post('http://localhost:8080/instruments/_search', data='{"query":{"query_string":{"query":"ZSPH","default_field":"symbol"}}}')
print(r.headers)
r = requests.post('http://localhost:8080/instruments/_search', data='{"query":{"query_string":{"query":"ZSPH","default_field":"symbol"}}}')
print(r.headers)

