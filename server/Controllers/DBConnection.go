package Controllers

import (
	"fmt"
	"github.com/spf13/viper"
	"goST/Models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var StripeDB *gorm.DB

func ConnectStripeDB() {
	StripeDB, _ = gorm.Open(mysql.Open(getStripeDSNString()), &gorm.Config{
		PrepareStmt: true,
	})
	sqlDB, _ := StripeDB.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	_ = StripeDB.AutoMigrate(
		&Models.Account{},
		&Models.Balance{},
		&Models.Invoice{},
		&Models.Payout{},
	)
}

func getStripeDSNString() string {
	viper.SetConfigName("config.env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	dbAddress := viper.GetString("ADDRESS")
	dbPort := viper.GetString("PORT")
	dbUser := viper.GetString("USER")
	dbPassword := viper.GetString("PASSWORD")
	dbName := viper.GetString("DATABASE")
	dbCharset := viper.GetString("CHARSET")
	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbAddress + ":" + dbPort + ")/" + dbName + "?charset=" + dbCharset + "&parseTime=True&loc=Local"
	return dsn
}
