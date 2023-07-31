terraform {
  required_providers {
    tmdb = {
      source = "hashicorp.com/edu/tmdb"
    }
  }
}

provider "tmdb" {}

data "tmdb_movies" "example" {}
