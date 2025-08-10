# terraform-provider-lambdazip

[![CI](https://github.com/winebarrel/terraform-provider-lambdazip/actions/workflows/ci.yml/badge.svg)](https://github.com/winebarrel/terraform-provider-lambdazip/actions/workflows/ci.yml)
[![terraform docs](https://img.shields.io/badge/terraform-docs-%35835CC?logo=terraform)](https://registry.terraform.io/providers/winebarrel/lambdazip/latest/docs)

Terraform provider creating zip file for AWS Lambda.

Update `base64sha256` attribute only when `triggers` attribute is updated.

## Usage

```
./
|-- lambda-src/
|   |-- index.js
|   |-- node_modules/
|   |-- package-lock.json
|   `-- package.json
`-- main.tf
```

```tf
terraform {
  required_providers {
    lambdazip = {
      source  = "winebarrel/lambdazip"
      version = ">= 0.10.0"
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
  base_dir      = "lambda-src"
  sources       = ["**"]
  excludes      = [".env"]
  output        = "lambda.zip"
  before_create = "npm i"
  triggers      = data.lambdazip_files_sha256.triggers.map
  # use_temp_dir      = true
  # compression_level = 9

  # triggers = {
  #   for i in [
  #     "lambda/index.js",
  #     "lambda/package.json",
  #     "lambda/package-lock.json",
  #   ] : i => filesha256(i)
  # }
}

resource "aws_lambda_function" "app" {
  filename         = lambdazip_file.app.output
  function_name    = "my_func"
  role             = aws_iam_role.lambda_app_role.arn
  handler          = "index.handler"
  source_code_hash = lambdazip_file.app.base64sha256
  runtime          = "nodejs22.x"
}

resource "aws_iam_role" "lambda_app_role" {
  name = "lambda-app-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_app_role" {
  role       = aws_iam_role.lambda_app_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}
```

### Specify contents directly

```tf
locals {
  index_js = <<-EOT
    exports.handler = async () => {
      console.log("hello, world");
    };
  EOT
}

resource "lambdazip_file" "node_program" {
  base_dir = "lambda"
  output   = "lambda.zip"

  contents = {
    "index.js" = local.index_js
  }
}
```

### Golang example

```tf
data "lambdazip_files_sha256" "triggers" {
  files = ["lambda/*.go", "lambda/go.mod", "lambda/go.sum"]
}

resource "lambdazip_file" "app" {
  base_dir      = "lambda"
  sources       = ["bootstrap"]
  output        = "lambda.zip"
  before_create = "GOOS=linux GOARCH=amd64 go build -o bootstrap main.go"
  triggers      = data.lambdazip_files_sha256.triggers.map
}

resource "aws_lambda_function" "app" {
  filename         = lambdazip_file.app.output
  function_name    = "my_func"
  role             = aws_iam_role.lambda_app_role.arn
  handler          = "my-handler"
  source_code_hash = lambdazip_file.app.base64sha256
  runtime          = "provided.al2023"
}
```

> [!note]
> Go v1.21 or higher recommended.
> see https://go.dev/blog/rebuild


## Run locally for development

```sh
cp lambdazip.tf.sample lambdazip.tf
make
make tf-plan
make tf-apply
```
