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
