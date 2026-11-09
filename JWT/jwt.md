# 简介
JWT 全名 JSON Web Tokens，关于 JWT 的详细介绍可以在[jwt.io](JSON Web Token Introduction - jwt.io)查看，它是一种开放的，安全的，紧凑的，以 JSON 对象为载体在服务双方之间传输信息的方式，它的特点就是安全性高，内容防篡改，消耗低。



## 结构

在 RFC 标准中，JWT 由以下三个部分组成：

- **Header 头部**
- **Payload 载荷**
- **Signature 签名**

然后每一个部分用一个点`.`来分隔，最后组成一个字符串，格式就是`header.payload.signature`吗，这就是一个 JWT 令牌的标准结构，接下来就一个个讲解每个结构的作用。

### 头部

头部只是声明一些基本信息，通常由两部分组成，令牌的类型，和签名所使用的加密算法，例如下方：

```
{
  "alg": "HS256",
  "typ": "JWT"
}
```

以上的信息大致就是，令牌的类型为 JWT，签名部分所使用的加密算法为 HS256，最后再将 JSON 对象通过`Base64Url`编码成字符串，该字符串就是 JWT 的头部。

### 载荷

JWT 的第二部分是载荷部分，主要包含声明(`claims`)部分，声明部分通常是关于一个实体的数据，比如一个用户。关于声明的类型总共有三种：

- `reigstered`：`Registered claims`代表着 一些预定义的声明，一些并不强制使用但是仍然推荐使用，例如：`iss`(issuer 签发者)，`exp`(expiration time 过期时间) ，`aud`(audience 受众)。
- `public`：`Public claims`是可以由使用 JWT 的人随意定义的，最好要避免和其他声明部分冲突。
- `private claims`：这部分的声明同样也是自定义的，通常用于在服务双方共享一些信息。

一个载荷示例如下：

```
{
  "sub": "1234567890",
  "name": "John Doe",
  "admin": true
}
```

该 JSON 对象将会通过`Base64Url`被编码成字符串，从而组成 JWT 的第二部分。

### 签名

在获得了编码的头部和编码的载荷部分后，就可以通过头部所指明的签名算法根据前两个部分的内容再加上密钥进行加密签名，所以一旦 JWT 的内容有任何变化，解密时得到的签名都会不一样，同时如果是使用私钥，也可以对 JWT 的签发者进行验证。



```
sign = HMACSHA256(
  base64UrlEncode(header) + "." +
  base64UrlEncode(payload),
  secret)
```

例如下方的例子：



```
Header
{
  "alg": "HS256",
  "typ": "JWT"
}

Payload
{
  "alg": "HS256",
  "typ": "JWT"
}

Verify Signature
HMACSHA256(
  base64UrlEncode(header) + "." +
  base64UrlEncode(payload),
  your secret
)
```

最后得到的输出就是一个由三个`base64Url`字符串组成且由`.`分隔的字符串，大概长下面这样



```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.
eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.
cThIIoDvwdueQB468K5xDc5633seEFoqwxjF_xSJyQQ
```

## 工作原理

在身份验证中，当用户使用凭据成功登录时，将返回一个 JSON Web 令牌。由于令牌是凭证，因此必须非常小心地防止出现安全问题。一般来说，令牌的保存时间不应超过所需的时间。然后无论何时用户想要访问受保护的路由和资源，在发起请求时就必须携带上 token，通常都是在请求头中的`Authorization` header 中的`Bearer schema`，例如下方：



```
Authorization: Bearer <token>
```

服务器在收到 JWT 后，会对其进行有效性验证，例如内容有篡改，token 已过期等等，如果验证通过就可以顺利的访问资源。虽然 JWT 中可以携带一些基本信息，但是依旧建议信息不要太大。



## 选择签名算法

可用的签名算法有好几种，在使用之前应该先了解下它们之间的区别以便更好的去选择签名算法，它们之间最大的不同就是**对称加密**和**非对称加密**。

最简单的对称加密算法`HSA`，让任何`[]byte`都可以用作有效的密钥，所以计算速度稍微快一点。在生产者和消费者双方都是可以被信任的时候，对称加密算法的效率是最高的。不过由于签名和验证都使用相同的密钥，因此无法轻松的分发用于验证的密钥，毕竟签名的密钥也是同一个，签名泄露了则 JWT 的安全性就毫无意义。

非对称加密签名方法，例如`RSA`，使用不同的密钥来进行签名和验证 token，这使得生成带有私钥的令牌成为可能，同时也允许任何使用公钥验证的人正常访问。

不同的签名算法所需要的密钥的类型也不同，下面给出一些常见签名算法的类型：

- `HMAC`：对称加密，需要类型`[]byte`的值用于签名和验证。 (`HS256`,`HS384`,`HS512`)
- `RSA`：非对称加密，需要`*rsa.PrivateKey`类型的值用于签名，和`*rsa.PublicKey`类型的值用于验证。(`RS256`,`RS384`,`RS512`)
- `ECDSA`：非对称加密，需要`*ecdsa.PrivateKey`类型的值用于签名，和`*ecdsa.PublicKey`类型的值用于验证。(`ES256`,`ES384`,`ES512`)
- `EdDSA`：非对称加密，需要`ed25519.PrivateKey`类型的值用于签名和`ed25519.PublicKey` 类型的值用于验证。(`Ed25519`)



## 完整示例

```go
package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)

func sign(input, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(input))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func main() {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"user1","exp":1735689600}`))

	input := header + "." + payload
	token := input + "." + sign(input, "secret")

	parts := strings.Split(token, ".")
	fmt.Println("segments:", len(parts))
	fmt.Println("signature matches:", sign(parts[0]+"."+parts[1], "secret") == parts[2])
	fmt.Println("wrong key rejected:", sign(parts[0]+"."+parts[1], "wrong") != parts[2])
}
```

预期输出：

```
segments: 3
signature matches: true
wrong key rejected: true
```



## 常见错误与踩坑

- **校验时必须固定签名算法。** 拒绝与预期不符的 `alg`，绝不接受 `alg: none`。
- **要校验 `exp` / `nbf`。** 大多数库会代劳，但前提是你使用解析后的 claims，而不是自己解码 payload。
- **不要硬编码密钥。** 从环境变量或密钥管理服务加载，并定期轮换。
- **令牌要短有效期。** 访问令牌配合刷新机制；否则很难撤销。
- **payload 是可读的。** 它只是 base64 而不是加密，敏感信息不要放进 claims。

