
resource "aws_subnet" "deis_private" {
    vpc_id = "${aws_vpc.deis.id}"

    cidr_block = "10.0.22.0/24"
    availability_zone = "${var.availability_zone}"
}

resource "aws_route_table" "deis_private" {
    vpc_id = "${aws_vpc.deis.id}"

    route {
        cidr_block = "0.0.0.0/0"
        instance_id = "${aws_instance.nat.id}"
    }
}

resource "aws_route_table_association" "deis_private" {
    subnet_id = "${aws_subnet.deis_private.id}"
    route_table_id = "${aws_route_table.deis_private.id}"
}

resource "aws_security_group" "nat" {
    name = "deis-nat"
    description = "Allow services from the Deis private subnet through NAT - Managed by Deis Terraform"

    ingress {
        from_port = 0
        to_port = 65535
        protocol = "tcp"
        cidr_blocks = ["${aws_subnet.deis_private.cidr_block}"]
    }
    ingress {
        from_port = 0
        to_port = 65535
        protocol = "tcp"
        cidr_blocks = ["${aws_subnet.deis_private.cidr_block}"]
    }

    ingress {
        from_port = 0
        to_port = 65535
        protocol = "udp"
        cidr_blocks = ["${aws_subnet.deis_private.cidr_block}"]
    }
    ingress {
        from_port = 0
        to_port = 65535
        protocol = "udp"
        cidr_blocks = ["${aws_subnet.deis_private.cidr_block}"]
    }
    tags {
        Name = "deis-nat"
    }

    vpc_id = "${aws_vpc.deis.id}"
}

resource "aws_instance" "nat" {
    ami = "${var.nat_ami}"
    instance_type = "m1.small"
    key_name = "${var.key_name}"
    security_groups = ["${aws_security_group.nat.id}"]
    subnet_id = "${aws_subnet.deis_public.id}"
    associate_public_ip_address = true
    source_dest_check = false
    tags {
        Name = "deis-nat"
    }
}

resource "aws_eip" "nat" {
    instance = "${aws_instance.nat.id}"
    vpc = true
}

resource "aws_instance" "bastion" {
    ami = "${var.bastion_ami}"
    instance_type = "t2.micro"
    key_name = "${var.key_name}"
    security_groups = ["${aws_security_group.deis_ssh.id}"]
    subnet_id = "${aws_subnet.deis_public.id}"
    tags {
        Name = "deis-bastion"
    }
}

resource "aws_eip" "bastion" {
    instance = "${aws_instance.bastion.id}"
    vpc = true
}

output "bastion-ip" {
    value = "${aws_instance.bastion.public_ip}"
}

output "nat-ip" {
    value = "${aws_instance.nat.public_ip}"
}
