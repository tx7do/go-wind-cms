<script setup lang="ts">
import { ref } from 'vue'
import { ContentViewer } from '@/components/ContentViewer'

const markdownContent = ref(`# ContentViewer Demo

Welcome to the **ContentViewer** component demonstration.

## Features

This component supports various Markdown features:

### Text Formatting

- **Bold text**
- *Italic text*
- ***Bold and italic***
- ~~Strikethrough~~

### Links and Images

[Visit GitHub](https://github.com)

![Sample Image](/logo.png)

### Code

Inline code: \`const viewer = new ContentViewer()\`

Code block:
\`\`\`javascript
function renderContent(content) {
  return marked(content);
}
\`\`\`

### Quotes

> This is a blockquote.
> It can span multiple lines.

### Lists

#### Unordered List
- First item
- Second item
  - Nested item 1
  - Nested item 2
- Third item

#### Ordered List
1. First step
2. Second step
3. Third step

#### Task List
- [x] Completed task
- [ ] Pending task
- [ ] Another pending task

### Tables

| Feature | Support | Notes |
|---------|---------|-------|
| Markdown | ✅ | Full GFM support |
| HTML | ✅ | Direct rendering |
| Text | ✅ | Preformatted display |

### Horizontal Rule

---

### Additional Content

Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.

#### Mathematical Expression (if supported)

E = mc²

#### Emoji (if supported)

😀 😁 😂 😃 😄 😅 😆 😇 😈 😉 😊 😍 😘 😗 😙 😚 😜 😝

👋 👍 👎 👌 ✌️ ✋ 👏 🙌 👐 🤝 💪

❤️ 🧡 💛 💚 💙 💜 🖤 🤍 🤎 💔 💕 💖 💗

🎉 🎊 🎈 🎁 🎀 🎂 🎃 🎄 ✨ 🌟 ⭐ 💫 🚀 🚁 🚀

## Diagrams

### Flowchart

\`\`\`mermaid
graph TD
    A[Start] --> B{Decision}
    B -->|Yes| C[Process A]
    B -->|No| D[Process B]
    C --> E[End]
    D --> E
    style A fill:#90EE90
    style E fill:#FFB6C6
\`\`\`

### Sequence Diagram

\`\`\`mermaid
sequenceDiagram
    participant User
    participant Client
    participant Server
    User->>Client: Click Button
    Client->>Server: Send Request
    Server-->>Client: Return Response
    Client-->>User: Display Result
\`\`\`

### State Diagram

\`\`\`mermaid
stateDiagram-v2
    [*] --> Idle
    Idle --> Loading: Start
    Loading --> Success: Complete
    Loading --> Error: Fail
    Success --> [*]
    Error --> Idle: Retry
\`\`\`

### Class Diagram

\`\`\`mermaid
classDiagram
    class Animal {
        name: string
        eat()
        sleep()
    }
    class Dog {
        breed: string
        bark()
    }
    class Cat {
        color: string
        meow()
    }
    Animal <|-- Dog
    Animal <|-- Cat
\`\`\`

### Gantt Chart

\`\`\`mermaid
gantt
    title Project Timeline
    dateFormat YYYY-MM-DD
    section Design
    UI Design :des1, 2024-01-01, 30d
    API Design :des2, 2024-01-15, 20d
    section Development
    Frontend :dev1, 2024-02-01, 45d
    Backend :dev2, 2024-02-01, 60d
    section Testing
    Testing :test1, 2024-04-01, 30d
\`\`\`

### Pie Chart

\`\`\`mermaid
pie title Browser Distribution
    "Chrome" : 45
    "Firefox" : 25
    "Safari" : 20
    "Edge" : 10
\`\`\`
`)

const htmlContent = ref(`
<h1>HTML Content</h1>
<p>This is <strong>HTML</strong> content rendered directly.</p>
<ul>
  <li>Item 1</li>
  <li>Item 2</li>
  <li>Item 3</li>
</ul>
<p>It supports all standard HTML elements with beautiful styling.</p>
`)

const textContent = ref(`Plain Text Content

This is plain text content.
It preserves line breaks
and whitespace.

  Indentation is preserved too.

Special characters: <tag> & "quotes" 'apostrophes'`)

const activeTab = ref('markdown')

