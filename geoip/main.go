package main

import (
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"net"
)

func main() {
	db, err := geoip2.Open("./GeoLite2-City.mmdb")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	ip := net.ParseIP("84.17.41.94")
	record, err := db.City(ip)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", record)

	//fmt.Printf("Portuguese (BR) city name: %v\n", record.City.Names["pt-BR"])
	//if len(record.Subdivisions) > 0 {
	//	fmt.Printf("English subdivision name: %v\n", record.Subdivisions[0].Names["en"])
	//}
	//fmt.Printf("Russian country name: %v\n", record.Country.Names["ru"])
	//fmt.Printf("ISO country code: %v\n", record.Country.IsoCode)
	//fmt.Printf("Time zone: %v\n", record.Location.TimeZone)
	//fmt.Printf("Coordinates: %v, %v\n", record.Location.Latitude, record.Location.Longitude)

}
