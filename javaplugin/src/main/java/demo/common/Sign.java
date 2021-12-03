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
package demo.common;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.util.DigestUtils;

import java.nio.charset.StandardCharsets;
import java.util.Arrays;

public class Sign  {
    private static final Logger LOGGER = LoggerFactory.getLogger(Sign.class);
    public static String call(String ts, String nonce, String id,String secret, byte[] body) {
        String s = ts+nonce+id+secret+ Arrays.toString(body);
        String sign = DigestUtils.md5DigestAsHex(s.getBytes(StandardCharsets.UTF_8));
        LOGGER.info("Sign {}", sign);
        return sign;
    }
}
