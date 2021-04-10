terraform {
  required_providers {
    conformity = {
      version = "0.3"
      source  = "trendmicro.com/cloudone/conformity"
    }
  }
}

variable "auth_token" {
  default = ""
}

provider "conformity" {
  region = "ap-southeast-2"
  auth_token = var.auth_token
}

module "psl" {
  source = "./group"

  group_name = "AmitTest"
}

//output "psl" {
//  value = module.psl.group
//}
