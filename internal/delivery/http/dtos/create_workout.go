package dtos

import (
	"errors"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// todo add exCount

type CreateWorkoutRequest struct {
	WorkoutName string      `json:"workout-name"`
	WorkoutType string      `json:"workout-type"`
	Exercises   []*Exercise `json:"exercises"`
}

type Exercise struct {
	Name                string `json:"name"`
	MuscleGroup         string `json:"muscle-group"`
	RestTime            string `json:"rest-time"` // todo найти решение
	ExecutionType       string `json:"execution-type"`
	ExercisesCombining  string `json:"exercises-combining"`
	NumberOfApproaches  int16  `json:"number-of-approaches"`
	NumberOfRepetitions int16  `json:"number-of-repetitions"`
}

var (
	ErrNameRequired              = errors.New("название обязательно")
	ErrNameLetterCount           = errors.New("название не должно быть меньше 3 или превышать 100 символов")
	ErrExNameBadSymbol           = errors.New("название содержит недопустимые символы (разрешены только буквы, цифры, пробелы, дефисы и градусы)")
	ErrWkNameBadSymbol           = errors.New("название содержит недопустимые символы (разрешены только буквы, цифры, пробелы и дефисы")
	ErrMuscleGroupInvalid        = errors.New("неверная группа мышц")
	ErrRepetitionsRange          = errors.New("количество повторений должно быть от 1 до 50")
	ErrApproachesRange           = errors.New("количество подходов должно быть от 1 до 20")
	ErrWorkoutTypeInvalid        = errors.New("неверный тип тренировки")
	ErrRestTimeInvalid           = errors.New("время отдыха должно быть до 10 минут")
	ErrExecutionTypeInvalid      = errors.New("неверный тип выполнения упражнения")
	ErrExercisesCombiningInvalid = errors.New("неверный тип комбинирования упражнений")
	ErrExercisesCountInvalid     = errors.New("количество упражнений должно быть от 1 до 30")
	ErrExercisesCount            = errors.New("должно быть хотя бы одно упражнение")
	ErrSpaceOrDash               = errors.New("нельзя использовать повторяющиеся пробелы, дефисы или градусы")
	ErrTimeFormat                = errors.New("неправильный формат времени отдыха")
)

var validationHash = map[string]bool{
	"Обычный подход": true,

	// workout type
	"Сплит":            true,
	"Фулл-боди":        true,
	"Тяни-толкай-ноги": true,

	// executionType
	"Дроп-сет": true, // после отказа снижение веса и продолжение подхода.
	"Дроп-сет с отдых-паузой": true, // после отказа короткий отдых, продолжение подхода с понижением веса.
	"Отдых-пауза":             true, // короткие паузы внутри подхода, вес остается прежним.

	// exercises combining
	"Суперсет": true, // два упражнения подряд без отдыха, на одну или разные группы мышц.
	"Трисет":   true, // три упражнения подряд без отдыха, на одну или разные группы мышц.
	"Чередование подходов": true, // два упражнения подряд с отдыхом, на одну или разные группы мышц.

	// muscle group
	"Грудь":      true,
	"Спина":      true,
	"Ноги":       true,
	"Бицепс":     true,
	"Трицепс":    true,
	"Плечи":      true,
	"Пресс":      true,
	"Кардио":     true,
	"Шея":        true,
	"Предплечье": true,
}

func (r *CreateWorkoutRequest) Validate() error {
	vErrs := &ValidationErrors{}

	if err := r.validateWorkoutType(); err != nil {
		vErrs.AddFieldError("workout-type", err)
	}

	if err := r.validateWorkoutName(); err != nil {
		vErrs.AddFieldError("workout-name", err)
	}

	if len(r.Exercises) == 0 {
		vErrs.AddFieldError("exercises", ErrExercisesCount)
	} else if len(r.Exercises) > 30 {
		vErrs.AddFieldError("exercises", ErrExercisesCountInvalid)
	}

	// если до этого момента нет ошибок проверяем упражнения
	if vErrs.FieldErrors == nil {
		for i, v := range r.Exercises {
			err := v.ValidateExercise(i)
			if err != nil {
				var vlErr *ListError
				if errors.As(err, &vlErr) {
					vErrs.AddErrToList(vlErr)
				}
			}
		}
	}

	if vErrs.ListErrors != nil || vErrs.FieldErrors != nil {
		return vErrs
	}
	return nil
}

func (r *CreateWorkoutRequest) validateWorkoutType() error {
	if r.WorkoutType == "" {
		return ErrWorkoutTypeInvalid
	}

	// Normalize case
	runes := []rune(strings.ToLower(r.WorkoutType))
	runes[0] = unicode.ToUpper(runes[0])
	r.WorkoutType = string(runes)
	if !validationHash[r.WorkoutType] {
		return ErrWorkoutTypeInvalid
	}
	return nil
}

func (r *CreateWorkoutRequest) validateWorkoutName() error {
	if len(r.WorkoutName) == 0 {
		return ErrNameRequired
	}

	rc := utf8.RuneCountInString(r.WorkoutName)
	if rc > 100 || rc < 3 {
		return ErrNameLetterCount
	}

	prevChar := rune(-1)

	for _, char := range r.WorkoutName {
		switch {
		case unicode.IsLetter(char):
			prevChar = -1
		case char == ' ' || char == '-':
			if prevChar == char {
				return ErrSpaceOrDash
			}
			prevChar = char
		case unicode.IsDigit(char):
			prevChar = -1
		default:
			return ErrWkNameBadSymbol
		}
	}

	runes := []rune(strings.ToLower(r.WorkoutName))
	runes[0] = unicode.ToUpper(runes[0])

	r.WorkoutName = string(runes)
	return nil
} // todo возможно стоит если есть несколько пробелов, вырезать их

func (e *Exercise) ValidateExercise(index int) error {
	vErr := &ListError{
		Index: index,
	}

	err := e.validateName()
	if err != nil {
		vErr.AddFieldError("name", err)
	}

	err = e.validateMuscleGroup()
	if err != nil {
		vErr.AddFieldError("muscle-group", err)
	}

	err = e.validateRepetitions()
	if err != nil {
		vErr.AddFieldError("number-of-repetitions", err)
	}

	err = e.validateApproaches()
	if err != nil {
		vErr.AddFieldError("number-of-approaches", err)
	}

	err = e.validateRestTime()
	if err != nil {
		vErr.AddFieldError("rest-time", err)
	}

	err = e.validateExecutionType()
	if err != nil {
		vErr.AddFieldError("execution-type", err)
	}

	err = e.validateExercisesCombining()
	if err != nil {
		vErr.AddFieldError("exercises-combining", err)
	}

	if vErr.FieldErrors != nil {
		return vErr
	}

	return nil
}

func (e *Exercise) validateName() error {
	if len(e.Name) == 0 {
		return ErrNameRequired
	}

	rc := utf8.RuneCountInString(e.Name)
	if rc > 100 || rc < 3 {
		return ErrNameLetterCount
	}

	hasLetter := false
	prevChar := rune(-1)

	for _, char := range e.Name {
		switch {
		case unicode.IsLetter(char):
			hasLetter = true
			prevChar = -1
		case char == ' ' || char == '-':
			if prevChar == char {
				return ErrSpaceOrDash
			}
			prevChar = char
		case unicode.IsDigit(char):
			prevChar = -1
		case char == '°':
			if prevChar == '°' {
				return ErrExNameBadSymbol
			}
			prevChar = char
		default:
			return ErrExNameBadSymbol
		}
	}

	if !hasLetter {
		return ErrNameRequired
	}

	runes := []rune(strings.ToLower(e.Name))
	runes[0] = unicode.ToUpper(runes[0])

	e.Name = string(runes)
	return nil
}

func (e *Exercise) validateMuscleGroup() error {
	if len(e.MuscleGroup) == 0 {
		return ErrMuscleGroupInvalid
	}

	runes := []rune(strings.ToLower(e.MuscleGroup))
	runes[0] = unicode.ToUpper(runes[0])
	e.MuscleGroup = string(runes)

	if !validationHash[e.MuscleGroup] {
		return ErrMuscleGroupInvalid
	}

	return nil
}

func (e *Exercise) validateRepetitions() error {
	if e.NumberOfRepetitions < 1 || e.NumberOfRepetitions > 50 {
		return ErrRepetitionsRange
	}

	return nil
}

func (e *Exercise) validateApproaches() error {
	if e.NumberOfApproaches < 1 || e.NumberOfApproaches > 20 {
		return ErrApproachesRange
	}

	return nil
}

func (e *Exercise) validateRestTime() error {
	if e.RestTime == "" {
		e.RestTime = "0m" // todo default value
		return nil
	}

	t, err := time.ParseDuration(e.RestTime)
	if err != nil {
		return ErrTimeFormat
	}
	if t > 10*time.Minute {
		return ErrRestTimeInvalid
	}

	return nil
}

func (e *Exercise) validateExecutionType() error {
	if e.ExecutionType == "" {
		e.ExecutionType = "Обычный подход" // todo default value
		return nil
	}

	runes := []rune(strings.ToLower(e.ExecutionType))
	runes[0] = unicode.ToUpper(runes[0])
	e.ExecutionType = string(runes)

	if !validationHash[e.ExecutionType] {
		return ErrExecutionTypeInvalid
	}
	return nil
}

func (e *Exercise) validateExercisesCombining() error {
	if e.ExercisesCombining == "" {
		e.ExercisesCombining = "Обычный подход" // todo default value
		return nil
	}

	runes := []rune(strings.ToLower(e.ExercisesCombining))
	runes[0] = unicode.ToUpper(runes[0])
	e.ExercisesCombining = string(runes)

	if !validationHash[e.ExercisesCombining] {
		return ErrExercisesCombiningInvalid
	}
	return nil
}
