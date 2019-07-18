package ip2location

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math/big"
	"net"
	"os"
	"strconv"
)

type Record struct {
	CountryShort       string
	CountryLong        string
	Region             string
	City               string
	ISP                string
	Latitude           float32
	Longitude          float32
	Domain             string
	ZipCode            string
	TimeZone           string
	NetSpeed           string
	IddCode            string
	AreaCode           string
	WeatherStationCode string
	WeatherStationName string
	MobileCountryCode  string
	MobileNetworkCode  string
	MobileBrand        string
	Elevation          float32
	UsageType          string
}

var (
	countryPosition            = [25]uint8{0, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2}
	regionPosition             = [25]uint8{0, 0, 0, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3}
	cityPosition               = [25]uint8{0, 0, 0, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4}
	ispPosition                = [25]uint8{0, 0, 3, 0, 5, 0, 7, 5, 7, 0, 8, 0, 9, 0, 9, 0, 9, 0, 9, 7, 9, 0, 9, 7, 9}
	latitudePosition           = [25]uint8{0, 0, 0, 0, 0, 5, 5, 0, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5}
	longitudePosition          = [25]uint8{0, 0, 0, 0, 0, 6, 6, 0, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6}
	domainPosition             = [25]uint8{0, 0, 0, 0, 0, 0, 0, 6, 8, 0, 9, 0, 10, 0, 10, 0, 10, 0, 10, 8, 10, 0, 10, 8, 10}
	zipcodePosition            = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 7, 7, 7, 0, 7, 7, 7, 0, 7, 0, 7, 7, 7, 0, 7}
	timezonePosition           = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 8, 7, 8, 8, 8, 7, 8, 0, 8, 8, 8, 0, 8}
	netspeedPosition           = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 11, 0, 11, 8, 11, 0, 11, 0, 11, 0, 11}
	iddcodePosition            = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 12, 0, 12, 0, 12, 9, 12, 0, 12}
	areacodePosition           = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10, 13, 0, 13, 0, 13, 10, 13, 0, 13}
	weatherstationcodePosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 14, 0, 14, 0, 14, 0, 14}
	weatherstationnamePosition = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10, 15, 0, 15, 0, 15, 0, 15}
	mccPosition                = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9, 16, 0, 16, 9, 16}
	mncPosition                = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10, 17, 0, 17, 10, 17}
	mobilebrandPosition        = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 11, 18, 0, 18, 11, 18}
	elevationPosition          = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 11, 19, 0, 19}
	usagetypePosition          = [25]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 12, 20}
)

var (
	maxIPV4Range = big.NewInt(4294967295)
	maxIPV6Range = big.NewInt(0)

	ErrInvalidAddress = errors.New("Invalid IP address")
	ErrInvalidFile    = errors.New("Invalid database file")
	ErrNotSupported   = errors.New("Unsupported feature for selected data file")
)

const (
	apiVersion string = "8.0.3"
)

const (
	countryshort uint32 = 1 << iota
	countrylong
	region
	city
	isp
	latitude
	longitude
	domain
	zipcode
	timezone
	netspeed
	iddcode
	areacode
	weatherstationcode
	weatherstationname
	mcc
	mnc
	mobilebrand
	elevation
	usagetype

	all uint32 = countryshort | countrylong | region | city | isp | latitude | longitude | domain | zipcode | timezone | netspeed | iddcode | areacode | weatherstationcode | weatherstationname | mcc | mnc | mobilebrand | elevation | usagetype
)

