mysql -uroot -p123456 << EOF
use tiktok;
show tables;
source ./tableStruct.sql;
describe users;
