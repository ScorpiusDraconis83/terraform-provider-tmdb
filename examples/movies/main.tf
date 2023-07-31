terraform {
  required_providers {
    tmdb = {
      source = "hashicorp.com/edu/tmdb"
    }
  }
}

provider "tmdb" {}

data "tmdb_popular_movies" "edu" {}

output "edu_coffees" {
  value = data.tmdb_popular_movies.edu
}