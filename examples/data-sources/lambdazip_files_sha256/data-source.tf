terraform {
  required_providers {
    lambdazip = {
      source  = "winebarrel/lambdazip"
      version = ">= 0.5.0"
    }
  }
}

data "lambdazip_files_sha256" "triggers" {
  files = [
    "lambda/*.js",
    "lambda/*.json",
  ]
}

resource "lambdazip_file" "app" {
  base_dir      = "lambda"
  sources       = ["**"]
  excludes      = [".env"]
  output        = "lambda.zip"
  before_create = "npm i"
  triggers      = data.lambdazip_files_sha256.map
}