type DB struct {
	f    *os.File
	meta ip2locationmeta

	countryPositionOffset            uint32
	regionPositionOffset             uint32
	cityPositionOffset               uint32
	ispPositionOffset                uint32
	domainPositionOffset             uint32
	zipcodePositionOffset            uint32
	latitudePositionOffset           uint32
	longitudePositionOffset          uint32
	timezonePositionOffset           uint32
	netspeedPositionOffset           uint32
	iddcodePositionOffset            uint32
	areacodePositionOffset           uint32
	weatherstationcodePositionOffset uint32
	weatherstationnamePositionOffset uint32
	mccPositionOffset                uint32
	mncPositionOffset                uint32
	mobilebrandPositionOffset        uint32
	elevationPositionOffset          uint32
	usagetypePositionOffset          uint32

	countryEnabled            bool
	regionEnabled             bool
	cityEnabled               bool
	ispEnabled                bool
	domainEnabled             bool
	zipcodeEnabled            bool
	latitudeEnabled           bool
	longitudeEnabled          bool
	timezoneEnabled           bool
	netspeedEnabled           bool
	iddcodeEnabled            bool
	areacodeEnabled           bool
	weatherstationcodeEnabled bool
	weatherstationnameEnabled bool
	mccEnabled                bool
	mncEnabled                bool
	mobilebrandEnabled        bool
	elevationEnabled          bool
	usagetypeEnabled          bool
}

// get IP type and calculate IP number; calculates index too if exists
func checkip(meta *ip2locationmeta, ip string) (iptype uint32, ipnum *big.Int, ipindex uint32) {
	iptype = 0
	ipnum = big.NewInt(0)
	ipnumtmp := big.NewInt(0)
	ipindex = 0
	ipaddress := net.ParseIP(ip)

	if ipaddress != nil {
		v4 := ipaddress.To4()

		if v4 != nil {
			iptype = 4
			ipnum.SetBytes(v4)
		} else {
			v6 := ipaddress.To16()

			if v6 != nil {
				iptype = 6
				ipnum.SetBytes(v6)
			}
		}
	}
	if iptype == 4 {
		if meta.ipv4indexbaseaddr > 0 {
			ipnumtmp.Rsh(ipnum, 16)
			ipnumtmp.Lsh(ipnumtmp, 3)
			ipindex = uint32(ipnumtmp.Add(ipnumtmp, big.NewInt(int64(meta.ipv4indexbaseaddr))).Uint64())
		}
	} else if iptype == 6 {
		if meta.ipv6indexbaseaddr > 0 {
			ipnumtmp.Rsh(ipnum, 112)
			ipnumtmp.Lsh(ipnumtmp, 3)
			ipindex = uint32(ipnumtmp.Add(ipnumtmp, big.NewInt(int64(meta.ipv6indexbaseaddr))).Uint64())
		}
	}
	return
}

// read byte
func readuint8(f *os.File, pos int64) (uint8, error) {
	var retval uint8
	data := make([]byte, 1)
	_, err := f.ReadAt(data, pos-1)
	if err != nil {
		return 0, err
	}
	retval = data[0]
	return retval, nil
}

// read unsigned 32-bit integer
func readuint32(f *os.File, pos uint32) (uint32, error) {
	pos2 := int64(pos)
	data := make([]byte, 4)
	_, err := f.ReadAt(data, pos2-1)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(data), nil
}

// read unsigned 128-bit integer
func readuint128(f *os.File, pos uint32) (*big.Int, error) {
	pos2 := int64(pos)
	retval := big.NewInt(0)
	data := make([]byte, 16)
	_, err := f.ReadAt(data, pos2-1)
	if err != nil {
		return nil, err
	}

	// little endian to big endian
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
	retval.SetBytes(data)
	return retval, nil
}

// read string
func readstr(f *os.File, pos uint32) (string, error) {
	pos2 := int64(pos)
	var retval string
	lenbyte := make([]byte, 1)
	_, err := f.ReadAt(lenbyte, pos2)
	if err != nil {
		return "", err
	}
	strlen := lenbyte[0]
	data := make([]byte, strlen)
	_, err = f.ReadAt(data, pos2+1)
	if err != nil {
		return "", err
	}
	retval = string(data[:strlen])
	return retval, nil
}

