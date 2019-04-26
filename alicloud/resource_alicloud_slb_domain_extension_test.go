package alicloud

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	"strconv"
	"testing"
)

func TestAccAlicloudSlbDomainExtension(t *testing.T) {
	var domainExtension slb.DescribeDomainExtensionsResponse
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_slb_domain_extension.example",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbDomainExtensionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSlbDomainExtensionBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbDomainExtensionExists("alicloud_slb_domain_extension.example", &domainExtension),
					resource.TestCheckResourceAttrSet("alicloud_slb_domain_extension.example", "load_balancer_id"),
					resource.TestCheckResourceAttrSet("alicloud_slb_domain_extension.example", "frontend_port"),
					resource.TestCheckResourceAttr("alicloud_slb_domain_extension.example", "domain", "www.test.com"),
					resource.TestCheckResourceAttrSet("alicloud_slb_domain_extension.example", "server_certificate_id"),
				),
			},
		},
	})
}

func testAccCheckSlbDomainExtensionExists(n string, domain *slb.DescribeDomainExtensionsResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if ds.Primary.ID == "" {
			return fmt.Errorf("No SLB Rule ID is set")
		}
		loadBalancerId := ds.Primary.Attributes["load_balancer_id"]
		port,err := strconv.Atoi(ds.Primary.Attributes["frontend_port"])
		if err != nil {
			return err
		}
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		slbService := SlbService{client}
		domainExtension, err := slbService.DescribeDomainExtensionAttribute(loadBalancerId,port,ds.Primary.ID)
		if err != nil {
			return err
		}
		domain = domainExtension

		return nil
	}
}

func testAccCheckSlbDomainExtensionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	slbService := SlbService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_slb_domain_extension" {
			continue
		}
		loadBalancerId := rs.Primary.Attributes["load_balancer_id"]
		port,err := strconv.Atoi(rs.Primary.Attributes["frontend_port"])
		if err != nil {
			return err
		}
		// Try to find the Slb server group
		if _, err := slbService.DescribeDomainExtensionAttribute(loadBalancerId,port,rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("SLB DomainExtension %s still exist", rs.Primary.ID)
	}

	return nil
}

