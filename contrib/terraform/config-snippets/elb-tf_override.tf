
resource "aws_elb" "deis_web_elb" {
    instances = ["${aws_instance.deis_node.*.id}"]
}
