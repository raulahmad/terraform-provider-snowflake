package datasources

import (
	"database/sql"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/snowflake"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var systemGetPrivateLinkConfigSchema = map[string]*schema.Schema{
	"account_name": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The name of your Snowflake account.",
	},

	"account_url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The URL used to connect to Snowflake through AWS PrivateLink or Azure Private Link.",
	},

	"ocsp_url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The OCSP URL corresponding to your Snowflake account that uses AWS PrivateLink or Azure Private Link.",
	},

	"aws_vpce_id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The AWS VPCE ID for your account.",
	},

	"azure_pls_id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The Azure Private Link Service ID for your account.",
	},
}

func SystemGetPrivateLinkConfig() *schema.Resource {
	return &schema.Resource{
		Read:   ReadSystemGetPrivateLinkConfig,
		Schema: systemGetPrivateLinkConfigSchema,
	}
}

// ReadSystemGetPrivateLinkConfig implements schema.ReadFunc
func ReadSystemGetPrivateLinkConfig(d *schema.ResourceData, meta interface{}) error {
	db := meta.(*sql.DB)

	sel := snowflake.SystemGetPrivateLinkConfigQuery()
	row := snowflake.QueryRow(db, sel)
	rawConfig, err := snowflake.ScanPrivateLinkConfig(row)

	if err == sql.ErrNoRows {
		// If not found, mark resource to be removed from statefile during apply or refresh
		log.Print("[DEBUG] system_get_privatelink_config not found")
		d.SetId("")
		return nil
	}

	config, err := rawConfig.GetStructuredConfig()
	if err != nil {
		log.Printf("[DEBUG] system_get_privatelink_config failed to decode")
		d.SetId("")
		return nil
	}

	d.SetId(config.AccountName)
	accNameErr := d.Set("account_name", config.AccountName)
	if accNameErr != nil {
		return accNameErr
	}
	accUrlErr := d.Set("account_url", config.AccountURL)
	if accUrlErr != nil {
		return accUrlErr
	}
	ocspUrlErr := d.Set("ocsp_url", config.OCSPURL)
	if ocspUrlErr != nil {
		return ocspUrlErr
	}

	if config.AwsVpceID != "" {
		awsVpceIdErr := d.Set("aws_vpce_id", config.AwsVpceID)
		if awsVpceIdErr != nil {
			return awsVpceIdErr
		}
	}

	if config.AzurePrivateLinkServiceID != "" {
		azurePlsIdErr := d.Set("azure_pls_id", config.AzurePrivateLinkServiceID)
		if azurePlsIdErr != nil {
			return azurePlsIdErr
		}
	}

	return nil
}
