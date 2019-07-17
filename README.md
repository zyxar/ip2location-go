[![Go Report Card](https://goreportcard.com/badge/github.com/zyxar/ip2location-go)](https://goreportcard.com/report/github.com/zyxar/ip2location-go)


IP2Location Go Package
======================

This Go package provides a fast lookup of country, region, city, latitude, longitude, ZIP code, time zone, ISP, domain name, connection type, IDD code, area code, weather station code, station name, mcc, mnc, mobile brand, elevation, and usage type from IP address by using IP2Location database. This package uses a file based database available at IP2Location.com. This database simply contains IP blocks as keys, and other information such as country, region, city, latitude, longitude, ZIP code, time zone, ISP, domain name, connection type, IDD code, area code, weather station code, station name, mcc, mnc, mobile brand, elevation, and usage type as values. It supports both IP address in IPv4 and IPv6.

This package can be used in many types of projects such as:

 - select the geographically closest mirror
 - analyze your web server logs to determine the countries of your visitors
 - credit card fraud detection
 - software export controls
 - display native language and currency
 - prevent password sharing and abuse of service
 - geotargeting in advertisement

The database will be updated in monthly basis for the greater accuracy. Free LITE databases are available at https://lite.ip2location.com/ upon registration.

The paid databases are available at https://www.ip2location.com under Premium subscription package.


Installation
=======

```
go get github.com/zyxar/ip2location-go
```

Example
=======

```go
package main

import (
	"fmt"
	"os"
	"github.com/zyxar/ip2location-go"
)

func main() {
	db, err := ip2location.NewDB("./IPV6-COUNTRY-REGION-CITY-LATITUDE-LONGITUDE-ZIPCODE-TIMEZONE-ISP-DOMAIN-NETSPEED-AREACODE-WEATHER-MOBILE-ELEVATION-USAGETYPE.BIN")
	if err != nil {
		fmt.Println(err.Error())
		os.Eixt(1)
	}
	defer db.Close()

	fmt.Printf("api version: %s\n", ip2location.APIVersion())

	record, err := ip2location.GetAll("8.8.8.8")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("country_short: %s\n", record.CountryShort)
	fmt.Printf("country_long: %s\n", record.CountryLong)
	fmt.Printf("region: %s\n", record.Region)
	fmt.Printf("city: %s\n", record.City)
	fmt.Printf("isp: %s\n", record.ISP)
	fmt.Printf("latitude: %f\n", record.Latitude)
	fmt.Printf("longitude: %f\n", record.Longitude)
}
```

Dependencies
============

The complete database is available at https://www.ip2location.com under subscription package.


IPv4 BIN vs IPv6 BIN
====================

Use the IPv4 BIN file if you just need to query IPv4 addresses.
Use the IPv6 BIN file if you need to query BOTH IPv4 and IPv6 addresses.


Copyright
=========

Copyright (C) 2018 by IP2Location.com, support@ip2location.com
