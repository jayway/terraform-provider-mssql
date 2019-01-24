package resources

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jayway/terraform-provider-mssql/sql"
)

const usernameProp = "username"
const passwordProp = "password"

// Login is the mssql_login resource
func Login() *schema.Resource {
	return &schema.Resource{
		Create: loginCreate,
		Read:   loginRead,
		Update: loginUpdate,
		Delete: loginDelete,

		Schema: map[string]*schema.Schema{
			usernameProp: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			passwordProp: &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func loginCreate(d *schema.ResourceData, meta interface{}) error {
	connector := meta.(sql.Connector)
	username := d.Get(usernameProp).(string)
	password := d.Get(passwordProp).(string)

	err := connector.CreateLogin(username, password)
	if err != nil {
		return err
	}

	d.SetId(username)

	return loginRead(d, meta)
}
func loginRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}
func loginUpdate(d *schema.ResourceData, meta interface{}) error {
	return loginRead(d, meta)
}
func loginDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
