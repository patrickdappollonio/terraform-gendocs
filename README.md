# Terraform Gendocs

This is a small Go project that parses the Go syntax tree pretty much like the compiler does,
fetches the parameters from the Terraform Schema definition and retrieves them for later use.

Right now, when used with `terraform-provider-dreamcloud` we got the output below. Note that
it requires some fine-tuning -- since it's showing all the parameters as Go code -- as well as
some formatting, but it's definitely achievable.

For a more clear explanation, feel free to use the [AST viewer from Yuroyoro here](http://goast.yuroyoro.net/).

```text
File: resource_automount_key.go
main.docDefinition{Name:"name", Type:"String", Optional:false, Computed:false, Required:false, ForceNew:true}
main.docDefinition{Name:"auto_map", Type:"String", Optional:false, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"mount_info", Type:"String", Optional:false, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"tags", Type:"List", Optional:true, Computed:false, Required:false, ForceNew:false}
File: resource_cifs_qtree.go
main.docDefinition{Name:"cluster_name", Type:"String", Optional:false, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"vserver_name", Type:"String", Optional:false, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"volume_name", Type:"String", Optional:false, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"share_name", Type:"String", Optional:false, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"share_comment", Type:"String", Optional:true, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"tags", Type:"List", Optional:true, Computed:false, Required:false, ForceNew:false}
File: resource_mount_point.go
main.docDefinition{Name:"path", Type:"String", Optional:false, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"svm_name", Type:"String", Optional:false, Computed:true, Required:false, ForceNew:false}
main.docDefinition{Name:"automount_id", Type:"Int", Optional:false, Computed:true, Required:false, ForceNew:false}
main.docDefinition{Name:"volume_id", Type:"Int", Optional:false, Computed:true, Required:false, ForceNew:false}
main.docDefinition{Name:"qtree_id", Type:"Int", Optional:false, Computed:true, Required:false, ForceNew:false}
main.docDefinition{Name:"qtree_name", Type:"String", Optional:false, Computed:true, Required:false, ForceNew:false}
main.docDefinition{Name:"tags", Type:"List", Optional:true, Computed:false, Required:false, ForceNew:false}
File: resource_network_set.go
main.docDefinition{Name:"name", Type:"String", Optional:false, Computed:false, Required:false, ForceNew:true}
main.docDefinition{Name:"vlan_ids", Type:"List", Optional:true, Computed:false, Required:false, ForceNew:false}
File: resource_nfs_qtree.go
main.docDefinition{Name:"cluster_name", Type:"String", Optional:false, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"vserver_name", Type:"String", Optional:false, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"volume_name", Type:"String", Optional:false, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"qtree_name", Type:"String", Optional:false, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"export_policy_name", Type:"String", Optional:false, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"tags", Type:"List", Optional:true, Computed:false, Required:false, ForceNew:false}
File: resource_svm.go
main.docDefinition{Name:"cluster_name", Type:"String", Optional:false, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"vserver_name", Type:"String", Optional:false, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"root_volume_aggregate", Type:"String", Optional:false, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"data_lif_address", Type:"String", Optional:false, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"data_lif_gateway", Type:"String", Optional:false, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"ip_space_name", Type:"String", Optional:true, Computed:true, Required:false, ForceNew:false}
main.docDefinition{Name:"name_service_switch", Type:"String", Optional:true, Computed:true, Required:false, ForceNew:false}
main.docDefinition{Name:"cifs_allowed", Type:"String", Optional:true, Computed:true, Required:false, ForceNew:false}
main.docDefinition{Name:"ldap_client_file", Type:"String", Optional:true, Computed:true, Required:false, ForceNew:false}
main.docDefinition{Name:"ldap_client_ips", Type:"List", Optional:true, Computed:true, Required:false, ForceNew:false}
main.docDefinition{Name:"cifs_domain", Type:"String", Optional:true, Computed:true, Required:false, ForceNew:false}
main.docDefinition{Name:"org_unit", Type:"String", Optional:true, Computed:true, Required:false, ForceNew:false}
main.docDefinition{Name:"dns_domain_names", Type:"List", Optional:true, Computed:true, Required:false, ForceNew:false}
main.docDefinition{Name:"dns_servers", Type:"List", Optional:true, Computed:true, Required:false, ForceNew:false}
main.docDefinition{Name:"data_lif_node_name", Type:"String", Optional:true, Computed:true, Required:false, ForceNew:false}
main.docDefinition{Name:"data_lif_port_name", Type:"String", Optional:true, Computed:true, Required:false, ForceNew:false}
main.docDefinition{Name:"tags", Type:"List", Optional:true, Computed:false, Required:false, ForceNew:false}
File: resource_volume.go
main.docDefinition{Name:"volume_name", Type:"String", Optional:false, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"cluster_name", Type:"String", Optional:false, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"vserver_name", Type:"String", Optional:false, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"volume_size_gb", Type:"Int", Optional:false, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"export_policy_name", Type:"String", Optional:false, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"aggregate_type", Type:"String", Optional:false, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"caching", Type:"Bool", Optional:true, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"tags", Type:"List", Optional:true, Computed:false, Required:false, ForceNew:false}
File: provider.go
main.docDefinition{Name:"region", Type:"String", Optional:true, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"ssl", Type:"Bool", Optional:true, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"endpoint_override", Type:"String", Optional:true, Computed:false, Required:false, ForceNew:false}
File: resource_instance.go
main.docDefinition{Name:"record_primary_id", Type:"String", Optional:false, Computed:true, Required:false, ForceNew:false}
main.docDefinition{Name:"hostname", Type:"String", Optional:false, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"domain", Type:"String", Optional:false, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"instance_type", Type:"String", Optional:false, Computed:false, Required:false, ForceNew:true}
main.docDefinition{Name:"os", Type:"String", Optional:false, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"network_roles", Type:"List", Optional:false, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"chef_role", Type:"String", Optional:true, Computed:false, Required:false, ForceNew:true}
main.docDefinition{Name:"chef_org", Type:"String", Optional:true, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"firmware", Type:"String", Optional:true, Computed:false, Required:false, ForceNew:true}
main.docDefinition{Name:"location", Type:"String", Optional:true, Computed:false, Required:false, ForceNew:true}
main.docDefinition{Name:"machine_type", Type:"String", Optional:false, Computed:true, Required:false, ForceNew:false}
main.docDefinition{Name:"tags", Type:"List", Optional:true, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"disk", Type:"List", Optional:true, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"id", Type:"Int", Optional:false, Computed:true, Required:false, ForceNew:false}
main.docDefinition{Name:"capacity", Type:"Int", Optional:false, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"type", Type:"String", Optional:true, Computed:true, Required:false, ForceNew:false}
main.docDefinition{Name:"mode", Type:"String", Optional:true, Computed:true, Required:false, ForceNew:false}
main.docDefinition{Name:"name", Type:"String", Optional:true, Computed:true, Required:false, ForceNew:false}
main.docDefinition{Name:"vm_manager", Type:"String", Optional:true, Computed:true, Required:false, ForceNew:false}
main.docDefinition{Name:"vm_record_id", Type:"Int", Optional:false, Computed:true, Required:false, ForceNew:false}
main.docDefinition{Name:"media", Type:"String", Optional:true, Computed:true, Required:false, ForceNew:false}
main.docDefinition{Name:"partition_table", Type:"String", Optional:true, Computed:true, Required:false, ForceNew:false}
main.docDefinition{Name:"cluster", Type:"String", Optional:true, Computed:true, Required:false, ForceNew:false}
File: resource_network.go
main.docDefinition{Name:"name", Type:"String", Optional:false, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"vlan_id", Type:"Int", Optional:false, Computed:false, Required:false, ForceNew:true}
main.docDefinition{Name:"address_cidr", Type:"String", Optional:false, Computed:false, Required:false, ForceNew:true}
main.docDefinition{Name:"locations", Type:"List", Optional:true, Computed:true, Required:false, ForceNew:false}
main.docDefinition{Name:"domains", Type:"List", Optional:true, Computed:true, Required:false, ForceNew:false}
main.docDefinition{Name:"organizations", Type:"List", Optional:true, Computed:true, Required:false, ForceNew:false}
main.docDefinition{Name:"rhev_targets", Type:"List", Optional:true, Computed:true, Required:false, ForceNew:false}
main.docDefinition{Name:"vsphere_targets", Type:"List", Optional:true, Computed:true, Required:false, ForceNew:false}
main.docDefinition{Name:"oneview_uplink_sets", Type:"List", Optional:true, Computed:true, Required:false, ForceNew:false}
main.docDefinition{Name:"proxy", Type:"String", Optional:true, Computed:true, Required:false, ForceNew:false}
main.docDefinition{Name:"location_ids", Type:"List", Optional:false, Computed:true, Required:false, ForceNew:false}
main.docDefinition{Name:"domain_ids", Type:"List", Optional:false, Computed:true, Required:false, ForceNew:false}
main.docDefinition{Name:"organization_ids", Type:"List", Optional:false, Computed:true, Required:false, ForceNew:false}
File: resource_vm_profile.go
main.docDefinition{Name:"name", Type:"String", Optional:false, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"cores", Type:"Int", Optional:false, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"memory", Type:"Int", Optional:false, Computed:false, Required:false, ForceNew:false}
main.docDefinition{Name:"storage_capacity", Type:"Int", Optional:false, Computed:false, Required:false, ForceNew:false}
```
