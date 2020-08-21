create database `lovenotesdb`;

create table `lovenotesdb`.`notes` (
	`id` int not null auto_increment,
    `note` varchar(255) not null,
    primary key(`id`));

insert into `lovenotesdb`.`notes` (`note`) values("I love you");
insert into `lovenotesdb`.`notes` (`note`) values("I love you more");
insert into `lovenotesdb`.`notes` (`note`) values("Hey baby, I'm fine!");

select * from `lovenotesdb`.`notes`;