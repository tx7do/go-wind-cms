<script setup lang="ts">
import { Page } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { notification } from 'ant-design-vue';
import { ref } from 'vue';

import { useVbenForm } from '#/adapter/form';
import { changeMyPassword } from '#/api';

const submitLoading = ref(false);

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
      component: 'VbenInputPassword',
      fieldName: 'oldPassword',
      label: $t('page.user.form.oldPassword'),
      componentProps: {
        passwordStrength: true,
        placeholder: $t('ui.placeholder.input'),
      },
      rules: 'required',
    },
    {
      component: 'VbenInputPassword',
      fieldName: 'newPassword',
      label: $t('page.user.form.newPassword'),
      componentProps: {
        passwordStrength: true,
        placeholder: $t('ui.placeholder.input'),
      },
      rules: 'required',
    },
    {
      component: 'VbenInputPassword',
      fieldName: 'confirmPassword',
      label: $t('page.user.form.confirmPassword'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
      },
      rules: 'required',
    },
  ],
});

async function handleSubmit() {
  // 校验输入的数据
  const validate = await baseFormApi.validate();
  if (!validate.valid) {
    return;
  }

  setLoading(true);

  // 获取表单数据
  const values = await baseFormApi.getValues();

  if (values.newPassword !== values.confirmPassword) {
    notification.error({
      message: $t('ui.notification.password_mismatch'),
    });

    setLoading(false);
    return;
  }

  try {
    await changeMyPassword({
      oldPassword: values.oldPassword,
      newPassword: values.newPassword,
    });

    notification.success({
      message: $t('ui.notification.update_success'),
    });
  } catch {
    notification.error({
      message: $t('ui.notification.update_failed'),
    });
  } finally {
    setLoading(false);
  }
}

function setLoading(loading: boolean) {
  submitLoading.value = loading;
}
</script>

<template>
  <Page :title="$t('page.user.profile.tab.editPassword')">
    <BaseForm />
    <a-button
      type="primary"
      :loading="submitLoading"
      @click="handleSubmit"
    >
      {{ $t('page.user.button.updatePassword') }}
    </a-button>
  </Page>
</template>

<style scoped></style>
