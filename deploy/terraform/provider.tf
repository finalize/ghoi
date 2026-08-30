# 認証は ADC（Application Default Credentials）を使う。
# ここに鍵ファイルのパスを書かない。手元は gcloud auth application-default login、
# CI は Workload Identity Federation（PR 5）が同じ枠組みで認証情報を渡す。

provider "google" {
  project = var.project_id
  region  = var.region
}
