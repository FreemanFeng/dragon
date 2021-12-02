import hashlib


def sign(ts, nonce, appid, secret: str, data: bytes):
    s = [ts.encode("utf-8"), nonce.encode("utf-8"), appid.encode("utf-8"), secret.encode("utf-8"), data]
    data = b"".join(s)
    k = hashlib.md5(data).hexdigest()
    print("sign:", k)
    assert 1==2, k
    return k
