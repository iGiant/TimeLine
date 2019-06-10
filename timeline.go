package timeline

import (
	"fmt"
	"sort"
)

type (
	OffsetTime int
	// EventTime время события
	EventTime struct {
		// Begin - начало события, End - конец события
		Begin, End OffsetTime
	}
	// TimeLine - временная линия
	TimeLine struct {
		// EventTimes список времен событий
		EventTimes []EventTime
		// Day - начало и конец временной линии (например, рабочий день)
		Day EventTime
	}

	windowsType int
)

const (
	first windowsType = iota
	min
	exactTime
)

// String - переопределение для вывода времени события в формате hh:mm – hh:mm
func (e EventTime) String() string {
	return fmt.Sprintf("%s – %s", e.Begin, e.End)
}

func (o OffsetTime) String() string {
	return fmt.Sprintf("%.2d:%.2d", o/60, o%60)
}

// CreateTL - конструктор создания временной линии
// beginHour, beginMinute, endHour, endMinute - время начала и конца временной линии в часах и минутах
func CreateTL(beginHour, beginMinute, endHour, endMinute int) (TimeLine, error) {
	if beginHour < 0 || beginHour > 23 {
		beginHour = 0
	}
	if beginMinute < 0 || beginMinute > 59 {
		beginMinute = 0
	}
	if endHour < 0 || endHour > 23 {
		endHour = 23
	}
	if endMinute < 0 || endMinute > 59 {
		endMinute = 59
	}
	beginTL := time2tL(OffsetTime(beginHour), OffsetTime(beginMinute))
	endTL := time2tL(OffsetTime(endHour), OffsetTime(endMinute))
	if beginTL > endTL {
		return TimeLine{}, fmt.Errorf("the beginning of the period is later than the ending")
	}
	tl := TimeLine{Day: EventTime{Begin: beginTL, End: endTL}}
	_ = tl.addEvent(beginTL, beginTL, true)
	_ = tl.addEvent(endTL, endTL, true)
	return tl, nil
}

// Add - добавление периода события во временную линию
// hoursBegin, minutesBegin, hoursEnd, minutesEnd - время начала и конца события в часах и минутах
// doNotMatter - добавлять ли событие, если его время пересекается со временем добавленного ранее события
// возвращает ошибку, если время пересекается
func (tl *TimeLine) Add(hoursBegin, minutesBegin, hoursEnd, minutesEnd int) error {
	return tl.addEvent(time2tL(OffsetTime(hoursBegin), OffsetTime(minutesBegin)),
		time2tL(OffsetTime(hoursEnd), OffsetTime(minutesEnd)), false)
}

func (tl *TimeLine) AddAnyWay(hoursBegin, minutesBegin, hoursEnd, minutesEnd int) error {
	return tl.addEvent(time2tL(OffsetTime(hoursBegin), OffsetTime(minutesBegin)),
		time2tL(OffsetTime(hoursEnd), OffsetTime(minutesEnd)), true)
}

func (tl *TimeLine) addEvent(begin, end OffsetTime, doNotMatter bool) (err error) {
	for _, event := range tl.EventTimes {
		if (begin > event.Begin && begin < event.End) || (end > event.Begin && end < event.End) {
			err = fmt.Errorf("event intersects with other events")
			if doNotMatter {
				break
			} else {
				return
			}
		}
	}
	(*tl).EventTimes = append((*tl).EventTimes, EventTime{Begin: OffsetTime(begin), End: OffsetTime(end)})
	tl.sort()
	return
}

// AddDurationFirst - добавляет событие в первое свободное "окно"
func (tl *TimeLine) AddDurationFirst(duration int) (EventTime, error) {
	return tl.addDuration(0, 0, duration, first)
}

// AddDurationMin - добавляет событие в минимальное по размеру свободное "окно"
func (tl *TimeLine) AddDurationMin(duration int) (EventTime, error) {
	return tl.addDuration(0, 0, duration, min)
}

// AddDurationExactTime - добавляет событие в точное время и с известной длительностью
func (tl *TimeLine) AddDurationExactTime(beginHour, beginMinute, duration int) (EventTime, error) {
	return tl.addDuration(beginHour, beginMinute, duration, exactTime)
}

// AddDuration - добавляет время события в свободное "окно"
// duration - длительность события в минутах
// first - признак добавления (true - добавить в первое подходящее "окно", false - добавить в минимальное подходящее "окно"
func (tl *TimeLine) addDuration(beginH, beginM, duration int, windows windowsType) (EventTime, error) {
	var (
		err        error
		begin, end OffsetTime
	)
	if windows == exactTime {
		beginOffset := time2tL(OffsetTime(beginH), OffsetTime(beginM))
		if beginOffset >= tl.Day.Begin {
			begin, end = beginOffset, beginOffset+OffsetTime(duration)
			err = tl.addEvent(begin, end, false)
			return EventTime{begin, end}, err
		} else {
			return EventTime{}, fmt.Errorf("event starts before timeline begins")
		}
	}
	events := tl.GetEmpty()
	if len(events) == 0 {
		return EventTime{}, fmt.Errorf("no free period")
	}
	var (
		index        = -1
		min          OffsetTime
		tempDuration OffsetTime
		flag         = true
	)
	for i := range events {
		tempDuration = events[i].End - events[i].Begin
		if (tempDuration >= OffsetTime(duration)) && (flag || tempDuration < min) {
			if windows == first {
				begin, end = events[i].Begin, events[i].Begin+OffsetTime(duration)
				err = tl.addEvent(begin, end, true)
				return EventTime{begin, end}, err
			}
			index = i
			min = tempDuration
			flag = false
		}
	}
	if index == -1 {
		return EventTime{}, fmt.Errorf("no free period")
	}
	begin, end = events[index].Begin, events[index].Begin+OffsetTime(duration)
	err = tl.addEvent(begin, end, false)
	return EventTime{begin, end}, err
}

// GetEmpty() получить список свободных "окон" во временной линии
func (tl TimeLine) GetEmpty() []EventTime {
	events := make([]EventTime, 0, 20)
	var begin, end OffsetTime
	for i := 0; i < len(tl.EventTimes)-1; i++ {
		if tl.EventTimes[i].End < tl.EventTimes[i+1].Begin {
			begin = tl.EventTimes[i].End
			end = tl.EventTimes[i+1].Begin
			events = append(events, EventTime{Begin: begin, End: end})
		}
	}
	return events
}

// сортировка событий по времени начала
func (tl *TimeLine) sort() {
	sort.Slice((*tl).EventTimes, func(i, j int) bool { return (*tl).EventTimes[i].Begin < (*tl).EventTimes[j].Begin })
}

// конвертация времени из часов и минут в смещение в минутах от начала суток
func time2tL(hour, minute OffsetTime) OffsetTime {
	return OffsetTime(hour*60 + minute)
}

// конвертация времени из смещения в минутах от начала суток в часы и минуты
func tL2time(tl OffsetTime) (OffsetTime, OffsetTime) {
	absTL := tl
	return absTL / 60, absTL % 60
}
