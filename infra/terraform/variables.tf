variable "project_id" {
  description = "GCP project ID for the DevOps lab."
  type        = string
}

variable "region" {
  description = "Primary GCP region."
  type        = string
  default     = "asia-southeast1"
}

variable "artifact_repository" {
  description = "Artifact Registry repository ID."
  type        = string
  default     = "budget-api"
}
