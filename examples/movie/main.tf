terraform {
  required_providers {
    tmdb = {
      source = "hashicorp.com/edu/tmdb"
    }
  }
}

provider "tmdb" {}

data "tmdb_movie" "movie" {
    id = 108
}

output "movie" {
  value = data.tmdb_movie.movie
}