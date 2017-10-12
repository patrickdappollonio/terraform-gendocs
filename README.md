# Terraform Gendocs [![Build Status](https://travis-ci.org/patrickdappollonio/terraform-gendocs.svg?branch=master)](https://travis-ci.org/patrickdappollonio/terraform-gendocs)

This is a small Go project that parses the Go syntax tree pretty much like the compiler does,
fetches the parameters from the Terraform Schema definition and retrieves them for later use.

Right now, when used with `terraform-provider-oneview` we got the output below. For a more clear
explanation, feel free to use the [AST viewer from Yuroyoro here](http://goast.yuroyoro.net/).

There's a high chance this code won't work with a different way to write Go Terraform Providers,
but if it looks pretty much like the rest of the Terraform Providers, this should catch all the
different parameters and its attributes.

