package openstack

import (
	"errors"
	"os"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"identity_endpoint": &schema.Schema{
				Type:        schema.TypeString,
				DefaultFunc: envDefaultFunc("OS_AUTH_URL"),
				Required:    true,
			},

			"user_id": &schema.Schema{
				Type:        schema.TypeString,
				DefaultFunc: envDefaultFunc("OS_USERID"),
				Optional:    true,
			},

			"username": &schema.Schema{
				Type:        schema.TypeString,
				DefaultFunc: envDefaultFunc("OS_USERNAME"),
				Optional:    true,
			},

			"password": &schema.Schema{
				Type:        schema.TypeString,
				DefaultFunc: envDefaultFunc("OS_PASSWORD"),
				Required:    true,
			},

			"tenant_id": &schema.Schema{
				Type:        schema.TypeString,
				DefaultFunc: envDefaultFunc("OS_TENANT_ID"),
				Optional:    true,
			},

			"tenant_name": &schema.Schema{
				Type:        schema.TypeString,
				DefaultFunc: envDefaultFunc("OS_TENANT_NAME"),
				Optional:    true,
			},

			"domain_id": &schema.Schema{
				Type:        schema.TypeString,
				DefaultFunc: envDefaultFunc("OS_DOMAIN_ID"),
				Optional:    true,
			},

			"domain_name": &schema.Schema{
				Type:        schema.TypeString,
				DefaultFunc: envDefaultFunc("OS_DOMAIN_NAME"),
				Optional:    true,
			},

			"compute_api_version": &schema.Schema{
				Type:        schema.TypeString,
				DefaultFunc: envDefaultFunc("OS_COMPUTE_API_VERSION"),
				Optional:    true,
			},

			"block_storage_api_version": &schema.Schema{
				Type:        schema.TypeString,
				DefaultFunc: envDefaultFunc("OS_VOLUME_API_VERSION"),
				Optional:    true,
			},

			"networking_api_version": &schema.Schema{
				Type:        schema.TypeString,
				DefaultFunc: envDefaultFunc("OS_NETWORK_API_VERSION"),
				Optional:    true,
			},

			"object_storage_api_version": &schema.Schema{
				Type:        schema.TypeString,
				DefaultFunc: envDefaultFunc("OS_OBJECT_API_VERSION"),
				Optional:    true,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"openstack_instance":    resourceInstance(),
			"openstack_keypair":     resourceKeypair(),
			"openstack_floating_ip": resourceFloatingIP(),
			"openstack_secgroup":    resourceSecgroup(),
			"openstack_volume":      resourceVolume(),
		},

		ConfigureFunc: configureProvider,
	}
}

func envDefaultFunc(k string) schema.SchemaDefaultFunc {
	return func() (interface{}, error) {
		if v := os.Getenv(k); v != "" {
			return v, nil
		}

		// Default API versions
		switch k {
		case "OS_VOLUME_API_VERSION":
			return "1", nil
		case "OS_COMPUTE_API_VERSION":
			return "2", nil
		case "OS_NETWORK_API_VERSION":
			return "2", nil
		case "OS_OBJECT_API_VERSION":
			return "1", nil
		}

		// catchall
		return "", nil
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {

	user_id := d.Get("user_id").(string)
	username := d.Get("username").(string)
	if user_id == "" && username == "" {
		return nil, errors.New("At least one of user_id or username must be specified")
	}

	tenant_id := d.Get("tenant_id").(string)
	tenant_name := d.Get("tenant_name").(string)
	if tenant_id == "" && tenant_name == "" {
		return nil, errors.New("At least one of tenant_id or tenant_name must be specified")
	}

	config := Config{
		IdentityEndpoint:        d.Get("identity_endpoint").(string),
		UserID:                  d.Get("user_id").(string),
		Username:                d.Get("username").(string),
		Password:                d.Get("password").(string),
		TenantID:                d.Get("tenant_id").(string),
		TenantName:              d.Get("tenant_name").(string),
		DomainID:                d.Get("domain_id").(string),
		DomainName:              d.Get("domain_name").(string),
		BlockStorageAPIVersion:  d.Get("block_storage_api_version").(string),
		ComputeAPIVersion:       d.Get("compute_api_version").(string),
		NetworkingAPIVersion:    d.Get("networking_api_version").(string),
		ObjectStorageAPIVersion: d.Get("object_storage_api_version").(string),
	}

	if err := config.NewClient(); err != nil {
		return nil, err
	}

	return &config, nil
}
