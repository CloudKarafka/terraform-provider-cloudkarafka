package cloudkarafka

import (
	"github.com/84codes/go-api/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func ResourceInstance() *schema.Resource {
	return &schema.Resource{
		Create: ResourceCreate,
		Read:   ResourceRead,
		Update: ResourceUpdate,
		Delete: ResourceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the instance",
			},
			"plan": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the plan",
				ValidateFunc: validatePlanName(),
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the region you want to create your instance in",
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Tag the instances with optional tags",
			},
			"ca": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Broker CA",
			},
			"apikey": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "API key for the CloudAMQP instance",
			},
			"brokers": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Comma separated list of Kafka broker urls",
			},
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Username for accessing the Kafka cluster",
			},
			"password": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Password for accessing the Kafka cluster",
			},
			"topic_prefix": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			"kafka_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Kafka version",
			},
			"vpc_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the VPC to create your instance in",
			},
			"vpc_subnet": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Dedicated VPC subnet, shouldn't overlap with your current VPC's subnet",
			},
			"ready": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag describing if the resource is ready",
			},
			"keep_associated_vpc": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Keep associated VPC when deleting instance",
			},
		},
	}
}

func ResourceCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		api    = meta.(*api.API)
		keys   = instanceCreateAttributeKeys()
		params = make(map[string]interface{})
	)

	for _, k := range keys {
		if v := d.Get(k); v != nil || v != "" {
			params[k] = v
		}
		// CloudAMQP supports through go-api to fetch default version.
		// } else if k == "kafka_version" {
		// 	version, _ := api.DefaultKafkaVersion()
		// 	params[k] = version["default_kafka_version"]
		// }

		if k == "vpc_id" {
			if d.Get(k).(int) == 0 {
				delete(params, k)
			}
		}

		if k == "vpc_subnet" {
			if d.Get(k) == "" {
				delete(params, k)
			}
		}
	}

	data, err := api.CreateInstance(params)
	if err != nil {
		return err
	}

	d.SetId(data["id"].(string))
	return ResourceRead(d, meta)
}

func ResourceRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	data, err := api.ReadInstance(d.Id())
	if err != nil {
		return err
	}

	for k, v := range data {
		if validateInstanceJsonFields(k) {
			switch k {
			case "vpc":
				err = d.Set("vpc_id", v.(map[string]interface{})["id"])
				err = d.Set("vpc_subnet", v.(map[string]interface{})["subnet"])
			default:
				d.Set(k, v)
			}
		}
	}

	return nil
}

func ResourceUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		api    = meta.(*api.API)
		keys   = []string{"name", "plan", "tags"}
		params = make(map[string]interface{})
	)

	for _, k := range keys {
		params[k] = d.Get(k)
	}

	return api.UpdateInstance(d.Id(), params)
}

func ResourceDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*api.API)
	return api.DeleteInstance(d.Id(), d.Get("keep_associated_vpc").(bool))
}

func validateInstanceJsonFields(key string) bool {
	switch key {
	case "name",
		"plan",
		"region",
		"tags",
		"vpc",
		"ca",
		"apikey",
		"brokers",
		"username",
		"password",
		"topic_prefix",
		"kafka_version",
		"ready":
		return true
	}
	return false
}

func instanceCreateAttributeKeys() []string {
	return []string{
		"name",
		"plan",
		"region",
		"tags",
		"kafka_version",
		"vpc_subnet",
	}
}

func validatePlanName() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"ducky",
		"mouse-1", "mouse-3", "mouse-5", "mouse-7",
		"bat-1", "bat-3", "bat-5", "bat-7",
		"fox-1", "fox-3", "fox-5", "fox-7",
		"lion-1", "lion-3", "lion-5", "lion-7",
		"penguin-1", "penguin-3", "penguin-5", "penguin-7",
	}, true)
}
