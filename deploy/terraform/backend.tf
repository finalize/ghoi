# state の置き場。
#
# ここだけ Terraform では作れない（state を置く場所が無いと init できないため）。
# バケットは gcloud で一度だけ作ってある。手順は README に書いた。
#
# バージョニングを有効にしてあるので、state を壊しても前の版に戻せる。

terraform {
  backend "gcs" {
    bucket = "ghoi-507101-tfstate"
    prefix = "terraform/state"
  }
}
