Alter user 'forest'@'%' IDENTIFIED WITH "123456" BY 'forest';
GRANT ALL PRIVILEGES ON report.* TO 'forest'@'%';
FLUSH PRIVILEGES;