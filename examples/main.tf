terraform {
  required_providers {
    cleura = {
      source  = "hashicorp.com/johanlundberg/cleura"
      version = "0.1.0"
    }
  }
}


resource "cleura_openstack_user" "testuser" {
  name      = "testuserone"
  domain_id = "f6ba827d60094aae8068161719d7172c"
  enabled   = true
  projects = [
    {
      id = "5068f750207a4b1b81e91cb90cefd293"
      roles = [
        "member","swiftoperator"
      ]
    }
  ]
}
# resource "cleura_user" "testuser2" {
#   name = "testusertwo"
#   domain_id = "f6ba827d60094aae8068161719d7172c"
#   enabled = true
#   projects = [
#     {
#       id = "5068f750207a4b1b81e91cb90cefd293"
#       roles = [
#         "member","swiftoperator"
#       ]
#     }
#   ]
# }

# data "cleura_user" "user" {
#   id = "7afc6d525f3c4a88af879a4e9add8482"
# }
# # output "testoutput" {
# #   value = data.cleura_user.user
# # }