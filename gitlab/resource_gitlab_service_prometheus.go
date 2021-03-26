package gitlab

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabServicePrometheus() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabServicePrometheusCreate,
		Read:   resourceGitlabServicePrometheusRead,
		Update: resourceGitlabServicePrometheusUpdate,
		Delete: resourceGitlabServicePrometheusDelete,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"api_url": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"google_iap_audience_client_id": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"google_iap_service_account_json": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
		},
	}
}

func resourceGitlabServicePrometheusSetToState(d *schema.ResourceData, service *gitlab.PrometheusService) {
	d.SetId(fmt.Sprintf("%d", service.ID))
	d.Set("api_url", service.Properties.APIURL)

	d.Set("title", service.Title)
	d.Set("created_at", service.CreatedAt.String())
	d.Set("updated_at", service.UpdatedAt.String())
	d.Set("active", service.Active)
}

func resourceGitlabServicePrometheusCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)

	opts := &gitlab.SetPrometheusServiceOptions{
		APIURL:                      gitlab.String(d.Get("api_url").(string)),
		GoogleIAPAudienceClientID:   gitlab.String(d.Get("google_iap_audience_client_id").(string)),
		GoogleIAPServiceAccountJSON: gitlab.String(d.Get("google_iap_service_account_json").(string)),
	}

	log.Printf("[DEBUG] create gitlab prometheus service for project %s", project)

	_, err := client.Services.SetPrometheusService(project, opts)
	if err != nil {
		return err
	}

	return resourceGitlabServicePrometheusRead(d, meta)
}

func resourceGitlabServicePrometheusRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)

	log.Printf("[DEBUG] read gitlab prometheus service for project %s", project)

	service, _, err := client.Services.GetPrometheusService(project)
	if err != nil {
		return err
	}

	resourceGitlabServicePrometheusSetToState(d, service)

	return nil
}

func resourceGitlabServicePrometheusUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceGitlabServicePrometheusCreate(d, meta)
}

func resourceGitlabServicePrometheusDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)

	log.Printf("[DEBUG] delete gitlab prometheus service for project %s", project)

	_, err := client.Services.DeletePrometheusService(project)
	return err
}
