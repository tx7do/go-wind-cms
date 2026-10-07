import { useEffect, useState } from 'react';
import { Button, Input, Popover, App } from 'antd';
import { ThunderboltOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useGenerateContent } from '@/api/hooks/ai-content';

interface AiGenerateButtonProps {
  /** 场景：DESCRIPTION / ANNOUNCEMENT / REPLY / GENERAL / POST_TITLE / POST_SUMMARY / POST_BODY */
  scene: string;
  /** 生成结果回调（前端将文本填入目标字段） */
  onGenerate: (content: string) => void;
  /** 补充上下文（可选） */
  context?: string;
  /** 主题预填值（如用文章标题当摘要生成的主题；用户仍可改） */
  defaultTopic?: string;
  /** 按钮文字，默认"AI 生成" */
  label?: string;
  /** 按钮大小 */
  size?: 'small' | 'middle' | 'large';
  /** 是否禁用 */
  disabled?: boolean;
}

/**
 * AI 内容生成按钮：点击弹出主题输入，AI 按场景生成文本后通过 onGenerate 回调。
 * 放在 textarea / input 旁边即可使用；通用组件，不绑定具体表单。
 */
export default function AiGenerateButton({
  scene,
  onGenerate,
  context,
  defaultTopic,
  label,
  size = 'small',
  disabled = false,
}: AiGenerateButtonProps) {
  const { t, i18n } = useTranslation('aiGenerate');
  const { message } = App.useApp();
  const [topic, setTopic] = useState(defaultTopic ?? '');
  const [open, setOpen] = useState(false);

  // defaultTopic 随源字段（如标题）变化时，若用户未手动输入过则跟随
  useEffect(() => {
    if (defaultTopic) setTopic(defaultTopic);
  }, [defaultTopic]);

  const generateMutation = useGenerateContent();

  const handleGenerate = async () => {
    const trimmed = topic.trim();
    if (!trimmed) return;
    try {
      const resp = await generateMutation.mutateAsync({
        scene,
        topic: trimmed,
        context,
        lang: i18n.language,
      });
      onGenerate(resp.content ?? '');
      setOpen(false);
      setTopic('');
    } catch (error) {
      console.error('ai content generate failed:', error);
      message.error(error instanceof Error ? error.message : t('failed'));
    }
  };

  return (
    <Popover
      open={open}
      onOpenChange={setOpen}
      trigger="click"
      placement="bottomRight"
      content={
        <div className="flex w-72 flex-col gap-2">
          <Input
            value={topic}
            onChange={(e) => setTopic(e.target.value)}
            onPressEnter={() => handleGenerate()}
            placeholder={t('topicPlaceholder')}
            size="small"
            autoFocus
            disabled={generateMutation.isPending}
          />
          <Button
            type="primary"
            size="small"
            icon={<ThunderboltOutlined />}
            loading={generateMutation.isPending}
            disabled={!topic.trim()}
            onClick={() => handleGenerate()}
          >
            {t('generateNow')}
          </Button>
        </div>
      }
    >
      <Button size={size} disabled={disabled} icon={<ThunderboltOutlined />}>
        {label ?? t('button')}
      </Button>
    </Popover>
  );
}
