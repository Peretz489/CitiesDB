create table cityData (cityID int primary key unique not null,
					   cityName varchar(30),
					   region varchar(30),
					   district varchar(30),
					   population int,
					   foundation int)