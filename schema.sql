create table sublicenses (
	id integer primary key,
	name text,
	numberOfSeats integer,
	licenseKey text,
	expiryDate integer,
	activ integer
);
create table login_cred (
	username varchar(25) primary key,
	pwd_hash blob(60)
)
