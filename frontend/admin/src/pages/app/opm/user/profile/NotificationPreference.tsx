import { useCallback, useState } from 'react';
import { Button, Checkbox, DatePicker, Spin, Switch, Typography, App } from 'antd';
import { useTranslation } from 'react-i18next';
import { useQuery } from '@tanstack/react-query';
import { apiClient } from '@/api/client';
import dayjs, { type Dayjs } from 'dayjs';

const { Text, Title } = Typography;

/**
 * 个人中心「通知偏好」：实时推送静音时段 + 站内信分类退订（仅全员广播）。
 *
 * 后端语义（通知域）：
 *  - 静音时段只抑制 SSE 实时推送（角标即时性），收件行照常落库，
 *    用户打开收件中心仍能看到全部消息；
 *  - 分类退订只约束全员广播，点对点定向发送不受影响。
 *  - 全链 fail-open：偏好读不到 = 不静音。
 * 验证码等事务性出站邮件不经偏好层，不受本页任何配置影响。
 */

const minutesToDayjs = (minutes: number) => dayjs().startOf('day').add(minutes, 'minute');

const dayjsToMinutes = (value: Dayjs | undefined): number | undefined => {
  if (!value?.isValid?.()) return undefined;
  return value.hour() * 60 + value.minute();
};

// 从未知错误对象取本地化 message（request-client 拦截器已把 reason 译文写进 message）
const errMessage = (error: unknown, fallback: string): string => {
  if (error && typeof error === 'object' && 'message' in error) {
    const m = (error as { message?: unknown }).message;
    if (typeof m === 'string' && m) return m;
  }
  return fallback;
};

const NotificationPreference = () => {
  const prefQuery = useQuery({
    queryKey: ['myNotificationPreference'],
    queryFn: () => apiClient.notificationPreferenceService.GetMyNotificationPreference({}),
  });
  const categoriesQuery = useQuery({
    queryKey: ['myNotifiableCategories'],
    queryFn: () => apiClient.notificationPreferenceService.ListMyNotificationCategories({}),
  });

  if (prefQuery.isPending) {
    return <Spin style={{ display: 'block', margin: '48px auto' }} />;
  }

  // 数据就绪后才挂表单：initial values 一次给定，改动都在表单内部状态，
  // 保存成功靠 key 重挂以服务端返回为准（避免 effect 里同步 setState 的级联渲染）。
  const pref = prefQuery.data;
  return (
    <QuietHoursForm
      key={pref?.updatedAt ?? 'default'}
      initialEnabled={!!pref?.quietEnabled}
      initialRange={[
        minutesToDayjs(pref?.quietStartMinute ?? 1320),
        minutesToDayjs(pref?.quietEndMinute ?? 480),
      ]}
      initialMutedCategoryIds={pref?.mutedCategoryIds ?? []}
      categoryOptions={(categoriesQuery.data?.items ?? []).map((c) => ({
        label: c.name ?? '',
        value: c.id ?? 0,
      }))}
      onSaved={() => prefQuery.refetch()}
    />
  );
};

type QuietHoursFormProps = {
  initialEnabled: boolean;
  initialRange: [Dayjs | undefined, Dayjs | undefined];
  initialMutedCategoryIds: number[];
  categoryOptions: { label: string; value: number }[];
  onSaved: () => void;
};

const QuietHoursForm = ({
  initialEnabled,
  initialRange,
  initialMutedCategoryIds,
  categoryOptions,
  onSaved,
}: QuietHoursFormProps) => {
  const { t } = useTranslation('profile');
  const { message } = App.useApp();
  const [saving, setSaving] = useState(false);
  const [quietEnabled, setQuietEnabled] = useState<boolean>(initialEnabled);
  const [quietRange, setQuietRange] = useState<[Dayjs | undefined, Dayjs | undefined]>(initialRange);
  const [mutedCategoryIds, setMutedCategoryIds] = useState<number[]>(initialMutedCategoryIds);

  const handleSave = useCallback(async () => {
    const startMinute = dayjsToMinutes(quietRange[0]);
    const endMinute = dayjsToMinutes(quietRange[1]);
    if (quietEnabled && (startMinute === undefined || endMinute === undefined)) {
      message.warning(t('notification.quietTimeRequired'));
      return;
    }
    if (quietEnabled && startMinute === endMinute) {
      message.warning(t('notification.quietZeroWidth'));
      return;
    }

    setSaving(true);
    try {
      const saved = await apiClient.notificationPreferenceService.UpdateMyNotificationPreference({
        quietEnabled,
        quietStartMinute: startMinute,
        quietEndMinute: endMinute,
        mutedCategoryIds,
      });
      message.success(t('notification.saveSuccess'));
      // 以服务端返回为准刷新（重挂表单），而不是本地假设
      if (saved) {
        onSaved();
      }
    } catch (error) {
      console.error('save notification preference failed', error);
      message.error(errMessage(error, t('notification.saveFailed')));
    } finally {
      setSaving(false);
    }
  }, [quietEnabled, quietRange, mutedCategoryIds, message, t, onSaved]);

  return (
    <div style={{ maxWidth: 560 }}>
      <Title level={5} style={{ marginBottom: 4 }}>
        {t('notification.quietTitle')}
      </Title>
      <Text type="secondary">{t('notification.quietDesc')}</Text>
      <div style={{ marginTop: 12, marginBottom: 16 }}>
        <Switch
          checked={quietEnabled}
          onChange={(checked) => setQuietEnabled(checked)}
          checkedChildren={t('notification.on')}
          unCheckedChildren={t('notification.off')}
        />
        <span style={{ marginLeft: 8 }}>{t('notification.quietEnable')}</span>
      </div>
      {quietEnabled && (
        <div style={{ marginBottom: 16 }}>
          <Text type="secondary" style={{ display: 'block', marginBottom: 8 }}>
            {t('notification.quietTimeRange')}
          </Text>
          {/* antd6 的 TimePicker.RangePicker 类型面已移除，时间范围走 DatePicker.RangePicker + picker="time" */}
          <DatePicker.RangePicker
            picker="time"
            value={quietRange}
            onChange={(values) => setQuietRange([values?.[0] ?? undefined, values?.[1] ?? undefined])}
            format="HH:mm"
            minuteStep={5}
            allowClear={false}
            style={{ width: 260 }}
          />
        </div>
      )}

      <Title level={5} style={{ marginBottom: 4, marginTop: 24 }}>
        {t('notification.muteTitle')}
      </Title>
      <Text type="secondary">{t('notification.muteDesc')}</Text>
      <div style={{ marginTop: 12, marginBottom: 16 }}>
        {categoryOptions.length > 0 ? (
          <Checkbox.Group
            options={categoryOptions}
            value={mutedCategoryIds}
            onChange={(values) => setMutedCategoryIds(values as number[])}
            style={{ display: 'flex', flexWrap: 'wrap', gap: '8px 16px' }}
          />
        ) : (
          <Text type="secondary">{t('notification.noCategories')}</Text>
        )}
      </div>

      <Button type="primary" onClick={handleSave} loading={saving}>
        {t('notification.save')}
      </Button>

      <Text type="secondary" style={{ display: 'block', marginTop: 16, fontSize: 12 }}>
        {t('notification.transactionalNote')}
      </Text>
    </div>
  );
};

export default NotificationPreference;
