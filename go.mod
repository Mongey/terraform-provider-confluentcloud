module github.com/Mongey/terraform-provider-confluent-cloud

go 1.12

require (
	github.com/cgroschupp/go-client-confluent-cloud v0.0.0-20191204162755-5bbf166f5417
	github.com/hashicorp/terraform v0.12.1
)

replace github.com/cgroschupp/go-client-confluent-cloud => ../go-client-confluent-cloud
