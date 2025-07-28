package models

// CurrentWeatherResponse respuesta del clima actual
type CurrentWeatherResponse struct {
	Location LocationInfo `json:"location"`
	Current  CurrentInfo  `json:"current"`
}

// ForecastResponse respuesta del pronóstico
type ForecastResponse struct {
	Location LocationInfo `json:"location"`
	Current  CurrentInfo  `json:"current"`
	Forecast ForecastInfo `json:"forecast"`
}

// SearchLocationResponse respuesta de búsqueda de ubicaciones
type SearchLocationResponse []LocationSearchResult

// AstronomyResponse respuesta de datos astronómicos
type AstronomyResponse struct {
	Location  LocationInfo  `json:"location"`
	Astronomy AstronomyInfo `json:"astronomy"`
}

// LocationInfo información de ubicación
type LocationInfo struct {
	Name           string  `json:"name"`
	Region         string  `json:"region"`
	Country        string  `json:"country"`
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	TzID           string  `json:"tz_id"`
	LocaltimeEpoch int64   `json:"localtime_epoch"`
	Localtime      string  `json:"localtime"`
}

// CurrentInfo información del clima actual
type CurrentInfo struct {
	LastUpdatedEpoch int64       `json:"last_updated_epoch"`
	LastUpdated      string      `json:"last_updated"`
	TempC            float64     `json:"temp_c"`
	TempF            float64     `json:"temp_f"`
	IsDay            int         `json:"is_day"`
	Condition        Condition   `json:"condition"`
	WindMph          float64     `json:"wind_mph"`
	WindKph          float64     `json:"wind_kph"`
	WindDegree       int         `json:"wind_degree"`
	WindDir          string      `json:"wind_dir"`
	PressureMb       float64     `json:"pressure_mb"`
	PressureIn       float64     `json:"pressure_in"`
	PrecipMm         float64     `json:"precip_mm"`
	PrecipIn         float64     `json:"precip_in"`
	Humidity         int         `json:"humidity"`
	Cloud            int         `json:"cloud"`
	FeelslikeC       float64     `json:"feelslike_c"`
	FeelslikeF       float64     `json:"feelslike_f"`
	VisKm            float64     `json:"vis_km"`
	VisMiles         float64     `json:"vis_miles"`
	UV               float64     `json:"uv"`
	GustMph          float64     `json:"gust_mph"`
	GustKph          float64     `json:"gust_kph"`
	AirQuality       *AirQuality `json:"air_quality,omitempty"`
}

// Condition información de condición climática
type Condition struct {
	Text string `json:"text"`
	Icon string `json:"icon"`
	Code int    `json:"code"`
}

// AirQuality información de calidad del aire
type AirQuality struct {
	CO           float64 `json:"co"`
	NO2          float64 `json:"no2"`
	O3           float64 `json:"o3"`
	SO2          float64 `json:"so2"`
	PM25         float64 `json:"pm2_5"`
	PM10         float64 `json:"pm10"`
	USEPAIndex   int     `json:"us-epa-index"`
	GBDefraIndex int     `json:"gb-defra-index"`
}

// ForecastInfo información de pronóstico
type ForecastInfo struct {
	Forecastday []ForecastDay `json:"forecastday"`
}

// ForecastDay día de pronóstico
type ForecastDay struct {
	Date      string     `json:"date"`
	DateEpoch int64      `json:"date_epoch"`
	Day       DayInfo    `json:"day"`
	Astro     AstroInfo  `json:"astro"`
	Hour      []HourInfo `json:"hour"`
}

// DayInfo información del día
type DayInfo struct {
	MaxtempC      float64   `json:"maxtemp_c"`
	MaxtempF      float64   `json:"maxtemp_f"`
	MintempC      float64   `json:"mintemp_c"`
	MintempF      float64   `json:"mintemp_f"`
	AvgtempC      float64   `json:"avgtemp_c"`
	AvgtempF      float64   `json:"avgtemp_f"`
	MaxwindMph    float64   `json:"maxwind_mph"`
	MaxwindKph    float64   `json:"maxwind_kph"`
	TotalprecipMm float64   `json:"totalprecip_mm"`
	TotalprecipIn float64   `json:"totalprecip_in"`
	Avgvis_km     float64   `json:"avgvis_km"`
	Avgvis_miles  float64   `json:"avgvis_miles"`
	Avghumidity   float64   `json:"avghumidity"`
	Condition     Condition `json:"condition"`
	UV            float64   `json:"uv"`
}

// AstroInfo información astronómica del día
type AstroInfo struct {
	Sunrise          string `json:"sunrise"`
	Sunset           string `json:"sunset"`
	Moonrise         string `json:"moonrise"`
	Moonset          string `json:"moonset"`
	MoonPhase        string `json:"moon_phase"`
	MoonIllumination string `json:"moon_illumination"`
}

// HourInfo información por hora
type HourInfo struct {
	TimeEpoch    int64       `json:"time_epoch"`
	Time         string      `json:"time"`
	TempC        float64     `json:"temp_c"`
	TempF        float64     `json:"temp_f"`
	IsDay        int         `json:"is_day"`
	Condition    Condition   `json:"condition"`
	WindMph      float64     `json:"wind_mph"`
	WindKph      float64     `json:"wind_kph"`
	WindDegree   int         `json:"wind_degree"`
	WindDir      string      `json:"wind_dir"`
	PressureMb   float64     `json:"pressure_mb"`
	PressureIn   float64     `json:"pressure_in"`
	PrecipMm     float64     `json:"precip_mm"`
	PrecipIn     float64     `json:"precip_in"`
	Humidity     int         `json:"humidity"`
	Cloud        int         `json:"cloud"`
	FeelslikeC   float64     `json:"feelslike_c"`
	FeelslikeF   float64     `json:"feelslike_f"`
	WindchillC   float64     `json:"windchill_c"`
	WindchillF   float64     `json:"windchill_f"`
	HeatindexC   float64     `json:"heatindex_c"`
	HeatindexF   float64     `json:"heatindex_f"`
	DewpointC    float64     `json:"dewpoint_c"`
	DewpointF    float64     `json:"dewpoint_f"`
	WillItRain   int         `json:"will_it_rain"`
	ChanceOfRain interface{} `json:"chance_of_rain"` // Puede ser string o number
	WillItSnow   int         `json:"will_it_snow"`
	ChanceOfSnow interface{} `json:"chance_of_snow"` // Puede ser string o number
	VisKm        float64     `json:"vis_km"`
	VisMiles     float64     `json:"vis_miles"`
	GustMph      float64     `json:"gust_mph"`
	GustKph      float64     `json:"gust_kph"`
	UV           float64     `json:"uv"`
}

// LocationSearchResult resultado de búsqueda de ubicación
type LocationSearchResult struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Region  string  `json:"region"`
	Country string  `json:"country"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	URL     string  `json:"url"`
}

// AstronomyInfo información astronómica
type AstronomyInfo struct {
	Astro AstroInfo `json:"astro"`
}
