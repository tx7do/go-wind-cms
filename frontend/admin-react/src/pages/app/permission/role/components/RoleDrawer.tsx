import { useRef, useState, useEffect } from 'react';
import type { ProFormInstance } from '@ant-design/pro-components';
import {
  DrawerForm,
  ProFormText,
  ProFormDigit,
  ProFormTextArea,
  ProFormRadio,
  ProFormSelect,
  ProFormDependency,
} from '@ant-design/pro-components';
import { App, Tree, Spin, Checkbox } from 'antd';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import type { permissionservicev1_Role as Role } from '@/api/generated/admin/service/v1';
import { useCreateRole, useUpdateRole } from '@/api/hooks/role';
import { fetchListPermissionGroups } from '@/api/hooks/permission-group';
import { fetchListPermissions } from '@/api/hooks/permission';
import { fetchListOrgUnits } from '@/api/hooks/org-unit';
import { PaginationQuery } from '@/core';
import { getStatusOptions, getDataScopeOptions, buildPermissionTree, extractLeafIds, buildOrgUnitTree, getUserFieldPermissionOptions } from '../constants';

interface RoleDrawerProps {
  open: boolean;
  mode: 'create' | 'edit';
  data?: Role;
  onClose: () => void;
  onSuccess: () => void;
}

/**
 * 角色编辑/创建抽屉组件
 */
