package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

// DaySteps содержит данные о прогулке и персональные данные пользователя.
type DaySteps struct {
	Steps    int           // количество шагов.
	Duration time.Duration // продолжительность прогулки.
	personaldata.Personal
}

// Parse разбирает строку вида "678,0h50m" и заполняет поля структуры.
// Поля меняются только если все части строки корректны.
func (ds *DaySteps) Parse(datastring string) (err error) {
	parts := strings.Split(datastring, ",")
	if len(parts) != 2 {
		return fmt.Errorf("неверный формат данных %q: ожидалось 2 значения, получено %d", datastring, len(parts))
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("не удалось преобразовать количество шагов %q: %w", parts[0], err)
	}
	if steps <= 0 {
		return fmt.Errorf("количество шагов должно быть больше нуля, получено %d", steps)
	}

	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return fmt.Errorf("не удалось преобразовать продолжительность %q: %w", parts[1], err)
	}
	if duration <= 0 {
		return fmt.Errorf("продолжительность должна быть больше нуля, получено %v", duration)
	}

	ds.Steps = steps
	ds.Duration = duration
	return nil
}

// ActionInfo возвращает строку с информацией о прогулке:
// количество шагов, дистанцию и потраченные калории.
func (ds DaySteps) ActionInfo() (string, error) {
	calories, err := spentenergy.WalkingSpentCalories(ds.Steps, ds.Weight, ds.Height, ds.Duration)
	if err != nil {
		return "", fmt.Errorf("не удалось рассчитать калории: %w", err)
	}

	distance := spentenergy.Distance(ds.Steps, ds.Height)

	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		ds.Steps, distance, calories,
	), nil
}
