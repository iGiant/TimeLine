package timeline

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCreateTL(t *testing.T) {
	value, err := CreateTL(8, 20, 17, 5)
	if err != nil {
		t.Error(err.Error())
	}
	result := TimeLine{Day: EventTime{Begin: 8*60 + 20, End: 17*60 + 5}}
	if !reflect.DeepEqual(value, result) {
		t.Errorf("Ожидается %+v, получено %+v", result, value)
	}
}

func TestTimeLine_Add(t *testing.T) {
	value, err := CreateTL(8, 20, 17, 5)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = value.Add(12, 0, 12, 45)
	if err != nil {
		t.Error(err.Error())
	}
	result := TimeLine{Day: EventTime{Begin: 8*60 + 20, End: 17*60 + 5}}
	result.EventTimes = append(result.EventTimes, EventTime{Begin: 12*60 + 0, End: 12*60 + 45})
	if !reflect.DeepEqual(value, result) {
		t.Errorf("Ожидается %+v, получено %+v", result, value)
	}
	err = value.Add(12, 30, 13, 20)
	if err == nil {
		t.Error("event intersects with other events")
	}
	if !reflect.DeepEqual(value, result) {
		t.Errorf("Ожидается %+v, получено %+v", result, value)
	}
	err = value.Add(15, 30, 17, 5)
	if err != nil {
		t.Error(err.Error())
	}
	result.EventTimes = append(result.EventTimes, EventTime{Begin: 15*60 + 30, End: 17*60 + 5})
	if !reflect.DeepEqual(value, result) {
		t.Errorf("Ожидается %+v, получено %+v", result, value)
	}
}

func TestTimeLine_AddAnyWay(t *testing.T) {
	value, err := CreateTL(8, 20, 17, 5)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = value.AddAnyWay(12, 0, 12, 45)
	if err != nil {
		t.Error(err.Error())
	}
	result := TimeLine{Day: EventTime{Begin: 8*60 + 20, End: 17*60 + 5}}
	result.EventTimes = append(result.EventTimes, EventTime{Begin: 12*60 + 0, End: 12*60 + 45})
	if !reflect.DeepEqual(value, result) {
		t.Errorf("Ожидается %+v, получено %+v", result, value)
	}
	err = value.AddAnyWay(12, 30, 13, 20)
	if err != nil {
		t.Error(err.Error())
	}
	result.EventTimes = append(result.EventTimes, EventTime{Begin: 12*60 + 30, End: 13*60 + 20})
	if !reflect.DeepEqual(value, result) {
		t.Errorf("Ожидается %+v, получено %+v", result, value)
	}
}

func TestTimeLine_GetEmpty(t *testing.T) {
	value, _ := CreateTL(8, 20, 17, 5)

	_ = value.AddAnyWay(12, 0, 12, 45)
	_ = value.AddAnyWay(12, 30, 13, 20)
	result := []EventTime{{8*60 + 20, 12*60 + 0}, {13*60 + 20, 17*60 + 5}}
	if !reflect.DeepEqual(value.GetEmpty(), result) {
		t.Errorf("Ожидается %+v, получено %+v", result, value.GetEmpty())
	}
}

func TestTimeLine_AddDurationFirst(t *testing.T) {
	//
}

func TestTimeLine_AddDurationMin(t *testing.T) {
	//
}

func TestTimeLine_AddDurationExactTime(t *testing.T) {
	//
}

func TestEventTime_String(t *testing.T) {
	testsTable := []struct {
		value  EventTime
		result string
	}{
		{EventTime{Begin: 0*60 + 0, End: 23*60 + 59}, "00:00 – 23:59"},
		{EventTime{Begin: 0*60 + 1, End: 0*60 + 9}, "00:01 – 00:09"},
		{EventTime{Begin: 8*60 + 30, End: 11*60 + 20}, "08:30 – 11:20"},
		{EventTime{Begin: 22*60 + 5, End: 23*60 + 5}, "22:05 – 23:05"},
	}
	for _, table := range testsTable {
		if fmt.Sprint(table.value) != table.result {
			t.Errorf("Ожидается %q, получили %s", table.result, table.value)
		}
	}
}

func TestOffsetTime_String(t *testing.T) {
	testsTable :=
		[]struct {
			value  OffsetTime
			result string
		}{{0*60 + 0, "00:00"},
			{8*60 + 30, "08:30"},
			{9*60 + 00, "09:00"},
			{11*60 + 01, "11:01"},
			{12*60 + 0, "12:00"},
			{23*60 + 59, "23:59"},
		}
	for _, table := range testsTable {
		if fmt.Sprint(table.value) != table.result {
			t.Errorf("Ожидается %q, получили %s", table.result, table.value)
		}
	}
}