const testAccSlbDomainExtensionBasic =`
resource "alicloud_slb" "instance" {
  name                 = "tffTestDomainExtension"
  internet_charge_type = "PayByTraffic"
  internet             = "true"
}

resource "alicloud_slb_server_certificate" "foo" {
  name               = "tf-testAccSlbServerCertificate"
  server_certificate = "-----BEGIN CERTIFICATE-----\nMIIDdjCCAl4CCQCcm+erkcKN7DANBgkqhkiG9w0BAQsFADB9MQswCQYDVQQGEwJj\nbjELMAkGA1UECAwCYmoxEDAOBgNVBAcMB2JlaWppbmcxDzANBgNVBAoMBmFsaXl1\nbjELMAkGA1UECwwCc2MxFTATBgNVBAMMDHd3dy50ZXN0LmNvbTEaMBgGCSqGSIb3\nDQEJARYLMTIzQDEyMy5jb20wHhcNMTkwNDI2MDM0ODAxWhcNMjQwNDI1MDM0ODAx\nWjB9MQswCQYDVQQGEwJjbjELMAkGA1UECAwCYmoxEDAOBgNVBAcMB2JlaWppbmcx\nDzANBgNVBAoMBmFsaXl1bjELMAkGA1UECwwCc2MxFTATBgNVBAMMDHd3dy50ZXN0\nLmNvbTEaMBgGCSqGSIb3DQEJARYLMTIzQDEyMy5jb20wggEiMA0GCSqGSIb3DQEB\nAQUAA4IBDwAwggEKAoIBAQDKMKF5qmN/uoMjdH3D8aPRcUOA0s8rZpYhG8zbkF1j\n8gHYoB/FDvM7G7dfVsyjbMwLOxKvAhWvHHSpEz/t7gB+QdwrAMiMJwGmtCnXrh2E\nWiXgalMe1y4a/T5R7q+m4T1zFATf+kbnHWfkSGF4W7b6UBoaH+9StQ95CnqzNf/2\np/Of7+S0XzCxFXw8GIVzZk0xFe6lHJzaq06f3mvzrD+4rpO56tTUvrgTY/n61gsF\nZP7f0CJ2JQh6eNRFOEUSfxKu/Dy/+IsQxorCJY2Q59ZAf3rXrqDN104jw9PlwnLl\nqfZz3RMODN6BWjxE8rvRtT0qMfuAfv1gjBdWZN0hUYBRAgMBAAEwDQYJKoZIhvcN\nAQELBQADggEBAABzo82TxGp5poVkd5pLWj5ACgcBv8Cs6oH9D+4Jz9BmyuBUsQXh\n2aG0hQAe1mU61C9konsl/GTW8umJQ4M4lYEztXXwMf5PlBMGwebM0ZbSGg6jKtZg\nWCgJ3eP/FMmyXGL5Jji5+e09eObhUDVle4tdi0On97zBoz85W02rgWFAqZJwiEAP\nt+c7jX7uOSBq2/38iGStlrX5yB1at/gJXXiA5CL5OtlR3Okvb0/QH37efO1Nu39m\nlFi0ODPAVyXjVypAiLguDxPn6AtDTdk9Iw9B19OD4NrzNRWgSSX5vuxo/VcRcgWk\n3gEe9Ca0ZKN20q9XgthAiFFjl1S9ZgdA6Zc=\n-----END CERTIFICATE-----"
  private_key        = "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAyjCheapjf7qDI3R9w/Gj0XFDgNLPK2aWIRvM25BdY/IB2KAf\nxQ7zOxu3X1bMo2zMCzsSrwIVrxx0qRM/7e4AfkHcKwDIjCcBprQp164dhFol4GpT\nHtcuGv0+Ue6vpuE9cxQE3/pG5x1n5EhheFu2+lAaGh/vUrUPeQp6szX/9qfzn+/k\ntF8wsRV8PBiFc2ZNMRXupRyc2qtOn95r86w/uK6TuerU1L64E2P5+tYLBWT+39Ai\ndiUIenjURThFEn8Srvw8v/iLEMaKwiWNkOfWQH96166gzddOI8PT5cJy5an2c90T\nDgzegVo8RPK70bU9KjH7gH79YIwXVmTdIVGAUQIDAQABAoIBAE1J4a/8biR5S3/W\nG+03BYQeY8tuyjqw8FqfoeOcf9agwAvqybouSNQjeCk9qOQfxq/UWQQFK/zQR9gJ\nv7pX7GBXFK5rkj3g+0SaQhRsPmRFgY0Tl8qGPt2aSKRRNVv5ZeADmwlzRn86QmiF\nMp0rkfqFfDTYWEepZszCML0ouzuxsW/9tq7rvtSjsgATNt31B3vFa3D3JBi31jUh\n5nfR9A3bATze7mQw3byEDiVl5ASRDgYyur403P1fDnMy9DBHZ8NaPOsFF6OBpJal\nBJsG5z00hll5PFN2jfmBQKlvAeU7wfwqdaSnGHOfqf2DeTTaFjIQ4gUhRn/m6pLo\n6kXttLECgYEA9sng0Qz/TcPFfM4tQ1gyvB1cKnnGIwg1FP8sfUjbbEgjaHhA224S\nk3BxtX2Kq6fhTXuwusAFc6OVMAZ76FgrQ5K4Ci7+DTsrF28z4b8td+p+lO/DxgP9\nlTgN+ddsiTOV4fUef9Z3yY0Zr0CnBUMbQYRaV2UIbCdiB0G4V/bt9TsCgYEA0bya\nOo9wGI0RJV0bYP7qwO74Ra1/i1viWbRlS7jU37Q+AZstrlKcQ5CTPzOjKFKMiUzl\n4miWacZ0/q2n+Mvd7NbXGXTLijahnyOYKaHJYyh4oBymfkgAifRstE0Ki9gdvArb\n/I+emC0GvLSyfGN8UUeDJs4NmqdEXGqjo2JOV+MCgYALFv1MR5o9Y1u/hQBRs2fs\nPiGDIx+9OUQxYloccyaxEfjNXAIGGkcpavchIbgWiJ++PJ2vdquIC8TLeK8evL+M\n9M3iX0Q5UfxYvD2HmnCvn9D6Xl/cyRcfGnq+TGjrLW9BzSMGuZt+aiHKV0xqFx7l\nbc4leTvMqGRmURS4lzcQOwKBgQCDzA/i4sYfN25h21tcHXSpnsG3D2rJyQi5NCo/\nZjunA92/JqOTGuiFcLGHEszhhtY3ZXJET1LNz18vtzKJnpqrvOnYXlOVW/U+SqDQ\n8JDb1c/PVZGuY1KrXkR9HLiW3kz5IJ3S3PFdUVYdeTN8BQxXCyg4V12nJJtJs912\ny0zN3wKBgGDS6YttCN6aI4EOABYE8fI1EYQ7vhfiYsaWGWSR1l6bQey7KR6M1ACz\nZzMASNyytVt12yXE4/Emv6/pYqigbDLfL1zQJSLJ3EHJYTh2RxjR+AaGDudYFG/T\nliQ9YXhV5Iu2x1pNwrtFnssDdaaGpfA7l3xC00BL7Z+SAJyI4QKA\n-----END RSA PRIVATE KEY-----"
}

resource "alicloud_slb_listener" "https" {
  load_balancer_id          = "${alicloud_slb.instance.id}"
  backend_port              = 80
  frontend_port             = 443
  protocol                  = "https"
  sticky_session            = "on"
  sticky_session_type       = "insert"
  cookie                    = "testslblistenercookie"
  cookie_timeout            = 86400
  health_check              = "on"
  health_check_uri          = "/cons"
  health_check_connect_port = 20
  healthy_threshold         = 8
  unhealthy_threshold       = 8
  health_check_timeout      = 8
  health_check_interval     = 5
  health_check_http_code    = "http_2xx,http_3xx"
  bandwidth                 = 10
  ssl_certificate_id        = "${alicloud_slb_server_certificate.foo.id}"
}

resource "alicloud_slb_domain_extension" "example" {
  load_balancer_id      = "${alicloud_slb.instance.id}"
  frontend_port         = "${alicloud_slb_listener.https.frontend_port}"
  domain                = "www.test.com"
  server_certificate_id = "${alicloud_slb_server_certificate.foo.id}"
}
`