package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Device struct {
	PartInfo                  any                  `json:"partInfo,omitempty"`
	IsFirmwareUpdateMandatory bool                 `json:"isFirmwareUpdateMandatory,omitempty"`
	ProductType               ProductType          `json:"productType,omitempty"`
	SafeLocations             []SafeLocations      `json:"safeLocations,omitempty"`
	Owner                     string               `json:"owner,omitempty"`
	BatteryStatus             int                  `json:"batteryStatus,omitempty"`
	SerialNumber              string               `json:"serialNumber,omitempty"`
	LostModeMetadata          LostModeMetadata     `json:"lostModeMetadata,omitempty"`
	Capabilities              int                  `json:"capabilities,omitempty"`
	Identifier                string               `json:"identifier,omitempty"`
	Address                   Address              `json:"address,omitempty"`
	Location                  Location             `json:"location,omitempty"`
	ProductIdentifier         string               `json:"productIdentifier,omitempty"`
	IsAppleAudioAccessory     bool                 `json:"isAppleAudioAccessory,omitempty"`
	CrowdSourcedLocation      CrowdSourcedLocation `json:"crowdSourcedLocation,omitempty"`
	GroupIdentifier           any                  `json:"groupIdentifier,omitempty"`
	Role                      Role                 `json:"role,omitempty"`
	SystemVersion             string               `json:"systemVersion,omitempty"`
	Name                      string               `json:"name,omitempty"`
}
type ProductInformation struct {
	ManufacturerName  string `json:"manufacturerName,omitempty"`
	ModelName         string `json:"modelName,omitempty"`
	ProductIdentifier int    `json:"productIdentifier,omitempty"`
	VendorIdentifier  int    `json:"vendorIdentifier,omitempty"`
	AntennaPower      int    `json:"antennaPower,omitempty"`
}
type ProductType struct {
	Type               string             `json:"type,omitempty"`
	ProductInformation ProductInformation `json:"productInformation,omitempty"`
}
type Location struct {
	PositionType       string  `json:"positionType,omitempty"`
	VerticalAccuracy   int     `json:"verticalAccuracy,omitempty"`
	Longitude          float64 `json:"longitude,omitempty"`
	FloorLevel         int     `json:"floorLevel,omitempty"`
	IsInaccurate       bool    `json:"isInaccurate,omitempty"`
	IsOld              bool    `json:"isOld,omitempty"`
	HorizontalAccuracy float64 `json:"horizontalAccuracy,omitempty"`
	Latitude           float64 `json:"latitude,omitempty"`
	TimeStamp          int64   `json:"timeStamp,omitempty"`
	Altitude           int     `json:"altitude,omitempty"`
	LocationFinished   bool    `json:"locationFinished,omitempty"`
}
type Address struct {
	SubAdministrativeArea string   `json:"subAdministrativeArea,omitempty"`
	Label                 string   `json:"label,omitempty"`
	StreetAddress         string   `json:"streetAddress,omitempty"`
	CountryCode           string   `json:"countryCode,omitempty"`
	StateCode             any      `json:"stateCode,omitempty"`
	AdministrativeArea    string   `json:"administrativeArea,omitempty"`
	StreetName            string   `json:"streetName,omitempty"`
	FormattedAddressLines []string `json:"formattedAddressLines,omitempty"`
	MapItemFullAddress    string   `json:"mapItemFullAddress,omitempty"`
	FullThroroughfare     string   `json:"fullThroroughfare,omitempty"`
	AreaOfInterest        []any    `json:"areaOfInterest,omitempty"`
	Locality              string   `json:"locality,omitempty"`
	Country               string   `json:"country,omitempty"`
}
type SafeLocations struct {
	Type          int      `json:"type,omitempty"`
	ApprovalState int      `json:"approvalState,omitempty"`
	Name          any      `json:"name,omitempty"`
	Identifier    string   `json:"identifier,omitempty"`
	Location      Location `json:"location,omitempty"`
	Address       Address  `json:"address,omitempty"`
}
type LostModeMetadata struct {
	Email       string  `json:"email,omitempty"`
	Message     string  `json:"message,omitempty"`
	OwnerNumber string  `json:"ownerNumber,omitempty"`
	Timestamp   float64 `json:"timestamp,omitempty"`
}
type CrowdSourcedLocation struct {
	PositionType       string  `json:"positionType,omitempty"`
	VerticalAccuracy   int     `json:"verticalAccuracy,omitempty"`
	Longitude          float64 `json:"longitude,omitempty"`
	FloorLevel         int     `json:"floorLevel,omitempty"`
	IsInaccurate       bool    `json:"isInaccurate,omitempty"`
	IsOld              bool    `json:"isOld,omitempty"`
	HorizontalAccuracy float64 `json:"horizontalAccuracy,omitempty"`
	Latitude           float64 `json:"latitude,omitempty"`
	TimeStamp          int64   `json:"timeStamp,omitempty"`
	Altitude           int     `json:"altitude,omitempty"`
	LocationFinished   bool    `json:"locationFinished,omitempty"`
}
type Role struct {
	Name       string `json:"name,omitempty"`
	Emoji      string `json:"emoji,omitempty"`
	Identifier int    `json:"identifier,omitempty"`
}

