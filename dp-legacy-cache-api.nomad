job "dp-legacy-cache-api" {
  datacenters = ["eu-west-2"]
  region      = "eu"
  type        = "service"

  update {
    stagger          = "60s"
    min_healthy_time = "30s"
    healthy_deadline = "2m"
    max_parallel     = 1
    auto_revert      = true
  }

  group "web" {
    count = "{{WEB_TASK_COUNT}}"

    constraint {
      attribute = "${node.class}"
      value     = "web"
    }

    restart {
      attempts = 3
      delay    = "15s"
      interval = "1m"
      mode     = "delay"
    }

    network {
      port "http" {
        to = 29100
      }
    }

    service {
      name = "dp-legacy-cache-api"
      port = "http"
      tags = ["web"]

      check {
        type     = "http"
        path     = "/health"
        interval = "10s"
        timeout  = "2s"
      }
    }

    task "dp-legacy-cache-api-web" {
      driver = "docker"

      config {
        command = "./dp-legacy-cache-api"
        image   = "{{ECR_URL}}:concourse-{{REVISION}}"
        ports   = ["http"]
      }

      resources {
        cpu    = "{{WEB_RESOURCE_CPU}}"
        memory = "{{WEB_RESOURCE_MEM}}"
      }

      template {
        data = <<EOH
        # Configs based on environment (e.g. export BIND_ADDR=":{{ env "NOMAD_PORT_http" }}")
        # or static (e.g. export BIND_ADDR=":8080")

        # Secret configs read from vault
        {{ with (secret (print "secret/" (env "NOMAD_TASK_NAME"))) }}
        {{ range $key, $value := .Data }}
        export {{ $key }}="{{ $value }}"
        {{ end }}
        {{ end }}
        EOH

        destination = "secrets/app.env"
        env         = true
        splay       = "1m"
        change_mode = "restart"
      }

      vault {
        policies = ["dp-legacy-cache-api-web"]
      }
    }
  }

  group "publishing" {
    count = "{{PUBLISHING_TASK_COUNT}}"

    constraint {
      attribute = "${node.class}"
      value     = "publishing"
    }

    restart {
      attempts = 3
      delay    = "15s"
      interval = "1m"
      mode     = "delay"
    }

    network {
      port "http" {
        to = 29100
      }
    }

    service {
      name = "dp-legacy-cache-api"
      port = "http"
      tags = ["publishing"]

      check {
        type     = "http"
        path     = "/health"
        interval = "10s"
        timeout  = "2s"
      }
    }

    task "dp-legacy-cache-api-publishing" {
      driver = "docker"

      config {
        command = "./dp-legacy-cache-api"
        image   = "{{ECR_URL}}:concourse-{{REVISION}}"
        ports   = ["http"]
      }

      resources {
        cpu    = "{{PUBLISHING_RESOURCE_CPU}}"
        memory = "{{PUBLISHING_RESOURCE_MEM}}"
      }

      template {
        data = <<EOH
        # Configs based on environment (e.g. export BIND_ADDR=":{{ env "NOMAD_PORT_http" }}")
        # or static (e.g. export BIND_ADDR=":8080")

        # Secret configs read from vault
        {{ with (secret (print "secret/" (env "NOMAD_TASK_NAME"))) }}
        {{ range $key, $value := .Data }}
        export {{ $key }}={{ $value | toJSON }}
        {{ end }}
        {{ end }}
        EOH

        destination = "secrets/app.env"
        env         = true
        splay       = "1m"
        change_mode = "restart"
      }

      vault {
        policies = ["dp-legacy-cache-api-publishing"]
      }
    }
  }
}
