package controls

import (
	"path/filepath"
	"strconv"

	"github.com/athunlal/config"
	"github.com/athunlal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/razorpay/razorpay-go"
)

//>>>>>> Add brand <<<<<<<<<<<<<<<<<<<<
func AddBrands(c *gin.Context) {
	var addbrand models.Brand
	if c.Bind(&addbrand) != nil {
		c.JSON(400, gin.H{
			"Error": "Could not bind JSON data",
		})
		return
	}
	DB := config.DBconnect()
	result := DB.Create(&addbrand)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message":       "New Brand added Successfully",
		"Brand details": addbrand,
	})
}

// >>>>>>>>>> view brand <<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func ViewBrand(c *gin.Context) {
	var brandData []models.Brand
	db := config.DBconnect()
	result := db.Find(&brandData)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"Status":  "False",
			"Message": "Could not retrieve the brands",
		})
		return
	}
	c.JSON(200, gin.H{
		"Brands data": brandData,
	})
}

//>>>>>>>>>>>>>> Edit brand <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func EditBrand(c *gin.Context) {
	bid := c.Param("id")
	id, err := strconv.Atoi(bid)
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}
	var editbrands models.Brand
	if c.Bind(&editbrands) != nil {
		c.JSON(400, gin.H{
			"Error": "Error in binding the JSON data",
		})
		return
	}
	editbrands.ID = uint(id)
	DB := config.DBconnect()

	result := DB.Model(&editbrands).Updates(models.Brand{
		Brandname: editbrands.Brandname,
	})

	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "Successfully updated the Brand",
	})
}

//>>>>>>>>>>>> Add to cart <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func AddToCart(c *gin.Context) {
	type data struct {
		Product_id uint
		Quantity   uint
	}
	var bindData data
	var cartdata models.Cart
	var productdata models.Product
	if c.Bind(&bindData) != nil {
		c.JSON(400, gin.H{
			"Error": "Could not bind the JSON data",
		})
		return
	}

	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}

	db := config.DBconnect()
	var count int64

	//fetching the table products for checking stocks
	db.Table("products").Select("stock, price").Where("productid = ?", bindData.Product_id).Scan(&productdata)
	if bindData.Quantity >= productdata.Stock {
		c.JSON(404, gin.H{
			"Message": "Out of Stock",
		})
		return
	}

	//fetching the table carts for checking the product_id is exist
	db.Model(&cartdata).Where("product_id = ?", bindData.Product_id).Count(&count)
	if count > 0 && id == int(cartdata.Userid) {
		var sum uint

		//fetching the quantity form carts
		db.Table("carts").Where("product_id = ?", bindData.Product_id).Select("SUM(quantity)").Row().Scan(&sum)
		totalQuantity := sum + bindData.Quantity

		//updating the quatity to the carts
		db.Model(&cartdata).Where("product_id = ?", bindData.Product_id).Update("quantity", totalQuantity)
		c.JSON(200, gin.H{
			"Message": "Quantity added Successfully",
		})
		return
	}
	totalprice := productdata.Price * bindData.Quantity
	cartitems := models.Cart{
		Product_id: bindData.Product_id,
		Quantity:   bindData.Quantity,
		Price:      productdata.Price,
		Totalprice: totalprice,
		Userid:     uint(id),
	}
	result := db.Create(&cartitems)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "Added to the Cart Successfull",
	})
}

//>>>>>>>>>>>>>>> View Products <<<<<<<<<<<<<<<<<<<<<<<

func ViewCart(c *gin.Context) {
	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}
	type cartdata struct {
		Productname string
		Quantity    uint
		Totalprice  uint
		Image       string
		Price       string
	}
	var datas []cartdata
	db := config.DBconnect()
	result := db.Table("carts").Select("products.productname, carts.quantity, carts.price, carts.totalprice").Joins("INNER JOIN products ON products.productid=carts.product_id").Where("userid = ?", id).Scan(&datas)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	if datas != nil {
		c.JSON(200, gin.H{
			"Cart Items": datas,
		})
	} else {
		c.JSON(404, gin.H{
			"Message": "Cart is empty",
		})
	}
}

