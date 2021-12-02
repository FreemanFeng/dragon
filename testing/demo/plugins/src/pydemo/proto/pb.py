def onReadySending(m, c, r: dict):
    print("onReadySending", m, c, r)


def onReceived(m, c, r: dict):
    print("onReceived", m, c, r)


def onError(m, c, r: dict):
    print("onError", m, c, r)