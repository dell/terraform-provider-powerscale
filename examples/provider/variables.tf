/*
Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://mozilla.org/MPL/2.0/


Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

# PowerScale API endpoint URL including port 8080
# Example: "https://192.168.1.100:8080"
variable "endpoint" {
  description = "PowerScale API endpoint"
  type        = string
}

# Username for API authentication (usually "root")
variable "username" {
  description = "PowerScale API username"
  type        = string
}

# Password for API authentication
variable "password" {
  description = "PowerScale API password"
  type        = string
  sensitive   = true
}

# Skip SSL certificate verification
# Set to true when using self-signed certificates
# Default: false (SSL verification enabled)
variable "insecure" {
  description = "Skip SSL certificate verification"
  type        = bool
  default     = false
}

# API request timeout in milliseconds
# Increase this value for slow networks or large responses
# Default: 2000 (2 seconds)
variable "timeout" {
  description = "API request timeout in milliseconds"
  type        = number
  default     = 2000
}

# Authentication type for the API
# 0 = Basic authentication (sends credentials with each request)
# 1 = Session-based authentication (creates a session, more efficient)
# Default: 1 (session-based)
variable "auth_type" {
  description = "Authentication type: 0 = basic, 1 = session"
  type        = number
  default     = 1
}
