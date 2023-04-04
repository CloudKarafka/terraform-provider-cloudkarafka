# Terraform provider for CloudKarafka

## Development

First tell Terraform to use the local build of the provider

in `~/terraformrc` put this

```
dev_overrides {
  "hashicorp/cloudkarafka" = "/home/$USER/code/cloudkarafka/terraform-provider-cloudkarafka/"
}
direct {}
```

This will tell terraform to look in that directly for a binary called `terraform-provider-cloudkarafka`

To build that binary

``` shell
make
```
