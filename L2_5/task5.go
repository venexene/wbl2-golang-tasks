package main

type customError struct {
  msg string
}

func (e *customError) Error() string {
  return e.msg
}

func test() *customError {
  // ... do something
  return nil
}

func main() {
  var err error
  err = test()
  if err != nil {
    println("error")
    return
  }
  println("ok")
}

/*
Выведется: error. Функция test() возвращает ошибку со значением nil и типом *customError.
При проверке err != nil получим истину, так как в интерфейсе и значение и тип должны быть nil,
чтобы он считался nil.
*/