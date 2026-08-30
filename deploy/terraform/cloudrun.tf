locals {
  registry = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.app.repository_id}"
}

# Cloud Run が名乗るサービスアカウント。
#
# 既定のサービスアカウントは権限が広すぎるので、専用のものを作って最小から始める。
# いまは何の権限も付けていない（/healthz を返すのに要らないため）。
# Gemini を叩く PR、Cloud SQL に繋ぐ PR で、必要な role をここに足していく。
resource "google_service_account" "run" {
  project      = var.project_id
  account_id   = "ghoi-run"
  display_name = "Ghoi の Cloud Run サービス"
  description  = "Ghoi のアプリが名乗る身元。必要な権限だけを足していく"
}

resource "google_cloud_run_v2_service" "app" {
  project  = var.project_id
  name     = "ghoi"
  location = var.region

  # 手元から terraform destroy できるようにしておく。
  # 本番として使い始めたら true に変える。
  deletion_protection = false

  template {
    service_account = google_service_account.run.email

    containers {
      image = "${local.registry}/ghoi:${var.image_tag}"

      # Cloud Run はこの番号を環境変数 PORT でコンテナに渡す。
      # internal/config がそれを読んで待ち受ける。
      ports {
        container_port = 8080
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
        # リクエストを処理していない間は CPU を割り当てない。
        # 起動しっぱなしにしないことで、待機中の料金がほぼゼロになる。
        cpu_idle = true
      }
    }

    scaling {
      # 0 まで縮む。リクエストが来なければ料金は発生しない。
      min_instance_count = 0
      # 暴走したときの上限。個人利用なので低くしておく。
      max_instance_count = 2
    }
  }

  depends_on = [google_project_service.enabled]
}

# 誰でも見られるようにする。
#
# Cloud Run は既定で非公開なので、これを付けないとブラウザから 403 になる。
# いま公開されるのは /healthz だけ。語を引く API を足す前に、
# 認証（PR 14-15）を入れる。
resource "google_cloud_run_v2_service_iam_member" "public" {
  project  = var.project_id
  location = google_cloud_run_v2_service.app.location
  name     = google_cloud_run_v2_service.app.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}
