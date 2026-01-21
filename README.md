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
  ingest      Checks the ingest statistics of Ingest Pipelines
  query       Checks the total hits/results of an Elasticsearch query
  snapshot    Checks the status of Elasticsearch snapshots

Flags:
  -H, --hostname string    Hostname of the Elasticsearch instance (CHECK_ELASTICSEARCH_HOSTNAME) (default "localhost")
  -p, --port int           Port of the Elasticsearch instance (default 9200)
  -U, --username string    Username for HTTP Basic Authentication (CHECK_ELASTICSEARCH_USERNAME)
  -P, --password string    Password for HTTP Basic Authentication (CHECK_ELASTICSEARCH_PASSWORD)
  -S, --tls                Use a HTTPS connection
      --insecure           Skip the verification of the server's TLS certificate
      --ca-file string     Specify the CA File for TLS authentication (CHECK_ELASTICSEARCH_CA_FILE)
      --cert-file string   Specify the Certificate File for TLS authentication (CHECK_ELASTICSEARCH_CERT_FILE)
      --key-file string    Specify the Key File for TLS authentication (CHECK_ELASTICSEARCH_KEY_FILE)
  -t, --timeout int        Timeout in seconds for the CheckPlugin (default 30)
  -h, --help               help for check_elasticsearch
  -v, --version            version for check_elasticsearch
```

The check plugin respects the environment variables `HTTP_PROXY`, `HTTPS_PROXY` and `NO_PROXY`.

Various flags can be set with environment variables, refer to the help to see which flags.

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
[OK] - Cluster es-example-cluster is green | status=0 nodes=3 data_nodes=3 active_primary_shards=10 active_shards=20
```

Elasticsearch cluster with yellow status (not all nodes are running):

```
$ check_elasticsearch health -U exampleuser -P examplepassword -S --insecure
[WARNING] - Cluster es-example-cluster is yellow | status=1 nodes=2 data_nodes=2 active_primary_shards=10 active_shards=13```
```

### Query

Checks the total hits/counts of an Elasticsearch query (using a query_string query type).

The plugin can count the number of documents based on a provided query string
and then compare it to the given thresholds

With the `--msgkey` flag extracts a value from a given field and shows in in the output.
This is intended to show message/body/log field values in the plugin output.

The `--index` flag supports index patterns like `my-index-*` and `index1,index2`.

```
Usage:
  check_elasticsearch query [flags]

Flags:
  -q, --query string      The Elasticsearch query to run (query_string type syntax)
  -I, --index string      Name of the Index which will be used (default "_all")
  -k, --msgkey string     Name of a field to display in the output (e.g. a message body)
  -m, --msglen int        Maximum number of characters to display from the requested field (default 80)
  -w, --warning string    Warning count threshold for total hits (default "20")
  -c, --critical string   Critical count threshold for total hits (default "50")
  -h, --help              help for query
```

Examples:

Search for total hits without any message:

```
$ check_elasticsearch query -q "event.dataset:sample_web_logs and @timestamp:[now-5m TO now]" -I "kibana_sample_data_logs"
[CRITICAL] - Search query hits: 14074 | query_hits=14074c;20;50
```

Search for total hits with message:

```
$ check_elasticsearch query -q "event.dataset:sample_web_logs and @timestamp:[now-5m TO now]" -I "kibana_sample_data_logs" -k "message"
[CRITICAL] - Search query hits: 14074
30.156.16.163 - - [2018-09-01T12:44:53.756Z] "GET /wp-content/plugins/video-play
 | query_hits=14074c;20;50
```

### Ingest

Checks the ingest statistics of Ingest Pipelines. Thresholds check against errors of an Elasticsearch Ingest Pipeline.

```
Checks the ingest statistics of Ingest Pipelines

Usage:
  check_elasticsearch ingest [flags]

Flags:
      --pipeline stringArray     Name of the pipeline to check. Can be used multiple times and supports regex.
      --failed-warning string    Warning threshold for failed ingest operations. Use min:max for a range. (default "10")
      --failed-critical string   Critical threshold for failed ingest operations. Use min:max for a range. (default "20")
  -h, --help                     help for ingest
```

Examples:

```
check_elasticsearch ingest --failed-warning 5 --failed-critical 10
[WARNING] - Ingest operations may not be alright
  \_[WARNING] Number of failed ingest operations for mypipeline: 6; | pipelines.mypipeline.failed=6c

check_elasticsearch ingest --pipeline foobar
[OK] - Ingest operations alright
  \_[OK] Number of failed ingest operations for foobar: 5; | pipelines.foobar.failed=5c
```

### Snapshot

Checks status of Snapshots.

```
Checks the status of Elasticsearch snapshots
The plugin maps snapshot status to the following status codes:

SUCCESS, Exit code 0
PARTIAL, Exit code 1
FAILED, Exit code 2
IN_PROGRESS, Exit code 3

If there are multiple snapshots the plugin uses the worst status

Usage:
  check_elasticsearch snapshot [flags]

Flags:
  -a, --all                         Check all retrieved snapshots. If not set only the latest snapshot is checked
  -N, --number int                  Check latest N number snapshots. If not set only the latest snapshot is checked (default 1)
  -r, --repository string           Comma-separated list of snapshot repository names used to limit the request (default "*")
  -s, --snapshot string             Comma-separated list of snapshot names to retrieve. Wildcard (*) expressions are supported (default "*")
  -T, --no-snapshots-state string   Set exit code to return if no snapshots are found. Supported values are 0, 1, 2, 3, OK, Warning, Critical, Unknown (case-insensitive - default "Unknown")
  -h, --help                        help for snapshot
```

Examples:

```
$ check_elasticsearch snapshot
[OK] - All evaluated snapshots are in state SUCCESS

$ check_elasticsearch snapshot --all -r myrepo
[CRITICAL] - At least one evaluated snapshot is in state FAILED

$ check_elasticsearch snapshot --number 5 -s mysnapshot
[WARNING] - At least one evaluated snapshot is in state PARTIAL
```

## License

Copyright (c) 2022 [NETWAYS GmbH](mailto:info@netways.de)

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public
License as published by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied
warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not,
see [gnu.org/licenses](https://www.gnu.org/licenses/).
