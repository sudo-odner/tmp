-- Практика 1.1
-- Нам нужно получить список всех аэропортов, которые находятся в часовом поясе 'Europe/Moscow', 
-- и отсортировать их по названию аэропорта (airport_name) в алфавитном порядке.
-- Используем представление airports. 
-- Выведи: код аэропорта (airport_code), название (airport_name) и город (city).

select airport_code, airport_name, city from airports where timezone = 'Europe/Moscow' order by airport_name asc;

SELECT
    airport_code,
    airport_name,
    city
FROM airports
WHERE timezone = 'Europe/Moscow'
ORDER BY airport_name asc;


-- Вопрос 1.1 
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
