resource "alicloud_slb" "instance" {
   name                 = "${var.slb_name}"
   internet_charge_type = "${var.internet_charge_type}"
   internet             = "${var.internet}"
 }

 resource "alicloud_slb_server_certificate" "foo" {
   name               = "tf-testAccSlbServerCertificate"
   server_certificate = "${file("${path.module}/server_certificate.pem")}"
   private_key        = "${file("${path.module}/private_key.pem")}"
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

 resource "alicloud_slb_domain_extension" "example1" {
   load_balancer_id      = "${alicloud_slb.instance.id}"
   frontend_port         = "${alicloud_slb_listener.https.frontend_port}"
   domain                = "www.test.com"
   server_certificate_id = "${alicloud_slb_server_certificate.foo.id}"
 }
data "alicloud_slb_domain_extensions" "slb_domain_extensions" {
	load_balancer_id      = "${alicloud_slb.instance.id}"
	frontend_port         = "${alicloud_slb_listener.https.frontend_port}"
	output_file           = "examples/slb-domain-extension/test.json"
}