const emojiTestContent = ref(`# 🎉 Emoji 完整测试

## 😊 表情符号

😀 😁 😂 😃 😄 😅 😆 😇 😈 😉 😊 😍 😘 😗 😙 😚 😜 😝 😛 😳 😌 😔 😑 😒 😐 😓 😥 😬 😕 😟 😲 😞 😖 😢 😭 😱 😠 😡 😦 😧 😨 😬 😰 😢 😭 😱 😨 😠 😡 😤 😖 😆 😋 😷 😎 😴 😵 😲 😳 😦 😧 😨 😬 😰 😢 😭 😱

## 👋 手势符号

👋 👍 👎 👌 ✌️ ✋ 👏 🙌 👐 🤝 💪 ☝️ 👆 👇 👈 👉 🤟 ✌️ 🤘 🤞

## ❤️ 心形和爱

❤️ 🧡 💛 💚 💙 💜 🖤 🤍 🤎 💔 💕 💖 💗 💘 💝 💞 💓 💞 💝 💖 💗 💓 💕 💞

## 🎉 庆祝符号

🎉 🎊 🎈 🎁 🎀 🎂 🎃 🎄 🎆 🎇 ✨ 🌟 ⭐ 💫 🌠 ⚡ 💥 🔥

## 🚀 交通工具

🚀 🚁 🚂 🚃 🚄 🚅 🚆 🚇 🚈 🚉 🚊 🚝 🚞 🚋 🚌 🚎 🚐 🚑 🚒 🚓 🚔 🚕 🚖 🚗 🚘 🚙 🚚 🚛 🚜

## 🏠 物体符号

🏠 🏡 🏢 🏣 🏤 🏥 🏦 🏧 🏨 🏩 🏪 🏫 🏬 🏭 🏯 🏰 💒 🗼 🗽 🗿 ⛪ 🕌 🕍 ⛩️

## 🌍 自然符号

🌍 🌎 🌏 🌐 🌑 🌒 🌓 🌔 🌕 🌖 🌗 🌘 🌙 🌚 🌝 🌛 🌜 ⭐ 🌟 💫 ✨ ⚡ ☄️ 💥 🔥

## 🍎 食物符号

🍎 🍊 🍋 🍌 🍉 🍇 🍓 🍈 🍒 🍑 🍍 🥭 🍅 🍆 🥑 🥦 🥬 🥒 🌶️ 🌽 🥕 🍞 🥐 🥯

## ⚽ 运动符号

⚽ 🏀 🏈 ⚾ 🥎 🎾 🏐 🏉 🥏 🎳 🏓 🏸 🏒 🏑 🥍 🏏 ⛳ 🪃 🥅 ⛸️ 🎣 🎽 🎿 ⛷️

## 音乐和艺术 🎵

🎵 🎶 🎼 🎹 🎸 🎺 🎷 🎸 🎻 🎲 🎯 🎳 🎮 🎰

## 完成测试 ✅

所有 emoji 都应该在上面正确显示！😊✨🚀
`)

const basicUsageCode = `<script setup>
// import { ContentViewer } from '@/components/ContentViewer'

const content = \`# Hello World
This is **Markdown** content.\`
<\/script>

<template>
  <ContentViewer :content="content" type="markdown" />
</template>`

const htmlUsageCode = `<template>
  <ContentViewer
    :content="htmlContent"
    type="html"
  />
</template>

<script setup>
const htmlContent = '<h1>Title</h1><p>Content</p>'
<\/script>`

const blogUsageCode = `<template>
  <div class="blog-post">
    <h1>{{ post.title }}</h1>
    <ContentViewer
      :content="post.content"
      type="markdown"
    />
  </div>
</template>

<script setup>
// import { ContentViewer } from '@/components/ContentViewer'

const post = {
  title: 'My Blog Post',
  content: '# Hello\\n\\nThis is my post content.'
}
<\/script>`

