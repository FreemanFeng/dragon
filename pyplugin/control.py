#+=========================================================================
#|   Copyright 2019 Freeman Feng<freeman@nuxim.cn>
#|
#|   Licensed under the Apache License, Version 2.0 (the "License");
#|   you may not use this file except in compliance with the License.
#|   You may obtain a copy of the License at
#|
#|       http://www.apache.org/licenses/LICENSE-2.0
#|
#|   Unless required by applicable law or agreed to in writing, software
#|   distributed under the License is distributed on an "AS IS" BASIS,
#|   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#|   See the License for the specific language governing permissions and
#|   limitations under the License.
#+=========================================================================

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
