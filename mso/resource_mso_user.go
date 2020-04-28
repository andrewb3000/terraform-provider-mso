package mso

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceMSOUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceMSOUserCreate,
		Update: resourceMSOUserUpdate,
		Read:   resourceMSOUserRead,
		Delete: resourceMSOUserDelete,

		SchemaVersion: 1,

		Schema: (map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"user_password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"first_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"last_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"email": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"phone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"account_status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"roles": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"roleid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"access_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Required: true,
			},
		}),
	}
}

func resourceMSOUserCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Schema: Beginning Creation")
	msoClient := m.(*client.Client)
	user := d.Get("username").(string)
	userPassword := d.Get("user_password").(string)
	firstName := d.Get("first_name").(string)
	lastName := d.Get("last_name").(string)
	email := d.Get("email").(string)
	phone := d.Get("phone").(string)
	accountStatus := d.Get("account_status").(string)
	domain := d.Get("domain").(string)
	roles := make([]interface{}, 0, 1)
	if val, ok := d.GetOk("roles"); ok {
		tp := val.(*schema.Set).List()
		for _, val := range tp {

			map1 := make(map[string]interface{})
			inner := val.(map[string]interface{})
			if inner["roleid"] != "" {
				map1["roleId"] = fmt.Sprintf("%v", inner["roleid"])
			}
			if inner["access_type"] != "" {
				map1["accessType"] = fmt.Sprintf("%v", inner["access_type"])
			}
			
			roles = append(roles, map1)
		}
		
	}


	userApp := models.NewUser("", user, userPassword, firstName, lastName, email, phone, accountStatus, domain, roles)
	cont, err := msoClient.Save("api/v1/users", userApp)
	if err != nil {
		return err
	}

	id := models.StripQuotes(cont.S("id").String())
	d.SetId(fmt.Sprintf("%v", id))
	log.Printf("[DEBUG] %s: Schema Creation finished successfully", d.Id())

	return resourceMSOUserRead(d, m)
}

func resourceMSOUserUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Schema: Beginning Creation")
	msoClient := m.(*client.Client)
	user := d.Get("username").(string)
	userPassword := d.Get("user_password").(string)
	firstName := d.Get("first_name").(string)
	lastName := d.Get("last_name").(string)
	email := d.Get("email").(string)
	phone := d.Get("phone").(string)

	accountStatus := d.Get("account_status").(string)
	domain := d.Get("domain").(string)
	roles := make([]interface{}, 0, 1)
	if val, ok := d.GetOk("roles"); ok {

		tp := val.(*schema.Set).List()
		for _, val := range tp {

			map1 := make(map[string]interface{})
			inner := val.(map[string]interface{})
			map1["roleId"] = fmt.Sprintf("%v", inner["roleid"])

			if inner["access_type"] != "" {
				map1["accessType"] = fmt.Sprintf("%v", inner["access_type"])
			}

			roles = append(roles, map1)
		}

	}

	userApp := models.NewUser("", user, userPassword, firstName, lastName, email, phone, accountStatus, domain, roles)

	cont, err := msoClient.Put(fmt.Sprintf("api/v1/users/%s", d.Id()), userApp)
	if err != nil {
		return err
	}

	id := models.StripQuotes(cont.S("id").String())
	d.SetId(fmt.Sprintf("%v", id))
	log.Printf("[DEBUG] %s: Schema Creation finished successfully", d.Id())

	return resourceMSOUserRead(d, m)
	return nil

}

func resourceMSOUserRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	msoClient := m.(*client.Client)

	dn := d.Id()
	con, err := msoClient.GetViaURL("api/v1/users/" + dn)

	if err != nil {
		return err
	}

	d.SetId(models.StripQuotes(con.S("id").String()))
	d.Set("username", models.StripQuotes(con.S("username").String()))
	//d.Set("user_password", models.StripQuotes(con.S("password").String()))
	if con.Exists("firstName") {
		d.Set("first_name", models.StripQuotes(con.S("firstName").String()))
	}
	if con.Exists("lastName") {
		d.Set("last_name", models.StripQuotes(con.S("lastName").String()))
	}
	if con.Exists("emailAddress") {
		d.Set("email", models.StripQuotes(con.S("emailAddress").String()))
	}
	if con.Exists("phoneNumber") {
		d.Set("phone", models.StripQuotes(con.S("phoneNumber").String()))
	}
	if con.Exists("accountStatus") {
		d.Set("account_status", models.StripQuotes(con.S("accountStatus").String()))
	}
	if con.Exists("domain") {
		d.Set("domain", models.StripQuotes(con.S("domain").String()))
	}
	count, err := con.ArrayCount("roles")
	if err != nil {
		return fmt.Errorf("No Roles found")
	}

	roles := make([]interface{}, 0)
	for i := 0; i < count; i++ {
		rolesCont, err := con.ArrayElement(i, "roles")

		if err != nil {
			return fmt.Errorf("Unable to parse the roles list")
		}

		map1 := make(map[string]interface{})

		map1["roleid"] = models.StripQuotes(rolesCont.S("roleId").String())
		map1["access_type"] = models.StripQuotes(rolesCont.S("accessType").String())
		roles = append(roles, map1)
	}
	d.Set("roles", roles)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceMSOUserDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	msoClient := m.(*client.Client)
	dn := d.Id()
	err := msoClient.DeletebyId("api/v1/users/" + dn)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}

// func toStringList(configured interface{}) []string {
// 	vs := make([]string, 0, 1)
// 	val, ok := configured.(string)
// 	if ok && val != "" {
// 		vs = append(vs, val)
// 	}
// 	return vs
// }
