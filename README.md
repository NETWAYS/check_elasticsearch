# check_elasticsearch

Icinga check plugin to check the health status of an Elasticsearch cluster or the total hits/results of an Elasticsearch
query.

## Usage

```
Usage:
  check_elasticsearch [flags]
  check_elasticsearch [command]

Available Commands:
  health      Checks the health status of an Elasticsearch cluster
  query       Checks the total hits/results of an Elasticsearch query

Flags:
  -H, --hostname string   Hostname of the Elasticsearch instance (default "localhost")
  -p, --port int          Port of the Elasticsearch instance (default 9200)
  -U, --username string   Username if authentication is required
  -P, --password string   Password if authentication is required
  -S, --tls               Use a HTTPS connection
      --insecure          Skip the verification of the server's TLS certificate
  -t, --timeout int       Timeout in seconds for the CheckPlugin (default 30)
  -h, --help              help for check_elasticsearch
  -v, --version           version for check_elasticsearch
```

### Health

Checks the health status of an Elasticsearch cluster.

```
Usage:
  check_elasticsearch health

The cluster health status is:
  green = OK
  yellow = WARNING
  red = CRITICAL
```

Examples:

Elasticsearch cluster with green status (all nodes are running):

```
$ check_elasticsearch health -U exampleuser -P examplepassword -S --insecure
OK - Cluster es-example-cluster is green | status=0 nodes=3 data_nodes=3 active_primary_shards=10 active_shards=20
```

Elasticsearch cluster with yellow status (not all nodes are running):

```
$ check_elasticsearch health -U exampleuser -P examplepassword -S --insecure
WARNING - Cluster es-example-cluster is yellow | status=1 nodes=2 data_nodes=2 active_primary_shards=10 active_shards=13```
```

### Query

Checks the total hits/results of an Elasticsearch query.

Hint: The plugin is currently capable to return the total hits of documents based on a provided query string.

```
Usage:
  check_elasticsearch query [flags]

Flags:
  -q, --query string      The Elasticsearch query
  -I, --index string      Name of the Index which will be used (default "_all")
  -k, --msgkey string     Message of messagekey to display
  -m, --msglen int        Number of characters to display in the latest message (default 80)
  -w, --warning string    Warning threshold for total hits (default "20")
  -c, --critical string   Critical threshold for total hits (default "50")
  -h, --help              help for query
```

Examples:

Search for total hits without any message:

```
$ check_elasticsearch query -q "event.dataset:sample_web_logs and @timestamp:[now-5m TO now]" -I "kibana_sample_data_logs"
CRITICAL - Total hits: 14074 | total=14074;20;50
```

Search for total hits with message:

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
