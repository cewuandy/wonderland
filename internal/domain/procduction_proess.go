package domain

type ProductionProcessReq struct {
	Name             string `form:"name"`
	Quantity         int    `form:"quantity"`
	ClockEnable      bool   `form:"clockEnable"`
	StandClockEnable bool   `form:"standClockEnable"`
	WindmillEnable   bool   `form:"windmillEnable"`
	ACEnable         bool   `form:"acEnable"`
}
