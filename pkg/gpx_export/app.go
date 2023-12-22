package gpx_export

type App = string

const (
	Strava   App = "strava"
	GarminCN App = "garmin_cn"
)

type AppExport interface {
	Run() error
	Auth(user, pwd string)
}
