resource "a" "a" {
  c = "c"
  d = "d"
}

resource "a" "b" {
  # should move
}

resource "b" "a" {
  # should stay
}
