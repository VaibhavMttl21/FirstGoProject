data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    "--path", "./internal/model",
    "--dialect", "postgres",
  ]
}

locals {
  env = {
    for tuple in regexall(
      "(?m)^([A-Za-z_][A-Za-z0-9_]*)=(.*)$",
      file(".env")
    ) : tuple[0] => tuple[1]
  }
}

env "gorm" {
  src = data.external_schema.gorm.url

  url = local.env.DATABASE_URL

  dev = "docker://postgres/18/dev?search_path=public"

  migration {
    dir = "file://migrations"
  }
}