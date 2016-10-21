
##token-auth-apns-http2-proxy

*  A example of proxy to APNs with http2 and token-based authoriation
*  项目描述:  做一个APN推送项目，目的是实现一个服务，从前端HTTP1.1代理到后端苹果的HTTP2 APNs， 由于没有找到完整可用的基于token认证的例子, 经过探索，在现有开源库的基础上做修改， 实现了该服务。 
*  示例代码里的私钥，iss, kid需要替换成从苹果获取的相关信息

###github.com/dgrijalva/jwt-go 的改动
  
*  ecdsa_utils.go 新增函数
<pre>
// Parse PEM encoded Elliptic Curve Private Key Structure
func ParsePKCS8PrivateKeyFromPEM(key []byte) (*ecdsa.PrivateKey, error) {
    var err error

    // Parse PEM block
    var block *pem.Block
    if block, _ = pem.Decode(key); block == nil {
        return nil, ErrKeyMustBePEMEncoded
    }

    // Parse the key
    var parsedKey interface{}
    if parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
        return nil, err
    }

    var pkey *ecdsa.PrivateKey
    var ok bool
    if pkey, ok = parsedKey.(*ecdsa.PrivateKey); !ok {
        return nil, ErrNotECPrivateKey
    }

    return pkey, nil
}
</pre>

*  token.go
	注释 //"typ": "JWT",


###github.com/sideshow/apns2 改动

*   client.go setHeaders() 里增加
*   
<pre>
    if n.Authorization != "" {
        r.Header.Set("authorization", n.Authorization)
    }
</pre>
	
*  notification 增加
<pre>
    //The provider token that authorizes APNs to send push notifications for the specified topics.
    //The token is in Base64URL-encoded JWT format, specified as bearer <provider token>.
    //When the provider certificate is used to establish a connection, this request header is ignored.
    Authorization string
</pre>
