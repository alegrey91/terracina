resource "terracina_data" "a" {
}

resource "terracina_data" "b" {
  input = terracina_data.a.id
}

resource "terracina_data" "c" {
  triggers_replace = terracina_data.b
}

resource "terracina_data" "d" {
  input = [ terracina_data.b, terracina_data.c ]
}

output "d" {
  value = terracina_data.d
}
