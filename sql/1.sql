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

----------------------------------------------------------------------------------------------------
-- Практика: 1.4
-- Нам нужно найти пассажиров (passenger_id, passenger_name), у которых:
-- * Общее количество билетов в tickets строго больше 4 (COUNT(DISTINCT t.ticket_no) > 4 или COUNT(t.ticket_no) > 4).
-- * Средняя стоимость билета / перелета (стоимость берем либо из bookings.total_amount, либо из поля стоимости в segments, если оно там есть — проверь структуру segments или bookings) превышает 50 000 рублей.
-- Связи таблиц:
-- tickets.book_ref = bookings.book_ref
-- tickets.ticket_no = segments.ticket_no
--
-- Выведи: passenger_id, passenger_name, total_tickets, avg_price (округли через ROUND(..., 2)).
-- Сортировка по avg_price DESC.

select
    t.passenger_id as passenger_id,
    t.passenger_name as passenger_name,
    count(distinct t.ticket_no) as ticket_count,
    round(avg(s.price), 2) as avg_price
from tickets t
left join segments s on s.ticket_no = t.ticket_no
group by t.passenger_id, t.passenger_name
having count(distinct t.ticket_no) > 4 and avg(s.price) > 50000
order by avg_price desc;

-- Вопрос 1.4 (Теория / B-Tree Индексы)
-- [ ] Почему обычный B-Tree индекс `CREATE INDEX idx_tickets_name ON tickets(passenger_name);`
-- совершенно не поможет для запроса:
--
-- SELECT * FROM tickets WHERE passenger_name LIKE '%IVANOV%';
--
-- Ответ: Так как в B-Tree данные отсартированы, если мы будем искать с права `some%` мы быстро сможем
-- найти начало и вывести резултат. В случаях с `%some` или `%some%` мы не сможем сразу найти начало.
-- Для это B-Tree придется пройтись по всему дереву, что почти также или даже хуже чем Seq Scan. 

----------------------------------------------------------------------------------------------------
-- Практика: 1.5 (Поиск по строкам & HAVING)
-- Сделай запрос к tickets для поиска пассажиров по имени:
-- * Найди всех пассажиров, у которых passenger_name содержит слово 'SERGEI' (например, SERGEI PETROV
--   или PETROV SERGEI). Используй LIKE 'SERGEI %' или LIKE '% SERGEI' 
--   (или универсальный LIKE '%SERGEI%').
-- * Выведи passenger_id, passenger_name и количество их билетов (tickets_count).
-- * Оставь только тех пассажиров, у которых билетов строго больше 1 (HAVING > 1).
-- * Отсортируй по tickets_count DESC.

select * from tickets;

select
    passenger_id,
    passenger_name,
    count(distinct ticket_no) as tickers_count
from tickets
where passenger_name like '%SERGEI%'
group by passenger_id, passenger_name
having count(distinct ticket_no) > 1
order by tickers_count desc;

-- Вопрос 1.5 (Теория / Подзапросы vs JOIN)
-- Представь, что нам нужно получить список всех билетов, по которым перелет был из Москвы.
-- Мы можем написать это двумя способами:
--    Вариант А (JOIN):
    SELECT DISTINCT t.ticket_no 
    FROM tickets t
    JOIN segments s ON s.ticket_no = t.ticket_no
    WHERE s.departure_city = 'Moscow';

--    Вариант Б (Подзапрос / IN):
    SELECT t.ticket_no 
    FROM tickets t
    WHERE t.ticket_no IN (
        SELECT s.ticket_no 
        FROM segments s 
        WHERE s.departure_city = 'Moscow'
    );

-- Вопросы:
-- [-] Есть ли разница в производительности между эти двумя запросами в современном Postgres?
-- Ответ: Нет разницы не будет, под капотом Postgres объеденит в Semi-Join

----------------------------------------------------------------------------------------------------
-- Практика: 1.6.1
-- Для каждого билета из таблицы segments пронумеруй его перелеты по порядку.
-- * Выведи: ticket_no, flight_id, price, а также segment_number (через ROW_NUMBER() OVER (PARTITION BY ticket_no ORDER BY flight_id)).

select
    s.ticket_no,
    s.flight_id,
    s.price,
    row_number() over (partition by ticket_no order by flight_id)
from segments s limit 100;


