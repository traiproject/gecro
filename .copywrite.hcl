schema_version = 1

project {
  license        = "MIT"
  copyright_year = 2025

  copyright_holder = "TRAI"

  header_ignore = [
    ".git/**",                       // Git directory
    "bazel-*/**",                    // Bazel output directories (bazel-bin, bazel-out, etc.)
    "vendor/**",                     // Go vendor directory
    "**/*_gen.go",                   // Generated Go files (e.g., wire_gen.go)
    "**/*.pb.go",                    // Generated Protobuf Go files
    "**/*.pb.gw.go",                 // Generated Protobuf gRPC Gateway files
    "**/testdata/**",                // Test data directories
    "go.work.sum",                   // Go workspace sum file
    "MODULE.bazel.lock",             // Bazel lock file for Bzlmod
    "LICENSE",                       // The main LICENSE file itself
    "README.md",                     // README files (unless you want headers in them)
    "*.yaml",                        // Configuration files - typically don't get headers
    "*.yml",                         // Configuration files
    ".copywrite.hcl",                // The copywrite config file itself
    ".copywrite/**",                 // Directory for copywrite templates
    ".lefthook/**"
  ]
}
