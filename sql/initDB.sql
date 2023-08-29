drop table

CREATE TABLE users (
    userid int NOT NULL AUTO_INCREMENT,
    FirstName varchar(255) NOT NULL,
    LastName varchar(255),
    email_id varchar(255),
    `password` varchar(255),
    PRIMARY KEY (userid)
);

insert into users (FirstName,LastName,email_id,`password`) values ("Devansh","Khemani","abc@gmail.com","Dev@04");
insert into users(FirstName,LastName,email_id,`password`) values("Aditya","Goenka","xyz@gmail.com","Adi@89");


CREATE TABLE tasks_info (
    task_id int NOT NULL AUTO_INCREMENT,
    task_name varchar(255),
    task_description varchar(255),
    PRIMARY KEY (task_id)
);

insert into tasks_info (task_name,task_description) values ("Brushing","thoroughly");
insert into tasks_info (task_name,task_description) values ("studying","concentrately");
insert into tasks_info (task_name,task_description) values ("wasing","gently");


CREATE TABLE tasks_status (
    task_id int,
    user_id int,
    status varchar(255)
);


insert into tasks_status (task_id,user_id,status) values (1,1,"done");
insert into tasks_status (task_id,user_id,status) values (2,2,"in_progress");
insert into tasks_status (task_id,user_id,status) values (1,2,"not_done");
insert into tasks_status (task_id,user_id,status) values (3,2,"not_done");

select user_id,tasks_info.task_id,task_name,task_description,status 
from tasks_info left join 
tasks_status on 
tasks_info.task_id = tasks_status.task_id 
where user_id = 1;