-- Практика: 1.6.2 
-- Найди самый дорогой перелет/сегмент для каждого пассажира из tickets + segments.
-- * Через CTE или подзапрос посчитай DENSE_RANK() OVER (PARTITION BY t.passenger_id ORDER BY s.price DESC) AS price_rank.
-- * Отфильтруй и оставь только записи с price_rank = 1.
-- * Выведи: passenger_id, passenger_name, ticket_no, price.
with ranced_price as (
    select
        t.passenger_id as passenger_id,
        t.passenger_name as passenger_name,
        t.ticket_no as ticket_no,
        s.price as price,
        dense_rank() over (partition by t.passenger_id order by s.price desc) as price_rank
    from tickets t
    left join segments s on t.ticket_no = s.ticket_no
)
select passenger_id, passenger_name, ticket_no, price
from ranced_price
where price_rank = 1;

-- Вопрос: 1.6 
-- [+/-] Выполнится ли запрос:
SELECT price, ROW_NUMBER() OVER (ORDER BY price DESC) as rn 
FROM segments 
WHERE ROW_NUMBER() OVER (ORDER BY price DESC) <= 5;
-- Ответ: Нет, WHERE как и HAVING отрабатывают до оконных функция. Postgres отвыст ошибку
-- windows functions are not allowed in WHERE

----------------------------------------------------------------------------------------------------
-- Практика: 1.7.1
-- Напиши запрос, который найдет перелеты, цена которых выше среднего чека конкретного пассажира:
-- * Создай CTE passenger_stats, которая считат средний чек для каждого пассажира:
--   passenger_id, ROUND(AVG(s.price), 2) AS avg_passenger_price.
-- * В основном запросе соедини tickets, segments и твоё CTE.
-- * Отфильтруй только те сегменты, где s.price > avg_passenger_price.
-- * Выведи: passenger_name, ticket_no, flight_id, price, avg_passenger_price.
-- * Отсортируй по разнице (s.price - avg_passenger_price) DESC.
--

with passenger_stats as (
    select
        t.passenger_id as passenger_id,
        round(avg(s.price), 2) as avg_passenger_price
    from tickets t
    left join segments s on t.ticket_no = s.ticket_no
    group by t.passenger_id
)
select
    t.passenger_name as passenger_name,
    t.ticket_no as ticket_no,
    s.flight_id as flight_id,
    s.price as price,
    ps.avg_passenger_price as avg_passenger_price
from tickets t
left join segments s on t.ticket_no = s.ticket_no
join passenger_stats ps on t.passenger_id = ps.passenger_id
where s.price > ps.avg_passenger_price
order by (s.price - ps.avg_passenger_price);

-- Практика: 1.7.2
-- Посчитай динамику цены между перелетами по одному и тому же билету:
-- * Выведи: ticket_no, flight_id, price.
-- * Добавь prev_price — стоимость предыдущего перелета по этому же билету (LAG(price) OVER (PARTITION BY ticket_no ORDER BY flight_id)).
-- * Добавь price_diff — разницу (price - prev_price).
select
    ticket_no,
    flight_id,
    price,
    lag(price, 1, 0) over (partition by ticket_no order by flight_id) as prev_price,
    (price - lag(price, 1, 0) over (partition by ticket_no order by flight_id)) as price_diff
from segments
order by ticket_no, flight_id;

-- Вопрос: 1.7
-- Посмотри на этот запрос с зависимым подзапросом в SELECT:
SELECT 
    t.ticket_no,
    t.passenger_name,
    (
        SELECT COUNT(*) 
        FROM segments s 
        WHERE s.ticket_no = t.ticket_no
    ) AS segments_count
FROM tickets t;

-- [+] Как Postgres будет выполнять этот подзапрос для каждой строки таблицы tickets?
-- Ответ: Поиде Postgres может сделать Semi-Join, однако в нутри есть агрегатная функция. Как я помню
-- при вычислениях или агрегатных функций нелзя делать Semi-Join. Если подумать то в данном запросе
-- у нас просиходи что мы проходимся по каждому ticket_no и смотрим сколько упоминаний было в segments.
-- Так как мы не можем гарантировать что в s есть/нет повторений у нас агрегатная функция будет слопывать
-- если при Semi-Join это произойдет то внешня функция может слопнутся когда этого нам не нужно.
-- Хотя мы знаем что ticket_no это PK которй уникальный он может и не быть таковым.

-- [+/-] Как переписать его через LEFT JOIN + GROUP BY, чтобы он работал значительно быстрее на больших данных?
select 
    t.ticket_no,
    t.passenger_name,
    count(s.flight_id)
from tickets t
left join segments s on t.ticket_no = s.ticket_no
group by t.ticket_no, t.passenger_name;
