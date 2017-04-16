
resource "aws_launch_configuration" "deis_launch_config" {
    name = "deis-node-launch-config"
    image_id = "${lookup(var.coreos_images, format("%s-%s", var.region, var.virt_type))}"
    instance_type = "${var.instance_type}"
    security_groups = [ "${aws_security_group.deis_ssh.id}", "${aws_security_group.deis_control.id}", "${aws_security_group.deis_web.id}", "${aws_security_group.intra_vpc.id}" ]
    key_name = "${var.key_name}"
}

resource "aws_autoscaling_group" "deis_nodes" {
    availability_zones = ["${var.availability_zone}"]
    vpc_zone_identifier = ["${aws_subnet.deis_public.id}"]
    name = "deis-autoscale-nodes"
    max_size = "${var.cluster_size}"
    min_size = "${var.cluster_size}"
    health_check_grace_period = 300
    health_check_type = "ELB"
    desired_capacity = "${var.cluster_size}"
    force_delete = true
    launch_configuration = "${aws_launch_configuration.deis_launch_config.name}"
    load_balancers = ["${aws_elb.deis_web_elb.name}"]
    ebs_block_device = {
        device_name = "/dev/xvdf"
        volume_size = "${var.docker_volume_size}"
    }
    tag {
        key = "Name"
        value = "deis-cluster-node"
        propagate_at_launch = true
    }
}
