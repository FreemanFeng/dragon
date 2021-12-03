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

import com.dragon.httpserver.helper.HttpHelper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.HashMap;

public class OnSending {
    private static final Logger LOGGER = LoggerFactory.getLogger(OnSending.class);
    public static void call(HashMap<String, Object> m, HashMap<String, Object> c, HashMap<String, Object> r) {
        String method = HttpHelper.getMethod(c);
        String url = HttpHelper.getURL(m, c);
        HashMap<String, String> headers = HttpHelper.getHeaders(m, c);
        String body = HttpHelper.getBody(m, c);
        LOGGER.info("OnSending");
    }
}
