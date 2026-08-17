create table meta_datas (
   id int primary key not null,
   last_fetch timestamp,
   total_fetched integer,
   sync_status varchar
);