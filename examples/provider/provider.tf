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

  triggers = [
    filesha256("example/index.js"),
    filesha256("example/package.json"),
    filesha256("example/package-lock.json"),
  ]
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
