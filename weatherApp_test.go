package weatherApp

import (
	"testing"

	"github.com/golang/mock/gomock"
	mock_weatherAppGo "github.com/rijuCB/weatherAppGo/mocks"
	"github.com/stretchr/testify/require"
)

func TestWeatherAppAPI(t *testing.T) {
	var (
		ctrl = gomock.NewController(t)
		cli  = mock_weatherAppGo.NewMockIweather(ctrl)
	)

	cli.EXPECT().GetTemp().Return(4.2)
	cli.EXPECT().GetWind().Return(7.8)
	cli.EXPECT().GetRain().Return(5.6)
	a := CoallesceWeatherInfo(cli)
	require.Equal(t, WeatherData{4.2, 7.8, 5.6}, a)
}
