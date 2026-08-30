# 使う API を有効にする。
#
# GCP は既定でほとんどの API が無効なので、使う前に明示的に有効化する必要がある。
# ここに並べるのは「この PR までで実際に使うもの」だけにして、
# Cloud Run や Secret Manager は必要になった PR で足す。

locals {
  services = [
    # Artifact Registry。コンテナイメージの置き場
    "artifactregistry.googleapis.com",
    # 他の API を有効化したり、プロジェクトの情報を読むために要る
    "cloudresourcemanager.googleapis.com",
    # Cloud Run。アプリを動かす場所
    "run.googleapis.com",
    # Cloud Build。手元に Docker が無くてもイメージを作れる
    "cloudbuild.googleapis.com",
  ]
}

resource "google_project_service" "enabled" {
  for_each = toset(local.services)

  project = var.project_id
  service = each.value

  # terraform destroy でも API を無効化しない。
  # 無効化すると、その API で作った資源が壊れることがあるため。
  disable_on_destroy = false
}