const diagramContent = ref(`# Mermaid Diagrams Demo

## Flow Chart - User Authentication Flow

\`\`\`mermaid
graph TD
    A[User] --> B{Has Account?}
    B -->|Yes| C[Login]
    B -->|No| D[Register]
    C --> E{Valid Credentials?}
    D --> F{Email Verification}
    E -->|Yes| G[Dashboard]
    E -->|No| H[Error Message]
    F -->|Success| G
    F -->|Failed| D
    H --> C
    style A fill:#E1F5FE
    style G fill:#C8E6C9
    style H fill:#FFCDD2
\`\`\`

## Sequence Diagram - API Request Flow

\`\`\`mermaid
sequenceDiagram
    participant Browser
    participant Frontend as Frontend App
    participant API as API Server
    participant Database
    Browser->>Frontend: User Action
    Frontend->>API: HTTP Request
    API->>Database: Query Data
    Database-->>API: Return Data
    API-->>Frontend: JSON Response
    Frontend-->>Browser: Render UI
    note over Browser,Frontend: Request Complete
\`\`\`

## State Machine - Order Status

\`\`\`mermaid
stateDiagram-v2
    [*] --> Pending
    Pending --> Processing: Confirm
    Processing --> Shipping: Pack
    Shipping --> Delivered: Deliver
    Shipping --> Failed: Error
    Failed --> Pending: Retry
    Delivered --> [*]
    Processing --> Cancelled: Cancel
    Cancelled --> [*]
\`\`\`

## Class Diagram - Content Management System

\`\`\`mermaid
classDiagram
    class Content {
        id: string
        title: string
        body: string
        publish()
        archive()
    }
    class Post {
        category: string
        tags: array
        getRelated()
    }
    class Page {
        slug: string
        featured: boolean
        getMenu()
    }
    class User {
        name: string
        email: string
        authenticate()
    }
    Content <|-- Post
    Content <|-- Page
    User --> Content : creates
\`\`\`

## Gantt Chart - Development Timeline

\`\`\`mermaid
gantt
    title Content Hub Development
    dateFormat YYYY-MM-DD

    section Planning
    Requirements :req, 2024-01-01, 14d
    Design :design, 2024-01-08, 21d

    section Frontend
    UI Components :ui, 2024-01-22, 30d
    Pages Development :pages, 2024-02-15, 35d
    Integration :fe_int, 2024-03-15, 20d

    section Backend
    API Design :api_design, 2024-01-22, 14d
    Core Services :services, 2024-02-05, 45d
    Database :db, 2024-02-05, 30d
    Integration :be_int, 2024-03-15, 20d

    section Testing
    Unit Testing :unit, 2024-04-01, 15d
    Integration Testing :int_test, 2024-04-10, 15d
    UAT :uat, 2024-04-20, 20d

    section Deployment
    Staging :stage, 2024-05-10, 10d
    Production :prod, 2024-05-20, 5d
    crit Milestone :crit, 2024-05-25, 1d
\`\`\`

## Pie Chart - Website Traffic Sources

\`\`\`mermaid
pie title Website Traffic Distribution
    "Direct" : 35
    "Search Engines" : 30
    "Social Media" : 20
    "Referrals" : 10
    "Ads" : 5
\`\`\`

## Entity Relationship Diagram

\`\`\`mermaid
erDiagram
    USER ||--o{ POST : creates
    USER ||--o{ COMMENT : writes
    POST ||--o{ COMMENT : has
    POST ||--o{ TAG : contains
    USER {
        int id
        string email
        string name
    }
    POST {
        int id
        string title
        text content
        int user_id
    }
    COMMENT {
        int id
        text content
        int user_id
        int post_id
    }
    TAG {
        int id
        string name
    }
\`\`\``
)
</script>

<template>
  <div class="viewer-demo">
    <n-card title="ContentViewer Component Demo">
      <n-tabs v-model:value="activeTab" type="segment">
        <n-tab-pane name="markdown" tab="Markdown">
          <n-card :bordered="false" class="content-card">
            <ContentViewer :content="markdownContent" type="markdown" />
          </n-card>
        </n-tab-pane>

        <n-tab-pane name="emoji" tab="Emoji Test 😊">
          <n-card :bordered="false" class="content-card">
            <ContentViewer :content="emojiTestContent" type="markdown" />
          </n-card>
        </n-tab-pane>

        <n-tab-pane name="diagrams" tab="Diagrams 📊">
          <n-card :bordered="false" class="content-card">
            <ContentViewer :content="diagramContent" type="markdown" />
          </n-card>
        </n-tab-pane>

        <n-tab-pane name="html" tab="HTML">
          <n-card :bordered="false" class="content-card">
            <ContentViewer :content="htmlContent" type="html" />
          </n-card>
        </n-tab-pane>

        <n-tab-pane name="text" tab="Plain Text">
          <n-card :bordered="false" class="content-card">
            <ContentViewer :content="textContent" type="text" />
          </n-card>
        </n-tab-pane>
      </n-tabs>
    </n-card>

    <!-- Usage Examples -->
    <n-card title="Usage Examples" style="margin-top: 2rem">
      <n-space vertical :size="24">
        <div>
          <h3>Basic Usage</h3>
          <n-code language="vue" :code="basicUsageCode" />
        </div>

        <div>
          <h3>With HTML Content</h3>
          <n-code language="vue" :code="htmlUsageCode" />
        </div>

        <div>
          <h3>In a Blog Post</h3>
          <n-code language="vue" :code="blogUsageCode" />
        </div>
      </n-space>
    </n-card>
  </div>
</template>


<style scoped lang="less">
.viewer-demo {
  padding: 2rem;
  max-width: 1200px;
  margin: 0 auto;

  .content-card {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    margin-top: 1rem;
  }

  h3 {
    margin-bottom: 1rem;
    color: var(--color-text-primary);
  }
}
</style>

