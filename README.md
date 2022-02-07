# check_elasticsearch

Icinga check plugin to check the health status of an Elasticsearch cluster or the total hits/results of an Elasticsearch
query.

## Usage

### Health

Checks the health status of an Elasticsearch cluster.

```
Usage:
  check_elasticsearch health

Flags:
  -h, --help   help for health

Global Flags:
  -H, --hostname string   Hostname or ip address of elasticsearch node (default "localhost")
      --insecure          Allow use of self signed certificates when using SSL
  -P, --password string   Password if authentication is required
  -p, --port int          Port of elasticsearch node (default 9200)
  -S, --tls               Use secure connection
  -U, --username string   Username if authentication is required
```

#### Elasticsearch cluster with green status (all nodes are running)

```
$ check_elasticsearch health -U exampleuser -P examplepassword -S --insecure
OK - Cluster es-example-cluster is green | status=0 nodes=3 data_nodes=3 active_primary_shards=10 active_shards=20
```

#### Elasticsearch cluster with yellow status (not all nodes are running)

```
$ check_elasticsearch health -U exampleuser -P examplepassword -S --insecure
WARNING - Cluster es-example-cluster is yellow | status=1 nodes=2 data_nodes=2 active_primary_shards=10 active_shards=13```
```

### Query

Checks the total hits/results of an Elasticsearch query.<br>
The plugin is currently capable to return the total hits of documents based on a provided query string.

```
Usage:
  check_elasticsearch query [flags]

Flags:
  -q, --query string    Elasticsearch query
  -I, --index string    The index which will be used  (default "_all")
  -k, --msgkey string   Message of messagekey to display
  -m, --msglen int      Number of characters to display in latest message (default 80)
  -w, --warning uint    Warning threshold for total hits (default 20)
  -c, --critical uint   Critical threshold for total hits (default 50)
  -h, --help            help for query
```

#### Search for total hits without any message

```
$ check_elasticsearch query -q "event.dataset:sample_web_logs and @timestamp:[now-5m TO now]" -I "kibana_sample_data_logs"
CRITICAL - Total hits: 14074 | total=14074;20;50
```

#### Search for total hits with message

```
$ check_elasticsearch query -q "event.dataset:sample_web_logs and @timestamp:[now-5m TO now]" -I "kibana_sample_data_logs" -k "message"
CRITICAL - Total hits: 14074
30.156.16.163 - - [2018-09-01T12:44:53.756Z] "GET /wp-content/plugins/video-play
 | total=14074;20;50
```

## Further Documentation

* [Elasticsearch API Docs](https://www.elastic.co/guide/en/elasticsearch/reference/current/rest-apis.html)
* [Elasticsearch SDK for Go](https://github.com/elastic/go-elasticsearch)

## License

Copyright (c) 2022 [NETWAYS GmbH](mailto:info@netways.de)

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public
License as published by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied
warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not,
see [gnu.org/licenses](https://www.gnu.org/licenses/).
