# 使う版を固定する。
#
# terraform 本体の版は mise.toml でも固定してあり、ここはその下限を宣言している。
# provider は .terraform.lock.hcl で厳密に固定されるので、ここは範囲でよい。
# （lock ファイルはコミットする。しないと人や CI ごとに違う版が入る）

terraform {
  required_version = "~> 1.16"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 7.0"
    }
  }
}
