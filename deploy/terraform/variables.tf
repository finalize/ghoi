variable "project_id" {
  description = "GCP のプロジェクト ID（名前ではなく ID）"
  type        = string
}

variable "region" {
  description = "資源を置くリージョン。Cloud Run も Artifact Registry もここに作る"
  type        = string
  default     = "asia-northeast1"
}

variable "image_tag" {
  description = "Cloud Run に載せるイメージのタグ。git の短い SHA を使う"
  type        = string
}
