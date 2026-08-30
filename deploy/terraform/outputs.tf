output "project_id" {
  description = "適用先のプロジェクト"
  value       = var.project_id
}

output "region" {
  value = var.region
}

output "registry_url" {
  description = "イメージの push 先"
  value       = local.registry
}

output "service_url" {
  description = "Cloud Run の公開 URL"
  value       = google_cloud_run_v2_service.app.uri
}

output "service_account" {
  description = "アプリが名乗る身元。権限を足すときの相手"
  value       = google_service_account.run.email
}
