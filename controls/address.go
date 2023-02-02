package controls

import (
	"net/http"

	"strconv"

	"github.com/athunlal/config"

	"github.com/athunlal/models"
	"github.com/gin-gonic/gin"
)

//>>>>>>>>>>> Add addresses <<<<<<<<<<<<<<<<<<<<
func Addaddress(c *gin.Context) {
	id, err := strconv.Atoi(c.GetString("userid"))
	if err != nil {
		c.JSON(400, gin.H{
			"Error": "Error in string conversion",
		})
	}
	var userEnterData models.Address
	if c.Bind(&userEnterData) != nil {
		c.JSON(400, gin.H{
			"Error": "Error in Binding the JSON",
		})
	}
	db := config.DBconnect()
	db.Model(&models.Address{}).Where("userid = ?", id).Update("defaultadd", false)
	userEnterData.Userid = uint(id)
	result := db.Create(&userEnterData)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	db.Model(&userEnterData).Where("addressid = ?", userEnterData.Addressid).Update("defaultadd", true)
	c.JSON(200, gin.H{
		"Message": "Address added succesfully",
	})
}

//>>>>>>>>>>>>> show address <<<<<<<<<<<<<<<<<<<<<
func ShowAddress(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var userAddres models.Address

	db := config.DBconnect()
	result := db.Raw("SELECT * from addresses WHERE userid = ?", id).Scan(&userAddres)

	// if userAddres.User.ID == 0 {
	// 	c.JSON(404, gin.H{
	// 		"Message": "There is no address in this user id",
	// 	})
	// 	return
	// }
	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Address": gin.H{
			"Addressid":  userAddres.Addressid,
			"Userid":     userAddres.Userid,
			"Name":       userAddres.Name,
			"Phoneno":    userAddres.Phoneno,
			"Houseno":    userAddres.Houseno,
			"Area":       userAddres.Area,
			"Landmark":   userAddres.Landmark,
			"City":       userAddres.City,
			"Pincode":    userAddres.Pincode,
			"District":   userAddres.District,
			"State":      userAddres.State,
			"Country":    userAddres.Country,
			"Defaultadd": userAddres.Defaultadd,
		},
	})
}

//>>>>>>>>>>>>>> Edit Address <<<<<<<<<<<<<<<<<<<<<
func EditUserAddress(c *gin.Context) {
	id := c.Param("id")

	var userAddress models.Address
	if c.Bind(&userAddress) != nil {
		c.JSON(404, gin.H{
			"Error": "Error in binding JSON data",
		})
		return
	}
	str, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(500, gin.H{
			"Error": err,
		})
		return
	}
	userAddress.Userid = uint(str)
	db := config.DBconnect()

	result := db.Model(userAddress).Where("userid = ?", id).Updates(models.Address{
		Name:     userAddress.Name,
		Phoneno:  userAddress.Phoneno,
		Houseno:  userAddress.Houseno,
		Area:     userAddress.Area,
		Landmark: userAddress.Landmark,
		City:     userAddress.City,
		Pincode:  userAddress.Pincode,
		District: userAddress.District,
		State:    userAddress.State,
		Country:  userAddress.Country,
	})

	if result.Error != nil {
		c.JSON(404, gin.H{
			"Error": result.Error.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"Message":      "Successfully Updated the Address",
		"Updated data": userAddress,
	})
}
