output "project_id" {
  value = var.project_id
}

output "artifact_registry_repository" {
  value = google_artifact_registry_repository.budget_api.name
}

output "docker_registry_host" {
  value = "${var.region}-docker.pkg.dev"
}
