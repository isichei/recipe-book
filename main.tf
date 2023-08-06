provider "aws" {
  region = "eu-west-2"
}

### LAMBDAS
resource "aws_iam_role" "lambda_role" {
  name = "lambda-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}


# Attach a basic IAM policy to the role
resource "aws_iam_role_policy_attachment" "lambda_policy_attach" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
  role       = aws_iam_role.lambda_role.name
}

# Create a ZIP archive containing your Go binary
data "archive_file" "recipe_app_zip" {
  type        = "zip"
  source_file = "bin/bootstrap"
  output_path = "lambda-zips/aws-lambda-go.zip"
}

# Create the Lambda function
resource "aws_lambda_function" "recipe_app_lambda" {
  function_name   = "recipe-app"
  filename         = data.archive_file.recipe_app_zip.output_path
  role             = aws_iam_role.lambda_role.arn
  handler          = "bootstrap"
  runtime          = "go1.x"
  timeout          = 10
  source_code_hash = filebase64sha256(data.archive_file.recipe_app_zip.output_path)
}

# ---

# AWS API GATEWAY DEFINITIONS
resource "aws_api_gateway_rest_api" "app_rest_api" {
  name        = "recipe-api"
  description = "API Gateway for the recipe API"
}

### ROOT ###
# Create a method for the API's / path
resource "aws_api_gateway_method" "root_method" {
  rest_api_id   = aws_api_gateway_rest_api.app_rest_api.id
  resource_id   = aws_api_gateway_rest_api.app_rest_api.root_resource_id
  http_method   = "GET"
  authorization = "NONE"
}

# Create the lambda integration for / path
resource "aws_api_gateway_integration" "root_integration" {
  rest_api_id             = aws_api_gateway_rest_api.app_rest_api.id
  resource_id             = aws_api_gateway_rest_api.app_rest_api.root_resource_id
  http_method             = aws_api_gateway_method.root_method.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.recipe_app_lambda.invoke_arn
}
### ROOT ###

### SEARCH-RECIPES ###
# Create the resource for the /search-recipes path
resource "aws_api_gateway_resource" "search_recipes" {
    path_part = "search-recipes"
    rest_api_id = aws_api_gateway_rest_api.app_rest_api.id
    parent_id   = aws_api_gateway_rest_api.app_rest_api.root_resource_id
}

# Create a method for the API's /search-recipes path
resource "aws_api_gateway_method" "search_recipes_method" {
  rest_api_id   = aws_api_gateway_rest_api.app_rest_api.id
  resource_id   = aws_api_gateway_resource.search_recipes.id
  http_method   = "GET"
  authorization = "NONE"
}

# Create the lambda integration for /search-recipes path
resource "aws_api_gateway_integration" "search_recipes_integration" {
  rest_api_id             = aws_api_gateway_rest_api.app_rest_api.id
  resource_id             = aws_api_gateway_resource.search_recipes.id
  http_method             = aws_api_gateway_method.search_recipes_method.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.recipe_app_lambda.invoke_arn
}
### SEARCH-RECIPES ###

resource "aws_api_gateway_deployment" "app_deployment" {
  rest_api_id       = aws_api_gateway_rest_api.app_rest_api.id

  triggers = {
    redeployment = sha1(jsonencode(aws_api_gateway_rest_api.app_rest_api.body))
  }
  lifecycle {
    create_before_destroy = true
  }
  depends_on = [
    aws_api_gateway_method.root_method,
    aws_api_gateway_integration.root_integration,
    aws_api_gateway_resource.search_recipes,
    aws_api_gateway_method.search_recipes_method,
    aws_api_gateway_integration.search_recipes_integration
  ]
}

resource "aws_lambda_permission" "api_gw" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.recipe_app_lambda.arn}"
  principal     = "apigateway.amazonaws.com"

  source_arn    = "${aws_api_gateway_rest_api.app_rest_api.execution_arn}/*/*/*"
}

resource "aws_api_gateway_stage" "app_stage" {
  deployment_id = aws_api_gateway_deployment.app_deployment.id
  rest_api_id   = aws_api_gateway_rest_api.app_rest_api.id
  stage_name    = "dev"
}