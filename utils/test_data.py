#!/usr/bin/env python

from jsondiff import diff

# don't do this in production
import warnings
warnings.filterwarnings("ignore")

import json

from elasticsearch import Elasticsearch

es = Elasticsearch()
esres = es.search(index="instruments", body={"query": {"match_all": {}}})


proxy = Elasticsearch([
    {'host': 'localhost', 'port': 8080},
])
proxyres = proxy.search(index="instruments", body={"query": {"match_all": {}}})


print(diff(esres, proxyres))