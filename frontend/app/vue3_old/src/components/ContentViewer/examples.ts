/**
 * ContentViewer Component - Real-world Usage Examples
 *
 * Examples for using ContentViewer in a Content Hub frontend application
 */

// ============================================================
// Example 1: Blog Post Detail Page
// ============================================================
/*
<template>
  <div class="post-detail">
    <article class="post-article">
      <header class="post-header">
        <h1 class="post-title">{{ post.title }}</h1>
        <div class="post-meta">
          <n-avatar :src="post.author.avatar" round size="small" />
          <span class="author-name">{{ post.author.name }}</span>
          <span class="separator">·</span>
          <span class="post-date">{{ formatDate(post.createdAt) }}</span>
          <span class="separator">·</span>
          <span class="read-time">{{ post.readTime }} min read</span>
        </div>
        <n-space class="post-tags">
          <n-tag v-for="tag in post.tags" :key="tag.id" :bordered="false">
            {{ tag.name }}
          </n-tag>
        </n-space>
      </header>

      <div class="post-content">
        <ContentViewer
          :content="post.content"
          :type="post.contentType || 'markdown'"
        />
      </div>

      <footer class="post-footer">
        <n-space justify="space-between">
          <n-space>
            <n-button circle>
              <template #icon>
                <div class="i-carbon-thumbs-up" />
              </template>
            </n-button>
            <span>{{ post.likes }} likes</span>
          </n-space>
          <n-button type="primary" @click="sharePost">
            <template #icon>
              <div class="i-carbon-share" />
            </template>
            Share
          </n-button>
        </n-space>
      </footer>
    </article>
  </div>
</template>

<script setup lang="ts">
import { ContentViewer } from '@/components/ContentViewer'
import type { Post } from '@/typings'

interface Props {
  post: Post
}

const props = defineProps<Props>()

function formatDate(date: string) {
  return new Date(date).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

function sharePost() {
  // Share logic
}
</script>

<style scoped lang="less">
.post-detail {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem;
}

.post-article {
  background: var(--color-surface);
  border-radius: 12px;
  padding: 2rem;
}

.post-header {
  border-bottom: 1px solid var(--color-border);
  padding-bottom: 2rem;
  margin-bottom: 2rem;
}

.post-title {
  font-size: 2.5rem;
  font-weight: 700;
  margin-bottom: 1rem;
  color: var(--color-text-primary);
}

.post-meta {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: var(--color-text-secondary);
  font-size: 0.875rem;
  margin-bottom: 1rem;
}

.post-content {
  margin: 2rem 0;
}

.post-footer {
  border-top: 1px solid var(--color-border);
  padding-top: 2rem;
  margin-top: 2rem;
}
</style>
*/

