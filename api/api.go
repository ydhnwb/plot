package api

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ydhnwb/plot/entity"
)

type ParkingApi interface {
	Allocate(ctx *gin.Context)
	CreateParkingLot(ctx *gin.Context)
	Leave(ctx *gin.Context)
	Status(ctx *gin.Context)
	ByColor(ctx *gin.Context)
	BySlotIndex(ctx *gin.Context)
	ByNumber(ctx *gin.Context)
	Bulk(ctx *gin.Context)
}

type parking struct {
	Park entity.Park
}

func NewApi() ParkingApi {
	p := entity.Park{}
	return &parking{
		Park: p,
	}
}

func (c *parking) doByColor(color string) string {
	temps := []string{}

	for i := 0; i < len(c.Park.Parkings); i++ {
		if c.Park.Parkings[i].Color == color {
			temps = append(temps, c.Park.Parkings[i].Number)
		}
	}

	msg := strings.Join(temps, ", ")
	return msg
}

func (c *parking) doByNumber(number string) string {
	indexes := []string{}

	for i := 0; i < len(c.Park.Parkings); i++ {
		if c.Park.Parkings[i].Number == number {
			ix := fmt.Sprintf("%d", i+1)
			indexes = append(indexes, ix)
		}
	}

	if len(indexes) == 0 {
		return "Not found"
	}
	msg := strings.Join(indexes, ", ")
	return msg
}

func (c *parking) doByIndex(color string) string {
	indexes := []string{}

	for i := 0; i < len(c.Park.Parkings); i++ {
		if c.Park.Parkings[i].Color == color {
			ix := fmt.Sprintf("%d", i+1)
			indexes = append(indexes, ix)
		}
	}

	msg := strings.Join(indexes, ", ")
	return msg
}

func (c *parking) doCheckStatus() string {
	msgs := []string{}
	msgs = append(msgs, "Slot No. Registration No Colour")

	msg := "Slot No. Registration No Colour\n"

	for i := 0; i < len(c.Park.Parkings); i++ {
		if (c.Park.Parkings[i] != entity.Parking{}) {
			idx := i + 1
			msgs = append(msgs, fmt.Sprintf("%d %s %s", idx, c.Park.Parkings[i].Number, c.Park.Parkings[i].Color))
			msg += fmt.Sprintf("%d %s %s\n", idx, c.Park.Parkings[i].Number, c.Park.Parkings[i].Color)
		}

	}
	res := strings.Join(msgs, "\n")
	return res
}

func (c *parking) doCreateParkingLot(slot int) string {
	var parks = make([]entity.Parking, slot)
	c.Park.Slot = slot
	c.Park.Parkings = parks
	msg := fmt.Sprintf("Created a parking lot with %d slots", slot)
	return msg

}

func (c *parking) doAllocate(number string, color string) string {
	maxSlot := c.Park.Slot
	currIndex := -1

	for i := 0; i < maxSlot; i++ {
		if (c.Park.Parkings[i] == entity.Parking{}) {
			currIndex = i
			c.Park.Parkings[i] = entity.Parking{
				Number: number,
				Color:  color,
			}
			break
		}
	}

	idx := currIndex + 1
	return fmt.Sprintf("%d", idx)

}

func (c *parking) removeAt(position int) {
	c.Park.Parkings[position] = entity.Parking{}
}

func (c *parking) CreateParkingLot(ctx *gin.Context) {
	param := ctx.Param("slot")
	slot, err := strconv.Atoi(param)
	if err != nil {
		ctx.String(400, err.Error())
		return
	}

	c.doCreateParkingLot(slot)

	msg := fmt.Sprintf("Created a parking lot with %s slots", param)
	ctx.String(200, msg)
}

func (c *parking) Allocate(ctx *gin.Context) {
	number := ctx.Param("number")
	color := ctx.Param("color")

	idx := c.doAllocate(number, color)
	id, _ := strconv.ParseInt(idx, 10, 64)

	if id == 0 {
		ctx.String(200, "Sorry, parking lot is full")
		return
	}

	msg := "Allocated slot number: " + idx
	ctx.String(200, msg)

}

func (c *parking) Leave(ctx *gin.Context) {
	index := ctx.Param("index")
	i, _ := strconv.Atoi(index)

	c.removeAt(i - 1)
	msg := fmt.Sprintf("Slot number %d is free", i)
	ctx.String(200, msg)

}

func (c *parking) Status(ctx *gin.Context) {
	res := c.doCheckStatus()
	ctx.String(200, res)

}

func (c *parking) ByColor(ctx *gin.Context) {
	color := ctx.Param("color")
	msg := c.doByColor(color)
	ctx.String(200, msg)

}

func (c *parking) BySlotIndex(ctx *gin.Context) {
	color := ctx.Param("color")
	msg := c.doByIndex(color)
	ctx.String(200, msg)

}

func (c *parking) ByNumber(ctx *gin.Context) {
	number := ctx.Param("number")
	msg := c.doByNumber(number)
	ctx.String(200, msg)

}

func (c *parking) Bulk(ctx *gin.Context) {
	reqBody, _ := ioutil.ReadAll(ctx.Request.Body)
	commands := strings.Split(string(reqBody), "\n")

	res := []string{}

	for i := 0; i < len(commands); i++ {
		result := c.mapper(commands[i])
		if result != "" {
			res = append(res, result)
		}
	}

	response := strings.Join(res, "\n")

	response += "\n"
	ctx.String(200, response)
}

func (c *parking) getParam(command string) []string {
	return strings.Split(command, " ")
}

func (c *parking) mapper(command string) string {
	realCommand := strings.Split(command, " ")
	switch realCommand[0] {
	case "create_parking_lot":
		slot := c.getParam(command)
		slotI, _ := strconv.Atoi(slot[1])
		return c.doCreateParkingLot(slotI)

	case "park":
		param := c.getParam(command)
		idx := c.doAllocate(param[1], param[2])
		id, _ := strconv.ParseInt(idx, 10, 64)

		if id == 0 {
			return "Sorry, parking lot is full"
		}

		msg := "Allocated slot number: " + idx
		return msg

	case "leave":
		param := c.getParam(command)
		p, _ := strconv.ParseInt(param[1], 10, 64)
		c.removeAt(int(p) - 1)
		return fmt.Sprintf("Slot number %d is free", p)

	case "status":
		return c.doCheckStatus()

	case "registration_numbers_for_cars_with_colour":
		param := c.getParam(command)
		return c.doByColor(param[1])

	case "slot_numbers_for_cars_with_colour":
		param := c.getParam(command)
		return c.doByIndex(param[1])

	case "slot_number_for_registration_number":
		param := c.getParam(command)
		return c.doByNumber(param[1])
	default:
		return ""

	}

}
