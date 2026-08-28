# Step 7 — Terraform foundation

This Terraform layer intentionally provisions only the first cloud resource: Artifact Registry.

We are **not** creating GKE, VPCs or Cloud SQL yet. Those will be added after the local Kubernetes layer is understood.

## Authenticate

```bash
gcloud auth application-default login
```

## Run

```bash
cd infra/terraform
cp terraform.tfvars.example terraform.tfvars
# edit terraform.tfvars and set the real project ID

terraform init
terraform fmt -check
terraform validate
terraform plan
terraform apply
```

Verify:

```bash
gcloud artifacts repositories list --location=asia-southeast1
```

Destroy the resource when finished with the exercise:

```bash
terraform destroy
```

## Why this is IaC

The repository definition is version-controlled. `terraform plan` shows the intended change before applying it, and `terraform destroy` can remove what this configuration created.
