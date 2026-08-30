output "project_id" {
  description = "適用先のプロジェクト"
  value       = var.project_id
}

output "region" {
  value = var.region
}

output "registry_url" {
  description = "docker push の宛先。PR 4 で使う"
  value       = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.app.repository_id}"
}
