// 最小 OpenAI 兼容 mock（AI 模块 e2e 测试用）：
//   POST /v1/chat/completions  非流式(JSON) + 流式(SSE data: 分帧) + Function Call 两段式
//   POST /v1/embeddings        确定性 1536 维伪向量
// 语义：
//   - created 必须 int（go-openai 校验 int64）
//   - 问数生成（system 含「只读数据分析引擎」）返回固定合法 SQL
//   - 请求带 CMS 工具清单（query_*/get_content_overview）时走 Function Call 两段式：
//     首轮流式回 delta.tool_calls；携带工具结果（role=tool）的次轮回文本总结
// 启动：node mock/ai-server.cjs  （监听 127.0.0.1:18080）
const http = require('http');
const reply = (topic) => '【mock生成】已收到主题「' + String(topic).slice(0, 60) + '」，这是一段用于端到端验证的回复文本。';

function chatHandler(req, res) {
  let body = '';
  req.on('data', (c) => (body += c));
  req.on('end', () => {
    let topic = '', stream = false, isQueryGen = false, wantsTool = false, hasToolResult = false;
    let toolSummary = '';
    try {
      const j = JSON.parse(body);
      stream = !!j.stream;
      const msgs = j.messages || [];
      topic = msgs.map(m => m.content).join(' | ');
      isQueryGen = msgs.some(m => (m.content || '').includes('只读数据分析引擎'));
      wantsTool = (j.tools || []).some(t => {
        const n = ((t.function || {}).name) || '';
        return n === 'get_content_overview' || n.startsWith('query_');
      });
      hasToolResult = msgs.some(m => m.role === 'tool');
      const toolMsg = msgs.filter(m => m.role === 'tool').pop();
      toolSummary = String(toolMsg ? toolMsg.content : '').slice(0, 150);
    } catch {}
    const created = Math.floor(Date.now() / 1000);

    // ── Function Call 两段式（流式） ──
    if (stream && wantsTool && !hasToolResult) {
      res.writeHead(200, { 'Content-Type': 'text/event-stream', 'Cache-Control': 'no-cache' });
      const frame = {
        id: 'mock', object: 'chat.completion.chunk', created, model: 'mock-model',
        choices: [{ index: 0, delta: { role: 'assistant', tool_calls: [{ index: 0, id: 'call_mock_1', type: 'function', function: { name: 'get_content_overview', arguments: '{}' } }] }, finish_reason: 'tool_calls' }],
      };
      res.write('data: ' + JSON.stringify(frame) + '\n\n');
      res.write('data: [DONE]\n\n');
      res.end();
      return;
    }
    if (stream && wantsTool && hasToolResult) {
      res.writeHead(200, { 'Content-Type': 'text/event-stream', 'Cache-Control': 'no-cache' });
      const text = '【mock工具总结】已读取内容概览工具结果：' + toolSummary;
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
            usage: { prompt_tokens: 20, completion_tokens: 40, total_tokens: 60 } };
          res.write('data: ' + JSON.stringify(done) + '\n\n');
          res.write('data: [DONE]\n\n');
          res.end();
          clearInterval(timer);
        }
      }, 40);
      return;
    }

    // ── 非流式 ──
    if (!stream) {
      if (wantsTool && !hasToolResult) {
        res.writeHead(200, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({
          id: 'mock', object: 'chat.completion', created, model: 'mock-model',
          choices: [{ index: 0, message: { role: 'assistant', content: null, tool_calls: [{ id: 'call_mock_1', type: 'function', function: { name: 'get_content_overview', arguments: '{}' } }] }, finish_reason: 'tool_calls' }],
          usage: { prompt_tokens: 10, completion_tokens: 5, total_tokens: 15 },
        }));
        return;
      }
      if (wantsTool && hasToolResult) {
        res.writeHead(200, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({
          id: 'mock', object: 'chat.completion', created, model: 'mock-model',
          choices: [{ index: 0, message: { role: 'assistant', content: '【mock工具总结】已读取内容概览工具结果：' + toolSummary }, finish_reason: 'stop' }],
          usage: { prompt_tokens: 20, completion_tokens: 40, total_tokens: 60 },
        }));
        return;
      }
      const content = isQueryGen ? 'SELECT COUNT(*) AS user_count FROM sys_users' : reply(topic);
      res.writeHead(200, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify({
        id: 'mock', object: 'chat.completion', created, model: 'mock-model',
        choices: [{ index: 0, message: { role: 'assistant', content }, finish_reason: 'stop' }],
        usage: { prompt_tokens: 12, completion_tokens: 34, total_tokens: 46 },
      }));
      return;
    }

    // ── 流式（普通对话回显） ──
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
}

http.createServer((req, res) => {
  if (req.method === 'POST' && req.url.includes('/chat/completions')) {
    chatHandler(req, res);
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