const timeLayout = "2006-01-02 15:04:05"

func main() {
	pflag.String("device", "", "Device name to track")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.AutomaticEnv()

	fmt.Println("Starting to track...")
	fmt.Println("Please keep `Find My` app open on your device.")
	fmt.Println("Press Ctrl+C to stop tracking.")
	fmt.Println("")

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	jsonFile := home + "/Library/Caches/com.apple.findmy.fmipcore/Items.data"

	lastUpdate := make(map[string]time.Time)
	writers := make(map[string]*csv.Writer)

	// startt caffeinate
	cmd := exec.Command("caffeinate", "-di", "-w", strconv.Itoa(os.Getpid()))
	err = cmd.Start()

	stop := make(chan os.Signal, 1)

	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				time.Sleep(1 * time.Second)
				f, err := os.ReadFile(jsonFile)
				if err != nil {
					log.Fatal("Error opening file:", err)
				}

				var devices []Device
				err = json.Unmarshal(f, &devices)
				if err != nil {
					log.Fatalf("error decoding json: %v", err)
				}

				for _, d := range devices {
					rd := strings.ToLower(viper.GetString("device"))
					if rd != "" && strings.ToLower(d.Name) != rd {
						continue
					}

					ts := time.Unix(d.Location.TimeStamp/1000, d.Location.TimeStamp%1000)
					if lu, ok := lastUpdate[d.Identifier]; ok && lu.Before(ts) || lu.Equal(ts) {
						continue
					}
					lastUpdate[d.Identifier] = ts

					fmt.Printf("[%s] %s: %s\r\n", ts.Format(timeLayout), d.Name, d.Address.MapItemFullAddress)

					w, ok := writers[d.Identifier]
					if !ok {
						filename := strings.ReplaceAll(fmt.Sprintf("%s.csv", d.Name), " ", "_")
						filename = strings.ReplaceAll(filename, "’", "")
						var f *os.File
						//create a file csv writer
						new := false
						if os.IsNotExist(err) {
							f, err = os.Create(filename)
							if err != nil {
								log.Fatal(err)
							}
							new = true
						} else {
							f, err = os.OpenFile(filename, os.O_APPEND|os.O_RDWR, 0600)
							if err != nil {
								log.Fatal(err)
							}

							//get last line
							var line []byte
							var cursor int64 = 0
							stat, err := f.Stat()
							if err != nil {
								log.Fatal(err)
							}
							filesize := stat.Size()
							for {
								cursor -= 1
								_, err := f.Seek(cursor, io.SeekEnd)
								if err != nil {
									log.Fatal(err)
								}

								char := make([]byte, 1)
								_, err = f.Read(char)
								if err != nil {
									log.Fatal(err)
								}

								if cursor != -1 && (char[0] == 10 || char[0] == 13) { // stop if we find a line
									break
								}

								// prepend the new char to the line
								line = append(char, line...)

								if cursor == -filesize {
									break
								}
							}
							if len(line) > 0 {
								last := strings.Split(string(line), ",")
								if len(last) > 0 {
									lastUpdate[d.Identifier], err = time.ParseInLocation("2006-01-02 15:04:05", last[0], time.Now().Location())
									if err != nil {
										log.Fatal(err)
									}
								}
								if lastUpdate[d.Identifier].Equal(ts) {
									continue
								}
							}
						}
						defer f.Close()
						w = csv.NewWriter(f)
						writers[d.Identifier] = w
						if new {
							w.Write([]string{"time", "latitude", "longitude", "horizontalAccuracy", "street", "number", "city", "country"})
						}
					}
					w.Write([]string{
						ts.Format(timeLayout),
						fmt.Sprintf("%f", d.Location.Latitude),
						fmt.Sprintf("%f", d.Location.Longitude),
						fmt.Sprintf("%f", d.Location.HorizontalAccuracy),
						d.Address.StreetName,
						d.Address.StreetAddress,
						d.Address.Locality,
						d.Address.Country,
					})
					w.Flush()
					if err := w.Error(); err != nil {
						log.Fatal(err)
					}
				}
			}
		}
	}()

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	fmt.Println("Stopped by user")
	//stop caffeinate
	cmd.Process.Signal(syscall.SIGTERM)
}
