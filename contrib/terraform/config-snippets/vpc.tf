
resource "aws_vpc" "deis" {
    cidr_block = "10.0.0.0/16"
}

resource "aws_internet_gateway" "deis" {
    vpc_id = "${aws_vpc.deis.id}"
}

resource "aws_subnet" "deis_public" {
    vpc_id = "${aws_vpc.deis.id}"

    cidr_block = "10.0.21.0/24"
    availability_zone = "${var.availability_zone}"
}

resource "aws_route_table" "deis_public" {
    vpc_id = "${aws_vpc.deis.id}"

    route {
        cidr_block = "0.0.0.0/0"
        gateway_id = "${aws_internet_gateway.deis.id}"
    }
}

resource "aws_route_table_association" "deis_public" {
    subnet_id = "${aws_subnet.deis_public.id}"
    route_table_id = "${aws_route_table.deis_public.id}"
}

