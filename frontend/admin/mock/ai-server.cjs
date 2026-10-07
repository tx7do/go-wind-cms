// 最小 OpenAI 兼容 mock（AI 模块 e2e 测试用）：
//   POST /v1/chat/completions  非流式(JSON) + 流式(SSE data: 分帧)
//   POST /v1/embeddings        确定性 1536 维伪向量
// 语义：
//   - created 必须 int（go-openai 校验 int64）
//   - 问数生成（system 含「只读数据分析引擎」）返回固定合法 SQL，其余回显拼接
// 启动：node mock/ai-server.cjs  （监听 127.0.0.1:18080）
const http = require('http');
const reply = (topic) => '【mock生成】已收到主题「' + String(topic).slice(0, 60) + '」，这是一段用于端到端验证的回复文本。';

http.createServer((req, res) => {
  if (req.method === 'POST' && req.url.includes('/chat/completions')) {
    let body = '';
    req.on('data', (c) => (body += c));
    req.on('end', () => {
      let topic = '', stream = false, isQueryGen = false;
      try {
        const j = JSON.parse(body);
        stream = !!j.stream;
        const msgs = j.messages || [];
        topic = msgs.map(m => m.content).join(' | ');
        isQueryGen = msgs.some(m => (m.content || '').includes('只读数据分析引擎'));
      } catch {}
      const created = Math.floor(Date.now() / 1000);
      if (!stream) {
        const content = isQueryGen ? 'SELECT COUNT(*) AS user_count FROM sys_users' : reply(topic);
        res.writeHead(200, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({
          id: 'mock', object: 'chat.completion', created, model: 'mock-model',
          choices: [{ index: 0, message: { role: 'assistant', content }, finish_reason: 'stop' }],
          usage: { prompt_tokens: 12, completion_tokens: 34, total_tokens: 46 },
        }));
        return;
      }
      res.writeHead(200, { 'Content-Type': 'text/event-stream', 'Cache-Control': 'no-cache' });
      const text = reply(topic);
      const pieces = text.match(/[\s\S]{1,12}/g) || [];
      let i = 0;
      const timer = setInterval(() => {
        if (i < pieces.length) {
          const frame = { id: 'mock', object: 'chat.completion.chunk', created, model: 'mock-model',
            choices: [{ index: 0, delta: { content: pieces[i] }, finish_reason: null }] };
          res.write('data: ' + JSON.stringify(frame) + '\n\n');
          i++;
        } else {
          const done = { id: 'mock', object: 'chat.completion.chunk', created, model: 'mock-model',
            choices: [{ index: 0, delta: {}, finish_reason: 'stop' }],
            usage: { prompt_tokens: 12, completion_tokens: 34, total_tokens: 46 } };
          res.write('data: ' + JSON.stringify(done) + '\n\n');
          res.write('data: [DONE]\n\n');
          res.end();
          clearInterval(timer);
        }
      }, 40);
    });
    return;
  }
  if (req.method === 'POST' && req.url.includes('/v1/embeddings')) {
    let body = '';
    req.on('data', (c) => (body += c));
    req.on('end', () => {
      let inputs = [];
      try { const j = JSON.parse(body); inputs = Array.isArray(j.input) ? j.input : [String(j.input)]; } catch {}
      const created = Math.floor(Date.now() / 1000);
      const data = inputs.map((text, idx) => {
        const vec = new Array(1536).fill(0).map((_, d) => {
          let h = 0;
          for (let k = 0; k < text.length; k++) h = (h * 31 + text.charCodeAt(k) * (d + 1)) % 997;
          return (h % 1000) / 1000 - 0.5;
        });
        return { object: 'embedding', index: idx, embedding: vec };
      });
      res.writeHead(200, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify({ object: 'list', data, model: 'mock-embedding', usage: { prompt_tokens: inputs.join('').length, total_tokens: inputs.join('').length } }));
    });
    return;
  }
  res.writeHead(404); res.end('{}');
}).listen(18080, () => console.log('mock openai on 127.0.0.1:18080'));
