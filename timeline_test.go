package timeline

import (
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
