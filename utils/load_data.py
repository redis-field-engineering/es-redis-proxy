#!/usr/bin/env python


# don't do this in production
import warnings
warnings.filterwarnings("ignore")

import json

from elasticsearch import Elasticsearch

es = Elasticsearch()

with open('./data.json', encoding='utf-8') as myfile:
    for line in myfile:
        instrument = json.loads(line)
        res = es.index(index="instruments", id=instrument["symbol"], body=instrument)


res = es.search(index="instruments", body={"query": {"match_all": {}}})
print("Got %d Hits:" % res['hits']['total']['value'])
for hit in res['hits']['hits']:
    print("\t%(symbol)s %(company_name)s" % hit["_source"])
