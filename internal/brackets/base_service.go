package brackets

import (
	"context"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"strings"
)

// Конфигурация сервиса
type Config struct {
	Logger *zap.SugaredLogger
	// Список открывающих и закрывающих скобок в формате "{}()[]"
	Bkts string
}

func NewBaseService(cfg *Config) (error, Brackets) {
	if len(cfg.Bkts) == 0 {
		return errors.New("Необходимо указать список обрабатываемых скобок в формате \"{}()[]\""), nil
	}

	bkts := strings.Split(cfg.Bkts, "")
	if len(bkts)%2 != 0 {
		return errors.New("Количество скобок должно быть четным"), nil
	}
	// Создаем сервис
	svc := &baseService{
		logger: cfg.Logger,
	}
	// Заполняем оппозитные скобки
	svc.fillBkts(bkts)

	return nil, svc
}

type baseService struct {
	logger    *zap.SugaredLogger
	openBkts  []string
	closeBkts []string
}

// IsValid Проверяет, что строка имеет правильную структуру открывающих и закрывающих скобок
func (svc *baseService) Validate(ctx context.Context, str string) bool {
	if len(str) == 0 {
		return true
	}
	// Разделяем строку на символы
	chars := strings.Split(str, "")
	// Количество символов должно быть четным
	if len(chars)%2 != 0 {
		return false
	}

	// Будем использовать FIFO очередь для проверяемых символов
	var queue []string
	for _, char := range chars {
		// Если символ открывающий - добавляем в начало очереди
		if svc.isOpenBracket(char) {
			queue = append([]string{char}, queue...)
			continue
		}
		// Если символ закрывающий - берем из очереди первый символ и убеждаемся что они подходят друг другу
		if svc.isCloseBracket(char) {
			// Если очередь пуста
			if len(queue) == 0 {
				return false
			}
			// Получаем индекс закрывающей скобки
			bktIdx := svc.getCloseBktIdx(char)
			// Если последний символ в очереди не равен нужному открывающему
			if queue[0] != svc.openBkts[bktIdx] {
				return false
			}
			// Удаляем первый элемент очереди
			queue = queue[1:]
			continue
		}
		// В противном случае символ не скобка
		return false
	}
	// Если в очереди остались символы
	return len(queue) == 0
}

func (svc *baseService) Fix(ctx context.Context, str string) (error, string) {
	if len(str) == 0 {
		return errors.New("Была передана пустая строка"), str
	}
	// Разделяем строку на символы
	chars := strings.Split(str, "")

	// Будем использовать FIFO очередь для проверяемых символов
	var queue []string
	var result []string
	for _, char := range chars {
		// Если символ открывающий - добавляем в начало очереди и в результат
		if svc.isOpenBracket(char) {
			queue = append([]string{char}, queue...)
			result = append(result, char)
			continue
		}
		// Если символ закрывающий
		if svc.isCloseBracket(char) {
			// Получаем индекс закрывающей скобки
			bktIdx := svc.getCloseBktIdx(char)
			// Нужная нам открывающая скобка
			openBkt := svc.openBkts[bktIdx]

			// Если очередь пуста
			if len(queue) == 0 {
				result = append(result, openBkt, char)
				continue
			}
			// Может получиться так, что открывающий символ находится в глубине очереди
			// В таком случае необходимо закрыть все открытые до него скобки
			openBktInStack := false
			for idx := range queue {
				if queue[idx] == openBkt {
					openBktInStack = true
					// Проходим по очереди до нужного нам символа
					for i := 0; i <= idx; i++ {
						openStackBkt := queue[i]
						for openIdx, openBkt := range svc.openBkts {
							if openBkt == openStackBkt {
								result = append(result, svc.closeBkts[openIdx])
								break
							}
						}
					}
					// Очищаем очередь
					queue = queue[idx+1:]
					break
				}
			}
			// Если в очереди нет открывающей скобки - добавляем ее вместе с проверяемым символом
			if openBktInStack == false {
				result = append(result, openBkt, char)
			}
			continue
		}
		// В противном случае символ не скобка
		return errors.Errorf("Символ %s не описан в списке скобок", char), str
	}
	// Если в очереди остались скобки - дополняем результат закрывающими скобками
	for _, char := range queue {
		for idx, openBkt := range svc.openBkts {
			if openBkt == char {
				result = append(result, svc.closeBkts[idx])
				break
			}
		}
	}

	return nil, strings.Join(result, "")
}

// Проверяет, что строка является одной из открываемых скобок
func (svc *baseService) isOpenBracket(str string) bool {
	for _, bkt := range svc.openBkts {
		if str == bkt {
			return true
		}
	}
	return false
}

// Проверяет, что строка является одной из закрываемых скобок
func (svc *baseService) isCloseBracket(str string) bool {
	for _, bkt := range svc.closeBkts {
		if str == bkt {
			return true
		}
	}
	return false
}

// Получает индекс закрывающей скобки
func (svc *baseService) getCloseBktIdx(str string) int {
	for idx, bkt := range svc.closeBkts {
		if str == bkt {
			return idx
		}
	}
	return 0
}

// Заполняет список оппозитных скобок выравнивая по ключу
func (svc *baseService) fillBkts(bkts []string) {
	for i := 0; i < len(bkts); i += 2 {
		svc.openBkts = append(svc.openBkts, bkts[i])
		svc.closeBkts = append(svc.closeBkts, bkts[i+1])
	}
}
