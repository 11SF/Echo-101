package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Goods struct {
	gorm.Model
	Name  string  `json:"name"`
	Price float32 `json:"price"`
}

type StoreHandler struct {
	DB *gorm.DB
}

func main() {
	e := echo.New()

	store_hlr := StoreHandler{}
	store_hlr.Initialize()

	//Router
	e.GET("/goods", store_hlr.GetGoods)
	e.POST("/goods", store_hlr.AddGoods)
	e.PUT("/goods/:id", store_hlr.UpdateGoods)
	e.DELETE("/goods/:id", store_hlr.DeleteGoods)

	e.Logger.Fatal(e.Start(":8080"))
}

func (h *StoreHandler) Initialize() {
	db, err := gorm.Open(sqlite.Open("store.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Goods{})
	h.DB = db
}

/*		Handler function		*/
//Create
func (h *StoreHandler) AddGoods(c echo.Context) (err error) {
	goods_req := new(Goods)

	//Binding Goods data from request body to Goods struct type
	if err = c.Bind(&goods_req); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	h.DB.Create(&goods_req)
	return c.NoContent(http.StatusOK)
}

//Read
func (h *StoreHandler) GetGoods(c echo.Context) error {
	goods := []Goods{}

	h.DB.Find(&goods)
	return c.JSON(http.StatusOK, goods)
}

//Update
func (h *StoreHandler) UpdateGoods(c echo.Context) error {
	req_id := c.Param("id")

	goods := Goods{}
	if err := c.Bind(&goods); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	h.DB.Where("id = ?", req_id).Updates(&goods)
	return c.NoContent(http.StatusOK)
}

//Delete
func (h *StoreHandler) DeleteGoods(c echo.Context) error {
	req_id := c.Param("id")
	goods := Goods{}

	h.DB.Where("id = ?", req_id).Delete(&goods)
	return c.NoContent(http.StatusOK)
}
