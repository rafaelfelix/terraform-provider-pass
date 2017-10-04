package pass

import (
	"context"
	"log"

	"gopkg.in/yaml.v2"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/justwatchcom/gopass/store/root"
	"github.com/justwatchcom/gopass/store/secret"
	"github.com/pkg/errors"
)

func passPasswordResource() *schema.Resource {
	return &schema.Resource{
		Create: passPasswordResourceWrite,
		Update: passPasswordResourceWrite,
		Delete: passPasswordResourceDelete,
		Read:   passPasswordResourceRead,

		Schema: map[string]*schema.Schema{
			"path": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Full path where the pass data will be written.",
			},

			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "secret password.",
			},

			"data": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "additional secret data.",
			},
		},
	}
}

func passPasswordResourceWrite(d *schema.ResourceData, meta interface{}) error {
	path := d.Get("path").(string)

	st := meta.(*root.Store)

	passwd := d.Get("passwd").(string)

	data := d.Get("data").(map[string]interface{})
	dataYaml, err := yaml.Marshal(&data)

	sec := secret.New(passwd, string(dataYaml))
	err = st.Set(context.Background(), path, sec)
	if err != nil {
		return errors.Wrapf(err, "failed to write secret at %s", path)
	}

	d.SetId(path)

	return nil
}

func passPasswordResourceDelete(d *schema.ResourceData, meta interface{}) error {
	path := d.Id()

	st := meta.(*root.Store)
	log.Printf("[DEBUG] Deleting generic Vault from %s", path)
	err := st.Delete(context.Background(), path)
	if err != nil {
		return errors.Wrapf(err, "failed to delete password at %s", path)
	}

	return nil
}

func passPasswordResourceRead(d *schema.ResourceData, meta interface{}) error {
	path := d.Id()

	st := meta.(*root.Store)
	sec, err := st.Get(context.Background(), path)
	if err != nil {
		errors.Wrapf(err, "failed to retrieve password at %s", path)
	}

	d.Set("password", sec.Password())
	d.Set("data", sec.Data())

	return nil
}
