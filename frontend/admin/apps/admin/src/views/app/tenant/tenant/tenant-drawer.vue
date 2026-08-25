<script lang="ts" setup>
import { computed, ref } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { notification } from 'ant-design-vue';

import { useVbenForm } from '#/adapter/form';
import {
  apiClient,
  makeUpdateMask,
  tenantAuditStatusList,
  tenantStatusList,
  tenantTypeList,
} from '#/api';

const data = ref();

const getTitle = computed(() =>
  data.value?.create
    ? $t('ui.modal.create', { moduleName: $t('page.tenant.moduleName') })
    : $t('ui.modal.update', { moduleName: $t('page.tenant.moduleName') }),
);
// const isCreate = computed(() => data.value?.create);

const [BaseForm, baseFormApi] = useVbenForm({
  showDefaultActions: false,
  // 所有表单项共用，可单独在表单内覆盖
  commonConfig: {
    // 所有表单项
    componentProps: {
      class: 'w-full',
    },
  },
  schema: [
    {
      component: 'Input',
      fieldName: 'name',
      label: $t('page.tenant.name'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
      rules: 'required',
    },
    {
      component: 'Input',
      fieldName: 'code',
      label: $t('page.tenant.code'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
      rules: 'required',
      dependencies: {
        show: (_values) => {
          return data.value?.create;
        },
        triggerFields: ['type'],
      },
    },
    {
      component: 'Select',
      fieldName: 'type',
      label: $t('page.tenant.type'),
      defaultValue: 'PAID',
      componentProps: {
        placeholder: $t('ui.placeholder.select'),
        options: tenantTypeList,
        filterOption: (input: string, option: any) =>
          option.label.toLowerCase().includes(input.toLowerCase()),
        allowClear: true,
        showSearch: true,
      },
      rules: 'selectRequired',
      dependencies: {
        show: (_values) => {
          return data.value?.create;
        },
        triggerFields: ['type'],
      },
    },
    {
      component: 'Select',
      fieldName: 'auditStatus',
      label: $t('page.tenant.auditStatus'),
      defaultValue: 'APPROVED',
      componentProps: {
        placeholder: $t('ui.placeholder.select'),
        options: tenantAuditStatusList,
        filterOption: (input: string, option: any) =>
          option.label.toLowerCase().includes(input.toLowerCase()),
        allowClear: true,
        showSearch: true,
      },
      rules: 'selectRequired',
      dependencies: {
        show: (_values) => {
          return data.value?.create;
        },
        triggerFields: ['type'],
      },
    },
    {
      component: 'Select',
      fieldName: 'status',
      defaultValue: 'ON',
      label: $t('ui.table.status'),
      rules: 'selectRequired',
      dependencies: {
        show: (_values) => {
          return data.value?.create;
        },
        triggerFields: ['type'],
      },
      componentProps: {
        options: tenantStatusList,
        filterOption: (input: string, option: any) =>
          option.label.toLowerCase().includes(input.toLowerCase()),
        allowClear: true,
        showSearch: true,
      },
    },
    {
      component: 'Textarea',
      fieldName: 'remark',
      label: $t('ui.table.remark'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
    },
    {
      component: 'Input',
      fieldName: 'logoUrl',
      label: $t('page.tenant.logoUrl'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
    },
    {
      component: 'Input',
      fieldName: 'industry',
      label: $t('page.tenant.industry'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
    },

    {
      component: 'Divider',
      fieldName: 'divider1',
      hideLabel: true,
      dependencies: {
        show: (_values) => {
          return data.value?.create;
        },
        triggerFields: ['type'],
      },
      renderComponentContent() {
        return {
          default: () => $t('page.tenant.adminSetting'),
        };
      },
    },

    {
      component: 'Input',
      fieldName: 'user.username',
      label: $t('page.tenant.adminUserName'),
      rules: 'required',
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
      dependencies: {
        show: (_values) => {
          return data.value?.create;
        },
        triggerFields: ['type'],
      },
    },

    {
      component: 'VbenInputPassword',
      fieldName: 'password',
      label: $t('page.tenant.adminPassword'),
      rules: 'required',
      componentProps: {
        passwordStrength: true,
        placeholder: $t('ui.placeholder.input'),
      },
      dependencies: {
        show: (_values) => {
          return data.value?.create;
        },
        triggerFields: ['type'],
      },
    },

    {
      component: 'VbenInputPassword',
      fieldName: 'passwordConfirm',
      label: $t('page.tenant.adminPasswordConfirm'),
      rules: 'required',
      componentProps: {
        passwordStrength: true,
        placeholder: $t('ui.placeholder.input'),
      },
      dependencies: {
        show: (_values) => {
          return data.value?.create;
        },
        triggerFields: ['type'],
      },
    },

    {
      component: 'Input',
      fieldName: 'user.mobile',
      label: $t('page.tenant.adminMobile'),
      rules: 'required',
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
      dependencies: {
        show: (_values) => {
          return data.value?.create;
        },
        triggerFields: ['type'],
      },
    },

    {
      component: 'Input',
      fieldName: 'user.email',
      label: $t('page.tenant.adminEmail'),
      rules: 'required',
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
      dependencies: {
        show: (_values) => {
          return data.value?.create;
        },
        triggerFields: ['type'],
      },
    },
  ],
});

const [Drawer, drawerApi] = useVbenDrawer({
  onCancel() {
    drawerApi.close();
  },

  async onConfirm() {
    // 校验输入的数据
    const validate = await baseFormApi.validate();
    if (!validate.valid) {
      return;
    }

    setLoading(true);

    // 获取表单数据
    const values = await baseFormApi.getValues();

    await (data.value?.create
      ? createTenantWithAdminUser(values)
      : updateTenant(values));
  },

  onOpenChange(isOpen: boolean) {
    if (isOpen) {
      // 获取传入的数据
      data.value = drawerApi.getData<Record<string, any>>();

      // 为表单赋值
      baseFormApi.setValues(data.value?.row);

      setLoading(false);

    }
  },
});

function setLoading(loading: boolean) {
  drawerApi.setState({ confirmLoading: loading });
}

// async function createTenant(values: any) {
//   console.log('createTenant', values);
//
//   try {
//     await tenantStore.createTenant(values);
//
//     notification.success({
//       message: $t('ui.notification.create_success'),
//     });
//   } catch {
//     notification.error({
//       message: $t('ui.notification.create_failed'),
//     });
//   } finally {
//     // 关闭窗口
//     drawerApi.close();
//     setLoading(false);
//   }
// }

async function createTenantWithAdminUser(values: any) {
  // 检查密码和确认密码是否一致
  if (values.password !== values.passwordConfirm) {
    notification.error({
      message: $t('ui.notification.password_mismatch'),
    });
    setLoading(false);
    return;
  }

  // 检查租户编码是否存在
  try {
    await apiClient.tenantService.TenantExists({
      code: values.code,
      name: values.name,
    });
  } catch {
    notification.error({
      message: $t('page.tenant.tenant_code_exists'),
    });
    setLoading(false);
    return;
  }

  // 检查用户名是否存在
  try {
    await apiClient.userService.UserExists({
      username: values.user.username,
    });
  } catch {
    notification.error({
      message: $t('page.tenant.notification.user_username_exists'),
    });
    setLoading(false);
    return;
  }

  try {
    await apiClient.tenantService.CreateTenantWithAdminUser({
      tenant: {
        name: values.name,
        code: values.code,
        type: values.type,
        auditStatus: values.auditStatus,
        status: values.status,
        remark: values.remark,
      },
      user: values.user,
      password: values.password,
    });

    notification.success({
      message: $t('ui.notification.create_success'),
    });
  } catch {
    notification.error({
      message: $t('ui.notification.create_failed'),
    });
  } finally {
    // 关闭窗口
    drawerApi.close();
    setLoading(false);
  }
}

async function updateTenant(values: any) {
  // 更新租户时仅提交后端 Update 回调实际可写的展示性字段。
  // code/type/auditStatus/status 后端经专门平台端点修改（见 tenant_repo.go Update 注释），
  // 普通更新不写——若混入 payload 会被静默丢弃，故排除以避免“已保存”假反馈。
  const tenantPayload = {
    name: values.name,
    remark: values.remark,
    logoUrl: values.logoUrl,
    industry: values.industry,
  };

  try {
    await apiClient.tenantService.Update({
      id: data.value.row.id,
      data: { ...tenantPayload },
      updateMask: makeUpdateMask(Object.keys(tenantPayload)),
    });

    notification.success({
      message: $t('ui.notification.update_success'),
    });
  } catch {
    notification.error({
      message: $t('ui.notification.update_failed'),
    });
  } finally {
    // 关闭窗口
    drawerApi.close();
    setLoading(false);
  }
}
</script>

<template>
  <Drawer :title="getTitle" class="w-full max-w-[800px]">
    <BaseForm />
  </Drawer>
</template>
