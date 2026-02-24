terraform {
  required_providers {
    level27 = {
      source  = "registry.terraform.io/level27/level27"
      version = "~> 0.1"
    }
  }
}

# Set the API key via the LEVEL27_API_KEY environment variable (recommended):
#   export LEVEL27_API_KEY="your-api-key"
#
# Or set it directly (not recommended for production):
provider "level27" {
  # api_key = "your-api-key-here"
}
