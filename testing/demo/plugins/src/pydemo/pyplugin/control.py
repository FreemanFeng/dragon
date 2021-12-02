import json

from twisted.web import server, resource
from twisted.internet import reactor, endpoints

PORT = "8088"


class Control(resource.Resource):
    isLeaf = True
    callback = {"init": None, "run": None}

    def render_GET(self, request):
        p = request.path.decode("utf-8").split("/")
        request.setHeader(b"content-type", b"application/json;chartset=uft-8")
        f = self.callback["init"]
        if len(p) < 2 or p[1] != "init" or f is None:
            return "service not init!".encode("utf-8")
        r = f()
        keys = [k for k in r]
        content = json.dumps(keys)
        return content.encode("utf-8")

    def render_POST(self, request):
        p = request.path.decode("utf-8").split("/")
        request.setHeader(b"content-type", b"application/json;chartset=uft-8")
        f = self.callback["run"]
        if len(p) < 3 or p[1] != "run" or f is None:
            return "service not init!".encode("utf-8")
        service = p[2]
        k = request.content.read()
        print(service, "Run Callback Function", f)
        content = f(service, k)
        return content


if __name__ == "__main__":
    endpoints.serverFromString(reactor, "tcp:%s" % PORT).listen(server.Site(Control()))
    reactor.run()
