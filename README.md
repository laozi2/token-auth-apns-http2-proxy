# apns-http2-proxy-with-go
A example of proxy to APNs with http2 and token-based authoriation



1. 基于token认证的APNs推送服务代理
	示例代码里的私钥，iss, kid需要替换成从苹果获取的相关信息

2. 后台运行，可使用  daemonize
	git clone git://github.com/bmc/daemonize.git
	./configure
	make && make install

	启动 /usr/local/sbin/daemonize -l /var/lock/subsys/go_apn_push -u root -o  /usr/local/go_apns_log/go_apn_push.log -e /usr/local/go_apns_log/go_apn_push.log -c /usr/local/go_apns/  /usr/local/go_apns/go_apn_push

3. github.com/dgrijalva/jwt-go 的改动
	3.1 ecdsa_utils.go 新增函数

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

	3.2 token.go
	注释 //"typ": "JWT",


4. github.com/sideshow/apns2 改动
	4.1 client.go setHeaders() 里增加
	if n.Authorization != "" {
		r.Header.Set("authorization", n.Authorization)
	}
	
	4.2 notification 增加
	
	//The provider token that authorizes APNs to send push notifications for the specified topics.
	//The token is in Base64URL-encoded JWT format, specified as bearer <provider token>.
	//When the provider certificate is used to establish a connection, this request header is ignored.
	Authorization string
