resource "google_project_service" "artifact_registry" {
  project            = var.project_id
  service            = "artifactregistry.googleapis.com"
  disable_on_destroy = false
}

resource "google_artifact_registry_repository" "budget_api" {
  project       = var.project_id
  location      = var.region
  repository_id = var.artifact_repository
  description   = "Container images for the production-grade Budget Tracking API"
  format        = "DOCKER"

  depends_on = [google_project_service.artifact_registry]
}
