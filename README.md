# Terraform Gendocs [![Build Status](https://travis-ci.com/patrickdappollonio/terraform-gendocs.svg?token=EXg9HysCEtbxxFpp8VYg&branch=master)](https://travis-ci.com/patrickdappollonio/terraform-gendocs)

A small project using the `go/ast` package to parse the Go code in a simple Terraform provider code.
The code iterates over each one of the parameters and tries to parse the different names and attributes
from these to generate an array of different parameters.

Per Terraform definition, these parameters are encapsulated in resources. Each resource can have sub-resources
that this tool can parse as well.

Unfortunately, due to the usage of `go/ast` it's not possible to parse every single provider out there, but
if your provider follows some standard format like [`terraform-provider-oneview`](https://github.com/HewlettPackard/terraform-provider-oneview)
with the following assumptions, then this tool can handle the documentation for you.

### Provider assumptions

* You have a provider in a repository with the format `${hostname}/${username}/terraform-provider-${provider-name}`.
  It can be nested under multiple folders, but the important part is that the provider name follows the convention
  at the end.
* There's a `main.go` at the project's root directory.
* There's a folder with the name of the provider based on the folder name, so if the provider is called `abc`, the folder
  name will be `terraform-provider-abc` and the folder inside the project's root will be called `abc`.
* Inside the provider folder, you'll have a `provider.go` file which configures the different settings for your provider
  as well as the definition of each resource.
* The provider Go file defines the parameters for the provider itself using a `return &schema.Provider{}` format. The
  parameters are inlined in the same return statement.
* Each resource is named `${provider-name}_${resource_name}` and it's defined in a Go file called `resource_${resource-name}.go`,
  where the file name replaces dashes (`-`) for underscores (`_`).
* The resource Go file defines the parameters for the provider itself using a `return &schema.Resource{}` format. The
  parameters are inlined in the same return statement.

### Usage mode

`terraform-gendocs` expects a minimum of 2 parameters and a maximum of 3:

```
terraform-gendocs $IMPORT_PATH $FORMAT [$FILENAME]
```

* `$IMPORT_PATH` is the path to the project's import URI. Going back to the provider named `abc` above, this import
  path will be (assuming myself as the project owner) `github.com/patrickdappollonio/abc`.
* `$FORMAT` can be either `hcl` or `html`. When using `hcl` it'll generate a simple terraform file which contains
  all the different parameter names like this: `parameter_name # type: string, required: false, optional: true, computed: true, force-new: false`.
  If the parameter by itself holds a subresource, then the subresource parameters are going to be indented with one extra tab
  under the parent parameter. With `html` format on the other side you'll get a browsable UI with a jump-to-definition menu
  which jump to a table view of each resource with its own definition. Each table allows clicking on the title to sort as well.
* `$FILENAME` (optional), this parameter allows to override the exported file name. By default the exported file name will be
  `${provider-name}-docs.${ext}`. Using the `abc` example and outputting HTML, you'll get a file named `abc-docs.html`. Passing
  the optional filename will replace the `${provider-name}-docs` part with your own. Extension is maintained.

Some example calls will be:

```bash
terraform-gendocs github.com/HewlettPackard/terraform-provider-oneview html      # this will generate an HTML documentation for `terraform-provider-oneview` in an output file called `terraform-docs.html`
terraform-gendocs github.com/HewlettPackard/terraform-provider-oneview hcl mytf  # this will generate an HCL .tf documentation for `terraform-provider-oneview` in an output file called `mytf.tf`
```

### Example HTML UI

This is an example from the [`terraform-provider-aws`](https://github.com/terraform-providers/terraform-provider-aws) which
unfortunately it can't be properly parsed due to some resources not being inlined (see "provider assumptions" section, last bullet point)
-- although PRs are welcome if you want to improve it! But you can see most of the definitions as well as the jump-to-definition menu here.

![AWS Provider Docs](https://i.imgur.com/8X0mWR0.png)
