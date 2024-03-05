# Terraform provider for CloudKarafka

Available here: https://registry.terraform.io/providers/CloudKarafka/cloudkarafka/latest

> [!WARNING]  
> The CloudKarafka service will reach its [End of Life on January 27, 2025](https://www.cloudkarafka.com/blog/end-of-life-announcement.html).

## Development

First tell Terraform to use the local build of the provider

in `~/.terraformrc` put this

```
dev_overrides {
  "hashicorp/cloudkarafka" = "/home/$USER/code/cloudkarafka/terraform-provider-cloudkarafka/"
}
direct {}
```

This will tell terraform to look in that directly for a binary called `terraform-provider-cloudkarafka`

To build that binary

``` shell
go build
```
