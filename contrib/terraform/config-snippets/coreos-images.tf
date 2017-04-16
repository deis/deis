
variable "coreos_images" {
	default = {
	      eu-central-1-pv    = "ami-0c300d11"
	      ap-northeast-1-pv  = "ami-b128dcb1"
	      sa-east-1-pv       = "ami-2154ec3c"
	      ap-southeast-2-pv  = "ami-bbb5c581"
	      ap-southeast-1-pv  = "ami-fa0b3aa8"
	      us-east-1-pv       = "ami-343b195c"
	      us-west-2-pv       = "ami-0989a439"
	      us-west-1-pv       = "ami-83d533c7"
	      eu-west-1-pv       = "ami-57950a20"

	      eu-central-1-hvm   = "ami-0e300d13"
	      ap-northeast-1-hvm = "ami-af28dcaf"
	      sa-east-1-hvm      = "ami-2354ec3e"
	      ap-southeast-2-hvm = "ami-b9b5c583"
	      ap-southeast-1-hvm = "ami-f80b3aaa"
	      us-east-1-hvm      = "ami-323b195a"
	      us-west-2-hvm      = "ami-0789a437"
	      us-west-1-hvm      = "ami-8dd533c9"
	      eu-west-1-hvm      = "ami-55950a22"
	}
}
