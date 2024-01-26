table "customers" {
  schema = schema.main
  column "id" {
    null = true
    type = integer
  }
  column "first_name" {
    null = false
    type = text
  }
  column "last_name" {
    null = false
    type = text
  }
  primary_key {
    columns = [column.id]
  }
}
schema "main" {
}