// read float
func readfloat(f *os.File, pos uint32) (float32, error) {
	pos2 := int64(pos)
	var retval float32
	data := make([]byte, 4)
	_, err := f.ReadAt(data, pos2-1)
	if err != nil {
		return 0, err
	}
	buf := bytes.NewReader(data)
	err = binary.Read(buf, binary.LittleEndian, &retval)
	if err != nil {
		return 0, err
	}
	return retval, nil
}

func init() {
	maxIPV6Range.SetString("340282366920938463463374607431768211455", 10)
}

// initialize the component with the database path
func NewDB(dbpath string) (*DB, error) {
	f, err := os.Open(dbpath)
	if err != nil {
		return nil, err
	}
	var meta ip2locationmeta
	meta.databasetype, err = readuint8(f, 1)
	if err != nil {
		return nil, ErrInvalidFile
	}
	meta.databasecolumn, err = readuint8(f, 2)
	if err != nil {
		return nil, ErrInvalidFile
	}
	meta.databaseyear, err = readuint8(f, 3)
	if err != nil {
		return nil, ErrInvalidFile
	}
	meta.databasemonth, err = readuint8(f, 4)
	if err != nil {
		return nil, ErrInvalidFile
	}
	meta.databaseday, err = readuint8(f, 5)
	if err != nil {
		return nil, ErrInvalidFile
	}
	meta.ipv4databasecount, err = readuint32(f, 6)
	if err != nil {
		return nil, ErrInvalidFile
	}
	meta.ipv4databaseaddr, err = readuint32(f, 10)
	if err != nil {
		return nil, ErrInvalidFile
	}
	meta.ipv6databasecount, err = readuint32(f, 14)
	if err != nil {
		return nil, ErrInvalidFile
	}
	meta.ipv6databaseaddr, err = readuint32(f, 18)
	if err != nil {
		return nil, ErrInvalidFile
	}
	meta.ipv4indexbaseaddr, err = readuint32(f, 22)
	if err != nil {
		return nil, ErrInvalidFile
	}
	meta.ipv6indexbaseaddr, err = readuint32(f, 26)
	if err != nil {
		return nil, ErrInvalidFile
	}
	meta.ipv4columnsize = uint32(meta.databasecolumn << 2)              // 4 bytes each column
	meta.ipv6columnsize = uint32(16 + ((meta.databasecolumn - 1) << 2)) // 4 bytes each column, except IPFrom column which is 16 bytes

	dbt := meta.databasetype
	db := &DB{f: f, meta: meta}

	// since both IPv4 and IPv6 use 4 bytes for the below columns, can just do it once here
	if countryPosition[dbt] != 0 {
		db.countryPositionOffset = uint32(countryPosition[dbt]-1) << 2
		db.countryEnabled = true
	}
	if regionPosition[dbt] != 0 {
		db.regionPositionOffset = uint32(regionPosition[dbt]-1) << 2
		db.regionEnabled = true
	}
	if cityPosition[dbt] != 0 {
		db.cityPositionOffset = uint32(cityPosition[dbt]-1) << 2
		db.cityEnabled = true
	}
	if ispPosition[dbt] != 0 {
		db.ispPositionOffset = uint32(ispPosition[dbt]-1) << 2
		db.ispEnabled = true
	}
	if domainPosition[dbt] != 0 {
		db.domainPositionOffset = uint32(domainPosition[dbt]-1) << 2
		db.domainEnabled = true
	}
	if zipcodePosition[dbt] != 0 {
		db.zipcodePositionOffset = uint32(zipcodePosition[dbt]-1) << 2
		db.zipcodeEnabled = true
	}
	if latitudePosition[dbt] != 0 {
		db.latitudePositionOffset = uint32(latitudePosition[dbt]-1) << 2
		db.latitudeEnabled = true
	}
	if longitudePosition[dbt] != 0 {
		db.longitudePositionOffset = uint32(longitudePosition[dbt]-1) << 2
		db.longitudeEnabled = true
	}
	if timezonePosition[dbt] != 0 {
		db.timezonePositionOffset = uint32(timezonePosition[dbt]-1) << 2
		db.timezoneEnabled = true
	}
	if netspeedPosition[dbt] != 0 {
		db.netspeedPositionOffset = uint32(netspeedPosition[dbt]-1) << 2
		db.netspeedEnabled = true
	}
	if iddcodePosition[dbt] != 0 {
		db.iddcodePositionOffset = uint32(iddcodePosition[dbt]-1) << 2
		db.iddcodeEnabled = true
	}
	if areacodePosition[dbt] != 0 {
		db.areacodePositionOffset = uint32(areacodePosition[dbt]-1) << 2
		db.areacodeEnabled = true
	}
	if weatherstationcodePosition[dbt] != 0 {
		db.weatherstationcodePositionOffset = uint32(weatherstationcodePosition[dbt]-1) << 2
		db.weatherstationcodeEnabled = true
	}
	if weatherstationnamePosition[dbt] != 0 {
		db.weatherstationnamePositionOffset = uint32(weatherstationnamePosition[dbt]-1) << 2
		db.weatherstationnameEnabled = true
	}
	if mccPosition[dbt] != 0 {
		db.mccPositionOffset = uint32(mccPosition[dbt]-1) << 2
		db.mccEnabled = true
	}
	if mncPosition[dbt] != 0 {
		db.mncPositionOffset = uint32(mncPosition[dbt]-1) << 2
		db.mncEnabled = true
	}
	if mobilebrandPosition[dbt] != 0 {
		db.mobilebrandPositionOffset = uint32(mobilebrandPosition[dbt]-1) << 2
		db.mobilebrandEnabled = true
	}
	if elevationPosition[dbt] != 0 {
		db.elevationPositionOffset = uint32(elevationPosition[dbt]-1) << 2
		db.elevationEnabled = true
	}
	if usagetypePosition[dbt] != 0 {
		db.usagetypePositionOffset = uint32(usagetypePosition[dbt]-1) << 2
		db.usagetypeEnabled = true
	}

	return db, nil
}

