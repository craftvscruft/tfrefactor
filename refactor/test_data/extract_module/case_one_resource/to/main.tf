
resource "a" "b" {
}

module "mymodule" {
  source = "mymodule"
}

moved {
  from = a.a
  to   = module.mymodule.a.a
}
