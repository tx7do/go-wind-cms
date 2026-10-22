import { useState } from 'react';
import { Form, Input, Select, Button, Avatar, Row, Col, Descriptions, Tag, App } from 'antd';
import { useTranslation } from 'react-i18next';
import { useQueryClient } from '@tanstack/react-query';
import type { identityservicev1_User as User } from '@/api/generated/admin/service/v1';
import {
  useGetUserProfile,
  useUpdateUserProfile,
} from '@/api/hooks/user-profile';
import { useAuthStore, useUserStore } from '@/stores';
import { getCharColor, getPresetColor } from '@/utils/color';
import { formatDateTime } from '@/utils/date';
import { getGenderOptions } from '../constants';

const BaseSetting = () => {
  const { t } = useTranslation('profile');
  const { message } = App.useApp();
  const queryClient = useQueryClient();
  const [form] = Form.useForm();
  const [saving, setSaving] = useState(false);

  const { data: userData, isLoading } = useGetUserProfile();
  const user = userData as User | undefined | null;

  const firstChar = user?.username?.slice(0, 1).toUpperCase() || '?';
  const avatarColor = getCharColor(firstChar);

  const updateMutation = useUpdateUserProfile({
    onSuccess: () => {
      message.success(t('saveSuccess'));
      queryClient.invalidateQueries({ queryKey: ['getMe'] });
      // 同步用户 store：顶栏/问候语读 useUserStore 快照（auth store 的 userInfo 仅登录期存在），不刷会残留旧昵称
      queryClient.fetchQuery({ queryKey: ['getMe'] }).then((info) => {
        if (!info) return;
        useUserStore.getState().setUserInfo(info as unknown as BasicUserInfo);
        useAuthStore.getState().setUserInfo(info as unknown as UserInfo);
      }).catch((err) => {
        console.error('Failed to refresh user info after profile save:', err);
      });
    },
    onError: (err: Error) => {
      console.error('Failed to update profile:', err);
      message.error(err.message || t('saveFailed'));
    },
  });

  const handleSave = async () => {
    try {
      const values = await form.validateFields();
      if (!user?.id) return;
      setSaving(true);
      await updateMutation.mutateAsync({ id: user.id, values });
    } catch {
      // 表单校验失败，Antd 已就地标红字段
    } finally {
      setSaving(false);
    }
  };

  if (isLoading) return null;

  return (
    <div style={{ overflowX: 'hidden' }}>
      <Row gutter={24}>
        <Col span={14}>
          <Form
            form={form}
            layout="horizontal"
            labelCol={{ flex: '0 0 120px' }}
            wrapperCol={{ flex: 'auto' }}
            initialValues={{
              nickname: user?.nickname || '',
              realname: user?.realname || '',
              email: user?.email || '',
              mobile: user?.mobile || '',
              telephone: user?.telephone || '',
              gender: user?.gender,
              region: user?.region || '',
              address: user?.address || '',
              description: user?.description || '',
              locale: user?.locale || '',
            }}
          >
            <Form.Item name="nickname" label={t('field.nickname')} rules={[{ required: true }]}>
              <Input placeholder={t('field.nickname')} />
            </Form.Item>
            <Form.Item name="realname" label={t('field.realname')}>
              <Input placeholder={t('field.realname')} />
            </Form.Item>
            <Form.Item name="email" label={t('field.email')}>
              <Input placeholder={t('field.email')} />
            </Form.Item>
            <Form.Item name="mobile" label={t('field.mobile')}>
              <Input placeholder={t('field.mobile')} />
            </Form.Item>
            <Form.Item name="telephone" label={t('field.telephone')}>
              <Input placeholder={t('field.telephone')} />
            </Form.Item>
            <Form.Item name="gender" label={t('field.gender')}>
              <Select allowClear options={getGenderOptions(t)} />
            </Form.Item>
            <Form.Item name="region" label={t('field.region')}>
              <Input placeholder={t('field.region')} />
            </Form.Item>
            <Form.Item name="address" label={t('field.address')}>
              <Input placeholder={t('field.address')} />
            </Form.Item>
            <Form.Item name="description" label={t('field.description')}>
              <Input.TextArea rows={3} placeholder={t('field.description')} />
            </Form.Item>
            <Form.Item name="locale" label={t('field.locale')} tooltip={t('localeHint')}>
              <Select
                allowClear
                options={[
                  { value: 'zh-CN', label: '简体中文' },
                  { value: 'en-US', label: 'English' },
                ]}
              />
            </Form.Item>
            <Form.Item>
              <Button type="primary" loading={saving} onClick={handleSave}>
                {t('save')}
              </Button>
            </Form.Item>
          </Form>
        </Col>

        <Col span={10}>
          <div style={{ paddingTop: 20 }}>
            <div style={{ fontSize: 14, fontWeight: 500, marginBottom: 16 }}>{t('field.avatar')}</div>
            {/* CMS 生成层无 UploadAvatar/DeleteAvatar RPC，头像仅展示（上传走用户管理/媒体库） */}
            <Avatar
              src={user?.avatar || undefined}
              size={96}
              style={!user?.avatar ? { backgroundColor: avatarColor, fontSize: 40, fontWeight: 700 } : {}}
            >
              {firstChar}
            </Avatar>
          </div>

          {/* 账号归属信息（只读） */}
          <Descriptions column={1} size="small" style={{ marginTop: 24 }}>
            <Descriptions.Item label={t('info.username')}>
              {user?.username || '-'}
            </Descriptions.Item>
            <Descriptions.Item label={t('info.roleNames')}>
              {user?.roleNames?.length ? (
                user.roleNames.map((role) => (
                  <Tag key={role} color={getPresetColor(role)} style={{ marginRight: 4 }}>
                    {role}
                  </Tag>
                ))
              ) : '-'}
            </Descriptions.Item>
            <Descriptions.Item label={t('info.orgUnitNames')}>
              {user?.orgUnitNames?.length ? (
                user.orgUnitNames.map((org) => (
                  <Tag key={org} color={getPresetColor(org)} style={{ marginRight: 4 }}>
                    {org}
                  </Tag>
                ))
              ) : '-'}
            </Descriptions.Item>
            <Descriptions.Item label={t('info.positionNames')}>
              {user?.positionNames?.length ? (
                user.positionNames.map((pos) => (
                  <Tag key={pos} color={getPresetColor(pos)} style={{ marginRight: 4 }}>
                    {pos}
                  </Tag>
                ))
              ) : '-'}
            </Descriptions.Item>
            <Descriptions.Item label={t('info.tenantName')}>
              {user?.tenantName || '-'}
            </Descriptions.Item>
            <Descriptions.Item label={t('info.lastLoginAt')}>
              {formatDateTime(user?.lastLoginAt)}
            </Descriptions.Item>
            <Descriptions.Item label={t('info.lastLoginIp')}>
              {user?.lastLoginIp || '-'}
            </Descriptions.Item>
            <Descriptions.Item label={t('info.createdAt')}>
              {formatDateTime(user?.createdAt)}
            </Descriptions.Item>
          </Descriptions>
        </Col>
      </Row>
    </div>
  );
};

export default BaseSetting;
