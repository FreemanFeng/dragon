//   Copyright 2019 Freeman Feng<freeman@nuxim.cn>
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
package demo.proto.http;

import com.dragon.httpserver.enums.Decoder;
import com.dragon.httpserver.enums.Encoder;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.HashMap;

public class OnReceived  {
    private static final Logger LOGGER = LoggerFactory.getLogger(OnReceived.class);
    public static void call(HashMap<String, Object> m, HashMap<String, Object> c, HashMap<String, Object> r) {
        Object a = r.get("resp");
        HashMap<String, Object> ma = Decoder.INSTANCE.decodeMap(a);
        ma.put("hello2", "world");
        r.put("resp", ma);
        LOGGER.info("OnReceived");
    }
}
