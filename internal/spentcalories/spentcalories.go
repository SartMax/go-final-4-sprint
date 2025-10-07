package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	lenStep = 0.65
	mInKm   = 1000
)

func ParseTraining(input string) (int, string, time.Duration, error) {
	if input == "" {
		return 0, "", 0, fmt.Errorf("пустой ввод")
	}

	parts := strings.Split(input, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("неверный формат данных: ожидается 'steps,activity,duration'")
	}

	stepsStr := strings.TrimSpace(parts[0])
	activity := strings.TrimSpace(parts[1])
	durationStr := strings.TrimSpace(parts[2])

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("неверное количество шагов")
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("количество шагов должно быть положительным")
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("неверный формат продолжительности")
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("продолжительность должна быть больше 0")
	}

	return steps, activity, duration, nil
}

func Distance(steps int) float64 {
	return float64(steps) * lenStep / mInKm
}

func MeanSpeed(distance float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	return distance / duration.Hours()
}

func RunningSpentCalories(steps int, weight float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть больше 0")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть больше 0")
	}

	distance := Distance(steps)
	speed := MeanSpeed(distance, duration)

	calories := (0.035*weight + (speed*speed/1.75)*0.029*weight) * 2 * duration.Hours()

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть больше 0")
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть больше 0")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть больше 0")
	}

	distance := Distance(steps)
	speed := MeanSpeed(distance, duration)

	calories := (0.035*weight + (speed*speed/height)*0.029*weight) * duration.Hours()

	return calories, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := ParseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	distance := Distance(steps)
	speed := MeanSpeed(distance, duration)
	durationHours := duration.Hours()

	var calories float64
	var caloriesErr error

	switch activity {
	case "Бег":
		calories, caloriesErr = RunningSpentCalories(steps, weight, duration)
	case "Ходьба":
		calories, caloriesErr = WalkingSpentCalories(steps, weight, height, duration)
	default:
		errMsg := fmt.Sprintf("неизвестный тип тренировки: %s", activity)
		log.Println(errMsg)
		return "", fmt.Errorf(errMsg)
	}

	if caloriesErr != nil {
		log.Println(caloriesErr)
		return "", caloriesErr
	}

	info := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activity, durationHours, distance, speed, calories)

	return info, nil
}
