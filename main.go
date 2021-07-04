package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ydhnwb/plot/api"
)

var (
	handler api.ParkingApi = api.NewApi()
)

func main() {
	server := gin.Default()

	server.POST("/create_parking_lot/:slot", handler.CreateParkingLot)
	server.POST("/leave/:index", handler.Leave)
	server.POST("/park/:number/:color", handler.Allocate)
	server.GET("/status", handler.Status)
	server.GET("/cars_registration_numbers/colour/:color", handler.ByColor)
	server.GET("/cars_slot/colour/:color", handler.BySlotIndex)
	server.GET("/slot_number/car_registration_number/:number", handler.ByNumber)

	server.POST("/bulk", handler.Bulk)

	server.Run()

}
