/*=============================================================*/
/* Table: manager                                              */
/*=============================================================*/
DROP TABLE IF EXISTS manager;
CREATE TABLE manager
(
  id bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
  username varchar(50) NOT NULL COMMENT '用户名',
  password varchar(50) NOT NULL COMMENT '密码',
  status int(11) NOT NULL DEFAULT '0' COMMENT '状态，0：未启用，1：正常，2：冻结，9：关闭',
  type int(11) NOT NULL DEFAULT '1' COMMENT '管理员类型，0：其他，1：普通一级管理员，2：二级管理员，3：三级管理员，7：高级管理员，8：系统管理员，9：超级用户',
  realname varchar(100) COMMENT '真实姓名',
  email varchar(200) COMMENT '邮箱',
  mobile varchar(20) COMMENT '手机号',
  createTime bigint(15) NOT NULL COMMENT '创建时间',
  updateTime bigint(15) COMMENT '修改时间',
  remark varchar(200) COMMENT '备注',
  json varchar(500) COMMENT 'json',
  PRIMARY KEY (id),
  UNIQUE (username)
)
ENGINE = InnoDB;

ALTER TABLE manager COMMENT '管理员用户表';

/*==============================================================*/
/* Index: createTime                                         */
/*==============================================================*/
create index createTime on manager
(
   createTime
);

/*==============================================================*/
/* Index: status                                         */
/*==============================================================*/
create index status on manager
(
   status
);

/*==============================================================*/
/* Index: type                                         */
/*==============================================================*/
create index type on manager
(
   type
);

/*=============================================================*/
/* Table: user                                              */
/*=============================================================*/
DROP TABLE IF EXISTS user;
CREATE TABLE user
(
  id bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'id',
  username varchar(100) NOT NULL COMMENT '用户名',
  password varchar(100) NOT NULL COMMENT '密码',
  status int(11) NOT NULL DEFAULT '0' COMMENT '状态，0：未验证，1：正常，2：冻结，3：需认证状态，9：关闭',
  type int(11) NOT NULL DEFAULT '1' COMMENT '用户类型，0：其他，1：普通用户',
  nickname varchar(100) COMMENT '昵称',
  realname varchar(100) COMMENT '真实姓名',
  email varchar(200) COMMENT '邮箱',
  mobile varchar(20) COMMENT '手机号码',
  sex int(1) COMMENT '性别，0：女，1：男',
  createTime bigint(15) NOT NULL COMMENT '创建（注册）时间',
  firstLoginTime bigint(15) COMMENT '首次登录时间',
  firstLoginIP varchar(100) COMMENT '首次登录IP',
  lastLoginTime bigint(15) COMMENT '上次登录时间',
  lastLoginIP varchar(100) COMMENT '上次登录IP',
  validatedTime bigint(15) COMMENT '完成验证时间',
  signature varchar(500) COMMENT '签名',
  remark varchar(200) COMMENT '备注',
  json varchar(500) COMMENT 'json字段',
  PRIMARY KEY (id),
  UNIQUE (username)
)
ENGINE = InnoDB;

ALTER TABLE user COMMENT '用户表';

/*==============================================================*/
/* Index: createTime                                         */
/*==============================================================*/
create index createTime on user
(
   createTime
);

/*==============================================================*/
/* Index: status                                         */
/*==============================================================*/
create index status on user
(
   status
);

/*==============================================================*/
/* Index: type                                         */
/*==============================================================*/
create index type on user
(
   type
);