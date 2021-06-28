#!/usr/bin/env python

from jsondiff import diff

# don't do this in production
import warnings
warnings.filterwarnings("ignore")

import json

from elasticsearch import Elasticsearch

for x in ['{"query": {"match_all": {}}}','{"query":{"query_string":{"query":"MSBHF OR GOOG","default_field":"instrument"}}}']:

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

