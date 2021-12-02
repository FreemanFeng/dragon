from .common import sign
from .proto import http, pb
from .pyplugin.plugin import Plugin


class Demo(object):
    def __init__(self, port):
        self.plugin = Plugin(port, self.init, self.run)

    def init(self):
        services = {
            "http.OnReadySending": http.onReadySending,
            "http.OnReceived": http.onReceived,
            "http.OnError": http.onError,
            "pb.OnReadySending": pb.onReadySending,
            "pb.OnReceived": pb.onReceived,
            "pb.OnError": pb.onError,
            "Sign": sign.sign
        }
        print("init services list:", services)
        self.plugin.init(services)
        return services

    def run(self, name, b):
        print("run service:", name)
        return self.plugin.run(name, b)

    def serve(self):
        self.plugin.serve()


def main(port=8086):
    demo = Demo(port)
    demo.serve()