// ============================================================
// Example 2: Comment List with ContentViewer
// ============================================================
/*
<template>
  <div class="comments-section">
    <h2 class="comments-title">
      {{ $t('cms.comment.title') }} ({{ comments.length }})
    </h2>

    <div class="comments-list">
      <div
        v-for="comment in comments"
        :key="comment.id"
        class="comment-item"
      >
        <div class="comment-header">
          <n-avatar :src="comment.author.avatar" round size="medium" />
          <div class="comment-author-info">
            <strong class="author-name">{{ comment.author.name }}</strong>
            <span class="comment-time">{{ timeAgo(comment.createdAt) }}</span>
          </div>
        </div>

        <div class="comment-content">
          <ContentViewer
            :content="comment.content"
            type="markdown"
          />
        </div>

        <div class="comment-actions">
          <n-button text @click="likeComment(comment.id)">
            <template #icon>
              <div class="i-carbon-thumbs-up" />
            </template>
            {{ comment.likes }}
          </n-button>
          <n-button text @click="replyComment(comment.id)">
            <template #icon>
              <div class="i-carbon-reply" />
            </template>
            Reply
          </n-button>
        </div>

        <!-- Nested replies -->
        <div v-if="comment.replies?.length" class="comment-replies">
          <div
            v-for="reply in comment.replies"
            :key="reply.id"
            class="comment-item comment-reply"
          >
            <!-- Similar structure for replies -->
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ContentViewer } from '@/components/ContentViewer'
import type { Comment } from '@/typings'

interface Props {
  comments: Comment[]
}

const props = defineProps<Props>()

function timeAgo(date: string) {
  // Time ago logic
  return '2 hours ago'
}

function likeComment(id: string) {
  // Like logic
}

function replyComment(id: string) {
  // Reply logic
}
</script>

<style scoped lang="less">
.comments-section {
  max-width: 800px;
  margin: 2rem auto;
  padding: 2rem;
}

.comments-title {
  font-size: 1.5rem;
  font-weight: 600;
  margin-bottom: 2rem;
  color: var(--color-text-primary);
}

.comment-item {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  padding: 1.5rem;
  margin-bottom: 1rem;
}

.comment-header {
  display: flex;
  gap: 1rem;
  margin-bottom: 1rem;
}

.comment-author-info {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.comment-time {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
}

.comment-content {
  margin: 1rem 0;
}

.comment-actions {
  display: flex;
  gap: 1rem;
}

.comment-replies {
  margin-top: 1rem;
  margin-left: 3rem;
}

.comment-reply {
  margin-bottom: 0.5rem;
}
</style>
*/

// ============================================================
// Example 3: Page Content Display
// ============================================================
/*
<template>
  <div class="page-content">
    <div class="page-hero" v-if="page.featuredImage">
      <img :src="page.featuredImage" :alt="page.title" />
    </div>

    <div class="page-container">
      <h1 class="page-title">{{ page.title }}</h1>

      <div class="page-body">
        <ContentViewer
          :content="page.content"
          :type="page.contentType || 'markdown'"
        />
      </div>

      <div class="page-meta" v-if="page.updatedAt">
        <n-text depth="3">
          Last updated: {{ formatDate(page.updatedAt) }}
        </n-text>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ContentViewer } from '@/components/ContentViewer'
import type { Page } from '@/typings'

interface Props {
  page: Page
}

const props = defineProps<Props>()

function formatDate(date: string) {
  return new Date(date).toLocaleDateString()
}
</script>

<style scoped lang="less">
.page-content {
  min-height: 100vh;
}

.page-hero {
  width: 100%;
  height: 400px;
  overflow: hidden;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

.page-container {
  max-width: 900px;
  margin: 0 auto;
  padding: 3rem 2rem;
}

.page-title {
  font-size: 3rem;
  font-weight: 700;
  margin-bottom: 2rem;
  color: var(--color-text-primary);
}

.page-body {
  margin: 2rem 0;
}

.page-meta {
  margin-top: 3rem;
  padding-top: 2rem;
  border-top: 1px solid var(--color-border);
}
</style>
*/

// ============================================================
// Example 4: Category/Tag Description
// ============================================================
/*
<template>
  <div class="category-page">
    <div class="category-header">
      <h1 class="category-title">{{ category.name }}</h1>

      <div v-if="category.description" class="category-description">
        <ContentViewer
          :content="category.description"
          type="markdown"
        />
      </div>

      <div class="category-stats">
        <n-statistic label="Posts" :value="category.postCount" />
        <n-statistic label="Followers" :value="category.followers" />
      </div>
    </div>

    <div class="category-posts">
      <!-- Post list -->
    </div>
  </div>
</template>

<script setup lang="ts">
import { ContentViewer } from '@/components/ContentViewer'
import type { Category } from '@/typings'

interface Props {
  category: Category
}

const props = defineProps<Props>()
</script>
*/

