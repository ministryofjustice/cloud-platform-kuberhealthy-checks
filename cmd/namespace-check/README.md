## Namespace Check

The *Namespace* checks if all the vital namespaces exists in the cluster. It checks if the list of namespaces mentioned in the list `namespaces` exists in the cluster. If any of the namespaces doesnot exists, then this checks sends a status: `false` with message `Namespace check failed:namespaces [deleted-namespace] not found`

The KuberhealthyCheck custom resource definition can be found in https://github.com/ministryofjustice/cloud-platform-terraform-kuberhealthy/tree/main/resources. It is set to run every 30 sec (spec.runInterval), with a check timeout set to 2m (spec.timeout). If the check
does not complete within the given timeout it will report a timeout error on the status page.

