package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	stepLength = 0.65
	mInKm      = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	if strings.TrimSpace(data) == "" {
		return 0, 0, fmt.Errorf("пустой ввод")
	}

	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("неверный формат данных: ожидается 'steps,duration'")
	}

	stepsStr := strings.TrimSpace(parts[0])
	durationStr := strings.TrimSpace(parts[1])

	// Проверяем на наличие пробелов в начале/конце
	if stepsStr != parts[0] || durationStr != parts[1] {
		return 0, 0, fmt.Errorf("неверный формат данных")
	}

	if stepsStr == "" {
		return 0, 0, fmt.Errorf("количество шагов не может быть пустым")
	}

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, 0, fmt.Errorf("неверное количество шагов")
	}

	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть больше 0")
	}

	if durationStr == "" {
		return 0, 0, fmt.Errorf("продолжительность не может быть пустой")
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, fmt.Errorf("неверная продолжительность")
	}

	if duration <= 0 {
		return 0, 0, fmt.Errorf("продолжительность должна быть больше 0")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	if strings.TrimSpace(data) == "" {
		log.Println("пустой ввод")
		return ""
	}

	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	distanceMeters := float64(steps) * stepLength
	distanceKm := distanceMeters / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println(err)
		return ""
	}

	result := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, distanceKm, calories)

	return result
}
