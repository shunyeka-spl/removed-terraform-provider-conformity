terraform {
  required_providers {
    conformity = {
      version = "0.3"
      source  = "trendmicro.com/cloudone/conformity"
    }
  }
}

provider "conformity" {
  region = "ap-southeast-2"
  auth_token = <authtoken>
}

module "psl" {
  source = "./group"

  group_name = "AmitTest"
}

output "psl" {
  value = module.psl.group
}
