import json

from io import BytesIO
from twisted.web import server
from twisted.internet import reactor, endpoints
from .control import Control
import base64


class Plugin(object):
    def __init__(self, port, init, run):
        self.port = port
        self.funcs = dict()
        Control.callback["init"] = init
        Control.callback["run"] = run
        print("init port", port, "service", init, run)

    def init(self, funcs):
        self.funcs = funcs

    def run(self, name, b):
        if name not in self.funcs:
            err = "service %s is not recognized!" % name
            return err.encode("utf-8")
        try:
            ss = json.loads(b)
            if type(ss) is not list:
                return "list of bytes is wanted".encode("utf-8")
            return self.process(name, ss)
        except Exception as e:
            return str(e).encode("utf-8")

    def process(self, name, ss):
        n = int(len(ss) / 2)
        ms = list()
        rs = list()
        f = self.funcs[name]
        for i in range(n):
            k = ss[i * 2]
            v = ss[i * 2 + 1]
            if k == "MapBytes":
                try:
                    b = base64.b64decode(v)
                    x = json.loads(b)
                    m = {k: json.loads(base64.b64decode(v)) for k, v in x.items()}
                    ms.append(m)
                except Exception as e:
                    print("plugin func", name, "No.", i, "argument deserialized failed", e)
                    return str(e).encode("utf-8")
            if k == "MAP":
                try:
                    b = base64.b64decode(v)
                    m = json.loads(b)
                    ms.append(m)
                except Exception as e:
                    print("plugin func", name, "No.", i, "argument deserialized failed", e)
                    return str(e).encode("utf-8")
            if k == "ListBytes":
                try:
                    b = base64.b64decode(v)
                    x = json.loads(b)
                    m = [base64.b64decode(v) for v in x]
                    ms.append(m)
                except Exception as e:
                    print("plugin func", name, "No.", i, "argument deserialized failed", e)
                    return str(e).encode("utf-8")
            if k == "LIST":
                try:
                    b = base64.b64decode(v)
                    m = json.loads(b)
                    ms.append(m)
                except Exception as e:
                    print("plugin func", name, "No.", i, "argument deserialized failed", e)
                    return str(e).encode("utf-8")
            if k == "STR":
                try:
                    ms.append(v)
                except Exception as e:
                    print("plugin func", name, "No.", i, "argument decode failed", e)
                    return str(e).encode("utf-8")
            if k == "BYTES":
                ms.append(v.encode("utf-8"))
            if k == "INT" or k == "UINT" or k == "FLOAT":
                try:
                    ms.append(int(v))
                except Exception as e:
                    print("plugin func", name, "No.", i, "argument decode failed", e)
                    return str(e).encode("utf-8")
            if k == "BOOL":
                try:
                    if v == "true" or v == "TRUE" or v == "True":
                        ms.append(True)
                    else:
                        ms.append(False)
                except Exception as e:
                    print("plugin func", name, "No.", i, "argument decode failed", e)
                    return str(e).encode("utf-8")
        if name.find(".") > 0:
            f(*ms)
            for m in ms:
                try:
                    x = json.dumps(m, ensure_ascii=False).encode("utf-8")
                    rs.append(base64.b64encode(x).decode("utf-8"))
                except Exception as e:
                    return str(e).encode("utf-8")
            try:
                return json.dumps(rs, ensure_ascii=False).encode("utf-8")
            except Exception as e:
                return str(e).encode("utf-8")
        return f(*ms)

    def serve(self):
        endpoints.serverFromString(reactor, "tcp:%s" % str(self.port)).listen(server.Site(Control()))
        reactor.run()
