package enums

type FlightClass int

const (
	Economy  FlightClass = 0
	Business FlightClass = 1
)

func (flightClass FlightClass) String() string {
	if flightClass == 0 {
		return "Economy"
	} else {
		return "Business"
	}
}
