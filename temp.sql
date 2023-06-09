select
    j.*,
    (select finished_at from run r where r.job_id = j.id order by finished_at desc limit 1) as last_run,
    (select status from run r where r.job_id = j.id order by finished_at desc limit 1) as last_run_status
from job j;

select * from job;


insert into job(id, created_at, updated_at, deleted_at, name,
                requirements, script, type, schedule_unit, schedule_value)
values(100, current_date, null, null, 'Test Job', 'requests', 'print("Hello World")',
       'Scrapper', 'minute', 1);

insert into run(id, created_at, updated_at, deleted_at, container,
                output, message, finished_at, status, job_id)
values (10, current_date, null, null, 'non existent', 'aaaa', '', current_timestamp, 2, 100);


insert into run(id, created_at, updated_at, deleted_at, container,
                output, message, finished_at, status, job_id)
values (11, current_date, null, null, 'non existent', 'aaaa', '', current_timestamp, 3, 100);

insert into run(id, created_at, updated_at, deleted_at, container,
                output, message, finished_at, status, job_id)
values (12, current_date, null, null, 'non existent', 'aaaa', '', null, 1, 100);

insert into job(id, created_at, updated_at, deleted_at, name,
                requirements, script, type, schedule_unit, schedule_value, runner_type_id)
values(101, current_date, null, null, 'Test Job', 'requests', 'print("Hello World")',
       'Scrapper', 'second', 1, 1);

select * from job;
select * from run order by finished_at desc;
delete from run;
