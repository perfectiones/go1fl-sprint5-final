package trainings

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

// Training содержит данные о тренировке и персональные данные пользователя.
type Training struct {
	Steps        int           // количество шагов.
	TrainingType string        // тип тренировки: "Бег" или "Ходьба".
	Duration     time.Duration // продолжительность тренировки.
	personaldata.Personal
}

// Parse разбирает строку вида "3456,Ходьба,3h00m" и заполняет поля структуры.
// Поля меняются только если все части строки корректны.
func (t *Training) Parse(datastring string) (err error) {
	parts := strings.Split(datastring, ",")
	if len(parts) != 3 {
		return fmt.Errorf("неверный формат данных %q: ожидалось 3 значения, получено %d", datastring, len(parts))
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("не удалось преобразовать количество шагов %q: %w", parts[0], err)
	}
	if steps <= 0 {
		return fmt.Errorf("количество шагов должно быть больше нуля, получено %d", steps)
	}

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return fmt.Errorf("не удалось преобразовать продолжительность %q: %w", parts[2], err)
	}
	if duration <= 0 {
		return fmt.Errorf("продолжительность должна быть больше нуля, получено %v", duration)
	}

	t.Steps = steps
	t.TrainingType = parts[1]
	t.Duration = duration
	return nil
}

// ActionInfo возвращает строку с информацией о тренировке:
// тип, длительность, дистанцию, среднюю скорость и потраченные калории.
func (t Training) ActionInfo() (string, error) {
	var (
		calories float64
		err      error
	)

	switch t.TrainingType {
	case "Бег":
		calories, err = spentenergy.RunningSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	case "Ходьба":
		calories, err = spentenergy.WalkingSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
	if err != nil {
		return "", fmt.Errorf("не удалось рассчитать калории: %w", err)
	}

	distance := spentenergy.Distance(t.Steps, t.Height)
	meanSpeed := spentenergy.MeanSpeed(t.Steps, t.Height, t.Duration)

	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		t.TrainingType, t.Duration.Hours(), distance, meanSpeed, calories,
	), nil
}
