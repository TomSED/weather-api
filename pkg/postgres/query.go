package postgres

import (
	"database/sql"
	"time"
)

func (c *Client) InitTables() error {

	script := `CREATE TABLE IF NOT EXISTS public.weather (
		datasource varchar NOT NULL,
		temperature integer NOT NULL,
		windspeed integer NOT NULL,
		updateddate timestamp NOT NULL
	);`

	_, err := c.database.Exec(script)
	if err != nil {
		return err
	}

	return nil
}

type WeatherData struct {
	DataSource  string
	Temperature int
	WindSpeed   int
	UpdatedDate time.Time
}

// InsertWeatherData inserts weather data into database row
func (c *Client) InsertWeatherData(weatherData *WeatherData) error {
	query := `INSERT INTO public.weather VALUES ($1, $2, $3, $4);`

	_, err := c.database.Exec(query, weatherData.DataSource, weatherData.Temperature, weatherData.WindSpeed, weatherData.UpdatedDate)
	if err != nil {
		return err
	}

	return nil
}

// GetLatestWeatherData returns latest weather data sorted by updated date
func (c *Client) GetLatestWeatherData() (*WeatherData, error) {
	query := `SELECT datasource,
				temperature,
				windspeed,
				updateddate
			FROM public.weather
			ORDER BY updateddate desc
			LIMIT 1;`

	row := c.database.QueryRow(query)

	dataSource := ""
	temp := 0
	windSpeed := 0
	updatedDate := time.Time{}
	err := row.Scan(&dataSource, &temp, &windSpeed, &updatedDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	out := &WeatherData{
		DataSource:  dataSource,
		Temperature: temp,
		WindSpeed:   windSpeed,
		UpdatedDate: updatedDate,
	}

	return out, nil
}
