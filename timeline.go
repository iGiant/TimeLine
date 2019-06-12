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
	minimum
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
	beginTL := OffsetTime(beginHour*60 + beginMinute)
	endTL := OffsetTime(endHour*60 + endMinute)
	if beginTL > endTL {
		return TimeLine{}, fmt.Errorf("the beginning of the period is later than the ending")
	}
	tl := TimeLine{Day: EventTime{Begin: beginTL, End: endTL}}
	return tl, nil
}

// Add - добавление периода события во временную линию
// hoursBegin, minutesBegin, hoursEnd, minutesEnd - время начала и конца события в часах и минутах
// Если событие пересекается с ранее введенным, то оно не добавляется.
func (tl *TimeLine) Add(hoursBegin, minutesBegin, hoursEnd, minutesEnd int) error {
	return tl.addEvent(OffsetTime(hoursBegin*60+minutesBegin), OffsetTime(hoursEnd*60+minutesEnd), false)
}

// AddAnyWay - добавление периода события во временную линию
// hoursBegin, minutesBegin, hoursEnd, minutesEnd - время начала и конца события в часах и минутах
// Если событие пересекается с ранее введенным, то оно добавляется.
func (tl *TimeLine) AddAnyWay(hoursBegin, minutesBegin, hoursEnd, minutesEnd int) error {
	return tl.addEvent(OffsetTime(hoursBegin*60+minutesBegin), OffsetTime(hoursEnd*60+minutesEnd), true)
}

func (tl *TimeLine) addEvent(begin, end OffsetTime, doNotMatter bool) (err error) {
	if !doNotMatter {
		for _, event := range tl.EventTimes {
			if (begin > event.Begin && begin < event.End) || (end > event.Begin && end < event.End) {
				err = fmt.Errorf("event intersects with other events")
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
	return tl.addDuration(0, 0, duration, minimum)
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
		minBegin     OffsetTime
		minDuration  = tl.Day.End - tl.Day.Begin
		tempDuration OffsetTime
		lastEvent    OffsetTime
	)
	switch windows {
	case exactTime:
		{
			beginOffset := OffsetTime(beginH*60 + beginM)
			if beginOffset >= tl.Day.Begin {
				begin, end := beginOffset, beginOffset+OffsetTime(duration)
				_ = tl.addEvent(begin, end, true)
				return EventTime{begin, end}, nil
			}
			return EventTime{}, fmt.Errorf("event starts before timeline begins")
		}
	case first:
		{
			events := tl.GetEmpty()
			if len(events) == 0 {
				break
			}
			for _, event := range events {
				if event.End > lastEvent {
					lastEvent = event.End
				}
				tempDuration = lastEvent - event.Begin
				if tempDuration >= OffsetTime(duration) {
					_ = tl.addEvent(event.Begin, event.Begin+OffsetTime(duration), true)
					return EventTime{event.Begin, event.Begin + OffsetTime(duration)}, nil
				}
			}
		}
	case minimum:
		{
			events := tl.GetEmpty()
			if len(events) == 0 {
				break
			}
			for _, event := range events {
				if event.End > lastEvent {
					lastEvent = event.End
				}
				tempDuration = lastEvent - event.Begin
				if (tempDuration >= OffsetTime(duration)) && (tempDuration < minDuration) {
					minBegin = event.Begin
					minDuration = tempDuration
				}
			}
			if minBegin != 0 {
				_ = tl.addEvent(minBegin, minBegin+OffsetTime(duration), true)
				return EventTime{minBegin, minBegin + OffsetTime(duration)}, nil
			}
		}
	}
	return EventTime{}, fmt.Errorf("no free period")
}

// GetEmpty() получить список свободных "окон" во временной линии
func (tl TimeLine) GetEmpty() []EventTime {
	if len(tl.EventTimes) == 0 {
		return []EventTime{{tl.Day.Begin, tl.Day.End}}
	}
	events := make([]EventTime, 0)
	if tl.EventTimes[0].Begin > tl.Day.Begin {
		events = append(events, EventTime{tl.Day.Begin, tl.EventTimes[0].Begin})
	}
	for i := 0; i < len(tl.EventTimes)-1; i++ {
		if tl.EventTimes[i].End < tl.EventTimes[i+1].Begin {
			events = append(events, EventTime{Begin: tl.EventTimes[i].End, End: tl.EventTimes[i+1].Begin})
		}
	}
	if tl.EventTimes[len(tl.EventTimes)-1].End < tl.Day.End {
		events = append(events, EventTime{Begin: tl.EventTimes[len(tl.EventTimes)-1].End, End: tl.Day.End})
	}
	return events
}

// сортировка событий по времени начала
func (tl *TimeLine) sort() {
	sort.Slice((*tl).EventTimes, func(i, j int) bool { return (*tl).EventTimes[i].Begin < (*tl).EventTimes[j].Begin })
}
