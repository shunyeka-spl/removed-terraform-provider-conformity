terraform {
  required_providers {
    conformity = {
      version = "0.3"
      source  = "trendmicro.com/cloudone/conformity"
    }
  }
}

variable "group_name" {
  type    = string
  default = "TestAmit"
}

resource "conformity_group" "grp6" {
  name = "hellotest3"
  tags = ["h", "j"]
}

data "conformity_groups" "all" {}

# Returns all groups
output "all_groups" {
   value = data.conformity_groups.all.groups
}

# # Only returns packer spiced latte
//output "group" {
//  value = data.conformity_groups.all
//}