//>>>>>>>>>>>>>Remove cart <<<<<<<<<<<<<<<<<<<<<
func DeleteCart(c *gin.Context) {
	id := c.Param("id")
	userid, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}

	db := config.DBconnect()
	result := db.Exec("delete from carts where id= ? AND userid = ?", id, userid)
	count := result.RowsAffected
	if count == 0 {
		c.JSON(400, gin.H{
			"Message": "Cart not exist",
		})
		return
	}
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Cart Items": "Delete successfully",
	})
}

//>>>>>>>>> Add Image <<<<<<<<<<<<<<<<<<<<<<<

func AddImages(c *gin.Context) {
	imagepath, _ := c.FormFile("image")
	extension := filepath.Ext(imagepath.Filename)
	image := uuid.New().String() + extension
	c.SaveUploadedFile(imagepath, "./public/images"+image)
	pidconv := c.PostForm("product_id")
	pid, _ := strconv.Atoi(pidconv)
	db := config.DBconnect()

	var product models.Product
	result := db.First(&product, pid)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	imagedata := models.Image{
		Image:      image,
		Product_id: uint(pid),
	}
	result = db.Create(&imagedata)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "Image Added Successfully",
	})
}

//>>>>>>>>>>>>>>>> Payment <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func Payment(c *gin.Context) {
	//fetching user id from token
	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}
	type data struct {
		Method string
	}
	var bindData data
	var cartDate models.Cart

	//binding the data from posman
	if c.Bind(&bindData) != nil {
		c.JSON(400, gin.H{
			"Error": "Could not bind the JSON data",
		})
		return
	}

	db := config.DBconnect()

	//fetching the data from the table carts by id
	result := db.First(&cartDate, "userid = ?", id)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	//fetching the total amount from the table carts
	var total_amount float64
	result = db.Table("carts").Where("userid = ?", id).Select("SUM(totalprice)").Scan(&total_amount)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}

	if bindData.Method == "COD" {
		if result.Error != nil {
			c.JSON(400, gin.H{
				"Error": result.Error.Error(),
			})
			return
		}

		paymentData := models.Payment{
			PaymentMethod: bindData.Method,
			Totalamount:   uint(total_amount),
			User_id:       uint(id),
		}

		result = db.Create(&paymentData)
		if result.Error != nil {
			c.JSON(400, gin.H{
				"Error": result.Error.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"Message": "Payment Method COD",
			"Status":  "Completed",
		})

	} else if bindData.Method == "UPI" {
		//razor pay code
		client := razorpay.NewClient("rzp_test_mCL1zwPhJbeuND", "qUeHjny0jl14sphKqOFpyq9M")

		data := map[string]interface{}{
			"amount":   cartDate.Totalprice,
			"currency": "INR",
			"receipt":  "some_receipt_id",
		}
		body, err := client.Order.Create(data, nil)

		if err != nil {
			c.JSON(500, gin.H{
				"Error": err.Error(),
			})
			return
		} else {
			paymentData := models.Payment{
				PaymentMethod: bindData.Method,
				Totalamount:   uint(total_amount),
				User_id:       uint(id),
			}

			result = db.Create(&paymentData)
			if result.Error != nil {
				c.JSON(400, gin.H{
					"Error": result.Error.Error(),
				})
				return
			}
			orderID := body["id"].(string)
			amount := body["amount"].(float64)
			c.JSON(200, gin.H{
				"Order ID": orderID,
				"Amount":   amount,
				"Message":  "Payment Method UPI",
				"Status":   "Completed",
			})
		}
	} else {
		c.JSON(400, gin.H{
			"Error": "Payment field",
		})
		return
	}
}