// get api version
func APIVersion() string {
	return apiVersion
}

func (db *DB) Close() error { return db.f.Close() }

// get all fields
func (db *DB) GetAll(ip string) (*Record, error) { return db.query(ip, all) }

// get country code
func (db *DB) GetCountryShort(ip string) (*Record, error) { return db.query(ip, countryshort) }

// get country name
func (db *DB) GetCountryLong(ip string) (*Record, error) { return db.query(ip, countrylong) }

// get region
func (db *DB) GetRegion(ip string) (*Record, error) { return db.query(ip, region) }

// get city
func (db *DB) GetCity(ip string) (*Record, error) { return db.query(ip, city) }

// get isp
func (db *DB) GetIsp(ip string) (*Record, error) { return db.query(ip, isp) }

// get latitude
func (db *DB) GetLatitude(ip string) (*Record, error) { return db.query(ip, latitude) }

// get longitude
func (db *DB) GetLongitude(ip string) (*Record, error) { return db.query(ip, longitude) }

// get domain
func (db *DB) GetDomain(ip string) (*Record, error) { return db.query(ip, domain) }

// get zip code
func (db *DB) GetZipcode(ip string) (*Record, error) { return db.query(ip, zipcode) }

// get time zone
func (db *DB) GetTimezone(ip string) (*Record, error) { return db.query(ip, timezone) }

// get net speed
func (db *DB) GetNetSpeed(ip string) (*Record, error) { return db.query(ip, netspeed) }

// get idd code
func (db *DB) GetIddCode(ip string) (*Record, error) { return db.query(ip, iddcode) }

// get area code
func (db *DB) GetAreaCode(ip string) (*Record, error) { return db.query(ip, areacode) }

// get weather station code
func (db *DB) GetWeatherStationCode(ip string) (*Record, error) {
	return db.query(ip, weatherstationcode)
}

