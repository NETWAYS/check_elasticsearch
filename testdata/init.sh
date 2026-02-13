# Setup an index and a pipeline

URL=http://localhost:9201

curl -s -k -X PUT -H 'Content-Type: application/json' "${URL}/_ingest/pipeline/example-pipeline" -d'
{
  "description": "My optional pipeline description",
  "processors": [
    {
      "set": {
        "description": "My optional processor description",
        "field": "my-long-field",
        "value": 10
      }
    }
  ]
}
'

curl -s -k -X PUT -H 'Content-Type: application/json' "${URL}/msg"
