gometa是一个基于go的元数据编程框架，很适合做企业应用软件开发。
gometa目标是让用户通过简单地配置即可完成增删改查功能，而专注于业务逻辑的编写和开发。

当前,gometa支持
1.数据库:MongoDB
2.前端:YUI

gometa演示项目github地址:https://github.com/hongjinqiu/gometa-financedemo
演示项目url地址:http://112.126.81.24/

依赖:
1.mgo,(go访问MongoDB驱动)

    mkdir $GOPATH/src/gopkg.in
    cd $GOPATH/src/gopkg.in
    git clone https://github.com/go-mgo/mgo.git mgo.v2
    cd mgo.v2
    git branch v2 remotes/origin/v2
    git checkout v2

2.otto,(go的js解析)

    go get -v github.com/robertkrimen/otto

3.godec,(go精确数值计算)

    go get -v code.google.com/p/godec/dec

gometa设计说明
gometa由数据源模型,数据集模型,字段模型,关联模型,列表模型,表单模型,权限模型组成。

1.数据源模型简介
数据源模型由数据集,字段模型组成。
数据源模型通常定义了数据源ID,业务名称等数据结构。
大部分的数据源模型通常由单个主数据集组成，反映到关系数据库层面就是一个单表，用MongoDB来存储就是:{"A": {"id": xxx, "key1": value1, "key2": value2,....}}
少部分数据源模型由一个主数据集和1、2个分录数据集，反映到关系数据库层面就是一个主从表，两个表，用用MongoDB来存储就是:{"A": {"id": xxx, "key1": value1, "key2": value2,....}, "B": [{"id": xxx, "key1": value1, ...}, {"id": xxx, "key1": value1, ...}, ....]}
数据源模型是gometa编程框架的基石，定义好数据源模型后，就可以在关联模型，列表模型，表单模型中方便的引用。

2.数据集模型简介
数据集分为主数据集，分录数据集两种。
主数据集与分录数据集都定义了固定字段，业务字段。
通常每个数据集的固定字段都相同，这样在编程时，可以按固定字段进行方便的字段引用，例如，可以定义所有的数据集都包含createBy字段，进行字段访问时，便可以直接写死createBy,而不是取得id字段存储的值进行访问。
主数据集与分录数据集的不同在于，分录存储的是列表。

3.字段模型简介
字段模型定义了字段ID，显示名称，字段数据类型，字段长度，默认值表达式，是否只读，是否允许为空，........
字段之间允许进行继承，因而可以大大减少配置量。
为了减少开发的工作量，gometa定义了两个字段池文件fieldpool.xml，business_fieldpool.xml。
fieldpool.xml，属于字段类型字段池，定义了固定字段，STRING_FIELD,INT_FIELD等字段。
business_fieldpool.xml，属于业务字段池，开发一个系统时，通常把所有的字段都定义在这里，这样可以避免字段不一致的问题，比如开发中常见的问题，有的表，字段叫SOURCE_NO，另一张表，字段叫SOURCE_BILL_NO，其实是同一个字段。

4.关联模型简介
关联模型配置在字段模型里面，用来表示这个字段是个关联字段，关联到别的数据源模型中，对应到关系数据库层面就是外键。
关联模型的数据结构有：名称，关联的数据源模型ID选择器，关联的数据源模型ID，关联表达式，关联显示字段，关联取值字段，....
关联模型允许用关联表达式来进行约束，当表达式返回true时，这个关联模型才生效，以满足页面上需要打开不同表达式窗口的情况。

5.列表模型简介
列表模型用于数据源模型的页面展示，其数据结构有：数据源模型ID，渲染页面html地址，工具栏，权限，数据来源，列表字段，查询条件。
工具栏：由按钮组成，按钮的数据结构为：按钮名称，按钮响应方式(打开新页面，JS操作，.....)。
权限:当前只支持按单位，按个人查询。
数据来源：一般放空，系统自动从数据源模型配置里面读取。
列表字段：分为auto-column,id-column,string-column,dictionary-column,date-column,.......,一般只需要配置auto-column,字段名称与数据源模型中定义的字段名称一致即可。
查询条件：由查询字段组成，查询字段的定义控件类型，查询方式(eq,in,like,....),一般放空即可,系统会自动生成。

6.表单模型简介
表单模型用于数据源模型的表单展示，其数据结构有：数据源模型ID，渲染页面html地址，工具栏，权限，主数据集列表字段，分录列表字段。
工具栏：由按钮组成，按钮的数据结构为：按钮名称，按钮响应方式(打开新页面，JS操作，.....)。
权限:当前只支持按单位，按个人查询。
主数据集列表字段:通常配置为与数据源模型中字段名称相同即可，再配置占几列。
分录列表字段：通常配置为与数据源模型中字段名称相同即可，页面渲染为带编辑功能的表格。

7.权限模型简介
当前只支持按单位，按个人进行查询。

