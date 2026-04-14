data "external_schema" "gorm" {
  program = [
    "go", "run", "-mod=mod", "./loader.go"
  ]
}

env "dev" {
  src = data.external_schema.gorm.url
  dev = "docker://postgres/15?search_path=messanger"
  migrations {
    dir = "file://migrations"
  }
}

env "docker" {
  src = data.external_schema.gorm.url
  dev = "postgres://user:pass@postgres:5432/gav?search_path=messanger"
  migrations {
    dir = "file://migrations"
  }
}