package scheduler

import (
    "fmt"
    "sort"
    "strconv"
    "strings"
    "time"
)

// возвращает в формате "20060102" первую дату > now, 
// начиная с dstart, по правилу repeat
func NextDate(now time.Time, dstart, repeat string) (string, error) {
    date, err := time.Parse("20060102", dstart)
    if err != nil {
        return "", fmt.Errorf("invalid dstart %q: %w", dstart, err)
    }

    // Нормализую время - сбрасываю часы, минуты, секунды для корректного сравнения
    now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
    date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

    repeat = strings.TrimSpace(repeat)
    if repeat == "" {
        return "", fmt.Errorf("no repeat rule")
    }

    parts := strings.SplitN(repeat, " ", 2)
    rule := parts[0]
    params := ""
    if len(parts) == 2 {
        params = parts[1]
    }

    switch rule {
    case "d":
        if params == "" {
            return "", fmt.Errorf("invalid repeat format %q", repeat)
        }
        n, err := strconv.Atoi(params)
        if err != nil || n < 1 || n > 400 {
            return "", fmt.Errorf("invalid interval %q", params)
        }
        for {
            date = date.AddDate(0, 0, n)
            if date.After(now) {
                return date.Format("20060102"), nil
            }
        }

    case "y":
        if params != "" {
            return "", fmt.Errorf("invalid repeat format %q", repeat)
        }
        for {
            date = date.AddDate(1, 0, 0)
            if date.After(now) {
                return date.Format("20060102"), nil
            }
        }

    case "w":
        if params == "" {
            return "", fmt.Errorf("invalid repeat format %q", repeat)
        }

        var allowed [8]bool
        seen := make(map[int]bool)
        for _, tok := range strings.Split(params, ",") {
            v, err := strconv.Atoi(strings.TrimSpace(tok))
            if err != nil || v < 1 || v > 7 {
                return "", fmt.Errorf("invalid weekday %q", tok)
            }
            if !seen[v] {
                allowed[v] = true
                seen[v] = true
            }
        }

        // Проверяю до 1000 дней вперед
        for i := 0; i < 1000; i++ {
            if date.After(now) {
                wd := int(date.Weekday())
                if wd == 0 {
                    wd = 7
                }
                if allowed[wd] {
                    return date.Format("20060102"), nil
                }
            }
            date = date.AddDate(0, 0, 1)
        }
        return "", fmt.Errorf("no matching weekday found within reasonable time")

    case "m":
        if params == "" {
            return "", fmt.Errorf("invalid repeat format %q", repeat)
        }

        parts2 := strings.SplitN(params, " ", 2)
        daysList := strings.Split(parts2[0], ",")
        monthsList := []string{}
        if len(parts2) == 2 {
            monthsList = strings.Split(parts2[1], ",")
        }

        var dayOffsets []int
        seen := make(map[int]bool)

        for _, tok := range daysList {
            tok = strings.TrimSpace(tok)
            d, err := strconv.Atoi(tok)
            if err != nil || d == 0 || d > 31 || d < -2 {
                return "", fmt.Errorf("invalid month-day %q", tok)
            }
            if !seen[d] {
                dayOffsets = append(dayOffsets, d)
                seen[d] = true
            }
        }

        var allowedMonth [13]bool
        if len(monthsList) > 0 {
            for _, tok := range monthsList {
                tok = strings.TrimSpace(tok)
                m, err := strconv.Atoi(tok)
                if err != nil || m < 1 || m > 12 {
                    return "", fmt.Errorf("invalid month %q", tok)
                }
                allowedMonth[m] = true
            }
        }

        // Ограничиваю поиск до 1000 дней
        for i := 0; i < 1000; i++ {
            if date.After(now) {
                year, month := date.Year(), date.Month()
                if len(monthsList) > 0 && !allowedMonth[int(month)] {
                    date = date.AddDate(0, 0, 1)
                    continue
                }

                lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, date.Location()).Day()

                // Пересчитываю дни для текущего месяца
                var validDays []int
                for _, offset := range dayOffsets {
                    var day int
                    if offset > 0 {
                        day = offset
                    } else {
                        day = lastDay + offset + 1
                    }
                    if day >= 1 && day <= lastDay {
                        validDays = append(validDays, day)
                    }
                }

                // Сортирую дни месяца по возрастанию для выбора ближайшего
                sort.Ints(validDays)

                for _, day := range validDays {
                    candidate := time.Date(year, month, day, 0, 0, 0, 0, date.Location())
                    if candidate.After(now) {
                        return candidate.Format("20060102"), nil
                    }
                }
            }
            date = date.AddDate(0, 0, 1)
        }
        return "", fmt.Errorf("no valid date found within reasonable time")

    default:
        return "", fmt.Errorf("invalid repeat format %q", repeat)
    }
}

// Преобразую offset в реальный день месяца
func dayOffsetToDay(date time.Time, offset int) int {
    year, month, _ := date.Date()
    lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, date.Location()).Day()
    if offset > 0 {
        return offset
    }
    return lastDay + offset + 1
}
