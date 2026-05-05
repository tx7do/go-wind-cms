export type ContentType = 'markdown' | 'html' | 'text'

export interface Props {
  content?: string
  type?: ContentType
  class?: string
}
