package main

func main() {
	store, err := NewSqliteStore()
	if err != nil {
		panic(err)
	}
	server := NewAPIServer(":3000", *store)
	server.Run()
}
