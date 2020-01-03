package cloudkarafka

import (
	"fmt"
	"log"

	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var version string

func Provider() terraform.ResourceProvider {
	log.Printf("Terraform-Provider-CloudKarafka Version: %s", version)
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"apikey": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLOUDKARAFKA_APIKEY", nil),
				Description: "Key used to authentication to the CloudKarafka API",
			},
			"baseurl": &schema.Schema{
				Type:        schema.TypeString,
				Default:     "https://customer.cloudkarafka.com",
				Optional:    true,
				Description: "Base URL to CloudKarafka website",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"cloudkarafka_instance": resourceInstance(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	useragent := fmt.Sprintf("terraform-provider-cloudkarafka_v%s", version)
	log.Printf("[DEBUG] cloudkarafka::provider::configure useragent: %v", useragent)
	return api.New(d.Get("baseurl").(string), d.Get("apikey").(string), useragent), nil
}
