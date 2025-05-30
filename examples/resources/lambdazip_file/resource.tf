terraform {
  required_providers {
    lambdazip = {
      source  = "winebarrel/lambdazip"
      version = ">= 0.5.0"
    }
  }
}

resource "lambdazip_file" "app" {
  base_dir      = "lambda"
  sources       = ["**"]
  excludes      = [".env"]
  output        = "lambda.zip"
  before_create = "npm i"

  contents = {
    extra_file = "Zap Zap Zap"
  }

  triggers = merge(
    {
      for i in [
        "lambda/index.js",
        "lambda/package.json",
        "lambda/package-lock.json",
      ] : i => filesha256(i)
    },
    {
      extra_file = sha256("Zap Zap Zap")
    }
  )
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
