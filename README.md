# TimeLine
Нахождение временных промежутков, достаточных для добавления событий определенной длительности

Пример использования:
```
tl, err := CreateTL(8,20,17,5)
	if err != nil {
		panic(err)
	}
	err = tl.Add(12,0,12,45, true)
	if err != nil {
		fmt.Println(err)
	}
	err = tl.Add(9,30,12,20, true)
	if err != nil {
		fmt.Println(err)
	}
	err = tl.Add(8,20, 9,40, false)
	if err != nil {
		fmt.Println(err)
	}
	err = tl.Add(15,0, 16,0, true)
	fmt.Println("Рабочий день:", tl.WorkingDay)
	for _, event := range tl.Events {
		fmt.Println(event)
	}
	fmt.Println("Свободное время:")
	for _, event := range tl.GetEmpty() {
		fmt.Println(event)
	}
	event, err := tl.AddDuration(125)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Добавлено событие:\n%s\n", event)
	fmt.Println("Свободное время:")
	for _, event := range tl.GetEmpty() {
		fmt.Println(event)
	}
  ```
