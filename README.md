# gophp
unserialize php serialize output to map, with go

数组只能以int, string为key，不支持循环引用
object视为为普通php数组，即go map

## 用法

```
unserializeObj := gophp.Parse(`a:7:{i:0;b:1;i:3;b:0;i:4;i:123;i:5;d:3.1415;s:3:"str";s:43:"a ba a ba, "(){}?:"><?!@#$%^&*-=_+|'\中国";s:4:"null";N;s:3:"obj";O:3:"cls":1:{s:5:"prop1";s:3:"ppp";}}`)

json.Marshal(unserializeObj)
//会输出{"0":true,"3":false,"4":123,"5":3.1415,"null":null,"obj":{"prop1":"ppp"},"str":"a ba a ba, \"(){}?:\"\u003e\u003c?!@#$%^\u0026*-=_+|'\\中国"}
```
