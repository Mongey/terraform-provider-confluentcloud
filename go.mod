module github.com/Mongey/terraform-provider-confluentcloud

go 1.16

require (
	github.com/Shopify/sarama v1.33.0
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/cgroschupp/go-client-confluent-cloud v0.0.0-20210518145537-98176441a5a5
	github.com/go-resty/resty/v2 v2.2.0 // indirect
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/hashicorp/go-uuid v1.0.3
	github.com/hashicorp/terraform-plugin-docs v0.9.0
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.10.1
	github.com/hashicorp/yamux v0.0.0-20200609203250-aecfd211c9ce // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.4.0 // indirect
	github.com/oklog/run v1.1.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20201214200347-8c77b98c765d // indirect
	google.golang.org/grpc v1.34.0 // indirect
)

replace github.com/cgroschupp/go-client-confluent-cloud => github.com/Mongey/go-client-confluent-cloud v0.0.0-20210716182312-db34016c1db0
