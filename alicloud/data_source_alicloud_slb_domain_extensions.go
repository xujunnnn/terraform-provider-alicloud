package alicloud

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	"regexp"
)

func dataSourceAlicloudSlbDomainExtensions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSlbDomainExtensionsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
				MinItems: 1,
			},
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"frontend_port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameRegex,
				ForceNew:     true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"slb_domain_extensions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_certificate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudSlbDomainExtensionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := slb.CreateDescribeDomainExtensionsRequest()
	args.LoadBalancerId = d.Get("load_balancer_id").(string)
	args.ListenerPort = requests.NewInteger(d.Get("frontend_port").(int))

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeDomainExtensions(args)
	})
	if err != nil {
		return fmt.Errorf("DescribeDomainExtensions got an error: %#v", err)
	}
	resp, _ := raw.(*slb.DescribeDomainExtensionsResponse)
	if resp == nil {
		return fmt.Errorf("there is no SLB with the ID %s. Please change your search criteria and try again", args.LoadBalancerId)
	}

	var filteredDomainExtensionsTemp []slb.DomainExtension
	nameRegex, ok := d.GetOk("name_regex")
	if (ok && nameRegex.(string) != "") || len(idsMap) > 0 {
		var r *regexp.Regexp
		if nameRegex != "" {
			r = regexp.MustCompile(nameRegex.(string))
		}
		for _, domainExtension := range resp.DomainExtensions.DomainExtension {
			if r != nil && !r.MatchString(domainExtension.Domain) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[domainExtension.DomainExtensionId]; !ok {
					continue
				}
			}
			filteredDomainExtensionsTemp = append(filteredDomainExtensionsTemp, domainExtension)
		}
	} else {
		filteredDomainExtensionsTemp = resp.DomainExtensions.DomainExtension
	}
	return slbDomainExtensionDescriptionAttributes(d, filteredDomainExtensionsTemp)
}

func slbDomainExtensionDescriptionAttributes(d *schema.ResourceData, domainExtensions []slb.DomainExtension) error {
	var ids []string
	var s []map[string]interface{}
	for _, domainExtension := range domainExtensions {
		mapping := map[string]interface{}{
			"id":                    domainExtension.DomainExtensionId,
			"domain":                domainExtension.Domain,
			"server_certificate_id": domainExtension.ServerCertificateId,
		}
		ids = append(ids, domainExtension.DomainExtensionId)
		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("slb_domain_extensions", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