// get weather station name
func (db *DB) GetWeatherStationName(ip string) (*Record, error) {
	return db.query(ip, weatherstationname)
}

// get mobile country code
func (db *DB) GetMobileCountryCode(ip string) (*Record, error) { return db.query(ip, mcc) }

// get mobile network code
func (db *DB) GetMobileNetworkCode(ip string) (*Record, error) { return db.query(ip, mnc) }

// get mobile carrier brand
func (db *DB) GetMobileBrand(ip string) (*Record, error) { return db.query(ip, mobilebrand) }

// get elevation
func (db *DB) GetElevation(ip string) (*Record, error) { return db.query(ip, elevation) }

// get usage type
func (db *DB) GetUsageType(ip string) (*Record, error) { return db.query(ip, usagetype) }

// main query
func (db *DB) query(ip string, mode uint32) (*Record, error) {
	// check IP type and return IP number & index (if exists)
	iptype, ipno, ipindex := checkip(&db.meta, ip)
	if iptype == 0 {
		return nil, ErrInvalidAddress
	}

	var colsize uint32
	var baseaddr uint32
	var low uint32
	var high uint32
	var mid uint32
	var rowoffset uint32
	var rowoffset2 uint32
	var x Record
	ipfrom := big.NewInt(0)
	ipto := big.NewInt(0)
	maxip := big.NewInt(0)

	if iptype == 4 {
		baseaddr = db.meta.ipv4databaseaddr
		high = db.meta.ipv4databasecount
		maxip = maxIPV4Range
		colsize = db.meta.ipv4columnsize
	} else {
		baseaddr = db.meta.ipv6databaseaddr
		high = db.meta.ipv6databasecount
		maxip = maxIPV6Range
		colsize = db.meta.ipv6columnsize
	}

	// reading index
	if ipindex > 0 {
		low, _ = readuint32(db.f, ipindex)
		high, _ = readuint32(db.f, ipindex+4)
	}

	if ipno.Cmp(maxip) >= 0 {
		ipno = ipno.Sub(ipno, big.NewInt(1))
	}

	for low <= high {
		mid = ((low + high) >> 1)
		rowoffset = baseaddr + (mid * colsize)
		rowoffset2 = rowoffset + colsize

		if iptype == 4 {
			val, _ := readuint32(db.f, rowoffset)
			ipfrom = big.NewInt(int64(val))
			val, _ = readuint32(db.f, rowoffset2)
			ipto = big.NewInt(int64(val))
		} else {
			ipfrom, _ = readuint128(db.f, rowoffset)
			ipto, _ = readuint128(db.f, rowoffset2)
		}

		if ipno.Cmp(ipfrom) >= 0 && ipno.Cmp(ipto) < 0 {
			if iptype == 6 {
				rowoffset = rowoffset + 12 // coz below is assuming all columns are 4 bytes, so got 12 left to go to make 16 bytes total
			}

			if mode&countryshort == 1 && db.countryEnabled {
				val, _ := readuint32(db.f, rowoffset+db.countryPositionOffset)
				x.CountryShort, _ = readstr(db.f, val)
			}

			if mode&countrylong != 0 && db.countryEnabled {
				val, _ := readuint32(db.f, rowoffset+db.countryPositionOffset)
				x.CountryLong, _ = readstr(db.f, val+3)
			}

			if mode&region != 0 && db.regionEnabled {
				val, _ := readuint32(db.f, rowoffset+db.regionPositionOffset)
				x.Region, _ = readstr(db.f, val)
			}

			if mode&city != 0 && db.cityEnabled {
				val, _ := readuint32(db.f, rowoffset+db.cityPositionOffset)
				x.City, _ = readstr(db.f, val)
			}

			if mode&isp != 0 && db.ispEnabled {
				val, _ := readuint32(db.f, rowoffset+db.ispPositionOffset)
				x.ISP, _ = readstr(db.f, val)
			}

			if mode&latitude != 0 && db.latitudeEnabled {
				x.Latitude, _ = readfloat(db.f, rowoffset+db.latitudePositionOffset)
			}

			if mode&longitude != 0 && db.longitudeEnabled {
				x.Longitude, _ = readfloat(db.f, rowoffset+db.longitudePositionOffset)
			}

			if mode&domain != 0 && db.domainEnabled {
				val, _ := readuint32(db.f, rowoffset+db.domainPositionOffset)
				x.Domain, _ = readstr(db.f, val)
			}

			if mode&zipcode != 0 && db.zipcodeEnabled {
				val, _ := readuint32(db.f, rowoffset+db.zipcodePositionOffset)
				x.ZipCode, _ = readstr(db.f, val)
			}

			if mode&timezone != 0 && db.timezoneEnabled {
				val, _ := readuint32(db.f, rowoffset+db.timezonePositionOffset)
				x.TimeZone, _ = readstr(db.f, val)
			}

			if mode&netspeed != 0 && db.netspeedEnabled {
				val, _ := readuint32(db.f, rowoffset+db.netspeedPositionOffset)
				x.NetSpeed, _ = readstr(db.f, val)
			}

			if mode&iddcode != 0 && db.iddcodeEnabled {
				val, _ := readuint32(db.f, rowoffset+db.iddcodePositionOffset)
				x.IddCode, _ = readstr(db.f, val)
			}

			if mode&areacode != 0 && db.areacodeEnabled {
				val, _ := readuint32(db.f, rowoffset+db.areacodePositionOffset)
				x.AreaCode, _ = readstr(db.f, val)
			}

			if mode&weatherstationcode != 0 && db.weatherstationcodeEnabled {
				val, _ := readuint32(db.f, rowoffset+db.weatherstationcodePositionOffset)
				x.WeatherStationCode, _ = readstr(db.f, val)
			}

			if mode&weatherstationname != 0 && db.weatherstationnameEnabled {
				val, _ := readuint32(db.f, rowoffset+db.weatherstationnamePositionOffset)
				x.WeatherStationName, _ = readstr(db.f, val)
			}

			if mode&mcc != 0 && db.mccEnabled {
				val, _ := readuint32(db.f, rowoffset+db.mccPositionOffset)
				x.MobileCountryCode, _ = readstr(db.f, val)
			}

			if mode&mnc != 0 && db.mncEnabled {
				val, _ := readuint32(db.f, rowoffset+db.mncPositionOffset)
				x.MobileNetworkCode, _ = readstr(db.f, val)
			}

			if mode&mobilebrand != 0 && db.mobilebrandEnabled {
				val, _ := readuint32(db.f, rowoffset+db.mobilebrandPositionOffset)
				x.MobileBrand, _ = readstr(db.f, val)
			}

			if mode&elevation != 0 && db.elevationEnabled {
				val, _ := readuint32(db.f, rowoffset+db.elevationPositionOffset)
				vals, _ := readstr(db.f, val)
				f, _ := strconv.ParseFloat(vals, 32)
				x.Elevation = float32(f)
			}

			if mode&usagetype != 0 && db.usagetypeEnabled {
				val, _ := readuint32(db.f, rowoffset+db.usagetypePositionOffset)
				x.UsageType, _ = readstr(db.f, val)
			}

			return &x, nil
		} else {
			if ipno.Cmp(ipfrom) < 0 {
				high = mid - 1
			} else {
				low = mid + 1
			}
		}
	}
	return &x, nil
}

type ip2locationmeta struct {
	databasetype      uint8
	databasecolumn    uint8
	databaseday       uint8
	databasemonth     uint8
	databaseyear      uint8
	ipv4databasecount uint32
	ipv4databaseaddr  uint32
	ipv6databasecount uint32
	ipv6databaseaddr  uint32
	ipv4indexbaseaddr uint32
	ipv6indexbaseaddr uint32
	ipv4columnsize    uint32
	ipv6columnsize    uint32
}
