package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	"time"
)

func resourceAlicloudSlbDomainExtension() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSlbDomainExtensionCreate,
		Read:   resourceAliyunSlbDomainExtensionRead,
		Update: resourceAliyunSlbDomainExtensionUpdate,
		Delete: resourceAliyunSlbDomainExtensionDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"frontend_port": {
				Type:         schema.TypeInt,
				ValidateFunc: validateIntegerInRange(1, 65535),
				Required:     true,
				ForceNew:     true,
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"server_certificate_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAliyunSlbDomainExtensionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slb_id := d.Get("load_balancer_id").(string)
	port := d.Get("frontend_port").(int)
	req := slb.CreateCreateDomainExtensionRequest()
	req.LoadBalancerId = slb_id
	req.ListenerPort = requests.NewInteger(port)
	req.Domain = d.Get("domain").(string)
	req.ServerCertificateId = d.Get("server_certificate_id").(string)

	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.CreateDomainExtension(req)
	})
	if err != nil {
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "slb_domain_extension", req.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	addDebug(req.GetActionName(), err)
	response, _ := raw.(*slb.CreateDomainExtensionResponse)
	d.SetId(response.DomainExtensionId)
	return resourceAliyunSlbDomainExtensionUpdate(d, meta)
}

func resourceAliyunSlbDomainExtensionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	loadBalancerId := d.Get("load_balancer_id").(string)
	port := d.Get("frontend_port").(int)
	id := d.Id()
	domainExtension, err := slbService.DescribeDomainExtensionAttribute(loadBalancerId, port, id)

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("id", domainExtension.DomainExtensions.DomainExtension[0].DomainExtensionId)
	d.Set("load_balancer_id", loadBalancerId)
	d.Set("domain", domainExtension.DomainExtensions.DomainExtension[0].Domain)
	d.Set("server_certificate_id", domainExtension.DomainExtensions.DomainExtension[0].ServerCertificateId)
	d.Set("frontend_port", port)
	return nil
}

func resourceAliyunSlbDomainExtensionUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	if d.HasChange("server_certificate_id") {
		req := slb.CreateSetDomainExtensionAttributeRequest()
		req.DomainExtensionId = d.Id()
		req.ServerCertificateId = d.Get("server_certificate_id").(string)
		client := meta.(*connectivity.AliyunClient)
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.SetDomainExtensionAttribute(req)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), req.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(req.GetActionName(), raw)
		d.SetPartial("server_certificate_id")
	}
	d.Partial(false)
	return resourceAliyunSlbDomainExtensionRead(d, meta)
}

func resourceAliyunSlbDomainExtensionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	req := slb.CreateDeleteDomainExtensionRequest()
	req.DomainExtensionId = d.Id()
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DeleteDomainExtension(req)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{""}) {
				return nil
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), req.GetActionName(), AlibabaCloudSdkGoERROR))
		}
		addDebug(req.GetActionName(), raw)

		client := meta.(*connectivity.AliyunClient)
		slbService := SlbService{client}
		lbId := d.Get("load_balancer_id").(string)
		port := d.Get("frontend_port").(int)
		if _, err := slbService.DescribeDomainExtensionAttribute(lbId, port, d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(WrapError(err))
		}
		return resource.RetryableError(WrapErrorf(err, DeleteTimeoutMsg, d.Id(), req.GetActionName(), ProviderERROR))
	})
}
