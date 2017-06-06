job "dp-content-api" {
  datacenters = ["DATA_CENTER"]
  update {
    stagger = "10s"
    max_parallel = 1
  }
  group "dp-content-api" {
    task "dp-content-api" {
      artifact {
        source = "s3::S3_TAR_FILE"
        destination = "."
        // The Following options are needed if no IAM roles are provided
        // options {
        // aws_access_key_id = ""
        // aws_access_key_secret = ""
        // }
      }
      env {
        PORT = "CONTENT_API_PORT"
        S3_URL = "S3_CONTENT_URL"
        S3_BUCKET = "S3_CONTENT_BUCKET"
        DB_ACCESS = "DATABASE_URL"
        GENERATOR_URL = "DP_GENERATOR_URL"
        HUMAN_LOG = "HUMAN_LOG_FLAG"
      }
      driver = "exec"
      config {
        command = "bin/dp-content-api"
      }
      resources {
        cpu = 500
        memory = 350
        network {
          port "http" {
            static = { "CONTENT_API_PORT" }
          }
        }
      }
      service {
        port = "http"
        check {
          type     = "http"
          path     = "HEALTHCHECK_ENDPOINT"
          interval = "10s"
          timeout  = "2s"
        }
      }
    }
  }
}
