
resource "aws_instance" "deis_node" {
    subnet_id = "${aws_subnet.deis_private.id}"
    associate_public_ip_address = false
}

