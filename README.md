# Что нужно сделать

Реализовать интерфейс с методом для проверки правил флуд-контроля. Если за последние N секунд вызовов метода Check будет больше K, значит, проверка на флуд-контроль не пройдена.

- Интерфейс FloodControl располагается в файле main.go.

- Флуд-контроль может быть запущен на нескольких экземплярах приложения одновременно, поэтому нужно предусмотреть общее хранилище данных. Допустимо использовать любое на ваше усмотрение. 

# Необязательно, но было бы круто

Хорошо, если добавите поддержку конфигурации итоговой реализации. Параметры — на ваше усмотрение.

# Ход мыслей
### Хранилище
При частых запросах лучшим вариантом будет кэш, т.к. он обеспечивает быструю чтение/запись. Думал просто сделать map, но в проде так никто не делает и поэтому выбор был очевиден - Redis. 

### Алгоритм
При получении данных мы должны проверить, есть ли ключ в кеше или нет. Если ключа нету, то это означает что ранее вызовов данного пользователя не было и в данном случае мы устанавливаем время последнего вызова в time.Now() и счётчик в 1.

Если данные успешно получены из кэша, проверяем, истек ли временной интервал с момента последнего вызова, с помощью функции intervalHasExpired. Если да, то устанавливаем текущее время в качестве последнего вызова и счетчик вызовов в 1.
Если временной интервал не истек, увеличиваем счетчик вызовов на 1 и проверяем, не превысил ли этот счетчик максимальное допустимое количество вызовов (CallLimitCount). Если превысил, значит flood control пользователь не прошёл.

Если все проверки пройдены успешно, значит пользователь не флудит.

После написания кода пошёл искать бест практисы по флуд контролю и наткнулся на очень качественное и изящное решение - https://github.com/tensorush/flood-control и решил отрефакторить код.

# Запуск

Для запуска билдим образ и поднимаем докер композ

```
docker build -t flood-control:v1 .
docker-compose up
```