// ============================================================
// Example 5: Post Preview Card
// ============================================================
/*
<template>
  <n-card
    class="post-card"
    hoverable
    @click="navigateToPost(post.id)"
  >
    <div class="card-image" v-if="post.coverImage">
      <img :src="post.coverImage" :alt="post.title" />
    </div>

    <div class="card-content">
      <h3 class="card-title">{{ post.title }}</h3>

      <div class="card-excerpt">
        <ContentViewer
          :content="getExcerpt(post.content)"
          type="markdown"
        />
      </div>

      <div class="card-footer">
        <n-space>
          <n-avatar :src="post.author.avatar" size="small" round />
          <span class="author-name">{{ post.author.name }}</span>
        </n-space>
        <n-text depth="3">{{ formatDate(post.createdAt) }}</n-text>
      </div>
    </div>
  </n-card>
</template>

<script setup lang="ts">
import { ContentViewer } from '@/components/ContentViewer'
import type { Post } from '@/typings'

interface Props {
  post: Post
}

const props = defineProps<Props>()

function getExcerpt(content: string, maxLength: number = 150) {
  const text = content.replace(/[#*`>\-\[\]]/g, '').trim()
  return text.length > maxLength
    ? text.substring(0, maxLength) + '...'
    : text
}

function navigateToPost(id: string) {
  // Navigate to post detail
}

function formatDate(date: string) {
  return new Date(date).toLocaleDateString()
}
</script>

<style scoped lang="less">
.post-card {
  cursor: pointer;
  transition: all 0.3s ease;

  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
  }
}

.card-image {
  width: 100%;
  height: 200px;
  overflow: hidden;
  border-radius: 8px 8px 0 0;
  margin: -1rem -1rem 1rem;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

.card-title {
  font-size: 1.25rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
  color: var(--color-text-primary);
}

.card-excerpt {
  color: var(--color-text-secondary);
  margin-bottom: 1rem;

  :deep(p) {
    margin: 0;
    line-height: 1.6;
  }
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 1rem;
  border-top: 1px solid var(--color-border);
}
</style>
*/

// ============================================================
// Example 6: Loading State with ContentViewer
// ============================================================
/*
<template>
  <div class="content-container">
    <n-skeleton v-if="loading" text :repeat="5" />

    <ContentViewer
      v-else-if="content"
      :content="content"
      type="markdown"
    />

    <n-empty
      v-else
      :description="$t('component.viewer.empty')"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ContentViewer } from '@/components/ContentViewer'

const loading = ref(true)
const content = ref('')

onMounted(async () => {
  try {
    // Fetch content
    const response = await fetchContent()
    content.value = response.data
  } finally {
    loading.value = false
  }
})

async function fetchContent() {
  // API call
}
</script>
*/

export {}

