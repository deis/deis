
resource "aws_instance" "deis_node" {
    count = "${var.cluster_size}"
    ami = "${lookup(var.coreos_images, format("%s-%s", var.region, var.virt_type))}"
    instance_type = "${var.instance_type}"
    vpc_security_group_ids = [ "${aws_security_group.deis_ssh.id}", "${aws_security_group.deis_control.id}", "${aws_security_group.deis_web.id}", "${aws_security_group.intra_vpc.id}" ]
    key_name = "${var.key_name}"
    subnet_id = "${aws_subnet.deis_public.id}"
    associate_public_ip_address = true
    user_data = "${file("coreos-user-data.txt")}"
    ebs_block_device = {
        device_name = "/dev/xvdf"
        volume_size = "${var.docker_volume_size}"
    }
    tags {
        Name = "deis-cluster-node"
    }
}

output "node-private-ips" {
    value = "${join(", ", aws_instance.deis_node.*.private_ip)}"
}

output "node-public-ips" {
    value = "${join(", ", aws_instance.deis_node.*.public_ip)}"
}
