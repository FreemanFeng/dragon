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
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.HashMap;

public class OnError {
    private static final Logger LOGGER = LoggerFactory.getLogger(OnError.class);
    public static void call(HashMap<String, Object> m, HashMap<String, Object> c, HashMap<String, Object> r) {
        if (r.containsKey("code")) {
            Object a = r.get("code");
            Integer code = Decoder.INSTANCE.decodeInt(a);
            if (code == 400) {
                r.put("code", code);
            }
        }
        LOGGER.info("OnError");
    }
}
