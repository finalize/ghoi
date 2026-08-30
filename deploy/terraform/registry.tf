# コンテナイメージの置き場。
#
# PR 4 で Cloud Run に載せるイメージをここに push する。
resource "google_artifact_registry_repository" "app" {
  project       = var.project_id
  location      = var.region
  repository_id = "ghoi"
  description   = "Ghoi のコンテナイメージ"
  format        = "DOCKER"

  # 古いイメージが溜まり続けないようにする。保管にも料金がかかるため。
  cleanup_policies {
    id     = "古いものを消す"
    action = "DELETE"
    condition {
      older_than = "2592000s" # 30日
    }
  }

  # 直近のものは日数にかかわらず残す。上の削除より優先される。
  cleanup_policies {
    id     = "直近5つは残す"
    action = "KEEP"
    most_recent_versions {
      keep_count = 5
    }
  }

  depends_on = [google_project_service.enabled]
}
