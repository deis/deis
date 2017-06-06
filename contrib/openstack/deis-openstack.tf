provider "openstack" {
    user_name = "${var.username}"
    tenant_name = "${var.tenant}"
    password = "${var.password}"
    auth_url = "${var.auth_url}"
}

resource "openstack_networking_floatingip_v2" "deis" {
  count = "${var.instances}"
  pool = "${var.floatingip_pool}"
}

resource "openstack_compute_keypair_v2" "deis" {
  name = "${var.stack_name}"
  public_key = "${file(var.public_key_path)}"
}

resource "openstack_compute_secgroup_v2" "deis" {
  name = "${var.stack_name}"
  description = "Deis Security Group"
  rule {
    ip_protocol = "tcp"
    from_port = "22"
    to_port = "22"
    cidr = "0.0.0.0/0"
  }
  rule {
    ip_protocol = "tcp"
    from_port = "2222"
    to_port = "2222"
    cidr = "0.0.0.0/0"
  }
  rule {
    ip_protocol = "tcp"
    from_port = "80"
    to_port = "80"
    cidr = "0.0.0.0/0"
  }
  rule {
    ip_protocol = "icmp"
    from_port = "-1"
    to_port = "-1"
    cidr = "0.0.0.0/0"
  }
  rule {
    ip_protocol = "icmp"
    from_port = "-1"
    to_port = "-1"
    self = true
  }
  rule {
    ip_protocol = "tcp"
    from_port = "1"
    to_port = "65535"
    self = true
  }
  rule {
    ip_protocol = "udp"
    from_port = "1"
    to_port = "65535"
    self = true
  }
}

resource "openstack_compute_instance_v2" "deis" {
  name = "${var.stack_name}-${count.index+1}"
  count = "${var.instances}"
  image_name = "${var.image}"
  flavor_id = "${var.flavor}"
  key_pair = "${openstack_compute_keypair_v2.deis.name}"
  user_data = "${file(var.userdata)}"
  network {
    name = "${var.network_name}"
  }
  security_groups = [ "${openstack_compute_secgroup_v2.deis.name}" ]
  floating_ip = "${element(openstack_networking_floatingip_v2.deis.*.address, count.index)}"
}

output "msg" {
    value = "Your hosts are ready to go! Continue following the documentation to install and start Deis. Your hosts are: ${join(", ", openstack_compute_instance_v2.deis.*.network.fixed_ip_v4 )}"
}

#output "ip" {
#    value = "${join(", ", openstack_compute_instance_v2.deis.*.address})"
#}
