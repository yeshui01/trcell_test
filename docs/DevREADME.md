# server_wiki

服务器说明

## 协议格式规范
1. 协议号定义
    - 对客户端的请求协议 ECMsg+模块名字+OpCode,例如ECMsgPlayerLogin
    - 对客户端的推送协议 ECMsg+模块名字+Push+OpCode,例如
    ECMsgPlayerPushLogin
    - 服务器内部请求协议 ESMsg+模块名字+OpCode,例如ESMsgPlayerLogin
    - 服务器内部推送协议 ESMsg+模块名字+Push+OpCode,例如ESMsgPlayerPushLogin
---
2. proto定义
    - 对客户端的协议文件名 c_+模块名字.proto,例如 c_player.proto
    - 客户端请求协议 协议号+Req,例如 message ECMsgPlayerLoginReq
    - 客户端回复协议 协议号+Rsp,例如 message ECMsgPlayerLoginRsp
    - 客户端推送协议 协议号+Notify 例如 message ECMsgPlayerPushLoginNotify
    - 对服务器协议文件名 s_+模块名字.proto,例如 s_player.proto
    - 服务器请求协议 协议号+Req,例如 message ESMsgPlayerLoginReq
    - 服务器回复协议 协议号+Rep,例如 message ESMsgPlayerLoginRep
    - 服务器推送协议 协议号+Notify 例如 message ESMsgPlayerPushLoginNotify

## 文件格式规范
1. 代码文件命名格式
    - 文件名全部采用小写加下划线的格式命名
---
2. 文件夹命名
    - 业务层，go的包名和文件夹名一致,全部小写,不加下划线

## 代码格式规范
1. 函数命名
    - 私有函数小写开头,小驼峰格式,不要加下划线
    - 对外公开的函数接口,大写驼峰,不要加下划线
    - 函数参数保持简短,超过6个参数的,封装成结构集合的参数
2. 防御性编程
    - 凡是指针,例如配置,获取的时候必须判空
    - 客户端参数不要信任,需要自己验证判断是否合法,不要拿着就用!!!
3. 线程安全函数
    - 区分内部和对外的接口
    - 带有返回值的，一定是不安全的，需要考虑使用方式
---

## 玩家数据持久化模块
1. 定义数据表
2. 使用ormtools工具生成对应的orm关联代码
3. 定义玩家数据私有成员
4. 定义自己的模块,实现playermodule的业务层接口
5. 实例化模块,在player的InitModule中注册实例化
6. 模块访问,全部使用player的GetModule接口访问,不要直接自定义变量,这里需要统一管理

---
## 业务层容器使用
- 为了保证安全和可追踪,业务层的代码，容器不要使用标准的map和切片,统一使用ustd包里面的Map和Vector

---
## 业务层时间接口
- 统一使用trframe.GetFrameSysNowTime()获取当前时间戳,禁止用原生的time.Now().Unix()

---
## 策划表符号规范
- 文件名全部采用小写加下划线分割的格式,字段名也是,如果发现不是请让策划改正,例如 union_boss.xlsx