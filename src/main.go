package main

func main() {
	router := SetupRouter()
	router.Run("localhost:8082")
}
