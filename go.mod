module github.com/Mongey/terraform-provider-confluentcloud

go 1.16

require (
	cloud.google.com/go v0.74.0 // indirect
	cloud.google.com/go/storage v1.12.0 // indirect
	github.com/Shopify/sarama v1.30.1
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/aws/aws-sdk-go v1.36.18 // indirect
	github.com/cgroschupp/go-client-confluent-cloud v0.0.0-20210518145537-98176441a5a5
	github.com/fatih/color v1.10.0 // indirect
	github.com/go-resty/resty/v2 v2.2.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-uuid v1.0.2
	github.com/hashicorp/terraform-plugin-docs v0.5.1
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.10.0
	github.com/hashicorp/yamux v0.0.0-20200609203250-aecfd211c9ce // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.4.0 // indirect
	github.com/oklog/run v1.1.0 // indirect
	github.com/ulikunitz/xz v0.5.9 // indirect
	golang.org/x/tools v0.0.0-20201230224404-63754364767c // indirect
	google.golang.org/genproto v0.0.0-20201214200347-8c77b98c765d // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/cgroschupp/go-client-confluent-cloud => github.com/Mongey/go-client-confluent-cloud v0.0.0-20210716182312-db34016c1db0
