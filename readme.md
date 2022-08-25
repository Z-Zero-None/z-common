#### 一、Mysql
1.运行mysql

2.内部自动执行了apollo所需要的表脚本

3.修改对应的表数据
```
UPDATE `ApolloConfigDB`.`ServerConfig` SET `Key` = 'eureka.service.url', `Cluster` = 'default', `Value` = 'http://apollo-configservice:8080/eureka/', `Comment` = 'Eureka服务Url，多个service以英文逗号分隔', `IsDeleted` = b'0', `DataChange_CreatedBy` = 'default', `DataChange_CreatedTime` = '2022-08-22 16:32:32', `DataChange_LastModifiedBy` = '', `DataChange_LastTime` = '2022-08-22 17:22:59' WHERE `Id` = 1;
```
4.执行建表语句
```
CREATE TABLE `country` (
  `code` char(2) NOT NULL,
  `name` char(52) NOT NULL,
  `population` int(11) NOT NULL DEFAULT '0',
  `age` tinyint(4) DEFAULT '0' COMMENT '字段描述',
  PRIMARY KEY (`code`)
) ENGINE=InnoDB;
```

