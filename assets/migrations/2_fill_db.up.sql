insert into equipment(id, name, description, price)
values (uuid_in(overlay(overlay(md5(random()::text || ':' || clock_timestamp()::text) placing '4' from 13) placing to_hex(floor(random()*(11-8+1) + 8)::int)::text from 17)::cstring), 'equipment1', 'description1', 10000),
       (uuid_in(overlay(overlay(md5(random()::text || ':' || clock_timestamp()::text) placing '4' from 13) placing to_hex(floor(random()*(11-8+1) + 8)::int)::text from 17)::cstring), '110 montauk', '*REDACTED*', 10000);

insert into person(id, login, password, email, full_name, role)
values (uuid_in(overlay(overlay(md5(random()::text || ':' || clock_timestamp()::text) placing '4' from 13) placing to_hex(floor(random()*(11-8+1) + 8)::int)::text from 17)::cstring), 'vpbukhti', 'qwerty', 'vpbukhti@hask.ru', 'Бухитийчук Владимир', 'support'),
       (uuid_in(overlay(overlay(md5(random()::text || ':' || clock_timestamp()::text) placing '4' from 13) placing to_hex(floor(random()*(11-8+1) + 8)::int)::text from 17)::cstring), 'user1', 'qwerty', 'user1@example.com', 'name', 'employee'),
       (uuid_in(overlay(overlay(md5(random()::text || ':' || clock_timestamp()::text) placing '4' from 13) placing to_hex(floor(random()*(11-8+1) + 8)::int)::text from 17)::cstring), 'rvsosn', 'qwerty', 'rvsosnovsky@gmail.com', 'Сосновский Роман', 'employee');

insert into support(id, person_id)
select uuid_in(overlay(overlay(md5(random()::text || ':' || clock_timestamp()::text) placing '4' from 13) placing to_hex(floor(random()*(11-8+1) + 8)::int)::text from 17)::cstring), id from person
where login = 'vpbukhti';
