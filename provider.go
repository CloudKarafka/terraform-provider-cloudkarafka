package main

import (
	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
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
	return api.New(d.Get("baseurl").(string), d.Get("apikey").(string)), nil
}
