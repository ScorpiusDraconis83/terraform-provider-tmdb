terraform {
  required_providers {
    tmdb = {
      source = "hashicorp.com/edu/tmdb"
    }
  }
}

provider "tmdb" {}

data "tmdb_search" "results" {
   query = "Seven Samurai"
}

data "tmdb_movie" "movie" {
  id = data.tmdb_search.results.movies[0].id
}

output "movie" {
  value = data.tmdb_movie.movie
}