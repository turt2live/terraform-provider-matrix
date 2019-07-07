module github.com/turt2live/terraform-provider-matrix

go 1.12

require github.com/hashicorp/terraform v0.12.3

replace github.com/go-resty/resty => gopkg.in/resty.v1 v1.12.0

replace github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.4.2

replace google.golang.org/cloud => cloud.google.com/go v0.41.0
