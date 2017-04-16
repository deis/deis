
resource "aws_elb" "deis_web_elb" {
    subnets = ["${aws_subnet.deis_private.id}"]
}
