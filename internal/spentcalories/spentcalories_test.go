package spentcalories

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SpentCaloriesTestSuite struct {
	suite.Suite
}

func TestSpentCaloriesSuite(t *testing.T) {
	suite.Run(t, new(SpentCaloriesTestSuite))
}

func (suite *SpentCaloriesTestSuite) TestParseTraining() {
	tests := []struct {
		name         string
		input        string
		wantSteps    int
		wantDuration time.Duration
		wantErr      bool
	}{
		{
			name:         "корректный ввод с часами и минутами",
			input:        "3456,Ходьба,3h00m",
			wantSteps:    3456,
			wantDuration: 3 * time.Hour,
			wantErr:      false,
		},
		// ... остальные тестовые случаи
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {

			gotSteps, _, gotDuration, err := ParseTraining(tt.input)

			if tt.wantErr {
				assert.Error(suite.T(), err)
				assert.Equal(suite.T(), 0, gotSteps)
				assert.Equal(suite.T(), time.Duration(0), gotDuration)
				return
			}

			assert.NoError(suite.T(), err)
			assert.Equal(suite.T(), tt.wantSteps, gotSteps)
			assert.Equal(suite.T(), tt.wantDuration, gotDuration)
		})
	}
}

func (suite *SpentCaloriesTestSuite) TestDistance() {
	tests := []struct {
		name     string
		steps    int
		wantDist float64
	}{
		{
			name:     "нормальное количество шагов",
			steps:    1000,
			wantDist: 0.65, // 1000 * 0.65 / 1000 = 0.65 км
		},
		{
			name:     "большое количество шагов",
			steps:    10000,
			wantDist: 6.5, // 10000 * 0.65 / 1000 = 6.5 км
		},
		{
			name:     "маленькое количество шагов",
			steps:    100,
			wantDist: 0.065, // 100 * 0.65 / 1000 = 0.065 км
		},
		{
			name:     "ноль шагов",
			steps:    0,
			wantDist: 0,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {

			got := Distance(tt.steps)
			assert.Equal(suite.T(), tt.wantDist, got)
		})
	}
}

func (suite *SpentCaloriesTestSuite) TestMeanSpeed() {
	tests := []struct {
		name      string
		distance  float64
		duration  time.Duration
		wantSpeed float64
	}{
		{
			name:      "нормальная скорость - один час",
			distance:  3.9,
			duration:  1 * time.Hour,
			wantSpeed: 3.9,
		},
		{
			name:      "нормальная скорость - полчаса",
			distance:  1.95,
			duration:  30 * time.Minute,
			wantSpeed: 3.9,
		},
		{
			name:      "нормальная скорость - два часа",
			distance:  7.8,
			duration:  2 * time.Hour,
			wantSpeed: 3.9,
		},
		{
			name:      "маленькая скорость",
			distance:  0.65,
			duration:  2 * time.Hour,
			wantSpeed: 0.325,
		},
		{
			name:      "большая скорость",
			distance:  13.0,
			duration:  1 * time.Hour,
			wantSpeed: 13.0,
		},
		{
			name:      "нулевая продолжительность",
			distance:  1000,
			duration:  0,
			wantSpeed: 0,
		},
		{
			name:      "отрицательная продолжительность",
			distance:  1000,
			duration:  -1 * time.Hour,
			wantSpeed: 0,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {

			got := MeanSpeed(tt.distance, tt.duration)
			assert.Equal(suite.T(), tt.wantSpeed, got)
		})
	}
}

