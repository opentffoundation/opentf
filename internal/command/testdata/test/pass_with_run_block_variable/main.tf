variable "sample_test_value" {
  type    = string
  default = "nowhere"
}

provider "test" {
  region = "somewhere"
}