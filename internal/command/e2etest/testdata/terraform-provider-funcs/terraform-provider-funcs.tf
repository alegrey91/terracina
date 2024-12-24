# This test fixture is here primarily just to make sure that the
# terracina.io/builtin/terracina functions remain available for use. The
# actual behavior of these functions is the responsibility of
# ./internal/builtin/providers/terracina, and so it has more detailed tests
# whereas this one is focused largely just on whether these functions are
# callable at all.

terracina {
  required_providers {
    terracina = {
      source = "terracina.io/builtin/terracina"
    }
  }
}

output "tfvarsencode" {
  value = provider::terracina::encode_tfvars({
    a = "👋"
    b = "🐝"
    c = "👓"
  })
}

output "tfvarsdecode" {
  value = provider::terracina::decode_tfvars(
    <<-EOT
      boop = "👃"
      baaa = "🐑"
    EOT
  )
}

output "exprencode" {
  value = provider::terracina::encode_expr([1, 2, 3])
}