func (suite *SpentCaloriesTestSuite) TestRunningSpentCalories() {
	tests := []struct {
		name     string
		steps    int
		weight   float64
		duration time.Duration
		wantCal  float64
		wantErr  bool
	}{
		{
			name:     "нормальная нагрузка - один час",
			steps:    6000,
			weight:   75.0,
			duration: 1 * time.Hour,
			wantCal:  354.38,
			wantErr:  false,
		},
		{
			name:     "нормальная нагрузка - полчаса",
			steps:    3000,
			weight:   75.0,
			duration: 30 * time.Minute,
			wantCal:  177.19,
			wantErr:  false,
		},
		{
			name:     "высокая скорость",
			steps:    20000,
			weight:   75.0,
			duration: 1 * time.Hour,
			wantCal:  590.62,
			wantErr:  false,
		},
		{
			name:     "низкая скорость",
			steps:    1000,
			weight:   75.0,
			duration: 2 * time.Hour,
			wantCal:  29.53,
			wantErr:  false,
		},
		{
			name:     "другой вес",
			steps:    6000,
			weight:   60.0,
			duration: 1 * time.Hour,
			wantCal:  283.50,
			wantErr:  false,
		},
		{
			name:     "нулевая продолжительность",
			steps:    1000,
			weight:   75.0,
			duration: 0,
			wantCal:  0,
			wantErr:  true,
		},
		{
			name:     "отрицательная продолжительность",
			steps:    1000,
			weight:   75.0,
			duration: -1 * time.Hour,
			wantCal:  0,
			wantErr:  true,
		},
		{
			name:     "ноль шагов",
			steps:    0,
			weight:   75.0,
			duration: 1 * time.Hour,
			wantCal:  0,
			wantErr:  true,
		},
		{
			name:     "отрицательные шаги",
			steps:    -1000,
			weight:   75.0,
			duration: 1 * time.Hour,
			wantCal:  0,
			wantErr:  true,
		},
		{
			name:     "нулевой вес",
			steps:    1000,
			weight:   0,
			duration: 1 * time.Hour,
			wantCal:  0,
			wantErr:  true,
		},
		{
			name:     "отрицательный вес",
			steps:    1000,
			weight:   -75.0,
			duration: 1 * time.Hour,
			wantCal:  0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {

			gotCal, gotErr := RunningSpentCalories(tt.steps, tt.weight, tt.duration)

			if tt.wantErr {
				assert.Error(suite.T(), gotErr)
				assert.Equal(suite.T(), 0.0, gotCal)
				return
			}

			assert.NoError(suite.T(), gotErr)
			assert.InDelta(suite.T(), tt.wantCal, gotCal, 0.1)
		})
	}
}

func (suite *SpentCaloriesTestSuite) TestWalkingSpentCalories() {
	tests := []struct {
		name     string
		steps    int
		weight   float64
		height   float64
		duration time.Duration
		wantCal  float64
		wantErr  bool
	}{
		{
			name:     "нормальная нагрузка",
			steps:    6000,
			weight:   75.0,
			height:   1.75,
			duration: 1 * time.Hour,
			wantCal:  177.19,
			wantErr:  false,
		},
		// ... остальные тестовые случаи
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			gotCal, gotErr := WalkingSpentCalories(tt.steps, tt.weight, tt.height, tt.duration)

			if tt.wantErr {
				assert.Error(suite.T(), gotErr)
				assert.Equal(suite.T(), 0.0, gotCal)
				return
			}

			assert.NoError(suite.T(), gotErr)
			assert.InDelta(suite.T(), tt.wantCal, gotCal, 0.1)
		})
	}
}

func (suite *SpentCaloriesTestSuite) TestTrainingInfo() {
	tests := []struct {
		name    string
		input   string
		weight  float64
		height  float64
		want    string
		wantErr bool
	}{
		{
			name:    "ходьба - нормальная нагрузка",
			input:   "6000,Ходьба,1h00m",
			weight:  75.0,
			height:  1.75,
			want:    "Тип тренировки: Ходьба\nДлительность: 1.00 ч.\nДистанция: 3.90 км.\nСкорость: 3.90 км/ч\nСожгли калорий: 177.19\n",
			wantErr: false,
		},
		// ... остальные тестовые случаи
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			got, err := TrainingInfo(tt.input, tt.weight, tt.height)

			if tt.wantErr {
				assert.Error(suite.T(), err)
				assert.Empty(suite.T(), got)
				if tt.name == "неизвестный тип тренировки - проверка текста ошибки" {
					assert.Contains(suite.T(), err.Error(), "неизвестный тип тренировки")
				}
				return
			}

			assert.NoError(suite.T(), err)
			assert.Equal(suite.T(), tt.want, got)
		})
	}
}
