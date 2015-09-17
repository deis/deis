provider "digitalocean" {
    token = "${var.token}"
}

resource "digitalocean_droplet" "deis" {
    count = "${var.instances}"
    image = "coreos-stable"
    name = "${var.prefix}-${count.index+1}"
    region = "${var.region}"
    size = "${var.size}"
    backups = "False"
    ipv6 = "False"
    private_networking = "True"
    ssh_keys = ["${var.ssh_keys}"]
    user_data = "${file("${path.module}/../coreos/user-data")}"
}

output "msg" {
    value = "Your hosts are ready to go! Continue following the documentation to install and start Deis. Your hosts are: ${join(", ", digitalocean_droplet.deis.*.ipv4_address)}"
}

output "ip" {
    value = "${digitalocean_droplet.deis.0.ipv4_address}"
}
