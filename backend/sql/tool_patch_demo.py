# -*- coding: utf-8 -*-
"""构造可安全注入现有 dev 库的演示数据副本：
1. 三处固定 id 的种子插入（sys_tenants/sys_users/sys_user_credentials）追加
   ON CONFLICT DO NOTHING——保留现有登录体系；
2. 其余内容表由预清理脚本先行清空，演示数据完整接管。
"""
import io

SRC = 'sql/postgresql-demo-data.sql'
DST = 'sql/.demo-data-patched.sql'

with io.open(SRC, encoding='utf-8') as f:
    src = f.read()

NEUTRALIZE = [
    'INSERT INTO public.sys_tenants(',
    'INSERT INTO public.sys_users (',
    'INSERT INTO public.sys_user_credentials (',
]

for prefix in NEUTRALIZE:
    start = src.index(prefix)
    end = src.index('\n;\n', start)
    src = (src[:end]
           + '\nON CONFLICT DO NOTHING'
           + src[end:])
    print('neutralized:', prefix[:40])

with io.open(DST, 'w', encoding='utf-8', newline='\n') as f:
    f.write(src)
print('patched copy written:', DST)
