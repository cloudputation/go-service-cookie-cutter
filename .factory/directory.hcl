project "autodocs" {
  directory "CI_CONTROLS" {
    file "auto_push"
  }
  file "Dockerfile"
  file "GNUmakefile"
  file "README.md"
  directory "helpers" {
    file "docker_rebuild.sh"
    file "onsave.service-seed.sh"
    file "rebuild.sh"
  }
  file "main.go"
  directory "packages" {
    directory "api" {
      file "README.md"
      file "server.go"
      directory "v1" {
        file "health.go"
        file "system_status.go"
      }
    }
    directory "bootstrap" {
      file "README.md"
      file "bootstrap.go"
    }
    directory "cli" {
      file "README.md"
      file "cli.go"
    }
    directory "config" {
      file "README.md"
      file "config.go"
    }
    directory "logger" {
      file "README.md"
      file "logger.go"
    }
    directory "stats" {
      file "README.md"
      file "stats.go"
    }
  }
}
