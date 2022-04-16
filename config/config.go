package config

const (
	IndexName = "news112317"
	Mapping   = `
{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{
			"properties":{
				"ID":{
					"type":"keyword"
				},
				"Code":{
					"type":"long"
				},
				"Header":{
					"type":"keyword"
				},
				"Body":{
					"type":"keyword"
				},
				"PublishedAt":{
					"type":"date"
				},
				"Author":{
					"type":"keyword"
				},
				"Link":{
					"type":"keyword"
				},
				"Error":{
					"type":"keyword"
				}
			}
		}
}`
)
