terraform {
  required_providers {
    tmdb = {
      source = "hashicorp.com/edu/tmdb"
    }
  }
}

provider "tmdb" {}

data "tmdb_popular_movies" "movies" {}

output "popular_movies" {
  value = data.tmdb_popular_movies.movies
}