const RoleDrawer: React.FC<RoleDrawerProps> = ({ open, mode, data, onClose, onSuccess }) => {
  const { t } = useTranslation('role');
  const formRef = useRef<ProFormInstance>(null);
  const queryClient = useQueryClient();
  const { message } = App.useApp();

  const [confirmLoading, setConfirmLoading] = useState(false);
  const [treeData, setTreeData] = useState<any[]>([]);
  const [checkedKeys, setCheckedKeys] = useState<number[]>([]);
  const [treeVersion, setTreeVersion] = useState(0);
  const [treeLoading, setTreeLoading] = useState(false);
  const [unitTreeData, setUnitTreeData] = useState<any[]>([]);
  const [unitCheckedKeys, setUnitCheckedKeys] = useState<number[]>([]);
  const [unitTreeVersion, setUnitTreeVersion] = useState(0);
  // 字段权限：User 资源上勾选隐藏的字段（json_name）
  const [userHiddenFields, setUserHiddenFields] = useState<string[]>([]);

  // 加载权限树数据
  useEffect(() => {
    if (open) {
      setTreeLoading(true);
      Promise.all([
        fetchListPermissionGroups(new PaginationQuery({ formValues: { status: 'ON' } })),
        fetchListPermissions(new PaginationQuery({ formValues: { status: 'ON' } })),
      ])
        .then(([groupRes, permRes]) => {
          const groups = groupRes?.items || [];
          const permissions = permRes?.items || [];
          if (groups.length === 0 && permissions.length > 0) {
            // 无权限组时，直接罗列权限（平铺，不嵌套）
            setTreeData(permissions.map((p) => ({
              key: Number(p.id),
              title: p.name || p.code || String(p.id),
            })));
          } else {
            setTreeData(buildPermissionTree(groups, permissions));
          }
          setTreeVersion((v) => v + 1);
        })
        .catch(() => setTreeData([]))
        .finally(() => setTreeLoading(false));
    }
  }, [open]);

  // 加载组织单元树（SELECTED_UNITS 自定义授权集配置用）
  useEffect(() => {
    if (open) {
      fetchListOrgUnits(new PaginationQuery({ formValues: { status: 'ON' } }))
        .then((res) => {
          setUnitTreeData(buildOrgUnitTree((res?.items || []) as any[]));
          setUnitTreeVersion((v) => v + 1);
        })
        .catch(() => setUnitTreeData([]));
    }
  }, [open]);

  // 编辑模式填充表单
  useEffect(() => {
    if (open && mode === 'edit' && data) {
      setTimeout(() => {
        formRef.current?.setFieldsValue({
          name: data.name || '',
          code: data.code || '',
          sortOrder: (data as any).sortOrder ?? 1,
          status: data.status || 'ON',
          dataScope: (data as any).dataScope || 'ALL',
          description: (data as any).description || '',
        });
      }, 0);
      // 设置已勾选的权限
      const perms = (data as any).permissions;
      if (Array.isArray(perms)) {
        setCheckedKeys(perms.filter((v: any) => typeof v === 'number'));
      }
      // SELECTED_UNITS 时回填已配置的授权单元集
      if ((data as any).dataScope === 'SELECTED_UNITS') {
        const units = (data as any).orgUnits;
        if (Array.isArray(units)) {
          setUnitCheckedKeys(units.filter((v: any) => typeof v === 'number'));
        }
      }
      // 回填字段权限（User 资源的隐藏字段集）
      const fpEntry = Array.isArray((data as any).fieldPermissions)
        ? (data as any).fieldPermissions.find((e: any) => e?.resource === 'User')
        : undefined;
      setUserHiddenFields(
        Array.isArray(fpEntry?.hiddenFields)
          ? fpEntry.hiddenFields.filter((v: any) => typeof v === 'string')
          : [],
      );
    }
  }, [open, mode, data]);

  const createMutation = useCreateRole({
    onSuccess: () => {
      message.success(t('createSuccess'));
      onSuccess();
      onClose();
      queryClient.invalidateQueries({ queryKey: ['listRoles'] });
    },
    onError: (error: Error) => {
      message.error(error.message || t('createFailed'));
    },
  });

  const updateMutation = useUpdateRole({
    onSuccess: () => {
      message.success(t('updateSuccess'));
      onSuccess();
      onClose();
      queryClient.invalidateQueries({ queryKey: ['listRoles'] });
    },
    onError: (error: Error) => {
      message.error(error.message || t('updateFailed'));
    },
  });

  const handleSubmit = async (values: any) => {
    setConfirmLoading(true);
    try {
      const payload = {
        ...values,
        // 用叶子交集剥离权限组父节点 ID（父子联动会自动勾选父节点，
        // 组 ID 混入提交会被后端当作权限 ID，可能越权绑定）。
        permissions: extractLeafIds(checkedKeys, treeData),
      };

      // 仅 SELECTED_UNITS 档提交授权单元集（含清空场景）；
      // 其余档位不携带该字段，后端维持既有集不替换。
      if (values.dataScope === 'SELECTED_UNITS') {
        payload.orgUnits = unitCheckedKeys.filter((v) => typeof v === 'number');
      }

      // 字段权限始终随表单提交（含清空场景）：抽屉所见即保存后的最终态。
      payload.fieldPermissions = userHiddenFields.length > 0
        ? [{ resource: 'User', hiddenFields: userHiddenFields }]
        : [];

      if (mode === 'create') {
        await createMutation.mutateAsync({ data: payload });
      } else if (data?.id) {
        await updateMutation.mutateAsync({ id: data.id, values: payload });
      }
    } finally {
      setConfirmLoading(false);
    }
  };

  return (
    <DrawerForm
      formRef={formRef}
      title={mode === 'create' ? t('create') : t('edit')}
      open={open}
      onOpenChange={(visible) => {
        if (!visible) {
          formRef.current?.resetFields();
          setCheckedKeys([]);
          setUnitCheckedKeys([]);
          setUserHiddenFields([]);
          onClose();
        }
      }}
      initialValues={{
        sortOrder: 1,
        status: 'ON',
        dataScope: 'ALL',
      }}
      onFinish={handleSubmit}
      submitter={{
        searchConfig: {
          submitText: t('common:button.submit'),
          resetText: t('common:button.cancel'),
        },
        submitButtonProps: {
          loading: confirmLoading || createMutation.isPending || updateMutation.isPending,
        },
        resetButtonProps: { onClick: onClose },
      }}
      drawerProps={{ destroyOnHidden: true, onClose, size: 600 }}
    >
      <ProFormText
        name="name"
        label={t('name')}
        placeholder={t('namePlaceholder')}
        rules={[{ required: true, message: t('requiredName') }]}
        fieldProps={{ allowClear: true }}
      />

      <ProFormText
        name="code"
        label={t('code')}
        placeholder={t('codePlaceholder')}
        rules={[{ required: true, message: t('requiredCode') }]}
        fieldProps={{ allowClear: true }}
      />

      <ProFormDigit
        name="sortOrder"
        label={t('sortOrder')}
        placeholder={t('sortOrderPlaceholder')}
        rules={[{ required: true }]}
        fieldProps={{ precision: 0, min: 0 }}
      />

      <ProFormRadio.Group
        name="status"
        label={t('status')}
        rules={[{ required: true, message: t('requiredStatus') }]}
        options={getStatusOptions(t)}
        fieldProps={{ optionType: 'button', buttonStyle: 'solid' }}
      />

      <ProFormSelect
        name="dataScope"
        label={t('dataScope')}
        placeholder={t('dataScopePlaceholder')}
        options={getDataScopeOptions(t)}
        rules={[{ required: true, message: t('requiredDataScope') }]}
        fieldProps={{ allowClear: false }}
      />

      <ProFormDependency name={['dataScope']}>
        {({ dataScope }) =>
          dataScope === 'SELECTED_UNITS' ? (
            <div className="mb-6">
              <label className="block mb-2 text-sm font-medium text-[color:var(--ant-color-text)]">
                {t('orgUnits')}
              </label>
              {unitTreeData.length > 0 ? (
                <div className="rounded-lg border border-white/10 bg-white/5 p-2.5 dark:border-white/8 dark:bg-zinc-800/40">
                  <Tree
                    checkable
                    checkedKeys={unitCheckedKeys}
                    onCheck={(checked) => {
                      setUnitCheckedKeys(checked as number[]);
                    }}
                    treeData={unitTreeData}
                    defaultExpandAll
                    key={unitTreeVersion}
                    className="max-h-80 overflow-auto"
                  />
                </div>
              ) : (
                <div className="rounded-lg border border-dashed border-[color:var(--ant-color-border)] py-4 text-center text-sm text-[color:var(--ant-color-text-quaternary)]">
                  {t('noOrgUnitData')}
                </div>
              )}
            </div>
          ) : null
        }
      </ProFormDependency>

      <ProFormTextArea
        name="description"
        label={t('description')}
        placeholder={t('descriptionPlaceholder')}
        fieldProps={{ allowClear: true, rows: 2 }}
      />

      {/* 权限树 */}
      <div className="mb-6">
        <label className="block mb-2 text-sm font-medium text-[color:var(--ant-color-text)]">
          {t('permissions')}
        </label>
        <Spin spinning={treeLoading}>
          {treeData.length > 0 ? (
            <div className="rounded-lg border border-white/10 bg-white/5 p-2.5 dark:border-white/8 dark:bg-zinc-800/40">
              <Tree
                checkable
                checkedKeys={checkedKeys}
                onCheck={(checked) => {
                  setCheckedKeys(checked as number[]);
                }}
                treeData={treeData}
                defaultExpandAll
                key={treeVersion}
                className="max-h-80 overflow-auto"
              />
            </div>
          ) : (
            <div className="rounded-lg border border-dashed border-[color:var(--ant-color-border)] py-4 text-center text-sm text-[color:var(--ant-color-text-quaternary)]">
              {treeLoading ? '' : 'No permission data'}
            </div>
          )}
        </Spin>
      </div>

      {/* 字段权限：勾选 = 对该角色用户隐藏（黑名单语义，User 资源试点） */}
      <div className="mb-6">
        <label className="block mb-2 text-sm font-medium text-[color:var(--ant-color-text)]">
          {t('fieldPerm.title')}
        </label>
        <div className="rounded-lg border border-white/10 bg-white/5 p-2.5 dark:border-white/8 dark:bg-zinc-800/40">
          <p className="mb-2 text-xs text-[color:var(--ant-color-text-quaternary)]">
            {t('fieldPerm.hint')}
          </p>
          <Checkbox.Group
            value={userHiddenFields}
            onChange={(vals) => setUserHiddenFields(vals as string[])}
            options={getUserFieldPermissionOptions(t)}
            className="flex flex-col gap-2"
          />
        </div>
      </div>
    </DrawerForm>
  );
};

export default RoleDrawer;