// ============================================================
// Example 7: Math Formulas and Diagrams
// ============================================================
/*
<template>
  <div class="demo-container">
    <n-card title="ContentViewer - Math & Diagrams Demo" :bordered="false">
      <ContentViewer :content="demoContent" type="markdown" />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ContentViewer } from '@/components/ContentViewer'

const demoContent = \`
# Math Formulas & Diagrams Demo

## Inline Math Formulas

This is an inline math formula: $E = mc^2$

The quadratic formula is: $x = \\frac{-b \\pm \\sqrt{b^2 - 4ac}}{2a}$

## Block Math Formulas

$$
\\int_0^\\infty e^{-x^2} dx = \\frac{\\sqrt{\\pi}}{2}
$$

$$
\\sum_{i=1}^n i = \\frac{n(n+1)}{2}
$$

## Mermaid Flowchart

\\\`\\\`\\\`mermaid
graph TD
    A[Start] --> B{Decision}
    B -->|Yes| C[Process A]
    B -->|No| D[Process B]
    C --> E[End]
    D --> E
\\\`\\\`\\\`

## Mermaid Sequence Diagram

\\\`\\\`\\\`mermaid
sequenceDiagram
    participant User
    participant Browser
    participant Server

    User->>Browser: Click Login
    Browser->>Server: Send Credentials
    Server->>Server: Validate
    Server->>Browser: Return Token
    Browser->>Browser: Save Token
    Browser->>User: Show Dashboard
\\\`\\\`\\\`

## Mermaid State Diagram

\\\`\\\`\\\`mermaid
stateDiagram-v2
    [*] --> Idle
    Idle --> Processing: Start
    Processing --> Success: Complete
    Processing --> Error: Fail
    Success --> [*]
    Error --> Idle: Retry
\\\`\\\`\\\`

## Mermaid Class Diagram

\\\`\\\`\\\`mermaid
classDiagram
    class Post {
        id: number
        title: string
        content: string
        createdAt: Date
        publish()
        delete()
    }

    class Comment {
        id: number
        content: string
        author: string
        postId: number
        create()
    }

    Post "1" --> "*" Comment : has
\\\`\\\`\\\`

## Mathematical Expressions

### Fractions
$\\frac{a}{b}$ and $\\frac{x^2}{y^2}$

### Powers
$x^n$ and $(a+b)^2$

### Square Roots
$\\sqrt{x}$ and $\\sqrt[3]{8} = 2$

### Summation
$\\sum_{k=1}^{n} k = \\frac{n(n+1)}{2}$

### Integration
$\\int_{0}^{1} x^2 dx = \\frac{1}{3}$
\`
</script>

<style scoped lang="less">
.demo-container {
  max-width: 1000px;
  margin: 0 auto;
  padding: 2rem;
}
</style>
*/

// ============================================================
// Example 8: Advanced ContentViewer with All Features
// ============================================================
/*
<template>
  <div class="content-showcase">
    <n-tabs type="segment" :tabs="tabs" @update:value="activeTab = $event">
      <template #default>
        <n-tab-pane name="markdown">
          <ContentViewer :content="markdownContent" type="markdown" class="viewer-pane" />
        </n-tab-pane>

        <n-tab-pane name="html">
          <ContentViewer :content="htmlContent" type="html" class="viewer-pane" />
        </n-tab-pane>

        <n-tab-pane name="text">
          <ContentViewer :content="plainTextContent" type="text" class="viewer-pane" />
        </n-tab-pane>
      </template>
    </n-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ContentViewer } from '@/components/ContentViewer'
import type { TabsProps } from 'naive-ui'

const activeTab = ref('markdown')

const tabs: TabsProps['tabs'] = [
  { name: 'markdown', label: 'Markdown' },
  { name: 'html', label: 'HTML' },
  { name: 'text', label: 'Plain Text' }
]

const markdownContent = \`
# Complete Feature Demo

## Code Highlighting
\\\`\\\`\\\`typescript
interface User {
  id: number
  name: string
  email: string
  createdAt: Date
}

const users: User[] = []
\\\`\\\`\\\`

## Tables
| Feature | Support | Description |
|---------|:-------:|-------------|
| Code Highlight | ✅ | 200+ languages |
| Math Formulas | ✅ | KaTeX support |
| Diagrams | ✅ | Mermaid charts |
| Images | ✅ | With captions |

## Blockquotes
> This is a blockquote with **bold** and *italic* text.
>
> Supporting multiple lines.

## Lists
1. Numbered item 1
2. Numbered item 2
   - Sub item 1
   - Sub item 2
3. Numbered item 3
\`

const htmlContent = \`
<h1>HTML Content</h1>
<p>This is <strong>HTML</strong> content with <em>formatting</em>.</p>
<ul>
  <li>Item 1</li>
  <li>Item 2</li>
</ul>
\`

const plainTextContent = \`
This is plain text content.
It preserves whitespace and line breaks.

    Like this indented section.
\`
</script>

<style scoped lang="less">
.content-showcase {
  padding: 2rem;
}

.viewer-pane {
  padding: 2rem;
  background: var(--color-surface);
  border-radius: 8px;
}
</style>
*/

