# terraform-provider-lambdazip

[![CI](https://github.com/winebarrel/terraform-provider-lambdazip/actions/workflows/ci.yml/badge.svg)](https://github.com/winebarrel/terraform-provider-lambdazip/actions/workflows/ci.yml)
[![terraform docs](https://img.shields.io/badge/terraform-docs-%35835CC?logo=terraform)](https://registry.terraform.io/providers/winebarrel/lambdazip/latest/docs)

Terraform provider creating zip file for AWS Lambda.

Update `base64sha256` attribute only when `triggers` attribute is updated.

## Usage

```
./
|-- lambda/
|   |-- index.js
|   |-- node_modules/
|   |-- package-lock.json
|   `-- package.json
`-- main.rf
```

```tf
terraform {
  required_providers {
    lambdazip = {
      source = "winebarrel/lambdazip"
    }
  }
}

provider "lambdazip" {
}

resource "lambdazip_file" "app" {
  base_dir      = "lambda"
  source        = "**"
  excludes      = [".env"]
  output        = "lambda.zip"
  before_create = "npm i"

  triggers = {
    for i in [
      "lambda/index.js",
      "lambda/package.json",
      "lambda/package-lock.json",
    ] : i => filesha256(i)
  }
}

resource "aws_lambda_function" "app" {
  filename         = lambdazip_file.app.output
  function_name    = "my_func"
  role             = aws_iam_role.lambda_app_role.arn
  handler          = "index.handler"
  source_code_hash = lambdazip_file.app.base64sha256
  runtime          = "nodejs20.x"
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

## Run locally for development

```sh
cp lambdazip.tf.sample lambdazip.tf
make tf-plan
make tf-apply
```
