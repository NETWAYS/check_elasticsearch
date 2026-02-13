# Generate an example doc

URL=http://localhost:9201

curl -X POST -H 'Content-Type: application/json' "${URL}/msg/_doc?pipeline=example-pipeline" -d '{
  "body": "Example",
  "severityNumber": 4,
  "resource": {
    "service.name": "node1"
  },
  "attributes": {
   "team": "awesome"
  }
}'
