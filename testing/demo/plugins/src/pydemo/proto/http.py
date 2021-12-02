import json


def onReadySending(m, c, r: dict):
    print("onReadySending message templates", m, "config", c, "request attributes", r)
    ma = m["a"]
    print("before update a", m["a"])
    if type(ma) is not dict:
        print("a is map expected")
        return
    ma["hello2"] = "world"
    print("after update a", m["a"])


def onReceived(m, c, r: dict):
    print("onReceived message templates", m, "config", c, "response", r)
    mm = r["resp"]
    print("before update resp", r["resp"])
    mm["hello2"] = "world"
    print("after update resp", r["resp"])


def onError(m, c, r: dict):
    print("OnError message templates", m, "config", c, "response", r)
    if "code" in r:
        if r["code"] == "400":
            r["code"] = "500"
