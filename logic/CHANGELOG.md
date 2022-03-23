版本号`major.minor.patch`具体规则如下：
- major：主版本号，如有重大版本重构则该字段递增，通常各主版本间接口不兼容。
- minor：次版本号，各次版本号间接口保持兼容，如有接口新增或优化则该字段递增。
- patch：补丁号，如有功能改善或缺陷修复则该字段递增。

## version 3.1.4 @2021.10.28

优化 log 模块

配置文件更新  
新增
```toml
[log]
    Level="debug"
    Mode="console"
    Path=""
    Display="json"
```

## version 3.1.3 @2021.7.29

修复当 addr 为空产生 panic

## version 3.1.2 @2021.7.22

**Feature**
- 支持根据 server 使用指定位置的 comet rpc @3.1.2 2021.7.22


## example x.x.x @yy.mm.dd

**Feature**

**Bug Fixes**

**Improvement**

**Breaking Change**
