-- Практика: 1.1
-- Нам нужно получить список всех аэропортов, которые находятся в часовом поясе 'Europe/Moscow', 
-- и отсортировать их по названию аэропорта (airport_name) в алфавитном порядке.
-- Используем представление airports. 
--
-- Выведи: код аэропорта (airport_code), название (airport_name) и город (city).

select airport_code, airport_name, city from ariports where timezone = 'Eruope/Moscow' order by airport_name asc;

SELECT
    airport_code,
    airport_name,
    city
FROM airports
WHERE timezone = 'Europe/Moscow'
ORDER BY airport_name asc;


-- Вопрос: 1.1 
-- Представь, что мы часто выполняем поиск аэропортов по конкретному часовому поясу (timezone).

-- [+] Какой тип индекса Postgres создаст по умолчанию, если мы напишем:
-- CREATE INDEX idx_airports_timezone ON airports_data(timezone);?
-- Ответ: Postgres по умолчанию всегда создаст BTree
--
-- [+/-] Поможет ли этот индекс ускорить выборку, если мы ищем WHERE timezone = 'Europe/Moscow'?
-- Ответ: Да индекс может помочь, однако не факт что Postgres будет его использовать.
-- При малой селективности (когда данные не так часто повторяются) например scheduler Postgres для 
-- 30-40% выгоднее сделаить Seq Scan, чем делать Index Scan. Потому что при использовании индексов 
-- происходит больше Cache miss, так как индекс будет праыгать по рандомным страницам.

----------------------------------------------------------------------------------------------------
-- Практика: 1.2
-- Нам нужно найти ТОП-5 городов по количеству аэропортов.
-- Сортировка по убыванию количества аэропортов. 
-- Если количество одинаковое — отсортируй по названию города по алфавиту.
--     Подсказка: Используем таблицу/представление airports. Подумай про GROUP BY, HAVING / ORDER BY и LIMIT.
--
-- Выведи: Название города (city), Количество аэропортов в этом городе (airports_count)

select city, count(city) as airports_count from airports group by city order by airports_count desc, city asc limit 5;

SELECT city,
    count(city) as ariports_count
FROM airports
GROUP BY city
ORDER BY ariports_count desc, city asc
LIMIT 5;


-- Вопрос: 1.2
-- Представь ситуацию на проде в Go-сервисе бронирования:
-- Транзакция А (Покупатель): Зашла в систему, выполняет UPDATE билета, меняя статус места:
-- | BEGIN;
-- | UPDATE ticket_flights SET fare_conditions = 'Business' WHERE ticket_no = '12345' AND flight_id = 100;
--
-- Транзакция Б (Менеджер): Через секунду пытается выполнить простой SELECT:
-- || SELECT * FROM ticket_flights WHERE ticket_no = '12345' AND flight_id = 100;

-- [+] Заблокируется ли Транзакция Б (SELECT) в ожидании Транзакции А?
-- Не заблокируется, потому что тут нет блокировок на строки и используется стандартная изоляция
-- Read Commited.
--
-- [+] Какое значение fare_conditions (старое или 'Business') увидит Транзакция Б?k
-- Старое

----------------------------------------------------------------------------------------------------
-- Практика: 1.3
-- Нам нужно найти 3 самых дорогих бронирования из таблицы bookings, но с условием:
-- Нам интересны только те бронирования, где количество пассажиров (билетов) в одном бронировании больше 2.
--
-- Сортировка по total_amount от самых дорогих к дешевым (DESC). Лимит — 3 записи.
--     Подсказка по схеме:
--     Связь между таблицами: bookings.book_ref = tickets.book_ref.
--     Не забудь сгруппировать по book_ref (и total_amount) и отфильтровать количество билетов через HAVING.
--
-- Что нужно вывести:
--    book_ref (номер бронирования из bookings)
--    total_amount (полная стоимость из bookings)
--    tickets_count (количество билетов, привязанных к этой брони из tickets)

select 
    b.book_ref,
    b.total_amount,
    count(distinct t.passenger_id) as count_passenger
from bookings b
left join tickets t on t.book_ref = b.book_ref
group by b.book_ref, b.total_amount
having count(distinct t.passenger_id) > 2
order by b.total_amount desc
limit 3;


-- Вопросы: 1.3 
-- Представь реальную ситуацию на проде в Go-микросервисе:
-- У тебя есть ручка (HTTP handler) отмены бронирования и возврата денег.
-- В Go-коде открывается транзакция tx.Begin(). Код выполняет запрос:
--      SELECT book_ref, total_amount FROM bookings WHERE book_ref = '00000F' FOR UPDATE;
-- После этого Go-сервис делает внешний HTTP-запрос в банк (REST API) для возврата денег, 
-- который задерживается на 5 секунд.
-- В эти 5 секунд другой пользователь пытаешься прочитать статус этой же брони через обычный запрос:
--      SELECT book_ref, total_amount FROM bookings WHERE book_ref = '00000F';

-- [-] Заблокируется ли второй пользователь на 5 секунд из-за внешнего HTTP-запроса в банк?
-- Не заблокируется. Хоть в первой транзакции мы блокирем строку `FOR UPDATE`, она не блокирет обычный
-- `SELECT`. Он просто прочитает старые не актуальные данные.
--
-- [+/-] Почему делать HTTP-запросы внутри открытой транзакции базы данных — это катастрофический 
--     антипаттерн в Go?
-- В первую очередь это антипаттерн clean arcitecture слой BD одновременно отвечает за хранение и 
-- ходение в сеть + во время выоплненние внешних запросов у нас держится соединение и другие 
-- транзакции могут не пройти про ctx
